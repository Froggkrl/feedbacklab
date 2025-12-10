package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	t "innotech/internal/temporal"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "incident-queue", worker.Options{})

	w.RegisterWorkflow(t.CreateIncidentWorkflow)
	w.RegisterActivity(t.CreateGitlabIncidentActivity)
	w.RegisterActivity(t.SaveIncidentToDBActivity)

	log.Println("Temporal Worker started...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}

}
