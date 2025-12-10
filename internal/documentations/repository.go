package documentations

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for documentation data access operations.
type Repository interface {
	Create(ctx context.Context, d *postgres.Documentation) error
	GetByID(ctx context.Context, id int) (*postgres.Documentation, error)
	GetAll(ctx context.Context) ([]postgres.Documentation, error)
	Update(ctx context.Context, d *postgres.Documentation) error
	Delete(ctx context.Context, id int) error
}

type documentationRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &documentationRepository{db: db}
}

func (r *documentationRepository) Create(ctx context.Context, d *postgres.Documentation) error {
	query := `
		INSERT INTO documentations (project_id, file_path, version, uploaded_by)
		VALUES (:project_id, :file_path, :version, :uploaded_by)
		RETURNING id, date_created, date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, d, d)
}

func (r *documentationRepository) GetByID(ctx context.Context, id int) (*postgres.Documentation, error) {
	var d postgres.Documentation
	err := r.db.GetContext(ctx, &d, "SELECT * FROM documentations WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *documentationRepository) GetAll(ctx context.Context) ([]postgres.Documentation, error) {
	var docs []postgres.Documentation
	err := r.db.SelectContext(ctx, &docs, "SELECT * FROM documentations ORDER BY date_created DESC")
	return docs, err
}

func (r *documentationRepository) Update(ctx context.Context, d *postgres.Documentation) error {
	query := `
		UPDATE documentations
		SET file_path=:file_path, version=:version, uploaded_by=:uploaded_by
		WHERE id=:id
		RETURNING date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, &d.DateUpdated, d)
}

func (r *documentationRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM documentations WHERE id=$1", id)
	return err
}
