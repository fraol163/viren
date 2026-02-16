# Real-World Scenario: Solving Multi-File Logic Bugs

This is the definitive guide to using Viren's Context Ingestion engine to solve complex, cross-file bugs that standard "single-file" AI tools cannot see.

---

## 1. The Scenario: The "Invisible" Crash
You have a Python project with two files:
1.  `database.py`: Handles connections and queries.
2.  `app.py`: The entry point that calls the database logic.

Your app crashes with a `AttributeError: 'NoneType' object has no attribute 'cursor'` only when the database is under high load. You suspect a race condition, but looking at just one file doesn't reveal the cause.

---

## 2. The Solution: Multi-File Context Loading

### Step 1: The "Load" Protocol (`!l`)
Inside the Viren interactive shell:
1.  Type `!l`.
2.  Use the `fzf` menu to navigate to `database.py`.
3.  Press `Tab` to select it (it will be highlighted).
4.  Navigate to `app.py`.
5.  Press `Tab` to select it.
6.  Hit `Enter`.

**Internal Process**: Viren reads both files, identifies their relative paths, and injects them into the "System" context of the next prompt.

### Step 2: The Inquiry
`USER ❯ I'm getting a NoneType error on the cursor. Looking at how app.py initializes database.py, why is the connection failing under load?`

### Step 3: The Cross-File Discovery
Because the AI can see **both** files simultaneously, it can reason about the interaction:
- *"I see that in app.py, you are not checking the return value of `db.connect()`. In database.py, line 45, if the pool is full, it returns `None` instead of raising an exception."*

---

## 3. Project-Wide Logic: The Codedump (`!d`)

If your bug involves 20 files instead of 2, use `!d`.

1.  Run `!d`.
2.  Viren scans the entire directory.
3.  **Exclusion**: It will show you a list of every file. Use the search bar to find `node_modules`, `.git`, or `dist` and exclude them to save tokens.
4.  **Bundling**: Viren creates a project map.
5.  **Query**: `USER ❯ Find all locations where database connections are not properly closed.`

---

## 4. Native Support for Rich Documentation

Sometimes the answer isn't in your code, but in the **Documentation**.
- **PDF Datasheets**: Use `!l` to load a 200-page PDF manual for a hardware component. Viren parses the text, and you can ask: *"What is the register address for the SPI clock polarity?"*
- **Excel Metrics**: Load an `.xlsx` file of performance metrics. Ask: *"Based on these numbers, which endpoint is the primary bottleneck?"*

---

## 5. Pro-Tips for Large Contexts
- **Shallow Loads**: If you are in a massive repository, update your `config.json` with `shallow_load_dirs` to prevent Viren from scanning unnecessary depths.
- **Context Hygiene**: Use `!c` (Clear) once you have fixed the bug. Large contexts consume more tokens per prompt and can eventually slow down the model's reasoning.

**Viren provides the "God's Eye View" of your technical stack.**
