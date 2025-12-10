package messageattachments

import (
	"context"
	"innotech/internal/storage/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr(s string) *string { return &s }

func TestRepository_Create_WithReturning(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	att := &postgres.MessageAttachment{
		ChatID:     1,
		FilePath:   "/path/to/file",
		UploadedBy: "user1",
		FileType:   ptr("image/png"),
	}

	mock.ExpectPrepare(`INSERT INTO message_attachments.*RETURNING.*`)

	rows := sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).
		AddRow(42, now, now)

	mock.ExpectQuery(`INSERT INTO message_attachments.*RETURNING.*`).
		WithArgs(att.ChatID, att.FilePath, att.UploadedBy, att.FileType).
		WillReturnRows(rows)

	err = repo.Create(ctx, att)
	assert.NoError(t, err)
	assert.Equal(t, 42, att.ID)
	assert.Equal(t, now, att.DateCreated)
	assert.Equal(t, now, att.DateUpdated)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_WithReturning(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	att := &postgres.MessageAttachment{
		ID:       5,
		FilePath: "/new/path",
		FileType: ptr("application/pdf"),
	}

	mock.ExpectPrepare(`UPDATE message_attachments SET.*WHERE id = \?.*RETURNING date_updated`)

	rows := sqlmock.NewRows([]string{"date_updated"}).AddRow(now)

	mock.ExpectQuery(`UPDATE message_attachments SET.*WHERE id = \?.*RETURNING date_updated`).
		WithArgs(att.FilePath, att.FileType, att.ID).
		WillReturnRows(rows)

	err = repo.Update(ctx, att)
	assert.NoError(t, err)
	assert.Equal(t, now, att.DateUpdated)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByChatID_WithOrder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{
		"id", "chat_id", "file_path", "uploaded_by", "file_type",
		"date_created", "date_updated",
	}).AddRow(
		1, 42, "/path1", "user1", ptr("image/png"),
		now, now,
	).AddRow(
		2, 42, "/path2", "user2", ptr("text/plain"),
		now.Add(time.Hour), now.Add(time.Hour),
	)

	mock.ExpectQuery(`SELECT \* FROM message_attachments WHERE chat_id = \$1 ORDER BY date_created`).
		WithArgs(42).
		WillReturnRows(rows)

	list, err := repo.GetByChatID(ctx, 42)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "image/png", *list[0].FileType)
	assert.Equal(t, "text/plain", *list[1].FileType)
	assert.NoError(t, mock.ExpectationsWereMet())
}
