# Agent Skills Specification

This document provides a comprehensive overview of the Agent Skills specification. For detailed information on specific topics, see the related references.

## What is Agent Skills

Agent Skills is a lightweight, open format for extending AI agent capabilities with specialized knowledge and workflows. Originally developed by Anthropic and released as an open standard, it has been adopted by many AI development tools.

### Key Characteristics

- **Portable**: Skills are just files, easy to edit, version, and share
- **Self-documenting**: Anyone can read a SKILL.md to understand what a skill does
- **Extensible**: Skills range from simple text instructions to executable code, assets, and templates
- **Interoperable**: Same skill works across different skills-compatible agent products

### Who Benefits

- **Skill Authors**: Build capabilities once, deploy across multiple agent products
- **Compatible Agents**: Support for skills lets end users give agents new capabilities
- **Teams and Enterprises**: Capture organizational knowledge in portable, version-controlled packages

## Directory Structure

A skill is a directory containing at minimum a `SKILL.md` file:

```
skill-name/
└── SKILL.md          # Required
```

Optional directories:

```
my-skill/
├── SKILL.md          # Required: instructions + metadata
├── scripts/          # Optional: executable code
├── references/       # Optional: documentation
└── assets/           # Optional: templates, resources
```

## Progressive Disclosure

Skills use a three-tier loading strategy to manage context efficiently:

| Tier         | What's Loaded               | When                             | Token Cost               |
| ------------ | --------------------------- | -------------------------------- | ------------------------ |
| Catalog      | Name + description          | Session start                    | ~50-100 tokens per skill |
| Instructions | Full SKILL.md body          | When skill is activated          | <5000 tokens recommended |
| Resources    | Scripts, references, assets | When instructions reference them | Varies                   |

This approach keeps agents fast while giving them access to specialized knowledge on demand.

### How Skills Work

1. **Discovery**: At startup, agents load only the name and description of each available skill
2. **Activation**: When a task matches a skill's description, the agent reads the full SKILL.md instructions into context
3. **Execution**: The agent follows the instructions, optionally loading referenced files or executing bundled code as needed

## Discovery Locations

Agents typically scan multiple scopes:

| Scope   | Path                          | Purpose                       |
| ------- | ----------------------------- | ----------------------------- |
| Project | `<project>/.agents/skills/`   | Cross-client interoperability |
| Project | `<project>/.<client>/skills/` | Client native location        |
| User    | `~/.agents/skills/`           | Cross-client interoperability |
| User    | `~/.agents/<client>/skills/`  | Client native location        |

### What to Scan For

Within each skills directory, look for subdirectories containing a file named exactly `SKILL.md`:

```
~/.agents/skills/
├── pdf-processing/
│   ├── SKILL.md          # discovered
│   └── scripts/
│       └── extract.py
├── data-analysis/
│   └── SKILL.md          # discovered
└── README.md             # ignored (not a skill directory)
```

### Discovery Steps

1. Find directories with SKILL.md: Look for subdirectories containing a file named exactly `SKILL.md`
2. Skip common non-skill directories: `.git/`, `node_modules/`, etc.
3. Handle name collisions: Project-level skills override user-level skills
4. Trust considerations: Consider gating project-level skill loading on a trust check

## Activation Mechanisms

**Model-driven activation:**
- File-read activation: Model calls file-read tool with the SKILL.md path
- Dedicated tool activation: Register a tool (e.g., `activate_skill`) that takes a skill name

**User-explicit activation:**
- Slash commands (`/skill-name`) or mention syntax (`$skill-name`)

## Structured Wrapping Example

```xml
<skill_content name="pdf-processing">
# PDF Processing

## When to use this skill
Use this skill when the user needs to work with PDF files...

[rest of SKILL.md body]

Skill directory: /home/user/.agents/skills/pdf-processing
Relative paths in this skill are relative to the skill directory.

<skill_resources>
  <file>scripts/extract.py</file>
  <file>scripts/merge.py</file>
  <file>references/pdf-spec-summary.md</file>
</skill_resources>
</skill_content>
```

## Validation

Use the skills-ref reference library to validate skills:

```bash
npx skills-ref validate ./my-skill
```

This checks that your SKILL.md frontmatter is valid and follows all naming conventions.

## Supported Tools and Platforms

The Agent Skills specification is supported by numerous AI development tools:

| Category        | Tools                                                              |
| --------------- | ------------------------------------------------------------------ |
| IDEs            | VS Code, Cursor, Juno                                              |
| CLI Agents      | Claude Code, Claude CLI, Gemini CLI, OpenCode, OpenHands, Roo Code |
| Cloud Platforms | Mux (Coder), Letta, Databricks, Snowflake                          |
| Enterprise      | GitHub Copilot, Laravel Boost, Spring AI                           |
