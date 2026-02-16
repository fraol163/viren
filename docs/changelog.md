# Viren Project Changelog: Evolution of the Neural Interface

All notable changes to the Viren project are documented here. This project adheres to **Semantic Versioning (SemVer)**. We prioritize stability and performance in every release.

---

## [1.0.0] - 2026-02-16
### The "Viren" Rebrand & Engine Rebuild
This release marks the official transition from the legacy "Cha" prototype to the production-ready **Viren** suite. The core engine has been completely re-engineered in Go.

### Added
- **Asynchronous Core**: Implemented a non-blocking I/O model using Go routines for sub-100ms response times.
- **Multi-Platform Multiplexer**: Native support for 11 providers including DeepSeek, Groq, Anthropic, and Amazon Bedrock.
- **Codedump v2**: A completely rewritten directory scanning engine that respects `.gitignore` and supports `fzf` multi-selection.
- **Native PDF/Office Support**: Integrated text extraction for `.pdf`, `.docx`, `.xlsx`, and `.csv` files.
- **Shell Record (`!x`)**: A native sub-shell capture system for real-time debugging and log analysis.
- **Advanced UI Themes**: Introduced `DeepSpace`, `Neon`, `Paper`, and `Matrix` ANSI themes.
- **Domain Logic Engine**: Added 40+ specialized modes (Zenith, Code Whisperer, etc.).
- **Smart Export**: Regex-based code detection with intelligent filename suggestion.
- **Ollama Integration**: First-class support for 100% local, offline AI reasoning.

### Changed
- **Installer Logic**: Moved build artifacts to `/tmp` and improved dependency recovery for Linux, macOS, and Android.
- **Flag Parsing**: Restructured `main.go` to ensure `-h` and `-v` flags do not trigger the onboarding wizard.
- **Documentation**: Established a 10-file comprehensive documentation suite in the `docs/` folder.

### Fixed
- Fixed a memory leak occurring during extremely long streaming responses.
- Resolved an issue where `fzf` would hang if invoked in a non-interactive shell.
- Corrected the CGO build flags for ARM64 macOS systems.

---

## [3.9.3] - Legacy (The "Cha" Era)
*The final stable release of the original interactive CLI prototype.*

### Added
- Initial "Bang-Command" (`!`) implementation.
- Support for OpenAI and Anthropic API protocols.
- Basic JSON history logging.
- Support for standard Unix pipes.

### Fixed
- Corrected a bug in the token estimation logic for long files.
- Improved the cleanup of temporary files after an export session.

---

## [2.5.0] - Legacy Prototype
- Initial implementation of the `!l` (Load) command.
- Support for basic ANSI coloring.
- Simple configuration file (`config.json`).

---

## Release Legend
- **Added**: New features available to the end-user.
- **Changed**: Updates to existing functionality or logic.
- **Deprecated**: Features that will be removed in future versions.
- **Removed**: Features that have been purged from the binary.
- **Fixed**: Bug fixes and stability improvements.
- **Security**: Specific patches for data privacy or credential handling.

**Viren follows a continuous delivery model for documentation and stability.**
