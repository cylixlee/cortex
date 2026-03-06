package session

import (
	"context"
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type Message struct {
	Role    string
	Content string
}

type Session struct {
	ID        string
	UserID    string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionStore interface {
	Create(ctx context.Context, s *Session) error
	Get(ctx context.Context, id string) (*Session, error)
	Update(ctx context.Context, s *Session) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]*Session, error)
}
