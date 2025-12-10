package userprojects

import (
	"context"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for user-project business logic operations.
type Service interface {
	Create(ctx context.Context, up *postgres.UserProject) error
	Get(ctx context.Context, userID string, projectID int) (*postgres.UserProject, error)
	GetAll(ctx context.Context) ([]postgres.UserProject, error)
	Update(ctx context.Context, up *postgres.UserProject) error
	Delete(ctx context.Context, userID string, projectID int) error
}

type userProjectService struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &userProjectService{repo: repo}
}

func (s *userProjectService) Create(ctx context.Context, up *postgres.UserProject) error {
	return s.repo.Create(ctx, up)
}

func (s *userProjectService) Get(ctx context.Context, userID string, projectID int) (*postgres.UserProject, error) {
	return s.repo.Get(ctx, userID, projectID)
}

func (s *userProjectService) GetAll(ctx context.Context) ([]postgres.UserProject, error) {
	return s.repo.GetAll(ctx)
}

func (s *userProjectService) Update(ctx context.Context, up *postgres.UserProject) error {
	return s.repo.Update(ctx, up)
}

func (s *userProjectService) Delete(ctx context.Context, userID string, projectID int) error {
	return s.repo.Delete(ctx, userID, projectID)
}
