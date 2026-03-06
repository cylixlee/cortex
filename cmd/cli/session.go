package main

import (
	"context"
	"sync"
	"time"

	"github.com/cylixlee/cortex/pkg/session"
)

type sessionStore struct {
	sessions map[string]*session.Session
	mu       sync.RWMutex
}

func newSessionStore() *sessionStore {
	s := &sessionStore{
		sessions: make(map[string]*session.Session),
	}
	s.sessions["default"] = &session.Session{
		ID:        "default",
		UserID:    "default",
		Messages:  []session.Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s
}

func (s *sessionStore) Create(ctx context.Context, sess *session.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
	return nil
}

func (s *sessionStore) Get(ctx context.Context, id string) (*session.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if sess, ok := s.sessions[id]; ok {
		return sess, nil
	}
	return nil, session.ErrSessionNotFound
}

func (s *sessionStore) Update(ctx context.Context, sess *session.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
	return nil
}

func (s *sessionStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

func (s *sessionStore) ListByUser(ctx context.Context, userID string) ([]*session.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*session.Session
	for _, sess := range s.sessions {
		if sess.UserID == userID {
			result = append(result, sess)
		}
	}
	return result, nil
}
