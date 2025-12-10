// Package config provides configuration loading and management for the application.
package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration settings.
type Config struct {
	AppPort       string
	HealthPort    string
	MigrationsDir string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSL       string
	DatabaseURL string

	// MinIO
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
	AppPort         int
	HealthPort      int
	DatabaseURL     string
	ParsedDBURL     *url.URL
	MigrationsDir   string
	SwaggerUsername string
	SwaggerPassword string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucket     string
	MinioUseSSL     bool
}

// Load reads configuration from environment variables and returns a Config instance.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		HealthPort:    getEnv("HEALTH_PORT", "8081"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "./migrations"),

		// DB
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "feedback"),
		DBPassword: getEnv("DB_PASSWORD", "feedback"),
		DBName:     getEnv("DB_NAME", "innotech"),
		DBSSL:      getEnv("DB_SSL", "disable"),

		// MinIO
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioBucket:    getEnv("MINIO_BUCKET", "feedback-files"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
	cfg := &Config{}

	appPort, err := getEnvInt("APP_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}
	cfg.AppPort = appPort

	healthPort, err := getEnvInt("HEALTH_PORT", 8081)
	if err != nil {
		return nil, fmt.Errorf("invalid HEALTH_PORT: %w", err)
	}
	cfg.HealthPort = healthPort

	cfg.DatabaseURL = getEnv("DATABASE_URL",
		"postgres://feedback:feedback@db:5432/innotech?sslmode=disable",
	)

	parsedURL, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_URL: %w", err)
	}
	cfg.ParsedDBURL = parsedURL

	cfg.MigrationsDir = getEnv("MIGRATIONS_DIR", "./migrations")

	cfg.SwaggerUsername = getEnv("SWAGGER_USERNAME", "swagger")
	cfg.SwaggerPassword = getEnv("SWAGGER_PASSWORD", "swagger")

	cfg.MinioEndpoint = getEnv("MINIO_ENDPOINT", "localhost:9000")
	cfg.MinioAccessKey = getEnv("MINIO_ACCESS_KEY", "minioadmin")
	cfg.MinioSecretKey = getEnv("MINIO_SECRET_KEY", "minioadmin123")
	cfg.MinioBucket = getEnv("MINIO_BUCKET", "feedback-files")
	minioUseSSL, err := getEnvBool("MINIO_SSL", false)
	if err != nil {
		return nil, fmt.Errorf("invalid MINIO_SSL: %w", err)
	}
	cfg.MinioUseSSL = minioUseSSL

	// Сборка строки подключения
	cfg.DatabaseURL = buildPostgresURL(cfg)

	log.Println("config loaded successfully")
	return cfg
	log.Println("config loaded and parsed successfully")
	return cfg, nil
}

func buildPostgresURL(cfg *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSL,
	)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	if val, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	}
	return fallback, nil
}

func getEnvBool(key string, fallback bool) (bool, error) {
	if val, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return false, err
		}
		return parsed, nil
	}
	return fallback, nil
}
