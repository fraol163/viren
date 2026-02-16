# Real-World Example: Smart Code Export

Viren eliminates the tedious "copy-paste" cycle. This guide shows you how to turn an AI conversation into a production-ready file in seconds.

## Scenario
You are building a Go microservice and need a boilerplate for a middleware that handles JWT authentication.

### Step 1: Generation
Ask Viren for the code:
`USER ❯ Write a Go middleware for JWT authentication using the golang-jwt library.`

### Step 2: Detection
The AI generates a response containing a code block like this:
```go
package middleware
import "github.com/golang-jwt/jwt/v5"
// ... logic here
```

### Step 3: Triggering Export
You have two options:
1. **Manual**: Type `!e` after the message is finished.
2. **Auto**: If configured, Viren will detect the code and ask: *"Detected code. Export? (y/N)"*

### Step 4: Intelligent Filenaming
Viren doesn't just name every file `code.txt`. It scans the code's contents:
- It sees `package middleware` and `func Auth(...)`.
- The `fzf` menu will appear with suggestions like:
    - `auth_middleware.go`
    - `jwt_auth.go`
    - `middleware.go`
    - `[Custom Name]`

### Step 5: Finalization
Select `auth_middleware.go` and hit `Enter`.
Viren saves the file to your current working directory and adds it to your "Recently Created Files" list.

---

## Advanced Exporting
- **Multi-Block Export**: If the AI provides three files (e.g., `main.go`, `types.go`, `util.go`), Viren will cycle through each one, letting you pick a name for each.
- **TURN Export**: Use `!e` and select "Turn Export" to save the entire conversation (User prompts + AI responses) as a single Markdown file for documentation.
- **Block Selection**: Use `fzf` to pick only the *specific* functions or blocks you want to keep, discarding the rest.
