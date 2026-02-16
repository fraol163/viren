# VIREN: The Advanced Neural Command Interface

<p align="center">
<pre>
██╗      ██╗   ██╗██╗██████╗ ███████╗███╗   ██╗
╚██╗     ██║   ██║██║██╔══██╗██╔════╝████╗  ██║
 ╚██╗    ██║   ██║██║██████╔╝█████╗  ██╔██╗ ██║
 ██╔╝    ╚██╗ ██╔╝██║██╔══██╗██╔══╝  ██║╚██╗██║
██╔╝      ╚████╔╝ ██║██║  ██║███████╗██║ ╚████║
╚═╝        ╚═══╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═══╝
</pre>
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
    *   [Linux & macOS](#linux--macos)
    *   [Windows Installation](#windows-installation)
    *   [Android (Termux)](#android-termux)
6.  [Quick Start Guide](#8-quick-start-guide)
7.  [Interactive Command Reference](#9-interactive-command-reference)
8.  [Advanced Context Ingestion](#10-advanced-context-ingestion)
9.  [Logic & Behavior](#11-logic--behavior)
    *   [AI Personalities](#ai-personalities)
    *   [Domain Modes](#domain-modes)
10. [Configuration & Personalization](#12-configuration--personalization)
11. [API Platform Integration](#13-api-platform-integration)
12. [Local & Open Source Setup (Ollama)](#14-local--open-source-setup)
13. [Performance Benchmarks](#15-performance-benchmarks)
14. [Comparison: Viren vs. Others](#16-comparison-viren-vs-others)
15. [Development & Building](#17-development--building)
16. [Project Roadmap](#18-project-roadmap)
17. [Privacy & Security Guarantee](#19-privacy--security-guarantee)
18. [Detailed Use Cases](#20-detailed-use-cases)
19. [Troubleshooting Strategy](#21-troubleshooting-strategy)
20. [Contributing](#22-contributing)
21. [Inspired By](#23-inspired-by)
22. [License](#24-license)

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
*   **Intelligent Tokenization**: Built-in estimation tools allow you to check the "weight" of a file before sending it to the AI.
*   **ANSI UI Framework**: A custom-built rendering engine that provides high-fidelity borders, boxes, and syntax highlighting.
*   **Session Persistence**: Every conversation is serialized to a local JSON database.
*   **Dynamic Logic Injection**: Viren modifies the "System Prompt" in real-time based on your selection.

---

## 6. The Neural Difference

What separates Viren from a simple "Chat-with-API" script?

### Contextual Awareness
Viren "sees" your environment. When you use the **Shell Record (`!x`)** feature, Viren opens a sub-process, watches your terminal output, and captures errors that normally would require manual copy-pasting.

### Zero Cloud Dependency
Viren does not have a backend. Your requests go directly from your local binary to the AI provider.

---

## 7. Installation Deep-Dive

### Linux & macOS
The recommended way to install Viren is via our production-ready shell script.

```bash
curl -sSL https://raw.githubusercontent.com/fraol163/viren/main/install.sh | bash
```

### Windows Installation
Viren is fully supported on Windows. For the best experience, use **Windows Terminal** with **PowerShell**.

#### Method 1: Pre-built Binary
1. Download the latest `viren-windows-amd64.exe` from the [Releases](https://github.com/fraol163/viren/releases) page.
2. Rename it to `viren.exe`.
3. Add the directory containing `viren.exe` to your system's `PATH`.

#### Method 2: Build from Source (PowerShell)
1. Install [Go](https://go.dev/dl/).
2. Install [Git](https://git-scm.com/downloads).
3. Install [fzf](https://github.com/junegunn/fzf#windows) (via `choco install fzf` or `scoop install fzf`).
4. Run:
```powershell
git clone https://github.com/fraol163/viren.git
cd viren
go build -o viren.exe cmd/viren/main.go
```

#### Method 3: WSL2
You can use the Linux installation method inside WSL2 (Ubuntu/Debian).

### Android (Termux)
1. Install Termux.
2. Run the Linux installation script. Viren detects Termux and disables CGO for maximum compatibility.

---

## 8. Quick Start Guide

1.  **Install**: Run the installation steps above.
2.  **Configure**: Set your API key: `export OPENAI_API_KEY="your-key"` (PowerShell: `$env:OPENAI_API_KEY="your-key"`).
3.  **Launch**: Type `viren`.
4.  **Onboard**: Enter your name and role.
5.  **Chat**: Type your first query.

---

## 9. Interactive Command Reference

Inside the `viren` shell, use these "Bang Commands":

| Command | Category | Functional Capability |
| :--- | :--- | :--- |
| `!q` | Session | **Quit**: Safely exits and saves your history. |
| `!h` | UI | **Help**: Opens the interactive dashboard. |
| `!c` | Context | **Clear**: Wipes the screen and context. |
| `!m` | Logic | **Model**: Switch between LLMs. |
| `!p` | Logic | **Platform**: Switch between providers. |
| `!u` | Tone | **Personality**: Switch between 7 AI personas. |
| `!v` | Logic | **Mode**: Apply specialized domain prompts. |
| `!z` | UI | **Theme**: Instantly change colors. |
| `!x` | Ingestion | **Shell Record**: Capture terminal output. |
| `!d` | Ingestion | **Codedump**: Bundle project for review. |
| `!l` | Ingestion | **Load**: Inject specific files (Go, Py, PDF, etc.). |
| `!s` | Ingestion | **Scrape**: Feed a website URL to the AI. |
| `!w` | Ingestion | **Web Search**: Live search via Brave. |
| `!a` | Session | **History**: Browse/restore conversations. |

---

## 10. Advanced Context Ingestion

### Codedump Protocol (`!d`)
Automatically bundles your source code while respecting `.gitignore`. Perfect for architectural reviews or finding project-wide bugs.

---

## 11. Logic & Behavior

Viren allows you to precisely tune how the AI thinks and speaks through Personalities and Domain Modes.

### AI Personalities
Switch using `!u`. Each personality alters the tone, verbosity, and style of the AI.
- **Analytical**: Logical, systematic, and data-driven.
- **Creative**: Artistic, imaginative, and metaphor-heavy.
- **Focused**: Goal-oriented, concise, and productivity-focused.
- **Empathetic**: Emotionally intelligent and supportive.
- **Playful**: Fun, energetic, and full of wit.
- **Balanced**: Versatile and adaptive (Default).
- **Rick Sanchez**: Cynical, sarcastic, and unhinged (Rick and Morty style).

### Domain Modes
Switch using `!v`. Each mode applies a heavy system mandate for specific technical tasks.
- **Standard**: Professional balanced tone.
- **Zenith**: Advanced reasoning with cosmic metaphors.
- **Code Whisperer**: Expert refactoring and idiomatic patterns.
- **Socratic**: Never gives direct answers; leads you through questions.
- **Complexity Analyzer**: Calculates Big-O time and space complexity.
- **DSA Mode**: Data Structures & Algorithms specialist.
- **CyberSec**: Security audit and penetration testing focus.
- **... and 30+ more** (Physics, Chemistry, Finance, Marketing, etc.)

---

## 12. Configuration & Personalization

Settings are stored in `~/.viren/config.json`.
*   **`user_profile`**: Your name, role, and professional ambition.
*   **`shallow_load_dirs`**: Folders Viren will not scan recursively.
*   **`current_theme`**: Set your preferred default look.

---

## 13. API Platform Integration

Viren supports OpenAI, Brave, OpenRouter, Gemini, Anthropic, DeepSeek, Groq, xAI, Mistral, Together AI, and AWS Bedrock. See [docs/api_keys.md](./docs/api_keys.md) for full setup details.

---

## 14. Local & Open Source Setup

Viren is a bridge to local LLMs via **Ollama**.
1. Install Ollama.
2. In Viren, type `!p ollama`.
3. Reasoning happens entirely on your local hardware.

---

## 15. Performance Benchmarks

Viren starts in under 100ms and can scan 1,000 code files in under 400ms.

---

## 16. Comparison: Viren vs. Others

Viren is faster than Python-based tools and more context-aware than standard web interfaces.

---

## 17. Development & Building

### Prerequisites
- **Go**: 1.21+
- **FZF**: Required for menus.

```bash
make build
```

---

## 18. Project Roadmap

Phase 2 will include Llama.cpp integration and multimodal (image) ingestion.

---

## 19. Privacy & Security Guarantee

Viren is data-transparent. We do not have servers. Your data stays on your disk.

---

## 20. Detailed Use Cases

Use Viren for project auditing, bug hunting, fact-checking via web search, and automated code generation.

---

## 21. Troubleshooting Strategy

Ensure your API keys are exported and `fzf` is in your system path.

---

## 22. Contributing

We welcome pull requests for new domain modes and platform adapters. See [CONTRIBUTING.md](./docs/contributing.md).

---

## 23. Inspired By

Viren is the spiritual successor to the **Cha** project, originally created by [Mehmetmhy](https://github.com/mehmetmhy). Viren takes the interactive experience to a professional engineering standard.

---

## 24. License

Viren is released under the **MIT License**.

---

<p align="center">
  <strong>Built with precision by developers, for developers.</strong>
</p>
