package postgres

import "time"

// Module represents a module in the database.
type Module struct {
	ID                int       `db:"id" json:"id"`
	ProjectID         int       `db:"project_id" json:"project_id"`
	Name              string    `db:"name" json:"name"`
	Description       *string   `db:"description" json:"description,omitempty"`
	ResponsibleUserID *string   `db:"responsible_user_id" json:"responsible_user_id,omitempty"`
	DateCreated       time.Time `db:"date_created" json:"date_created"`
	DateUpdated       time.Time `db:"date_updated" json:"date_updated"`
}
