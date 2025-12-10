package transport

// CreateProjectDTO represents the data structure for creating a project.
type CreateProjectDTO struct {
	Name            string  `json:"name" validate:"required,min=2,max=255"`
	Description     *string `json:"description,omitempty"`
	GitlabProjectID *int    `json:"gitlab_project_id,omitempty"`
	MattermostTeam  *string `json:"mattermost_team,omitempty"`
}

// UpdateProjectDTO represents the data structure for updating a project.
type UpdateProjectDTO struct {
	Name            string  `json:"name" validate:"required,min=2,max=255"`
	Description     *string `json:"description,omitempty"`
	GitlabProjectID *int    `json:"gitlab_project_id,omitempty"`
	MattermostTeam  *string `json:"mattermost_team,omitempty"`
}
