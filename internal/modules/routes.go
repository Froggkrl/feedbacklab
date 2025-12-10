package modules

import (
	"innotech/internal/storage/transport"
	"innotech/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers HTTP routes for module operations.
func RegisterRoutes(app *fiber.App, h *Handler) {
	api := app.Group("/api/modules")

	api.Get("/", h.GetAll)
	api.Get("/:id", h.GetByID)
	api.Post("/", middleware.ValidateBody[transport.CreateModuleDTO](h.Create))
	api.Put("/:id", middleware.ValidateBody[transport.UpdateModuleDTO](h.Update))
	api.Delete("/:id", h.Delete)
}
