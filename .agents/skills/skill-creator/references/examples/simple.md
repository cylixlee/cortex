# Simple Skill Example

A straightforward skill with clear structure and minimal resources.

```mdx
---
name: markdown-formatter
description: Formats and lints Markdown files. Use when editing, writing, or reviewing Markdown documents, or when the user mentions formatting Markdown or checking Markdown syntax.
---

# Markdown Formatter

## When to Use This Skill

Use this skill when working with Markdown files:
- Writing new Markdown documents
- Editing existing Markdown
- Reviewing Markdown content
- User asks to format or lint Markdown

## How to Format Markdown

1. Install required tools:
   ```bash
   uvx markdownlint@latest --fix .
   ```

2. Common fixes applied:
   - Consistent heading hierarchy
   - Proper list indentation
   - Code block language hints
   - Link validation

## Example

Input:
```markdown
# title

some text

## sub

more text
```

Output:
```markdown
# Title

Some text

## Sub

More text
```

## Common Issues

- **Broken links**: Verify all links are accessible
- **Missing alt text**: Ensure images have descriptive alt text
- **Long lines**: Wrap lines at 80-100 characters
```

## Key Characteristics

This example demonstrates:
- Required frontmatter with valid name and description
- Clear "When to Use This Skill" section
- Step-by-step instructions
- Concrete input/output examples
- Common edge cases section
