package service

import (
	"context"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/cylixlee/cortex/internal/models"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/google/uuid"
)

type ChatService struct {
	conversationRepo *repository.ConversationRepository
	messageRepo      *repository.MessageRepository
	userRepo         *repository.UserRepository
	sessionManager   *llm.SessionManager
	embeddingClient  llm.Embedder
	chunkRepo        *repository.ChunkRepository
	client           *llm.Client
}

func NewChatService(
	conversationRepo *repository.ConversationRepository,
	messageRepo *repository.MessageRepository,
	userRepo *repository.UserRepository,
	client *llm.Client,
	embeddingClient llm.Embedder,
	chunkRepo *repository.ChunkRepository,
) *ChatService {
	return &ChatService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		userRepo:         userRepo,
		sessionManager:   llm.NewSessionManager(client),
		embeddingClient:  embeddingClient,
		chunkRepo:        chunkRepo,
		client:           client,
	}
}

type SendMessageInput struct {
	Message        string
	ConversationID string
}

type SendMessageOutput struct {
	ConversationID uuid.UUID
	Stream         func(ctx context.Context, message string, callback func(content string, err error) bool) error
}

func (s *ChatService) GetOrCreateConversation(userID uuid.UUID, conversationID *string) (*models.Conversation, error) {
	if conversationID != nil && *conversationID != "" {
		convID, err := uuid.Parse(*conversationID)
		if err == nil {
			conversation, err := s.conversationRepo.FindByID(convID)
			if err == nil && conversation.UserID == userID {
				return conversation, nil
			}
		}
	}

	title := "New Chat"
	return s.createNewConversation(userID, title)
}

func (s *ChatService) createNewConversation(userID uuid.UUID, title string) (*models.Conversation, error) {
	conversation := &models.Conversation{
		UserID: userID,
		Title:  title,
	}

	if err := s.conversationRepo.Create(conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *ChatService) GenerateTitle(userID uuid.UUID, firstMessage string) (*models.Conversation, error) {
	title := firstMessage
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", "")
	title = strings.TrimSpace(title)
	if len(title) > 50 {
		title = title[:50] + "..."
	}
	if title == "" {
		title = "New Chat"
	}
	return s.createNewConversation(userID, title)
}

func (s *ChatService) GenerateSmartTitle(ctx context.Context, conversationID uuid.UUID, firstUserMsg, firstAssistantMsg string) (string, error) {
	prompt := "请根据以下对话内容生成一个简短的中文标题（不超过30字）：\n\n用户：" + firstUserMsg + "\n\n助手：" + firstAssistantMsg

	messages := []*schema.Message{
		schema.SystemMessage("你是一个对话标题生成器。请根据用户和助手的对话内容生成一个简洁的中文标题，不超过30字，直接返回标题，不要有任何解释或引号。"),
		schema.UserMessage(prompt),
	}

	resp, err := s.client.GetChatModel().Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	title := resp.Content
	title = strings.ToValidUTF8(title, "")
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", "")
	title = strings.TrimSpace(title)
	if len(title) > 30 {
		title = title[:30]
	}

	if err := s.UpdateConversationTitle(conversationID, title); err != nil {
		return "", err
	}

	return title, nil
}

func (s *ChatService) UpdateConversationTitle(conversationID uuid.UUID, title string) error {
	conversation, err := s.conversationRepo.FindByID(conversationID)
	if err != nil {
		return err
	}

	conversation.Title = title
	return s.conversationRepo.Update(conversation)
}

func (s *ChatService) SaveUserMessage(conversationID uuid.UUID, content string) (*models.Message, error) {
	message := &models.Message{
		ConversationID: conversationID,
		Role:           models.MessageRoleUser,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) SaveAssistantMessage(conversationID uuid.UUID, content string) (*models.Message, error) {
	message := &models.Message{
		ConversationID: conversationID,
		Role:           models.MessageRoleAssistant,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ChatService) UpdateConversationTimestamp(conversationID uuid.UUID) error {
	conversation, err := s.conversationRepo.FindByID(conversationID)
	if err != nil {
		return err
	}

	conversation.UpdatedAt = time.Now()
	return s.conversationRepo.Update(conversation)
}

func (s *ChatService) GetSession(ctx context.Context, conversationID uuid.UUID) (*llm.Session, error) {
	sessionID := conversationID.String()
	return s.sessionManager.GetOrCreate(ctx, sessionID)
}

func (s *ChatService) RetrieveContext(ctx context.Context, query string, userID uuid.UUID, topK int) ([]string, error) {
	embeddings, err := s.embeddingClient.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, err
	}

	chunks, err := s.chunkRepo.SearchByEmbeddingForUser(embeddings[0], userID, topK)
	if err != nil {
		return nil, err
	}

	contexts := make([]string, len(chunks))
	for i, chunk := range chunks {
		contexts[i] = chunk.Content
	}
	return contexts, nil
}
