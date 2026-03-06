package agent

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/cylixlee/cortex/pkg/session"
)

type mockSessionStore struct {
	mu       sync.Mutex
	sessions map[string]*session.Session
}

func newMockSessionStore() *mockSessionStore {
	return &mockSessionStore{
		sessions: make(map[string]*session.Session),
	}
}

func (s *mockSessionStore) Create(ctx context.Context, sess *session.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
	return nil
}

func (s *mockSessionStore) Get(ctx context.Context, id string) (*session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[id]; ok {
		return sess, nil
	}
	return nil, session.ErrSessionNotFound
}

func (s *mockSessionStore) Update(ctx context.Context, sess *session.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
	return nil
}

func (s *mockSessionStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

func (s *mockSessionStore) ListByUser(ctx context.Context, userID string) ([]*session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []*session.Session
	for _, sess := range s.sessions {
		if sess.UserID == userID {
			result = append(result, sess)
		}
	}
	return result, nil
}

func TestNewAgent_CreatesInstance(t *testing.T) {
	mockLLM := &mockLLMClient{}
	mockStore := newMockSessionStore()
	mockUI := &mockUserInterface{}

	a := NewAgent(mockLLM, mockStore, mockUI)

	if a == nil {
		t.Error("expected non-nil Agent")
	}
	if a.llmClient != mockLLM {
		t.Error("expected llmClient to be set")
	}
	if a.sessionStore != mockStore {
		t.Error("expected sessionStore to be set")
	}
	if a.ui != mockUI {
		t.Error("expected ui to be set")
	}
}

func TestAgent_Chat_CallsOnError_WhenSessionNotExist(t *testing.T) {
	mockLLM := &mockLLMClient{}
	mockStore := newMockSessionStore()
	mockUI := &mockUserInterface{}

	a := NewAgent(mockLLM, mockStore, mockUI)
	ctx := context.Background()

	a.Chat(ctx, "non-existent-session", "Hello")

	if len(mockUI.errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(mockUI.errors))
	}
	if mockUI.errors[0] != session.ErrSessionNotFound {
		t.Errorf("expected ErrSessionNotFound, got %v", mockUI.errors[0])
	}
}

func TestAgent_Chat_AppendsUserMessage(t *testing.T) {
	mockLLM := &mockLLMClient{}
	mockStore := newMockSessionStore()
	mockUI := &mockUserInterface{}

	now := time.Now()
	mockStore.Create(context.Background(), &session.Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []session.Message{},
		CreatedAt: now,
		UpdatedAt: now,
	})

	a := NewAgent(mockLLM, mockStore, mockUI)
	ctx := context.Background()

	a.Chat(ctx, "session-1", "Hello")

	sess, _ := mockStore.Get(ctx, "session-1")
	if len(sess.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(sess.Messages))
	}
	if sess.Messages[0].Role != "user" {
		t.Errorf("expected role 'user', got '%s'", sess.Messages[0].Role)
	}
	if sess.Messages[0].Content != "Hello" {
		t.Errorf("expected content 'Hello', got '%s'", sess.Messages[0].Content)
	}
}

func TestAgent_Chat_AppendsAssistantMessage(t *testing.T) {
	mockLLM := &mockLLMClient{
		responses: []llm.ChatResponse{
			{Delta: "Hi there!", IsFinal: true},
		},
	}
	mockStore := newMockSessionStore()
	mockUI := &mockUserInterface{}

	now := time.Now()
	mockStore.Create(context.Background(), &session.Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []session.Message{},
		CreatedAt: now,
		UpdatedAt: now,
	})

	a := NewAgent(mockLLM, mockStore, mockUI)
	ctx := context.Background()

	a.Chat(ctx, "session-1", "Hello")

	sess, _ := mockStore.Get(ctx, "session-1")
	if len(sess.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(sess.Messages))
	}
	if sess.Messages[1].Role != "assistant" {
		t.Errorf("expected role 'assistant', got '%s'", sess.Messages[1].Role)
	}
	if sess.Messages[1].Content != "Hi there!" {
		t.Errorf("expected content 'Hi there!', got '%s'", sess.Messages[1].Content)
	}
}

func TestAgent_Chat_CallsOnResponse(t *testing.T) {
	mockLLM := &mockLLMClient{
		responses: []llm.ChatResponse{
			{Delta: "Hello", IsFinal: false},
			{Delta: " World", IsFinal: true, Usage: llm.Usage{TotalTokens: 10}},
		},
	}
	mockStore := newMockSessionStore()
	mockUI := &mockUserInterface{}

	now := time.Now()
	mockStore.Create(context.Background(), &session.Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []session.Message{},
		CreatedAt: now,
		UpdatedAt: now,
	})

	a := NewAgent(mockLLM, mockStore, mockUI)
	ctx := context.Background()

	a.Chat(ctx, "session-1", "Hello")

	if len(mockUI.responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(mockUI.responses))
	}
	if mockUI.responses[0].Delta != "Hello" {
		t.Errorf("expected first delta 'Hello', got '%s'", mockUI.responses[0].Delta)
	}
	if !mockUI.responses[1].IsFinal {
		t.Error("expected second response to be final")
	}
}

func TestAgent_Chat_CallsOnError_OnFailure(t *testing.T) {
	mockLLM := &mockLLMClient{
		err: context.DeadlineExceeded,
	}
	mockStore := newMockSessionStore()
	mockUI := &mockUserInterface{}

	now := time.Now()
	mockStore.Create(context.Background(), &session.Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []session.Message{},
		CreatedAt: now,
		UpdatedAt: now,
	})

	a := NewAgent(mockLLM, mockStore, mockUI)
	ctx := context.Background()

	a.Chat(ctx, "session-1", "Hello")

	if len(mockUI.errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(mockUI.errors))
	}
	if mockUI.errors[0] != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", mockUI.errors[0])
	}
}

type mockLLMClient struct {
	responses []llm.ChatResponse
	err       error
}

func (m *mockLLMClient) Stream(ctx context.Context, req llm.ChatRequest) (<-chan llm.ChatResponse, <-chan error) {
	respChan := make(chan llm.ChatResponse, 10)
	errChan := make(chan error, 10)

	go func() {
		defer close(respChan)
		defer close(errChan)
		for _, resp := range m.responses {
			select {
			case respChan <- resp:
			case <-ctx.Done():
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
		if m.err != nil {
			errChan <- m.err
		}
	}()

	return respChan, errChan
}

type mockUserInterface struct {
	mu        sync.Mutex
	responses []llm.ChatResponse
	errors    []error
}

func (m *mockUserInterface) OnResponse(resp llm.ChatResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = append(m.responses, resp)
}

func (m *mockUserInterface) OnError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors = append(m.errors, err)
}
