# Core Concepts

This document describes the core concepts of Eino ADK.

## Agent

The Agent is the core interface that defines an entity capable of processing inputs and generating outputs through tools.

```go
type Agent interface {
    Name(ctx context.Context) string
    Description(ctx context.Context) string
    Run(ctx context.Context, input *AgentInput, options ...AgentRunOption) *AsyncIterator[*AgentEvent]
}
```

- **Name**: Unique identifier for the Agent
- **Description**: Used for task distribution among multiple agents
- **Run**: Asynchronous execution, returns event stream

## ResumableAgent

Supports interruption and resume. Required for human-in-the-loop scenarios.

```go
type ResumableAgent interface {
    Agent
    Resume(ctx context.Context, info *ResumeInfo, opts ...AgentRunOption) *AsyncIterator[*AgentEvent]
}
```

## Runner

The execution engine that manages Agent lifecycle.

```go
type Runner struct {
    a               Agent
    enableStreaming bool
    store           CheckPointStore
}
```

**Key Methods**:
- `Run(ctx, messages)`: Execute with message input
- `Query(ctx, query)`: Convenience method with string query
- `Resume(ctx, checkpointID)`: Resume from interruption
- `ResumeWithParams(ctx, checkpointID, params)`: Resume with specific data

## AgentEvent

Events generated during Agent execution.

```go
type AgentEvent struct {
    AgentName string
    RunPath   []RunStep      // Execution path tracking
    Output    *AgentOutput   // Message output
    Action    *AgentAction   // Actions (Exit/Interrupt/Transfer)
    Err       error
}
```

### AgentAction

```go
type AgentAction struct {
    Exit            bool
    Interrupted     *InterruptInfo
    TransferToAgent *TransferToAgentAction
    BreakLoop      *BreakLoopAction
    CustomizedAction any
}
```

## Session

Manages session-level state and event history.

```go
type runSession struct {
    Values    map[string]any       // Key-value store
    Events    []*agentEventWrapper  // Event history
    LaneEvents *laneEvents          // Parallel lane events
}
```

**Key Methods**:
- `GetSessionValues(ctx)`: Get all session values
- `AddSessionValue(ctx, key, value)`: Set a value
- `GetSessionValue(ctx, key)`: Get a value

## ChatModelAgent

The most common Agent implementation based on LLM with tool calling.

```go
type ChatModelAgent struct {
    name        string
    description string
    instruction string
    model       model.ToolCallingChatModel
    toolsConfig ToolsConfig
    // ...
}
```

### ReAct Pattern

ChatModelAgent uses ReAct (Reasoning + Acting) pattern:

1. **ChatModel**: Generate tool calls based on input
2. **ToolNode**: Execute tools and return results
3. **Loop**: Continue until no more tool calls or exit condition

## FlowAgent

Manages sub-agents and event forwarding.

Key features:
- SetSubAgents: Configure sub-agents
- HistoryRewriter: Rewrite conversation history
- Event forwarding between agents

## Interrupt

Interruption mechanism for human-in-the-loop.

Types:
- **Basic Interrupt**: `adkInterrupt(ctx, message)`
- **Stateful Interrupt**: `StatefulInterrupt(ctx, message, state)`
- **Composite Interrupt**: `CompositeInterrupt(ctx, message, state, subSignals...)`

## AgentTool

Wraps an Agent as a Tool for use by other Agents.

```go
agentTool := adk.NewAgentTool(ctx, agent,
    adk.WithFullChatHistoryAsInput(),
)
```

## Checkpoint

State persistence for interruption recovery.

Required interface:
```go
type CheckPointStore interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte) error
}
```

## RunPath

Tracks execution path in multi-agent scenarios.

```go
type RunStep struct {
    Address Address
    Name    string
}
```

Used to identify which sub-agent generated an event in complex workflows.
