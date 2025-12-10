package ticketattachments

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
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()
	att := &postgres.TicketAttachment{
		TicketID:    1,
		FilePath:    "/path/to/file",
		UploadedBy:  "user1",
		FileType:    ptr("image/png"),
		Description: ptr("Screenshot"),
	}

	mock.ExpectPrepare(`INSERT INTO ticket_attachments.*RETURNING.*`)

	rows := sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).
		AddRow(42, now, now)

	mock.ExpectQuery(`INSERT INTO ticket_attachments.*RETURNING.*`).
		WithArgs(att.TicketID, att.FilePath, att.UploadedBy, att.FileType, att.Description).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), att)

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

	now := time.Now().UTC()
	att := &postgres.TicketAttachment{
		ID:          5,
		FilePath:    "/new/path",
		FileType:    ptr("application/pdf"),
		Description: ptr("Updated doc"),
	}

	mock.ExpectPrepare(`UPDATE ticket_attachments SET.*WHERE id = \?.*RETURNING date_updated`)

	rows := sqlmock.NewRows([]string{"date_updated"}).AddRow(now)

	mock.ExpectQuery(`UPDATE ticket_attachments SET.*WHERE id = \?.*RETURNING date_updated`).
		WithArgs(att.FilePath, att.FileType, att.Description, att.ID).
		WillReturnRows(rows)

	err = repo.Update(context.Background(), att)

	assert.NoError(t, err)
	assert.Equal(t, now, att.DateUpdated)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByTicketID_WithOrder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{
		"id", "ticket_id", "file_path", "uploaded_by", "file_type",
		"description", "date_created", "date_updated",
	}).AddRow(
		1, 42, "/path1", "user1", ptr("image/png"),
		ptr("First file"), now, now,
	).AddRow(
		2, 42, "/path2", "user2", ptr("text/plain"),
		ptr("Second file"), now.Add(time.Hour), now.Add(time.Hour),
	)

	mock.ExpectQuery(`SELECT \* FROM ticket_attachments WHERE ticket_id = \$1 ORDER BY date_created`).
		WithArgs(42).
		WillReturnRows(rows)

	list, err := repo.GetByTicketID(context.Background(), 42)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "First file", *list[0].Description)
	assert.Equal(t, "Second file", *list[1].Description)
	assert.NoError(t, mock.ExpectationsWereMet())
}
