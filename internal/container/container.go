// Package container provides dependency injection container for the application.
package container

import (
	"innotech/config"
	"innotech/internal/contract"
	"innotech/internal/documentations"
	"innotech/internal/files"
	"innotech/internal/health"
	"innotech/internal/messageattachments"
	"innotech/internal/projects"
	"innotech/internal/ticketattachments"
	"innotech/internal/ticketchats"
	"innotech/internal/tickets"
	user_projects "innotech/internal/userprojects"
	"innotech/pkg/db"
	"innotech/pkg/i18n"
	"innotech/pkg/logger"
	minio_client "innotech/pkg/minio"
	"log"

	"github.com/jmoiron/sqlx"
	"go.temporal.io/sdk/client"
	"log/slog"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// Container holds all application dependencies and services.
type Container struct {
	Config         *config.Config
	DB             *sqlx.DB
	TemporalClient client.Client
	Log            *slog.Logger

	Config                    *config.Config
	DB                        *sqlx.DB
	I18nBundle                *goi18n.Bundle
	HealthHandler             *health.Handler
	TicketHandler             *tickets.Handler
	TicketChatsHandler        *ticketchats.Handler
	TicketAttachmentsHandler  *ticketattachments.Handler
	MessageAttachmentsHandler *messageattachments.Handler
	ContractHandler           *contract.Handler
	ProjectHandler            *projects.Handler
	DocumentationHandler      *documentations.Handler
	UserProjectHandler        *user_projects.Handler
	FileHandler               *files.Handler
}

// New creates and initializes a new Container with all dependencies.
func New() *Container {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalln("failed to load config: ", err)
	}

	// global logger
	logg := logger.NewLogger()
	slog.SetDefault(logg)

	// DB
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// TEMPORAL CLIENT
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("Temporal connection failed: %v", err)
	}

	// SERVICES + REPOS + HANDLERS
	healthService := health.NewSelfHealthService()
	healthHandler := health.NewHandler(healthService)

	ticketRepo := tickets.NewRepository(database)
	ticketService := tickets.NewService(ticketRepo, temporalClient, logg)
	ticketHandler := tickets.NewHandler(ticketService)
	ticketService := tickets.NewService(ticketRepo)
	ticketHandler := tickets.NewHandler(ticketService, logger.Global)

	chatRepo := ticketchats.NewRepository(database)
	chatService := ticketchats.NewService(chatRepo)
	chatHandler := ticketchats.NewHandler(chatService)

	attachRepo := ticketattachments.NewRepository(database)
	attachService := ticketattachments.NewService(attachRepo)
	attachHandler := ticketattachments.NewHandler(attachService)

	msgAttachRepo := messageattachments.NewRepository(database)
	msgAttachService := messageattachments.NewService(msgAttachRepo)
	msgAttachHandler := messageattachments.NewHandler(msgAttachService)

	contractRepo := contract.NewRepository(database)
	contractService := contract.NewService(contractRepo)
	contractHandler := contract.NewHandler(contractService)

	projectRepo := projects.NewRepository(database)
	projectService := projects.NewService(projectRepo)
	projectHandler := projects.NewHandler(projectService)

	docRepo := documentations.NewRepository(database)
	docService := documentations.NewService(docRepo)
	docHandler := documentations.NewHandler(docService)

	userProjectRepo := user_projects.NewRepository(database)
	userProjectService := user_projects.NewService(userProjectRepo)
	userProjectHandler := user_projects.NewHandler(userProjectService)

	bundle := i18n.InitBundle()
	if err := i18n.LoadTranslations(bundle, "./locales"); err != nil {
		log.Printf("warning: failed to load translations: %v", err)
	}

	minioClient, err := minio_client.New(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}
	fileService := files.NewService(minioClient, logger.Global)
	fileHandler := files.NewHandler(fileService, logger.Global)

	return &Container{
		Config:         cfg,
		DB:             database,
		TemporalClient: temporalClient,
		Log:            logg,

		Config:                    cfg,
		DB:                        database,
		I18nBundle:                bundle,
		HealthHandler:             healthHandler,
		TicketHandler:             ticketHandler,
		TicketChatsHandler:        chatHandler,
		TicketAttachmentsHandler:  attachHandler,
		MessageAttachmentsHandler: msgAttachHandler,
		ContractHandler:           contractHandler,
		ProjectHandler:            projectHandler,
		DocumentationHandler:      docHandler,
		UserProjectHandler:        userProjectHandler,
		FileHandler:               fileHandler,
	}
}
