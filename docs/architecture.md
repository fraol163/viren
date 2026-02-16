# Viren: Full Technical Architecture

Viren is built for speed, modularity, and extensibility. This document provides a deep-dive into the internal mechanics of the binary, explaining how we achieve sub-100ms startup times while maintaining complex state and multi-platform support.

## 1. High-Level Design Philosophy

Viren follows the **Unix Philosophy**: small, sharp, and composable.
- **Concurrency**: We use Go's `goroutines` to ensure that network I/O never blocks the UI.
- **Static Typing**: Every data structure (from a single chat message to a full system configuration) is strictly typed in `pkg/types/types.go`.
- **Stateless Networking**: The `platform` module treats each API request as a unique event, allowing for seamless switching between models mid-conversation.

---

## 2. Core Module Breakdown

### A. The Command Interface (`cmd/viren/main.go`)
The "Bootstrap" layer.
- **Execution Flow**: It first parses CLI flags (like `-p` or `-m`). If a direct query is provided (e.g., `viren "hello"`), it executes in "Single-Shot" mode. If not, it enters the "Interactive Loop."
- **Onboarding**: It checks for the existence of `~/.viren/config.json`. If missing, it triggers the `config.RunOnboarding` wizard.
- **Signal Control**: It manages a global `context.CancelFunc` to ensure that hitting `Ctrl+C` kills the active AI stream without terminating the whole application.

### B. The Chat Manager (`internal/chat/`)
The "Orchestrator."
- **History Management**: It maintains an in-memory slice of `ChatMessage` structs. It is responsible for appending the "System Prompt" (which includes the user's role and the current mode) to every request.
- **Modes & Personalities**: This module contains the logic for the 40+ domain modes. When you switch to `!v aerospace`, the Chat Manager re-calculates the entire system prompt and resets the context window.
- **Export Logic**: It contains the regex-based detection for code blocks and the intelligent filename generator that scans code for keywords (like `package` or `func`) to suggest names.

### C. The Platform Manager (`internal/platform/`)
The "Translator."
- **Standardization**: Every AI provider uses a slightly different JSON schema. The Platform Manager abstracts this away. It takes a Viren `ChatRequest` and converts it into an `OpenAIRequest`, `AnthropicRequest`, etc.
- **Streaming (SSE)**: It implements a custom Server-Sent Events parser. This allows Viren to render the AI's response character-by-character as it arrives, rather than waiting for the entire block.
- **Authentication**: It fetches keys from the environment variables, ensuring that no secrets are ever persisted in the binary.

### D. The UI Engine (`internal/ui/`)
The "Renderer."
- **ANSI Engine**: We avoid heavy TUI libraries (like Bubbletea) to keep the binary size small and the startup fast. Instead, we use raw ANSI escape codes for coloring, borders, and animations.
- **Fzf Bridge**: Viren spawns `fzf` as a subprocess for all selection menus. This provides a lightning-fast, searchable interface that users already know and love.
- **OCR Logic**: Utilizes `tesseract` bindings (if CGO is enabled) to convert images loaded via `!l` into text context.

### E. The Config System (`internal/config/`)
The "Persistence" layer.
- **Schema Management**: It handles the merging of the `DefaultConfig` with the user's `config.json`.
- **Safety**: It implements the "Shallow Load" logic, preventing the application from scanning system-critical folders.

---

## 3. The Lifecycle of a Prompt

To understand Viren, follow a single prompt from input to output:

1.  **Input**: User types "Refactor this code" and presses `Enter`.
2.  **Bang Detection**: The `ChatManager` checks if the input starts with `!`. If it doesn't, it proceeds.
3.  **Context Enrichment**: 
    - The `ChatManager` looks at the `Messages` slice.
    - It prepends the `SystemPrompt` (calculated from Mode + Personality + UserProfile).
    - It appends any files previously loaded via `!l`.
4.  **Dispatch**: The `PlatformManager` picks the current platform (e.g., `deepseek`).
5.  **Transformation**: The prompt is wrapped in the DeepSeek-specific JSON format.
6.  **Stream**: An HTTP POST is opened. As chunks arrive, the `UIEngine` renders them.
7.  **Auto-Command**: If the response contains a ` ```bash ` block, the `UIEngine` triggers a "DETECTED COMMAND" prompt.
8.  **Finalize**: Once the stream ends, the `ChatManager` saves the turn to `~/.viren/tmp/` for future sessions.

---

## 4. Performance Optimizations

### Why is it so fast?
1.  **Go standard library**: We use `net/http` and `encoding/json` almost exclusively.
2.  **Lazy Loading**: We don't load the history database or the 40+ modes into memory until the user actually requests them.
3.  **No Interpreter**: Unlike Python or Node.js tools, Viren is a compiled static binary. There is no runtime to start up.
4.  **Subprocess fzf**: By delegating search to `fzf`, we get world-class performance without writing complex search algorithms in Go.

---

## 5. Security Architecture

- **Environment-Only Keys**: Viren does not store API keys in its config file. This prevents "Key Leakage" if you share your `config.json`.
- **Local History**: History files are stored with `0600` permissions (read/write only for the current user).
- **Process Isolation**: When running shell commands via `!x`, they run in a standard sub-shell with the user's existing permissions.

**Viren is built to be the most technically robust AI CLI available.**
