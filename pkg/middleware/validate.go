package middleware

import (
	"innotech/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateBody creates a middleware that validates request body against the provided type.
func ValidateBody[T any](next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T

		logger.Debug("starting request body validation",
			"method", c.Method(),
			"path", c.Path(),
			"ip", c.IP(),
		)

		// Парсинг тела запроса
		if err := c.BodyParser(&body); err != nil {
			logger.Error("request body parsing failed",
				"method", c.Method(),
				"path", c.Path(),
				"error", err.Error(),
				"content_type", c.Get("Content-Type"),
			)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON: " + err.Error(),
			})
		}

		// Валидация данных
		if err := validate.Struct(&body); err != nil {
			errors := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				errors[e.Field()] = e.Tag()
			}

			logger.Warn("request validation failed",
				"method", c.Method(),
				"path", c.Path(),
				"validation_errors", errors,
				"error_count", len(errors),
			)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"validation_errors": errors,
			})
		}

		logger.Debug("request validation passed",
			"method", c.Method(),
			"path", c.Path(),
		)

		c.Locals("body", &body)
		return next(c)
	}
}
