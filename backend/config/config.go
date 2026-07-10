package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string

	// Database (assembled from parts)
	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPassword string
	DBSSLMode string

	// Redis (assembled from parts)
	RedisHost string
	RedisPort string

	JWTSecret  string
	JWTExpiryH int

	GoogleClientID     string
	GoogleClientSecret string
	GitHubClientID     string
	GitHubClientSecret string

	OpenAIKey string
	StripeKey    string
	StripeWebhookSecret string

	GCSBucket      string
	GCSCredentials string

	SendGridKey    string
	SendGridFrom   string

	FrontendURL   string
	FastAPIURL    string
	UploadsDir    string
	MigrationsDir string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "fastlearn_db"),
		DBUser:     getEnv("DB_USER", "fastlearn"),
		DBPassword: getEnv("DB_PASSWORD", "dev_password"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiryH: 24,

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),

		OpenAIKey: os.Getenv("OPENAI_API_KEY"),
		StripeKey:           os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),

		GCSBucket:      getEnv("GCS_BUCKET", "fastlearn-videos"),
		GCSCredentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),

		SendGridKey:  os.Getenv("SENDGRID_API_KEY"),
		SendGridFrom: getEnv("SENDGRID_FROM", "noreply@fastlearn.ai"),

		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
		FastAPIURL:    getEnv("FASTAPI_URL", "http://localhost:8001"),
		UploadsDir:    getEnv("UPLOADS_DIR", "./uploads/video"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "./db/migrations"),
	}
}

// DatabaseURL assembles the postgres DSN from individual parts.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

// RedisURL assembles the redis URL from individual parts.
func (c *Config) RedisURL() string {
	return fmt.Sprintf("redis://%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
