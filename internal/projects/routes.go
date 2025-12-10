package projects

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/logger"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for project operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/projects")

	// Middleware для логирования всех запросов
	api.Use(func(c *fiber.Ctx) error {
		logger.Info("HTTP request to projects API",
			"method", c.Method(),
			"path", c.Path(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)
		return c.Next()
	})

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Post("/", middleware.ValidateBody[transport.CreateProjectDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateProjectDTO](h.Update))
	api.Delete("/:id", h.Delete)

	logger.Info("project routes registration completed",
		"total_routes", 5,
		"base_path", "/api/projects",
	)
}
