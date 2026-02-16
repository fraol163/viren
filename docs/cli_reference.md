# Viren CLI Reference: The Exhaustive Manual

Viren is a high-performance, dual-mode application. It functions as both an immersive, interactive shell and a traditional Unix-style command-line utility. This document provides the exhaustive list of flags, commands, environment variables, and piping patterns.

---

## 1. Binary Invocation & Global Flags

The `viren` binary supports the following flags when launched from your shell:

### Metadata & Info
- `-h, --help`: Displays the full help dashboard, including the ASCII logo and a summary of interactive commands.
- `-v, --version`: Outputs the semantic version (e.g., `v1.0.0`), the build timestamp, and the git commit hash.

### Session Management
- `-c, --continue`: Automatically resumes the most recent conversation stored in `~/.viren/tmp/`.
- `-a, --history`: Opens the interactive history manager. 
    - **Argument**: Adding `exact` (e.g., `viren -a exact`) disables fuzzy matching for session titles.
- `-nh, --no-history`: Prevents Viren from writing the current session to the local history database. Use this for highly sensitive or one-off queries.

### Logic & Platform Overrides
- `-p, --platform <name>`: Forces Viren to start with a specific provider (e.g., `viren -p groq`).
- `-m, --model <name>`: Forces Viren to start with a specific model (e.g., `viren -m claude-3-opus`).
- `-o, --all <p|m>`: A shorthand format to set both at once (e.g., `viren -o "openai|gpt-4o"`).

---

## 2. Ingestion Flags (Direct Execution Mode)

Direct mode allows you to use Viren as a single-shot utility. The output is printed to `stdout` and the program exits.

### Load File (`-l`)
Injects files or URLs as context for a single query.
- **Usage**: `viren -l <file_path> "prompt"`
- **Supported Formats**: `.txt`, `.go`, `.py`, `.js`, `.json`, `.pdf`, `.docx`, `.xlsx`, `.csv`.
- **Example**: `viren -l main.go "Explain the concurrency logic here"`

### Codedump (`-d`)
Bundles a directory into a structured context stream.
- **Usage**: `viren -d <dir_path> "prompt"`
- **Example**: `viren -d ./internal "Find all potential memory leaks"`

### Web Search (`-w`)
Searches the live web via Brave Search before answering.
- **Usage**: `viren -w "query" "instruction"`
- **Example**: `viren -w "latest rust releases" "Should I update my project?"`

---

## 3. The Power of Unix Pipes

Viren is a "First-Class Citizen" of the Unix pipeline. It treats `stdin` as context and its arguments as instructions.

### Piping Context IN
```bash
# Debug a log file
tail -n 50 /var/log/syslog | viren "What is the cause of these errors?"

# Explain a Git diff
git diff HEAD~1 | viren "Summarize these changes for a technical blog post"

# Ingest data
cat users.json | viren "Convert this to a SQL insert script"
```

### Piping Response OUT
```bash
# Generate and save a script
viren "Write a bash script to clean up Docker images" > cleanup.sh

# Count tokens in a response
viren "Write a short story" | wc -w
```

---

## 4. Environment Variables

Viren looks for these variables to manage logic and security without touching your `config.json`.

- **`OPENAI_API_KEY`**: Required for all OpenAI models.
- **`ANTHROPIC_API_KEY`**: Required for Claude models.
- **`DEEP_SEEK_API_KEY`**: Required for DeepSeek models.
- **`GEMINI_API_KEY`**: Required for Google Gemini models.
- **`GROQ_API_KEY`**: Required for Groq models.
- **`BRAVE_API_KEY`**: Required for the `!w` command.
- **`VIREN_DEFAULT_PLATFORM`**: Overrides the starting provider.
- **`VIREN_DEFAULT_MODEL`**: Overrides the starting model.
- **`EDITOR`**: Defines the binary used for `!e` (Export) and `!t` (Editor) modes.

---

## 5. Interactive Commands (Bang Commands)

Once inside the Viren shell, use these commands to control the application state:

| Command | functional definition |
| :--- | :--- |
| `!q` | **Quit**: Exits and serializes history. |
| `!h` | **Help**: Opens the fzf command dashboard. |
| `!c` | **Clear**: Resets the context and clears the screen. |
| `!m` | **Model**: Searchable menu of available LLMs. |
| `!p` | **Platform**: Searchable menu of AI providers. |
| `!u` | **Personality**: Switch between tone templates (Creative, Focused, etc.). |
| `!v` | **Domain Mode**: Apply specialized system prompts (Zenith, Code Whisperer). |
| `!z` | **Theme**: Instant ANSI color palette switch. |
| `!x` | **Shell Record**: Ingest terminal output for debugging. |
| `!d` | **Codedump**: Bundle your project directory for context. |
| `!l` | **Load**: Select files/URLs to inject into context. |
| `!s` | **Scrape**: Extract clean text from a URL. |
| `!w` | **Web Search**: Perform a live Brave Search. |
| `!a` | **History**: Interactively browse and restore sessions. |
| `!y` | **Clipboard**: copy responses to system clipboard. |
| `!b` | **Backtrack**: Revert the last N turns of conversation. |
| `!t` | **Editor**: Open long-form input in your preferred text editor. |

---

## 6. Exit Codes
Viren uses standard exit codes for automation scripting:
- `0`: Successful execution.
- `1`: General Error (Configuration missing, API failure).
- `126`: Command invoked cannot execute.
- `130`: Terminated by User (Ctrl+C).

**Viren is designed to be the ultimate automation partner for the modern engineer.**
