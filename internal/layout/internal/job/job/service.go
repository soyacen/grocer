package job

import (
	"context"
	"log/slog"
)

type Service struct {
	repo *Repository
}

func (s *Service) Run(ctx context.Context) error {
	slog.Info("implement logic here")
	return nil
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
