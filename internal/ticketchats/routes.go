package ticketchats

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for ticket chat operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/ticket_chats")

	api.Get("/:id", h.GetByID)
	api.Get("/ticket/:ticket_id", h.GetByTicketID)
	api.Post("/", middleware.ValidateBody[transport.CreateTicketChatDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateTicketChatDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
