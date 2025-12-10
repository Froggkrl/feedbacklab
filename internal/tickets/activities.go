package tickets

import (
	"context"
	"log/slog"
)

func CreateIncidentActivity(ctx context.Context, input CreateIncidentInput) (string, error) {
	slog.Info("activity: creating gitlab incident",
		slog.Int("ticket_id", input.TicketID),
		slog.String("title", input.Title),
	)

	// TODO: здесь будет вызов GitLab API

	fakeURL := "https://gitlab.example.com/incident/" + input.Title

	return fakeURL, nil
}
