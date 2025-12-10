package projects

import (
	"context"
	"innotech/internal/storage/postgres"
	"innotech/pkg/logger"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for project data access operations.
type Repository interface {
	Create(ctx context.Context, p *postgres.Project) error
	GetByID(ctx context.Context, id int) (*postgres.Project, error)
	GetAll(ctx context.Context) ([]postgres.Project, error)
	Update(ctx context.Context, p *postgres.Project) error
	Delete(ctx context.Context, id int) error
}

type projectRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	logger.Info("project repository initialized")
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, p *postgres.Project) error {
	logger.Debug("repo: create project start",
		"name", p.Name,
		"gitlab_project_id", p.GitlabProjectID,
	)

	query := `
		INSERT INTO projects (name, description, gitlab_project_id, mattermost_team)
		VALUES (:name, :description, :gitlab_project_id, :mattermost_team)
		RETURNING id, date_created, date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		logger.Error("repo: prepare named context failed",
			"error", err.Error(),
			"name", p.Name,
		)
		return err
	}
	if err := stmt.GetContext(ctx, p, p); err != nil {
		logger.Error("repo: insert project failed",
			"error", err.Error(),
			"name", p.Name,
		)
		return err
	}

	logger.Info("repo: project created",
		"id", p.ID,
		"name", p.Name,
	)
	return nil
}

func (r *projectRepository) GetByID(ctx context.Context, id int) (*postgres.Project, error) {
	logger.Debug("repo: get project by id", "id", id)

	var p postgres.Project
	err := r.db.GetContext(ctx, &p, "SELECT * FROM projects WHERE id=$1", id)
	if err != nil {
		logger.Error("repo: get by id failed",
			"id", id,
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug("repo: project retrieved successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return &p, nil
}

func (r *projectRepository) GetAll(ctx context.Context) ([]postgres.Project, error) {
	logger.Debug("repo: get all projects")

	var ps []postgres.Project
	err := r.db.SelectContext(ctx, &ps, "SELECT * FROM projects ORDER BY date_created DESC")
	if err != nil {
		logger.Error("repo: select all failed",
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Info("repo: projects list retrieved",
		"count", len(ps),
	)
	return ps, nil
}

func (r *projectRepository) Update(ctx context.Context, p *postgres.Project) error {
	logger.Debug("repo: update project",
		"id", p.ID,
		"name", p.Name,
	)

	query := `
		UPDATE projects
		SET name=:name, description=:description, gitlab_project_id=:gitlab_project_id, mattermost_team=:mattermost_team
		WHERE id=:id
		RETURNING date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		logger.Error("repo: prepare named context failed",
			"error", err.Error(),
			"id", p.ID,
		)
		return err
	}
	if err := stmt.GetContext(ctx, &p.DateUpdated, p); err != nil {
		logger.Error("repo: update failed",
			"id", p.ID,
			"error", err.Error(),
		)
		return err
	}

	logger.Info("repo: project updated",
		"id", p.ID,
		"name", p.Name,
	)
	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id int) error {
	logger.Warn("repo: delete project", "id", id)

	_, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id=$1", id)
	if err != nil {
		logger.Error("repo: delete failed",
			"id", id,
			"error", err.Error(),
		)
		return err
	}

	logger.Info("repo: project deleted", "id", id)
	return nil
}
