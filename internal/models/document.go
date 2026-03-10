package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SkillID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"skill_id"`
	Filename  string         `gorm:"type:varchar(512);not null" json:"filename"`
	Language  string         `gorm:"type:varchar(50)" json:"language"`
	Content   string         `gorm:"type:text" json:"content"`
	Size      int64          `gorm:"type:bigint" json:"size"`
	Chunks    []Chunk        `gorm:"foreignKey:DocumentID" json:"chunks,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
