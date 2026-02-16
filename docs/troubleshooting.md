# Troubleshooting & Systems Diagnostics Manual

Viren is a high-performance tool that interacts with multiple cloud APIs, local subprocesses, and the host file system. When you encounter unexpected behavior, use this exhaustive guide to identify and resolve the root cause.

---

## 1. Network & API Connectivity

### Error: "API Key not found"
**Symptom**: Viren fails to start a chat and prints an error regarding a missing key.
- **Diagnostic**: Viren checks your environment variables (`OPENAI_API_KEY`, etc.).
- **Fix**: Run `export PLATFORM_API_KEY="your-key"`. 
- **Verify**: Type `env | grep API_KEY` to ensure the variable is actually exported to your current shell.
- **Persistent Fix**: Add the export line to your `~/.zshrc` or `~/.bashrc` and restart the terminal.

### Error: "429 Too Many Requests"
**Symptom**: The AI stops mid-stream or fails to respond with a 429 status code.
- **Cause**: You have hit the rate limit of your AI provider or have an insufficient account balance.
- **Fix**: Check your provider's billing dashboard (e.g., OpenAI or Anthropic). Most APIs require pre-paid credits.
- **Fix**: Switch to a faster platform like **Groq** if you are performing many small queries.

### Error: "Connection Refused" (Local Ollama)
**Symptom**: chat fails when platform is set to `ollama`.
- **Cause**: Viren tried to talk to `localhost:11434` but found no active listener.
- **Fix**: Ensure the Ollama service is running. Run `ollama serve` in a separate terminal or check your system tray icon.

---

## 2. UI & Terminal Interaction

### The screen is filled with "Strange Characters"
**Symptom**: You see raw ANSI escape codes like `[38;2;0;255;255m` instead of colored text.
- **Cause**: Your terminal emulator is legacy or does not support 24-bit TrueColor.
- **Fix**: Upgrade to a modern terminal like **Alacritty**, **Kitty**, **iTerm2**, or **Windows Terminal**.
- **Fix**: If you must stay on your current terminal, try `export TERM=xterm-256color`.

### Fzf menus are not appearing
**Symptom**: Typing `!m` or `!l` causes the app to hang or return an error.
- **Cause**: The `fzf` binary is missing from your system `$PATH`.
- **Fix (macOS)**: `brew install fzf`
- **Fix (Linux)**: `sudo apt install fzf` or `sudo pacman -S fzf`
- **Verify**: Type `fzf --version` in your terminal. Viren requires `fzf` for all interactive selections.

---

## 3. Build & Installation Hurdles

### Build fails with "CGO_ENABLED" errors
**Symptom**: `make build` returns errors about missing header files or GCC.
- **Cause**: Viren's OCR features require a C compiler. 
- **Fix (Linux)**: `sudo apt install build-essential`
- **Fix (macOS)**: `xcode-select --install`
- **Fix (Android/Termux)**: Always use `CGO_ENABLED=0 go build ...` as Termux has non-standard library paths.

### "viren: command not found"
**Symptom**: After running the installer, typing `viren` does nothing.
- **Cause**: The symlink was created in `/usr/local/bin`, but that directory is not in your shell's `$PATH`.
- **Fix**: Add `export PATH="/usr/local/bin:$PATH"` to your shell profile.

---

## 4. Performance & Context Issues

### Startup is slow (> 500ms)
**Diagnostic**: Viren is designed to start in under 100ms.
- **Cause**: Your `config.json` might be extremely large, or you are running Viren on a network-mounted filesystem (NFS).
- **Fix**: Check your `shallow_load_dirs` setting in `config.json`. Ensure Viren isn't scanning massive folders on launch.

### The AI is "Hallucinating" or giving wrong answers
**Diagnostic**: The prompt context might be cluttered or the system instructions are too weak.
- **Fix**: Use `!c` to clear the history and start fresh.
- **Fix**: Switch to **Zenith Mode** (`!v zenith`). This forces the model to use Chain-of-Thought reasoning, which significantly reduces errors.

---

## 5. Frequent Questions (FAQ)

**Q: Does Viren store my data in the cloud?**
A: **No.** Viren is a local-first binary. Your history and configs stay on your disk.

**Q: Can I use Viren without an internet connection?**
A: **Yes.** Install **Ollama**, download a model, and switch Viren to the `ollama` platform (`!p ollama`).

**Q: How do I change the default AI model?**
A: Open `~/.viren/config.json` and update the `"default_model"` string to your preferred model name (e.g., `claude-3-5-sonnet`).

---

## 6. How to Report a Bug
If your issue is not resolved by this guide:
1.  Run `viren -v` to get your exact version and commit hash.
2.  Capture the error output from your terminal.
3.  Open a detailed issue on our GitHub repository.
4.  Specify your OS (Linux, macOS, Windows) and shell (Zsh, Bash, Fish).

**We are committed to making Viren the most stable tool in your technical arsenal.**
