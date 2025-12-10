package tickets

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

func TestTicketRepository_Create_WithReturning(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()
	ticket := &postgres.Ticket{
		ProjectID:  1,
		ModuleID:   new(int),
		ContractID: 2,
		CreatedBy:  "user1",
		Title:      "Test Ticket",
		Message:    "Details",
		Status:     "open",
	}

	mock.ExpectPrepare(`INSERT INTO tickets \(project_id, module_id, contract_id, created_by, assigned_to, title, message, status\).*RETURNING.*`)

	rows := sqlmock.NewRows([]string{"id", "date_created", "date_updated"}).
		AddRow(123, now, now)

	mock.ExpectQuery(`INSERT INTO tickets \(project_id, module_id, contract_id, created_by, assigned_to, title, message, status\).*RETURNING.*`).
		WithArgs(ticket.ProjectID, ticket.ModuleID, ticket.ContractID, ticket.CreatedBy, nil, ticket.Title, ticket.Message, ticket.Status).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), ticket)

	assert.NoError(t, err)
	assert.Equal(t, 123, ticket.ID)
	assert.Equal(t, now, ticket.DateCreated)
	assert.Equal(t, now, ticket.DateUpdated)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketRepository_Update_WithReturning(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()
	ticket := &postgres.Ticket{
		ID:       1,
		Title:    "Updated",
		Message:  "New details",
		Status:   "in_progress",
		ModuleID: new(int),
	}

	mock.ExpectPrepare(`UPDATE tickets SET.*WHERE id=.*RETURNING date_updated`)

	rows := sqlmock.NewRows([]string{"date_updated"}).AddRow(now)

	mock.ExpectQuery(`UPDATE tickets SET.*WHERE id=.*RETURNING date_updated`).
		WillReturnRows(rows)

	err = repo.Update(context.Background(), ticket)

	assert.NoError(t, err)
	assert.Equal(t, now, ticket.DateUpdated)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketRepository_GetByID_HandlesNullableFields(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	now := time.Now().UTC()
	moduleID := 5

	rows := sqlmock.NewRows([]string{
		"id", "project_id", "module_id", "contract_id", "created_by",
		"assigned_to", "title", "message", "status", "gitlab_issue_url",
		"mattermost_thread_url", "date_created", "date_updated",
	}).AddRow(
		1, 2, moduleID, 3, "user1",
		"user2", "Title", "Message", "open", "http://gitlab.com/1",
		"http://mattermost.com/1", now, now,
	)

	mock.ExpectQuery(`SELECT \* FROM tickets WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(rows)

	ticket, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, ticket.ID)
	assert.Equal(t, moduleID, *ticket.ModuleID)
	assert.Equal(t, "user2", *ticket.AssignedTo)
	assert.Equal(t, "http://gitlab.com/1", *ticket.GitlabIssueURL)
	assert.NoError(t, mock.ExpectationsWereMet())
}
