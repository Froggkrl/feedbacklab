package postgres

import "time"

// MessageAttachment represents a file attachment for a chat message in the database.
type MessageAttachment struct {
	ID          int       `db:"id" json:"id"`
	ChatID      int       `db:"chat_id" json:"chat_id"`
	FilePath    string    `db:"file_path" json:"file_path"`
	UploadedBy  string    `db:"uploaded_by" json:"uploaded_by"`
	FileType    *string   `db:"file_type" json:"file_type,omitempty"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
