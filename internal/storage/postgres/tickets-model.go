package postgres

import "time"

// Ticket represents a ticket in the database.
type Ticket struct {
	ID                  int       `db:"id" json:"id"`
	ProjectID           int       `db:"project_id" json:"project_id"`
	ModuleID            *int      `db:"module_id" json:"module_id,omitempty"`
	ContractID          int       `db:"contract_id" json:"contract_id"`
	CreatedBy           string    `db:"created_by" json:"created_by"`
	AssignedTo          *string   `db:"assigned_to" json:"assigned_to,omitempty"`
	Title               string    `db:"title" json:"title"`
	Message             string    `db:"message" json:"message"`
	Status              string    `db:"status" json:"status"`
	GitlabIssueURL      *string   `db:"gitlab_issue_url" json:"gitlab_issue_url,omitempty"`
	MattermostThreadURL *string   `db:"mattermost_thread_url" json:"mattermost_thread_url,omitempty"`
	DateCreated         time.Time `db:"date_created" json:"date_created"`
	DateUpdated         time.Time `db:"date_updated" json:"date_updated"`
}
