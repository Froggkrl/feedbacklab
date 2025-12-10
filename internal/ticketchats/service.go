package ticketchats

import (
	"context"
	"errors"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for ticket chat business logic operations.
type Service interface {
	Create(ctx context.Context, chat *postgres.TicketChat) error
	GetByID(ctx context.Context, id int) (*postgres.TicketChat, error)
	GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketChat, error)
	Update(ctx context.Context, chat *postgres.TicketChat) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, chat *postgres.TicketChat) error {
	if chat.Message == "" {
		return errors.New("message cannot be empty")
	}
	return s.repo.Create(ctx, chat)
}

func (s *service) GetByID(ctx context.Context, id int) (*postgres.TicketChat, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketChat, error) {
	return s.repo.GetByTicketID(ctx, ticketID)
}

func (s *service) Update(ctx context.Context, chat *postgres.TicketChat) error {
	return s.repo.Update(ctx, chat)
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
