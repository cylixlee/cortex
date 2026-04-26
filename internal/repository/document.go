package repository

import (
	"github.com/open-portfolios/cortex/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(doc *models.Document) error {
	return r.db.Create(doc).Error
}

func (r *DocumentRepository) CreateBatch(docs []models.Document) error {
	return r.db.Create(&docs).Error
}

func (r *DocumentRepository) FindBySkillID(skillID uuid.UUID) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("skill_id = ?", skillID).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) DeleteBySkillID(skillID uuid.UUID) error {
	return r.db.Where("skill_id = ?", skillID).Delete(&models.Document{}).Error
}
