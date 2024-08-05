package dbsource

import (
	"context"

	"github.com/VadimGossip/tj-drs-storage/internal/domain"
)

type Service interface {
	GetTaskRequests(ctx context.Context, limit int64) ([]domain.TaskRequest, error)
}

type service struct {
	repository Repository
}

var _ Service = (*service)(nil)

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) GetTaskRequests(ctx context.Context, limit int64) ([]domain.TaskRequest, error) {
	return s.repository.GetTaskRequests(ctx, limit)
}
