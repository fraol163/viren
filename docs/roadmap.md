# Viren Project Roadmap: The Future of Neural Workflows

Viren is currently in **Version 1.0.0 (The Rebrand)**. This roadmap serves as our strategic north star, detailing the transition from a high-performance CLI to a fully-realized "Neural Operating System" layer. We are building for a future where AI is not a destination, but an invisible utility.

---

## Phase 1: The Foundation (COMPLETED)
The goal of Phase 1 was to build a tool stable enough for daily production use by elite developers.
- [x] **Static Go Engine**: Full transition from legacy scripts to a compiled, type-safe binary.
- [x] **Multi-Provider Bridge**: Unified access to 11+ providers (OpenAI, Anthropic, DeepSeek, etc.).
- [x] **Context Mastery**: High-speed implementation of Codedump (`!d`) and File Loading (`!l`).
- [x] **Sub-100ms Startup**: Aggressive optimization of the configuration and initialization sequence.
- [x] **Interactive Persistence**: A robust local JSON history database.
- [x] **Universal ANSI UI**: Themes and borders that work in any modern terminal.

---

## Phase 2: Local Intelligence & Multimodality (Q1 - Q2 2026)
We are currently moving away from cloud dependency and toward "Edge Intelligence."
- **Direct Llama.cpp Integration**: Move beyond Ollama by allowing Viren to run `.gguf` files directly using C-bindings. This removes the need for a separate background service.
- **Multimodal Context Expansion**:
    - **Image Reasoning**: Support for loading images (`!l bug_screenshot.png`) and asking the AI to diagnose visual glitches, UI alignment issues, or OCR errors.
    - **Enhanced OCR**: Deepening Tesseract support for multi-language text extraction and better layout recognition.
- **Advanced PDF Parsing**: Native support for reading tables, charts, and mathematical notations in technical whitepapers.

---

## Phase 3: The "Agentic" Evolution (Q3 - Q4 2026)
This phase transforms Viren from a "Chatbot" into an "Action Engine."
- **Autonomous Shell Agent**: A new mode (`!agent`) where Viren can generate a multi-step plan, execute shell commands (with human "Y/N" confirmation), and check the results to complete a task (e.g., "Deploy this Go app to AWS using Terraform").
- **Local RAG (Retrieval Augmented Generation)**: 
    - Integration of a local vector database (like Faiss or Chroma).
    - Viren will index your entire `~/Projects` folder, allowing you to ask: "Where did I implement that JWT logic three months ago across all my repos?"
- **Logic Plugin System**: A Lua-based scripting engine where users can share custom "Domain Modes" and "Command Bangs."

---

## Phase 4: Ecosystem Integration (2027)
Making Viren an invisible, seamless part of your existing dev stack.
- **The "Editor Bridge"**:
    - **Viren-Vim/Neovim**: A plugin that allows you to "pipe" a visual selection directly into an active Viren buffer.
    - **Viren-VSCode**: A specialized terminal pane with integrated history syncing and "Click-to-Apply" code fixes.
- **Encrypted Profile Sync**: An optional, zero-knowledge, end-to-end encrypted method to sync your Neural Profile and Chat History between multiple devices (Laptop, Desktop, Cloud Server).
- **System-Wide Knowledge**: Allowing Viren to (safely) ingest system logs (`dmesg`, `syslog`) and active process stats to act as a real-time "Systems Architect."

---

## Phase 5: Normalization & Voice (Visionary)
- **Voice-to-CLI**: High-speed, local Whisper integration for hands-free terminal queries.
- **Natural Language OS**: A mode where the entire shell experience is abstracted through Viren, allowing users to perform complex Unix tasks without knowing specific flags.

---

## Our Engineering Principles (The "Viren Code")
1.  **Speed Above All**: If a feature adds 10ms to the startup time, it must be optional or backgrounded.
2.  **Privacy is Non-Negotiable**: We will never build a feature that requires a Viren-owned cloud server.
3.  **Keyboard-Centric**: Every feature must be fully usable without a mouse.
4.  **Open Standards**: We will always prioritize open-source models and standards (like Llama.cpp and GGUF).

---

## How to Contribute to the Roadmap
Viren is a community project. Your feedback determines our priority.
- **Join the RFCs**: Look for "Request for Comments" issues on GitHub.
- **Benchmark**: Help us test the performance of Phase 2 local models on different hardware (M1, RTX 4090, Raspberry Pi).
- **Build**: We prioritize features that come with a working Proof of Concept.

**Viren is not just a tool; it is a decade-long project to redefine how humans interact with machines.**
