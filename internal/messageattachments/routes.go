package messageattachments

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for message attachment operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/message_attachments")

	api.Get("/:id", h.GetByID)
	api.Get("/chat/:chat_id", h.GetByChatID)
	api.Post("/", middleware.ValidateBody[transport.CreateMessageAttachmentDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateMessageAttachmentDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
