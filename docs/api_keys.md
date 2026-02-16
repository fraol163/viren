# Global AI Platform Integration & Credential Security Manual

Viren is a high-performance multiplexer for the world's most advanced Large Language Models. To achieve seamless connectivity, you must provide valid API credentials for the platforms you wish to utilize. This guide provides an exhaustive, production-grade breakdown of every supported platform, specific model recommendations, acquisition links, and the security protocols for managing your digital secrets.

---

## 1. Documentation Philosophy: Security First

Viren follows a **"Zero-Key Persistence"** architecture. Unlike many consumer AI apps that store your API keys in a local JSON file or a cloud database, Viren reads your secrets directly from your operating system's environment. This is the industry-standard method for avoiding "Secret Leakage," ensuring that even if you share your `config.json` or your conversation history, your financial credentials remain safe.

---

## 2. Primary Reasoning & Coding Platforms

### OpenAI (The Industry Benchmark)
OpenAI's GPT models are the foundation of modern conversational AI. They offer the most stable API performance and are the primary testing ground for many of Viren's features.
- **Acquisition Link**: [https://platform.openai.com/api-keys](https://openai.com/api/)
- **Required Env Variable**: `OPENAI_API_KEY`
- **Recommended Models**:
    - `gpt-4o`: The best for complex logic, multi-file debugging, and architectural design.
    - `gpt-4o-mini`: Extremely fast and cost-effective. Use this for single-file scripts.
    - `o1-preview`: A reasoning-heavy model that "thinks" before replying. Essential for algorithmic problem solving.
- **Viren Pro-Tip**: OpenAI uses a tiered rate-limit system. If you are a new user, you must pre-pay at least $5 to reach Tier 1 and avoid aggressive rate limiting.

### Anthropic (The Logic Specialist)
Anthropic's Claude 3.5 models are widely considered the gold standard for software engineering. They possess a nuanced understanding of code safety and maintain coherence over massive contexts.
- **Acquisition Link**: [https://console.anthropic.com/settings/keys](https://console.anthropic.com/)
- **Required Env Variable**: `ANTHROPIC_API_KEY`
- **Recommended Models**:
    - `claude-3-5-sonnet`: Currently the highest-rated model for technical tasks and large-scale project refactors.
    - `claude-3-opus`: For high-end creative work and nuanced philosophical inquiries.
- **Viren Pro-Tip**: Claude is exceptionally effective when used with Viren's `!d` (Codedump) because its context window is stable up to 200k tokens.

### DeepSeek (Efficiency & Price Leader)
DeepSeek has rapidly become a developer favorite due to its extreme performance in coding tasks and its highly aggressive pricing (often 1/10th the cost of GPT-4o).
- **Acquisition Link**: [https://platform.deepseek.com/api_keys](https://api-docs.deepseek.com/)
- **Required Env Variable**: `DEEP_SEEK_API_KEY`
- **Recommended Models**:
    - `deepseek-chat` (V3): A general-purpose powerhouse that rivals GPT-4o.
    - `deepseek-coder`: Specifically tuned for the most difficult programming challenges.
- **Viren Pro-Tip**: DeepSeek's API is 100% compatible with the OpenAI format, making it one of the most stable connectors in Viren.

---

## 3. Infrastructure, Speed & Search Platforms

### Brave Search API (The Live Internet Bridge)
Standard LLMs are limited by their training cutoff dates (e.g., they might not know about Go 1.24). Viren solves this by using Brave Search to provide "Live Context."
- **Acquisition Link**: [https://api.search.brave.com/app/dashboard](https://brave.com/search/api/)
- **Required Env Variable**: `BRAVE_API_KEY`
- **Use Case**: Enabling the `!w` (Web Search) command. This allows Viren to find library updates or documentation released *today*.

### Groq (Extreme Inference Speed)
Groq uses specialized LPU (Language Processing Unit) hardware to run open-source models (like Llama and Mixtral) at over 300 tokens per second.
- **Acquisition Link**: [https://console.groq.com/keys](https://console.groq.com/keys)
- **Required Env Variable**: `GROQ_API_KEY`
- **Recommended Models**:
    - `llama-3.1-70b-versatile`: Deep reasoning at extreme speed.
    - `mixtral-8x7b-32768`: Excellent for fast summaries.
- **Viren Pro-Tip**: The experience of using Viren with Groq is near-instant. It is ideal for rapid-fire Q&A where latency is the enemy.

---

## 4. Enterprise & Aggregator Platforms

### Amazon Bedrock
Enterprise-grade access for users already within the AWS ecosystem.
- **Acquisition**: [https://aws.amazon.com/bedrock](https://aws.amazon.com/bedrock)
- **Required Env Variables**: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.
- **Developer Note**: Requires AWS IAM permissions to be configured correctly in your AWS Console.

---

## 5. Setting Environment Variables: OS-Specific Protocols

### Linux & macOS (Bash/Zsh)
1.  Open your profile: `nano ~/.zshrc` (or `~/.bashrc`).
2.  Add the following block:
    ```bash
    # Viren API Configuration
    export OPENAI_API_KEY="sk-..."
    export ANTHROPIC_API_KEY="sk-ant-..."
    export DEEP_SEEK_API_KEY="sk-..."
    export BRAVE_API_KEY="bs-..."
    ```
3.  Save and exit (`Ctrl+O`, `Enter`, `Ctrl+X`).
4.  Reload the shell: `source ~/.zshrc`.

### Windows (PowerShell - Recommended)
Run these commands in an administrative PowerShell window:
```powershell
[System.Environment]::SetEnvironmentVariable('OPENAI_API_KEY', 'your-key-here', [System.EnvironmentVariableTarget]::User)
```

---

## 6. Billing & Rate Limit Strategies

Every AI provider has different billing mechanics. Understanding these can save you hundreds of dollars in API costs.

### Pre-paid vs. Post-paid
- **OpenAI/Anthropic**: Usually requires a pre-paid balance. If your balance hits zero, Viren will return a `401 Unauthorized` or `402 Payment Required` error.
- **Google Gemini**: Often offers a "Pay-as-you-go" model with a free tier.

### Managing Token Usage
Viren is context-aware. Use the `!c` (Clear) command frequently. If you have a conversation history of 50 messages, every new message sends those 50 messages back to the API, multiplying your token cost.

---

## 7. Model Comparison Metrics

| Model | Logic Score | Speed (tokens/sec) | Cost |
| :--- | :--- | :--- | :--- |
| `gpt-4o` | 9.5/10 | 60-80 | High |
| `claude-3.5-sonnet` | 9.8/10 | 50-70 | Medium |
| `llama-3.1-70b` | 8.5/10 | 300+ (on Groq) | Low |
| `deepseek-chat` | 9.2/10 | 40-60 | Ultra-Low |

---

## 8. Regional Data Residency Policies

When using Viren in high-security environments, you must be aware of where your data is being processed.
- **Anthropic**: Offers specific endpoints for EU-resident data.
- **Amazon Bedrock**: Allows you to pin your model requests to a specific AWS region (e.g., `us-east-1` or `eu-central-1`) to ensure compliance with local laws.

---

## 9. Advanced IAM Role Configuration (AWS)

If you are using Amazon Bedrock, you should not use your root account keys. Create an IAM user with the following policy:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "bedrock:InvokeModel",
            "Resource": "*"
        }
    ]
}
```

---

## 10. Troubleshooting Authentication

- **"Variable not exported"**: Run `env | grep API_KEY`. If nothing shows up, your shell profile was not sourced correctly.
- **"Invalid API Key"**: Check for trailing spaces or special characters when you copy-pasted the key.

**Your API keys are the keys to your digital workshop. Protect them.**
