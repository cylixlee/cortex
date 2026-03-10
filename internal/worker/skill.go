package worker

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cylixlee/cortex/internal/config"
	"github.com/cylixlee/cortex/internal/models"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/cylixlee/cortex/internal/service"
	"github.com/cylixlee/cortex/internal/workflows"
	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/cylixlee/cortex/pkg/storage"
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

type TaskMessage struct {
	SkillID  string `json:"skill_id"`
	UserID   string `json:"user_id"`
	Filename string `json:"filename"`
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

func (w *SkillWorker) InitTaskStatus(skillID string) {
	ctx := context.Background()
	id, err := uuid.Parse(skillID)
	if err != nil {
		return
	}
	w.publishStatus(ctx, id, "processing", 0)
}

func (w *SkillWorker) processTask(ctx context.Context, taskKey string) {
	log.Printf("Processing task: %s", taskKey)

	skillID, err := uuid.Parse(taskKey)
	if err != nil {
		log.Printf("Failed to parse skill ID: %v", err)
		return
	}

	w.publishStatus(ctx, skillID, "processing", 0)

	skill, err := w.skillRepo.FindByID(skillID)
	if err != nil {
		log.Printf("Failed to find skill: %v", err)
		w.publishStatus(ctx, skillID, "failed", 0)
		return
	}

	if err := w.skillService.UpdateSkillStatus(ctx, skillID, models.SkillStatusProcessing, 10); err != nil {
		log.Printf("Failed to update status: %v", err)
		return
	}

	obj, err := w.minioClient.Client().GetObject(ctx, w.minioClient.Bucket(), skill.StoragePath, minio.GetObjectOptions{})
	if err != nil {
		w.skillService.UpdateSkillStatusWithError(ctx, skillID, models.SkillStatusFailed, 0, err.Error())
		return
	}

	objBytes, err := io.ReadAll(obj)
	if err != nil {
		w.skillService.UpdateSkillStatusWithError(ctx, skillID, models.SkillStatusFailed, 0, err.Error())
		return
	}

	docs, err := w.skillService.ProcessSkill(ctx, skillID, objBytes)
	if err != nil {
		log.Printf("Failed to process skill: %v", err)
		return
	}

	if err := w.skillService.UpdateSkillStatus(ctx, skillID, models.SkillStatusProcessing, 50); err != nil {
		log.Printf("Failed to update status: %v", err)
		return
	}

	codeContext := buildCodeContext(docs)

	result, err := w.workflow.Run(ctx, codeContext)
	if err != nil {
		w.skillService.UpdateSkillStatusWithError(ctx, skillID, models.SkillStatusFailed, 60, err.Error())
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
	if err := w.createSkillPackage(ctx, skillID, result); err != nil {
		log.Printf("Failed to create skill package: %v", err)
		w.skillService.UpdateSkillStatusWithError(ctx, skillID, models.SkillStatusFailed, 90, err.Error())
		return
	}
	log.Printf("Skill package created successfully for skillID=%s", skillID.String())

	if err := w.skillService.UpdateSkillStatus(ctx, skillID, models.SkillStatusCompleted, 100); err != nil {
		log.Printf("Failed to update status: %v", err)
		return
	}

	w.publishStatus(ctx, skillID, "completed", 100)
	log.Printf("Task completed: %s", taskKey)
}

func (w *SkillWorker) publishStatus(ctx context.Context, skillID uuid.UUID, status string, progress int) {
	key := fmt.Sprintf("skill:status:%s", skillID.String())
	data := map[string]interface{}{
		"status":     status,
		"progress":   progress,
		"updated_at": time.Now().Unix(),
	}
	jsonData, _ := json.Marshal(data)
	w.redisClient.Set(ctx, key, jsonData, 24*time.Hour)
}

func (w *SkillWorker) GetTaskStatus(ctx context.Context, skillID string) (string, int, error) {
	key := fmt.Sprintf("skill:status:%s", skillID)
	data, err := w.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", 0, err
	}

	var status map[string]interface{}
	json.Unmarshal([]byte(data), &status)

	progress := int(status["progress"].(float64))
	return status["status"].(string), progress, nil
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
