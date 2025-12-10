package postgres

import "time"

// UserProject represents a user-project relationship in the database.
type UserProject struct {
	UserID      string    `db:"user_id" json:"user_id"`
	ProjectID   int       `db:"project_id" json:"project_id"`
	Role        string    `db:"role" json:"role"`
	Permissions []string  `db:"permissions" json:"permissions"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
