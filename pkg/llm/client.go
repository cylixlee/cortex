package llm

import (
	"context"
)

// Message represents a chat message.
//
// The Role field indicates who sent the message: "user", "assistant", or "system". The Content field contains the
// actual message text. This type is provider-agnostic and can be mapped to any LLM API's message format.
type Message struct {
	Role    string
	Content string
}

// ChatRequest represents a request to the LLM chat API.
//
// This struct contains all parameters needed to make a chat completion request. The Model field specifies which model
// to use. Messages contains the conversation history. Temperature controls randomness (0.0-2.0), and MaxTokens limits
// the response length.
type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature float64
	MaxTokens   int
}

// ChatResponse represents a streaming response from the LLM.
//
// Responses are streamed incrementally. Delta contains the new content for this chunk. IsFinal is true for the last
// chunk, which also includes [Usage] statistics. The Usage field is only meaningful when IsFinal is true.
type ChatResponse struct {
	Delta   string
	IsFinal bool
	Usage   Usage
}

// Usage represents token usage statistics for a complete response.
//
// These values are only available after the entire response has been generated. PromptTokens counts input tokens,
// CompletionTokens counts output tokens, and TotalTokens is their sum.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Client defines the interface for LLM streaming chat.
//
// Implementations communicate with various LLM providers (OpenAI, Anthropic, local models, etc.). The Stream method
// returns channels for receiving streaming responses and errors, following Go concurrency patterns.
type Client interface {
	// Stream sends a chat request and returns channels for streaming responses.
	//
	// The respChan delivers [ChatResponse] chunks as they arrive. The errChan delivers any errors that occur during the
	// request. Both channels are closed when the stream ends. The caller should use select to receive from both
	// channels concurrently.
	Stream(ctx context.Context, req ChatRequest) (<-chan ChatResponse, <-chan error)
}
