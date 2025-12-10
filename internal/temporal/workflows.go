package temporal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type CreateIncidentInput struct {
	TicketID    int
	Title       string
	Description string
	AssigneeID  int
}

func CreateIncidentWorkflow(ctx workflow.Context, input CreateIncidentInput) error {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 3,
		RetryPolicy: &workflow.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result GitlabIssueResult
	
	err := workflow.ExecuteActivity(ctx, CreateGitlabIncidentActivity, input).Get(ctx, &result)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, SaveIncidentToDBActivity, input.TicketID, result).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
