// Package transport contains data transfer objects (DTOs) for API requests and responses.
package transport

// CreateDocumentationDTO represents the data structure for creating a documentation entry.
type CreateDocumentationDTO struct {
	ProjectID  int     `json:"project_id" validate:"required"`
	FilePath   string  `json:"file_path" validate:"required"`
	Version    *string `json:"version,omitempty" validate:"omitempty"`
	UploadedBy *string `json:"uploaded_by,omitempty" validate:"omitempty,uuid4"`
}

// UpdateDocumentationDTO represents the data structure for updating a documentation entry.
type UpdateDocumentationDTO struct {
	FilePath   string  `json:"file_path" validate:"required"`
	Version    *string `json:"version,omitempty" validate:"omitempty"`
	UploadedBy *string `json:"uploaded_by,omitempty" validate:"omitempty,uuid4"`
}
