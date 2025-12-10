package postgres

import "time"

// TicketAttachment represents a file attachment for a ticket in the database.
type TicketAttachment struct {
	ID          int       `db:"id" json:"id"`
	TicketID    int       `db:"ticket_id" json:"ticket_id"`
	FilePath    string    `db:"file_path" json:"file_path"`
	UploadedBy  string    `db:"uploaded_by" json:"uploaded_by"`
	FileType    *string   `db:"file_type" json:"file_type,omitempty"`
	Description *string   `db:"description" json:"description,omitempty"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
