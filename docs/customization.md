# Comprehensive Configuration & Logic Tuning

Viren is designed to be a transparent extension of your technical personality. Its sophisticated configuration engine allows you to manipulate every aspect of the application, from the core behavioral logic of the AI to the visual aesthetic of the terminal interface.

---

## 1. The Configuration File: `config.json`

Viren stores its persistent state in `~/.viren/config.json`. 

### Full Schema Reference
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
  "shallow_load_dirs": [
    "/",
    "/home/",
    "/usr/",
    "/tmp/"
  ],
  "user_profile": {
    "name": "Alex",
    "role": "Lead Backend Engineer",
    "environment": "Arch Linux / Tmux / Neovim",
    "ambition": "Building low-latency distributed systems"
  }
}
```

---

## 2. Behavioral Personalities (`!u`)

Viren allows you to switch between distinct personalities that alter the AI's tone, verbosity, and style.

- **Analytical**: (ID: `analytical`) Logical, systematic, and data-driven. Focuses on proofs and reasoning.
- **Creative**: (ID: `creative`) Artistic, imaginative, and metaphor-heavy. Great for UI/UX brainstorming.
- **Focused**: (ID: `focused`) Goal-oriented and concise. Removes all "filler" text for maximum speed.
- **Empathetic**: (ID: `empathetic`) Emotionally intelligent and supportive. Good for explaining complex topics.
- **Playful**: (ID: `playful`) Fun, energetic, and full of wit.
- **Balanced**: (ID: `balanced`) The default professional standard.
- **Rick Sanchez**: (ID: `rick`) Cynical, sarcastic, and unhinged. Uses the distinct voice of the scientist from Rick and Morty.

---

## 3. Domain Modes (`!v`)

A "Mode" in Viren is a specialized behavioral template that re-architects the entire system prompt for specific technical fields.

### Core Modes
- **Standard**: Professional balanced tone for general queries.
- **Zenith**: High-level reasoning using cosmic metaphors and philosophical depth.
- **Code Whisperer**: Expert level pair-programmer. Strips all fluff and focuses on idiomatic code.
- **Socratic**: Pedagogy mode. The AI never answers directly; it only asks questions to guide you.
- **Complexity Analyzer**: Specifically tuned to calculate and explain Big-O time and space complexity.

### Specialized Technical Modes
- **AlgoThink**: Software engineering and system design patterns.
- **DSA Mode**: Data Structures & Algorithms specialist.
- **CyberSec**: Security auditing, penetration testing, and vulnerability research.
- **Physics/Chemistry/Bio**: Specialized prompts for hard sciences.
- **FinancePro**: Market analysis, risk modeling, and economic forecasting.

---

## 4. Safety & Performance: `shallow_load_dirs`

Recursive file scanning is dangerous in root directories. Any path listed in the `shallow_load_dirs` array will **never be scanned deeper than 1 level**. 
- **Recommendation**: Always include your system's root (`/`) and home folder (`/home/user`) to protect against accidental recursive scanning during a `!d` or `!l` command.

---

## 5. Neural Profile Injections

The `user_profile` section allows the AI to understand your context from the first message.
- **Role Tuning**: If set to "DevOps," the AI will prioritize automation and infrastructure.
- **Env Awareness**: If set to "Vim," the AI will suggest terminal commands over GUI instructions.

---

## 6. Windows Installation Details

Viren is fully production-ready for Windows users.

### Environment variables
Windows users should set API keys via PowerShell:
```powershell
[System.Environment]::SetEnvironmentVariable('OPENAI_API_KEY', 'your-key', [System.EnvironmentVariableTarget]::User)
```

### Path configuration
Ensure the directory containing `viren.exe` is added to your User `PATH` environment variable so you can launch it from any terminal.

**Viren is your digital workshop. Configure it to be the ultimate extension of your technical mind.**
