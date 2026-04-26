package worker

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/open-portfolios/cortex/internal/config"
	"github.com/open-portfolios/cortex/internal/models"
	"github.com/open-portfolios/cortex/internal/repository"
	"github.com/open-portfolios/cortex/internal/service"
	"github.com/open-portfolios/cortex/internal/workflows"
	"github.com/open-portfolios/cortex/pkg/llm"
	"github.com/open-portfolios/cortex/pkg/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type SkillWorker struct {
	skillRepo     *repository.SkillRepository
	documentRepo  *repository.DocumentRepository
	chunkRepo     *repository.ChunkRepository
	referenceRepo *repository.ReferenceRepository
	skillService  *service.SkillService
	minioClient   *storage.MinIOClient
	chatClient    *llm.Client
	workflow      *workflows.SkillWorkflow
	redisClient   *redis.Client
	taskChannel   chan string
}

func NewSkillWorker(
	skillRepo *repository.SkillRepository,
	documentRepo *repository.DocumentRepository,
	chunkRepo *repository.ChunkRepository,
	referenceRepo *repository.ReferenceRepository,
	minioClient *storage.MinIOClient,
	chatClient *llm.Client,
	cfg *config.Config,
) (*SkillWorker, error) {
	embeddingClient, err := llm.NewEmbeddingClient(
		context.Background(),
		cfg.EmbeddingProvider,
		cfg.EmbeddingBaseURL,
		cfg.EmbeddingAPIKey,
		cfg.EmbeddingModel,
	)
	if err != nil {
		return nil, err
	}

	skillService := service.NewSkillService(
		skillRepo,
		documentRepo,
		chunkRepo,
		referenceRepo,
		minioClient,
		embeddingClient,
	)

	workflow, err := workflows.NewSkillWorkflow(chatClient.GetChatModel())
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(redisOpts)

	return &SkillWorker{
		skillRepo:     skillRepo,
		documentRepo:  documentRepo,
		chunkRepo:     chunkRepo,
		referenceRepo: referenceRepo,
		skillService:  skillService,
		minioClient:   minioClient,
		chatClient:    chatClient,
		workflow:      workflow,
		redisClient:   redisClient,
		taskChannel:   make(chan string, 100),
	}, nil
}

func (w *SkillWorker) Start(ctx context.Context) {
	log.Println("Skill worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Skill worker stopped")
			return
		case taskKey := <-w.taskChannel:
			w.processTask(ctx, taskKey)
		}
	}
}

func (w *SkillWorker) EnqueueTask(skillID string) {
	w.taskChannel <- skillID
}

func (w *SkillWorker) processTask(ctx context.Context, taskKey string) {
	log.Printf("Processing task: %s", taskKey)

	skillID, err := uuid.Parse(taskKey)
	if err != nil {
		log.Printf("Failed to parse skill ID: %v", err)
		return
	}

	w.publishStage(ctx, skillID, models.StageExtracting)

	skill, err := w.skillRepo.FindByID(skillID)
	if err != nil {
		log.Printf("Failed to find skill: %v", err)
		w.publishStage(ctx, skillID, models.StageFailed)
		return
	}

	obj, err := w.minioClient.Client().GetObject(ctx, w.minioClient.Bucket(), skill.StoragePath, minio.GetObjectOptions{})
	if err != nil {
		w.skillService.UpdateSkillStageWithError(ctx, skillID, models.StageFailed, err.Error())
		w.publishStage(ctx, skillID, models.StageFailed)
		return
	}
	defer obj.Close()

	objBytes, err := io.ReadAll(obj)
	if err != nil {
		w.skillService.UpdateSkillStageWithError(ctx, skillID, models.StageFailed, err.Error())
		w.publishStage(ctx, skillID, models.StageFailed)
		return
	}

	docs, err := w.skillService.ProcessSkill(ctx, skillID, objBytes)
	if err != nil {
		log.Printf("Failed to process skill: %v", err)
		w.publishStage(ctx, skillID, models.StageFailed)
		return
	}

	w.publishStage(ctx, skillID, models.StageAnalyzing)

	codeContext := buildCodeContext(docs)

	result, err := w.workflow.Run(ctx, codeContext)
	if err != nil {
		w.skillService.UpdateSkillStageWithError(ctx, skillID, models.StageFailed, err.Error())
		w.publishStage(ctx, skillID, models.StageFailed)
		return
	}

	for _, ref := range result.References {
		refModel := models.Reference{
			SkillID:  skillID,
			Filename: ref.Filename,
			Content:  ref.Content,
		}
		w.skillService.SaveReferences(ctx, skillID, []models.Reference{refModel})
	}

	skill.Description = result.Overview
	w.skillRepo.Update(skill)

	log.Printf("Creating skill package for skillID=%s", skillID.String())
	w.publishStage(ctx, skillID, models.StageGenerating)

	if err := w.createSkillPackage(ctx, skillID, result); err != nil {
		log.Printf("Failed to create skill package: %v", err)
		w.skillService.UpdateSkillStageWithError(ctx, skillID, models.StageFailed, err.Error())
		w.publishStage(ctx, skillID, models.StageFailed)
		return
	}
	log.Printf("Skill package created successfully for skillID=%s", skillID.String())

	skill, err = w.skillRepo.FindByID(skillID)
	if err != nil {
		log.Printf("Failed to find skill: %v", err)
		return
	}
	skill.Status = models.SkillStatusCompleted
	skill.Stage = models.StageCompleted
	if err := w.skillRepo.Update(skill); err != nil {
		log.Printf("Failed to update skill: %v", err)
		return
	}

	w.publishStage(ctx, skillID, models.StageCompleted)
	log.Printf("Task completed: %s", taskKey)
}

func (w *SkillWorker) publishStage(ctx context.Context, skillID uuid.UUID, stage models.SkillStage) {
	channel := fmt.Sprintf("skill:stage:%s", skillID.String())
	data := map[string]interface{}{
		"stage": int(stage),
		"name":  stage.String(),
	}
	jsonData, _ := json.Marshal(data)
	w.redisClient.Publish(ctx, channel, jsonData)
}

func (w *SkillWorker) GetTaskStage(ctx context.Context, skillID string) (models.SkillStage, error) {
	skill, err := w.skillRepo.FindByID(uuid.MustParse(skillID))
	if err != nil {
		return 0, err
	}
	return skill.Stage, nil
}

func buildCodeContext(docs []models.Document) string {
	var result string
	for _, doc := range docs {
		result += fmt.Sprintf("=== %s (%s) ===\n%s\n\n", doc.Filename, doc.Language, doc.Content)
	}
	return result
}

func (w *SkillWorker) GetMinioClient() *storage.MinIOClient {
	return w.minioClient
}

func (w *SkillWorker) createSkillPackage(ctx context.Context, skillID uuid.UUID, result *workflows.WorkflowResult) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	skillMD := "# SKILL.md\n\n" + result.Overview
	if err := addFileToZip(zipWriter, "SKILL.md", skillMD); err != nil {
		return err
	}

	for _, ref := range result.References {
		if err := addFileToZip(zipWriter, "references/"+ref.Filename, ref.Content); err != nil {
			return err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	packagePath := fmt.Sprintf("skills/%s/package.zip", skillID.String())
	_, err := w.minioClient.Client().PutObject(ctx, w.minioClient.Bucket(), packagePath, bytes.NewReader(buf.Bytes()), int64(buf.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	err = w.skillRepo.UpdateStoragePath(skillID, packagePath)
	if err != nil {
		return err
	}

	return nil
}

func addFileToZip(w *zip.Writer, filename, content string) error {
	f, err := w.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	return err
}
