package main

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/cylixlee/cortex/pkg/llm"
)

type llmClient struct{}

func newLLMClient() *llmClient {
	return &llmClient{}
}

func (c *llmClient) Stream(ctx context.Context, req llm.ChatRequest) (<-chan llm.ChatResponse, <-chan error) {
	respChan := make(chan llm.ChatResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)

		var userMessage string
		for i := len(req.Messages) - 1; i >= 0; i-- {
			if req.Messages[i].Role == "user" {
				userMessage = req.Messages[i].Content
				break
			}
		}

		if userMessage == "" {
			return
		}

		for word := range strings.SplitSeq(userMessage, "") {
			select {
			case <-ctx.Done():
				return
			default:
			}

			respChan <- llm.ChatResponse{
				Delta:   word,
				IsFinal: false,
			}

			delay := time.Duration(20+rand.Intn(60)) * time.Millisecond
			time.Sleep(delay)
		}

		respChan <- llm.ChatResponse{
			Delta:   "",
			IsFinal: true,
			Usage: llm.Usage{
				PromptTokens:     len(userMessage),
				CompletionTokens: len(userMessage),
				TotalTokens:      len(userMessage) * 2,
			},
		}
	}()

	return respChan, errChan
}
