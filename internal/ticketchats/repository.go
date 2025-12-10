package ticketchats

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for ticket chat data access operations.
type Repository interface {
	Create(ctx context.Context, chat *postgres.TicketChat) error
	GetByID(ctx context.Context, id int) (*postgres.TicketChat, error)
	GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketChat, error)
	Update(ctx context.Context, chat *postgres.TicketChat) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, chat *postgres.TicketChat) error {
	query := `
		INSERT INTO ticket_chats (ticket_id, sender_id, sender_role, message, message_type, mattermost_message_id)
		VALUES (:ticket_id, :sender_id, :sender_role, :message, :message_type, :mattermost_message_id)
		RETURNING id, date_created, date_updated;
	`

	return r.db.QueryRowxContext(ctx, query, chat).
		StructScan(chat)
}

func (r *repository) GetByID(ctx context.Context, id int) (*postgres.TicketChat, error) {
	var chat postgres.TicketChat
	err := r.db.GetContext(ctx, &chat,
		`SELECT * FROM ticket_chats WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *repository) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketChat, error) {
	var chats []postgres.TicketChat
	err := r.db.SelectContext(ctx, &chats,
		`SELECT * FROM ticket_chats WHERE ticket_id = $1 ORDER BY date_created`,
		ticketID,
	)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *repository) Update(ctx context.Context, chat *postgres.TicketChat) error {
	query := `
		UPDATE ticket_chats
		SET message = :message,
		    message_type = :message_type,
		    mattermost_message_id = :mattermost_message_id
		WHERE id = :id
		RETURNING date_updated;
	`

	return r.db.QueryRowxContext(ctx, query, chat).
		Scan(&chat.DateUpdated)
}

func (r *repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM ticket_chats WHERE id = $1`,
		id,
	)
	return err
}
