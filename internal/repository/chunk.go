package repository

import (
	"github.com/cylixlee/cortex/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChunkRepository struct {
	db *gorm.DB
}

func NewChunkRepository(db *gorm.DB) *ChunkRepository {
	return &ChunkRepository{db: db}
}

func (r *ChunkRepository) Create(chunk *models.Chunk) error {
	return r.db.Create(chunk).Error
}

func (r *ChunkRepository) CreateBatch(chunks []models.Chunk) error {
	return r.db.Create(&chunks).Error
}

func (r *ChunkRepository) FindBySkillID(skillID uuid.UUID) ([]models.Chunk, error) {
	var chunks []models.Chunk
	err := r.db.Where("skill_id = ?", skillID).Order("document_id, chunk_index").Find(&chunks).Error
	return chunks, err
}

func (r *ChunkRepository) DeleteBySkillID(skillID uuid.UUID) error {
	return r.db.Where("skill_id = ?", skillID).Delete(&models.Chunk{}).Error
}

func (r *ChunkRepository) SearchByEmbedding(embedding []float64, skillID uuid.UUID, limit int) ([]models.Chunk, error) {
	var chunks []models.Chunk
	err := r.db.Where("skill_id = ?", skillID).
		Order("embedding <-> ?").
		Limit(limit).
		Find(&chunks, embedding).Error
	return chunks, err
}
