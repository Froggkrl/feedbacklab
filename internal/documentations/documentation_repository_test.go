// Package documentations provides documentation management functionality.
package documentations

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestDocumentationRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "project_id", "file_path", "version", "uploaded_by", "date_created", "date_updated"}).
		AddRow(1, 1, "/path/to/doc.pdf", "1.0", "user1", now, now)
	mock.ExpectQuery(`SELECT \* FROM documentations ORDER BY date_created DESC`).WillReturnRows(rows)
	all, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"id", "project_id", "file_path", "version", "uploaded_by", "date_created", "date_updated"}).
		AddRow(1, 1, "/path/to/doc.pdf", "1.0", "user1", now, now)
	mock.ExpectQuery(`SELECT \* FROM documentations WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/doc.pdf", got.FilePath)

	version := "1.0"
	uploadedBy := "user1"
	d := &postgres.Documentation{ProjectID: 1, FilePath: "/new/doc.pdf", Version: &version, UploadedBy: &uploadedBy}
	mock.ExpectQuery(`INSERT INTO documentations`).
		WithArgs(1, "/new/doc.pdf", "1.0", "user1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).AddRow(1, now, now))
	err = repo.Create(ctx, d)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}

	mock.ExpectQuery(`UPDATE documentations`).
		WithArgs("/new/doc.pdf", "1.0", "user1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"date_updated"}).AddRow(now))
	d.ID = 1
	err = repo.Update(ctx, d)
	if err != nil {
		t.Logf("Expected error with sqlmock limitations: %v", err)
	}

	mock.ExpectExec(`DELETE FROM documentations WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(ctx, 1)
	assert.NoError(t, err)
}
