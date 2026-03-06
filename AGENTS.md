# Cortex AI Agent - Development Guide

## Project Overview

Cortex is an AI Agent system with multi-turn conversation, tool calling (MCP, Skills), and sub-agent collaboration capabilities. It provides two interaction modes:
- **Web Service**: RESTful API / WebSocket via Gin with multi-user authentication and session persistence
- **CLI Tool**: Command-line interaction via Cobra with local configuration and session storage

## Project Structure

```
.
├── cmd
│   ├── cli                 # CLI program entry
│   └── server              # HTTP service entry
├── pkg
│   ├── agent               # Agent core scheduling logic
│   ├── session            # Session storage interface
│   ├── llm                # LLM client interface & implementations
│   ├── subagent           # Sub-agent scheduling
│   ├── mcp                # Model Communication Protocol abstraction
│   └── skills             # Skill registration and execution
├── internal
│   ├── config              # Configuration loading (Viper)
│   ├── db                  # Database operations (GORM)
│   ├── middleware          # HTTP middleware (Gin)
│   ├── repository          # Repository layer implementation
│   └── service             # Business services
├── api                     # API definitions (OpenAPI/Swagger)
├── docs                    # Documentation
├── web                     # Frontend static files
├── .devcontainer           # Dev container configuration
├── Taskfile.yml            # Task runner
├── go.mod
└── go.sum
```

## Tech Stack

| Area          | Technology            |
| ------------- | --------------------- |
| Web Framework | Gin                   |
| Database      | PostgreSQL + GORM     |
| Cache         | Redis                 |
| CLI           | Cobra                 |
| TUI           | Bubble Tea            |
| Config        | Viper                 |
| Logging       | Zap / Logrus          |
| Auth          | JWT + Redis blacklist |
| Streaming     | SSE or WebSocket      |
| Task Queue    | Asynq / channel       |
| Build Tool    | Task                  |
| Dev Env       | VSCode Dev Containers |

## Minimal Core API Design

See [docs/minimal.md](docs/minimal.md) for the Minimal Core API design, which focuses on session management and remote model calling only.

## Module Responsibilities

### pkg/agent
Provides core agent scheduling logic:
- Handles main flow orchestration
- Exposes extensible interfaces (Client, SessionStore, UserInterface) for different scenarios

### pkg/session, pkg/llm
- Session storage interface definition
- LLM client interface definition

### pkg/subagent, pkg/mcp, pkg/skills
- Sub-agent scheduling
- MCP abstraction
- Skill registration and execution

### cmd/server
- Initializes HTTP service, database, Redis
- Injects core dependencies
- Provides user authentication, session management, and chat endpoints

### cmd/cli
- Implements local commands
- Injects local storage (SQLite/file) session implementation
- Interfaces with core

### internal
Contains private implementations:
- Database repository layer
- Configuration loading
- HTTP middleware

## Development Guidelines

- Follow standard Go project layout
- Use dependency injection for testability
- Define interfaces in `pkg` for external implementations
- Keep business logic separate from infrastructure
- Configure development environment via .devcontainer
