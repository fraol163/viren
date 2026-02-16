# Deep Customization & Logic Tuning

Viren is designed to be a transparent extension of your technical personality. This document explains how to manipulate the internal configuration engine to fit your specific hardware, OS, and professional requirements.

---

## 1. The Core Config: `config.json`
Viren stores its persistent state in `~/.viren/config.json`. 

### Full Schema Definition
```json
{
  "default_model": "gpt-4o",
  "current_platform": "openai",
  "exit_key": "!q",
  "model_switch": "!m",
  "editor_input": "!t",
  "clear_history": "!c",
  "help_key": "!h",
  "export_chat": "!e",
  "backtrack": "!b",
  "web_search": "!w",
  "show_search_results": true,
  "num_search_results": 5,
  "scrape_url": "!s",
  "copy_to_clipboard": "!y",
  "load_files": "!l",
  "platform_switch": "!p",
  "all_models": "!o",
  "code_dump": "!d",
  "shell_record": "!x",
  "multi_line": "\\",
  "preferred_editor": "vim",
  "current_theme": "deepspace",
  "current_mode": "standard",
  "current_personality": "balanced",
  "mute_notifications": false,
  "enable_session_save": true,
  "shallow_load_dirs": [
    "/",
    "/home/",
    "/usr/",
    "/etc/",
    "/var/",
    "/tmp/"
  ],
  "user_profile": {
    "name": "Jane Doe",
    "role": "Lead DevOps Engineer",
    "environment": "Arch Linux / Kitty / Tmux",
    "ambition": "Automating cloud infrastructure at scale"
  }
}
```

---

## 2. Advanced Parameter In-Depth

### `shallow_load_dirs` (Safety Mechanism)
Recursive file scanning is dangerous in root directories. 
- **The Problem**: If you run `!d` (Codedump) in `/home/user`, Viren might try to read your browser cache, logs, and downloads, leading to massive context usage and system lag.
- **The Solution**: Any path listed in this array will **never be scanned deeper than 1 level**. 
- **Recommendation**: Always include your system's root and your user's home directory here.

### `user_profile` (Contextual Tuning)
The values here are injected into the "Neural Profile" segment of every prompt.
- **Role**: Changes the AI's technical vocabulary.
- **Environment**: If you are a "Vim" user, the AI will provide `:w` commands instead of telling you to click a "Save" icon.
- **Ambition**: Provides a long-term goal for the AI to keep in mind (e.g., "Refactor for readability over speed").

---

## 3. The Logic Engine: Domain Modes (`!v`)
A "Mode" in Viren is a specialized set of system instructions that override the default behavior.

| Mode | Functional Shift |
| :--- | :--- |
| **Zenith** | High-level reasoning. It uses COSMIC metaphors and focuses on philosophical first principles. |
| **Code Whisperer** | Expert Software Engineer. Removes all "Sure, I can help" filler. Focuses purely on idiomatic code. |
| **Socratic** | Guiding Teacher. Will never give an answer. Will only ask questions to lead you to the truth. |
| **DSA Mode** | Data Structures & Algorithms. Focuses on Big-O notation, memory safety, and optimization. |
| **CyberSec** | Focuses on vulnerability assessment, penetration testing patterns, and OWASP standards. |

---

## 4. Visual customization: Theming Engine (`!z`)
Viren uses a custom-built ANSI renderer. You can switch themes instantly without restarting.

- **DeepSpace**: Default. Optimized for OLED and high-contrast dark terminals.
- **Neon**: Vibrant colors for low-light environments.
- **Paper**: Light mode. Optimized for high-glare environments or daylight coding.
- **Matrix**: Monochrome green. High focus, minimal distractions.

---

## 5. Overriding Configuration
Sometimes you need to change settings for just one session. Viren supports **Environment Overrides**:

```bash
# Start Viren with a specific model and theme for this run only
VIREN_DEFAULT_PLATFORM=anthropic VIREN_DEFAULT_MODEL=claude-3-opus viren
```

**Viren is your tool. Bend it to your workflow.**
