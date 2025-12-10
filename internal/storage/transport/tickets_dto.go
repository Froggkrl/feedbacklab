package transport

// CreateTicketDTO represents the data structure for creating a ticket.
type CreateTicketDTO struct {
	ProjectID           int     `json:"project_id" validate:"required"`
	ModuleID            *int    `json:"module_id,omitempty"`
	ContractID          int     `json:"contract_id" validate:"required"`
	CreatedBy           string  `json:"created_by" validate:"required,uuid4"`
	AssignedTo          *string `json:"assigned_to,omitempty" validate:"omitempty,uuid4"`
	Title               string  `json:"title" validate:"required,min=3,max=255"`
	Message             string  `json:"message" validate:"required"`
	Status              string  `json:"status" validate:"omitempty,oneof=open in_progress resolved closed"`
	GitlabIssueURL      *string `json:"gitlab_issue_url,omitempty" validate:"omitempty,url"`
	MattermostThreadURL *string `json:"mattermost_thread_url,omitempty" validate:"omitempty,url"`
}

// UpdateTicketDTO represents the data structure for updating a ticket.
type UpdateTicketDTO struct {
	AssignedTo          *string `json:"assigned_to,omitempty" validate:"omitempty,uuid4"`
	Title               string  `json:"title" validate:"required,min=3,max=255"`
	Message             string  `json:"message" validate:"required"`
	Status              string  `json:"status" validate:"required,oneof=open in_progress resolved closed"`
	GitlabIssueURL      *string `json:"gitlab_issue_url,omitempty" validate:"omitempty,url"`
	MattermostThreadURL *string `json:"mattermost_thread_url,omitempty" validate:"omitempty,url"`
}
