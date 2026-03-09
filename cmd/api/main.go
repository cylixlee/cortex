package main

import (
	"log"

	"github.com/cylixlee/cortex/internal/config"
	"github.com/cylixlee/cortex/internal/handlers"
	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	llmClient, err := llm.NewClient(cfg.ChatBaseURL, cfg.ChatAPIKey, "deepseek-chat")
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}
	chatHandler := handlers.NewChatHandler(llmClient)

	r := gin.Default()

	r.Use(corsMiddleware())

	r.POST("/api/chat", chatHandler.Chat)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
