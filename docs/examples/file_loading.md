# Real-World Example: Multi-File Context Ingestion

This guide covers the most powerful feature of Viren: giving the AI access to your actual code files to solve real bugs.

## Scenario: The "Mystery" Crash
You have a Python project where a function in `utils.py` is called by `main.py`, but it crashes with a `KeyError` only under certain conditions.

### Step 1: Loading the Context
Instead of copy-pasting both files, use the Load command:
`USER ❯ !l`

Viren opens an `fzf` menu.
1. Scroll to `main.py`, press `Tab` to select.
2. Scroll to `utils.py`, press `Tab` to select.
3. Press `Enter`.

### Step 2: The Inquiry
Viren has now "read" both files. Ask your question:
`USER ❯ In utils.py, the function fetch_data() is causing a KeyError. Looking at how it's called in main.py, why is this happening?`

### Step 3: Architectural Analysis
Because the AI sees **both** files, it can perform cross-file reasoning:
- *"The AI notices that main.py passes a dictionary without the 'id' key, which fetch_data in utils.py expects on line 42."*

### Step 4: The Fix
The AI suggests a fix. Use `!e` to save the updated `utils.py` directly.

---

## Best Practices for File Loading
- **Relative Paths**: Viren handles relative paths automatically. You can load files from subdirectories (e.g., `src/auth/logic.go`).
- **PDF/Documentation**: Use `!l` to load a PDF manual for a library you are using. Viren's native PDF parser will extract the text so you can ask: *"According to the manual, what is the correct retry policy for this client?"*
- **Directory Dumps (`!d`)**: If you have 50 files, don't use `!l`. Use `!d` to bundle the entire project into a single stream. Viren will automatically respect your `.gitignore` so it doesn't waste tokens on `dist/` or `__pycache__/`.
