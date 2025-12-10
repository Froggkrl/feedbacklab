// Package contract provides contract management functionality.
package contract

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestContractRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "project_id", "client_name", "start_date", "end_date", "description"}).
		AddRow(1, 1, "Client", time.Now(), time.Now().AddDate(1, 0, 0), "Desc")
	mock.ExpectQuery(`SELECT \* FROM contracts ORDER BY id`).WillReturnRows(rows)
	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"id", "project_id", "client_name", "start_date", "end_date", "description"}).
		AddRow(1, 1, "Client", time.Now(), time.Now().AddDate(1, 0, 0), "Desc")
	mock.ExpectQuery(`SELECT \* FROM contracts WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Client", got.ClientName)

	startDate := time.Now()
	endDate := time.Now().AddDate(1, 0, 0)
	c := &Contract{ProjectID: 1, ClientName: "New Client", StartDate: startDate, EndDate: endDate, Description: "D"}
	mock.ExpectQuery(`INSERT INTO contracts \(project_id, client_name, start_date, end_date, description\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(1, "New Client", startDate, endDate, "D").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	err = repo.Create(c)
	assert.NoError(t, err)

	mock.ExpectExec(`UPDATE contracts SET client_name=\$1, start_date=\$2, end_date=\$3, description=\$4 WHERE id=\$5`).
		WithArgs("New Client", startDate, endDate, "D", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	c.ID = 1
	err = repo.Update(c)
	assert.NoError(t, err)

	mock.ExpectExec(`DELETE FROM contracts WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(1)
	assert.NoError(t, err)
}
