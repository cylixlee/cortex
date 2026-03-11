package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cylixlee/cortex/internal/config"
	"github.com/cylixlee/cortex/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var chatTimeout = 120 * time.Second

func InitChatTimeout(cfg *config.Config) {
	if cfg.ChatTimeout > 0 {
		chatTimeout = time.Duration(cfg.ChatTimeout) * time.Second
	}
}

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

type ChatRequest struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id"`
}

type SSEConversationStart struct {
	ConversationID string `json:"conversation_id"`
}

func writeSSE(w http.ResponseWriter, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("data: " + string(b) + "\n\n"))
	return err
}

func writeSSEDone(w http.ResponseWriter) {
	w.Write([]byte("data: [DONE]\n\n"))
}

func writeSSEError(w http.ResponseWriter, errMsg string) {
	w.Write([]byte("data: {\"error\": \"" + errMsg + "\"}\n\n"))
}

func (h *ChatHandler) Chat(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), chatTimeout)
	defer cancel()

	conversationID := req.ConversationID
	var conversation *service.ConversationOutput
	var err error

	if conversationID != "" {
		_, err := uuid.Parse(conversationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
			return
		}
		output, err := h.chatService.GetOrCreateConversation(userID, &conversationID)
		if err == nil {
			conversation = &service.ConversationOutput{
				ID:        output.ID,
				Title:     output.Title,
				CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt: output.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			}
		}
	}

	if conversation == nil {
		conv, err := h.chatService.GenerateTitle(userID, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
		conversation = &service.ConversationOutput{
			ID:        conv.ID,
			Title:     conv.Title,
			CreatedAt: conv.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: conv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	_, err = h.chatService.SaveUserMessage(conversation.ID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user message"})
		return
	}

	session, err := h.chatService.GetSession(ctx, conversation.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat session"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	writeSSE(c.Writer, SSEConversationStart{ConversationID: conversation.ID.String()})
	flusher.Flush()

	var assistantContent string

	err = session.Send(ctx, req.Message, func(content string, err error) bool {
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
		h.chatService.SaveAssistantMessage(conversation.ID, assistantContent)
		h.chatService.UpdateConversationTimestamp(conversation.ID)
	}

	if err != nil {
		writeSSEError(c.Writer, err.Error())
		flusher.Flush()
	}

	writeSSEDone(c.Writer)
	flusher.Flush()
}
