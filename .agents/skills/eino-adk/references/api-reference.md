# API Reference

This document provides detailed API documentation for Eino ADK.

## Agent Creation

### NewChatModelAgent

Creates a ChatModel-based Agent.

```go
func NewChatModelAgent(ctx context.Context, config *ChatModelAgentConfig) (Agent, error)
```

**Config Fields**:
| Field         | Type                       | Required | Description                          |
| ------------- | -------------------------- | -------- | ------------------------------------ |
| Name          | string                     | Yes      | Agent name                           |
| Description   | string                     | Yes      | Agent description (for task routing) |
| Instruction   | string                     | No       | System prompt                        |
| Model         | model.ToolCallingChatModel | Yes      | LLM with tool calling                |
| ToolsConfig   | ToolsConfig                | No       | Tool configuration                   |
| MaxIterations | int                        | No       | Max iterations (default: 20)         |
| Middlewares   | []AgentMiddleware          | No       | Middleware list                      |
| Exit          | tool.BaseTool              | No       | Custom exit tool                     |

### NewSequentialAgent

Creates an Agent that executes sub-agents sequentially.

```go
func NewSequentialAgent(ctx context.Context, config *SequentialAgentConfig) (Agent, error)
```

**Config Fields**:
| Field     | Type    | Description                    |
| --------- | ------- | ------------------------------ |
| Name      | string  | Agent name                     |
| SubAgents | []Agent | Sub-agents to execute in order |

### NewParallelAgent

Creates an Agent that executes sub-agents in parallel.

```go
func NewParallelAgent(ctx context.Context, config *ParallelAgentConfig) (Agent, error)
```

### NewLoopAgent

Creates an Agent that loops through sub-agents.

```go
func NewLoopAgent(ctx context.Context, config *LoopAgentConfig) (Agent, error)
```

**Config Fields**:
| Field         | Type    | Description           |
| ------------- | ------- | --------------------- |
| SubAgents     | []Agent | Sub-agents to iterate |
| MaxIterations | int     | Maximum iterations    |

## Runner

### NewRunner

Creates a Runner for executing an Agent.

```go
func NewRunner(ctx context.Context, config RunnerConfig) *Runner
```

**Config Fields**:
| Field           | Type            | Description             |
| --------------- | --------------- | ----------------------- |
| Agent           | Agent           | Agent to execute        |
| EnableStreaming | bool            | Enable streaming events |
| CheckPointStore | CheckPointStore | Optional, for resume    |

### Runner.Run

Execute Agent with message input.

```go
func (r *Runner) Run(ctx context.Context, messages []schema.Message, options ...AgentRunOption) *AsyncIterator[*AgentEvent]
```

### Runner.Query

Convenience method for single query execution.

```go
func (r *Runner) Query(ctx context.Context, query string) *AsyncIterator[*AgentEvent]
```

### Runner.Resume

Resume from checkpoint.

```go
func (r *Runner) Resume(ctx context.Context, checkpointID string) (*AsyncIterator[*AgentEvent], error)
```

### Runner.ResumeWithParams

Resume with specific target data.

```go
func (r *Runner) ResumeWithParams(ctx context.Context, checkpointID string, params *ResumeParams) (*AsyncIterator[*AgentEvent], error)
```

## AgentTool

### NewAgentTool

Wraps an Agent as a Tool.

```go
func NewAgentTool(ctx context.Context, agent Agent, opts ...AgentToolOption) tool.BaseTool
```

**Options**:
| Option                       | Description                             |
| ---------------------------- | --------------------------------------- |
| WithFullChatHistoryAsInput() | Pass full conversation history to agent |
| WithAgentInputSchema(schema) | Custom input schema                     |

## Session Values

### AddSessionValue

Store a value in session.

```go
func AddSessionValue(ctx context.Context, key string, value any)
```

### GetSessionValue

Retrieve a value from session.

```go
func GetSessionValue(ctx context.Context, key string) (any, bool)
```

### GetSessionValues

Retrieve all session values.

```go
func GetSessionValues(ctx context.Context) map[string]any
```

## Call Options

### WithSessionValues

Pass session values to Agent execution.

```go
runner.Run(ctx, messages, adk.WithSessionValues(map[string]any{
    "user_id": "123",
}))
```

### WithSkipTransferMessages

Skip transfer messages in history.

```go
runner.Run(ctx, messages, adk.WithSkipTransferMessages())
```

## Interruption

### StatefulInterrupt

Trigger a stateful interruption.

```go
func StatefulInterrupt(ctx context.Context, info *InterruptInfo, state any) error
```

### GetInterruptState

Get interruption state.

```go
func GetInterruptState[T any](ctx context.Context) (bool, bool, T)
```

### GetResumeContext

Get resume context data.

```go
func GetResumeContext[T any](ctx context.Context) (bool, bool, T)
```

## SetSubAgents

Set sub-agents for a composite Agent.

```go
func SetSubAgents(ctx context.Context, agent Agent, subAgents []Agent) (Agent, error)
```

## Tool Creation Helpers

### utils.InferTool

Wrap a function as InvokableTool.

```go
import "github.com/cloudwego/eino/components/tool/utils"

tool, err := utils.InferTool(ctx, "tool_name", "description",
    func(ctx context.Context, input *MyInput) (string, error) {
        // implementation
        return result, nil
    })
```

## Events

### GetMessage

Extract message from AgentEvent.

```go
func GetMessage(e *AgentEvent) (Message, *AgentEvent, error)
```
