package handlers

import (
	"net/http"

	"github.com/cylixlee/cortex/internal/models"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConversationHandler struct {
	conversationRepo *repository.ConversationRepository
	messageRepo      *repository.MessageRepository
}

func NewConversationHandler(conversationRepo *repository.ConversationRepository, messageRepo *repository.MessageRepository) *ConversationHandler {
	return &ConversationHandler{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
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
	ID        uuid.UUID          `json:"id"`
	Role      models.MessageRole `json:"role"`
	Content   string             `json:"content"`
	CreatedAt string             `json:"created_at"`
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

	conversations, err := h.conversationRepo.FindByUserID(userID, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list conversations"})
		return
	}

	var response []ConversationResponse
	for _, conv := range conversations {
		response = append(response, ConversationResponse{
			ID:        conv.ID,
			Title:     conv.Title,
			CreatedAt: conv.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: conv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
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

	conversation := &models.Conversation{
		UserID: userID,
		Title:  req.Title,
	}

	if err := h.conversationRepo.Create(conversation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	c.JSON(http.StatusCreated, ConversationResponse{
		ID:        conversation.ID,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: conversation.UpdatedAt.Format("2006-01-02T15:04:05Z"),
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

	conversation, err := h.conversationRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	if conversation.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	messages, err := h.messageRepo.FindByConversationID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, MessageResponse{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, ConversationDetailResponse{
		ID:        conversation.ID,
		Title:     conversation.Title,
		Messages:  messageResponses,
		CreatedAt: conversation.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: conversation.UpdatedAt.Format("2006-01-02T15:04:05Z"),
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

	conversation, err := h.conversationRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	if conversation.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	h.messageRepo.DeleteByConversationID(id)
	h.conversationRepo.Delete(id)

	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted"})
}
