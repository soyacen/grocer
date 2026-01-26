package cronjob

import "context"

type Service struct {
	repo *Repository
}

func (s *Service) Run(ctx context.Context) error {
	return nil
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
