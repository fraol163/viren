# Contributing to the Viren Ecosystem

Viren is a community-driven project built by developers, for developers. We are obsessed with performance, privacy, and minimalist elegance. If you share these values, we welcome your contributions. This document outlines our technical standards and the workflow for submitting improvements.

---

## 1. Our Core Technical Standards

### A. Performance is a Feature
Viren's primary competitive advantage is its sub-100ms startup time. 
- **The 5ms Rule**: Any pull request that increases the cold-start time of Viren by more than 5 milliseconds will be rejected unless the feature is critical and cannot be backgrounded.
- **Dependency Policy**: We strictly prefer the Go standard library (`net/http`, `encoding/json`, `os`). Every external dependency must be justified by a 10x improvement in functionality that cannot be achieved natively.

### B. Clean & Idiomatic Go
We follow the principles outlined in **Effective Go**.
- **Formatting**: Always run `go fmt ./...` before committing.
- **Vetting**: Your code must pass `go vet ./...` without warnings.
- **Concurrency**: Use channels and waitgroups properly. Avoid "leaking" goroutines that stay alive after a chat stream finishes.

### C. UI Consistency
Viren uses a specific ANSI design language.
- **Coloring**: Never hard-code ANSI escape sequences in logic files. Use the methods provided in the `Terminal` struct.
- **Feedback**: Use `PrintSuccess()`, `PrintError()`, and `PrintInfo()` to ensure the user experience remains consistent across all commands.

---

## 2. Areas for High-Impact Contribution

### New Domain Modes (`internal/chat/modes.go`)
We want Viren to be an expert in every technical niche. If you are a specialist in:
- Kernel Development
- Quantitative Finance
- FPGA/Hardware Engineering
- High-Frequency Trading
- Theoretical Physics
Help us write and test a hyper-optimized "System Prompt" for that specific domain.

### Platform Adapters (`internal/platform/`)
As new AI providers emerge (or existing ones change their API schemas), we need to update our translation layers. If you find a new provider with a unique model, help us build an adapter for it.

### Documentation & Real-World Examples
Clear documentation is as valuable as clean code. Improving the `docs/` or adding new scenarios to `docs/examples/` is the fastest way to help new users.

---

## 3. The Development Workflow

### Step 1: Fork & Environment Setup
```bash
git clone https://github.com/your-username/viren.git
cd viren
# Ensure you have Go 1.21+ installed
go version
```

### Step 2: Creating a Feature Branch
Use descriptive branch names:
- `feat/add-mistral-large`
- `fix/fzf-window-resize`
- `docs/update-installation-guide`

### Step 3: Build & Local Testing
Use the provided `Makefile` to ensure your build environment matches production.
```bash
make build
# Run the binary to test your changes
./bin/viren
```

### Step 4: Submitting a Pull Request
When opening a PR, your description must answer:
1.  **What** is changed?
2.  **Why** is this change beneficial for a developer?
3.  **Performance Impact**: Did you run `time ./bin/viren -v`?

---

## 4. UI/UX Philosophy
- **Keyboard First**: Every feature must be accessible without a mouse.
- **High Density**: Don't waste vertical space. Keep prompts and responses compact.
- **Non-Intrusive**: Avoid pop-ups or interrupting the user's typing flow.

---

## 5. Community & Ethics
- **Professionalism**: Be respectful and technical in all discussions.
- **Privacy First**: Never suggest a feature that requires a centralized Viren database or tracking.
- **Credit**: If you are inspired by or port logic from another project, always provide proper attribution in the code and the `README.md`.

---

## 6. Financial Contributions
Currently, Viren is a labor of love. If you wish to support the project financially, please contact the lead maintainers regarding sponsorship opportunities or bounty programs for specific features.

**Thank you for helping us build the future of the command line.**
