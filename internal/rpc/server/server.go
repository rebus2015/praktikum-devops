package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	pb "github.com/rebus2015/praktikum-devops/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rebus2015/praktikum-devops/internal/config"
	"github.com/rebus2015/praktikum-devops/internal/model"
	"github.com/rebus2015/praktikum-devops/internal/rpc/interceptors"
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

func (s *MetricsRPCServer) AddMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var response pb.AddMetricsResponse
	log.Println("Incoming request Updates, before decoder")

	var metrics []*model.Metrics
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

func (s *MetricsRPCServer) Run() error {
	host, _, err := net.SplitHostPort(s.cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("grpc server err(SplitHostPort): %w", err)
	}
	listen, err := net.Listen("tcp", host+s.cfg.PortRPC)
	if err != nil {
		return fmt.Errorf("start RPC server error: %w", err)
	}
	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.GzipInterceptor,
			interceptors.RsaInterceptor(s.cfg.CryptoKey),
			interceptors.HashInterceptor(s.cfg.Key),
		))
	// создаём gRPC-сервер без зарегистрированной службы
	s.srv = grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServer(s.srv, &MetricsRPCServer{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.srv.Serve(listen); err != nil {
		log.Fatal(err)
		return fmt.Errorf("run server err: %w", err)
	}
	return nil
}
