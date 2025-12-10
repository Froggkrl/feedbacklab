package ticketattachments

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for ticket attachment operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/ticket_attachments")

	api.Get("/:id", h.GetByID)
	api.Get("/ticket/:ticket_id", h.GetByTicketID)
	api.Post("/", middleware.ValidateBody[transport.CreateTicketAttachmentDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateTicketAttachmentDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
