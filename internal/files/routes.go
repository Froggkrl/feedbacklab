package files

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers HTTP routes for files operations.
func RegisterRoutes(router fiber.Router, h *Handler) {
	files := router.Group("/files")
	files.Post("/upload", h.Upload)
}
