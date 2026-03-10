package config

import (
	"log"
	"os"
	"strconv"

	dotenv "github.com/joho/godotenv"
)

func init() {
	if err := dotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

type Config struct {
	ChatProvider   string
	ChatBaseURL    string
	ChatAPIKey     string
	ChatModel      string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiryHours int

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool

	RedisURL string

	EmbeddingProvider string
	EmbeddingBaseURL  string
	EmbeddingAPIKey   string
	EmbeddingModel    string
}

func Load() *Config {
	jwtExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if jwtExpiry == 0 {
		jwtExpiry = 24
	}

	minioUseSSL := os.Getenv("MINIO_USE_SSL") == "true" || os.Getenv("MINIO_USE_SSL") == "1"

	return &Config{
		ChatProvider:   os.Getenv("CHAT_PROVIDER"),
		ChatBaseURL:    os.Getenv("CHAT_BASE_URL"),
		ChatAPIKey:     os.Getenv("CHAT_API_KEY"),
		ChatModel:      os.Getenv("CHAT_MODEL"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiryHours: jwtExpiry,

		MinIOEndpoint:  os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinIOBucket:    os.Getenv("MINIO_BUCKET"),
		MinIOUseSSL:    minioUseSSL,

		RedisURL: os.Getenv("REDIS_URL"),

		EmbeddingProvider: os.Getenv("EMBEDDING_PROVIDER"),
		EmbeddingBaseURL:  os.Getenv("EMBEDDING_BASE_URL"),
		EmbeddingAPIKey:   os.Getenv("EMBEDDING_API_KEY"),
		EmbeddingModel:    os.Getenv("EMBEDDING_MODEL"),
	}
}
