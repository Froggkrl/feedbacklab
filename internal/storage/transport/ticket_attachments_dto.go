package transport

// CreateTicketAttachmentDTO represents the data structure for creating a ticket attachment.
type CreateTicketAttachmentDTO struct {
	TicketID    int     `json:"ticket_id" validate:"required"`
	FilePath    string  `json:"file_path" validate:"required"`
	UploadedBy  string  `json:"uploaded_by" validate:"required,uuid4"`
	FileType    *string `json:"file_type,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateTicketAttachmentDTO represents the data structure for updating a ticket attachment.
type UpdateTicketAttachmentDTO struct {
	FilePath    string  `json:"file_path" validate:"omitempty"`
	FileType    *string `json:"file_type,omitempty"`
	Description *string `json:"description,omitempty"`
}
