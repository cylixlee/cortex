---
name: skill-creator
description: Creates new Agent Skills following the Agent Skills specification. Use when the user asks to create a new skill, add a skill, define a skill, or build a skill for a specific capability.
---

# Skill Creator

This skill guides the creation of new Agent Skills that conform to the Agent Skills specification. Follow these instructions whenever you need to create a new skill for any capability.

## When to Use This Skill

- User asks to create a new skill
- User requests adding a skill for a specific capability
- User wants to define or build a skill
- User needs to package knowledge or workflows as a reusable skill

## Agent Skills Core Concepts

Skills are portable, self-documenting packages that extend AI agent capabilities. Each skill is a directory containing a `SKILL.md` file with instructions.

For complete specification details, see [references/spec.md](references/spec.md).

### Progressive Disclosure Architecture

Skills use a three-tier loading strategy for efficient context management:

| Tier         | What is Loaded              | When            | Token Cost               |
| ------------ | --------------------------- | --------------- | ------------------------ |
| Catalog      | name + description          | Session start   | ~50-100 tokens           |
| Instructions | Full SKILL.md body          | Skill activated | <5000 tokens recommended |
| Resources    | Scripts, references, assets | When referenced | Varies                   |

## Creating a New Skill 

Follow these six steps to create a skill:

### Step 1: Determine Skill Purpose

Before creating files, clearly define:
- What capability does this skill provide?
- When should an agent use this skill?
- What keywords would trigger this skill?

### Step 2: Create Directory Structure

```
skill-name/
├── SKILL.md          # Required
├── scripts/          # Optional: executable code
├── references/       # Optional: documentation
└── assets/           # Optional: templates, resources
```

The directory name must match the skill name.

### Step 3: Write Frontmatter

The SKILL.md file must start with YAML frontmatter containing required fields.

Required fields:
- `name`: Skill identifier (1-64 characters, lowercase letters and hyphens only)
- `description`: What the skill does AND when to use it (1-1024 characters)

Optional fields:
- `license`: License name or file
- `compatibility`: Environment requirements
- `metadata`: Additional key-value pairs
- `allowed-tools`: Pre-approved tools (experimental)

For detailed field specifications, see [references/frontmatter.md](references/frontmatter.md).

### Step 4: Write Instruction Content

The Markdown body after frontmatter contains skill instructions. There are no format restrictions, but include:
- When to use this skill
- Step-by-step instructions
- Examples of inputs and outputs
- Common edge cases

For writing guidelines, see [references/content.md](references/content.md).

### Step 5: Add Optional Resources

If your skill requires scripts, references, or assets, add them to the appropriate subdirectories:

- `scripts/`: Executable code that the skill instructions may invoke
- `references/`: Supplementary documentation
- `assets/`: Templates, data files, or other resources

For script support details, see [references/scripts/overview.md](references/scripts/overview.md).

### Step 6: Validate

Verify your skill follows the specification:

```bash
npx skills-ref validate ./skill-name
```

For validation details, see [references/spec.md](references/spec.md).

## Quick Reference

### Valid Name Patterns

```
Valid:   pdf-processing, data-analysis, code-review
Invalid: PDF-Processing (uppercase), -pdf (starts with hyphen), pdf--processing (consecutive hyphens)
```

### Good Description Example

```yaml
description: Extracts text and tables from PDF files, fills PDF forms, and merges multiple PDFs. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction.
```

## Additional Resources

- Specification: [references/spec.md](references/spec.md)
- Frontmatter fields: [references/frontmatter.md](references/frontmatter.md)
- Content writing: [references/content.md](references/content.md)
- Script support overview: [references/scripts/overview.md](references/scripts/overview.md)
- Tool runners: [references/scripts/tool-runner.md](references/scripts/tool-runner.md)
- Python scripts: [references/scripts/python.md](references/scripts/python.md)
- Deno scripts: [references/scripts/deno.md](references/scripts/deno.md)
- Bun scripts: [references/scripts/bun.md](references/scripts/bun.md)
- Examples: [references/examples/](references/examples/)
