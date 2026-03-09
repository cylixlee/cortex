# Examples

This document describes common example scenarios for Eino ADK.

## Hello World

The simplest example showing basic Agent creation and execution.

```go
import (
    "context"
    "errors"
    "fmt"
    "io"
    "os"

    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/components/model"
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/utils"
    "github.com/cloudwego/eino/schema"
    "github.com/cloudwego/eino/components/model/openai"
)

func main() {
    ctx := context.Background()

    // 1. Initialize model
    model, _ := openai.NewChatModel(ctx, &openai.ChatModelConfig{
        APIKey: os.Getenv("OPENAI_API_KEY"),
        Model:  os.Getenv("OPENAI_MODEL"),
    })

    // 2. Create Agent
    agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "hello_agent",
        Description: "A friendly greeting assistant",
        Instruction: "You are a friendly assistant.",
        Model:       model,
    })

    // 3. Create Runner
    runner := adk.NewRunner(ctx, adk.RunnerConfig{
        Agent:           agent,
        EnableStreaming: true,
    })

    // 4. Execute
    events := runner.Run(ctx, []adk.Message{
        schema.UserMessage("Hello! Introduce yourself."),
    })

    // 5. Process events
    for {
        event, ok := events.Next()
        if !ok {
            break
        }
        if event.Output != nil && event.Output.MessageOutput != nil {
            if stream := event.Output.MessageOutput.MessageStream; stream != nil {
                for {
                    chunk, err := stream.Recv()
                    if errors.Is(err, io.EOF) {
                        break
                    }
                    if err != nil {
                        // handle error
                        break
                    }
                    fmt.Printf("Agent: %s\n", chunk.Content)
                }
            }
        }
    }
}
```

## ChatModel with Tools

Agent with tool calling capability.

```go
// Define a tool
type GetWeatherInput struct {
    City string `json:"city"`
}

getWeatherTool, err := utils.InferTool(ctx, "get_weather",
    "Get weather information for a city",
    func(ctx context.Context, in *GetWeatherInput) (string, error) {
        return fmt.Sprintf("Weather in %s: Sunny", in.City), nil
    })

// Create Agent with tools
agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "weather_agent",
    Description: "Agent that provides weather information",
    Instruction: "You are a weather assistant. Use get_weather tool when needed.",
    Model:       model,
    ToolsConfig: adk.ToolsConfig{
        Tools: []tool.BaseTool{getWeatherTool},
    },
})
```

## Workflow

### Sequential Workflow

Execute sub-agents in sequence.

```go
agent := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
    Name:        "ResearchAgent",
    Description: "A sequential research workflow",
    SubAgents: []adk.Agent{
        researchAgent,  // First: research topic
        writerAgent,    // Second: write summary
    },
})
```

### Parallel Workflow

Execute sub-agents concurrently.

```go
agent := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
    Name:        "AnalysisAgent",
    Description: "Parallel analysis",
    SubAgents: []adk.Agent{
        analysisAgent1,
        analysisAgent2,
        analysisAgent3,
    },
})
```

### Loop Workflow

Iterate until condition met.

```go
agent := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
    Name:          "ReflectiveAgent",
    SubAgents:      []adk.Agent{reflectAgent},
    MaxIterations:  3,
})
```

## Multi-Agent

### Supervisor Pattern

Supervisor coordinates multiple sub-agents.

```go
supervisor := supervisor.New(ctx, &supervisor.Config{
    Supervisor: supervisorAgent,
    SubAgents: []adk.Agent{
        researcherAgent,
        coderAgent,
        reviewerAgent,
    },
})
```

### Plan-Execute-Replan Pattern

```go
entryAgent := planexecute.New(ctx, &planexecute.Config{
    Planner:       planAgent,    // Creates plan
    Executor:      execAgent,    // Executes plan
    Replanner:     replanAgent,  // Revises if needed
    MaxIterations: 20,
})
```

### Layered Supervisor

Multiple levels of supervisors.

```go
// Layer 1: Individual task supervisors
devSup := supervisor.New(ctx, &supervisor.Config{Supervisor: devSupervisorAgent, SubAgents: []adk.Agent{...}})
testSup := supervisor.New(ctx, &supervisor.Config{Supervisor: testSupervisorAgent, SubAgents: []adk.Agent{...}})

// Layer 2: Top-level coordinator
mainSup := supervisor.New(ctx, &supervisor.Config{
    Supervisor: mainSupervisorAgent,
    SubAgents:  []adk.Agent{devSup, testSup},
})
```

## Agent Tool

Wrap an Agent as a Tool for another Agent.

```go
// Create sub-agent
subAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "researcher",
    Description: "Researches topics",
    Model:       model,
})

// Wrap as tool
researcherTool := adk.NewAgentTool(ctx, subAgent,
    adk.WithFullChatHistoryAsInput(),
)

// Use in parent agent
parentAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "coordinator",
    Instruction: "Use researcher tool for research tasks.",
    Model:       model,
    ToolsConfig: adk.ToolsConfig{
        Tools: []tool.BaseTool{researcherTool},
    },
})
```

## Human-in-the-Loop

### Type Definitions

```go
type ApprovalInfo struct {
    ToolName        string
    ArgumentsInJSON string
    Approved        bool
}

type ReviewInfo struct {
    Requester string
    Content   string
}

type FollowUpInfo struct {
    Question string
}
```

### Approval

Tool execution requires user approval.

```go
type ApprovableTool struct {
    tool.InvokableTool
}

func (t ApprovableTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
    wasInterrupted, _, storedArgs := tool.GetInterruptState[string](ctx)
    
    if !wasInterrupted {
        // First call: trigger interruption for approval
        return "", tool.StatefulInterrupt(ctx, &ApprovalInfo{
            ToolName:        "BookTicket",
            ArgumentsInJSON: args,
        }, args)
    }
    
    // Resume: process approval result
    isResumeTarget, hasData, data := tool.GetResumeContext[*ApprovalInfo](ctx)
    if hasData && data.Approved {
        return t.InvokableTool.InvokableRun(ctx, storedArgs, opts...)
    }
    return "Rejected by user", nil
}
```

### Follow-up

Agent requests user clarification.

```go
// In tool implementation
if missingInfo {
    return "", tool.StatefulInterrupt(ctx, &FollowUpInfo{
        Question: "What is your budget?",
    }, currentArgs)
}
```

## Session Values

Share data between tools in an Agent execution.

```go
// Tool A: Store value
func toolA(ctx context.Context, in *Input) (string, error) {
    adk.AddSessionValue(ctx, "user-name", in.Name)
    return in.Name, nil
}

// Tool B: Retrieve value
func toolB(ctx context.Context, in *Input) (string, error) {
    userName, ok := adk.GetSessionValue(ctx, "user-name")
    if !ok {
        return "User not found", nil
    }
    return fmt.Sprintf("Hello %s", userName), nil
}
```

## HTTP SSE Service

Stream Agent events as Server-Sent Events.

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           agent,
    EnableStreaming: true,
})

events := runner.Run(ctx, messages)

w.Header().Set("Content-Type", "text/event-stream")
for {
    event, ok := events.Next()
    if !ok {
        break
    }
    // Convert event to SSE format
    fmt.Fprintf(w, "data: %s\n\n", eventJSON)
    w.Flush()
}
```

## Checkpoint and Resume

Persist state for recovery.

```go
// Implement custom CheckPointStore
type myCheckPointStore struct {
    data map[string][]byte
    mu   sync.Mutex
}

func (m *myCheckPointStore) Get(ctx context.Context, key string) ([]byte, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    v, ok := m.data[key]
    if !ok {
        return nil, nil
    }
    return v, nil
}

func (m *myCheckPointStore) Set(ctx context.Context, key string, value []byte) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = value
    return nil
}

// Create store
store := &myCheckPointStore{data: make(map[string][]byte)}

// Create runner with checkpoint
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           agent,
    EnableStreaming: true,
    CheckPointStore: store,
})

// On interruption
// ... Agent triggers StatefulInterrupt ...

// Later: Resume
iter, err := runner.Resume(ctx, checkpointID)
```

## Deep Agent

Prebuilt agent with task planning and TODO management.

```go
deepAgent := deep.New(ctx, &deep.Config{
    ChatModel: model,
    SubAgents: []adk.Agent{writerAgent},
})
```
