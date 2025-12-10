package documentations

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for documentation operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Создать новую документацию
// @Tags Documentations
// @Accept json
// @Produce json
// @Param documentation body transport.CreateDocumentationDTO true "Documentation Data"
// @Success 201 {object} postgres.Documentation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documentations [post]
// Create handles the creation of a new documentation entry.
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateDocumentationDTO)
	d := postgres.Documentation{
		ProjectID:  dto.ProjectID,
		FilePath:   dto.FilePath,
		Version:    dto.Version,
		UploadedBy: dto.UploadedBy,
	}
	if err := h.service.Create(c.Context(), &d); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(d)
}

// GetByID godoc
// @Summary Получить документацию по ID
// @Tags Documentations
// @Produce json
// @Param id path int true "Documentation ID"
// @Success 200 {object} postgres.Documentation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/documentations/{id} [get]
// GetByID retrieves a documentation entry by its ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	d, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(d)
}

// GetAll godoc
// @Summary Получить всю документацию
// @Tags Documentations
// @Produce json
// @Success 200 {array} postgres.Documentation
// @Failure 500 {object} map[string]string
// @Router /api/documentations [get]
// GetAll retrieves all documentation entries.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	docs, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(docs)
}

// Update godoc
// @Summary Обновить документацию
// @Tags Documentations
// @Accept json
// @Produce json
// @Param id path int true "Documentation ID"
// @Param documentation body transport.UpdateDocumentationDTO true "Documentation Data"
// @Success 200 {object} postgres.Documentation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documentations/{id} [put]
// Update handles the update of an existing documentation entry.
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	dto := c.Locals("body").(*transport.UpdateDocumentationDTO)
	d := postgres.Documentation{
		ID:         id,
		FilePath:   dto.FilePath,
		Version:    dto.Version,
		UploadedBy: dto.UploadedBy,
	}
	if err := h.service.Update(c.Context(), &d); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(d)
}

// Delete godoc
// @Summary Удалить документацию
// @Tags Documentations
// @Param id path int true "Documentation ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documentations/{id} [delete]
// Delete handles the deletion of a documentation entry.
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
