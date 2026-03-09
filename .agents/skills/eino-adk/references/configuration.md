# Configuration

This document describes configuration options for Eino ADK components.

## ChatModelAgentConfig

Configuration for ChatModelAgent.

| Field            | Type                       | Required | Default | Description                        |
| ---------------- | -------------------------- | -------- | ------- | ---------------------------------- |
| Name             | string                     | Yes      | -       | Agent name                         |
| Description      | string                     | Yes      | -       | Agent description for task routing |
| Instruction      | string                     | No       | -       | System prompt                      |
| Model            | model.ToolCallingChatModel | Yes      | -       | LLM with tool calling              |
| ToolsConfig      | ToolsConfig                | No       | -       | Tool configuration                 |
| GenModelInput    | GenModelInput              | No       | -       | Custom input generation            |
| Exit             | tool.BaseTool              | No       | -       | Custom exit tool                   |
| OutputKey        | string                     | No       | -       | Output storage key                 |
| MaxIterations    | int                        | No       | 20      | Maximum execution iterations       |
| Middlewares      | []AgentMiddleware          | No       | -       | Middleware list                    |
| ModelRetryConfig | *ModelRetryConfig          | No       | -       | Model retry configuration          |

## ToolsConfig

Tool configuration for Agent.

| Field              | Type            | Description                           |
| ------------------ | --------------- | ------------------------------------- |
| Tools              | []tool.BaseTool | Available tools                       |
| Strict             | bool            | Enable strict tool matching           |
| ReturnDirectly     | map[string]bool | Tools that trigger immediate return   |
| EmitInternalEvents | bool            | Emit internal events during execution |

### ReturnDirectly

Tools that return results directly without going back to the model:

```go
ToolsConfig: adk.ToolsConfig{
    Tools: []tool.BaseTool{searchTool, calcTool},
    ReturnDirectly: map[string]bool{
        "search": true,  // Returns directly after search
    },
}
```

## RunnerConfig

Configuration for Runner.

| Field           | Type            | Description               |
| --------------- | --------------- | ------------------------- |
| Agent           | Agent           | Agent to execute          |
| EnableStreaming | bool            | Enable streaming events   |
| CheckPointStore | CheckPointStore | For interruption recovery |

## AgentMiddleware

Middleware for ChatModelAgent.

| Field                 | Type                   | Description                  |
| --------------------- | ---------------------- | ---------------------------- |
| AdditionalInstruction | string                 | Append to system instruction |
| AdditionalTools       | []tool.BaseTool        | Additional tools available   |
| BeforeChatModel       | func(ctx, state) error | Hook before model call       |
| AfterChatModel        | func(ctx, state) error | Hook after model call        |
| WrapToolCall          | compose.ToolMiddleware | Wrap tool calls              |

### Example: Adding Instructions

```go
agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "my_agent",
    Model:       model,
    Middlewares: []adk.AgentMiddleware{
        {
            AdditionalInstruction: "Always double-check calculations.",
        },
    },
})
```

### Example: Custom Tool Wrap

```go
agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "my_agent",
    Model:       model,
    Middlewares: []adk.AgentMiddleware{
        {
            WrapToolCall: func(next compose.InvokableToolEndpoint) compose.InvokableToolEndpoint {
                return func(ctx context.Context, args string, opts ...tool.Option) (string, error) {
                    // Log before call
                    log.Printf("Calling tool with args: %s", args)
                    result, err := next(ctx, args, opts...)
                    // Log after call
                    log.Printf("Tool result: %s", result)
                    return result, err
                }
            },
        },
    },
})
```

## SequentialAgentConfig

Configuration for SequentialAgent.

| Field       | Type    | Description                    |
| ----------- | ------- | ------------------------------ |
| Name        | string  | Agent name                     |
| Description | string  | Agent description              |
| SubAgents   | []Agent | Sub-agents to execute in order |

## ParallelAgentConfig

Configuration for ParallelAgent.

| Field       | Type    | Description                       |
| ----------- | ------- | --------------------------------- |
| Name        | string  | Agent name                        |
| Description | string  | Agent description                 |
| SubAgents   | []Agent | Sub-agents to execute in parallel |

## LoopAgentConfig

Configuration for LoopAgent.

| Field          | Type    | Description             |
| -------------- | ------- | ----------------------- |
| Name           | string  | Agent name              |
| Description    | string  | Agent description       |
| SubAgents      | []Agent | Sub-agents to iterate   |
| MaxIterations  | int     | Maximum iterations      |
| BreakCondition | string  | Condition to break loop |

## Supervisor Config

Configuration for Supervisor prebuilt agent.

| Field       | Type                       | Description               |
| ----------- | -------------------------- | ------------------------- |
| Name        | string                     | Supervisor name           |
| Description | string                     | Supervisor description    |
| Agents      | []Agent                    | Sub-agents to supervise   |
| Model       | model.ToolCallingChatModel | Model for decision making |

## PlanExecute Config

Configuration for Plan-Execute-Replan agent.

| Field         | Type  | Description            |
| ------------- | ----- | ---------------------- |
| Planner       | Agent | Creates execution plan |
| Executor      | Agent | Executes plan steps    |
| Replanner     | Agent | Revises plan if needed |
| MaxIterations | int   | Maximum iterations     |
| MaxRetries    | int   | Maximum replan retries |

## ModelRetryConfig

Configuration for model call retries.

| Field           | Type          | Default | Description            |
| --------------- | ------------- | ------- | ---------------------- |
| MaxRetries      | int           | 3       | Maximum retry attempts |
| InitialInterval | time.Duration | 1s      | Initial retry interval |
| Multiplier      | float64       | 2.0     | Backoff multiplier     |
| MaxInterval     | time.Duration | 30s     | Maximum retry interval |

## AgentTool Options

Options for NewAgentTool.

| Option                       | Description                             |
| ---------------------------- | --------------------------------------- |
| WithFullChatHistoryAsInput() | Pass full conversation history to agent |
| WithAgentInputSchema(schema) | Use custom input schema                 |
| WithSessionValues(values)    | Pre-set session values                  |

## Call Options

Options for Agent execution.

| Option                   | Description         |
| ------------------------ | ------------------- |
| WithSessionValues(map)   | Pass session values |
| WithInput(input)         | Custom AgentInput   |
| WithCallbacks(callbacks) | Execution callbacks |
