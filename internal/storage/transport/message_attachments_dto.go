package transport

// CreateMessageAttachmentDTO represents the data structure for creating a message attachment.
type CreateMessageAttachmentDTO struct {
	ChatID     int     `json:"chat_id" validate:"required"`
	FilePath   string  `json:"file_path" validate:"required"`
	UploadedBy string  `json:"uploaded_by" validate:"required,uuid4"`
	FileType   *string `json:"file_type,omitempty"`
}

// UpdateMessageAttachmentDTO represents the data structure for updating a message attachment.
type UpdateMessageAttachmentDTO struct {
	FilePath string  `json:"file_path" validate:"omitempty"`
	FileType *string `json:"file_type,omitempty"`
}
