# Skill with Multiple Resources Example

A skill that references multiple supporting documents.

```mdx
---
name: api-design
description: Design RESTful APIs following best practices for HTTP APIs. Use when creating new API endpoints, designing API contracts, or reviewing API implementations.
---

# API Design

## When to Use This Skill

- Creating new API endpoints
- Designing API contracts
- Reviewing API implementations
- Writing API documentation

## Core Principles

### Resource-Oriented Design

Design around resources (nouns), not actions (verbs):
```
/users          # User collection
/users/{id}     # Specific user
/orders         # Order collection
```

### HTTP Methods

Use HTTP methods semantically:
- GET: Retrieve resources
- POST: Create new resources
- PUT: Replace entire resource
- PATCH: Partial update
- DELETE: Remove resources

### Status Codes

Use appropriate HTTP status codes:
- 200: Success
- 201: Created
- 204: No Content
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Server Error

## Best Practices

See [references/best-practices.md](references/best-practices.md) for detailed guidelines.

## Error Handling

See [references/error-handling.md](references/error-handling.md) for error response patterns.

## Versioning

See [references/versioning.md](references/versioning.md) for API versioning strategies.

## Examples

See [references/examples/](references/examples/) for complete API examples.
```

## Directory Structure

```
api-design/
├── SKILL.md
└── references/
    ├── best-practices.md
    ├── error-handling.md
    ├── versioning.md
    └── examples/
        ├── user-api.md
        └── order-api.md
```

## Key Characteristics

This example demonstrates:
- Referencing external documentation in references/ directory
- Using relative links to other reference files
- Condensed core content with links to detailed docs
- Supporting progressive disclosure of information

## Progressive Disclosure Pattern

This skill uses progressive disclosure effectively:
1. **Catalog**: name + description shown at session start
2. **Instructions**: Core principles and quick reference
3. **Resources**: Detailed docs loaded only when needed
