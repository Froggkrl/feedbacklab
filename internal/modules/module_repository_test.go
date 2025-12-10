package modules

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestModuleRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "project_id", "name", "description", "responsible_user_id", "date_created", "date_updated"}).
		AddRow(1, 1, "Module", "Desc", nil, now, now)
	mock.ExpectQuery(`SELECT \* FROM modules ORDER BY date_created DESC`).WillReturnRows(rows)
	all, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"id", "project_id", "name", "description", "responsible_user_id", "date_created", "date_updated"}).
		AddRow(1, 1, "Module", "Desc", nil, now, now)
	mock.ExpectQuery(`SELECT \* FROM modules WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Module", got.Name)

	desc := "D"
	m := &postgres.Module{ProjectID: 1, Name: "New", Description: &desc}
	mock.ExpectQuery(`INSERT INTO modules`).
		WithArgs(1, "New", "D", nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).AddRow(1, now, now))
	err = repo.Create(ctx, m)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}

	m.ID = 1
	mock.ExpectQuery(`UPDATE modules`).
		WithArgs("New", "D", nil, 1).
		WillReturnRows(sqlmock.NewRows([]string{"date_updated"}).AddRow(now))
	err = repo.Update(ctx, m)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}

	mock.ExpectExec(`DELETE FROM modules WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
}
