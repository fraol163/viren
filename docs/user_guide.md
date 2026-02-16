# The Master Viren User Guide: From Apprentice to Neural Architect

Welcome to the definitive guide for **Viren**. This document is designed to take you from a basic user to a master of terminal-based AI automation. Viren is not just a chat tool; it is a context-injection engine that allows you to treat Large Language Models (LLMs) as first-class citizens of your Command Line Interface (CLI).

---

## 1. The Viren Philosophy: Documentation Logic

**Why is this guide so detailed?**
Viren is a tool for developers. Developers value precision, depth, and edge-case coverage. While the interface is designed to be intuitive, its advanced capabilities—like native OCR parsing and sub-shell recording—require a deeper explanation to extract their full value. This guide provides the high-density information needed to build a "Neural Workflow" that saves hours of manual labor every week.

---

## 2. Initial Setup & Onboarding Mechanics

### The First Launch Sequence
When you execute `viren` for the first time, the binary detects the absence of a `~/.viren/` directory and initiates the **Neural Profile Initialization**.

1.  **Identity Mapping**: The AI will ask for your name. This is used to customize the greeting and provide a sense of persistent partnership.
2.  **Role Classification**: This is the most critical step. If you identify as a "Senior Rust Engineer," the AI will skip basic syntax explanations and focus on ownership, lifetimes, and zero-cost abstractions. If you identify as a "Product Manager," it will focus on strategy and communication.
3.  **Environmental Context**: Listing your primary OS (e.g., "Fedora") and editor ("Neovim") allows the AI to provide copy-pasteable commands that work specifically for your system.

---

## 3. The Interactive Command Layer (Bang Commands)

Inside the Viren shell, you can communicate in natural language. However, to control the application's internal state machine, you must use **Bang Commands** (`!`).

### A. Session & History Management
- **`!q` (Safe Quit)**: Exits the loop and ensures the current `Messages` slice is saved to the JSON history.
- **`!a` (History Manager)**: Opens a searchable `fzf` menu of every session you have ever had.
- **`!c` (Context Purge)**: Clears the current session's memory and wipes the terminal screen.

### B. Logical & Platform Multiplexing
- **`!m` (Model Switcher)**: Fetches the live model list from your current provider.
- **`!p` (Platform Switcher)**: Change your entire infrastructure backend. Switch from cloud (Anthropic) to local (Ollama) in under 3 seconds.
- **`!o` (Global Selector)**: A master list of every model from every platform you have an API key for.

### C. Behavioral Tuning
- **`!u` (Personality Matrix)**: Viren comes with 7 pre-tuned personas.
    - *Analytical*: Focuses on data structures and logic proofs.
    - *Playful*: Injects humor and wit into technical explanations.
    - *Rick Sanchez*: A cynical, brilliant, and sarcastic persona.
- **`!v` (Domain Mode)**: This changes the core "System Mandate." 
    - *Code Whisperer*: Transforms Viren into a pair-programmer.
    - *Socratic*: Transforms Viren into a teacher.

---

## 4. Masterclass: Advanced Context Ingestion

### The "Load" Protocol (`!l`)
The `!l` command is for precision context injection.
- **Multiple Selection**: Use `Tab` in the fzf menu to pick multiple files.
- **Rich Format Parsing**:
    - **PDFs**: Viren uses a native parser to extract text.
    - **Excel/CSV**: It converts row/column data into structured Markdown.
    - **Images (OCR)**: Perfect for screenshots of failing logs.

### The Codedump Protocol (`!d`)
For project-wide architectural analysis.
1.  **Automatic Ignoring**: Viren reads your `.gitignore`.
2.  **Manual Pruning**: An `fzf` menu appears, letting you exclude specific files.
3.  **Bundling**: Viren creates a single, high-density XML-style context package.

---

## 5. Pro-Developer Workflows: The Action Engine

### The Diagnostic Loop (`!x`)
This is the "Killer Feature" for systems engineers.
1.  Launch `!x`.
2.  Perform your failing workflow.
3.  Type `exit`.
4.  Viren captures the output and asks the AI for a fix.

### The Export Pipeline (`!e`)
Viren is a **Producer**, not just a **Chatter**.
1.  Once the AI generates a code block, type `!e`.
2.  Viren's "Smart Guess" engine looks at the code and suggests a name.
3.  Hit `Enter` to save the file to your disk.

---

## 6. High-Density Prompt Engineering within Viren

To get the most out of Viren, you must understand how it builds prompts. Viren uses a "Layered Context" strategy:

1.  **Layer 1: The Domain Mandate**: (Set via `!v`) Defines the AI's core logic.
2.  **Layer 2: The Neural Profile**: (Set in config) Defines your expertise level.
3.  **Layer 3: The Active Context**: (Files loaded via `!l` or `!d`) Provides the data.
4.  **Layer 4: The Immediate Query**: Your specific question.

---

## 7. Performance Optimization Workflows

### Managing the Context Window
Every AI model has a limit.
- **Symptom**: The AI starts "hallucinating" or forgetting.
- **Solution**: Use `!c` to purge and re-load.

---

## 8. Detailed Mode Descriptions

### `!v zenith`
**Purpose**: Advanced Logical Reasoning.
**Logic**: Injects instructions forcing the model to use Chain-of-Thought (CoT) and avoid jumping to conclusions.

### `!v codewhisperer`
**Purpose**: Idiomatic Refactoring.
**Logic**: Strips all conversational filler. The model will only provide code blocks and high-level comments.

### `!v timeanalyzer`
**Purpose**: Performance Bottleneck Detection.
**Logic**: Forces the model to calculate Big-O complexity for every function it sees.

---

## 9. Keyboard Shortcuts for Maximum Throughput

- **`Ctrl+C`**: Cancels an active AI stream.
- **`Ctrl+D`**: Safely quits.
- **`Up Arrow`**: Cycles through your *previous prompts*.
- **`\`**: Triggers multi-line input mode.

---

## 10. Frequently Asked Questions (FAQ)

**Q: Can I use Viren with local models?**
A: Yes. Use `!p ollama`.

**Q: Does Viren support Windows?**
A: Yes. We recommend using **Windows Terminal**.

**Q: Where is my history saved?**
A: Everything is stored in `~/.viren/tmp/`.

---

## 11. Security Best Practices

1.  **Environment Variables**: Never put your API keys in a text file.
2.  **Codedump Caution**: Always review the file list in `!d`.
3.  **VPN**: When using `!w` (Web Search), use a VPN.

---

## 12. Troubleshooting Common Scenarios

### Scenario: Viren is not seeing my API keys
- **Fix**: Ensure you have `export`ed them in your shell profile and ran `source ~/.zshrc`.

### Scenario: The ASCII art looks broken
- **Fix**: Check your terminal font. Use a "Nerd Font" for the best compatibility with symbols and boxes.

---

## 13. Mastering Single-Shot Mode

Viren isn't just for chatting. It can be used as a high-speed utility in your build scripts.
`viren -l config.yaml "Convert this to JSON format" > config.json`

---

## 14. Conclusion: The Neural Advantage

Viren is more than a CLI; it is a force-multiplier for your technical ability. By mastering the commands in this guide, you are becoming an engineer who is contextually aware, high-speed, and ready for the era of AI-integrated development.

**Master the interface. Master the machine. Viren is your new technical partner.**
