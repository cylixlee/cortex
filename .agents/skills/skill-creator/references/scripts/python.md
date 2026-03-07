# Python Scripts

Python scripts can declare dependencies inline using the PEP 723 format. This allows scripts to be self-contained and run without a separate requirements file.

## PEP 723 Format

The script must start with a special comment block that declares dependencies:

```python
#!/usr/bin/env python3
# /// script
# dependencies = [
#   "beautifulsoup4",
#   "requests>=2.28.0",
# ]
# ///

import requests
from bs4 import BeautifulSoup
# ... rest of script
```

## How It Works

The `# /// script` marker begins the dependency declaration block, and `# ///` ends it. Tools like uv and pip can read this header and create an isolated environment with the specified dependencies before running the script.

## Running Scripts

With uv:

```bash
uv run script.py
```

With pipx:

```bash
pipx run script.py
```

## Complete Example

```python
#!/usr/bin/env python3
"""Extracts text content from HTML files."""

# /// script
# dependencies = [
#   "beautifulsoup4>=4.12.0",
# ]
# ///

from pathlib import Path
import sys
from bs4 import BeautifulSoup


def extract_text(html_path: Path) -> str:
    with open(html_path, "r", encoding="utf-8") as f:
        soup = BeautifulSoup(f.read(), "html.parser")
    return soup.get_text()


def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <html-file>", file=sys.stderr)
        sys.exit(1)

    html_file = Path(sys.argv[1])
    if not html_file.exists():
        print(f"Error: File '{html_file}' not found.", file=sys.stderr)
        sys.exit(1)

    text = extract_text(html_file)
    print(text)


if __name__ == "__main__":
    main()
```

## Key Points

1. **Shebang line**: Include `#!/usr/bin/env python3` for direct execution
2. **Dependency format**: Use pip-style specifiers (package, package>=version, etc.)
3. **No import outside block**: Dependencies declared in the block are available globally
4. **Isolation**: Each run uses an isolated environment

## Best Practices

- Always include a docstring describing what the script does
- Use argparse for command-line argument parsing
- Exit with appropriate error codes (0 for success, non-zero for errors)
- Print errors to stderr, normal output to stdout
- Handle missing files and invalid input gracefully
