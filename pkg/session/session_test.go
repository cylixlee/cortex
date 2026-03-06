package session

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockSessionStore struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

func NewMockSessionStore() *mockSessionStore {
	return &mockSessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *mockSessionStore) Create(ctx context.Context, sess *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
	return nil
}

func (s *mockSessionStore) Get(ctx context.Context, id string) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[id]; ok {
		return sess, nil
	}
	return nil, ErrSessionNotFound
}

func (s *mockSessionStore) Update(ctx context.Context, sess *Session) error {
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

func (s *mockSessionStore) ListByUser(ctx context.Context, userID string) ([]*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []*Session
	for _, sess := range s.sessions {
		if sess.UserID == userID {
			result = append(result, sess)
		}
	}
	return result, nil
}

func TestMessage_Fields(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello",
	}

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got '%s'", msg.Role)
	}
	if msg.Content != "Hello" {
		t.Errorf("expected content 'Hello', got '%s'", msg.Content)
	}
}

func TestSession_Fields(t *testing.T) {
	now := time.Now()
	sess := Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []Message{{Role: "user", Content: "Hi"}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if sess.ID != "session-1" {
		t.Errorf("expected ID 'session-1', got '%s'", sess.ID)
	}
	if sess.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", sess.UserID)
	}
	if len(sess.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(sess.Messages))
	}
	if sess.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if sess.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestMockSessionStore_Create(t *testing.T) {
	store := NewMockSessionStore()
	ctx := context.Background()
	now := time.Now()

	sess := &Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []Message{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := store.Create(ctx, sess)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMockSessionStore_Get(t *testing.T) {
	store := NewMockSessionStore()
	ctx := context.Background()
	now := time.Now()

	sess := &Session{
		ID:        "session-1",
		UserID:    "user-1",
		Messages:  []Message{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	store.Create(ctx, sess)

	got, err := store.Get(ctx, "session-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got.ID != "session-1" {
		t.Errorf("expected ID 'session-1', got '%s'", got.ID)
	}
}

func TestMockSessionStore_GetNotFound(t *testing.T) {
	store := NewMockSessionStore()
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("expected ErrSessionNotFound, got %v", err)
	}
}
