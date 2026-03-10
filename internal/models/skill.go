package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillStatus string

const (
	SkillStatusPending    SkillStatus = "pending"
	SkillStatusProcessing SkillStatus = "processing"
	SkillStatusCompleted  SkillStatus = "completed"
	SkillStatusFailed     SkillStatus = "failed"
)

type Skill struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Status       SkillStatus    `gorm:"type:varchar(20);default:pending" json:"status"`
	StoragePath  string         `gorm:"type:varchar(512)" json:"storage_path"`
	Progress     int            `gorm:"type:int;default:0" json:"progress"`
	ErrorMessage string         `gorm:"type:text" json:"error_message,omitempty"`
	Documents    []Document     `gorm:"foreignKey:SkillID" json:"documents,omitempty"`
	References   []Reference    `gorm:"foreignKey:SkillID" json:"references,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Skill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
