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

type SkillStage int

const (
	StagePending    SkillStage = 1 // 上传中
	StageExtracting SkillStage = 2 // 解析文件中
	StageAnalyzing  SkillStage = 3 // AI 分析中
	StageGenerating SkillStage = 4 // 生成文档中
	StageCompleted  SkillStage = 5 // 完成
	StageFailed     SkillStage = 6 // 失败
)

func (s SkillStage) String() string {
	switch s {
	case StagePending:
		return "pending"
	case StageExtracting:
		return "extracting"
	case StageAnalyzing:
		return "analyzing"
	case StageGenerating:
		return "generating"
	case StageCompleted:
		return "completed"
	case StageFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func (s SkillStage) Description() string {
	switch s {
	case StagePending:
		return "上传中"
	case StageExtracting:
		return "解析文件中"
	case StageAnalyzing:
		return "AI 分析中"
	case StageGenerating:
		return "生成文档中"
	case StageCompleted:
		return "已完成"
	case StageFailed:
		return "处理失败"
	default:
		return "未知状态"
	}
}

type Skill struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Overview     string         `gorm:"type:text" json:"overview"`
	Status       SkillStatus    `gorm:"type:varchar(20);default:pending" json:"status"`
	Stage        SkillStage     `gorm:"type:int;default:1" json:"stage"`
	StoragePath  string         `gorm:"type:varchar(512)" json:"storage_path"`
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
