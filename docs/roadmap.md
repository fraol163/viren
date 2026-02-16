# Project Roadmap & Future Vision

Viren is currently in **Version 1.0.0 (The Rebrand)**. This roadmap outlines the strategic direction of the project, focusing on becoming the most essential developer tool in the AI era.

---

## Phase 1: Foundation (Completed)
- [x] **Static Go Engine**: Full transition from legacy scripts to a high-concurrency Go binary.
- [x] **Multi-Provider Bridge**: Support for OpenAI, Anthropic, Gemini, Groq, DeepSeek, and more.
- [x] **Context Mastery**: Implementation of Codedump (`!d`) and File Loading (`!l`).
- [x] **Sub-100ms Startup**: Optimization of the configuration and initialization logic.
- [x] **Fzf Integration**: Seamless, lightning-fast interactive menus.

---

## Phase 2: Local Intelligence (Current Focus)
- **Direct Llama.cpp Integration**: Move beyond Ollama to allow Viren to run `.gguf` files directly with zero external dependencies.
- **Multimodal Context**: 
    - **Image Reasoning**: Send a screenshot of a bug or a UI mockup directly to GPT-4o or Gemini for analysis.
    - **Native OCR**: Deepening the Tesseract integration for better text extraction from images.
- **PDF Intelligence**: Improved parsing for complex multi-column scientific papers and technical documentation.

---

## Phase 3: The "Agentic" Shift (Short-Term)
- **Autonomous Shell Agent**: A mode where Viren can generate and execute a sequence of commands (with human confirmation) to set up environments or fix complex system errors.
- **Local RAG (Alpha)**: Implementation of a local vector database (like Faiss) to allow Viren to index massive repositories (10,000+ files) without overflowing the model's context window.
- **Plugin System**: A Lua-based scripting engine to allow users to write their own custom commands and domain modes.

---

## Phase 4: Ecosystem Integration (Mid-Term)
- **Editor Sync (The Bridge)**: 
    - **Viren-Vim**: A plugin to sync the current buffer to Viren.
    - **Viren-VSCode**: A dedicated terminal extension.
- **Encrypted Sync**: An optional, E2EE (End-to-End Encrypted) cloud backup for history and profiles, allowing you to use Viren on your laptop and your server with the same context.

---

## Phase 5: Normalization (Long-Term)
- **Voice-to-CLI**: High-speed, local whisper integration for hands-free terminal queries.
- **System-Wide Knowledge**: Allowing Viren to (safely) ingest system logs and environment stats to act as a real-time "Systems Architect."

---

## How We Prioritize
1.  **Speed**: If a feature slows down `viren -v`, it must be moved to a background thread.
2.  **Privacy**: We will never implement a feature that requires a Viren-owned cloud server.
3.  **Developer Utility**: Features are judged by how many keystrokes they save.

**Viren is not just a tool; it is a decade-long project to redefine the command line.**
