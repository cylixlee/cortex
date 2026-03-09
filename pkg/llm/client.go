package llm

import (
	"context"
	"io"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type Client struct {
	chatModel model.ToolCallingChatModel
}

type Session struct {
	client   *Client
	id       string
	messages []adk.Message
	runner   *adk.Runner
}

func NewClient(baseURL, apiKey, modelName string) (*Client, error) {
	chatModel, err := deepseek.NewChatModel(context.Background(), &deepseek.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		chatModel: chatModel,
	}, nil
}

func (c *Client) CreateSession(ctx context.Context, sessionID string) (*Session, error) {
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "chat_agent",
		Description: "A helpful assistant",
		Instruction: "You are a helpful assistant.",
		Model:       c.chatModel,
	})
	if err != nil {
		return nil, err
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	return &Session{
		client:   c,
		id:       sessionID,
		messages: make([]adk.Message, 0),
		runner:   runner,
	}, nil
}

func (s *Session) Send(ctx context.Context, userMsg string, handler func(content string, err error) bool) error {
	userMessage := schema.UserMessage(userMsg)
	s.messages = append(s.messages, userMessage)

	events := s.runner.Run(ctx, s.messages)

	var assistantContent string

	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			if stream := event.Output.MessageOutput.MessageStream; stream != nil {
				for {
					chunk, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						handler("", err)
						return err
					}
					assistantContent += chunk.Content
					if !handler(chunk.Content, nil) {
						return nil
					}
				}
			}
		}
	}

	if assistantContent != "" {
		s.messages = append(s.messages, schema.AssistantMessage(assistantContent, nil))
	}

	return nil
}

type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]*Session
	client   *Client
}

func NewSessionManager(client *Client) *SessionManager {
	return &SessionManager{
		client:   client,
		sessions: make(map[string]*Session),
	}
}

func (m *SessionManager) GetOrCreate(ctx context.Context, sessionID string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[sessionID]; exists {
		return session, nil
	}

	session, err := m.client.CreateSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	m.sessions[sessionID] = session
	return session, nil
}
