# Real-World Scenario: Smart Code Export & Project Building

Viren eliminates the manual "Copy-Paste" loop that plagues most AI workflows. This guide demonstrates how to turn an AI conversation into a production-ready file structure using the context-aware export engine.

---

## 1. The Scenario: Building a New Module
Imagine you need to create a complex Go middleware for your project that handles both JWT validation and rate limiting.

---

## 2. Step-by-Step Workflow

### Step 1: The Request
Start Viren and ask for the implementation:
`USER ❯ Create a Go middleware that validates a JWT from the header and uses a Redis store for rate limiting.`

### Step 2: The AI Response
The AI will generate one or more code blocks. Viren's internal regex engine (`codeBlockRegex`) is watching for these in the stream.

### Step 3: The Export Trigger
Once the message finishes, you have two options:
1.  **Manual**: Type `!e`.
2.  **Auto**: If you see the detected command prompt, hit `y`.

### Step 4: The Intelligent Filename Menu
Viren opens an `fzf` menu. Unlike other tools, Viren scans the code content to provide suggestions:
- `auth_middleware.go` (Detected `func Auth`)
- `middleware.go` (Detected `package middleware`)
- `[Custom Name]` (Type your own)

### Step 5: Finalization
Select a name and hit `Enter`. Viren performs an atomic write to your local filesystem.

---

## 3. Advanced Export Modes

When you type `!e`, you are presented with three professional modes:

### A. Turn Export
Saves the **entire conversation** (your prompts + the AI's answers) as a single, formatted Markdown file.
- **Best for**: Generating documentation or "Post-Mortem" reports for a debugging session.

### B. Block Export
Only extracts the raw code blocks (e.g., `bash`, `go`, `javascript`) and prompts you to save each one individually.
- **Best for**: Quickly generating multiple files (like `main.go` and `Makefile`) from a single AI response.

### C. Manual Export
Pops open your configured text editor (`$EDITOR`) with the conversation context, allowing you to edit and prune the text before Viren saves it to a file.

---

## 4. Configuration Integration
Viren checks your `config.json` for the `preferred_editor` setting. If you set this to `nvim`, Viren will use Neovim's power for the manual export step.

```json
{
  "preferred_editor": "nvim"
}
```

---

## 5. Security & Safety Features
- **Collision Detection**: If you try to save a file that already exists, Viren will warn you and ask for a confirmation or a new filename.
- **Recently Created Files**: Viren tracks the last 10 files you exported. You can see this list in the history manager to verify your progress.

**Stop copy-pasting. Start building with Viren.**
