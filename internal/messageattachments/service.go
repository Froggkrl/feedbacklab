package messageattachments

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for message attachment business logic operations.
type Service interface {
	Create(ctx context.Context, att *postgres.MessageAttachment) error
	GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error)
	GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error)
	Update(ctx context.Context, att *postgres.MessageAttachment) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, att *postgres.MessageAttachment) error {
	if att.FilePath == "" {
		return errors.New("file_path cannot be empty")
	}
	if att.ChatID == 0 {
		return errors.New("chat_id is required")
	}
	return s.repo.Create(ctx, att)
}

func (s *service) GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error) {
	return s.repo.GetByChatID(ctx, chatID)
}

func (s *service) Update(ctx context.Context, att *postgres.MessageAttachment) error {
	return s.repo.Update(ctx, att)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
