package documentations

import (
	"context"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for documentation business logic operations.
type Service interface {
	Create(ctx context.Context, d *postgres.Documentation) error
	GetByID(ctx context.Context, id int) (*postgres.Documentation, error)
	GetAll(ctx context.Context) ([]postgres.Documentation, error)
	Update(ctx context.Context, d *postgres.Documentation) error
	Delete(ctx context.Context, id int) error
}

type documentationService struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &documentationService{repo: repo}
}

func (s *documentationService) Create(ctx context.Context, d *postgres.Documentation) error {
	return s.repo.Create(ctx, d)
}

func (s *documentationService) GetByID(ctx context.Context, id int) (*postgres.Documentation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *documentationService) GetAll(ctx context.Context) ([]postgres.Documentation, error) {
	return s.repo.GetAll(ctx)
}

func (s *documentationService) Update(ctx context.Context, d *postgres.Documentation) error {
	return s.repo.Update(ctx, d)
}

func (s *documentationService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
