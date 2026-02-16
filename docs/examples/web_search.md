# Real-World Scenario: Mastering Live Web Context with Brave Search

This guide provides a comprehensive walkthrough of using Viren's `!w` (Web Search) command to bridge the gap between an AI model's training data and the current state of the world.

---

## 1. The Challenge: The Knowledge Cutoff
AI models, no matter how powerful, are static snapshots. If you are a developer working with:
- A library released 2 months ago.
- A critical security CVE announced this morning.
- The latest documentation for a fast-moving project like Go or Rust.
A standard model will likely hallucinate or admit it doesn't know. 

---

## 2. The Solution: Live Context Ingestion

### Step 1: Secure Your Credentials
Viren requires a Brave Search API key to perform live lookups.
1.  Go to [api.search.brave.com](https://api.search.brave.com/app/dashboard).
2.  Register for a Free or Pro account.
3.  Add the key to your environment:
    ```bash
    export BRAVE_API_KEY="your-bs-key-here"
    ```

### Step 2: The Direct CLI Pattern
Use this for quick fact-checks without entering the interactive shell.
```bash
viren -w "Go 1.24 crypto library changes" "Summarize the top 3 changes for developers"
```
**Process Breakdown**:
1.  Viren hits the Brave API.
2.  It retrieves the top 5 most relevant URLs.
3.  It scrapes the textual content of those pages.
4.  It feeds the scraped text + your question to the AI.
5.  **Result**: You get an answer based on facts that are only hours old.

---

## 3. The Interactive Pattern: Research Mode

If you are already inside a Viren session, you can switch to "Research Mode" dynamically.

1.  **Launch**: `viren`
2.  **Search**: `USER ❯ !w`
3.  **Prompt**: `Enter search query: Go 1.24 release notes`
4.  **Analysis**: Viren will show a "Searching..." animation. Once finished, it will output: *"Web search results loaded into context."*
5.  **Query**: `USER ❯ Based on the search results, how do I implement the new TLS 1.3 features?`

---

## 4. Advanced Combinations: Web + Export

The true power of Viren is combining these tools. 

**Workflow Example**:
1.  `!w`: Find the latest API spec for a new service.
2.  `USER ❯`: "Write a client implementation for this spec."
3.  `!e`: Automatically detect the code block and save it as `client.go`.

---

## 5. Pro-Tips for High-Density Search
- **Specificity**: Instead of `!w rust`, use `!w Rust 1.80 stabilization of LazyLock`.
- **Locale Control**: Viren respects your `config.json` search settings. If you need international results, update the `search_country` and `search_lang` fields.
- **Cost Management**: Each search uses tokens. Use `!c` (Clear) if you have finished a research task to prevent the large search context from being sent with every subsequent prompt.

**Viren turns your AI into a live researcher.**
