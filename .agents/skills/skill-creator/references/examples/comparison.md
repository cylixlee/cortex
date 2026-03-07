# Description Comparison Examples

Demonstrates what makes a good vs bad description.

## Bad Examples

### Example 1: Too Vague

```yaml
---
name: pdf-tool
description: Helps with PDFs.
---
```

**Problems:**
- Too short (only 16 characters)
- No specific trigger keywords
- Doesn't explain when to use the skill
- Doesn't describe what the skill actually does

### Example 2: Missing Context

```yaml
---
name: data-processor
description: Processes data efficiently.
---
```

**Problems:**
- Doesn't say what kind of data
- Doesn't explain what "processing" means
- No keywords to help agent match to tasks
- Too generic to be useful

### Example 3: Wrong Format

```yaml
---
name: API-Tool
description: This skill is for working with APIs. Use it when you need to do API things.
---
```

**Problems:**
- Name has uppercase (invalid)
- Description uses "things" (too vague)
- Circular explanation ("API-Tool" → "APIs")

---

## Good Examples

### Example 1: Comprehensive Description

```yaml
---
name: pdf-processing
description: Extracts text and tables from PDF files, fills PDF forms, and merges multiple PDFs. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction.
---
```

**Why it works:**
- Lists specific capabilities (extract, fill, merge)
- Includes trigger keywords (PDF, documents, forms, extraction)
- Explains when to use (working with PDFs)
- Length: 186 characters (within 1024 limit)

### Example 2: Action-Oriented

```yaml
---
name: git-hooks
description: Sets up and manages Git hooks for pre-commit and pre-push workflows. Use when initializing a new repository, adding Git hooks, or configuring automated checks before commits and pushes.
---
```

**Why it works:**
- Clear action words (sets up, manages)
- Specific use cases (initializing repo, adding hooks)
- Trigger keywords (Git hooks, pre-commit, pre-push)

### Example 3: Technical with Context

```yaml
---
name: docker-build
description: Builds and optimizes Docker images with multi-stage builds, layer caching, and best practices. Use when creating Dockerfiles, building container images, or optimizing container build times.
---
```

**Why it works:**
- Technical details (multi-stage, layer caching)
- Clear use cases
- Keywords for DevOps tasks (Dockerfile, container)

---

## Writing Guidelines

A good description should answer:

1. **What does this skill do?** - List specific capabilities
2. **When to use it?** - Trigger conditions
3. **What keywords match?** - Words that indicate the skill is needed

### Template

```
{Action verbs} {specific capabilities}. Use when {trigger conditions}, or when {keywords}.
```

### Length

- Minimum: 50 characters (enough to be useful)
- Maximum: 1024 characters
- Recommended: 100-300 characters

### Keywords

Include keywords that help agents match user requests:
- Task types: "extract", "analyze", "format", "validate"
- Domain terms: "PDF", "CSV", "API", "database"
- Action phrases: "working with", "creating", "debugging"
