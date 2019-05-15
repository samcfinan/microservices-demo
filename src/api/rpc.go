package main

import (
	"context"

	pb "github.com/samcfinan/microservices-demo/src/api/genproto"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// gRPC Handler
func (fe *frontendServer) CheckName(ctx context.Context, nr *pb.NameRequest) (*pb.NameResponse, error) {
	resp, err := pb.NewNameServiceClient(fe.nameSvcConn).CheckName(ctx, nr)
	return resp, err
}

// HTTP Handler
func (fe *frontendServer) getNameLength(ctx context.Context, name string) (*pb.NameResponse, error) {
	resp, err := pb.NewNameServiceClient(fe.nameSvcConn).CheckName(ctx, &pb.NameRequest{Name: name})

	return resp, err
}

func (fe *frontendServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (fe *frontendServer) Watch(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
