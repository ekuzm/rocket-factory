package health

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *Server) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Server) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc.ServerStreamingServer[grpc_health_v1.HealthCheckResponse]) error {
	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func RegisterService(server *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(server, &Server{})
}
