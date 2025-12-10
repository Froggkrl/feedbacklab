package contract

import "time"

// Contract represents a contract entity.
type Contract struct {
	ID          int       `db:"id" json:"id"`
	ProjectID   int       `db:"project_id" json:"project_id"`
	ClientName  string    `db:"client_name" json:"client_name"`
	StartDate   time.Time `db:"start_date" json:"start_date"`
	EndDate     time.Time `db:"end_date" json:"end_date"`
	Description string    `db:"description" json:"description"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
