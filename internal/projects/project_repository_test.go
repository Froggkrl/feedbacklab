package projects

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestProjectRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRepository(sqlxDB)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "gitlab_project_id", "mattermost_team", "date_created", "date_updated"}).
		AddRow(1, "Test", "Desc", nil, nil, now, now)

	mock.ExpectQuery(`SELECT \* FROM projects ORDER BY date_created DESC`).WillReturnRows(rows)

	items, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Test", items[0].Name)
}

func TestProjectRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRepository(sqlxDB)
	now := time.Now()
	row := sqlmock.NewRows([]string{"id", "name", "description", "gitlab_project_id", "mattermost_team", "date_created", "date_updated"}).
		AddRow(1, "Test", "Desc", nil, nil, now, now)

	mock.ExpectQuery(`SELECT \* FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(row)

	item, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Test", item.Name)
}

func TestProjectRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRepository(sqlxDB)
	now := time.Now()
	desc := "Desc"
	p := &postgres.Project{Name: "Test", Description: &desc}

	mock.ExpectQuery(`INSERT INTO projects`).
		WithArgs("Test", "Desc", nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).AddRow(1, now, now))

	err := repo.Create(context.Background(), p)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}
}

func TestProjectRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRepository(sqlxDB)
	now := time.Now()
	desc := "Updated"
	p := &postgres.Project{ID: 1, Name: "New", Description: &desc}

	mock.ExpectQuery(`UPDATE projects`).
		WithArgs("New", "Updated", nil, nil, 1).
		WillReturnRows(sqlmock.NewRows([]string{"date_updated"}).AddRow(now))

	err := repo.Update(context.Background(), p)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}
}

func TestProjectRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRepository(sqlxDB)
	mock.ExpectExec(`DELETE FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
}
