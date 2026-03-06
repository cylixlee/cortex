package llm

import (
	"context"
)

type Message struct {
	Role    string
	Content string
}

type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature float64
	MaxTokens   int
}

type ChatResponse struct {
	Delta   string
	IsFinal bool
	Usage   Usage
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type Client interface {
	Stream(ctx context.Context, req ChatRequest) (<-chan ChatResponse, <-chan error)
}
