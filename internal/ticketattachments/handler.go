// Package ticketattachments provides ticket attachment management functionality.
package ticketattachments

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for ticket attachment operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary создать вложение тикета
// @Tags TicketAttachments
// @Accept json
// @Produce json
// @Param attachment body transport.CreateTicketAttachmentDTO true "Attachment"
// @Success 201 {object} postgres.TicketAttachment
// @Router /ticket_attachments/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateTicketAttachmentDTO)

	att := postgres.TicketAttachment{
		TicketID:    dto.TicketID,
		FilePath:    dto.FilePath,
		UploadedBy:  dto.UploadedBy,
		FileType:    dto.FileType,
		Description: dto.Description,
	}

	if err := h.service.Create(c.Context(), &att); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(att)
}

// GetByID godoc
// @Summary получить вложение по ID
// @Tags TicketAttachments
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} postgres.TicketAttachment
// @Router /ticket_attachments/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	att, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(att)
}

// GetByTicketID godoc
// @Summary получить все вложения тикета
// @Tags TicketAttachments
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Success 200 {array} postgres.TicketAttachment
// @Router /ticket_attachments/ticket/{ticket_id} [get]
func (h *Handler) GetByTicketID(c *fiber.Ctx) error {
	ticketID, err := strconv.Atoi(c.Params("ticket_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ticket id"})
	}
	list, err := h.service.GetByTicketID(c.Context(), ticketID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// Update godoc
// @Summary обновить вложение
// @Tags TicketAttachments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param attachment body transport.UpdateTicketAttachmentDTO true "Attachment"
// @Success 200 {object} postgres.TicketAttachment
// @Router /ticket_attachments/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	dto := c.Locals("body").(*transport.UpdateTicketAttachmentDTO)

	att := postgres.TicketAttachment{
		ID:          id,
		FilePath:    dto.FilePath,
		FileType:    dto.FileType,
		Description: dto.Description,
	}

	if err := h.service.Update(c.Context(), &att); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(att)
}

// Delete godoc
// @Summary удалить вложение
// @Tags TicketAttachments
// @Param id path int true "ID"
// @Success 204
// @Router /ticket_attachments/{id} [delete]
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
