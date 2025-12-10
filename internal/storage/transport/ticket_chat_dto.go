package transport

// CreateTicketChatDTO represents the data structure for creating a ticket chat message.
type CreateTicketChatDTO struct {
	TicketID            int     `json:"ticket_id" validate:"required"`
	SenderID            string  `json:"sender_id" validate:"required,uuid4"`
	SenderRole          string  `json:"sender_role" validate:"required,oneof=customer manager developer tester"`
	Message             string  `json:"message" validate:"required,min=1"`
	MessageType         string  `json:"message_type" validate:"oneof=text file system"`
	MattermostMessageID *string `json:"mattermost_message_id,omitempty"`
}

// UpdateTicketChatDTO represents the data structure for updating a ticket chat message.
type UpdateTicketChatDTO struct {
	Message     string `json:"message" validate:"omitempty,min=1"`
	MessageType string `json:"message_type" validate:"omitempty,oneof=text file system"`
}
