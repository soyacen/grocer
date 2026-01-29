package service

import (
	"context"
	"fmt"

	"github.com/soyacen/grocer/internal/layout/api/http"
	"github.com/soyacen/grocer/internal/layout/internal/http/http/repository"
)

type Service struct {
	repo repository.Repository
}

func (s *Service) Run(ctx context.Context) error {
	fmt.Println("implement logic here")
	return nil
}

func (s *Service) SayHello(ctx context.Context, req *http.HelloRequest) (*http.HelloReply, error) {
	return &http.HelloReply{
		Message: fmt.Sprintf("Hello %s!", req.GetName()),
	}, nil
}

func NewService(repo repository.Repository) http.GreeterService {
	return &Service{repo: repo}
}
