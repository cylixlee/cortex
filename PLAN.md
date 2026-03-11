# Cortex Project Implementation Plan

## Overview

This document tracks the implementation plans for Cortex project improvements beyond Phase 1-4.

---

## Phase 1: Vector Semantic Search (RAG Enhancement)

### Background

The project currently lacks vector-based semantic search capability. When users ask questions about their uploaded code, the system cannot retrieve relevant context from the knowledge base. This feature will enable users to optionally enable "Skill Retrieval" during chat, similar to DeepSeek's "Deep Think" toggle.

### Architecture

```
User Message + Enable RAG Flag
         │
         ▼
┌────────────────────────┐
│  Handler Layer         │
│  (chat.go)            │
└────────────────────────┘
         │
         ▼ (if enable_rag = true)
┌────────────────────────┐
│  ChatService           │
│  RetrieveContext()    │
└────────────────────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌────────┐ ┌──────────────┐
│Embed   │ │Chunk Repo    │
│Query   │ │SearchByEmbed │
└────────┘ └──────────────┘
    │         │
    └────┬────┘
         ▼
┌────────────────────────┐
│  Context + Message     │
│  → LLM                 │
└────────────────────────┘
```

### Implementation Details

#### 1. Frontend Changes

**File**: `web/src/views/ChatView.vue`

- Add a toggle button "Skill 检索" (Skill Retrieval) in the message input area
- Store toggle state locally (not persisted to backend)
- Pass `enable_rag: boolean` in chat API request

**File**: `web/src/api/chat.ts` (or similar)

```typescript
interface ChatRequest {
  message: string;
  conversation_id?: string;
  enable_rag?: boolean;  // New field
}
```

#### 2. Backend Handler Changes

**File**: `internal/handlers/chat.go`

```go
type ChatRequest struct {
    Message        string `json:"message" binding:"required"`
    ConversationID string `json:"conversation_id"`
    EnableRAG      bool   `json:"enable_rag"`  // New field
}
```

In the `Chat()` handler:
- Check `req.EnableRAG`
- If enabled, call `chatService.RetrieveContext()` to get relevant chunks
- Prepend context to user message before sending to LLM

#### 3. Service Layer Changes

**File**: `internal/service/chat.go`

Add new fields to `ChatService` struct:

```go
type ChatService struct {
    // ... existing fields
    embeddingClient llm.Embedder
    chunkRepo       *repository.ChunkRepository
}
```

Add new method:

```go
func (s *ChatService) RetrieveContext(ctx context.Context, query string, topK int) ([]string, error) {
    // 1. Embed the query string
    embeddings, err := s.embeddingClient.EmbedStrings(ctx, []string{query})
    if err != nil {
        return nil, err
    }

    // 2. Search chunks by embedding (search across all skills, or specific skill)
    chunks, err := s.chunkRepo.SearchByEmbedding(embeddings[0], topK)
    if err != nil {
        return nil, err
    }

    // 3. Extract content
    contexts := make([]string, len(chunks))
    for i, chunk := range chunks {
        contexts[i] = chunk.Content
    }
    return contexts, nil
}
```

**Note**: Need to discuss search scope - should it search:
- All chunks across all skills (global search)
- User's own skills only
- Specific skill (requires frontend skill selector)

#### 4. Repository Layer

**File**: `internal/repository/chunk.go`

The `SearchByEmbedding` method already exists but needs review:

```go
func (r *ChunkRepository) SearchByEmbedding(embedding []float64, skillID uuid.UUID, limit int) ([]models.Chunk, error)
```

Current implementation requires `skillID`. We may need to create an overloaded method or modify to support global search (without skillID filter).

#### 5. Dependency Injection

**File**: `cmd/api/main.go`

Update `NewChatService()` call to include:
- embeddingClient
- chunkRepo

#### 6. Database Schema

If using global search without skillID filter, no schema changes needed.

If searching by skill, may need to track user's skills or add skill selector in frontend.

### Simplified Message Format

Instead of using system prompt, simply prepend context to user message:

```
Based on the following context:

<chunk1 content>

---

<chunk2 content>

---

Question: <user's original message>
```

This approach requires no changes to Session or Client logic.

### File Changes Summary

| # | File | Changes |
|---|------|---------|
| 1 | `web/src/views/ChatView.vue` | Add RAG toggle button |
| 2 | `web/src/api/chat.ts` | Add `enable_rag` field to request type |
| 3 | `internal/handlers/chat.go` | Add `EnableRAG` field, implement RAG logic |
| 4 | `internal/service/chat.go` | Add embeddingClient, chunkRepo, RetrieveContext() |
| 5 | `internal/repository/chunk.go` | Optionally modify SearchByEmbedding for global search |
| 6 | `cmd/api/main.go` | Update ChatService initialization |

### Open Questions

~~1. **Search Scope**: Should RAG search all chunks globally, or only user's own skills?~~
~~2. **Chunk Limit**: What is a reasonable default for `topK`? (Suggested: 5)~~
~~3. **Fallback Behavior**: If embedding service fails, should chat fail or continue without RAG?~~

### Decisions Made

1. **Search Scope**: Search only user's own skills (filtered by user_id)
2. **Chunk Limit**: topK = 5 (default)
3. **Fallback Behavior**: Graceful degradation - if embedding fails, continue chat without RAG (log error, don't block user)

### Additional Implementation Details

#### Repository Layer - Search Scope Implementation

Since the search must be scoped to user's skills, we need to:

1. First get all SkillIDs belonging to the current user
2. Then search chunks filtered by those SkillIDs

**File**: `internal/repository/chunk.go`

Add a new method for user-scoped search:

```go
func (r *ChunkRepository) SearchByEmbeddingForUser(embedding []float64, userID uuid.UUID, limit int) ([]models.Chunk, error) {
    // First get user's skill IDs
    var skillIDs []uuid.UUID
    r.db.Model(&models.Skill{}).Where("user_id = ?", userID).Pluck("id", &skillIDs)

    if len(skillIDs) == 0 {
        return []models.Chunk{}, nil
    }

    var chunks []models.Chunk
    err := r.db.Where("skill_id IN ?", skillIDs).
        Order("embedding <-> ?").
        Limit(limit).
        Find(&chunks, embedding).Error
    return chunks, err
}
```

**Verified**: Skill model already has `user_id` field (line 70 in skill.go).

---

## Phase 2: Admin Dashboard UI

### Background

Currently the backend has basic admin structures but no frontend admin interface. This phase adds management capabilities for users and skills.

### Planned Features

- User management (list, disable/enable)
- Skill management (list, view details, delete)
- Statistics dashboard

### Files to Create

- `web/src/views/admin/Users.vue`
- `web/src/views/admin/Skills.vue`
- `web/src/views/admin/Dashboard.vue`

### Priority: Medium

---

## Phase 3: Agent Framework Format Optimization (Phase 5)

### Background

The system generates standard SKILL.md format. This phase adds support for exporting to specific agent framework formats (OpenAI GPTs, Claude, LangChain, etc.).

### Planned Features

- Multi-format export options
- JSON Schema format
- Function Calling definitions

### Priority: Low

---

## Implementation Order

| Order | Phase | Priority |
|-------|-------|----------|
| 1 | Vector Semantic Search | High |
| 2 | Admin Dashboard UI | Medium |
| 3 | Agent Framework Format | Low |

---

*Last Updated: 2026-03-11*
