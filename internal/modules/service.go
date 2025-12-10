package modules

import (
	"context"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for module business logic operations.
type Service interface {
	Create(ctx context.Context, m *postgres.Module) error
	GetByID(ctx context.Context, id int) (*postgres.Module, error)
	GetAll(ctx context.Context) ([]postgres.Module, error)
	Update(ctx context.Context, m *postgres.Module) error
	Delete(ctx context.Context, id int) error
}

type moduleService struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &moduleService{repo: repo}
}

func (s *moduleService) Create(ctx context.Context, m *postgres.Module) error {
	return s.repo.Create(ctx, m)
}

func (s *moduleService) GetByID(ctx context.Context, id int) (*postgres.Module, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *moduleService) GetAll(ctx context.Context) ([]postgres.Module, error) {
	return s.repo.GetAll(ctx)
}

func (s *moduleService) Update(ctx context.Context, m *postgres.Module) error {
	return s.repo.Update(ctx, m)
}

func (s *moduleService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
