// Package tickets provides ticket management functionality.
package tickets

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for ticket operations.
type Handler struct {
	service Service
}

func NewHandler(service Service, log *slog.Logger) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary создать тикет
// @Tags Tickets
// @Accept json
// @Produce json
// @Param ticket body transport.CreateTicketDTO true "Ticket"
// @Success 201 {object} postgres.Ticket
// @Failure 400 {object} map[string]string
// @Router /tickets [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateTicketDTO)

	t := postgres.Ticket{
		ProjectID:           dto.ProjectID,
		ModuleID:            dto.ModuleID,
		ContractID:          dto.ContractID,
		CreatedBy:           dto.CreatedBy,
		AssignedTo:          dto.AssignedTo,
		Title:               dto.Title,
		Message:             dto.Message,
		Status:              dto.Status,
		GitlabIssueURL:      dto.GitlabIssueURL,
		MattermostThreadURL: dto.MattermostThreadURL,
	}

	if err := h.service.Create(c.Context(), &t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(t)
}

// GetByID godoc
// @Summary получить тикет по ID
// @Tags Tickets
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} postgres.Ticket
// @Failure 404 {object} map[string]string
// @Router /tickets/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	t, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(t)
}

// GetAll godoc
// @Summary получить все тикеты
// @Tags Tickets
// @Produce json
// @Success 200 {object} postgres.Ticket
// @Failure 404 {object} map[string]string
// @Router /tickets/ [get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	tickets, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tickets)
}

// Update godoc
// @Summary обновить тикет
// @Tags Tickets
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ticket body transport.UpdateTicketDTO true "Ticket"
// @Success 200 {object} postgres.Ticket
// @Router /tickets/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	dto := c.Locals("body").(*transport.UpdateTicketDTO)

	t := postgres.Ticket{
		ID:                  id,
		Title:               dto.Title,
		Message:             dto.Message,
		Status:              dto.Status,
		AssignedTo:          dto.AssignedTo,
		GitlabIssueURL:      dto.GitlabIssueURL,
		MattermostThreadURL: dto.MattermostThreadURL,
	}

	if err := h.service.Update(c.Context(), &t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(t)
}

// Delete godoc
// @Summary удалить тикет
// @Tags Tickets
// @Param id path int true "ID"
// @Success 204
// @Router /tickets/{id} [delete]
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
