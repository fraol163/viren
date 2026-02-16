# Master User Guide: Becoming a Viren Power User

Welcome to Viren. This guide is designed to take you from a basic user to a master of terminal-based AI automation. Viren is not just a chat tool; it is a context-injection engine that allows you to treat LLMs as a first-class citizen of your operating system.

---

## 1. Documentation Philosophy
**Why is Viren documentation structured this way?**
We follow the **High-Density Information** principle. Developers don't want fluff; they want facts and syntax. Our documentation is designed to be browsed via `grep` or read quickly during a build cycle. Every file in the `docs/` folder is a modular component of the Viren ecosystem.

---

## 2. Core Workflow: The First 5 Minutes

### Initial Setup
1.  **Install**: `curl -sSL ... | bash`
2.  **Keys**: Export your `OPENAI_API_KEY` or `ANTHROPIC_API_KEY`.
3.  **First Run**: Type `viren`.
4.  **Onboarding**: Fill in your profile. Be specific about your role (e.g., "Embedded Systems Engineer") as this directly changes how the AI explains complex concepts.

### Your First Session
Try asking: *"What is the best way to handle concurrency in this environment?"*
Then, try switching the logic:
1.  Type `!v`.
2.  Select `Socratic Mode`.
3.  Ask the same question. Notice how the AI no longer gives you the answer but asks you questions to help you figure it out.

---

## 3. Master Command Table

| Command | Usage | Pro-Tip |
| :--- | :--- | :--- |
| `!q` | Exit | Use this to ensure your session is saved to history. |
| `!h` | Help | Use the search bar in the fzf menu to find commands instantly. |
| `!c` | Clear | Run this when you're switching from "Debugging" to "Creative Writing." |
| `!m` | Model | Check for newer models (like `o1-preview`) regularly. |
| `!p` | Platform | Use `!p groq` for instant replies, `!p anthropic` for deep logic. |
| `!u` | Personality | Try `Rick Sanchez` mode for a humorous, sarcastic coding session. |
| `!v` | Domain Mode | Use `Code Whisperer` for purely idiomatic refactors. |
| `!z` | Theme | Use `Paper` for bright environments, `DeepSpace` for OLED screens. |
| `!x` | Shell Record | Perfect for fixing failing `Makefile` or `docker-compose` errors. |
| `!d` | Codedump | Always exclude `node_modules` and `.git` to save tokens. |
| `!l` | Load | You can load PDFs! Try loading a datasheet or a technical spec. |
| `!s` | Scrape | Extract clean text from any URL without the ads and banners. |
| `!w` | Web Search | Use this to verify facts that occurred in the last 24 hours. |
| `!a` | History | You can restore a session from 3 weeks ago in 2 seconds. |
| `!y` | Clipboard | Use "Turn Copy" to share a full conversation with a teammate. |

---

## 4. Context Injection Mastery

### The "Load" Protocol (`!l`)
The `!l` command is the heart of Viren's productivity.
- **Multiple Files**: Press `Tab` in the fzf menu to select 5-10 files. Viren will ingest them all.
- **OCR Support**: If you load a `.jpg` or `.png`, and have `tesseract` installed, Viren will extract the text from the image. This is great for screenshots of error messages or whiteboard diagrams.
- **PDF Parsing**: Viren includes a native PDF text extractor. Load a manual, and the AI becomes an expert on that documentation.

### The Codedump (`!d`)
Don't copy-paste your code. Use `!d` to give the AI a "God's Eye View" of your project.
- **Respecting .gitignore**: Viren is smart. It won't load your `.env` files or secret keys if they are ignored in git.
- **Exclusion**: Use the interactive menu to prune your context. The smaller the context, the faster and more accurate the AI.

---

## 5. Advanced Features: The Developer's Edge

### Shell Recording (`!x`)
This is Viren's most powerful feature for systems work.
1.  Launch `!x`.
2.  Run `npm install`. If it crashes with a massive stack trace, don't scroll up to copy it.
3.  Type `exit`.
4.  Viren hands the *entire* log to the AI and asks for a fix.
5.  The AI analyzes the environment and the error simultaneously.

### Smart Export (`!e`)
Viren's export engine is context-aware.
- It detects the language of the code block.
- It suggests a name (e.g., `main.go` or `style.css`).
- It checks if the file exists and prompts for an overwrite or a new name.

---

## 6. Offline & Local AI (Ollama)
Viren is a bridge to the future of local AI.
- **Privacy**: If you are working on sensitive government or corporate code, use `!p ollama`.
- **Performance**: Local models have zero network latency.
- **Requirement**: You must have the [Ollama](https://ollama.com) service running in the background.

---

## 7. Troubleshooting Common Workflows
- **AI is Hallucinating?**: Switch to `Zenith Mode` (`!v`) to force better reasoning.
- **Response is too long?**: Switch to `Focused` personality (`!u`).
- **Context limit reached?**: Use `!c` to clear the history and start fresh with only the essential files loaded.

**Viren is your new technical partner. Use it wisely.**
