package http

import (
	"github.com/soyacen/grocer/internal/layout/internal/http/http/repository"
	"github.com/soyacen/grocer/internal/layout/internal/http/http/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	fx.Provide(repository.NewRepository, service.NewService),
)
