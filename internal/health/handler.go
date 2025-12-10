// Package health provides health check functionality for the application.
package health

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for health check operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// CheckHealth godoc
// @Summary проверка здоровья сервиса
// @Description возвращает статус OK
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) CheckHealth(c *fiber.Ctx) error {
	ctx := context.Background()
	if err := h.service.Check(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
	})
}
