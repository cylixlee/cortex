package session

import (
	"context"
	"errors"
	"time"
)

// ErrSessionNotFound is returned when the requested session does not exist in the session store.
var ErrSessionNotFound = errors.New("session not found")

// Message represents a single message in a conversation.
//
// The Role field indicates who sent the message: "user", "assistant", or "system". The Content field contains the
// actual message text.
type Message struct {
	Role    string
	Content string
}

// Session represents a conversation session with multiple messages.
//
// A Session belongs to a user (identified by UserID) and maintains a history of all messages exchanged. The ID field
// uniquely identifies the session. CreatedAt and UpdatedAt timestamps track session lifecycle.
type Session struct {
	ID        string
	UserID    string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionStore defines the interface for session storage operations.
//
// Implementations can provide various storage backends such as in-memory maps for testing, Redis for caching,
// PostgreSQL for persistence, or any other storage mechanism. All methods accept context for cancellation and timeout
// control.
type SessionStore interface {
	// Create stores a new session in the store. The session ID should be unique and set by the caller.
	Create(ctx context.Context, s *Session) error

	// Get retrieves a session by its ID. Returns [ErrSessionNotFound] if the session does not exist.
	Get(ctx context.Context, id string) (*Session, error)

	// Update modifies an existing session. The session must already exist in the store.
	Update(ctx context.Context, s *Session) error

	// Delete removes a session from the store by its ID.
	Delete(ctx context.Context, id string) error

	// ListByUser returns all sessions belonging to a specific user, sorted by UpdatedAt in descending order (most
	// recent first).
	ListByUser(ctx context.Context, userID string) ([]*Session, error)
}
