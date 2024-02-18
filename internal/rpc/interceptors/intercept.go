package interceptors

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rebus2015/praktikum-devops/internal/model"
	pb "github.com/rebus2015/praktikum-devops/internal/rpc/proto"
	"github.com/rebus2015/praktikum-devops/internal/signer"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const methodPing string = "/proto.Metrics/Ping"

func GzipInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info.FullMethod == methodPing {
		return handler(ctx, req)
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("gzip")
		if len(values) == 0 {
			return handler(ctx, req)
		}
	}
	data, ok := req.(*pb.AddMetricsRequest)
	if !ok {
		return nil, status.Errorf(codes.Canceled, "%v", "gzip interceptor error: corrupted data")
	}

	gz, err := gzip.NewReader(bytes.NewReader(data.Metrics))
	if err != nil {
		log.Printf("Failed to create gzip reader: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "gzip interceptor error: %v", err)
	}
	defer func() {
		if err := gz.Close(); err != nil {
			log.Errorf("error occured when closing gzip: %v", err)
		}
	}()
	data.Metrics, err = io.ReadAll(gz)
	if err != nil {
		log.Printf("Failed to read bytes from gzip reader: %v", err.Error())
		return nil, status.Errorf(codes.Internal,
			"gzip interceptor error, read all from gz: %v", err)
	}

	return handler(ctx, data)
}

func RsaInterceptor(key *rsa.PrivateKey) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if key == nil || info.FullMethod == methodPing {
			return handler(ctx, req)
		}
		data, ok := req.(*pb.AddMetricsRequest)
		if !ok {
			return nil, status.Errorf(codes.Canceled, "%v", "gzip interceptor error: corrupted data")
		}
		nextData, err := signer.DecryptMessage(key, data.Metrics)
		if err != nil {
			log.Printf("Failed to create gzip reader: %v", err.Error())
			return nil, status.Errorf(codes.Internal, "gzip interceptor error: %v", err)
		}
		data.Metrics = nextData
		return handler(ctx, data)
	}
}

func HashInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if key == "" || info.FullMethod == methodPing {
			return handler(ctx, req)
		}

		var err error
		data, ok := req.(*pb.AddMetricsRequest)
		if !ok {
			return nil, status.Errorf(codes.Canceled, "%v", "hash interceptor error: corrupted data")
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("single")
			if len(values) == 0 {
				return handler(ctx, req)
			}
		}

		reader := bytes.NewReader(data.Metrics)
		log.Println("Incoming request Updates, before decoder")

		bodyBytes, _ := io.ReadAll(reader)
		metrics, err := getMetrics(bodyBytes)
		if err != nil {
			log.Printf("Failed to Decode incoming metricList %v, error: %v", string(bodyBytes), err)
			return nil, status.Errorf(codes.InvalidArgument, "Failed to Decode incoming metricList %v", err)
		}
		log.Printf("Try to update metrics: %v", metrics)
		for i := range metrics {
			if key != "" {
				pass, err := checkMetric(metrics[i], key)
				if err != nil || !pass {
					log.Printf("check Metrics error %v, error: %v", string(bodyBytes), err)
					return nil, status.Errorf(codes.InvalidArgument, "Failed to Decode incoming metricList %v", err)
				}
			}
		}
		return handler(ctx, data)
	}
}

func getMetrics(body []byte) ([]*model.Metrics, error) {
	var metrics []*model.Metrics
	x := bytes.TrimLeft(body, " \t\r\n")

	isArray := len(x) > 0 && x[0] == '['
	isObject := len(x) > 0 && x[0] == '{'

	switch {
	case isArray:
		err := json.Unmarshal(body, &metrics)
		if err != nil {
			return nil, fmt.Errorf("unmarshall error:%w", err)
		}
	case isObject:
		var metric *model.Metrics
		err := json.Unmarshal(body, metric)
		if err != nil {
			return nil, fmt.Errorf("unmarshall error:%w", err)
		}
		metrics = append(metrics, metric)
	default:
		return nil, fmt.Errorf("unmarshall error: couldn't define the type of structure")
	}

	return metrics, nil
}

func SubnetCheckInterceptor(s subnet) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if s == nil {
			return handler(ctx, req)
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("X-Real-IP")
			if len(values) == 0 {
				return nil, status.Errorf(codes.Code(403), "%v", "Request is from Untrusted subnet")
			}
			for _, addr := range values {
				if s.CheckIP(addr) {
					return handler(ctx, req)
				}
			}
		}
		return nil, status.Errorf(codes.Code(403), "%v", "Request is from Untrusted subnet")
	}
}

type subnet interface {
	CheckIP(ipAddr string) bool
}

// checkMetric внутренняя функция проверки целостности метрики.
func checkMetric(metric *model.Metrics, key string) (bool, error) {
	if metric.ID == "" {
		return false, fmt.Errorf("metric.ID is empty /n Body: %v", metric)
	}
	if metric.MType == "" {
		return false, fmt.Errorf("metric.MType is empty /n Body: %v", metric)
	}
	if metric.Hash != "" {
		hashObject := signer.NewHashObject(key)
		passed, err := hashObject.Verify(metric)
		if err != nil {
			log.Printf(
				"Incoming Metric verification error: \nBody: %v, \n error: %v",
				metric,
				err)
			return false, fmt.Errorf("incoming Metric verification error: \nBody: %v, \n error: %w",
				metric,
				err)
		}
		if !passed {
			log.Printf(
				"Error: Incoming Metric could not pass signature verification: \nBody: %v",
				metric)

			return false, fmt.Errorf("error: Incoming Metric could not pass signature verification: \nBody: %v",
				metric)
		}
	}
	return true, nil
}
