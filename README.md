# VIREN: The Advanced Neural Command Interface

<p align="center">
  <img src="docs/logo.png" alt="Viren Logo" width="200" />
</p>

<p align="center">
  <strong>Unleashing elite machine intelligence directly into your terminal environment.</strong>
</p>

---

## 1. What is Viren?

**Viren** is a specialized, high-performance Command Line Interface (CLI) built from the ground up in GoLang. It is designed for software engineers, system architects, and researchers who require an immediate, context-aware bridge to the world's most powerful Large Language Models (LLMs) without the friction of a web browser.

Viren is defined by four core pillars:
1.  **Velocity**: A "zero-wait" architecture with sub-100ms startup times.
2.  **Context**: Native ingestion of local source code, documentation, and active terminal sessions.
3.  **Privacy**: A local-first philosophy where your data never touches a "Viren Cloud."
4.  **Ergonomics**: A keyboard-centric design utilizing fuzzy finding and pipe integration.

As the spiritual successor to the **Cha** project, Viren takes the interactive AI experience to its logical conclusion: a tool that is as stable and essential as `git`, `docker`, or `vim`.

---

## 2. Table of Contents

1.  [Overview](#3-overview)
2.  [The Vision](#4-the-vision)
3.  [Core Technical Features](#5-core-technical-features)
4.  [The Neural Difference](#6-the-neural-difference)
5.  [Installation Deep-Dive](#7-installation-deep-dive)
6.  [Quick Start Guide](#8-quick-start-guide)
7.  [Interactive Command Reference](#9-interactive-command-reference)
8.  [Advanced Context Ingestion](#10-advanced-context-ingestion)
9.  [Configuration & Personalization](#11-configuration--personalization)
10. [API Platform Integration](#12-api-platform-integration)
11. [Local & Open Source Setup (Ollama)](#13-local--open-source-setup)
12. [Performance Benchmarks](#14-performance-benchmarks)
13. [Comparison: Viren vs. Others](#15-comparison-viren-vs-others)
14. [Development & Building](#16-development--building)
15. [Project Roadmap](#17-project-roadmap)
16. [Privacy & Security Guarantee](#18-privacy--security-guarantee)
17. [Contributing](#19-contributing)
18. [Inspired By](#20-inspired-by)
19. [License](#21-license)

---

## 3. Overview

Viren is engineered for the "Flow State." Developers spend their lives in the terminal; switching to a Chrome tab to ask an AI about a specific bug breaks the neural loop of coding. Viren brings the AI to the code.

By leveraging Go's native concurrency and efficient memory management, Viren can handle massive codebases (via the Codedump feature) and inject them into the context window of elite models like Claude 3.5 Sonnet or GPT-4o. It acts as a "Universal API Translator," allowing you to switch between 11+ providers and 40+ specialized logic modes with a single command.

---

## 4. The Vision

### Normalization of Intelligence
We believe AI should not be a "destination" but a utility. Just as you don't think twice about running `ls` to see files, you shouldn't think twice about running `viren` to refactor a function. Our vision is to make high-level reasoning an invisible layer of the Unix environment.

### The Keyboard-First Era
Web UIs are optimized for engagement and mouse clicks. Viren is optimized for **throughput**. By using fuzzy-selection (`fzf`) for every menu, we ensure that you can change models, themes, or platforms in under 2 seconds without ever moving your hands from the home row.

---

## 5. Core Technical Features

*   **Asynchronous Engine**: Viren uses background goroutines to handle animations, API requests, and terminal rendering simultaneously.
*   **Intelligent Tokenization**: Built-in estimation tools allow you to check the "weight" of a file before sending it to the AI, helping you manage costs and context limits.
*   **ANSI UI Framework**: A custom-built rendering engine that provides high-fidelity borders, boxes, and syntax highlighting without requiring a GPU-accelerated terminal.
*   **Session Persistence**: Every conversation is serialized to a local JSON database, allowing for instant "checkpointing" and resumption of complex tasks.
*   **Dynamic Logic Injection**: Viren modifies the "System Prompt" in real-time based on your "Domain Mode" selection, ensuring the AI behaves exactly like the expert you need (e.g., a Database Admin or a Frontend Guru).

---

## 6. The Neural Difference

What separates Viren from a simple "Chat-with-API" script?

### Contextual Awareness
Viren "sees" your environment. When you use the **Shell Record (`!x`)** feature, Viren opens a sub-process, watches your terminal output, and captures errors that normally would require manual copy-pasting.

### Zero Cloud Dependency
Viren does not have a backend. There is no `api.viren.sh`. Your requests go directly from your local binary to the AI provider. This architectural choice ensures that your proprietary code never passes through a third-party's logging server.

---

## 7. Installation Deep-Dive

### The Automated Installer
The recommended way to install Viren is via our production-ready shell script. It is designed to be idempotent and safe.

```bash
curl -sSL https://raw.githubusercontent.com/fraol163/viren/main/install.sh | bash
```

### What happens during installation?
1.  **OS/Arch Detection**: The script identifies if you are on x86_64, ARM64 (M1/M2 Macs), or Android (Termux).
2.  **Dependency Audit**: It checks for `go` (>=1.21) and `fzf`.
3.  **Automatic Recovery**: If `fzf` is missing, the script attempts to install it via `apt`, `brew`, or `dnf` automatically.
4.  **Temp Workspace**: To keep your filesystem clean, the script clones Viren into `/tmp/viren-build-randomID/`.
5.  **CGO Tuning**: It automatically disables CGO for Termux builds to prevent common linker errors while enabling it for macOS/Linux to support optional OCR features.
6.  **Binary Stripping**: The resulting binary is optimized for size and performance before being moved to `/usr/local/bin`.

---

## 8. Quick Start Guide

1.  **Install**: Run the curl command above.
2.  **Configure**: Set your API key: `export OPENAI_API_KEY="your-key"`.
3.  **Launch**: Type `viren`.
4.  **Onboard**: Enter your name and role (Backend, Frontend, etc.).
5.  **Chat**: Type "How do I optimize a Go map lookup?"
6.  **Export**: Type `!e` to save the code the AI gives you.

---

## 9. Interactive Command Reference

Inside the `viren` shell, use these "Bang Commands" for total control:

| Command | Category | Functional Capability |
| :--- | :--- | :--- |
| `!q` | Session | **Quit**: Safely exits and saves your history. |
| `!h` | UI | **Help**: Opens the interactive dashboard. |
| `!c` | Context | **Clear**: Wipes the screen and the context window. |
| `!m` | Logic | **Model**: Switch between GPT-4, Claude, DeepSeek, etc. |
| `!p` | Logic | **Platform**: Switch between OpenAI, Groq, Ollama, etc. |
| `!u` | Tone | **Personality**: Switch between 7 AI personas. |
| `!v` | Logic | **Mode**: Apply one of 40+ specialized domain prompts. |
| `!z` | UI | **Theme**: Instantly change colors (Neon, Paper, Matrix). |
| `!x` | Ingestion | **Shell Record**: Capture your terminal output for the AI. |
| `!d` | Ingestion | **Codedump**: Bundle your whole project for review. |
| `!l` | Ingestion | **Load**: Inject specific files (Go, Py, PDF, Images). |
| `!s` | Ingestion | **Scrape**: Feed a website URL to the AI. |
| `!w` | Ingestion | **Web Search**: Live search via Brave Search API. |
| `!a` | Session | **History**: Browse and restore past conversations. |
| `!y` | Utility | **Clipboard**: Copy responses to your system clipboard. |
| `!b` | Logic | **Backtrack**: Revert the last 1 or more messages. |

---

## 10. Advanced Context Ingestion

### The Codedump Protocol (`!d`)
Viren is designed for "Repository-Level Intelligence." When you run `!d`, Viren:
1.  Loads your `.gitignore` to avoid binary files and `node_modules`.
2.  Presents an `fzf` menu allowing you to exclude specific files manually.
3.  Minifies and bundles the remaining source code into a structured context package.
4.  Allows you to ask questions like: *"Where is the race condition in my websocket handler?"* or *"Summarize the architecture of this entire project."*

### Multi-File Reasoning
You can use `!l` to load multiple files at once. Viren will tag each file clearly in the prompt, allowing the AI to understand the relationship between your `main.go` and your `types.go`, for example.

---

## 11. Configuration & Personalization

### The Neural Profile
Viren uses a `user_profile` section in `config.json` to tailor its responses.
*   **Role-Based Tuning**: If you are a "System Architect," the AI will prioritize scalability and patterns. If you are a "Junior Developer," it will focus on explaining syntax and logic.
*   **Env-Awareness**: List your editor (Vim, VS Code) so the AI provides relevant keyboard shortcuts and plugin suggestions.

### Config File Location
Settings are stored in `~/.viren/config.json`.
- **`shallow_load_dirs`**: List paths (like `/`, `/home`) that Viren should never scan recursively. This protects you from accidental context overflows.
- **`exit_key`**: Customize how you quit the app.

---

## 12. API Platform Integration

Viren supports every major AI platform on the market. Below are the official links to obtain keys:

| Provider | Purpose | Link |
| :--- | :--- | :--- |
| **OpenAI** | General Purpose | [openai.com/api/](https://openai.com/api/) |
| **Brave Search** | Live Internet | [brave.com/search/api/](https://brave.com/search/api/) |
| **OpenRouter** | Aggregator | [openrouter.ai](https://openrouter.ai/settings/keys) |
| **Google Gemini**| Multimodal | [ai.google.dev](https://ai.google.dev/gemini-api/docs/api-key) |
| **Anthropic** | Reasoning | [console.anthropic.com](https://console.anthropic.com/) |
| **DeepSeek** | Coding Efficiency| [deepseek.com](https://api-docs.deepseek.com/) |
| **Groq** | Inference Speed | [console.groq.com](https://console.groq.com/keys) |
| **xAI** | Real-time Info | [console.x.ai](https://x.ai/api) |
| **Mistral** | Efficiency | [console.mistral.ai](https://docs.mistral.ai/getting-started/quickstart) |
| **Together AI** | Open Source | [together.ai](https://docs.together.ai/docs/quickstart) |
| **AWS Bedrock** | Enterprise | [aws.amazon.com/bedrock](https://aws.amazon.com/bedrock) |

---

## 13. Local & Open Source Setup

For users who cannot send code to the cloud due to security constraints, Viren fully supports **Ollama**.

1.  **Download Ollama**: Visit [ollama.com](https://ollama.com).
2.  **Pull a Model**: `ollama run deepseek-coder:6.7b`.
3.  **Connect Viren**: Launch Viren and type `!p ollama`.
4.  **Offline Privacy**: All reasoning now happens 100% on your machine. No data leaves your local network.

---

## 14. Performance Benchmarks

Viren is built for speed. Here is how it compares to standard tools:

| Action | Viren (Go) | Python-Based CLIs | Web-Based AI |
| :--- | :--- | :--- | :--- |
| **Startup Time** | < 100ms | 1.2s - 2.5s | 3s - 8s (Browser) |
| **Codedump Scan**| 400ms (1k files)| 3.5s | Manual Upload |
| **Memory Usage** | ~15MB | 120MB+ | 800MB+ (Chrome) |
| **Syntax Highlighting**| Real-time | Delayed | Real-time |

---

## 15. Comparison: Viren vs. Others

### Viren vs. ChatGPT/Claude Web
- **Web**: Manual copy-pasting, distracting UI, no file system access.
- **Viren**: Direct file ingestion, shell recording, keyboard-driven.

### Viren vs. Other CLI Tools
- **Others**: Often complex setup, Python dependencies, slow startup.
- **Viren**: Single static binary, ultra-fast, built-in fuzzy finder (`fzf`).

---

## 16. Development & Building

### Prerequisites
- **Go**: 1.21 or higher.
- **FZF**: Required for the interactive UI.
- **Make**: For standard build tasks.

### Build Instructions
```bash
git clone https://github.com/fraol163/viren.git
cd viren
make build
```
The binary will be generated in `./bin/viren`.

### Running Tests
```bash
make test
```

---

## 17. Project Roadmap

*   **Phase 1**: (Completed) Core speed optimization, 11-provider support, Codedump.
*   **Phase 2**: (In Progress) Direct `.gguf` support (Llama.cpp integration), Multimodal image ingestion.
*   **Phase 3**: (Planning) Autonomous Agent Mode, local RAG (Vector Database), VS Code/Vim plugins.

---

## 18. Privacy & Security Guarantee

Viren is a **Data-Transparent** application.
*   **Encryption**: We recommend environment variables so keys are never written to disk in plain text.
*   **No Logs**: Viren does not log your prompts to our servers because we do not have servers.
*   **Source Integrity**: The project is open-source. You can verify every line of networking code in `internal/platform/`.

---

## 19. Contributing

We are actively seeking contributions for:
- New **Domain Modes** (System Prompts for specific niches).
- Improved **Documentation** and examples.
- Better **Windows Support** optimization.

Please see [CONTRIBUTING.md](./docs/contributing.md) for details.

---

## 20. Inspired By

Viren is the spiritual successor to the **Cha** project, originally created by [Mehmetmhy](https://github.com/mehmetmhy). Viren was inspired by the elegance of Cha's interactive design but was rebuilt from scratch to meet the performance and stability needs of modern software engineering. We thank the original Cha creator for pioneering this workflow.

---

## 21. License

Viren is released under the **MIT License**. You are free to use, modify, and distribute it. See [LICENSE](./LICENSE) for details.

---

<p align="center">
  <strong>Built with precision by developers, for developers.</strong>
</p>
