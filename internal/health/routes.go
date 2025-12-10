package health

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers HTTP routes for health check operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/health", h.CheckHealth)
}
