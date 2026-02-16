# Troubleshooting & Systems Diagnostics

Viren is a complex tool interacting with multiple cloud APIs, local subprocesses, and the file system. When something goes wrong, use this guide to identify and resolve the issue.

---

## 1. Network & API Failures

### Error: "API Key not found"
**Diagnostic**: Viren checked your environment and found nothing for the selected platform.
- **Fix**: Run `export PLATFORM_API_KEY="your-key"`. 
- **Verify**: Run `env | grep API_KEY` to ensure it is correctly exported. 
- **Persistent Fix**: Add the export line to your `~/.zshrc` or `~/.bashrc`.

### Error: "429 Too Many Requests"
**Diagnostic**: You have hit the rate limit of your AI provider or have an insufficient balance.
- **Fix**: Check your provider's billing dashboard. OpenAI and Anthropic require pre-paid credits for API access.

### Error: "Connection Refused" (Local Ollama)
**Diagnostic**: Viren tried to talk to `localhost:11434` but found no listener.
- **Fix**: Ensure the Ollama service is running. Run `ollama serve` in a separate terminal.

---

## 2. UI & Interaction Issues

### The screen is filled with "Strange Characters" (ANSI Artifacts)
**Diagnostic**: Your terminal emulator does not support ANSI escape codes or 24-bit color.
- **Fix**: Use a modern terminal like **Alacritty**, **Kitty**, **iTerm2**, or **Windows Terminal**.
- **Fix**: If on an old system, try `export TERM=xterm-color`.

### Fzf menus are not appearing
**Diagnostic**: The `fzf` binary is missing from your `$PATH`.
- **Fix**: Install fzf.
    - **macOS**: `brew install fzf`
    - **Linux**: `sudo apt install fzf`
- **Verify**: Type `fzf --version` in your terminal.

---

## 3. Build & Installation Issues

### Build fails with "CGO_ENABLED" errors
**Diagnostic**: Your system lacks a C compiler (like `gcc` or `clang`) or is trying to cross-compile incorrectly.
- **Fix (Linux)**: `sudo apt install build-essential`
- **Fix (macOS)**: `xcode-select --install`
- **Fix (Android)**: Always use `CGO_ENABLED=0 go build ...` on Termux.

### "viren: command not found" after installation
**Diagnostic**: `/usr/local/bin` is not in your shell's `$PATH`.
- **Fix**: Run `export PATH="/usr/local/bin:$PATH"` or check the installer logs for the symlink location.

---

## 4. Frequent Questions (FAQ)

### Q: Does Viren support GPT-4o with Images?
**A**: Currently, Viren handles text-based ingestion. Support for direct image uploading (`!l image.png`) is planned for Phase 2.

### Q: Why is my history not saving?
**A**: Check if `"enable_session_save"` is set to `true` in your `config.json`. Also, ensure Viren has write permissions to `~/.viren/tmp/`.

### Q: How do I wipe everything?
**A**: Run `./install.sh --safe-uninstall`. This will remove the binary, the config, and all conversation logs.

---

## 5. Contact & Reporting
If your issue persists:
1.  Run `viren -v` to get your version info.
2.  Open an issue on GitHub with your OS details and the exact error message.
3.  Include whether you were in Interactive mode or Single-Shot mode.

**We are committed to making Viren the most stable tool in your kit.**
