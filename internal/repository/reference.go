package repository

import (
	"github.com/open-portfolios/cortex/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReferenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) *ReferenceRepository {
	return &ReferenceRepository{db: db}
}

func (r *ReferenceRepository) Create(ref *models.Reference) error {
	return r.db.Create(ref).Error
}

func (r *ReferenceRepository) CreateBatch(refs []models.Reference) error {
	return r.db.Create(&refs).Error
}

func (r *ReferenceRepository) FindBySkillID(skillID uuid.UUID) ([]models.Reference, error) {
	var refs []models.Reference
	err := r.db.Where("skill_id = ?", skillID).Find(&refs).Error
	return refs, err
}

func (r *ReferenceRepository) DeleteBySkillID(skillID uuid.UUID) error {
	return r.db.Where("skill_id = ?", skillID).Delete(&models.Reference{}).Error
}
