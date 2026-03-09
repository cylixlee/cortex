package handlers

import (
	"net/http"

	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	llmClient *llm.Client
}

func NewChatHandler(llmClient *llm.Client) *ChatHandler {
	return &ChatHandler{
		llmClient: llmClient,
	}
}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	err := h.llmClient.StreamWithHandler(c.Request.Context(), req.Message, func(content string, err error) bool {
		if err != nil {
			return false
		}
		if content != "" {
			c.Writer.Write([]byte("data: "))
			c.Writer.Write([]byte(content))
			c.Writer.Write([]byte("\n\n"))
			flusher.Flush()
		}
		return true
	})

	if err != nil {
		c.Writer.Write([]byte("data: [ERROR] "))
		c.Writer.Write([]byte(err.Error()))
		c.Writer.Write([]byte("\n\n"))
		flusher.Flush()
	}

	c.Writer.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
}
