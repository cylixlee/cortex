package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/open-portfolios/cortex/internal/models"
	"github.com/open-portfolios/cortex/internal/service"
	"github.com/open-portfolios/cortex/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type SkillHandler struct {
	skillService *service.SkillService
	worker       *worker.SkillWorker
}

func NewSkillHandler(skillService *service.SkillService, worker *worker.SkillWorker) *SkillHandler {
	return &SkillHandler{
		skillService: skillService,
		worker:       worker,
	}
}

type UploadSkillResponse struct {
	SkillID string `json:"skill_id"`
	Status  string `json:"status"`
}

type SkillResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      string        `json:"status"`
	Stage       int           `json:"stage"`
	Skill       *SkillContent `json:"skill,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type SkillContent struct {
	Overview   string              `json:"overview"`
	References []ReferenceResponse `json:"references"`
}

type ReferenceResponse struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type SkillListResponse struct {
	Skills []SkillResponse `json:"skills"`
}

const maxFileSize = 100 << 20 // 100MB

func (h *SkillHandler) Upload(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 100MB limit"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	input := service.UploadSkillInput{
		UserID:      userID,
		Name:        name,
		FileContent: fileContent,
	}

	output, err := h.skillService.UploadSkill(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload skill"})
		return
	}

	h.worker.EnqueueTask(output.SkillID.String())

	c.JSON(http.StatusAccepted, UploadSkillResponse{
		SkillID: output.SkillID.String(),
		Status:  "processing",
	})
}

func (h *SkillHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	skills, err := h.skillService.GetUserSkills(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list skills"})
		return
	}

	response := SkillListResponse{
		Skills: make([]SkillResponse, 0, len(skills)),
	}
	for _, skill := range skills {
		response.Skills = append(response.Skills, SkillResponse{
			ID:        skill.ID,
			Name:      skill.Name,
			Status:    string(skill.Status),
			Stage:     int(skill.Stage),
			CreatedAt: skill.CreatedAt,
			UpdatedAt: skill.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *SkillHandler) Get(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	skill, err := h.skillService.GetSkill(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	if skill.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var refs []ReferenceResponse
	for _, ref := range skill.References {
		refs = append(refs, ReferenceResponse{
			Filename: ref.Filename,
			Content:  ref.Content,
		})
	}

	response := SkillResponse{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		Status:      string(skill.Status),
		Stage:       int(skill.Stage),
		CreatedAt:   skill.CreatedAt,
		UpdatedAt:   skill.UpdatedAt,
	}

	if skill.Status == models.SkillStatusCompleted {
		response.Skill = &SkillContent{
			Overview:   skill.Description,
			References: refs,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *SkillHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	skill, err := h.skillService.GetSkill(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	if skill.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = h.skillService.DeleteSkill(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted"})
}

func (h *SkillHandler) Status(c *gin.Context) {
	idStr := c.Param("id")

	stage, err := h.worker.GetTaskStage(c.Request.Context(), idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stage": int(stage),
		"name":  stage.String(),
	})
}

func (h *SkillHandler) SSEStatus(c *gin.Context) {
	idStr := c.Param("id")

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			stage, err := h.worker.GetTaskStage(c.Request.Context(), idStr)
			if err != nil {
				fmt.Fprintf(c.Writer, "data: {\"stage\":0,\"name\":\"error\"}\n\n")
				c.Writer.Flush()
				return
			}

			fmt.Fprintf(c.Writer, "data: {\"stage\":%d,\"name\":\"%s\"}\n\n", int(stage), stage.String())
			c.Writer.Flush()

			if stage == models.StageCompleted || stage == models.StageFailed {
				return
			}
		}
	}
}

func (h *SkillHandler) Download(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	skill, err := h.skillService.GetSkill(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	if skill.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if skill.Status != models.SkillStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill not ready for download"})
		return
	}

	obj, err := h.worker.GetMinioClient().Client().GetObject(c.Request.Context(), h.worker.GetMinioClient().Bucket(), skill.StoragePath, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file from storage"})
		return
	}
	defer obj.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", skill.Name))
	c.Header("Content-Type", "application/zip")

	_, err = io.Copy(c.Writer, obj)
	if err != nil {
		log.Printf("Failed to send file: %v", err)
	}
}
