package llm

import (
	"context"
	"io"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type Client struct {
	chatModel model.ToolCallingChatModel
}

func NewClient(baseURL, apiKey, modelName string) (*Client, error) {
	chatModel, err := deepseek.NewChatModel(context.Background(), &deepseek.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		chatModel: chatModel,
	}, nil
}

func (c *Client) StreamChat(ctx context.Context, prompt string) *adk.AsyncIterator[*adk.AgentEvent] {
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "chat_agent",
		Description: "A helpful assistant",
		Instruction: "You are a helpful assistant.",
		Model:       c.chatModel,
	})
	if err != nil {
		return nil
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	messages := []adk.Message{
		schema.UserMessage(prompt),
	}

	return runner.Run(ctx, messages)
}

type StreamHandler func(content string, err error) bool

func (c *Client) StreamWithHandler(ctx context.Context, prompt string, handler StreamHandler) error {
	iter := c.StreamChat(ctx, prompt)
	if iter == nil {
		return io.ErrUnexpectedEOF
	}

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			if stream := event.Output.MessageOutput.MessageStream; stream != nil {
				for {
					chunk, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						handler("", err)
						return err
					}
					if !handler(chunk.Content, nil) {
						return nil
					}
				}
			}
		}
	}
	return nil
}
