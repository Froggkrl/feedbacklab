// Package projects provides project management functionality.
package projects

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"innotech/pkg/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for project operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	logger.Info("project handler initialized")
	return &Handler{service: service}
}

// Create godoc
// @Summary Создать новый проект
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body transport.CreateProjectDTO true "Project Data"
// @Success 201 {object} postgres.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects [post]
// Create handles the creation of a new project.
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateProjectDTO)
	p := postgres.Project{
		Name:            dto.Name,
		Description:     dto.Description,
		GitlabProjectID: dto.GitlabProjectID,
		MattermostTeam:  dto.MattermostTeam,
	}

	logger.Info("handler: create project request",
		"name", p.Name,
		"path", c.Path(),
	)

	if err := h.service.Create(c.Context(), &p); err != nil {
		logger.Error("handler: create project failed",
			"error", err.Error(),
			"name", p.Name,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: project created",
		"id", p.ID,
		"name", p.Name,
	)
	return c.Status(fiber.StatusCreated).JSON(p)
}

// GetByID godoc
// @Summary Получить проект по ID
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} postgres.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/projects/{id} [get]
// GetByID retrieves a project by its ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Warn("handler: invalid id param",
			"param", c.Params("id"),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	logger.Debug("handler: get project by id", "id", id)

	p, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		logger.Error("handler: get by id failed",
			"id", id,
			"error", err.Error(),
		)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Debug("handler: project retrieved successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return c.JSON(p)
}

// GetAll godoc
// @Summary Получить все проекты
// @Tags Projects
// @Produce json
// @Success 200 {array} postgres.Project
// @Failure 500 {object} map[string]string
// @Router /api/projects [get]
// GetAll retrieves all projects.
func (h *Handler) GetAll(c *fiber.Ctx) error {
	logger.Debug("handler: get all projects")

	ps, err := h.service.GetAll(c.Context())
	if err != nil {
		logger.Error("handler: get all failed",
			"error", err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: projects list retrieved",
		"count", len(ps),
	)
	return c.JSON(ps)
}

// Update godoc
// @Summary Обновить проект
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body transport.UpdateProjectDTO true "Project Data"
// @Success 200 {object} postgres.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects/{id} [put]
// Update handles the update of an existing project.
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Warn("handler: invalid id param",
			"param", c.Params("id"),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	dto := c.Locals("body").(*transport.UpdateProjectDTO)
	p := postgres.Project{
		ID:              id,
		Name:            dto.Name,
		Description:     dto.Description,
		GitlabProjectID: dto.GitlabProjectID,
		MattermostTeam:  dto.MattermostTeam,
	}

	logger.Info("handler: update project request",
		"id", p.ID,
		"name", p.Name,
	)

	if err := h.service.Update(c.Context(), &p); err != nil {
		logger.Error("handler: update project failed",
			"id", p.ID,
			"error", err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: project updated successfully",
		"id", p.ID,
		"name", p.Name,
	)
	return c.JSON(p)
}

// Delete godoc
// @Summary Удалить проект
// @Tags Projects
// @Param id path int true "Project ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/projects/{id} [delete]
// Delete handles the deletion of a project.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Warn("handler: invalid id param",
			"param", c.Params("id"),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	logger.Warn("handler: delete project request", "id", id)

	if err := h.service.Delete(c.Context(), id); err != nil {
		logger.Error("handler: delete project failed",
			"id", id,
			"error", err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("handler: project deleted successfully", "id", id)
	return c.SendStatus(fiber.StatusNoContent)
}
