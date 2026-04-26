package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/open-portfolios/cortex/internal/models"
	"github.com/open-portfolios/cortex/internal/repository"
	"github.com/open-portfolios/cortex/pkg/llm"
	"github.com/open-portfolios/cortex/pkg/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type SkillService struct {
	skillRepo       *repository.SkillRepository
	documentRepo    *repository.DocumentRepository
	chunkRepo       *repository.ChunkRepository
	referenceRepo   *repository.ReferenceRepository
	minioClient     *storage.MinIOClient
	embeddingClient llm.Embedder
}

func NewSkillService(
	skillRepo *repository.SkillRepository,
	documentRepo *repository.DocumentRepository,
	chunkRepo *repository.ChunkRepository,
	referenceRepo *repository.ReferenceRepository,
	minioClient *storage.MinIOClient,
	embeddingClient llm.Embedder,
) *SkillService {
	return &SkillService{
		skillRepo:       skillRepo,
		documentRepo:    documentRepo,
		chunkRepo:       chunkRepo,
		referenceRepo:   referenceRepo,
		minioClient:     minioClient,
		embeddingClient: embeddingClient,
	}
}

type UploadSkillInput struct {
	UserID      uuid.UUID
	Name        string
	FileContent []byte
}

type UploadSkillOutput struct {
	SkillID uuid.UUID
}

func (s *SkillService) UploadSkill(ctx context.Context, input UploadSkillInput) (*UploadSkillOutput, error) {
	storagePath := fmt.Sprintf("skills/%s/%s.zip", input.UserID.String(), uuid.New().String())

	skill := &models.Skill{
		UserID:      input.UserID,
		Name:        input.Name,
		Status:      models.SkillStatusPending,
		StoragePath: storagePath,
	}

	if err := s.skillRepo.Create(skill); err != nil {
		return nil, err
	}

	_, err := s.minioClient.Client().PutObject(ctx, s.minioClient.Bucket(), storagePath, bytes.NewReader(input.FileContent), int64(len(input.FileContent)), minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &UploadSkillOutput{
		SkillID: skill.ID,
	}, nil
}

func (s *SkillService) GetSkill(ctx context.Context, skillID uuid.UUID) (*models.Skill, error) {
	return s.skillRepo.FindByIDWithRelations(skillID)
}

func (s *SkillService) GetUserSkills(ctx context.Context, userID uuid.UUID) ([]models.Skill, error) {
	return s.skillRepo.FindByUserID(userID)
}

func (s *SkillService) DeleteSkill(ctx context.Context, skillID uuid.UUID) error {
	skill, err := s.skillRepo.FindByID(skillID)
	if err != nil {
		return err
	}

	if err := s.skillRepo.DeleteWithRelations(skillID); err != nil {
		return err
	}

	if skill.StoragePath != "" {
		s.minioClient.Client().RemoveObject(ctx, s.minioClient.Bucket(), skill.StoragePath, minio.RemoveObjectOptions{})
	}

	return nil
}

func (s *SkillService) UpdateSkillStage(ctx context.Context, skillID uuid.UUID, stage models.SkillStage) error {
	return s.skillRepo.UpdateStage(skillID, stage)
}

func (s *SkillService) UpdateSkillStageWithError(ctx context.Context, skillID uuid.UUID, stage models.SkillStage, errMsg string) error {
	return s.skillRepo.UpdateStageWithError(skillID, stage, errMsg)
}

func (s *SkillService) ProcessSkill(ctx context.Context, skillID uuid.UUID, fileContent []byte) ([]models.Document, error) {
	if err := s.skillRepo.UpdateStage(skillID, models.StageExtracting); err != nil {
		return nil, err
	}

	docs, err := s.extractAndScanFiles(ctx, skillID, fileContent)
	if err != nil {
		s.skillRepo.UpdateStageWithError(skillID, models.StageFailed, err.Error())
		return nil, err
	}

	if err := s.generateEmbeddings(ctx, skillID, docs); err != nil {
		s.skillRepo.UpdateStageWithError(skillID, models.StageFailed, err.Error())
		return nil, err
	}

	return docs, nil
}

func (s *SkillService) extractAndScanFiles(ctx context.Context, skillID uuid.UUID, zipContent []byte) ([]models.Document, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
	if err != nil {
		return nil, fmt.Errorf("failed to read zip: %w", err)
	}

	var documents []models.Document

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name))
		if !isCodeFile(ext) {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			continue
		}
		defer rc.Close()

		content, err := io.ReadAll(rc)
		if err != nil {
			continue
		}

		doc := models.Document{
			SkillID:  skillID,
			Filename: file.Name,
			Language: getLanguage(ext),
			Content:  string(content),
			Size:     int64(len(content)),
		}

		if err := s.documentRepo.Create(&doc); err != nil {
			continue
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

func (s *SkillService) generateEmbeddings(ctx context.Context, skillID uuid.UUID, documents []models.Document) error {
	const chunkSize = 512
	const chunkOverlap = 50

	var allChunks []models.Chunk

	for _, doc := range documents {
		content := doc.Content
		for i := 0; i < len(content); i += chunkSize - chunkOverlap {
			end := i + chunkSize
			if end > len(content) {
				end = len(content)
			}

			chunkText := content[i:end]
			chunkIndex := len(allChunks)

			embeddings, err := s.embeddingClient.EmbedStrings(ctx, []string{chunkText})
			if err != nil {
				continue
			}

			chunk := models.Chunk{
				SkillID:    skillID,
				DocumentID: doc.ID,
				Content:    chunkText,
				Embedding:  models.Float64ToVector(embeddings[0]),
				ChunkIndex: chunkIndex,
			}

			allChunks = append(allChunks, chunk)
		}
	}

	if len(allChunks) > 0 {
		return s.chunkRepo.CreateBatch(allChunks)
	}

	return nil
}

func isCodeFile(ext string) bool {
	codeExts := map[string]bool{
		".go": true, ".java": true, ".python": true, ".js": true,
		".ts": true, ".jsx": true, ".tsx": true, ".c": true,
		".cpp": true, ".h": true, ".hpp": true, ".cs": true,
		".rb": true, ".php": true, ".swift": true, ".kt": true,
		".rs": true, ".vue": true, ".svelte": true,
	}
	return codeExts[ext]
}

func getLanguage(ext string) string {
	langMap := map[string]string{
		".go":     "go",
		".java":   "java",
		".py":     "python",
		".js":     "javascript",
		".ts":     "typescript",
		".jsx":    "javascript",
		".tsx":    "typescript",
		".c":      "c",
		".cpp":    "cpp",
		".h":      "c",
		".hpp":    "cpp",
		".cs":     "csharp",
		".rb":     "ruby",
		".php":    "php",
		".swift":  "swift",
		".kt":     "kotlin",
		".rs":     "rust",
		".vue":    "vue",
		".svelte": "svelte",
	}
	return langMap[ext]
}

func (s *SkillService) SaveReferences(ctx context.Context, skillID uuid.UUID, refs []models.Reference) error {
	for i := range refs {
		refs[i].SkillID = skillID
	}
	return s.referenceRepo.CreateBatch(refs)
}

func (s *SkillService) GetReferences(ctx context.Context, skillID uuid.UUID) ([]models.Reference, error) {
	return s.referenceRepo.FindBySkillID(skillID)
}
