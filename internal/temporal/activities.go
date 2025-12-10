package temporal

import (
	"context"
	"fmt"
)

type GitlabIssueResult struct {
	IID    int
	WebURL string
}

func CreateGitlabIncidentActivity(ctx context.Context, input CreateIncidentInput) (GitlabIssueResult, error) {
	// TODO: подставить GitLab API

	fmt.Println("Creating Gitlab Incident:", input.Title)

	return GitlabIssueResult{
		IID:    123,
		WebURL: "https://gitlab.com/group/project/-/issues/123",
	}, nil
}

func SaveIncidentToDBActivity(ctx context.Context, ticketID int, result GitlabIssueResult) error {

	fmt.Println("Saving to DB ticket:", ticketID, result.WebURL)
	return nil
}
