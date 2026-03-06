package agent

import (
	"context"
	"strings"
	"time"

	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/cylixlee/cortex/pkg/session"
)

type UserInterface interface {
	OnResponse(resp llm.ChatResponse)
	OnError(err error)
}

type Agent struct {
	llmClient    llm.Client
	sessionStore session.SessionStore
	ui           UserInterface
}

func NewAgent(llmClient llm.Client, sessionStore session.SessionStore, ui UserInterface) *Agent {
	return &Agent{
		llmClient:    llmClient,
		sessionStore: sessionStore,
		ui:           ui,
	}
}

func (a *Agent) Chat(ctx context.Context, sessionID, userInput string) {
	sess, err := a.sessionStore.Get(ctx, sessionID)
	if err != nil || sess == nil {
		a.ui.OnError(session.ErrSessionNotFound)
		return
	}

	sess.Messages = append(sess.Messages, session.Message{
		Role:    "user",
		Content: userInput,
	})

	req := llm.ChatRequest{
		Messages: convertMessages(sess.Messages),
	}

	respChan, errChan := a.llmClient.Stream(ctx, req)

	var assistantContent strings.Builder

	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				return
			}
			a.ui.OnResponse(resp)
			assistantContent.WriteString(resp.Delta)
			if resp.IsFinal {
				sess.Messages = append(sess.Messages, session.Message{
					Role:    "assistant",
					Content: assistantContent.String(),
				})
				sess.UpdatedAt = time.Now()
				a.sessionStore.Update(ctx, sess)
			}
		case err, ok := <-errChan:
			if !ok {
				return
			}
			a.ui.OnError(err)
			sess.UpdatedAt = time.Now()
			a.sessionStore.Update(ctx, sess)
			return
		case <-ctx.Done():
			return
		}
	}
}

func convertMessages(msgs []session.Message) []llm.Message {
	result := make([]llm.Message, len(msgs))
	for i, msg := range msgs {
		result[i] = llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}
