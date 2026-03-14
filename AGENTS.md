# Cortex - Agent Coding Guide

> This document defines the architectural vision and development guidelines for the Cortex project.

## Project Vision

Cortex is an **AI-powered source code analysis platform** that transforms any codebase into a standardized **Agent Skill** package. It serves as both an intelligent code assistant and a skill factory.

### Core Workflow

1. **Upload & Analyze**: User uploads a source code ZIP → triggers multi-stage Agent orchestration
2. **Skill Generation**:
   - **Overview Agent**: Generates `SKILL.md` (standardized capability definition)
   - **API Retrieval Agent**: Generates detailed `references/` documentation
3. **Dual Output**:
   - **Internal (RAG)**: Vectorized into knowledge base for contextual Q&A
   - **External (Download)**: Package as downloadable Skill for downstream Agents

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│              Frontend (React + shadcn/ui)           │
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
│   Orchestration | Multi-stage Analysis | Streaming  │
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

## Technology Stack

| Layer        | Technology                                         |
| ------------ | -------------------------------------------------- |
| Frontend     | React 19 + TypeScript + shadcn/ui + Tailwind CSS 4 |
| Backend      | Go + Gin                                           |
| AI Framework | **Eino** (Workflow Orchestration)                  |
| Database     | PostgreSQL + pgvector                              |
| Cache        | Redis                                              |
| File Storage | MinIO (S3-compatible)                              |
| Deployment   | Docker                                             |

## Core Modules

### 1. Frontend

- React 19 with TypeScript
- **shadcn/ui** component library (radix-maia style)
- Tailwind CSS 4 for styling
- Zustand for state management
- React Router 7 for routing
- Chat interface with Markdown rendering and code highlighting
- Project upload with progress tracking (SSE)
- Skill report detail page (`SKILL.md` + `references/`)
- One-click Skill package download (ZIP)
- SSE streaming response display

### 2. API Gateway (Go + Gin)

- Request routing with JWT authentication
- Rate limiting and CORS configuration
- Request/response logging

### 3. Chat Service

- Session management and context window maintenance
- Vector-based skill document retrieval (RAG)
- Streaming LLM message responses via SSE

### 4. Knowledge Base Service (Skill Factory)

Multi-stage workflow:
```
Source Upload (ZIP)
       ↓
[Stage 1: Overview Agent] → Generate SKILL.md
       ↓
[Stage 2: API Retrieval Agent] → Generate references/*.md
       ↓
[Dual Output]
  ├── Path A: Chunk → Embedding → Vector DB (RAG)
  └── Path B: Package → Object Storage (Download)
```

### 5. User Service

- User registration and login
- JWT token management with refresh support

## Data Storage

| Storage        | Purpose                                     |
| -------------- | ------------------------------------------- |
| PostgreSQL     | Users, conversations, messages, vector data |
| Redis          | Session cache, rate limiting, task status   |
| Object Storage | Source ZIP files, generated Skill packages  |

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
│   └── storage/            # File storage utilities
└── web/                    # Frontend (React + shadcn/ui)
```

## Core Features

| Feature                | Description                                         |
| ---------------------- | --------------------------------------------------- |
| **Streaming Chat**     | SSE real-time push, ChatGPT-like experience         |
| **Skill Factory**      | Auto-generate standardized Agent Skill packages     |
| **Dual Output**        | Same analysis serves both RAG and external download |
| **Structured Reports** | Auto-generated SKILL.md + references/               |
| **Skill-based RAG**    | Highly structured retrieval for accurate answers    |
| **One-click Download** | Export Skill packages as ZIP                        |

## Design Principles

- **Skill Standardization**: `SKILL.md` includes clear capability descriptions, I/O specifications, and usage examples
- **Modularity**: Services are independent for easy extension
- **Streaming First**: All LLM responses support SSE streaming
- **Eino as Core**: Leverage Eino's workflow orchestration for multi-stage analysis
- **PostgreSQL Unified**: Relational + vector data in one database
- **Separation of concerns**: handlers → services → repository layers

## Agent Restrictions

> **IMPORTANT**: You are NOT permitted to execute any shell commands directly, except for `task check` defined in Taskfile for static analysis. This restriction is in place to prevent accidental modification or corruption of the development environment.

All code modifications should be performed through file edit tools only.
