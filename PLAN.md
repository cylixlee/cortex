# Cortex Implementation Plan

## Current Status

- **Phase 1**: 基础聊天功能 (Chat + LLM) ✅ 100%
- **Phase 2**: 用户认证 + 对话历史 ✅ 100%
- **Phase 3**: Skill 生成工作流 ✅ 95%
- **Phase 4**: 前端技能页面 + 下载功能 ✅ 100%
- **Phase 5**: 优化 Skill 格式兼容性 ⏳ 未开始

---

# Plan 1: 代码质量改进 (Code Quality Improvements)

## 1.1 已完成的修复

| 问题 | 修复位置 | 状态 |
|------|----------|------|
| 权限控制不足 | `handlers/auth.go` - ListUsers/DeleteUser 管理员检查 | ✅ |
| 敏感信息泄露 | 统一错误响应格式 | ✅ |
| 文件上传无限制 | 添加 100MB 大小限制 | ✅ |
| CORS 通配符 | 改为配置化 `CORS_ALLOWED_ORIGINS` | ✅ |
| MinIO 资源泄漏 | 添加 defer obj.Close() | ✅ |
| uuid.Parse 错误忽略 | 添加错误处理 | ✅ |
| SessionManager 内存泄漏 | 添加 30 分钟过期清理 | ✅ |
| 登录无速率限制 | 添加内存速率限制 | ✅ |
| LLM 无超时 | 添加可配置超时 CHAT_TIMEOUT | ✅ |
| DeleteSkill 无事务 | 添加事务删除 | ✅ |

---

## 1.2 待完成的修复

### Middleware 包重构

**问题**: `pkg/middleware/ratelimit.go` 和 `pkg/middleware/login_ratelimit.go` 包含下划线

**方案**: 重构为 `pkg/middleware/ratelimit/` 包结构

```
pkg/middleware/ratelimit/
├── ratelimit.go    # 公共接口定义
├── login.go        # 登录速率限制（内存实现）
└── redis.go        # Redis 实现（可选）
```

### 配置项更新

在 `.env` 中添加:

```bash
# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173

# Chat Timeout (秒)
CHAT_TIMEOUT=120
```

---

# Plan 2: Vue 前端改进 (Frontend Improvements)

## 2.1 问题清单

| 问题 | 严重程度 | 修复方案 |
|------|----------|----------|
| API_BASE 不统一 | 🔴 高 | 创建 `api/client.ts` 统一配置 |
| getToken 逻辑重复 | 🔴 高 | 提取到 `api/auth.ts` |
| fetchWithAuth 逻辑重复 | 🔴 高 | 统一使用 |
| 未使用文件 | 🟡 中 | 删除 `counter.ts`, `chat.ts` |
| formatDate 重复 | 🟡 中 | 提取到 `utils/date.ts` |
| 类型定义分散 | 🟡 中 | 创建 `src/types/` |
| 缺少公共组件 | 🟡 中 | 提取 Badge, Loading, Empty |
| 硬编码 API 地址 | 🟡 中 | 使用 `.env` |
| 缺少 Markdown 库 | 🟡 中 | 添加 `marked` |

## 2.2 实施步骤

### Step 1: 统一 API 配置

创建 `web/src/api/client.ts`:

```typescript
// web/src/api/client.ts
const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api/v1'

export { API_BASE }
```

更新所有 API 文件导入 `API_BASE`

### Step 2: 删除未使用文件

- `web/src/stores/counter.ts`
- `web/src/api/chat.ts`

### Step 3: 提取公共类型

创建 `web/src/types/`:

```
web/src/types/
├── api.ts          # API 响应类型
├── skill.ts        # Skill 相关类型
├── conversation.ts # 会话类型
└── index.ts        # 统一导出
```

### Step 4: 提取公共组件

创建 `web/src/components/`:

```
web/src/components/
├── StatusBadge.vue    # 状态徽章
├── LoadingSpinner.vue # 加载中
└── EmptyState.vue    # 空状态
```

### Step 5: 环境变量配置

创建 `web/.env`:

```bash
VITE_API_BASE=http://localhost:8080/api/v1
```

### Step 6: 添加 Markdown 支持

```bash
cd web && pnpm add marked
```

---

# Plan 3: 架构优化 (Architecture Improvements)

## 3.1 Go 后端 - DDD 演进

### 当前结构 (扁平三层)

```
internal/
├── handlers/  → services → repository
```

### 建议结构 (引入 Domain 层)

```
internal/
├── domain/           # 领域层 - 业务实体和规则
│   ├── skill/        # Skill 聚合
│   │   ├── entity.go
│   │   └── errors.go
│   └── user/        # User 聚合
├── handlers/         # 接口层
├── services/         # 应用层
└── repository/       # 基础设施层
```

### 错误定义原则

错误应该定义在产生它的包里，而非中央 errors 包：

```go
// internal/domain/skill/errors.go
var (
    ErrSkillNotFound = errors.New("skill not found")
    ErrSkillNotReady = errors.New("skill not ready")
)

// internal/service/skill.go
if err := r.db.First(&skill, id).Error; err != nil {
    return nil, fmt.Errorf("find skill %s: %w", id, domain.ErrSkillNotFound)
}
```

## 3.2 可观测性 (Observability)

### 添加 OpenTelemetry

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/otlp
```

### 添加 Prometheus 指标

```bash
go get github.com/prometheus/client_golang/prometheus
```

### 健康检查增强

```go
r.GET("/health", func(c *gin.Context) {
    // 检查 DB
    // 检查 Redis
    // 检查 MinIO
    c.JSON(200, gin.H{
        "status": "ok",
        "db": "ok",
        "redis": "ok",
        "minio": "ok",
    })
})
```

---

# Implementation Order

## Phase A: Middleware 重构 (Priority: High)

1. 重构 `pkg/middleware/ratelimit/` 包结构
2. 更新 `cmd/api/main.go` 导入路径
3. 添加 `.env` 配置项

## Phase B: Vue 前端改进 (Priority: High)

1. 创建 `api/client.ts` 统一 API 配置
2. 删除未使用文件
3. 创建 `types/` 目录
4. 创建公共组件
5. 配置环境变量
6. 添加 Markdown 支持

## Phase C: 架构演进 (Priority: Medium)

1. 引入 domain 层（可选）
2. 添加可观测性
3. 增强健康检查

---

# Original Phase 3 Implementation Plan

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
