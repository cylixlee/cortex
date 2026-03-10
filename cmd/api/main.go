package main

import (
	"context"
	"log"

	"github.com/cylixlee/cortex/internal/config"
	"github.com/cylixlee/cortex/internal/graceful"
	"github.com/cylixlee/cortex/internal/handlers"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/cylixlee/cortex/internal/service"
	"github.com/cylixlee/cortex/internal/worker"
	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/cylixlee/cortex/pkg/middleware"
	"github.com/cylixlee/cortex/pkg/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	cfg := config.Load()

	if err := repository.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	DB = repository.DB

	if err := repository.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	llmClient, err := llm.NewClient(
		context.Background(),
		cfg.ChatProvider,
		cfg.ChatBaseURL,
		cfg.ChatAPIKey,
		cfg.ChatModel,
	)
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	minioClient, err := storage.NewMinIOClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	if err := minioClient.CreateBucketIfNotExists(context.Background()); err != nil {
		log.Fatalf("Failed to create MinIO bucket: %v", err)
	}

	userRepo := repository.NewUserRepository(DB)
	conversationRepo := repository.NewConversationRepository(DB)
	messageRepo := repository.NewMessageRepository(DB)
	skillRepo := repository.NewSkillRepository(DB)
	documentRepo := repository.NewDocumentRepository(DB)
	chunkRepo := repository.NewChunkRepository(DB)
	referenceRepo := repository.NewReferenceRepository(DB)

	embeddingClient, err := llm.NewEmbeddingClient(
		context.Background(),
		cfg.EmbeddingProvider,
		cfg.EmbeddingBaseURL,
		cfg.EmbeddingAPIKey,
		cfg.EmbeddingModel,
	)
	if err != nil {
		log.Fatalf("Failed to create embedding client: %v", err)
	}

	skillService := service.NewSkillService(
		skillRepo,
		documentRepo,
		chunkRepo,
		referenceRepo,
		minioClient,
		embeddingClient,
	)

	skillWorker, err := worker.NewSkillWorker(
		skillRepo,
		documentRepo,
		chunkRepo,
		referenceRepo,
		minioClient,
		llmClient,
		cfg,
	)
	if err != nil {
		log.Fatalf("Failed to create skill worker: %v", err)
	}

	go skillWorker.Start(context.Background())

	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	conversationService := service.NewConversationService(conversationRepo, messageRepo)
	chatService := service.NewChatService(conversationRepo, messageRepo, userRepo, llmClient)

	authHandler := handlers.NewAuthHandler(userService)
	conversationHandler := handlers.NewConversationHandler(conversationService)
	chatHandler := handlers.NewChatHandler(chatService)
	skillHandler := handlers.NewSkillHandler(skillService, skillWorker)

	r := gin.Default()

	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authProtected.GET("/me", authHandler.Me)
		}

		conversations := v1.Group("/conversations")
		conversations.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			conversations.GET("", conversationHandler.List)
			conversations.POST("", conversationHandler.Create)
			conversations.GET("/:id", conversationHandler.Get)
			conversations.DELETE("/:id", conversationHandler.Delete)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware())
		{
			admin.GET("/users", authHandler.ListUsers)
			admin.DELETE("/users/:id", authHandler.DeleteUser)
		}

		chat := v1.Group("/chat")
		chat.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			chat.POST("", chatHandler.Chat)
		}

		skills := v1.Group("/skills")
		skills.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			skills.POST("/upload", skillHandler.Upload)
			skills.GET("", skillHandler.List)
			skills.GET("/:id", skillHandler.Get)
			skills.DELETE("/:id", skillHandler.Delete)
			skills.GET("/:id/status", skillHandler.Status)
			skills.GET("/:id/status/sse", skillHandler.SSEStatus)
			skills.GET("/:id/download", skillHandler.Download)
		}
	}

	log.Println("Server starting on :8080")
	graceful.Run(r, ":8080")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
