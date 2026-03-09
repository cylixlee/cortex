package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type Message struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Role           MessageRole    `gorm:"type:varchar(20);not null" json:"role"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
