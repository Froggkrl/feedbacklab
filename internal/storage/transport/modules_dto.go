package transport

// CreateModuleDTO represents the data structure for creating a module.
type CreateModuleDTO struct {
	ProjectID         int     `json:"project_id" validate:"required"`
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	Description       *string `json:"description,omitempty"`
	ResponsibleUserID *string `json:"responsible_user_id,omitempty" validate:"omitempty,uuid4"`
}

// UpdateModuleDTO represents the data structure for updating a module.
type UpdateModuleDTO struct {
	Name              string  `json:"name" validate:"required,min=2,max=255"`
	Description       *string `json:"description,omitempty"`
	ResponsibleUserID *string `json:"responsible_user_id,omitempty" validate:"omitempty,uuid4"`
}
