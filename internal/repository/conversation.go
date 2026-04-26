package repository

import (
	"github.com/open-portfolios/cortex/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(conversation *models.Conversation) error {
	return r.db.Create(conversation).Error
}

func (r *ConversationRepository) FindByID(id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.First(&conversation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("updated_at DESC").Find(&conversations).Error
	return conversations, err
}

func (r *ConversationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Conversation{}, "id = ?", id).Error
}

func (r *ConversationRepository) Update(conversation *models.Conversation) error {
	return r.db.Save(conversation).Error
}
