package postgres

import "time"

// Project represents a project in the database.
type Project struct {
	ID              int       `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Description     *string   `db:"description" json:"description,omitempty"`
	GitlabProjectID *int      `db:"gitlab_project_id" json:"gitlab_project_id,omitempty"`
	MattermostTeam  *string   `db:"mattermost_team" json:"mattermost_team,omitempty"`
	DateCreated     time.Time `db:"date_created" json:"date_created"`
	DateUpdated     time.Time `db:"date_updated" json:"date_updated"`
}
