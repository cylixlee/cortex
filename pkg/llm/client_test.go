package llm

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockLLMClient struct {
	responses []ChatResponse
	err       error
}

func (m *mockLLMClient) Stream(ctx context.Context, req ChatRequest) (<-chan ChatResponse, <-chan error) {
	respChan := make(chan ChatResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)
		for _, resp := range m.responses {
			select {
			case respChan <- resp:
			case <-ctx.Done():
				return
			}
		}
		if m.err != nil {
			errChan <- m.err
		}
	}()

	return respChan, errChan
}

func TestLLMClient_Stream_ReturnsChannels(t *testing.T) {
	client := &mockLLMClient{
		responses: []ChatResponse{
			{Delta: "Hello", IsFinal: false},
			{Delta: " world", IsFinal: true, Usage: Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15}},
		},
	}

	ctx := context.Background()
	respChan, errChan := client.Stream(ctx, ChatRequest{Model: "gpt-4"})

	if respChan == nil {
		t.Error("expected non-nil response channel")
	}
	if errChan == nil {
		t.Error("expected non-nil error channel")
	}
}

func TestLLMClient_Stream_ReceivesResponses(t *testing.T) {
	client := &mockLLMClient{
		responses: []ChatResponse{
			{Delta: "Hello", IsFinal: false},
			{Delta: " World", IsFinal: true, Usage: Usage{TotalTokens: 10}},
		},
	}

	ctx := context.Background()
	respChan, _ := client.Stream(ctx, ChatRequest{Model: "gpt-4"})

	var received []ChatResponse
	for resp := range respChan {
		received = append(received, resp)
	}

	if len(received) != 2 {
		t.Errorf("expected 2 responses, got %d", len(received))
	}
	if received[0].Delta != "Hello" {
		t.Errorf("expected first delta 'Hello', got '%s'", received[0].Delta)
	}
	if !received[1].IsFinal {
		t.Error("expected second response to be final")
	}
}

func TestLLMClient_Stream_Error(t *testing.T) {
	client := &mockLLMClient{
		responses: []ChatResponse{},
		err:       context.DeadlineExceeded,
	}

	ctx := context.Background()
	_, errChan := client.Stream(ctx, ChatRequest{Model: "gpt-4"})

	select {
	case err := <-errChan:
		if err != context.DeadlineExceeded {
			t.Errorf("expected context.DeadlineExceeded, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for error")
	}
}

func TestLLMClient_Stream_AggregatesContent(t *testing.T) {
	client := &mockLLMClient{
		responses: []ChatResponse{
			{Delta: "Hello", IsFinal: false},
			{Delta: " World", IsFinal: false},
			{Delta: "!", IsFinal: true, Usage: Usage{TotalTokens: 10}},
		},
	}

	ctx := context.Background()
	respChan, _ := client.Stream(ctx, ChatRequest{Model: "gpt-4"})

	var fullContent string
	for resp := range respChan {
		fullContent += resp.Delta
	}

	if fullContent != "Hello World!" {
		t.Errorf("expected 'Hello World!', got '%s'", fullContent)
	}
}

func TestLLMClient_Stream_EmptyResponses(t *testing.T) {
	client := &mockLLMClient{
		responses: []ChatResponse{},
	}

	ctx := context.Background()
	respChan, _ := client.Stream(ctx, ChatRequest{Model: "gpt-4"})

	resp, ok := <-respChan
	if ok {
		t.Error("expected channel to be closed")
	}
	_ = resp
}

func TestLLMClient_Stream_ErrorReturned(t *testing.T) {
	testErr := errors.New("test error")
	client := &mockLLMClient{
		responses: []ChatResponse{},
		err:       testErr,
	}

	ctx := context.Background()
	_, errChan := client.Stream(ctx, ChatRequest{Model: "gpt-4"})

	err := <-errChan
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}
