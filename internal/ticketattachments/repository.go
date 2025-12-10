package ticketattachments

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for ticket attachment data access operations.
type Repository interface {
	Create(ctx context.Context, att *postgres.TicketAttachment) error
	GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error)
	GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error)
	Update(ctx context.Context, att *postgres.TicketAttachment) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, att *postgres.TicketAttachment) error {
	query := `
		INSERT INTO ticket_attachments (ticket_id, file_path, uploaded_by, file_type, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, date_created, date_updated
	`

	return r.db.QueryRowxContext(ctx, query,
		att.TicketID,
		att.FilePath,
		att.UploadedBy,
		att.FileType,
		att.Description,
	).StructScan(att)
}

func (r *repository) GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error) {
	var att postgres.TicketAttachment
	err := r.db.GetContext(ctx, &att, `SELECT * FROM ticket_attachments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error) {
	var list []postgres.TicketAttachment
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM ticket_attachments WHERE ticket_id = $1 ORDER BY date_created`, ticketID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) Update(ctx context.Context, att *postgres.TicketAttachment) error {
	query := `
		UPDATE ticket_attachments
		SET file_path = $1,
		    file_type = $2,
		    description = $3
		WHERE id = $4
		RETURNING date_updated
	`

	return r.db.QueryRowxContext(ctx, query,
		att.FilePath,
		att.FileType,
		att.Description,
		att.ID,
	).Scan(&att.DateUpdated)
}

func (r *repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ticket_attachments WHERE id = $1`, id)
	return err
}
