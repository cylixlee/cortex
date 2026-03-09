# Cortex - AI Source Code Analysis Platform Architecture Document

## Project Overview

Cortex is a combination of an **"Agent Skill Factory"** and an **intelligent code assistant**. It is not just a code Q&A tool, but a platform that can automatically transform any source code repository into a standardized **Agent Skill**.

**Core Workflow:**
1. **Upload & Analyze**: User uploads a source code ZIP, system automatically triggers multi-stage Agent orchestration.
2. **Skill Generation**:
   - **Overview Agent** distills core concepts and generates `SKILL.md` (standardized capability definition file).
   - **API Retrieval Agent** explores details and generates `references/` directory with detailed reference documentation.
3. **Dual Output**:
   - **Internal Empowerment (RAG)**: Generated Skill documents are vectorized and stored in knowledge base, serving as precise context for user conversations, enabling Q&A based on "official documentation-level" understanding.
   - **External Reuse (Download)**: Complete Skill package (`SKILL.md` + `references/`) can be downloaded by users with one click, directly used as standard Skill modules for downstream Agent development.

Through Cortex, developers can quickly transform obscure legacy code into AI-understandable assets that can be invoked by other Agents.

---

## Overall Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (Vue)                          │
│        Code Chat Interface | Project Upload | Skill Report  │
│                      Detail | Skill Download                 │
└─────────────────────────────────────────────────────────────┘
                               │ HTTP/SSE
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (Go)                          │
│              Auth | Rate Limit | Routing | Logging          │
└─────────────────────────────────────────────────────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         ▼                     ▼                     ▼
 ┌───────────────┐    ┌───────────────┐    ┌───────────────┐
 │  Chat Service │    │  Code KB Svc  │    │  User Service │
 │   (Eino)      │    │   (Eino)      │    │               │
 └───────────────┘    └───────────────┘    └───────────────┘
         │                     │                     │
         └─────────────────────┼─────────────────────┘
                               ▼
 ┌─────────────────────────────────────────────────────────────┐
 │                  Eino Workflow Engine                        │
 │   Workflow Orchestration | Multi-stage Code Analysis |      │
 │   Skill Generation | Semantic Retrieval | Streaming Output  │
 └─────────────────────────────────────────────────────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         ▼                     ▼                     ▼
 ┌───────────────┐    ┌───────────────┐    ┌───────────────┐
 │  PostgreSQL   │    │     Redis     │    │   Vector DB   │
 │  + pgvector   │    │   (Cache)     │    │  (pgvector)   │
 └───────────────┘    └───────────────┘    └───────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   Object Storage    │
                    │ (Skill Packages:    │
                    │  SKILL.md + refs/)  │
                    └─────────────────────┘
```

---

## Core Modules

### 1. Frontend

| Technology        | Description             |
| ----------------- | ----------------------- |
| **Vue**           | Main framework          |
| TypeScript        | Type safety             |
| SSE               | Streaming communication |
| Markdown Renderer | Report rendering        |

**Responsibilities:**
- Chat interface (with Markdown and code highlighting support)
- Project upload and status viewing
- **Skill Report Detail Page**: Displays `SKILL.md` and `references/` content
- **Skill Download**: One-click download of generated Skill packages (ZIP)
- Streaming response rendering

---

### 2. API Gateway

| Technology | Description    |
| ---------- | -------------- |
| Go         | Main language  |
| Gin        | Web framework  |
| JWT        | Authentication |

**Responsibilities:**
- Request routing and distribution
- User authentication and authorization
- Rate limiting
- Request/response logging

---

### 3. Chat Service

| Technology | Description            |
| ---------- | ---------------------- |
| **Eino**   | Workflow orchestration |
| SSE        | Streaming output       |
| Redis      | Session cache          |

**Core Workflow:**
```
User Question → Vector Search (Skill Docs) → Build Prompt → LLM Call → Streaming Response
```

**Responsibilities:**
- Create and manage conversation sessions
- Maintain conversation context window
- Answer questions based on retrieved "Skill documents"
- Message persistence

---

### 4. Knowledge Base Service (Skill Factory) ⭐ Core Refactoring

| Technology | Description                                          |
| ---------- | ---------------------------------------------------- |
| **Eino**   | **Multi-stage Analysis & Skill Generation Workflow** |
| LLM        | Code analysis and Skill definition                   |
| pgvector   | Vector storage                                       |

**Core Workflow (Multi-Stage Skill Generation):**
```
1. Source Upload (ZIP)
   ↓
2. [Stage 1: Overview Agent]
   ├→ Traverse codebase structure
   ├→ Distill core concepts, main usage, input/output specifications
   └→ Generate SKILL.md (standard Agent Skill description file)
   ↓
3. [Stage 2: API Retrieval Agent]
   ├→ Based on SKILL.md guidance
   ├→ Deep dive into specific module/API implementation details
   └→ Generate references/xxx.md (detailed reference documentation)
   ↓
4. [Product Processing: Dual Output]
   ├─> Path A (Internal RAG):
   │    Chunk SKILL.md + references/*.md → Embedding → Store in Vector DB
   ├─> Path B (External Download):
   │    Package SKILL.md + references/ directory → Store in Object Storage (for user download)
   ↓
5. Completion Notification
```

**Responsibilities:**
- Execute multi-stage code analysis workflow
- **Generate standardized Skill packages**: Strictly follow `SKILL.md` (capability definition) + `references/` (implementation details) structure
- **Dual-path data flow**: Serve both internal retrieval augmentation and external skill reuse
- Semantic retrieval

---

### 5. User Service

| Technology | Description          |
| ---------- | -------------------- |
| JWT        | Token authentication |
| bcrypt     | Password encryption  |

**Responsibilities:**
- User registration/login
- Token management
- User preference configuration

---

### 6. AI Service Layer

| Technology | Description                |
| ---------- | -------------------------- |
| Eino LLM   | Unified model call wrapper |

**Responsibilities:**
- Unified LLM calls through Eino framework
- Streaming call wrapper
- Error retry logic

---

## Data Storage

| Storage        | Purpose                                                                 |
| -------------- | ----------------------------------------------------------------------- |
| PostgreSQL     | Users, conversations, messages, vector data                             |
| Redis          | Session cache, rate limiting, task status                               |
| Object Storage | Source ZIP files, **generated Skill packages** (SKILL.md + references/) |

---

## Core Workflows

### Chat Workflow
```
┌─────────┐    ┌─────────    ┌─────────┐    ┌─────────┐
│ User    │ →  │ Vector   │ →  │ RAG     │ →  │ LLM     │
│ Query   │    │ Search   │    │ Prompt  │    │ Response│
│         │    │ (Skill   │    │ Build   │    │         │
│         │    │  Docs)   │    │         │    │         │
└─────────┘    └─────────┘    └─────────┘    └─────────┘
```

### Code-to-Skill Workflow
```
┌─────────┐
│ File    │
│ Upload  │
└────┬────┘
     │
     ▼
┌─────────────────────────────────────────┐
│  Eino Multi-stage Analysis Workflow     │
│                                         │
│  1. [Overview Agent]                    │
│     └─> Generate SKILL.md (Core          │
│         Capability Definition)          │
│                                         │
│  2. [API Retrieval Agent] (Depends on   │
│     Stage 1)                            │
│     └─> Generate references/xxx.md       │
│                                         │
│  3. [Dual Output Processing]             │
│     ├─> Path A: Chunk → Embedding →      │
│     │   Vector DB (for RAG)              │
│     └─> Path B: Package → Object        │
│         Storage (for download)           │
└─────────────────────────────────────────┘
```

---

## Technology Stack Overview

| Layer             | Technology Choice         |
| ----------------- | ------------------------- |
| **Frontend**      | Vue + TypeScript          |
| **Backend**       | Go + Gin                  |
| **LLM Framework** | **Eino**                  |
| **Database**      | PostgreSQL + pgvector     |
| **Cache**         | Redis                     |
| **File Storage**  | Object Storage (MinIO/S3) |
| **Deployment**    | Docker                    |

---

## Project Structure

```
cortex/
├── cmd/                    # Application entry points
├── internal/               # Internal business logic
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP handlers
│   ├── services/           # Business logic layer
│   ├── workflows/          # Eino workflow definitions (Skill generation/conversation)
│   ├── models/             # Data models
│   ├── repository/        # Data access layer
│   └── middleware/         # Middleware
├── pkg/                    # Shared packages
│   ├── llm/                # LLM client wrapper (Eino)
│   ├── vector/             # Vector database client
│   └── storage/            # File storage utilities
├── web/                    # Frontend project (Vue)
└── deployments/            # Deployment configs (Docker/K8s)
```

---

## Core Features

| Feature                   | Description                                                                                                   |
| ------------------------- | ------------------------------------------------------------------------------------------------------------- |
| **Streaming Chat**        | SSE real-time push, ChatGPT-like experience                                                                   |
| **Skill Factory Pattern** | Upload code and automatically generate standardized **Agent Skill packages** (`SKILL.md` + `refs`)            |
| **Dual Value Output**     | Same analysis result serves both internal RAG knowledge base and external download for downstream Agent reuse |
| **Structured Reports**    | Auto-generated `SKILL.md` (capability definition) and `references/` (implementation details)                  |
| **Skill-based RAG**       | Highly structured Skill documents as retrieval content, improving answer accuracy                             |
| **One-click Download**    | Export generated Skill packages as ZIP for other Agent projects                                               |

---

## Development Priority

1. **Phase 1**: Basic chat functionality (Chat + LLM integration)
2. **Phase 2**: User authentication + conversation history
3. **Phase 3**: **Skill generation workflow** (Source upload → Generate SKILL.md/refs → Storage)
4. **Phase 4**: Frontend Skill detail page + **download functionality** implementation
5. **Phase 5**: Optimize Skill description format to ensure compatibility with mainstream Agent frameworks

---

## Design Principles

- **Skill Standardization**: Generated `SKILL.md` must include clear "capability descriptions", "input/output definitions", and "usage examples" so both humans and AI can understand.
- **Modularity**: Services are independent for easy extension
- **Streaming First**: All LLM responses support SSE streaming
- **Eino as Core**: Fully leverage Eino's workflow orchestration capabilities for multi-stage analysis and dual output
- **PostgreSQL Unified**: Relational data + vector data stored together
- **Go Performance**: Leverage Go's concurrency advantages for high-throughput requests
- **Knowledge as a Service (KaaS)**: Provide not just Q&A, but also reusable knowledge assets (Skills)
