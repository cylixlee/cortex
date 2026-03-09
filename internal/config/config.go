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
	ChatProvider string
	ChatBaseURL  string
	ChatAPIKey   string
	ChatModel    string
}

func Load() *Config {
	return &Config{
		ChatProvider: os.Getenv("CHAT_PROVIDER"),
		ChatBaseURL:  os.Getenv("CHAT_BASE_URL"),
		ChatAPIKey:   os.Getenv("CHAT_API_KEY"),
		ChatModel:    os.Getenv("CHAT_MODEL"),
	}
}
