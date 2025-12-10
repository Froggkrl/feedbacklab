// Package postgres contains database model definitions for PostgreSQL.
package postgres

import "time"

// Documentation represents a documentation file in the database.
type Documentation struct {
	ID          int       `db:"id" json:"id"`
	ProjectID   int       `db:"project_id" json:"project_id"`
	FilePath    string    `db:"file_path" json:"file_path"`
	Version     *string   `db:"version" json:"version,omitempty"`
	UploadedBy  *string   `db:"uploaded_by" json:"uploaded_by,omitempty"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
