package repository

import (
	"github.com/open-portfolios/cortex/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	vec := models.Float64ToVector(embedding)
	err := r.db.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "embedding <-> ?", Vars: []interface{}{vec}},
		}).
		Where("skill_id = ?", skillID).
		Limit(limit).
		Find(&chunks).Error
	return chunks, err
}

func (r *ChunkRepository) SearchByEmbeddingForUser(embedding []float64, userID uuid.UUID, limit int) ([]models.Chunk, error) {
	var skillIDs []uuid.UUID
	err := r.db.Model(&models.Skill{}).Where("user_id = ?", userID).Pluck("id", &skillIDs).Error
	if err != nil {
		return nil, err
	}

	if len(skillIDs) == 0 {
		return []models.Chunk{}, nil
	}

	var chunks []models.Chunk
	vec := models.Float64ToVector(embedding)
	err = r.db.
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "embedding <-> ?", Vars: []interface{}{vec}},
		}).
		Where("skill_id IN ?", skillIDs).
		Limit(limit).
		Find(&chunks).Error
	return chunks, err
}
