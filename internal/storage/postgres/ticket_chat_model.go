package postgres

import "time"

// TicketChat represents a chat message in a ticket in the database.
type TicketChat struct {
	ID                  int       `db:"id" json:"id"`
	TicketID            int       `db:"ticket_id" json:"ticket_id"`
	SenderID            string    `db:"sender_id" json:"sender_id"`
	SenderRole          string    `db:"sender_role" json:"sender_role"`
	Message             string    `db:"message" json:"message"`
	MessageType         string    `db:"message_type" json:"message_type"`
	MattermostMessageID *string   `db:"mattermost_message_id" json:"mattermost_message_id,omitempty"`
	DateCreated         time.Time `db:"date_created" json:"date_created"`
	DateUpdated         time.Time `db:"date_updated" json:"date_updated"`
}
