// Package modules provides module management functionality.
package modules

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for module operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Создать новый модуль
// @Tags Modules
// @Accept json
// @Produce json
// @Param module body transport.CreateModuleDTO true "Module Data"
// @Success 201 {object} postgres.Module
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/modules [post]
// Create handles the creation of a new module.
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateModuleDTO)
	m := postgres.Module{
		ProjectID:         dto.ProjectID,
		Name:              dto.Name,
		Description:       dto.Description,
		ResponsibleUserID: dto.ResponsibleUserID,
	}
	if err := h.service.Create(c.Context(), &m); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(m)
}

// GetByID godoc
// @Summary Получить модуль по ID
// @Tags Modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} postgres.Module
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/modules/{id} [get]
// GetByID retrieves a module by its ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	m, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(m)
}

// GetAll godoc
// @Summary Получить все модули
// @Tags Modules
// @Produce json
// @Success 200 {array} postgres.Module
// @Failure 500 {object} map[string]string
// @Router /api/modules [get]
// GetAll retrieves all modules.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	ms, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ms)
}

// Update godoc
// @Summary Обновить модуль
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path int true "Module ID"
// @Param module body transport.UpdateModuleDTO true "Module Data"
// @Success 200 {object} postgres.Module
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/modules/{id} [put]
// Update handles the update of an existing module.
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	dto := c.Locals("body").(*transport.UpdateModuleDTO)
	m := postgres.Module{
		ID:                id,
		Name:              dto.Name,
		Description:       dto.Description,
		ResponsibleUserID: dto.ResponsibleUserID,
	}
	if err := h.service.Update(c.Context(), &m); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(m)
}

// Delete godoc
// @Summary Удалить модуль
// @Tags Modules
// @Param id path int true "Module ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/modules/{id} [delete]
// Delete handles the deletion of a module.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
