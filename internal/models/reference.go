package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reference struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SkillID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"skill_id"`
	Filename  string         `gorm:"type:varchar(255);not null" json:"filename"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Reference) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
