// Package llm provides vendor-neutral LLM client interfaces and types.
//
// This package defines the Client interface and related types (ChatRequest, ChatResponse, Usage) for interacting with
// language models. It is designed to be provider-agnostic, allowing implementations for OpenAI, Anthropic, local
// models, or any other LLM API that supports streaming responses.
package llm
