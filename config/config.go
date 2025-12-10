package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
}

func Load() *Config {
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
	}

	// Сборка строки подключения
	cfg.DatabaseURL = buildPostgresURL(cfg)

	log.Println("config loaded successfully")
	return cfg
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
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
