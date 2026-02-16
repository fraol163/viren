# Viren: The Neural Architect's Internal Manual

Viren is a high-performance, concurrent Command Line Interface (CLI) engineered for developers who demand absolute speed, total privacy, and extreme modularity. This document provides an exhaustive, low-level breakdown of the internal mechanics, module interactions, and design decisions that define the Viren ecosystem.

---

## 1. High-Level Design Philosophy: The Go Advantage

The Viren architecture is rooted in the **Unix Philosophy**: build small, sharp tools that communicate via standard interfaces. By choosing GoLang as our primary engine, we achieve several critical architectural goals:

### A. Non-Blocking, Asynchronous Core
Unlike many AI wrappers written in Python or JavaScript, Viren is built in **GoLang**. It leverages Go's native `goroutines` and `channels` to ensure that network I/O, file system scanning, and terminal rendering never block each other. While an API request is in flight, the UI engine continues to run high-refresh animations and monitor for user signals (like `Ctrl+C`).

### B. Static Typing & Data Sovereignty
Every data structure in Viren—from a single chat message fragment to a project-wide project dump—is strictly typed in `pkg/types/types.go`. This "Contract-First" approach ensures that the Platform Manager always sends perfectly valid JSON to third-party APIs, eliminating the "undefined" errors common in dynamic toolchains.

### C. Zero-Runtime Initialization
Viren intentionally avoids heavy TUI (Terminal User Interface) libraries. Instead, we have built a custom ANSI compositor. This allows Viren to bypass the complex "handshake" and screen-redrawing sequences required by frameworks, resulting in a cold-start time that is practically instantaneous (<100ms).

---

## 2. Core Module Deep-Dive

### The Bootstrap Layer: `cmd/viren/main.go`
The primary entry point of the binary.
- **Flag Orchestration**: Uses the Go `flag` package to parse CLI arguments. It distinguishes between **Interactive Mode** (the default shell) and **Single-Shot Mode** (where output is typically redirected to a file or piped to another Unix tool).
- **Graceful Signal Management**: Captures `SIGINT` and `SIGTERM`. If a user hits `Ctrl+C` while the AI is streaming, Viren sends a cancellation signal to the networking context, immediately kills the stream, and returns the user to the prompt instead of crashing the process.
- **Bootstrap Hook**: Checks for the existence of the `~/.viren/` directory structure. If it is a fresh install, it triggers the `config.RunOnboarding` logic before proceeding to the main loop.

### The Brain: `internal/chat/`
The "State Machine" and orchestrator of conversation.
- **Context Bundling Engine**: This is the most technically complex module. It must assemble a multi-part prompt containing:
    1.  **System Mandate**: The base instructions defined by the current Domain Mode (e.g., Zenith).
    2.  **Neural Profile**: The user's name, role, and environmental goals.
    3.  **Injected Context**: Text extracted from loaded files, PDFs, or directory dumps.
    4.  **Short-Term Memory**: The actual message history, trimmed to fit the model's token limit.
- **Logic Matrix**: Contains the registry of 40+ modes. When a user switches modes via `!v`, the Chat Manager re-calculates the entire system prompt and triggers a "Context Refresh" event.
- **Smart Export Engine**: Uses advanced regular expressions to detect code blocks in real-time. It implements a naming algorithm that scans for keywords (like `package`, `func`, `class`, or `import`) to provide high-probability filename suggestions.

### The Translator: `internal/platform/`
The "Communication Hub."
- **Unified API Adapter**: Every AI provider (OpenAI, Anthropic, Google, DeepSeek) uses a different JSON schema. The Platform Manager abstracts these differences. It takes a generic Viren `ChatRequest` and translates it into the provider-specific payload.
- **Streaming Implementation (SSE)**: Implements a custom Server-Sent Events client. It parses the incoming byte stream from the API in real-time, extracting text "deltas" and passing them immediately to the UI layer for low-latency rendering.
- **Local Bridge**: Manages connections to local LLM servers like **Ollama**. It treats `localhost` as a first-class citizen, ensuring that the local experience is just as smooth as the cloud experience.

---

## 3. Data Persistence Logic

Viren utilizes a local-first persistence strategy to ensure that your technical history is never lost or leaked.

### A. The History Database
Conversations are stored in `~/.viren/tmp/` as high-density JSON files. 
- **Encryption in Transit**: While files are plain JSON, they are created with `0600` permissions.
- **Atomicity**: Viren uses a "Write-Ahead" pattern where the temporary session is updated after every successful turn, preventing data loss during a power failure.

### B. Configuration Merging
On startup, the `config` module performs a three-way merge:
1.  **Hard-coded Defaults**: The baseline "Viren Experience."
2.  **Global Config**: Your settings in `config.json`.
3.  **Ephemeral Flags**: Overrides provided via CLI (e.g., `-p groq`).

---

## 4. Concurrency Patterns & Memory Management

### Goroutine Ownership
- **Stream Listener**: A dedicated routine that watches the HTTP response body.
- **Animation Timer**: A background ticker that updates the loading spinner at 10Hz.
- **Signal Watcher**: A routine that blocks on the OS signal channel.

### Memory Pooling
To handle the **Codedump (`!d`)** of large projects, Viren uses a `sync.Pool` for `bytes.Buffer` objects. This significantly reduces the overhead of the Go garbage collector when processing thousands of files simultaneously.

---

## 5. Tokenization Engine Details

Viren includes a native Go implementation of the **Tiktoken** algorithm (specifically the `cl100k_base` encoding).
- **Offline Count**: You can run `viren -t file.go` to get a precise token count without any internet connection.
- **Cost Estimation**: The Chat Manager uses this engine to calculate when the conversation history needs to be truncated to avoid "Context Overflow" errors from the AI providers.

---

## 6. Binary Stripping & Optimization

The production Viren binary is optimized for minimal size:
- **LDFLAGS**: We use `-s -w` to strip debug symbols and the DWARF table.
- **Static Linking**: The binary is statically linked so it can run on systems without the Go runtime installed.

---

## 7. Directory Structure & Module Ownership

- `cmd/viren/`: Entry point, CLI flags, interactive loop logic.
- `internal/chat/`: Conversation state, bang-command parsing, file format extraction.
- `internal/config/`: JSON i/o, default config generation, onboarding wizard UI.
- `internal/platform/`: Provider-specific JSON mappers, SSE streaming clients.
- `internal/ui/`: Custom ANSI rendering engine, fzf bridge, theme registry.
- `internal/util/`: File system helpers, path cleaners, hashing functions.
- `pkg/types/`: Global schema definitions used by all internal packages.

---

## 8. Detailed Use Cases for the Architect

### Case 1: Legacy Code Migration
Using `!d` on a monolithic repo allows the architect to ask: *"Identify all tightly coupled dependencies between module A and module B and suggest a decoupling strategy."*

### Case 2: Security Audit
By loading `schema.sql` and `api_handlers.go`, the architect can ask: *"Scan these files for potential SQL injection points or lack of input validation."*

---

## 9. Error Handling & Recovery

Viren uses a multi-tiered error recovery system:
1. **Network Retry**: Transient HTTP errors (503, 504) trigger an immediate exponential backoff retry.
2. **Context Compression**: If a prompt is too large, the Chat Manager attempts to strip the oldest messages to fit the limit before failing.
3. **Panic Catching**: Critical modules are wrapped in `recover()` blocks to ensure a single failed operation doesn't crash the entire interactive session.

---

## 10. The IPC Roadmap

Our 2026-2027 roadmap includes:
1.  **Shared Memory IPC**: A local socket bridge allowing external tools to read the active Viren context.
2.  **Plugin Runtime**: Support for loading external `.wasm` files to extend Viren's command set.
3.  **Local RAG Engine**: Integration of a vector database for hyper-scale codebase indexing.

---

## 11. Developer Environment Prerequisites

To build Viren from source, you need:
- **Go 1.21+**: The compiler.
- **GCC/Clang**: If enabling OCR support via CGO.
- **Make**: For running the automated build targets.
- **Git**: For version tracking and repo management.

---

## 12. Conclusion

**Viren is the intersection of high-performance systems engineering and cutting-edge artificial intelligence. It is built to be a permanent part of your technical infrastructure.**
