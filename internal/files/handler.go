// Package files - package to files entity
package files

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for files operations.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new Handler instance.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Upload godoc
// @Summary Загрузить файл
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /files/upload [post]
// Upload POST /files/upload
func (h *Handler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("no file provided", "err", err)
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	savePath := "/tmp/" + file.Filename

	if err := c.SaveFile(file, savePath); err != nil {
		h.logger.Error("cannot save uploaded file", "err", err)
		return fiber.NewError(fiber.StatusInternalServerError, "file save error")
	}

	url, err := h.service.UploadFile(c.Context(), file.Filename, savePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "upload to storage failed")
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}
