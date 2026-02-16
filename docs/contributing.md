# Contributing to the Viren Ecosystem

We are building Viren to be the industry-standard AI CLI. To achieve this, we need contributors who are obsessed with performance, clean code, and developer ergonomics.

---

## 1. Our Technical Standards

### Performance is a Requirement
If your pull request increases the startup time of Viren by more than 5ms, it will be scrutinized. We use Go specifically for its low-overhead execution.
- Avoid large external dependencies.
- Use the standard library whenever possible (`net/http`, `encoding/json`).
- Optimize file scanning logic using concurrency where appropriate.

### Clean & Idiomatic Go
Follow the principles of "Effective Go."
- Run `go fmt ./...` before every commit.
- Use `go vet ./...` to catch common mistakes.
- Ensure all exported functions have clear comments.

---

## 2. Areas for Contribution

### Domain Modes (`internal/chat/modes.go`)
We want to expand our "Expert Modes." If you are a specialist in a field (e.g., Quantum Physics, Kernel Development, High-Frequency Trading), help us write a hyper-optimized system prompt for that domain.

### Platform Support (`internal/platform/`)
Help us integrate new AI providers as they emerge. Each provider should implement the standard `SendChatRequest` flow.

### Documentation & Examples
The most impactful way to help new users is by improving the `docs/` and `docs/examples/` folders. Clear, concise guides are as valuable as code.

---

## 3. The Development Workflow

1.  **Fork & Clone**:
    ```bash
    git clone https://github.com/your-username/viren.git
    cd viren
    ```
2.  **Create a Branch**:
    ```bash
    git checkout -b feat/my-new-feature
    ```
3.  **Build and Test**:
    Use the `Makefile` to keep your environment consistent.
    ```bash
    make build
    make test
    ```
4.  **Linting**:
    If you have `golangci-lint` installed, please run it.
5.  **Submit PR**:
    Ensure your PR description explains **why** the change is beneficial, not just what was changed.

---

## 4. UI Consistency
Viren uses a specific ANSI design language.
- Use `Terminal.PrintSuccess`, `Terminal.PrintError`, and `Terminal.PrintInfo` for all user feedback.
- Do not hard-code colors; use the variables defined in the current `Theme`.

---

## 5. Community & Ethics
- Be professional and respectful in all interactions.
- Prioritize **Privacy**. Never suggest a feature that requires a centralized Viren database or tracking.

**Thank you for helping us build the future of the command line.**
