package grpc

import (
	"github.com/soyacen/grocer/internal/layout/internal/grpc/grpc/repository"
	"github.com/soyacen/grocer/internal/layout/internal/grpc/grpc/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc",
	fx.Provide(repository.NewRepository, service.NewService),
)
