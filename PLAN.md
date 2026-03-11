# Cortex Phase 2 Implementation Plan

## Overview

Phase 2 Goal: User Authentication + Conversation History Persistence

This phase adds user authentication (JWT) and persists chat history to PostgreSQL, completing Phase 1's basic chat functionality and laying the foundation for Phase 3 (Knowledge Base Service).

---

## 1. Infrastructure Layer

### 1.1 Configuration Update

**File**: `internal/config/config.go`

Add the following fields to the Config struct:

```go
type Config struct {
    // Existing
    ChatProvider string
    ChatBaseURL  string
    ChatAPIKey   string
    ChatModel    string

    // New - Database
    DatabaseURL string  // postgres://user:pass@localhost:5432/cortex

    // New - JWT Auth
    JWTSecret       string
    JWTExpiryHours  int
}
```

Also update `.env` file with new configuration variables.

### 1.2 Project Structure

```
internal/
├── domain/
│   ├── user.go           # User domain entity
│   ├── conversation.go   # Conversation domain entity
│   └── message.go        # Message domain entity
├── repository/
│   ├── db.go             # GORM connection configuration
│   ├── user.go           # User repository
│   ├── conversation.go   # Conversation repository
│   └── message.go        # Message repository
├── service/
│   ├── user.go           # User service
│   └── conversation.go   # Conversation service
├── auth/
│   ├── claims.go         # JWT Claims definition
│   └── token.go          # Token generation/validation
├── handlers/
│   ├── auth.go           # Auth Handler (new)
│   ├── conversation.go  # Conversation Handler (new)
│   └── chat.go           # Chat Handler (modify)
pkg/
└── middleware/
    └── auth.go           # JWT middleware
```

### 1.3 Database Schema

| Table Name      | Fields                                                                                                            |
| --------------- | ----------------------------------------------------------------------------------------------------------------- |
| `users`         | id (UUID), email (VARCHAR unique), password_hash (VARCHAR), role (VARCHAR default 'user'), created_at, updated_at |
| `conversations` | id (UUID), user_id (UUID FK), title (VARCHAR), created_at, updated_at                                             |
| `messages`      | id (UUID), conversation_id (UUID FK), role (VARCHAR 'user'/'assistant'), content (TEXT), created_at               |

---

## 2. Backend Modules

### 2.1 User Authentication Module

**New Files**:

- `internal/auth/claims.go` - JWT Claims struct with UserID, Email, Role
- `internal/auth/token.go` - Token generation and validation functions
- `pkg/middleware/auth.go` - JWT authentication middleware
- `internal/handlers/auth.go` - Register, Login, Refresh, Me handlers

**Dependencies**:

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get github.com/google/uuid
```

**API Endpoints**:

| Method | Path                    | Description          | Auth |
| ------ | ----------------------- | -------------------- | ---- |
| POST   | `/api/v1/auth/register` | User registration    | No   |
| POST   | `/api/v1/auth/login`    | User login           | No   |
| POST   | `/api/v1/auth/refresh`  | Refresh access token | No   |
| GET    | `/api/v1/auth/me`       | Get current user     | Yes  |

**Request/Response Formats**:

Register:

```json
// Request
{ "email": "user@example.com", "password": "password123" }
// Response 201
{ "user_id": "uuid" }
```

Login:

```json
// Request
{ "email": "user@example.com", "password": "password123" }
// Response 200
{ "access_token": "jwt...", "refresh_token": "jwt..." }
```

### 2.2 Conversation Management Module

**New Files**:

- `internal/repository/conversation.go` - Conversation CRUD operations
- `internal/repository/message.go` - Message CRUD operations
- `internal/service/conversation.go` - Business logic for conversations
- `internal/handlers/conversation.go` - Conversation REST handlers

**API Endpoints**:

| Method | Path                        | Description                    | Auth |
| ------ | --------------------------- | ------------------------------ | ---- |
| GET    | `/api/v1/conversations`     | List user's conversations      | Yes  |
| POST   | `/api/v1/conversations`     | Create new conversation        | Yes  |
| GET    | `/api/v1/conversations/:id` | Get conversation with messages | Yes  |
| DELETE | `/api/v1/conversations/:id` | Delete conversation            | Yes  |

**Request/Response Formats**:

List Conversations:

```json
// Response 200
[
  {
    "id": "uuid",
    "title": "Conversation 1",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
]
```

Create Conversation:

```json
// Request
{ "title": "My Conversation" }
// Response 201
{ "id": "uuid", "title": "My Conversation", "created_at": "timestamp" }
```

Get Conversation:

```json
// Response 200
{
  "id": "uuid",
  "title": "My Conversation",
  "messages": [
    {
      "id": "uuid",
      "role": "user",
      "content": "Hello",
      "created_at": "timestamp"
    },
    {
      "id": "uuid",
      "role": "assistant",
      "content": "Hi!",
      "created_at": "timestamp"
    }
  ]
}
```

### 2.3 Chat Endpoint Modification

**Modify**: `internal/handlers/chat.go`

Changes:

- Add JWT authentication middleware
- Accept optional `conversation_id` in request
- If no conversation_id, create new conversation
- Persist user message and assistant response to database

**Modified Request Format**:

```json
// Request
{ "message": "Hello", "conversation_id": "uuid (optional)" }
// Response: SSE stream as before
```

---

## 3. Frontend Structure

### 3.1 New Files

```
web/src/
├── api/
│   ├── auth.ts            # login, register, refresh, me
│   ├── conversation.ts    # list, create, get, delete
│   └── chat.ts           # send message (modify)
├── stores/
│   ├── user.ts           # user state (token, user info)
│   └── chat.ts          # chat state (modify: conversation_id)
├── views/
│   ├── Login.vue         # Login page
│   ├── Register.vue     # Register page
│   ├── ConversationList.vue # Conversation list page
│   └── ChatView.vue     # Chat page (modify: load history)
└── router/
    └── index.ts          # Add auth guards
```

### 3.2 API Functions

**auth.ts**:

```typescript
interface LoginResponse {
  access_token: string;
  refresh_token: string;
}
export async function login(
  email: string,
  password: string,
): Promise<LoginResponse>;
export async function register(
  email: string,
  password: string,
): Promise<{ user_id: string }>;
export async function getCurrentUser(): Promise<User>;
```

**conversation.ts**:

```typescript
export async function listConversations(): Promise<Conversation[]>;
export async function createConversation(title: string): Promise<Conversation>;
export async function getConversation(
  id: string,
): Promise<ConversationWithMessages>;
export async function deleteConversation(id: string): Promise<void>;
```

**chat.ts** (modify):

- Add `Authorization: Bearer <token>` header
- Add optional `conversation_id` parameter

### 3.3 Router Guards

Add navigation guards to protect routes:

- `/chat` - requires authentication
- `/conversations` - requires authentication
- `/login` and `/register` - redirect to chat if already authenticated

---

## 4. Implementation Order

### Phase 2A: Database Infrastructure (Priority: High)

1. Add GORM dependency: `go get gorm.io/gorm` and `go get gorm.io/driver/postgres`
2. Update config with DatabaseURL
3. Create domain models (user, conversation, message)
4. Implement repositories (user, conversation, message)
5. Test database connection

### Phase 2B: User Authentication (Priority: High)

1. Implement JWT auth module (claims.go, token.go)
2. Create auth middleware
3. Implement user repository and service
4. Create auth handlers (register, login)
5. Wire dependencies in main.go

### Phase 2C: Conversation Persistence (Priority: High)

1. Implement conversation and message repositories
2. Create conversation service
3. Create conversation handlers
4. Modify chat handler to persist messages
5. Test full chat flow with persistence

### Phase 2D: Frontend Integration (Priority: Medium)

1. Create Login.vue and Register.vue
2. Add auth API functions
3. Create user store
4. Add router guards
5. Create ConversationList.vue
6. Modify ChatView.vue to load history
7. Test end-to-end flow

---

## 5. Technology Stack

| Layer            | Technology                         |
| ---------------- | ---------------------------------- |
| Database         | PostgreSQL + GORM                  |
| Auth             | JWT (github.com/golang-jwt/jwt/v5) |
| Password Hashing | bcrypt (golang.org/x/crypto)       |
| Frontend State   | Pinia                              |
| Frontend Routing | Vue Router                         |

---

## 6. Considerations

### Database Connection

- Use connection pooling with appropriate MaxOpenConns and MaxIdleConns
- Implement retry logic with exponential backoff for connection
- Use context for all database operations

### Security

- Always hash passwords with bcrypt (cost >= 12 for production)
- Use environment variables for secrets (JWT_SECRET, DATABASE_URL)
- Implement rate limiting on auth endpoints
- Validate all input in handlers

### Error Handling

- Return generic error messages to clients (don't leak internal errors)
- Log detailed errors server-side
- Use proper HTTP status codes

### Session Management

- Store JWT in httpOnly cookie or localStorage (consider CSRF implications)
- Implement token refresh mechanism
- Consider token blacklisting for logout
