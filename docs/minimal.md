# Minimal Core API Design

This document describes the Minimal Core API for Cortex, focusing on session management and remote model calling only.

## 1. Overview

The Minimal Core API provides a vendor-neutral abstraction layer for:
- **Session Management**: Multi-turn conversation persistence
- **LLM Client**: Streaming chat with remote models

Key design principles:
- **Streaming only**: All LLM interactions use streaming
- **Pure abstraction**: Internal formats are vendor-neutral; adapters handle OpenAI/Anthropic/etc. conversion
- **UserInterface pattern**: Separates rendering logic from agent logic

## 2. Data Model

### Message

```go
type Message struct {
    Role    string // "user" | "assistant" | "system"
    Content string
}
```

| Field   | Type   | Description                           |
| ------- | ------ | ------------------------------------- |
| Role    | string | Message sender: user/assistant/system |
| Content | string | Message content                       |

### Session

```go
type Session struct {
    ID        string
    UserID    string
    Messages  []Message
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

| Field     | Type      | Description                |
| --------- | --------- | -------------------------- |
| ID        | string    | Unique session identifier  |
| UserID    | string    | Owner user identifier      |
| Messages  | []Message | Conversation history       |
| CreatedAt | time      | Session creation timestamp |
| UpdatedAt | time      | Last message timestamp     |


### SessionStore

```go
type SessionStore interface {
    Create(ctx context.Context, s *Session) error
    Get(ctx context.Context, id string) (*Session, error)
    Update(ctx context.Context, s *Session) error
    Delete(ctx context.Context, id string) error
    ListByUser(ctx context.Context, userID string) ([]*Session, error)
}
```

**Design Rationale**:
- Interface-based design allows multiple implementations (in-memory, Redis, SQLite, PostgreSQL)
- `ListByUser` enables listing all sessions for a specific user
- All operations accept `context.Context` for cancellation/timeout support

## 3. LLM Client

### ChatRequest

```go
type ChatRequest struct {
    Model       string
    Messages    []Message
    Temperature float64
    MaxTokens   int
}
```

| Field       | Type      | Description                                  |
| ----------- | --------- | -------------------------------------------- |
| Model       | string    | Model identifier (e.g., "gpt-4", "claude-3") |
| Messages    | []Message | Conversation history                         |
| Temperature | float64   | Sampling temperature (0.0 - 2.0)             |
| MaxTokens   | int       | Maximum tokens to generate                   |

### ChatResponse

```go
type ChatResponse struct {
    Delta   string // incremental content
    IsFinal bool   // true if this is the last chunk
    Usage   Usage  // token usage (only valid when IsFinal is true)
}
```

| Field   | Type   | Description                                         |
| ------- | ------ | --------------------------------------------------- |
| Delta   | string | Incremental content for this chunk                  |
| IsFinal | bool   | Whether this is the final chunk                     |
| Usage   | Usage  | Token usage statistics (meaningful only when final) |

### Usage

```go
type Usage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}
```

| Field            | Type | Description              |
| ---------------- | ---- | ------------------------ |
| PromptTokens     | int  | Tokens in the prompt     |
| CompletionTokens | int  | Tokens in the completion |
| TotalTokens      | int  | Total tokens used        |

**Why `Usage` only on final chunk**: Token counts are only available after the entire response is generated. Streaming responses don't include usage until the stream ends.

### Client

```go
type Client interface {
    Stream(ctx context.Context, req ChatRequest) (<-chan ChatResponse, <-chan error)
}
```

**Design Rationale**:
- Channel-based streaming aligns with Go concurrency patterns
- Caller can use `select` to handle both response and error channels simultaneously
- Error returned from function indicates channel creation failure; stream errors are sent via error channel

## 4. Agent

### UserInterface

```go
type UserInterface interface {
    OnResponse(resp ChatResponse)
    OnError(err error)
}
```

| Method     | Description                        |
| ---------- | ---------------------------------- |
| OnResponse | Called for each streaming response |
| OnError    | Called when an error occurs        |

**Why this pattern**:
- Separates rendering logic from agent logic
- TUI implementations can do typewriter effects
- Web implementations can push SSE events
- Allows external control over throttling, interruption, etc.
- No need to switch UI within a single conversation

### Agent

```go
type Agent struct {
    llmClient    Client
    sessionStore SessionStore
    ui           UserInterface
}

func NewAgent(llmClient Client, sessionStore SessionStore, ui UserInterface) *Agent

func (a *Agent) Chat(ctx context.Context, sessionID, userInput string)
```

**Workflow**:
1. Load session from store (or create new if doesn't exist)
2. Append user message to session messages
3. Call LLM client with messages
4. Stream responses back via `ui.OnResponse()`
5. Append assistant response to session
6. Update session in store

## 5. Directory Structure

```
pkg/
├── session/
│   ├── session.go    # Session, Message structs + SessionStore interface
│   └── memory.go     # In-memory implementation (for testing/dev)
├── llm/
│   ├── client.go     # ChatRequest, ChatResponse, Usage + Client interface
│   └── openai.go     # OpenAI Compatible implementation
└── agent/
    └── agent.go      # Agent core scheduling logic + extensible interfaces
```

**Design Rationale**:
- `agent` package contains core scheduling logic
- Only exposes extensible interfaces (Client, SessionStore, UserInterface) for different scenarios
- Swap Client to switch models (OpenAI, Anthropic, local models)
- Swap SessionStore to switch databases (in-memory, SQLite, PostgreSQL)
- Swap UserInterface to switch UI (TUI with typewriter effect, Web with SSE)

## 6. Future Extensions

The following features are out of scope for Minimal Core API but planned for future versions:

- **Tool Calling**: Support for function/tool execution during conversation
- **System Prompt Management**: Dynamic system prompt updates
- **MCP Integration**: Model Context Protocol support
- **Skills**: Custom skill registration and execution
- **Sub-agents**: Agent-to-agent collaboration
