package handlers

import (
	"net/http"

	"github.com/cylixlee/cortex/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler(conversationService *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

type CreateConversationRequest struct {
	Title string `json:"title" binding:"required"`
}

type ConversationResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
}

type ConversationDetailResponse struct {
	ID        uuid.UUID         `json:"id"`
	Title     string            `json:"title"`
	Messages  []MessageResponse `json:"messages"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

func (h *ConversationHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	conversations, err := h.conversationService.ListByUser(userID, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list conversations"})
		return
	}

	response := []ConversationResponse{}
	for _, conv := range conversations {
		response = append(response, ConversationResponse{
			ID:        conv.ID,
			Title:     conv.Title,
			CreatedAt: conv.CreatedAt,
			UpdatedAt: conv.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *ConversationHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.conversationService.Create(userID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	c.JSON(http.StatusCreated, ConversationResponse{
		ID:        output.ID,
		Title:     output.Title,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	})
}

func (h *ConversationHandler) Get(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	output, err := h.conversationService.Get(userID, id)
	if err != nil {
		switch err {
		case service.ErrConversationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		case service.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		}
		return
	}

	var messageResponses []MessageResponse
	for _, msg := range output.Messages {
		messageResponses = append(messageResponses, MessageResponse{
			ID:        msg.ID,
			Role:      string(msg.Role),
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, ConversationDetailResponse{
		ID:        output.ID,
		Title:     output.Title,
		Messages:  messageResponses,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	})
}

func (h *ConversationHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	err = h.conversationService.Delete(userID, id)
	if err != nil {
		switch err {
		case service.ErrConversationNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		case service.ErrAccessDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted"})
}
