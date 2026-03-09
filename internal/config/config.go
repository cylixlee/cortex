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
}

func Load() *Config {
	jwtExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if jwtExpiry == 0 {
		jwtExpiry = 24
	}

	return &Config{
		ChatProvider:   os.Getenv("CHAT_PROVIDER"),
		ChatBaseURL:    os.Getenv("CHAT_BASE_URL"),
		ChatAPIKey:     os.Getenv("CHAT_API_KEY"),
		ChatModel:      os.Getenv("CHAT_MODEL"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiryHours: jwtExpiry,
	}
}
