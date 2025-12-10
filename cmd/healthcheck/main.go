// Package main is the entry point for the healthcheck service.
package main

import (
	"fmt"
	"innotech/config"
	"innotech/internal/health"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	target := fmt.Sprintf("http://app:%d/api/health", cfg.AppPort)
	service := health.NewHTTPHealthService(target, 3*time.Second)
	handler := health.NewHandler(service)

	app := fiber.New()
	health.RegisterRoutes(app, handler)

	port := cfg.HealthPort

	log.Printf("healthcheck service running on port %d (target: %s)\n", port, target)

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start healthcheck service: %v", err)
	}
}
