# Scripts Support Overview

Skills can include executable scripts that agents can run to perform specialized tasks. This document provides an overview; see related files for details.

## When to Use Scripts

Include scripts in your skill when:
- The task requires executing code (compiling, linting, testing, etc.)
- You want to provide reusable, automated tooling
- Complex commands need to be encapsulated in executable files

## Script Directories

Place scripts in the `scripts/` subdirectory of your skill:

```
skill-name/
├── SKILL.md
└── scripts/
    ├── script1.py
    ├── script2.sh
    └── ...
```

## Types of Script Support

### 1. Tool Runners

Use tool runners to execute one-off commands without creating script files. See [tool-runner.md](tool-runner.md) for details.

| Tool     | Description                                   | Command Example                                     |
| -------- | --------------------------------------------- | --------------------------------------------------- |
| uvx      | Runs Python packages in isolated environments | `uvx ruff@0.8.0 check .`                            |
| pipx     | Runs Python packages in isolated environments | `pipx run 'black==24.10.0' .`                       |
| npx      | Runs npm packages                             | `npx eslint@9 --fix .`                              |
| bunx     | Bun's equivalent of npx                       | `bunx eslint@9 --fix .`                             |
| deno run | Runs scripts directly from URLs               | `deno run npm:create-vite@6 my-app`                 |
| go run   | Compiles and runs Go packages                 | `go run golang.org/x/tools/cmd/goimports@v0.28.0 .` |

### 2. Self-Contained Scripts

Several languages support inline dependency declarations, allowing scripts to declare their own dependencies:

- **Python**: PEP 723 format (see [python.md](python.md))
- **Deno**: Native dependency imports (see [deno.md](deno.md))
- **Bun**: Native dependency imports (see [bun.md](bun.md))

## Script Design Best Practices

### Avoid Interactive Prompts

Agents operate in non-interactive shells. Scripts must:
- Accept all necessary input via arguments, environment variables, or files
- Never wait for user input
- Exit with appropriate error codes

### Document Usage

Include `--help` output that provides:
- Brief description of what the script does
- All available flags and arguments
- Example usage

```python
#!/usr/bin/env python3
"""Description of what this script does."""

import argparse

def main():
    parser = argparse.ArgumentParser(description="Description")
    parser.add_argument("-i", "--input", required=True, help="Input file path")
    parser.add_argument("-o", "--output", help="Output file path")
    # ...
```

### Write Helpful Error Messages

When errors occur, say:
- What went wrong
- What was expected
- What to try next

```python
if not input_file.exists():
    print(f"Error: Input file '{input_file}' not found.", file=sys.stderr)
    print("Use --input to specify the input file.", file=sys.stderr)
    sys.exit(1)
```

### Use Structured Output

Prefer machine-parseable formats:
- JSON for complex data
- CSV or TSV for tabular data
- Plain text only when truly appropriate

### Separate Data from Diagnostics

- Send structured data to stdout
- Send human-readable diagnostics to stderr

## Invoking Scripts from SKILL.md

When your SKILL.md instructs the agent to run a script, use relative paths from the skill directory:

```markdown
To extract text from a PDF, run:

    python scripts/extract.py --input document.pdf --output output.txt
```

## Related Documents

- Tool runners: [tool-runner.md](tool-runner.md)
- Python scripts: [python.md](python.md)
- Deno scripts: [deno.md](deno.md)
- Bun scripts: [bun.md](bun.md)
