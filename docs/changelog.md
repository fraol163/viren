# Viren Project Changelog

All notable changes to the Viren project will be documented in this file. This project adheres to **Semantic Versioning (SemVer)**.

---

## [1.0.0] - 2026-02-16
### Rebrand & Evolution
- **The Viren Shift**: Official project rename from "Cha" to "Viren." This marks the transition from a prototype to a production-ready software suite.
- **Engine Rebuild**: The core internal loop was rewritten in Go to eliminate the overhead found in hybrid script approaches.
- **Performance**: Achieved <100ms startup times across all tested Unix environments.

### New Features
- **Codedump v2**: Reimplemented directory scanning to respect `.gitignore` and added an `fzf` multi-select menu for precise context control.
- **Shell Record (`!x`)**: Introduced a native sub-shell capture system that allows the AI to "see" your terminal output directly.
- **Theme Matrix**: Added 4 distinct ANSI themes: `DeepSpace`, `Neon`, `Paper`, and `Matrix`.
- **Domain Modes**: Expanded the logic engine to include 40+ specialized expert modes.
- **Version Bumping**: Added a semantic versioning tool to the installer for developers.

### Improvements & Fixes
- **Installer Security**: Moved all temporary build files to `/tmp` to ensure system cleanliness.
- **Flag Parsing**: Fixed a critical bug where `-h` or `-v` would trigger the onboarding wizard.
- **Documentation**: Established a 10-file comprehensive documentation suite in the `docs/` folder.
- **OCR Integration**: Optional Tesseract support for image-to-text context injection.

---

## [3.9.3] - Legacy (The "Cha" Era)
- Final stable release of the original interactive CLI prototype.
- Initial implementation of the "bang-command" system.
- Support for OpenAI and Anthropic API protocols.
- Basic session persistence and JSON history logging.

---

### Legend
- **Added**: For new features.
- **Changed**: For changes in existing functionality.
- **Deprecated**: For soon-to-be removed features.
- **Removed**: For now removed features.
- **Fixed**: For any bug fixes.
- **Security**: In case of vulnerabilities or security updates.

**Viren follows a continuous delivery model for documentation and stability.**
