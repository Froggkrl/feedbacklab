package tickets

import (
	"context"
	"log/slog"
	"strconv"

	"go.temporal.io/sdk/client"
)

// Service defines the interface for ticket business logic operations.
type Service interface {
	Create(ctx context.Context, t *postgres.Ticket) error
	GetByID(ctx context.Context, id int) (*postgres.Ticket, error)
	GetAll(ctx context.Context) ([]postgres.Ticket, error)
	Update(ctx context.Context, t *postgres.Ticket) error
	Delete(ctx context.Context, id int) error
}

type ticketService struct {
	repo           Repository
	temporalClient client.Client
	log            *slog.Logger
}

func NewService(repo Repository, temporalClient client.Client, log *slog.Logger) Service {
	return &ticketService{
		repo:           repo,
		temporalClient: temporalClient,
		log:            log,
	}
}

func (s *ticketService) Create(ctx context.Context, t *Ticket) error {
	s.log.Info("service: creating ticket")

	// стандартная логика
// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &ticketService{repo: repo}
}

func (s *ticketService) Create(ctx context.Context, t *postgres.Ticket) error {
	t.Status = "open"
	err := s.repo.Create(ctx, t)
	if err != nil {
		s.log.Error("service: failed to create ticket", slog.String("error", err.Error()))
		return err
	}

	// запускаем Temporal workflow НЕ ЛОМАЯ интерфейс
	workflowOptions := client.StartWorkflowOptions{
		ID:        "incident_workflow_ticket_" + strconv.Itoa(t.ID),
		TaskQueue: "incident-task-queue",
	}

	_, wfErr := s.temporalClient.ExecuteWorkflow(
		ctx,
		workflowOptions,
		"CreateIncidentWorkflow", // имя workflow
		CreateIncidentInput{
			TicketID:    t.ID,
			Title:       t.Title,
			Description: t.Message,
			Assignee:    t.AssignedTo, // *string
		},
	)

	if wfErr != nil {
		s.log.Error("service: failed to start workflow", slog.String("error", wfErr.Error()))
		// не возвращаем ошибку → тикет должен создаваться ВСЕГДА
	}

	return nil
}

func (s *ticketService) GetByID(ctx context.Context, id int) (*postgres.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ticketService) GetAll(ctx context.Context) ([]postgres.Ticket, error) {
	return s.repo.GetAll(ctx)
}

func (s *ticketService) Update(ctx context.Context, t *postgres.Ticket) error {
	return s.repo.Update(ctx, t)
}

func (s *ticketService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
