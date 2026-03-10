package llm

import (
	"context"
	"errors"

	openai "github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
)

type Embedder interface {
	EmbedStrings(ctx context.Context, texts []string) ([][]float64, error)
}

type EmbeddingClient struct {
	embedder embedding.Embedder
}

func NewEmbeddingClient(ctx context.Context, provider, baseURL, apiKey, model string) (Embedder, error) {
	var embedder embedding.Embedder
	var err error

	switch provider {
	case "openai":
		embedder, err = openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
			APIKey: apiKey,
			Model:  model,
		})
	default:
		return nil, errors.New("unsupported embedding provider: " + provider)
	}

	if err != nil {
		return nil, err
	}

	return &EmbeddingClient{
		embedder: embedder,
	}, nil
}

func (c *EmbeddingClient) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	return c.embedder.EmbedStrings(ctx, texts)
}
