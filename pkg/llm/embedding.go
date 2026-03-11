package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	openai "github.com/cloudwego/eino-ext/components/embedding/openai"
)

type Embedder interface {
	EmbedStrings(ctx context.Context, texts []string) ([][]float64, error)
}

type EmbeddingClient struct {
	embedder Embedder
}

type SiliconFlowRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type SiliconFlowResponse struct {
	Object string `json:"object"`
	Model  string `json:"model"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type SiliconFlowEmbedder struct {
	client  *http.Client
	baseURL string
	apiKey  string
	model   string
}

type OpenAICompatEmbedder struct {
	embedder *openai.Embedder
}

func (o *OpenAICompatEmbedder) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	return (*o.embedder).EmbedStrings(ctx, texts)
}

func NewSiliconFlowEmbedder(baseURL, apiKey, model string) Embedder {
	return &SiliconFlowEmbedder{
		client:  &http.Client{},
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
	}
}

func (e *SiliconFlowEmbedder) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	reqBody := SiliconFlowRequest{
		Model: e.model,
		Input: texts,
	}

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", e.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SiliconFlowResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	embeddings := make([][]float64, len(result.Data))
	for i, d := range result.Data {
		embeddings[i] = d.Embedding
	}
	return embeddings, nil
}

func NewEmbeddingClient(ctx context.Context, provider, baseURL, apiKey, model string) (Embedder, error) {
	var embedder Embedder

	switch provider {
	case "openai":
		e, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
			APIKey:  apiKey,
			Model:   model,
			BaseURL: baseURL,
		})
		if err != nil {
			return nil, err
		}
		embedder = &OpenAICompatEmbedder{embedder: e}
	case "siliconflow":
		embedder = NewSiliconFlowEmbedder(baseURL, apiKey, model)
	default:
		return nil, errors.New("unsupported embedding provider: " + provider)
	}

	return &EmbeddingClient{
		embedder: embedder,
	}, nil
}

func (c *EmbeddingClient) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	return c.embedder.EmbedStrings(ctx, texts)
}
