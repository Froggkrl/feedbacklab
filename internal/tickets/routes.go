package tickets

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for ticket operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/tickets")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)

	api.Post("/", middleware.ValidateBody[transport.CreateTicketDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateTicketDTO](h.Update))

	api.Delete("/:id", h.Delete)
}
