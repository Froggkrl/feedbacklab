package messageattachments

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for message attachment data access operations.
type Repository interface {
	Create(ctx context.Context, att *postgres.MessageAttachment) error
	GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error)
	GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error)
	Update(ctx context.Context, att *postgres.MessageAttachment) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, att *postgres.MessageAttachment) error {
	query := `
		INSERT INTO message_attachments (chat_id, file_path, uploaded_by, file_type)
		VALUES (:chat_id, :file_path, :uploaded_by, :file_type)
		RETURNING id, date_created, date_updated;
	`

	return r.db.QueryRowxContext(ctx, query, att).
		StructScan(att)
}

func (r *repository) GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error) {
	var att postgres.MessageAttachment
	err := r.db.GetContext(ctx, &att,
		`SELECT * FROM message_attachments WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error) {
	var list []postgres.MessageAttachment
	err := r.db.SelectContext(ctx, &list,
		`SELECT * FROM message_attachments WHERE chat_id = $1 ORDER BY date_created`,
		chatID,
	)
	return list, err
}

func (r *repository) Update(ctx context.Context, att *postgres.MessageAttachment) error {
	query := `
		UPDATE message_attachments
		SET file_path = :file_path,
			file_type = :file_type
		WHERE id = :id
		RETURNING date_updated;
	`

	return r.db.QueryRowxContext(ctx, query, att).
		Scan(&att.DateUpdated)
}

func (r *repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM message_attachments WHERE id = $1`,
		id,
	)
	return err
}
