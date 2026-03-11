package repository

import (
	"github.com/cylixlee/cortex/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) UpdateStoragePath(id uuid.UUID, path string) error {
	return r.db.Model(&models.Skill{}).Where("id = ?", id).Update("storage_path", path).Error
}

func (r *SkillRepository) Create(skill *models.Skill) error {
	return r.db.Create(skill).Error
}

func (r *SkillRepository) FindByID(id uuid.UUID) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.First(&skill, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *SkillRepository) FindByIDWithRelations(id uuid.UUID) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.Preload("Documents").Preload("References").First(&skill, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *SkillRepository) FindByUserID(userID uuid.UUID) ([]models.Skill, error) {
	var skills []models.Skill
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&skills).Error
	return skills, err
}

func (r *SkillRepository) Update(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *SkillRepository) UpdateStatus(id uuid.UUID, status models.SkillStatus, progress int) error {
	return r.db.Model(&models.Skill{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   status,
		"progress": progress,
	}).Error
}

func (r *SkillRepository) UpdateStatusWithError(id uuid.UUID, status models.SkillStatus, progress int, errMsg string) error {
	return r.db.Model(&models.Skill{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"progress":      progress,
		"error_message": errMsg,
	}).Error
}

func (r *SkillRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Skill{}, "id = ?", id).Error
}

func (r *SkillRepository) DeleteWithRelations(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("skill_id = ?", id).Delete(&models.Chunk{}).Error; err != nil {
			return err
		}

		if err := tx.Where("skill_id = ?", id).Delete(&models.Document{}).Error; err != nil {
			return err
		}

		if err := tx.Where("skill_id = ?", id).Delete(&models.Reference{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.Skill{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}
