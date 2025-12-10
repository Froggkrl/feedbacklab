// Package middleware provides HTTP middleware functions for the application.
package middleware

import (
	"github.com/gofiber/fiber/v2"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// I18nMiddleware creates a middleware that handles internationalization by detecting the user's language.
func I18nMiddleware(bundle *goi18n.Bundle) fiber.Handler {
	return func(c *fiber.Ctx) error {
		acceptLanguage := c.Get("Accept-Language", "ru")

		supportedLanguages := []language.Tag{
			language.Russian,
			language.English,
		}

		matcher := language.NewMatcher(supportedLanguages)
		tag, _ := language.MatchStrings(matcher, acceptLanguage)

		localizer := goi18n.NewLocalizer(bundle, tag.String())

		c.Locals("localizer", localizer)
		c.Locals("language", tag.String())

		return c.Next()
	}
}
