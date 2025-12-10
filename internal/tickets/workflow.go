package tickets

import "context"

type CreateIncidentInput struct {
	TicketID    int
	Title       string
	Description string
	Assignee    *string
}

func CreateIncidentWorkflow(ctx context.Context, input CreateIncidentInput) (string, error) {

	url, err := CreateIncidentActivity(ctx, input)
	if err != nil {
		return "", err
	}

	return url, nil
}
