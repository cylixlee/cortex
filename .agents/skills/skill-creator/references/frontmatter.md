# Frontmatter Field Reference

The SKILL.md file must contain YAML frontmatter followed by Markdown content. This document details all available frontmatter fields.

## Required Fields

### name

The skill identifier. Must satisfy all of the following:

- 1-64 characters
- Unicode lowercase alphanumeric characters and hyphens only (`a-z` and `-`)
- Must not start or end with `-`
- Must not contain consecutive hyphens (`--`)
- Must match the parent directory name

**Valid examples:**

```yaml
name: pdf-processing
name: data-analysis
name: code-review
```

**Invalid examples:**

```yaml
name: PDF-Processing  # uppercase not allowed
name: -pdf           # cannot start with hyphen
name: pdf--processing  # cannot have consecutive hyphens
```

### description

A description of what the skill does AND when to use it. Must be 1-1024 characters.

The description should include specific keywords that help agents identify relevant tasks.

**Good example:**

```yaml
description: Extracts text and tables from PDF files, fills PDF forms, and merges multiple PDFs. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction.
```

**Poor example:**

```yaml
description: Helps with PDFs.
```

## Optional Fields

### license

Specifies the license applied to the skill. Can be either the name of a license or the name of a bundled license file.

```yaml
license: Apache-2.0
license: LICENSE.txt
```

### compatibility

Indicates specific environment requirements. Use only if your skill has special requirements.

Can indicate:
- Intended product (e.g., "Claude Code")
- Required system packages
- Network access needs
- Other environment-specific needs

```yaml
compatibility: Requires Python 3.10+ and pdftotext CLI tool
compatibility: Designed for Claude Code with uv package manager
```

### metadata

A map from string keys to string values. Clients can use this to store additional properties not defined by the Agent Skills specification.

```yaml
metadata:
  author: example-org
  version: "1.0"
  tags: ["pdf", "document", "extraction"]
  created: "2024-01-15"
```

### allowed-tools

A space-delimited list of tools that are pre-approved to run. This field is experimental and support may vary between agent implementations.

```yaml
allowed-tools: Bash(git:*) Bash(jq:*) Read
```

## SKILL.md Template

For a complete SKILL.md template with frontmatter and body structure, see [content.md](content.md).

## Validation

Run validation to check frontmatter:

```bash
npx skills-ref validate ./skill-name
```

Common validation errors:
- Missing required `name` or `description` field
- Invalid name format (uppercase, special characters, wrong length)
- Description too short to be useful
