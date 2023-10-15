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
	"github.com/rebus2015/praktikum-devops/internal/signer"
	pb "github.com/rebus2015/praktikum-devops/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GzipInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
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
		if key == nil {
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
		return handler(ctx, nextData)
	}
}

func HashInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if key == "" {
			return handler(ctx, req)
		}

		var err error
		data, ok := req.(*pb.AddMetricsRequest)
		if !ok {
			return nil, status.Errorf(codes.Canceled, "%v", "hash interceptor error: corrupted data")
		}
		reader := bytes.NewReader(data.Metrics)
		log.Println("Incoming request Updates, before decoder")

		var metrics []*model.Metrics
		bodyBytes, _ := io.ReadAll(reader)
		err = json.Unmarshal(bodyBytes, &metrics)
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
