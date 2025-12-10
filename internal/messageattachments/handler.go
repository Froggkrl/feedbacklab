// Package messageattachments provides message attachment management functionality.
package messageattachments

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for message attachment operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary создать вложение сообщения
// @Tags MessageAttachments
// @Accept json
// @Produce json
// @Param attachment body transport.CreateMessageAttachmentDTO true "Attachment"
// @Success 201 {object} postgres.MessageAttachment
// @Router /message_attachments/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateMessageAttachmentDTO)

	att := postgres.MessageAttachment{
		ChatID:     dto.ChatID,
		FilePath:   dto.FilePath,
		UploadedBy: dto.UploadedBy,
		FileType:   dto.FileType,
	}

	if err := h.service.Create(c.Context(), &att); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(att)
}

// GetByID godoc
// @Summary получить вложение по ID
// @Tags MessageAttachments
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} postgres.MessageAttachment
// @Router /message_attachments/{id} [get]
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

// GetByChatID godoc
// @Summary получить все вложения по chat_id
// @Tags MessageAttachments
// @Produce json
// @Param chat_id path int true "Chat ID"
// @Success 200 {array} postgres.MessageAttachment
// @Router /message_attachments/chat/{chat_id} [get]
func (h *Handler) GetByChatID(c *fiber.Ctx) error {
	chatID, err := strconv.Atoi(c.Params("chat_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chat id"})
	}
	list, err := h.service.GetByChatID(c.Context(), chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// Update godoc
// @Summary обновить вложение
// @Tags MessageAttachments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param attachment body transport.UpdateMessageAttachmentDTO true "Attachment"
// @Success 200 {object} postgres.MessageAttachment
// @Router /message_attachments/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	dto := c.Locals("body").(*transport.UpdateMessageAttachmentDTO)

	att := postgres.MessageAttachment{
		ID:       id,
		FilePath: dto.FilePath,
		FileType: dto.FileType,
	}

	if err := h.service.Update(c.Context(), &att); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(att)
}

// Delete godoc
// @Summary удалить вложение
// @Tags MessageAttachments
// @Param id path int true "ID"
// @Success 204
// @Router /message_attachments/{id} [delete]
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
