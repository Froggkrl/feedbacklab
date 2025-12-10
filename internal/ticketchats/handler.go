// Package ticketchats provides ticket chat message management functionality.
package ticketchats

import (
	"innotech/internal/storage/postgres"
	"innotech/internal/storage/transport"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for ticket chat operations.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary создание сообщения тикета
// @Tags TicketChats
// @Accept json
// @Produce json
// @Param chat body transport.CreateTicketChatDTO true "Chat"
// @Success 201 {object} postgres.TicketChat
// @Failure 400 {object} map[string]string
// @Router /ticket_chats/ [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	dto := c.Locals("body").(*transport.CreateTicketChatDTO)

	chat := postgres.TicketChat{
		TicketID:            dto.TicketID,
		SenderID:            dto.SenderID,
		SenderRole:          dto.SenderRole,
		Message:             dto.Message,
		MessageType:         dto.MessageType,
		MattermostMessageID: dto.MattermostMessageID,
	}

	if err := h.service.Create(c.Context(), &chat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(chat)
}

// GetByID godoc
// @Summary получить сообщение по ID
// @Tags TicketChats
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} postgres.TicketChat
// @Failure 404 {object} map[string]string
// @Router /ticket_chats/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	chat, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(chat)
}

// GetByTicketID godoc
// @Summary получить все сообщения тикета
// @Tags TicketChats
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Success 200 {array} postgres.TicketChat
// @Router /ticket_chats/ticket/{ticket_id} [get]
func (h *Handler) GetByTicketID(c *fiber.Ctx) error {
	ticketID, err := strconv.Atoi(c.Params("ticket_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ticket id"})
	}
	chats, err := h.service.GetByTicketID(c.Context(), ticketID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chats)
}

// Update godoc
// @Summary обновить сообщение
// @Tags TicketChats
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param chat body transport.UpdateTicketChatDTO true "Chat"
// @Success 200 {object} postgres.TicketChat
// @Router /ticket_chats/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	dto := c.Locals("body").(*transport.UpdateTicketChatDTO)

	chat := postgres.TicketChat{
		ID:          id,
		Message:     dto.Message,
		MessageType: dto.MessageType,
	}

	if err := h.service.Update(c.Context(), &chat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(chat)
}

// Delete godoc
// @Summary удалить сообщение
// @Tags TicketChats
// @Param id path int true "ID"
// @Success 204
// @Router /ticket_chats/{id} [delete]
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
