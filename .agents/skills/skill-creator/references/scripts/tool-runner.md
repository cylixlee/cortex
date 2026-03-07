# Tool Runners

Tool runners execute one-off commands without creating local script files. They are useful for running tools or utilities that are already installed or available as packages.

## uvx

Runs Python packages in isolated environments without requiring local installation.

```bash
uvx ruff@0.8.0 check .
uvx black@latest format .
uvx httpie@0.9.0 POST https://api.example.com/data
```

**Use cases:** Running Python CLI tools, linters, formatters, HTTP clients.

## pipx

Similar to uvx, runs Python packages in isolated environments. Pre-installed on some systems.

```bash
pipx run 'black==24.10.0' .
pipx run 'cookiecutter==1.7.3' my-template
pipx run 'pytest==8.0.0' tests/
```

**Use cases:** Running specific versions of Python tools.

## npx

Runs npm packages from the Node.js ecosystem.

```bash
npx eslint@9 --fix .
npx typescript@latest --noEmit
npx create-react-app@latest my-app
```

**Use cases:** Running Node.js CLI tools, running one-off npm packages.

## bunx

Bun's equivalent of npx, typically faster for Bun projects.

```bash
bunx eslint@9 --fix .
bunx tailwindcss init
bunx prisma generate
```

**Use cases:** Running Node.js tools in Bun environments, faster execution.

## deno run

Runs scripts directly from URLs, supports both Deno and npm packages.

```bash
deno run npm:create-vite@6 my-app
deno run https://example.com/script.ts
deno run -A npm:deno-githooks@1.0.0/hook.ts
```

**Use cases:** Running Deno scripts, running npm packages with Deno.

## go run

Compiles and runs Go packages directly.

```bash
go run golang.org/x/tools/cmd/goimports@latest -w .
go run github.com/google/ko@latest build .
go run github.com/cosmtrek/air@latest
```

**Use cases:** Running Go CLI tools, compiling and executing Go packages.

## Choosing a Tool Runner

| Language | Recommended Tool Runner |
| -------- | ----------------------- |
| Python   | uvx (preferred) or pipx |
| Node.js  | npx or bunx             |
| Deno     | deno run                |
| Go       | go run                  |

## Best Practices

1. **Pin versions**: Specify versions to ensure reproducibility
   ```bash
   # Good
   uvx ruff@0.8.0 check .
   
   # Avoid (may produce different results)
   uvx ruff check .
   ```

2. **Use -y for yes to all**: Some tools prompt for confirmation
   ```bash
   npx -y create-next-app@latest my-app
   ```

3. **Check availability**: Tool runners may not be installed on all systems
   - uvx requires uv
   - bunx requires Bun
   - deno requires Deno

4. **Consider offline usage**: Some tool runners download packages each time; for offline use, consider embedding scripts directly in the skill
