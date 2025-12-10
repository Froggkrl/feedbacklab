package ticketchats

import (
	"context"
	"database/sql"
	"innotech/internal/storage/postgres"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create_WithReturning(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func(mockDB *sql.DB) {
		err := mockDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(mockDB)

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()
	chat := &postgres.TicketChat{
		TicketID:    1,
		SenderID:    "user1",
		SenderRole:  "customer",
		Message:     "test message",
		MessageType: "text",
	}

	rows := sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).
		AddRow(42, now, now)

	mock.ExpectQuery(`INSERT INTO ticket_chats.*RETURNING.*`).
		WillReturnRows(rows)

	ctx := context.Background()
	err = repo.Create(ctx, chat)

	assert.NoError(t, err)
	assert.Equal(t, 42, chat.ID)
	assert.Equal(t, now, chat.DateCreated)
	assert.Equal(t, now, chat.DateUpdated)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_WithReturning(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()
	chat := &postgres.TicketChat{
		ID:          5,
		Message:     "updated",
		MessageType: "edited",
	}

	rows := sqlmock.NewRows([]string{"date_updated"}).AddRow(now)

	mock.ExpectQuery(`UPDATE ticket_chats SET.*WHERE id = .*RETURNING date_updated`).
		WithArgs(chat.Message, chat.MessageType, nil, chat.ID).
		WillReturnRows(rows)

	ctx := context.Background()
	err = repo.Update(ctx, chat)

	assert.NoError(t, err)
	assert.Equal(t, now, chat.DateUpdated)
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
		"id", "ticket_id", "sender_id", "sender_role", "message",
		"message_type", "mattermost_message_id", "date_created", "date_updated",
	}).AddRow(
		1, 42, "user1", "customer", "msg1",
		"text", nil, now, now,
	).AddRow(
		2, 42, "user2", "support", "msg2",
		"text", "mm-123", now.Add(time.Hour), now.Add(time.Hour),
	)

	mock.ExpectQuery(`SELECT \* FROM ticket_chats WHERE ticket_id = \$1 ORDER BY date_created ASC`).
		WithArgs(42).
		WillReturnRows(rows)

	ctx := context.Background()
	list, err := repo.GetByTicketID(ctx, 42)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "mm-123", *list[1].MattermostMessageID)
	assert.Nil(t, list[0].MattermostMessageID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{
		"id", "ticket_id", "sender_id", "sender_role", "message",
		"message_type", "mattermost_message_id", "date_created", "date_updated",
	}).AddRow(
		1, 42, "user1", "customer", "test message",
		"text", "mm-123", now, now,
	)

	mock.ExpectQuery(`SELECT \* FROM ticket_chats WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.Background()
	chat, err := repo.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, chat.ID)
	assert.Equal(t, "test message", chat.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	mock.ExpectExec(`DELETE FROM ticket_chats WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	err = repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
