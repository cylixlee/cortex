# Content Writing Guide

This guide covers how to write effective SKILL.md instruction content. The Markdown body after frontmatter contains the actual skill instructions.

## Complete Template

A full SKILL.md file combines frontmatter with content body:

```mdx
---
name: skill-name
description: A description of what this skill does and when to use it.
license: Apache-2.0
metadata:
  author: your-name
  version: "1.0"
---

# Skill Title

## When to Use This Skill
Describe specific trigger conditions and keywords...

## Prerequisites
What is needed before starting (optional)...

## How to [Action]
Step-by-step instructions...

## Examples

### Example 1: [Scenario]
Input:
```

Some input```

Output:
```

Expected output```

### Example 2: [Another Scenario]
...

## Common Edge Cases
- Situation 1: How to handle...
- Situation 2: What to do if...

## Error Handling
What to do when things go wrong...

## References
Link to additional documentation...
```

For more detailed structure guidance, see the section below.

## Writing Principles

### Clarity Over Brevity

Write clear, unambiguous instructions. Agents cannot ask clarifying questions, so anticipate edge cases and provide specific guidance.

### Include Trigger Conditions

Start by explaining when this skill should be activated. Help the agent understand:
- What types of tasks trigger this skill
- What keywords or phrases indicate this skill is relevant
- What the skill should NOT be used for

### Provide Step-by-Step Instructions

For complex skills, break down the process into numbered steps. This makes it easier for agents to follow and execute correctly.

### Include Examples

Show concrete examples of:
- What input looks like
- What output should look like
- Common variations or edge cases

### Document Edge Cases

Anticipate common problems or unusual scenarios. Explain how to handle them.

## Recommended Structure

A well-structured SKILL.md typically includes these sections:

## Length Guidelines

- Keep SKILL.md under 5000 tokens when activated
- Move detailed reference material to separate files in `references/` directory
- Use the main SKILL.md for essential instructions, not exhaustive documentation

## Content Tips

**Do:**
- Use clear, direct language
- Include specific commands or code snippets
- Specify exact file formats or naming conventions
- Add warnings about common mistakes

**Don't:**
- Use vague language like "do it properly"
- Assume context the agent might not have
- Include unnecessary backstory or motivation
- Over-explain basics the agent already knows

## Progressive Disclosure in Practice

Structure your content to support progressive disclosure:

1. **Catalog level** (description): What the skill does and when to use it
2. **Instructions level**: Essential steps to complete the task
3. **Resources level**: Detailed reference material in separate files

Example of linking to resources:

```markdown
For detailed guidance on X, see [references/x-details.md](references/x-details.md).
```
