# Deno Scripts

Deno scripts have native support for dependency imports. Unlike Node.js, Deno imports directly from URLs, making scripts self-contained.

## Native Import Syntax

Deno scripts import packages directly from URLs:

```typescript
import * as cheerio from "npm:cheerio@1.0.0";
import { readFileStr } from "https://deno.land/std@0.208.0/fs/mod.ts";
```

## Shebang for Direct Execution

Include a shebang line to run the script directly:

```typescript
#!/usr/bin/env -S deno run

import * as cheerio from "npm:cheerio@1.0.0";
// ... rest of script
```

## Running Scripts

```bash
deno run script.ts
deno run --allow-read script.ts
deno run -A https://example.com/script.ts
```

Flags:
- `--allow-read`: Allow file system read operations
- `--allow-write`: Allow file system write operations
- `-A` or `--allow-all`: Allow all permissions

## Complete Example

```typescript
#!/usr/bin/env -S deno run --allow-read

import * as cheerio from "npm:cheerio@1.0.0";

interface ExtractOptions {
  selector: string;
  attribute?: string;
}

function extractHtml(filePath: string, options: ExtractOptions): string {
  const html = Deno.readTextFileSync(filePath);
  const $ = cheerio.load(html);
  
  if (options.attribute) {
    return $(options.selector).attr(options.attribute) || "";
  }
  
  return $(options.selector).text();
}

function main() {
  const args = Deno.args;
  
  if (args.length < 2) {
    console.error(`Usage: ${Deno.mainModule} <html-file> <selector> [--attribute <attr>]`);
    Deno.exit(1);
  }
  
  const [file, selector] = args;
  const attributeFlag = args.indexOf("--attribute");
  const attribute = attributeFlag > -1 ? args[attributeFlag + 1] : undefined;
  
  try {
    const result = extractHtml(file, { selector, attribute });
    console.log(result);
  } catch (error) {
    console.error(`Error: ${error.message}`);
    Deno.exit(1);
  }
}

main();
```

## NPM Packages

Deno can run npm packages using the `npm:` prefix:

```typescript
import express from "npm:express@4.18.2";
import lodash from "npm:lodash@4.17.21";
```

## Best Practices

1. **Permissions**: Request only the permissions your script needs
   ```typescript
   // Request specific permissions at runtime
   const read = await Deno.permissions.query({ name: "read", path: "." });
   ```

2. **Error handling**: Always handle errors and exit with appropriate codes

3. **Output**: Print errors to stderr, data to stdout

4. **TypeScript**: Use TypeScript for better error prevention and documentation

5. **Imports**: Pin versions for reproducibility
   ```typescript
   // Good
   import express from "npm:express@4.18.2";
   
   // Avoid (may break with updates)
   import express from "npm:express";
   ```
