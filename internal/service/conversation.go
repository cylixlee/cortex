package service

import (
	"errors"

	"github.com/cylixlee/cortex/internal/models"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrAccessDenied         = errors.New("access denied")
)

type ConversationService struct {
	conversationRepo *repository.ConversationRepository
	messageRepo      *repository.MessageRepository
}

func NewConversationService(conversationRepo *repository.ConversationRepository, messageRepo *repository.MessageRepository) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

type ConversationOutput struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type MessageOutput struct {
	ID        uuid.UUID          `json:"id"`
	Role      models.MessageRole `json:"role"`
	Content   string             `json:"content"`
	CreatedAt string             `json:"created_at"`
}

type ConversationDetailOutput struct {
	ID        uuid.UUID       `json:"id"`
	Title     string          `json:"title"`
	Messages  []MessageOutput `json:"messages"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func (s *ConversationService) ListByUser(userID uuid.UUID, limit, offset int) ([]ConversationOutput, error) {
	conversations, err := s.conversationRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []ConversationOutput
	for _, conv := range conversations {
		result = append(result, ConversationOutput{
			ID:        conv.ID,
			Title:     conv.Title,
			CreatedAt: conv.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: conv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return result, nil
}

func (s *ConversationService) Create(userID uuid.UUID, title string) (*ConversationOutput, error) {
	conversation := &models.Conversation{
		UserID: userID,
		Title:  title,
	}

	if err := s.conversationRepo.Create(conversation); err != nil {
		return nil, err
	}

	return &ConversationOutput{
		ID:        conversation.ID,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: conversation.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *ConversationService) Get(userID, conversationID uuid.UUID) (*ConversationDetailOutput, error) {
	conversation, err := s.conversationRepo.FindByID(conversationID)
	if err != nil {
		return nil, ErrConversationNotFound
	}

	if conversation.UserID != userID {
		return nil, ErrAccessDenied
	}

	messages, err := s.messageRepo.FindByConversationID(conversationID)
	if err != nil {
		return nil, err
	}

	var messageOutputs []MessageOutput
	for _, msg := range messages {
		messageOutputs = append(messageOutputs, MessageOutput{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &ConversationDetailOutput{
		ID:        conversation.ID,
		Title:     conversation.Title,
		Messages:  messageOutputs,
		CreatedAt: conversation.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: conversation.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *ConversationService) Delete(userID, conversationID uuid.UUID) error {
	conversation, err := s.conversationRepo.FindByID(conversationID)
	if err != nil {
		return ErrConversationNotFound
	}

	if conversation.UserID != userID {
		return ErrAccessDenied
	}

	s.messageRepo.DeleteByConversationID(conversationID)
	return s.conversationRepo.Delete(conversationID)
}
