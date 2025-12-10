package userprojects

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for user-project data access operations.
type Repository interface {
	Create(ctx context.Context, up *postgres.UserProject) error
	Get(ctx context.Context, userID string, projectID int) (*postgres.UserProject, error)
	GetAll(ctx context.Context) ([]postgres.UserProject, error)
	Update(ctx context.Context, up *postgres.UserProject) error
	Delete(ctx context.Context, userID string, projectID int) error
}

type userProjectRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &userProjectRepository{db: db}
}

func (r *userProjectRepository) Create(ctx context.Context, up *postgres.UserProject) error {
	query := `
		INSERT INTO user_projects (user_id, project_id, permissions)
		VALUES ($1, $2, $3)
	`

	return r.db.QueryRowxContext(ctx, query,
		up.UserID,
		up.ProjectID,
		up.Permissions,
	).Scan(&up.DateCreated, &up.DateUpdated)
}

func (r *userProjectRepository) Get(ctx context.Context, userID string, projectID int) (*postgres.UserProject, error) {
	var up postgres.UserProject
	err := r.db.GetContext(ctx, &up, `
		SELECT * FROM user_projects WHERE user_id = $1 AND project_id = $2
	`, userID, projectID)
	if err != nil {
		return nil, err
	}
	return &up, nil
}

func (r *userProjectRepository) GetAll(ctx context.Context) ([]postgres.UserProject, error) {
	var ups []postgres.UserProject
	err := r.db.SelectContext(ctx, &ups, "SELECT * FROM user_projects")
	return ups, err
}

func (r *userProjectRepository) Update(ctx context.Context, up *postgres.UserProject) error {
	query := `
		UPDATE user_projects
		SET permissions = $1
		WHERE user_id = $2 AND project_id = $3
	`

	return r.db.QueryRowxContext(ctx, query,
		up.Permissions,
		up.UserID,
		up.ProjectID,
	).Scan(&up.DateUpdated)
}

func (r *userProjectRepository) Delete(ctx context.Context, userID string, projectID int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM user_projects WHERE user_id = $1 AND project_id = $2
	`, userID, projectID)
	return err
}
