# Cortex - Agent Coding Guide

> This document defines the architectural vision and development guidelines for the Cortex project. All coding agents should follow these principles when implementing features.

## Project Vision

Cortex is an **AI-powered source code analysis platform** that transforms any codebase into a standardized **Agent Skill** package. It serves as both an intelligent code assistant and a skill factory.

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

### Core Workflow

1. **Upload & Analyze**: User uploads a source code ZIP → triggers multi-stage Agent orchestration
2. **Skill Generation**:
   - **Overview Agent**: Generates `SKILL.md` (standardized capability definition)
   - **API Retrieval Agent**: Generates detailed `references/` documentation
3. **Dual Output**:
   - **Internal (RAG)**: Vectorized into knowledge base for contextual Q&A
   - **External (Download)**: Package as downloadable Skill for downstream Agents

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                   Frontend (Vue)                    │
│        Chat UI | Upload | Skill Details | Download  │
└─────────────────────────────────────────────────────┘
                         │ HTTP/SSE
                         ▼
┌─────────────────────────────────────────────────────┐
│               API Gateway (Go + Gin)                │
│              Auth | Rate Limit | Routing            │
└─────────────────────────────────────────────────────┘
                         │
    ┌─────────────────────┼─────────────────────┐
    ▼                     ▼                     ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  Chat Service │  │   Code KB     │  │  User Service │
│    (Eino)     │  │   (Eino)      │  │               │
└───────────────┘  └───────────────┘  └───────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│            Eino Workflow Engine                     │
│   Orchestration | Multi-stage Analysis | Streaming │
└─────────────────────────────────────────────────────┘
                         │
    ┌─────────────────────┼─────────────────────┐
    ▼                     ▼                     ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│  PostgreSQL   │  │     Redis     │  │   Vector DB   │
│  + pgvector   │  │   (Cache)     │  │  (pgvector)   │
└───────────────┘  └───────────────┘  └───────────────┘
                         │
                         ▼
              ┌─────────────────────┐
              │   Object Storage    │
              │ (Skill Packages)    │
              └─────────────────────┘
```

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Technology Stack

| Layer        | Technology                              |
| ------------ | --------------------------------------- |
| Frontend     | Vue + TypeScript + SSE                  |
| Backend      | Go + Gin                                |
| AI Framework | **Eino** (Workflow Orchestration)       |
| Database     | PostgreSQL + pgvector (unified storage) |
| Cache        | Redis                                   |
| File Storage | Object Storage (MinIO/S3)               |
| Deployment   | Docker                                  |

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Core Modules

### 1. Frontend

- Chat interface with Markdown rendering and code highlighting
- Project upload with progress tracking
- Skill report detail page (`SKILL.md` + `references/`)
- One-click Skill package download (ZIP)
- SSE streaming response display

### 2. API Gateway (Go + Gin)

- Request routing
- JWT authentication
- Rate limiting
- Request/response logging

### 3. Chat Service

- Session management
- Context window maintenance
- Vector-based skill document retrieval
- Streaming message responses

### 4. Knowledge Base Service (Skill Factory) ⭐

Multi-stage workflow:

```
1. Source Upload (ZIP)
      ↓
2. [Stage 1: Overview Agent]
   → Generate SKILL.md (capability definition)
      ↓
3. [Stage 2: API Retrieval Agent] (depends on Stage 1)
   → Generate references/*.md
      ↓
4. [Dual Output]
   → Path A: Chunk → Embedding → Vector DB (RAG)
   → Path B: Package → Object Storage (Download)
      ↓
5. Completion Notification
```

### 5. User Service

- Registration/Login
- JWT token management
- User preferences

### 6. AI Service Layer

- Unified LLM invocation via Eino
- Streaming support
- Error retry logic

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Data Storage

| Storage        | Purpose                                     |
| -------------- | ------------------------------------------- |
| PostgreSQL     | Users, conversations, messages, vector data |
| Redis          | Session cache, rate limiting, task status   |
| Object Storage | Source ZIP files, generated Skill packages  |

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Core Workflows

### Chat Workflow

```
User Query → Vector Search (Skill Docs) → Build Prompt → LLM → Streaming Response
```

### Code-to-Skill Workflow

```
File Upload → Eino Multi-stage Analysis
  ├── [Overview Agent] → SKILL.md
  ├── [API Retrieval Agent] → references/*.md
  └── [Dual Output]
        ├── Path A: Chunk → Embedding → Vector DB
        └── Path B: Package → Object Storage
```

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Project Structure

```
cortex/
├── cmd/                    # Application entry points
├── internal/               # Internal business logic
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP handlers
│   ├── services/           # Business logic layer
│   ├── workflows/          # Eino workflow definitions
│   ├── models/             # Data models
│   ├── repository/         # Data access layer
│   └── middleware/         # Middleware
├── pkg/                    # Shared packages
│   ├── llm/                # LLM client wrapper (Eino)
│   ├── vector/             # Vector database client
│   └── storage/            # File storage utilities
├── web/                    # Frontend (Vue)
└── deployments/            # Deployment configs (Docker/K8s)
```

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Core Features

| Feature                | Description                                         |
| ---------------------- | --------------------------------------------------- |
| **Streaming Chat**     | SSE real-time push, ChatGPT-like experience         |
| **Skill Factory**      | Auto-generate standardized Agent Skill packages     |
| **Dual Output**        | Same analysis serves both RAG and external download |
| **Structured Reports** | Auto-generated SKILL.md + references/               |
| **Skill-based RAG**    | Highly structured retrieval for accurate answers    |
| **One-click Download** | Export Skill packages as ZIP                        |

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Development Phases

1. **Phase 1**: Basic chat functionality (Chat + LLM integration)
2. **Phase 2**: User authentication + conversation history
3. **Phase 3**: Skill generation workflow (Upload → SKILL.md/refs → storage)
4. **Phase 4**: Frontend skill details page + download functionality
5. **Phase 5**: Optimize Skill format for主流 Agent framework compatibility

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Design Principles

- **Skill Standardization**: `SKILL.md` must include clear capability descriptions, I/O specifications, and usage examples
- **Modularity**: Services are independent for easy extension
- **Streaming First**: All LLM responses support SSE streaming
- **Eino as Core**: Leverage Eino's workflow orchestration for multi-stage analysis
- **PostgreSQL Unified**: Relational + vector data in one database
- **Go Performance**: Leverage Go's concurrency for high-throughput
- **Knowledge as a Service (KaaS)**: Provide both Q&A and reusable knowledge assets

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.

## Important Notes for Coding Agents

1. **Do not implement all features at once** - follow the development phases
2. **Focus on core workflows first** - chat and LLM integration
3. **Use Eino for all AI-related orchestration** - not raw LLM calls
4. **Keep streaming support** - all user-facing LLM responses must be streamed
5. **Maintain skill package structure** - always generate both SKILL.md and references/
6. **PostgreSQL with pgvector** - use vector extension for semantic search
7. **Separation of concerns** - handlers → services → repository layers
8. **Use `task check` for validation** - you are not permitted to run any commands directly except `task check`

> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
> NOTE: You don't have any right to run command except `task check` defined in Taskfile in order to perform static analysis.
