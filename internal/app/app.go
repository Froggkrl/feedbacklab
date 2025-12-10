// Package app provides application initialization and routing.
package app

import (
	// Import swagger docs for API documentation.
	_ "innotech/docs"
	"innotech/internal/container"
	"innotech/internal/contract"
	"innotech/internal/documentations"
	"innotech/internal/files"
	"innotech/internal/projects"
	user_projects "innotech/internal/userprojects"
	"strconv"

	"innotech/internal/health"
	"innotech/internal/messageattachments"
	"innotech/internal/ticketattachments"
	"innotech/internal/ticketchats"
	"innotech/internal/tickets"
	"innotech/pkg/middleware"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	swagger "github.com/swaggo/fiber-swagger"
)

// Start initializes and starts the HTTP server with all registered routes.
func Start(container *container.Container) {
	app := fiber.New()

	app.Use(middleware.I18nMiddleware(container.I18nBundle))

	if container.Config.SwaggerUsername != "" && container.Config.SwaggerPassword != "" {
		authMiddleware := basicauth.New(basicauth.Config{
			Users: map[string]string{
				container.Config.SwaggerUsername: container.Config.SwaggerPassword,
			},
		})
		app.Use("/swagger", authMiddleware)
		app.Use("/docs", authMiddleware)
	}

	// Swagger UI endpoint
	app.Get("/swagger/*", swagger.WrapHandler)

	// Static Swagger HTML for offline viewing
	app.Static("/docs", "./docs")

	health.RegisterRoutes(app, container.HealthHandler)

	tickets.RegisterRoutes(app, container.TicketHandler)
	ticketchats.RegisterRoutes(app, container.TicketChatsHandler)
	ticketattachments.RegisterRoutes(app, container.TicketAttachmentsHandler)
	messageattachments.RegisterRoutes(app, container.MessageAttachmentsHandler)
	contract.RegisterRoutes(app, container.ContractHandler)
	projects.RegisterRoutes(app, container.ProjectHandler)
	documentations.RegisterRoutes(app, container.DocumentationHandler)
	user_projects.RegisterRoutes(app, container.UserProjectHandler)

	files.RegisterRoutes(app, container.FileHandler)

	log.Printf(" Server running on port %d\n", container.Config.AppPort)
	if err := app.Listen(":" + strconv.Itoa(container.Config.AppPort)); err != nil {
		log.Fatalf("failed to start feedbacklab: %v", err)
	}
}
