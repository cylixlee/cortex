package agent

import (
	"context"
	"strings"
	"time"

	"github.com/cylixlee/cortex/pkg/llm"
	"github.com/cylixlee/cortex/pkg/session"
)

// UserInterface defines the interface for handling agent responses and errors.
//
// This abstraction separates rendering logic from agent logic, enabling different UI implementations. For example, a
// TUI might use typewriter effects, while a web interface would push SSE events. The interface also allows external
// control over throttling, interruption, and other UI-specific behaviors.
type UserInterface interface {
	// OnResponse is called for each streaming response chunk from the LLM. The response may be a partial chunk
	// (IsFinal=false) or the final chunk (IsFinal=true) which includes usage statistics.
	OnResponse(resp llm.ChatResponse)

	// OnError is called when an error occurs during chat processing. This includes LLM API errors, session store
	// errors, and context errors.
	OnError(err error)
}

// Agent is the core struct that orchestrates conversations.
//
// It coordinates between an [llm.Client], a [session.SessionStore], and a [UserInterface] to handle the complete chat
// workflow. The agent loads conversation history from the store, sends messages to the LLM, streams responses back to
// the UI, and persists the updated conversation to the store.
type Agent struct {
	llmClient    llm.Client
	sessionStore session.SessionStore
	ui           UserInterface
}

// NewAgent creates a new Agent instance with the given dependencies.
//
// The llmClient handles communication with the language model. The sessionStore provides persistence for conversation
// history. The ui receives streaming responses and error notifications. All three dependencies can be swapped for
// different implementations (e.g., different LLM providers, storage backends).
func NewAgent(llmClient llm.Client, sessionStore session.SessionStore, ui UserInterface) *Agent {
	return &Agent{
		llmClient:    llmClient,
		sessionStore: sessionStore,
		ui:           ui,
	}
}

// Chat handles a chat interaction for the given session.
//
// This method loads the session from the store, appends the user's input as a new [session.Message], sends the
// conversation to the LLM for streaming response, delivers each response chunk via ui.OnResponse(), and finally saves
// the complete assistant response back to the session store. If the session does not exist, it calls ui.OnError with
// [session.ErrSessionNotFound]. The method returns when the stream completes, an error occurs, or the context is
// cancelled.
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
