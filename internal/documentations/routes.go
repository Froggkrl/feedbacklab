package documentations

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for documentation operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/documentations")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Post("/", middleware.ValidateBody[transport.CreateDocumentationDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateDocumentationDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
