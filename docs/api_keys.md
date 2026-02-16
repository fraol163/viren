# Comprehensive Guide to AI Platforms & API Keys

Viren acts as a high-performance multiplexer for the world's most advanced Large Language Models. To enable this connectivity, you must provide valid API credentials for each platform. This guide details how to acquire these keys, the specific environment variables required, and the recommended models for each provider.

---

## 1. Primary AI Platforms

### OpenAI (GPT-4o, o1)
The industry standard for general-purpose reasoning and high-speed instruction following.
- **Acquisition**: [platform.openai.com/api-keys](https://openai.com/api/)
- **Env Variable**: `OPENAI_API_KEY`
- **Top Models**: `gpt-4o`, `gpt-4o-mini`, `o1-preview`, `gpt-4-turbo`.
- **Developer Tip**: Ensure your account has a positive credit balance to avoid `429 Insufficient Balance` errors.

### Anthropic (Claude 3.5)
Widely considered the best model for nuanced coding, architectural design, and complex writing.
- **Acquisition**: [console.anthropic.com/settings/keys](https://console.anthropic.com/)
- **Env Variable**: `ANTHROPIC_API_KEY`
- **Top Models**: `claude-3-5-sonnet-20240620`, `claude-3-opus-20240229`.
- **Developer Tip**: Use Claude when you need extremely long context (up to 200k tokens).

### DeepSeek (DeepSeek-V3)
The current leader in price-to-performance ratio for coding. DeepSeek models are highly optimized for technical tasks.
- **Acquisition**: [platform.deepseek.com/api_keys](https://api-docs.deepseek.com/)
- **Env Variable**: `DEEP_SEEK_API_KEY`
- **Top Models**: `deepseek-chat`, `deepseek-coder`.
- **Developer Tip**: DeepSeek is excellent for "Codedump" analysis because of its aggressive pricing.

### Google Gemini (Pro & Flash)
Google's native models, featuring massive context windows and high multimodal capabilities.
- **Acquisition**: [aistudio.google.com/app/apikey](https://ai.google.dev/gemini-api/docs/api-key)
- **Env Variable**: `GEMINI_API_KEY`
- **Top Models**: `gemini-1.5-pro`, `gemini-1.5-flash`.
- **Developer Tip**: Check Gemini for free-tier access, which is often available for development testing.

---

## 2. Infrastructure & Search Platforms

### Brave Search API
Required for the `!w` (Web Search) capability. This allows Viren to browse the live internet.
- **Acquisition**: [api.search.brave.com/app/dashboard](https://brave.com/search/api/)
- **Env Variable**: `BRAVE_API_KEY`
- **Why?**: Standard LLMs have a knowledge cutoff. Using Brave allows Viren to feed the model the latest documentation or news from today.

### Groq (Ultra-Fast Inference)
Specializes in extreme speed by using LPU (Language Processing Unit) hardware.
- **Acquisition**: [console.groq.com/keys](https://console.groq.com/keys)
- **Env Variable**: `GROQ_API_KEY`
- **Top Models**: `llama-3.1-70b-versatile`, `mixtral-8x7b-32768`.
- **Why?**: Use Groq when you want the AI response to appear instantly (<1 second).

### OpenRouter (The Aggregator)
A single API that connects you to every model (Llama, Qwen, Mistral, etc.) without needing 20 different accounts.
- **Acquisition**: [openrouter.ai/settings/keys](https://openrouter.ai/settings/keys)
- **Env Variable**: `OPENROUTER_API_KEY`
- **Why?**: If you only want to manage one wallet and one key, OpenRouter is the best choice.

---

## 3. Specialized & Enterprise Platforms

### xAI (Grok)
Focuses on real-time awareness and less restrictive output filters.
- **Acquisition**: [console.x.ai](https://x.ai/api)
- **Env Variable**: `XAI_API_KEY`

### Mistral AI
European-based models known for high efficiency and transparency.
- **Acquisition**: [console.mistral.ai](https://docs.mistral.ai/getting-started/quickstart)
- **Env Variable**: `MISTRAL_API_KEY`

### Together AI
The home of open-source model inference at scale.
- **Acquisition**: [api.together.ai](https://docs.together.ai/docs/quickstart)
- **Env Variable**: `TOGETHER_API_KEY`

### Amazon Bedrock
Enterprise-grade access via AWS infrastructure.
- **Acquisition**: [aws.amazon.com/bedrock](https://aws.amazon.com/bedrock)
- **Env Variables**: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.

---

## 4. Setting Environment Variables (OS Specific)

### Linux & macOS (Bash/Zsh)
Open your profile file:
```bash
nano ~/.zshrc  # or ~/.bashrc
```
Add the following block:
```bash
# Viren API Credentials
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export DEEP_SEEK_API_KEY="sk-..."
export BRAVE_API_KEY="bs-..."
```
Reload the shell:
```bash
source ~/.zshrc
```

### Windows (PowerShell)
Run this command for each key:
```powershell
[System.Environment]::SetEnvironmentVariable('OPENAI_API_KEY', 'your-key-here', [System.EnvironmentVariableTarget]::User)
```
*Note: Restart your terminal after setting these.*

---

## 5. Security & Rotation Best Practices
1.  **Never hard-code keys**: Viren will never ask you to type a key into the `config.json`.
2.  **Use Scoped Keys**: If a provider allows it, create a key specifically for Viren so you can revoke it without breaking other apps.
3.  **Monitor Usage**: Check your provider dashboards weekly to ensure no unauthorized usage has occurred.

**Viren treats your secrets with the same respect as your code.**
