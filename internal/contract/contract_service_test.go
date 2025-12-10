package contract

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestContractService_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRepository(sqlxDB)
	svc := NewService(repo)

	now := time.Now()
	startDate := now
	endDate := now.AddDate(1, 0, 0)

	rows := sqlmock.NewRows([]string{"id", "project_id", "client_name", "start_date", "end_date", "description"}).
		AddRow(1, 1, "Client", startDate, endDate, "Agreement")
	mock.ExpectQuery(`SELECT \* FROM contracts ORDER BY id`).WillReturnRows(rows)

	items, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Client", items[0].ClientName)

	row := sqlmock.NewRows([]string{"id", "project_id", "client_name", "start_date", "end_date", "description"}).
		AddRow(1, 1, "Client", startDate, endDate, "Agreement")
	mock.ExpectQuery(`SELECT \* FROM contracts WHERE id=\$1`).WithArgs(1).WillReturnRows(row)

	got, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Client", got.ClientName)

	contract := &Contract{
		ProjectID:   1,
		ClientName:  "New Client",
		StartDate:   startDate,
		EndDate:     endDate,
		Description: "New Agreement",
	}
	mock.ExpectQuery(`INSERT INTO contracts \(project_id, client_name, start_date, end_date, description\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(1, "New Client", startDate, endDate, "New Agreement").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	err = svc.Create(contract)
	assert.NoError(t, err)
	assert.Equal(t, 2, contract.ID)

	contract.ID = 2
	contract.ClientName = "Updated Client"
	mock.ExpectExec(`UPDATE contracts SET client_name=\$1, start_date=\$2, end_date=\$3, description=\$4 WHERE id=\$5`).
		WithArgs("Updated Client", startDate, endDate, "New Agreement", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = svc.Update(contract)
	assert.NoError(t, err)

	mock.ExpectExec(`DELETE FROM contracts WHERE id=\$1`).WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
	err = svc.Delete(2)
	assert.NoError(t, err)
}
