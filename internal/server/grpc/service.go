package grpc

import (
	"golang.org/x/net/context"
	ggrpc "google.golang.org/grpc"

	"github.com/smakimka/mtrcscollector/internal/model"
	"github.com/smakimka/mtrcscollector/internal/server/config"
	"github.com/smakimka/mtrcscollector/internal/server/grpc/interceptors"
	"github.com/smakimka/mtrcscollector/internal/storage"
	pb "github.com/smakimka/mtrcscollector/protobuf/server"
)

func NewServer(cfg *config.Config, storage storage.Storage) *ggrpc.Server {
	inters := []ggrpc.UnaryServerInterceptor{}

	if cfg.TrustedSubnet != nil {
		subnetInterseptor := interceptors.NewSubnetInterseptor(cfg.TrustedSubnet)
		inters = append(inters, subnetInterseptor.AllowTrusted)
	}

	s := ggrpc.NewServer(ggrpc.ChainUnaryInterceptor(inters...))
	service := &Service{s: storage}

	pb.RegisterMetricsCollectorServer(s, service)

	return s
}

type Service struct {
	pb.UnimplementedMetricsCollectorServer
	s storage.Storage
}

func (s *Service) Update(ctx context.Context, in *pb.UpdateMetrics) (*pb.Response, error) {
	var response pb.Response

	data := model.MetricsData{}
	for _, metric := range in.Metrics {
		data = append(data, model.MetricData{
			Delta: &metric.Delta,
			Value: &metric.Value,
			Name:  metric.Name,
			Kind:  metric.Kind,
		})
	}

	err := s.s.UpdateMetrics(ctx, data)
	if err != nil {
		response.Ok = false
		response.Detail = err.Error()
		return &response, nil
	}

	response.Ok = true
	return &response, nil
}
