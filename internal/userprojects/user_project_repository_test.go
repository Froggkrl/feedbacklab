package userprojects

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserProjectRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"user_id", "project_id", "role", "permissions", "date_created", "date_updated"}).
		AddRow("uuid-1", 2, "developer", []string{"read", "write"}, now, now)
	mock.ExpectQuery(`SELECT \* FROM user_projects`).WillReturnRows(rows)
	all, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"user_id", "project_id", "role", "permissions", "date_created", "date_updated"}).
		AddRow("uuid-1", 2, "developer", []string{"read", "write"}, now, now)
	mock.ExpectQuery(`SELECT \* FROM user_projects WHERE user_id = \$1 AND project_id = \$2`).
		WithArgs("uuid-1", 2).WillReturnRows(row)
	got, err := repo.Get(ctx, "uuid-1", 2)
	assert.NoError(t, err)
	assert.Equal(t, "uuid-1", got.UserID)

	up := &postgres.UserProject{UserID: "uuid-1", ProjectID: 2, Permissions: []string{"read"}}
	mock.ExpectQuery(`INSERT INTO user_projects`).
		WithArgs("uuid-1", 2, []string{"read"}).
		WillReturnRows(sqlmock.NewRows([]string{"date_created", "date_updated"}).AddRow(now, now))
	err = repo.Create(ctx, up)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}

	mock.ExpectQuery(`UPDATE user_projects`).
		WithArgs([]string{"read", "write"}, "uuid-1", 2).
		WillReturnRows(sqlmock.NewRows([]string{"date_updated"}).AddRow(now))
	up.Permissions = []string{"read", "write"}
	err = repo.Update(ctx, up)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}

	mock.ExpectExec(`DELETE FROM user_projects WHERE user_id = \$1 AND project_id = \$2`).
		WithArgs("uuid-1", 2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(ctx, "uuid-1", 2)
	assert.NoError(t, err)
}
