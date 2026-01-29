package service

import (
	"context"
	"fmt"

	"github.com/soyacen/grocer/internal/layout/api/grpc"
	"github.com/soyacen/grocer/internal/layout/internal/grpc/grpc/repository"
)

type Service struct {
	grpc.UnimplementedGreeterServer
	repo repository.Repository
}

func (s *Service) SayHello(ctx context.Context, req *grpc.HelloRequest) (*grpc.HelloReply, error) {
	return &grpc.HelloReply{
		Message: fmt.Sprintf("Hello %s!", req.GetName()),
	}, nil
}

func NewService(repo repository.Repository) grpc.GreeterServer {
	return &Service{repo: repo}
}
