package tickets

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for ticket data access operations.
type Repository interface {
	Create(ctx context.Context, t *postgres.Ticket) error
	GetByID(ctx context.Context, id int) (*postgres.Ticket, error)
	GetAll(ctx context.Context) ([]postgres.Ticket, error)
	Update(ctx context.Context, t *postgres.Ticket) error
	Delete(ctx context.Context, id int) error
}

type ticketRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, t *postgres.Ticket) error {
	query := `
		INSERT INTO tickets (project_id, module_id, contract_id, created_by, assigned_to, title, message, status)
		VALUES (:project_id, :module_id, :contract_id, :created_by, :assigned_to, :title, :message, :status)
		RETURNING id, date_created, date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, t, t)
}

func (r *ticketRepository) GetByID(ctx context.Context, id int) (*postgres.Ticket, error) {
	var t postgres.Ticket
	err := r.db.GetContext(ctx, &t, "SELECT * FROM tickets WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ticketRepository) GetAll(ctx context.Context) ([]postgres.Ticket, error) {
	var tickets []postgres.Ticket
	err := r.db.SelectContext(ctx, &tickets, "SELECT * FROM tickets ORDER BY date_created DESC")
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) Update(ctx context.Context, t *postgres.Ticket) error {
	query := `
		UPDATE tickets
		SET title=:title, message=:message, status=:status, assigned_to=:assigned_to, module_id=:module_id
		WHERE id=:id
		RETURNING date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, &t.DateUpdated, t)
}

func (r *ticketRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tickets WHERE id=$1", id)
	return err
}
