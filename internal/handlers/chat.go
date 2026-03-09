package handlers

import (
	"net/http"
	"strings"

	"github.com/cylixlee/cortex/internal/models"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type ChatHandler struct {
	sessionManager   *llm.SessionManager
	conversationRepo *repository.ConversationRepository
	messageRepo      *repository.MessageRepository
	userRepo         *repository.UserRepository
}

func NewChatHandler(client *llm.Client, conversationRepo *repository.ConversationRepository, messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *ChatHandler {
	return &ChatHandler{
		sessionManager:   llm.NewSessionManager(client),
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		userRepo:         userRepo,
	}
}

type ChatRequest struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var conversation *models.Conversation
	var err error

	if req.ConversationID != "" {
		convID, err := uuid.Parse(req.ConversationID)
		if err == nil {
			conversation, err = h.conversationRepo.FindByID(convID)
			if err != nil || conversation.UserID != userID {
				conversation = nil
			}
		}
	}

	if conversation == nil {
		title := req.Message
		title = strings.ReplaceAll(title, "\n", " ")
		title = strings.ReplaceAll(title, "\r", "")
		title = strings.TrimSpace(title)
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		if title == "" {
			title = "New Chat"
		}
		conversation = &models.Conversation{
			UserID: userID,
			Title:  title,
		}
		if err := h.conversationRepo.Create(conversation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
	}

	userMessage := &models.Message{
		ConversationID: conversation.ID,
		Role:           models.MessageRoleUser,
		Content:        req.Message,
		CreatedAt:      time.Now(),
	}
	if err := h.messageRepo.Create(userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user message"})
		return
	}

	sessionID := conversation.ID.String()
	session, err := h.sessionManager.GetOrCreate(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	c.Writer.Write([]byte("data: {\"conversation_id\": \"" + conversation.ID.String() + "\"}\n\n"))
	flusher.Flush()

	var assistantContent string

	err = session.Send(c.Request.Context(), req.Message, func(content string, err error) bool {
		if err != nil {
			return false
		}
		if content != "" {
			assistantContent += content
			c.Writer.Write([]byte("data: "))
			c.Writer.Write([]byte(content))
			c.Writer.Write([]byte("\n\n"))
			flusher.Flush()
		}
		return true
	})

	if assistantContent != "" {
		assistantMessage := &models.Message{
			ConversationID: conversation.ID,
			Role:           models.MessageRoleAssistant,
			Content:        assistantContent,
			CreatedAt:      time.Now(),
		}
		h.messageRepo.Create(assistantMessage)

		conversation.UpdatedAt = time.Now()
		h.conversationRepo.Update(conversation)
	}

	if err != nil {
		c.Writer.Write([]byte("data: [ERROR] "))
		c.Writer.Write([]byte(err.Error()))
		c.Writer.Write([]byte("\n\n"))
		flusher.Flush()
	}

	c.Writer.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
}
