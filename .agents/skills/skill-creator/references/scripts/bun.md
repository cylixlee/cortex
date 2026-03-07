# Bun Scripts

Bun scripts support native dependency imports similar to Deno, with a focus on speed and compatibility with the JavaScript ecosystem.

## Native Import Syntax

Bun scripts can import packages using various sources:

```typescript
import * as cheerio from "cheerio@1.0.0";
import { glob } from "bun:glob";
import { readFile } from "bun:fs";
```

## Shebang for Direct Execution

Include a shebang line to run the script directly:

```typescript
#!/usr/bin/env bun

import { echo } from "bun:shell";
// ... rest of script
```

## Running Scripts

```bash
bun run script.ts
./script.ts  # if shebang is set
```

## Complete Example

```typescript
#!/usr/bin/env bun

import { readFileSync, writeFileSync, existsSync } from "bun:fs";
import { glob } from "bun:glob";

interface MinifyOptions {
  removeComments: boolean;
  collapseWhitespace: boolean;
}

function minifyHtml(content: string, options: MinifyOptions): string {
  let result = content;
  
  if (options.removeComments) {
    result = result.replace(/<!--[\s\S]*?-->/g, "");
  }
  
  if (options.collapseWhitespace) {
    result = result.replace(/\s+/g, " ");
  }
  
  return result;
}

function main() {
  const args = Bun.argv.slice(2);
  
  if (args.length < 1) {
    console.error(`Usage: ${Bun.mainModule} <input-file> [--output <output-file>]`);
    process.exit(1);
  }
  
  const inputFile = args[0];
  const outputFlag = args.indexOf("--output");
  const outputFile = outputFlag > -1 ? args[outputFlag + 1] : inputFile;
  
  if (!existsSync(inputFile)) {
    console.error(`Error: File '${inputFile}' not found.`);
    process.exit(1);
  }
  
  try {
    const content = readFileSync(inputFile, "utf-8");
    const minified = minifyHtml(content, {
      removeComments: true,
      collapseWhitespace: true
    });
    writeFileSync(outputFile, minified);
    console.log(`Minified ${inputFile} -> ${outputFile}`);
  } catch (error) {
    console.error(`Error: ${error.message}`);
    process.exit(1);
  }
}

main();
```

## Built-in Modules

Bun provides built-in modules for common operations:

| Module      | Description             |
| ----------- | ----------------------- |
| `bun:fs`    | File system operations  |
| `bun:path`  | Path manipulation       |
| `bun:shell` | Shell command execution |
| `bun:glob`  | Glob pattern matching   |
| `bun:http`  | HTTP server/client      |

## Best Practices

1. **Use built-in modules**: Prefer Bun's built-in modules for performance
   ```typescript
   // Good
   import { readFileSync } from "bun:fs";
   
   // Avoid (adds Node.js compatibility overhead)
   import { readFileSync } from "fs";
   ```

2. **Version pins**: Pin package versions for reproducibility
   ```typescript
   // Good
   import cheerio from "cheerio@1.0.0";
   ```

3. **Error handling**: Exit with appropriate error codes

4. **Output**: Print errors to stderr, data to stdout

5. **Permissions**: Bun has no runtime permissions model; scripts have full file system access by default
