# Cortex Phase 3 Implementation Plan

## Overview

Phase 3 Goal: Skill Factory - Code Upload -> AI Analysis -> Skill Package Generation

This phase implements the core "Skill Factory" workflow:

1. **File Upload**: Accept source code ZIP files
2. **Stage 1**: Overview Agent -> Generate SKILL.md
3. **Stage 2**: API Retrieval Agent -> Generate references/\*.md
4. **Dual Output**: Vector storage (RAG) + Object storage (Download)

---

## 1. Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                      Upload (ZIP)                        │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                    API Gateway (Gin)                     │
│                   /api/v1/skills/*                      │
└─────────────────────────────────────────────────────────┘
                            │
              ┌─────────────┴─────────────┐
              ▼                           ▼
┌─────────────────────────┐   ┌─────────────────────────┐
│   Handler (HTTP)        │   │   Redis (Task Queue)    │
│   Returns skill_id/task_id│  │   Push task to queue    │
└─────────────────────────┘   └─────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────┐
│              Worker (Async Consumer)                    │
│  1. Extract ZIP                                          │
│  2. Scan source files                                    │
│  3. Eino Workflow (Stage 1 + Stage 2)                   │
│  4. Chunk → Embedding → pgvector                        │
│  5. Package → MinIO                                      │
│  6. Update status → Redis/DB                            │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│               SSE/WebSocket Status Notification         │
└─────────────────────────────────────────────────────────┘
```

---

## 2. Technology Stack

| Layer           | Technology                                          |
| --------------- | --------------------------------------------------- |
| Object Storage  | MinIO (S3 compatible)                              |
| Vector Store    | pgvector (PostgreSQL)                             |
| Task Queue      | Redis                                              |
| Embedding       | Eino (via eino-ext, e.g., Doubao)                 |
| AI Framework    | Eino (Sequential Workflow + Embedder)              |
| File Processing | Go standard library (archive/zip)                  |

---

## 3. Configuration

### 3.1 Config Updates

**File**: `internal/config/config.go`

```go
type Config struct {
    // Existing
    ChatProvider   string
    ChatBaseURL    string
    ChatAPIKey     string
    ChatModel      string
    DatabaseURL    string
    JWTSecret      string
    JWTExpiryHours int

    // New - MinIO Object Storage
    MinIOEndpoint  string
    MinIOAccessKey string
    MinIOSecretKey string
    MinIOBucket    string
    MinIOUseSSL    bool

    // New - Redis
    RedisURL string

    // New - Embedding
    EmbeddingProvider string
    EmbeddingBaseURL  string
    EmbeddingAPIKey   string
    EmbeddingModel    string
}
```

### 3.2 Environment Variables

Update `.env` with:

```bash
# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=cortex-skills
MINIO_USE_SSL=false

# Redis
REDIS_URL=redis://localhost:6379

# Embedding (example: Doubao)
EMBEDDING_PROVIDER=doubao
EMBEDDING_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
EMBEDDING_API_KEY=your-api-key
EMBEDDING_MODEL=embedding-model-name
```

### 3.3 Dependencies

```bash
go get github.com/minio/minio-go/v7
go get github.com/redis/go-redis/v9
```

Note: Embedding uses Eino's built-in `embedding.Embedder` interface (from `github.com/cloudwego/eino/components/embedding`), no additional token counting library needed.

---

## 4. Database Schema

### 4.1 New Tables

| Table Name   | Fields                                                                                                                             |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------------- |
| `skills`     | id (UUID), user_id (UUID FK), name (VARCHAR), description (TEXT), status (VARCHAR), storage_path (VARCHAR), created_at, updated_at |
| `documents`  | id (UUID), skill_id (UUID FK), filename (VARCHAR), language (VARCHAR), content (TEXT), created_at                                  |
| `chunks`     | id (UUID), skill_id (UUID FK), document_id (UUID FK), content (TEXT), embedding (vector), chunk_index (INT)                        |
| `references` | id (UUID), skill_id (UUID FK), filename (VARCHAR), content (TEXT), created_at                                                      |

### 4.2 File Structure

```
internal/models/
├── skill.go       # Skill model
├── document.go    # Source document model
├── chunk.go       # Vector chunk model
└── reference.go  # Reference document model
```

---

## 5. Backend Modules

### 5.1 Project Structure

```
internal/
├── config/
│   └── config.go                 # Update: add MinIO/Redis/Embedding config
├── models/
│   ├── skill.go                 # Skill model
│   ├── document.go              # Document model
│   ├── chunk.go                 # Chunk model
│   └── reference.go            # Reference model
├── repository/
│   ├── skill.go                 # Skill CRUD
│   ├── document.go              # Document CRUD
│   ├── chunk.go                 # Chunk CRUD
│   └── reference.go             # Reference CRUD
├── service/
│   ├── skill.go                 # Core business logic
│   ├── embedding.go             # Vector embedding
│   └── storage.go               # MinIO storage
├── handlers/
│   └── skill.go                 # HTTP handlers
├── worker/
│   └── skill.go                 # Async task worker
├── workflows/
│   └── skill.go                 # Eino workflow
└── pkg/
    └── storage/
        └── minio.go             # MinIO client
```

### 5.2 API Endpoints

| Method | Path                          | Description                 | Auth |
| ------ | ----------------------------- | --------------------------- | ---- |
| POST   | `/api/v1/skills/upload`       | Upload ZIP file             | Yes  |
| GET    | `/api/v1/skills`              | List user's skills          | Yes  |
| GET    | `/api/v1/skills/:id`          | Get skill details           | Yes  |
| GET    | `/api/v1/skills/:id/download` | Download skill package      | Yes  |
| DELETE | `/api/v1/skills/:id`          | Delete skill                | Yes  |
| GET    | `/api/v1/skills/:id/status`   | Get processing status (SSE) | Yes  |

### 5.3 Request/Response Formats

**POST /api/v1/skills/upload**:

```json
// Request (multipart/form-data)
{
  "file": (ZIP file),
  "name": "my-awesome-lib"
}

// Response 202
{
  "skill_id": "uuid",
  "status": "processing"
}
```

**GET /api/v1/skills/:id**:

```json
// Response 200
{
  "id": "uuid",
  "name": "my-awesome-lib",
  "description": "A description",
  "status": "completed",
  "skill": {
    "overview": "# SKILL.md content...",
    "references": [
      { "filename": "api.md", "content": "..." },
      { "filename": "types.md", "content": "..." }
    ]
  },
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

**GET /api/v1/skills/:id/status** (SSE):

```json
// Streaming
{"status": "processing", "progress": 25}
{"status": "processing", "progress": 50}
{"status": "processing", "progress": 75}
{"status": "completed", "progress": 100}
```

---

## 6. Eino Workflow Design

### 6.1 Two-Stage Sequential Workflow

```go
// Stage 1: Overview Agent
overviewAgent := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "overview_agent",
    Instruction: overviewSystemPrompt,
    Model:       chatModel,
})

// Stage 2: API Retrieval Agent (depends on Stage 1 output)
apiRetrievalAgent := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "api_retrieval_agent",
    Instruction: apiRetrievalSystemPrompt,
    Model:       chatModel,
})

// Sequential Workflow
workflow := workflow.NewSequential(ctx, workflow.SequentialConfig{
    Nodes: []workflow.Node{
        {Name: "overview", Agent: overviewAgent},
        {Name: "api_retrieval", Agent: apiRetrievalAgent},
    },
})
```

### 6.2 Prompts

**Overview Agent** (generate SKILL.md):

- Analyze source code structure
- Identify main capabilities and features
- Generate standardized SKILL.md format

**API Retrieval Agent** (generate references/):

- Based on SKILL.md output
- Extract detailed API signatures
- Generate reference documentation (api.md, types.md, etc.)

---

## 7. Embedding Strategy

- Use Eino's built-in `embedding.Embedder` interface for vector generation
- Use factory pattern (like `pkg/llm/client.go`) to support multiple embedding providers
- Chunk size: 512 tokens (approximate, use simple character-based splitting)
- Overlap: 50 tokens
- Store embeddings in pgvector

### Factory Pattern Implementation

Similar to `pkg/llm/client.go`, create `pkg/llm/embedding.go`:

```go
package llm

import (
    "context"
    "errors"
    
    "github.com/cloudwego/eino/components/embedding"
    "github.com/cloudwego/eino-ext/components/embedding/doubao"
)

func NewEmbedder(ctx context.Context, provider, baseURL, apiKey, model string) (embedding.Embedder, error) {
    switch provider {
    case "doubao":
        return doubao.NewEmbeddingModel(ctx, &doubao.EmbeddingModelConfig{
            BaseURL: baseURL,
            APIKey:  apiKey,
            Model:   model,
        })
    // Future: add more providers (openai, azure, etc.)
    default:
        return nil, errors.New("unsupported embedding provider: " + provider)
    }
}
```

### Chunk Model

```go
type Chunk struct {
    ID           uuid.UUID
    SkillID      uuid.UUID
    DocumentID   uuid.UUID
    Content      string
    Embedding    []float64  // Eino uses []float64
    ChunkIndex   int
}
```

---

## 8. Frontend Modules

### 8.1 File Structure

```
web/src/
├── api/
│   └── skill.ts                 # Skill API client
├── views/
│   ├── SkillList.vue            # Skill list page
│   ├── Upload.vue               # Upload page
│   └── SkillDetail.vue          # Skill detail page
├── components/
│   ├── SkillCard.vue            # Skill card component
│   ├── FileUploader.vue         # File upload component
│   └── MarkdownRenderer.vue     # Markdown renderer
└── router/
    └── index.ts                 # Update: add routes
```

### 8.2 Routes

```typescript
{
  path: '/skills',
  name: 'skills',
  component: SkillList,
  meta: { requiresAuth: true },
},
{
  path: '/skills/upload',
  name: 'skill-upload',
  component: Upload,
  meta: { requiresAuth: true },
},
{
  path: '/skills/:id',
  name: 'skill-detail',
  component: SkillDetail,
  meta: { requiresAuth: true },
},
```

---

## 9. Implementation Order

### Phase 3A: Infrastructure (Priority: High)

1. Update `internal/config/config.go` with new fields
2. Update `.env` with MinIO/Redis/Embedding config
3. Add MinIO/Redis dependencies
4. Create `pkg/storage/minio.go` - MinIO client
5. Update `docker-compose.yml` - add MinIO and Redis services

### Phase 3B: Database Layer (Priority: High)

1. Create `internal/models/skill.go`
2. Create `internal/models/document.go`
3. Create `internal/models/chunk.go`
4. Create `internal/models/reference.go`
5. Update `internal/repository/db.go` - add migrations

### Phase 3C: Repository Layer (Priority: High)

1. Create `internal/repository/skill.go`
2. Create `internal/repository/document.go`
3. Create `internal/repository/chunk.go`
4. Create `internal/repository/reference.go`

### Phase 3D: Service Layer (Priority: High)

1. Create `internal/service/storage.go` - MinIO operations
2. Create `internal/service/embedding.go` - Vector embedding
3. Create `internal/service/skill.go` - Core business logic

### Phase 3E: Workflow & Worker (Priority: High)

1. Create `internal/workflows/skill.go` - Eino workflow
2. Create `internal/worker/skill.go` - Async task worker

### Phase 3F: Handler Layer (Priority: High)

1. Create `internal/handlers/skill.go` - HTTP handlers
2. Update `cmd/api/main.go` - wire dependencies and routes

### Phase 3G: Frontend (Priority: Medium)

1. Create `web/src/api/skill.ts` - API client
2. Create `web/src/views/SkillList.vue`
3. Create `web/src/views/Upload.vue`
4. Create `web/src/views/SkillDetail.vue`
5. Create components: SkillCard, FileUploader, MarkdownRenderer
6. Update `web/src/router/index.ts`

---

## 10. Docker Updates

### 10.1 docker-compose.yml

Add MinIO and Redis services:

```yaml
services:
  # Existing: PostgreSQL
  postgres:
    image: pgvector/pgvector:pg16
    # ... existing config

  # New: MinIO
  minio:
    image: minio/minio:latest
    container_name: cortex-minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  # New: Redis
  redis:
    image: redis:7-alpine
    container_name: cortex-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  minio_data:
  redis_data:
```

---

## 11. Considerations

### 11.1 Error Handling

- Implement retry logic for MinIO and Redis connections
- Use context for all async operations
- Return generic error messages to clients
- Log detailed errors server-side

### 11.2 Security

- Validate uploaded files (check ZIP format, max size)
- Sanitize filenames to prevent path traversal
- Use UUIDs for all resource identifiers
- Implement proper CORS settings

### 11.3 Performance

- Process large ZIP files in chunks
- Use connection pooling for MinIO and Redis
- Implement proper timeouts for LLM calls
- Use streaming for large file uploads/downloads

### 11.4 Testing

- Unit tests for repositories and services
- Integration tests for workflow
- E2E tests for upload → process → download flow
