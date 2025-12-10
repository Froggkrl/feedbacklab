// Package main is the entry point for the FeedbackLab API application.
package main

import (
	"innotech/internal/app"
	"innotech/internal/container"
	"innotech/pkg/db"
	"innotech/pkg/logger"
	"log"
)

// @title FeedbackLab API
// @version 0.1
// @description API для управления проектами, тикетами и документациейs
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	logger.Init()
	logger.Info("starting FeedbackLab API application")

	c := container.New()
	logger.Info("container initialized")

	if err := db.Migrate(c.Config.DatabaseURL, c.Config.MigrationsDir); err != nil {
		logger.Error("database migration failed",
			"error", err,
			"database_url", c.Config.DatabaseURL,
			"migrations_dir", c.Config.MigrationsDir,
		)
		log.Fatalf("Migration failed: %v", err)
	}
	logger.Info("database migrations completed successfully")

	logger.Info("starting application server",
		"version", "0.1",
	)

	app.Start(c)
}

// todo 1) сделать парсер в конфиге, 2) прикрутить slog, 3) названия в docker(https://github.com/wagoodman/dive), 4) линтер(https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322), 5) переводы, 6) minio 7)ошибки.
// todo 1) Makfeile для основных команд, 2) openAPI 3, 3) swagger static(html), 4) закрыть свагер, 5) clock-permission, 6)поразбираться с keycloak - получение токена(будто user), получение user info с scope permissions.
