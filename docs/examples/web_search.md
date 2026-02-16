# Real-World Example: Advanced Web Search

This guide demonstrates how to use Viren's `!w` command to overcome the knowledge cutoff limits of standard AI models.

## The Problem
Imagine you are working with a new library, for example, **"Go 1.24"**, which was released after the AI model's training data was collected. If you ask a standard model about it, it will either hallucinate or tell you it doesn't know.

## The Viren Solution: Live Context Ingestion

### Step 1: Configuration
Ensure you have your Brave Search API key set:
```bash
export BRAVE_API_KEY="your-key-here"
```

### Step 2: The Direct CLI Method
If you just want a quick summary without entering interactive mode:
```bash
viren -w "Go 1.24 release notes and new features" "Summarize the top 5 changes for developers"
```

**What happens behind the scenes:**
1. Viren uses the Brave API to find the most relevant web pages.
2. It scrapes the text content from the top results.
3. It bundles this "Live Context" with your prompt.
4. The AI processes the web data and gives you an accurate answer based on *current* facts.

### Step 3: The Interactive Method
If you are already inside a chat:
1. Type `!w`.
2. Viren will ask for a query. Type `Go 1.24 new features`.
3. Viren performs the search and injects the results into your current conversation.
4. You can now ask follow-up questions: *"Show me a code example using the new feature X."*

---

## Pro-Tips for Web Search
- **Be Specific**: Instead of `!w react`, use `!w React 19 concurrent rendering changes`.
- **Combine with Exports**: After searching for a solution, use `!e` to save the generated code immediately.
- **Multilingual**: Viren's search uses your `config.json` locale settings, but you can override them for global results.
