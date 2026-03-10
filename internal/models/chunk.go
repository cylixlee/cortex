package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chunk struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SkillID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"skill_id"`
	DocumentID uuid.UUID      `gorm:"type:uuid;not null;index" json:"document_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Embedding  []float64      `gorm:"type:jsonb" json:"embedding,omitempty"`
	ChunkIndex int            `gorm:"type:int;not null" json:"chunk_index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Chunk) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
