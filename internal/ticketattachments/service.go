package ticketattachments

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for ticket attachment business logic operations.
type Service interface {
	Create(ctx context.Context, att *postgres.TicketAttachment) error
	GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error)
	GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error)
	Update(ctx context.Context, att *postgres.TicketAttachment) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, att *postgres.TicketAttachment) error {
	if att.FilePath == "" {
		return errors.New("file_path cannot be empty")
	}
	// Pass the context to the repository
	return s.repo.Create(ctx, att)
}

func (s *service) GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error) {
	return s.repo.GetByTicketID(ctx, ticketID)
}

func (s *service) Update(ctx context.Context, att *postgres.TicketAttachment) error {
	return s.repo.Update(ctx, att)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
