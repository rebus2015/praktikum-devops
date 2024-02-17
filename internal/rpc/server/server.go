package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/rpc/interceptors"
	pb "github.com/rebus2015/praktikum-devops/internal/rpc/proto"
	"github.com/rebus2015/praktikum-devops/internal/storage"
	"github.com/rebus2015/praktikum-devops/internal/storage/dbstorage"
)

type MetricsRPCServer struct {
	srv *grpc.Server
	pb.UnimplementedMetricsServer
	metricStorage  storage.Repository
	postgreStorage dbstorage.SQLStorage
	cfg            config.Config
}

func NewRPCServer(storage storage.Repository,
	pgsStorage dbstorage.SQLStorage,
	conf config.Config) *MetricsRPCServer {
	return &MetricsRPCServer{
		metricStorage:  storage,
		postgreStorage: pgsStorage,
		cfg:            conf,
	}
}

func (s *MetricsRPCServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse
	log.Println("Incoming request Ping")
	// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
	if s.postgreStorage == nil {
		response.Status = 500
		response.Error = status.Error(codes.Internal,
			"Failed to ping database: nil reference exception: postgreStorage udefined").Error()
		return &response, fmt.Errorf("nil reference exception: postgreStorage udefined ")
	}
	if err := s.postgreStorage.Ping(ctx); err != nil {
		log.Printf("Cannot ping database because %s", err)
		response.Status = 500
		response.Error = status.Errorf(codes.Internal, "Failed to ping database because %s", err).Error()
		return &response, fmt.Errorf("failed to Decode incoming metricList %w", err)
	}
	response.Status = 200
	return &response, nil
}

func (s *MetricsRPCServer) AddMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse
	log.Println("Incoming request Updates, before decoder")

	var metrics = make([]*model.Metrics, 0)
	err := json.Unmarshal(in.Metrics, &metrics)
	if err != nil {
		log.Printf("Failed to Decode incoming metricList %v, error: %v", string(in.Metrics), err)
		response.Error = status.Errorf(codes.InvalidArgument, "Failed to Decode incoming metricList %v", err).Error()
		return &response, fmt.Errorf("failed to Decode incoming metricList %w", err)
	}

	err = s.metricStorage.AddMetrics(metrics)
	if err != nil {
		log.Printf("Error: [UpdateJSONMultipleMetricHandlerFunc] Add multiple metrics error: %v", err)
		response.Error = status.Errorf(codes.Internal, "Add multiple metrics error: %v", err).Error()
		return &response, fmt.Errorf("add multiple metrics error: %w", err)
	}

	return &response, nil
}

func (s *MetricsRPCServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	log.Println("Incoming request Updates, before decoder")

	mtype := in.Mname
	name := in.Mtype
	val := in.Val
	var err error
	switch mtype {
	case "gauge":
		_, err = s.metricStorage.AddGauge(name, val)
	case "counter":
		_, err = s.metricStorage.AddCounter(name, val)
	default:
		{
			log.Printf("Error: [UpdateMetricHandlerFunc] Update metric error: %v", err)
			response.Error = status.Errorf(codes.Internal, "Update single metric error: %v", err).Error()
			return &response, fmt.Errorf("update single metric error: %w", err)
		}
	}
	if err != nil {
		log.Printf("Error: [UpdateMetricHandlerFunc] Update metric error: %v", err)
		response.Error = status.Errorf(codes.Internal, "Update single metric error: %v", err).Error()
		return &response, fmt.Errorf("update single metric error: %w", err)
	}

	return &response, nil
}

func (s *MetricsRPCServer) Run() error {
	listen, err := net.Listen("tcp", s.cfg.RPCServerAddress)
	if err != nil {
		return fmt.Errorf("start RPC server error: %w", err)
	}

	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.SubnetCheckInterceptor(s.cfg),
			interceptors.GzipInterceptor,
			interceptors.RsaInterceptor(s.cfg.CryptoKey),
			interceptors.HashInterceptor(s.cfg.Key),
		))

	// регистрируем сервис
	pb.RegisterMetricsServer(s.srv, &MetricsRPCServer{
		metricStorage:  s.metricStorage,
		postgreStorage: s.postgreStorage,
		cfg:            s.cfg,
	})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.srv.Serve(listen); err != nil {
		log.Fatal(err)
		return fmt.Errorf("run server err: %w", err)
	}
	return nil
}

func (s *MetricsRPCServer) Shutdown() {
	s.srv.GracefulStop()
}
