# Skill with Scripts Example

A skill that includes executable scripts for complex operations.

```mdx
---
name: csv-analytics
description: Performs analytics on CSV files including aggregation, filtering, and statistical analysis. Use when working with CSV data, generating reports, or analyzing tabular data.
---

# CSV Analytics

## When to Use This Skill

- Working with CSV files
- Generating data reports
- Performing statistical analysis on tabular data
- User asks to analyze CSV or generate insights

## Prerequisites

Python 3.10+ with uv package manager installed.

## How to Analyze CSV

### Basic Statistics

Run the analytics script:

```bash
uv run scripts/analyze.py --input data.csv --summary
```

### Generate Report

```bash
uv run scripts/analyze.py --input data.csv --report report.json
```

### Filter Data

```bash
uv run scripts/analyze.py --input data.csv --filter "column > 100"
```

## Output Formats

The script supports JSON, CSV, and HTML output:
- `--format json`: Machine-parseable JSON
- `--format csv`: CSV with aggregated data
- `--format html`: HTML table for visualization

## Related Scripts

- `scripts/analyze.py`: Main analytics engine
- `scripts/transform.py`: Data transformation utilities
- `scripts/visualize.py`: Generate charts from data
```

## Directory Structure

```
csv-analytics/
├── SKILL.md
└── scripts/
    ├── analyze.py      # Main analytics script (PEP 723)
    ├── transform.py    # Data transformation
    └── visualize.py   # Visualization generation
```

## Key Characteristics

This example demonstrates:
- Using tool runners (uvx, uv run) in instructions
- Referencing scripts in the scripts/ directory
- Multiple script use cases (analyze, transform, visualize)
- Command-line argument documentation
- Output format options
