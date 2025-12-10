package tickets

import (
	"log"
	"log/slog"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	temporal client.Client
	log      *slog.Logger
}

func NewWorker(c client.Client, log *slog.Logger) *Worker {
	return &Worker{
		temporal: c,
		log:      log,
	}
}

func (w *Worker) Start() {
	w.log.Info("temporal worker: starting...")

	temporalWorker := worker.New(w.temporal, "incident-task-queue", worker.Options{})

	// Регистрируем workflow
	temporalWorker.RegisterWorkflow(CreateIncidentWorkflow)

	// Регистрируем activity
	temporalWorker.RegisterActivity(CreateIncidentActivity)

	// Запускаем цикл обработки
	if err := temporalWorker.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("unable to start Temporal worker: %v", err)
	}

	w.log.Info("temporal worker: stopped")
}
