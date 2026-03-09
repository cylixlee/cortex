package config

import (
	"log"
	"os"

	dotenv "github.com/joho/godotenv"
)

func init() {
	if err := dotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

type Config struct {
	ChatBaseURL string
	ChatAPIKey  string
}

func Load() *Config {
	return &Config{
		ChatBaseURL: os.Getenv("CHAT_BASE_URL"),
		ChatAPIKey:  os.Getenv("CHAT_API_KEY"),
	}
}

