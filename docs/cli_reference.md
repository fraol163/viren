# Viren CLI Reference: The Manual

Viren is a dual-mode application. It functions both as an immersive, interactive shell and as a traditional Unix-style CLI utility. This document provides the exhaustive list of flags, environment variables, and piping patterns.

---

## 1. Global Flags (The Binary)

| Flag | Argument | Functional Result |
| :--- | :--- | :--- |
| `-h, --help` | None | Displays ASCII branding, command table, and core capabilities. |
| `-v, --version` | None | Outputs semantic version, build timestamp, and git commit ID. |
| `-c, --continue` | None | Resumes from the most recent session in `~/.viren/tmp/`. |
| `-a, --history` | `exact` | Opens the fzf history manager. `exact` disables fuzzy matching. |
| `-p, --platform` | `name` | Overrides the default platform (e.g., `viren -p anthropic`). |
| `-m, --model` | `name` | Overrides the default model (e.g., `viren -m gpt-4o`). |
| `-o, --all` | `p\|m` | Shortcut for both (e.g., `viren -o "groq\|llama3"`). |
| `-nh, --no-history` | None | Runs Viren without writing anything to the local history database. |

---

## 2. Ingestion Flags (Direct Mode)

### Load File (`-l`)
Ingests files or URLs as context for a single query.
- **Syntax**: `viren -l <path> "your prompt"`
- **Example**: `viren -l error.log "What is the root cause?"`
- **Piping Example**: `cat data.csv | viren -l schema.sql "Generate an import script"`

### Codedump (`-d`)
Bundles a directory for a project-wide query.
- **Syntax**: `viren -d <dir_path> "your prompt"`
- **Example**: `viren -d ./internal "Refactor these modules for better concurrency"`

### Web search (`-w`)
Searches the live web via Brave before answering.
- **Syntax**: `viren -w "query" "instruction"`
- **Example**: `viren -w "Go 1.24 release" "What are the new crypto features?"`

---

## 3. Environment Variables

Viren looks for the following variables to manage logic and security:

- **`OPENAI_API_KEY`**: Required for OpenAI models.
- **`ANTHROPIC_API_KEY`**: Required for Claude models.
- **`GEMINI_API_KEY`**: Required for Google models.
- **`DEEP_SEEK_API_KEY`**: Required for DeepSeek models.
- **`BRAVE_API_KEY`**: Required for the `!w` (Web Search) command.
- **`VIREN_DEFAULT_PLATFORM`**: Sets the default provider on startup.
- **`VIREN_DEFAULT_MODEL`**: Sets the default model on startup.
- **`EDITOR`**: Defines which editor Viren opens for `!e` and `!t`.

---

## 4. Piping & Redirects (The Unix Way)

Viren is a "First-Class Citizen" of the Unix pipe.

### Reading from Stdin
You can pipe any text into Viren. It will treat the piped text as context and your argument as the instruction.
```bash
git diff HEAD~1 | viren "Write a detailed commit message for these changes"
```

### Writing to Stdout
By default, Viren streams to the terminal. You can redirect this to files or other tools.
```bash
viren "Write a python script to scrape news" > scraper.py
viren "List 10 colors" | grep "Blue"
```

---

## 5. Exit Codes
- `0`: Success.
- `1`: Error (Network, API, or File System).
- `130`: Termination by User (Ctrl+C).

**Viren is designed to be the ultimate automation partner.**
