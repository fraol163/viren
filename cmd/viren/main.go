package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/chzyer/readline"
	"github.com/fraol163/viren/internal/chat"
	"github.com/fraol163/viren/internal/config"
	"github.com/fraol163/viren/internal/platform"
	"github.com/fraol163/viren/internal/ui"
	"github.com/fraol163/viren/internal/util"
	"github.com/fraol163/viren/pkg/types"
	"github.com/google/uuid"
	"github.com/tiktoken-go/tokenizer"
)

var (

	version	= "v1.0.0"




	buildTime = "unknown"
	gitCommit = "unknown"
)

func init() {

	flag.Usage = func() {
		state := config.InitializeAppState()
		terminal := ui.NewTerminal(state.Config)
		terminal.ShowHelp()
		os.Exit(0)
	}
}

func main() {

	ui.EnableVirtualTerminalProcessing()

	state := config.InitializeAppState()

	stdoutStat, _ := os.Stdout.Stat()
	if (stdoutStat.Mode() & os.ModeCharDevice) == 0 {
		state.Config.IsPipedOutput = true
	}

	terminal := ui.NewTerminal(state.Config)
	terminal.ApplyTheme()
	chatManager := chat.NewManager(state)
	platformManager := platform.NewManager(state.Config)

	state.CurrentPersonality = state.Config.CurrentPersonality
	state.CurrentMode = state.Config.CurrentMode

	chatManager.UpdateFullSystemPrompt()

	// parse command line arguments
	var (
		helpFlag	= flag.Bool("h", false, "Show help")
		codedumpFlag	= flag.String("d", "", "Generate codedump file (optionally specify directory path)")
		platformFlag	= flag.String("p", "", "Switch platform (leave empty for interactive selection)")
		modelFlag	= flag.String("m", "", "Specify model to use")
		allModelsFlag	= flag.String("o", "", "Specify platform and model (format: platform|model)")
		exportCodeFlag	= flag.Bool("e", false, "Export code blocks from the last response")
		tokenFlag	= flag.String("t", "", "Estimate token count in file")
		loadFileFlag	= flag.String("l", "", "Load and display file content (supports text, PDF, DOCX, XLSX, CSV)")
		webSearchFlag	= flag.String("w", "", "Perform a web search and print the results")
		scrapeURLFlag	= flag.String("s", "", "Scrape a URL and print the content")
		continueFlag	= flag.Bool("c", false, "Continue from latest session")
		clearFlag	= flag.Bool("clear", false, "Clear latest session")
		historyFlag	= flag.Bool("a", false, "Search and load previous sessions")
		versionFlag	= flag.Bool("version", false, "Show version")
		vFlag		= flag.Bool("v", false, "Show version")
	)
	flag.StringVar(tokenFlag, "token", "", "Estimate token count in file")
	flag.BoolVar(continueFlag, "continue", false, "Continue from latest session")
	flag.BoolVar(historyFlag, "history", false, "Search and load previous sessions")
	flag.BoolVar(historyFlag, "hs", false, "Search and load previous sessions")

	noHistoryFlag := flag.Bool("nh", false, "Disable session saving for this run")
	flag.Bool("no-history", false, "Disable session saving for this run")

	flag.Parse()

	if *versionFlag || *vFlag {
		fmt.Printf("Viren %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		fmt.Printf("Git Commit: %s\n", gitCommit)
		return
	}

	if *helpFlag {
		terminal.ShowHelp()
		return
	}

	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".viren", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) && !state.Config.IsPipedOutput && terminal.IsTerminal() {
		err = config.RunOnboarding(terminal, state.Config)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("onboarding failed: %v", err))
		} else {

			state = config.InitializeAppState()
			chatManager = chat.NewManager(state)
			platformManager = platform.NewManager(state.Config)
			terminal = ui.NewTerminal(state.Config)

			state.CurrentPersonality = state.Config.CurrentPersonality
			state.CurrentMode = state.Config.CurrentMode

			chatManager.UpdateFullSystemPrompt()

			terminal.ClearTerminal()
			terminal.ApplyTheme()
		}
	}

	if flag.Lookup("no-history").Value.String() == "true" {
		*noHistoryFlag = true
	}
	remainingArgs := flag.Args()

	// Check if input is being piped
	var pipedInput string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {

		pipedBytes, err := io.ReadAll(os.Stdin)
		if err == nil && len(pipedBytes) > 0 {
			pipedInput = string(pipedBytes)
		}
	}

	if *clearFlag {
		if !state.Config.EnableSessionSave {
			terminal.PrintError("session save feature is disabled in config")
			return
		}

		fmt.Printf("\033[91mdelete all temp files? (y/N)\033[0m ")
		var response string
		_, err := fmt.Scanln(&response)

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("cancelled")
			return
		}

		tmpDir, err := util.GetTempDir()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("failed to get temp directory: %v", err))
			return
		}

		err = os.RemoveAll(tmpDir)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error clearing temporary files: %v", err))
			return
		}

		err = os.MkdirAll(tmpDir, 0755)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error recreating temporary directory: %v", err))
			return
		}

		return
	}

	if *historyFlag {

		if !state.Config.SaveAllSessions {
			terminal.PrintError("history search requires save_all_sessions to be enabled in config")
			return
		}

		exact := len(remainingArgs) > 0 && remainingArgs[0] == "exact"

		session, err := chatManager.ManageSessions(terminal, exact)
		if err != nil {
			if err.Error() == "selection cancelled" {
				return
			}
			terminal.PrintError(fmt.Sprintf("%v", err))
			return
		}

		chatManager.RestoreSessionState(session)

		err = platformManager.Initialize()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("failed to initialize client: %v", err))
			return
		}

		fmt.Printf("\033[91mrestored session from %s UTC\033[0m\n", time.Unix(session.Timestamp, 0).UTC().Format("2006-01-02 15:04:05"))

		for _, entry := range session.ChatHistory {
			if entry.User == state.Config.SystemPrompt {
				continue
			}

			if entry.User != "" {
				theme := terminal.GetTheme()
				fmt.Printf("%s USER \033[0m ❯ %s\n", theme.UserBox, entry.User)
			}

			if entry.Bot != "" {
				theme := terminal.GetTheme()
				fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, entry.Bot)
			}
		}

		return
	}

	if flag.Lookup("d").Value.String() != flag.Lookup("d").DefValue {
		targetDir := *codedumpFlag
		if targetDir == "" {
			targetDir = "."
		}

		if !isValidCodedumpDir(targetDir) {
			if targetDir != "." {
				terminal.PrintError("invalid directory path or permission denied")
				return
			}
		}

		codedump, err := terminal.CodeDumpFromDirForCLI(targetDir)
		if err != nil {

			if strings.Contains(err.Error(), "user cancelled") {
				return
			}
			terminal.PrintError(fmt.Sprintf("error generating codedump: %v", err))
			return
		}

		currentDir, err := os.Getwd()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error getting current directory: %v", err))
			return
		}
		filename := generateUniqueCodeDumpFilename(currentDir, codedump)
		err = os.WriteFile(filename, []byte(codedump), 0644)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error writing codedump file: %v", err))
			return
		}

		fmt.Println(filename)
		return
	}

	if *exportCodeFlag {
		err := handleExportCodeBlocks(chatManager, terminal)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error exporting code blocks: %v", err))
		}
		return
	}

	if *tokenFlag != "" {
		err := handleTokenCount(*tokenFlag, *modelFlag, terminal, state)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error counting tokens: %v", err))
		}
		return
	}

	if *allModelsFlag != "" {
		parts := strings.Split(*allModelsFlag, "|")
		if len(parts) != 2 {
			terminal.PrintError("invalid -o format: use platform|model (e.g., openai|gpt-4)")
			return
		}

		platformName := strings.TrimSpace(parts[0])
		modelName := strings.TrimSpace(parts[1])

		if platformName == "" || modelName == "" {
			terminal.PrintError("invalid -o format: platform and model cannot be empty")
			return
		}

		if platformName != "openai" {
			if _, exists := state.Config.Platforms[platformName]; !exists {
				terminal.PrintError(fmt.Sprintf("platform '%s' not found", platformName))
				return
			}
		}

		*platformFlag = platformName
		*modelFlag = modelName
	}

	finalPlatform := state.Config.CurrentPlatform
	finalModel := state.Config.CurrentModel

	if p := os.Getenv("VIREN_DEFAULT_PLATFORM"); p != "" {
		finalPlatform = p
	}
	if m := os.Getenv("VIREN_DEFAULT_MODEL"); m != "" {
		finalModel = m
	}

	if *platformFlag != "" {
		finalPlatform = *platformFlag
	}
	if *modelFlag != "" {
		finalModel = *modelFlag
	}

	sessionRestored := false
	if *continueFlag {

		if !state.Config.EnableSessionSave {
			terminal.PrintError("session save feature is disabled in config")
			return
		}

		var session *types.SessionFile
		var err error

		if len(remainingArgs) > 0 {
			customPath := remainingArgs[0]

			if _, statErr := os.Stat(customPath); statErr == nil {

				session, err = chatManager.LoadCustomHistoryFile(customPath)
				if err != nil {
					terminal.PrintError(fmt.Sprintf("%v", err))
					return
				}

				remainingArgs = remainingArgs[1:]
			} else {

				session, err = chatManager.LoadLatestSessionState()
				if err != nil {

					if !strings.Contains(err.Error(), "no session file found") {
						terminal.PrintError(fmt.Sprintf("error loading session: %v", err))
						return
					}
					terminal.PrintError("no previous session found to continue from")
					return
				}
			}
		} else {

			session, err = chatManager.LoadLatestSessionState()
			if err != nil {

				if !strings.Contains(err.Error(), "no session file found") {
					terminal.PrintError(fmt.Sprintf("error loading session: %v", err))
					return
				}
				terminal.PrintError("no previous session found to continue from")
				return
			}
		}

		chatManager.RestoreSessionState(session)
		sessionRestored = true

		finalPlatform = session.Platform
		finalModel = session.Model

		fmt.Printf("\033[91mrestored session from %s UTC\033[0m\n", time.Unix(session.Timestamp, 0).UTC().Format("2006-01-02 15:04:05"))

		for _, entry := range session.ChatHistory {
			if entry.User == state.Config.SystemPrompt {
				continue
			}

			if entry.User != "" {
				theme := terminal.GetTheme()
				fmt.Printf("%s USER \033[0m ❯ %s\n", theme.UserBox, entry.User)
			}

			if entry.Bot != "" {
				theme := terminal.GetTheme()
				fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, entry.Bot)
			}
		}
	}

	if !sessionRestored && (finalPlatform != state.Config.CurrentPlatform || finalModel != state.Config.CurrentModel) {

		if *platformFlag != "" {
			result, err := platformManager.SelectPlatform(finalPlatform, finalModel, terminal.FzfSelect)
			if err != nil {
				terminal.PrintError(fmt.Sprintf("%v", err))
				return
			}
			if result != nil {
				chatManager.SetCurrentPlatform(result["platform_name"].(string))
				chatManager.SetCurrentModel(result["picked_model"].(string))
				state.Config.CurrentPlatform = result["platform_name"].(string)
				state.Config.CurrentModel = result["picked_model"].(string)
			}
		} else {
			chatManager.SetCurrentPlatform(finalPlatform)
			chatManager.SetCurrentModel(finalModel)
			state.Config.CurrentPlatform = finalPlatform
			state.Config.CurrentModel = finalModel
		}

		config.SaveConfigToFile(state.Config)
	}

	err := platformManager.Initialize()
	if err != nil {
		terminal.PrintError(fmt.Sprintf("failed to initialize client: %v", err))
		return
	}

	if *webSearchFlag != "" {
		queries := splitByDelimiters(*webSearchFlag)
		prompt := strings.Join(flag.Args(), " ")

		// Combine results from multiple queries
		var allResults []string
		for _, query := range queries {
			results, err := terminal.WebSearch(query)
			if err != nil {
				terminal.PrintError(fmt.Sprintf("error during web search for '%s': %v", query, err))
				continue
			}
			allResults = append(allResults, results)
		}

		combinedResults := strings.Join(allResults, "\n\n---\n\n")

		if prompt == "" {
			fmt.Print(combinedResults)
			return
		}

		err := handleFlagWithPrompt(chatManager, platformManager, terminal, state, combinedResults, prompt, *noHistoryFlag)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error: %v", err))
		}
		return
	}

	if *scrapeURLFlag != "" {
		urls := splitByDelimiters(*scrapeURLFlag)
		prompt := strings.Join(flag.Args(), " ")

		// Combine content from multiple URLs
		var allContent []string
		for _, url := range urls {
			content, err := terminal.ScrapeURLs([]string{url})
			if err != nil {
				terminal.PrintError(fmt.Sprintf("error scraping URL '%s': %v", url, err))
				continue
			}
			allContent = append(allContent, content)
		}

		combinedContent := strings.Join(allContent, "\n\n---\n\n")

		if prompt == "" {
			fmt.Println(strings.TrimSpace(combinedContent))
			return
		}

		err := handleFlagWithPrompt(chatManager, platformManager, terminal, state, combinedContent, prompt, *noHistoryFlag)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error: %v", err))
		}
		return
	}

	if *loadFileFlag != "" {
		files := splitByDelimiters(*loadFileFlag)
		prompt := strings.Join(flag.Args(), " ")

		// Load content from all specified files
		var allContent []string
		for _, file := range files {

			if terminal.IsURL(file) {
				content, err := terminal.LoadFileContent([]string{file})
				if err != nil {
					terminal.PrintError(fmt.Sprintf("error loading URL '%s': %v", file, err))
					continue
				}
				allContent = append(allContent, content)
			} else {

				if _, err := os.Stat(file); os.IsNotExist(err) {
					terminal.PrintError(fmt.Sprintf("file does not exist: %s", file))
					continue
				}

				content, err := terminal.LoadFileContent([]string{file})
				if err != nil {
					terminal.PrintError(fmt.Sprintf("error loading file '%s': %v", file, err))
					continue
				}
				allContent = append(allContent, content)
			}
		}

		combinedContent := strings.Join(allContent, "\n\n---\n\n")

		if prompt == "" {
			fmt.Println(strings.TrimSpace(combinedContent))
			return
		}

		err := handleFlagWithPrompt(chatManager, platformManager, terminal, state, combinedContent, prompt, *noHistoryFlag)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error: %v", err))
		}
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigChan {
			if state.IsStreaming && state.StreamingCancel != nil {
				fmt.Print("\r\033[K")
				state.StreamingCancel()
			} else if state.IsExecutingCommand && state.CommandCancel != nil {
				fmt.Print("\r\033[K")
				state.CommandCancel()
			} else {
				os.Exit(0)
			}
		}
	}()

	if len(remainingArgs) > 0 || pipedInput != "" {
		var query string

		if pipedInput != "" && len(remainingArgs) > 0 {

			query = strings.TrimSpace(pipedInput) + " " + strings.Join(remainingArgs, " ")
		} else if pipedInput != "" {

			query = strings.TrimSpace(pipedInput)
		} else {

			query = strings.Join(remainingArgs, " ")
		}

		err := processDirectQuery(query, chatManager, platformManager, terminal, state, *exportCodeFlag, *noHistoryFlag)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("%v", err))
		}
		return
	}

	terminal.ApplyTheme()
	terminal.ShowLogo()
	runInteractiveMode(chatManager, platformManager, terminal, state, *noHistoryFlag)
}

func processDirectQuery(query string, chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, exportCode bool, noHistory bool) error {
	if handleSpecialCommands(query, chatManager, platformManager, terminal, state, noHistory, nil) {
		return nil
	}

	chatManager.AddUserMessage(query)

	// Pass animation context to SendChatRequest if needed
	var animationCancel context.CancelFunc
	if !state.Config.IsPipedOutput {
		var ctx context.Context
		ctx, animationCancel = context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "thinking")
	}

	response, err := platformManager.SendChatRequest(chatManager.GetMessages(), chatManager.GetCurrentModel(), &state.StreamingCancel, &state.IsStreaming, animationCancel, terminal)
	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			if animationCancel != nil {
				animationCancel()
			}
			return nil
		}
		return err
	}

	if animationCancel != nil {
		animationCancel()
	}

	if platformManager.IsReasoningModel(chatManager.GetCurrentModel()) {
		if state.Config.IsPipedOutput {
			fmt.Printf("%s\n", response)
		} else {
			theme := terminal.GetTheme()
			fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)
		}
	} else {

		if !state.Config.IsPipedOutput {
			fmt.Println()
		}
	}

	chatManager.AddAssistantMessage(response)
	chatManager.AddToHistory(query, response)

	if !state.Config.IsPipedOutput && terminal.IsTerminal() {
		shellRegex := regexp.MustCompile("(?s)```(?:bash|sh|shell)\\n(.*?)\\n```")
		matches := shellRegex.FindAllStringSubmatch(response, -1)
		if len(matches) > 0 {

			lastMatch := matches[len(matches)-1]
			command := strings.TrimSpace(lastMatch[1])

			fmt.Printf("\n\033[1;36mDETECTED COMMAND:\033[0m \033[93m%s\033[0m\n", command)
			fmt.Print("\033[1;36mEXECUTE? (y/N) ❯ \033[0m")

			var confirm string
			reader := bufio.NewReader(os.Stdin)
			confirm, _ = reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))

			if confirm == "y" || confirm == "yes" {
				handleShellCommand(command, chatManager, terminal, state)
			}
		}
	}

	if state.Config.EnableSessionSave && !noHistory {
		if err := chatManager.SaveSessionState(); err != nil {
			terminal.PrintError(fmt.Sprintf("warning: failed to save session: %v", err))
		}
	}

	if exportCode {
		filePaths, exportErr := chatManager.ExportCodeBlocks(terminal)
		if exportErr != nil {
			terminal.PrintError(fmt.Sprintf("error exporting code blocks: %v", exportErr))
		} else if len(filePaths) > 0 {
			for _, filePath := range filePaths {
				fmt.Println(filePath)
			}
		}
	}

	return nil
}

// runInteractiveMode runs the main interactive chat loop
func runInteractiveMode(chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, noHistory bool) {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          terminal.GetPrompt(),
		InterruptPrompt: "",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {

				if !state.Config.MuteNotifications {
					fmt.Printf("\033[93mpress ctrl+d to exit\033[0m\n")
				}
				continue
			}

			break
		}
		input := strings.TrimSpace(line)

		if input == "" {
			continue
		}

		if strings.HasSuffix(input, state.Config.MultiLine) && input != state.Config.MultiLine {

			input = strings.TrimSuffix(input, state.Config.MultiLine)
			input = strings.TrimRight(input, " \t")

			var lines []string
			if input != "" {
				lines = append(lines, input)
			}

			multiLineRl, err := readline.NewEx(&readline.Config{
				Prompt:      "... ",
				HistoryFile: "/dev/null",
			})
			if err != nil {
				terminal.PrintError(fmt.Sprintf("error creating multi-line input: %v", err))
				continue
			}

			multiLineActive := true
			for multiLineActive {
				line, err := multiLineRl.Readline()
				if err != nil {
					if err == readline.ErrInterrupt || err == io.EOF {
						multiLineRl.Close()
						multiLineActive = false
						break
					}
					break
				}

				if strings.HasSuffix(line, state.Config.MultiLine) {

					line = strings.TrimSuffix(line, state.Config.MultiLine)
					line = strings.TrimRight(line, " \t")
					lines = append(lines, line)
				} else {

					lines = append(lines, line)
					multiLineActive = false
					break
				}
			}

			multiLineRl.Close()

			input = strings.Join(lines, "\n")
			if strings.TrimSpace(input) == "" {
				continue
			}
		}

		if handleSpecialCommands(input, chatManager, platformManager, terminal, state, noHistory, rl) {
			continue
		}

		chatManager.AddUserMessage(input)

		ctx, animationCancel := context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "waiting")

		response, err := platformManager.SendChatRequest(chatManager.GetMessages(), chatManager.GetCurrentModel(), &state.StreamingCancel, &state.IsStreaming, animationCancel, terminal)

		animationCancel()

		if err != nil {
			if err.Error() == "request was interrupted" {
				chatManager.RemoveLastUserMessage()
				continue
			}
			terminal.PrintError(fmt.Sprintf("%v", err))
			continue
		}

		if platformManager.IsReasoningModel(chatManager.GetCurrentModel()) {
			if state.Config.IsPipedOutput {
				fmt.Printf("%s\n", response)
			} else {
				theme := terminal.GetTheme()
				fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)
			}
		} else {

			if !state.Config.IsPipedOutput {
				fmt.Println()
			}
		}

		chatManager.AddAssistantMessage(response)
		chatManager.AddToHistory(input, response)

		if !state.Config.IsPipedOutput {
			fmt.Printf("\033[38;2;0;0;0m%s\033[0m\n", strings.Repeat("┈", 60))
		}

		shellRegex := regexp.MustCompile("(?s)```(?:bash|sh|shell)\\n(.*?)\\n```")
		matches := shellRegex.FindAllStringSubmatch(response, -1)
		if len(matches) > 0 {

			lastMatch := matches[len(matches)-1]
			command := strings.TrimSpace(lastMatch[1])

			fmt.Printf("\n\033[1;36mDETECTED COMMAND:\033[0m \033[93m%s\033[0m\n", command)
			fmt.Print("\033[1;36mEXECUTE? (y/N) ❯ \033[0m")

			var confirm string

			reader := bufio.NewReader(os.Stdin)
			confirm, _ = reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))

			if confirm == "y" || confirm == "yes" {
				handleShellCommand(command, chatManager, terminal, state)
			}
		}

		if state.Config.EnableSessionSave && !noHistory {
			if err := chatManager.SaveSessionState(); err != nil {
				terminal.PrintError(fmt.Sprintf("warning: failed to save session: %v", err))
			}
		}
	}
}

func handleSpecialCommands(input string, chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, noHistory bool, rl *readline.Instance) bool {
	return handleSpecialCommandsInternal(input, chatManager, platformManager, terminal, state, false, noHistory, rl)
}

func handleSpecialCommandsInternal(input string, chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, fromHelp bool, noHistory bool, rl *readline.Instance) bool {
	configObj := state.Config

	switch {
	case input == configObj.ExitKey:
		os.Exit(0)
		return true

	case input == configObj.HelpKey || input == "help":
		selectedCommand := terminal.ShowHelpFzf()
		if selectedCommand == ">state" {
			err := handleShowState(chatManager, terminal, state)
			if err != nil {
				terminal.PrintError(fmt.Sprintf("error showing state: %v", err))
			}
			return true
		}
		if selectedCommand != "" {

			return handleSpecialCommandsInternal(selectedCommand, chatManager, platformManager, terminal, state, true, noHistory, rl)
		}
		return true

	case input == configObj.ClearHistory:
		terminal.ClearTerminal()
		chatManager.ClearHistory()
		terminal.ApplyTheme()
		terminal.ShowLogo()
		terminal.PrintInfo("history and screen cleared")
		return true

	case input == configObj.ModelSwitch:
		models, err := platformManager.ListModels()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error fetching models: %v", err))
			return true
		}

		selectedModel, err := terminal.FzfSelect(models, "model: ")
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error selecting model: %v", err))
			return true
		}

		if selectedModel != "" {
			chatManager.SetCurrentModel(selectedModel)
			state.Config.CurrentModel = selectedModel
			config.SaveConfigToFile(state.Config)
			if !configObj.MuteNotifications {
				terminal.PrintModelSwitch(selectedModel)
			}
		}
		return true

	case strings.HasPrefix(input, configObj.ModelSwitch+" "):
		modelName := strings.TrimPrefix(input, configObj.ModelSwitch+" ")
		chatManager.SetCurrentModel(modelName)
		state.Config.CurrentModel = modelName
		config.SaveConfigToFile(state.Config)
		if !configObj.MuteNotifications {
			terminal.PrintModelSwitch(modelName)
		}
		return true

	case input == configObj.PlatformSwitch:
		result, err := platformManager.SelectPlatform("", "", terminal.FzfSelect)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("%v", err))
		} else if result != nil {
			chatManager.SetCurrentPlatform(result["platform_name"].(string))
			chatManager.SetCurrentModel(result["picked_model"].(string))
			state.Config.CurrentPlatform = result["platform_name"].(string)
			state.Config.CurrentModel = result["picked_model"].(string)
			configObj.CurrentBaseURL = result["base_url"].(string)
			state.Config.CurrentBaseURL = result["base_url"].(string)
			config.SaveConfigToFile(state.Config)
			err = platformManager.Initialize()
			if err != nil {
				terminal.PrintError(fmt.Sprintf("error initializing client: %v", err))
			} else {
				if !configObj.MuteNotifications {
					terminal.PrintPlatformSwitch(result["platform_name"].(string), result["picked_model"].(string))
				}
			}
		}
		return true

	case strings.HasPrefix(input, configObj.PlatformSwitch+" "):
		platformName := strings.TrimPrefix(input, configObj.PlatformSwitch+" ")
		result, err := platformManager.SelectPlatform(platformName, "", terminal.FzfSelect)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("%v", err))
		} else if result != nil {
			chatManager.SetCurrentPlatform(result["platform_name"].(string))
			chatManager.SetCurrentModel(result["picked_model"].(string))
			state.Config.CurrentPlatform = result["platform_name"].(string)
			state.Config.CurrentModel = result["picked_model"].(string)
			configObj.CurrentBaseURL = result["base_url"].(string)
			state.Config.CurrentBaseURL = result["base_url"].(string)
			config.SaveConfigToFile(state.Config)
			err = platformManager.Initialize()
			if err != nil {
				terminal.PrintError(fmt.Sprintf("error initializing client: %v", err))
			} else {
				if !configObj.MuteNotifications {
					terminal.PrintPlatformSwitch(result["platform_name"].(string), result["picked_model"].(string))
				}
			}
		}
		return true

	case input == configObj.AllModels:
		return handleAllModels(chatManager, platformManager, terminal, state)

	case input == "!v":
		modes := chat.GetModes()
		var modeNames []string
		for _, m := range modes {
			modeNames = append(modeNames, fmt.Sprintf("%s (%s)", m.Name, m.ID))
		}

		selected, err := terminal.FzfSelect(modeNames, "mode: ")
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error selecting mode: %v", err))
			return true
		}

		if selected != "" {

			parts := strings.Split(selected, "(")
			if len(parts) < 2 {
				return true
			}
			modeID := strings.TrimSuffix(parts[len(parts)-1], ")")
			chatManager.SetMode(modeID)
			state.Config.CurrentMode = modeID
			config.SaveConfigToFile(state.Config)

			if rl != nil {
				rl.SetPrompt(terminal.GetPrompt())
			}

			terminal.PrintInfo(fmt.Sprintf("switched to %s", selected))
		}
		return true

	case input == "!z":
		themes := ui.GetThemes()
		var themeNames []string
		for _, t := range themes {
			themeNames = append(themeNames, fmt.Sprintf("%s (%s)", t.Name, t.ID))
		}

		selected, err := terminal.FzfSelect(themeNames, "theme: ")
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error selecting theme: %v", err))
			return true
		}

		if selected != "" {

			parts := strings.Split(selected, "(")
			if len(parts) < 2 {
				return true
			}
			themeID := strings.TrimSuffix(parts[len(parts)-1], ")")
			terminal.SetTheme(themeID)
			state.Config.CurrentTheme = themeID
			state.Config.UserProfile.Theme = themeID

			err = config.SaveConfigToFile(state.Config)
			if err != nil {
				terminal.PrintError(fmt.Sprintf("failed to save theme config: %v", err))
			}

			terminal.ApplyTheme()
			terminal.ClearTerminal()
			terminal.ApplyTheme()
			terminal.ShowLogo()

			if rl != nil {
				rl.SetPrompt(terminal.GetPrompt())
			}

			newTheme := terminal.GetTheme()
			fmt.Printf("\033[92m  ○\033[0m \033[1mTHEME UPDATED\033[0m  %s[\033[0m\033[1;96m%s\033[0m%s]\033[0m\n", newTheme.BorderColor, newTheme.Name, newTheme.BorderColor)
		}
		return true

	case input == "!onboard":
		err := config.RunOnboarding(terminal, state.Config)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("onboarding failed: %v", err))
		} else {

			newCfg := config.DefaultConfig()
			*state.Config = *newCfg

			state.CurrentPersonality = state.Config.CurrentPersonality
			state.CurrentMode = state.Config.CurrentMode
			state.CurrentTheme = state.Config.CurrentTheme

			chatManager.UpdateFullSystemPrompt()

			platformManager.Initialize()

			terminal.ClearTerminal()
			terminal.ApplyTheme()
			terminal.ShowLogo()
		}

		if rl != nil {
			rl.SetPrompt(terminal.GetPrompt())
		}
		return true

	case input == "!u":
		personalities := chat.GetPersonalities()
		var pNames []string
		for _, p := range personalities {
			pNames = append(pNames, fmt.Sprintf("%s (%s)", p.Name, p.ID))
		}

		selected, err := terminal.FzfSelect(pNames, "personality: ")
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error selecting personality: %v", err))
			return true
		}

		if selected != "" {
			parts := strings.Split(selected, "(")
			if len(parts) < 2 {
				return true
			}
			pID := strings.TrimSuffix(parts[len(parts)-1], ")")
			chatManager.SetPersonality(pID)

			err = config.SaveConfigToFile(state.Config)
			if err != nil {
				terminal.PrintError(fmt.Sprintf("failed to save personality config: %v", err))
			}

			terminal.PrintInfo(fmt.Sprintf("switched to %s personality", selected))
		}
		return true

	case input == configObj.LoadFiles:
		return handleFileLoad(chatManager, terminal, state, "")

	case strings.HasPrefix(input, configObj.LoadFiles+" "):
		dirPath := strings.TrimSpace(strings.TrimPrefix(input, configObj.LoadFiles+" "))
		return handleFileLoad(chatManager, terminal, state, dirPath)

	case input == configObj.CodeDump:
		return handleCodeDump(chatManager, terminal, state)

	case input == configObj.ShellRecord:
		if fromHelp {
			fmt.Printf("\033[93m%s - record shell session\033[0m\n", configObj.ShellRecord)
			return true
		}
		return handleShellRecord(chatManager, terminal, state)

	case strings.HasPrefix(input, configObj.ShellRecord+" "):
		command := strings.TrimPrefix(input, configObj.ShellRecord+" ")
		if fromHelp {
			fmt.Printf("\033[93m%s [command] - record shell session\033[0m\n", configObj.ShellRecord)
			return true
		}
		return handleShellCommand(command, chatManager, terminal, state)

	case strings.HasPrefix(input, configObj.EditorInput+" "):
		arg := strings.TrimSpace(strings.TrimPrefix(input, configObj.EditorInput+" "))

		if fromHelp {
			fmt.Printf("\033[93m%s - Advanced Text Editor Mode\033[0m\n", configObj.EditorInput)
			return true
		}

		if arg == "buff" {

			userInput, err := chatManager.HandleTerminalInput("")
			if err != nil {
				terminal.PrintError(fmt.Sprintf("%v", err))
				return true
			}
			chatManager.AddUserMessage(userInput)
			terminal.PrintInfo("Content loaded into buffer")
			return true
		}
		return true

	case input == configObj.EditorInput:
		if fromHelp {
			fmt.Printf("\033[93m%s - Advanced Text Editor Mode\033[0m\n", configObj.EditorInput)
			return true
		}

		editors := ui.GetAvailableEditors(configObj)
		if len(editors) == 0 {
			terminal.PrintError("No suitable editors found on this system")
			return true
		}

		selectedEditor, err := terminal.FzfSelect(editors, "select editor: ")
		if err != nil || selectedEditor == "" {
			return true
		}

		userInput, err := chatManager.HandleTerminalInput(selectedEditor)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("%v", err))
			return true
		}

		actions := []string{"SEND TO AI", "LOAD TO BUFFER"}
		choice, err := terminal.FzfSelect(actions, "action: ")
		if err != nil || choice == "" {
			return true
		}

		if choice == "SEND TO AI" {
			chatManager.AddUserMessage(userInput)

			ctx, animationCancel := context.WithCancel(context.Background())
			go terminal.ShowLoadingAnimation(ctx, "thinking")

			response, err := platformManager.SendChatRequest(chatManager.GetMessages(), chatManager.GetCurrentModel(), &state.StreamingCancel, &state.IsStreaming, animationCancel, terminal)

			animationCancel()

			if err != nil {
				if err.Error() == "request was interrupted" {
					chatManager.RemoveLastUserMessage()
					return true
				}
				terminal.PrintError(fmt.Sprintf("%v", err))
				return true
			}

			if platformManager.IsReasoningModel(chatManager.GetCurrentModel()) {
				theme := terminal.GetTheme()
				fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)
			} else {

				fmt.Println()
			}

			chatManager.AddAssistantMessage(response)
			chatManager.AddToHistory(userInput, response)
		} else {
			chatManager.AddUserMessage(userInput)
			terminal.PrintInfo("Content loaded into buffer")
		}
		return true

	case input == configObj.ExportChat || strings.HasPrefix(input, configObj.ExportChat+" "):

		targetFile := ""
		if strings.HasPrefix(input, configObj.ExportChat+" ") {
			targetFile = strings.TrimSpace(strings.TrimPrefix(input, configObj.ExportChat+" "))
		}
		err := handleExportChatInteractive(chatManager, terminal, state, targetFile)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error exporting chat: %v", err))
		}
		return true

	case input == configObj.Backtrack:
		backtrackedCount, err := chatManager.BacktrackHistory(terminal)
		if err != nil {
			terminal.PrintError(err.Error())
		} else {
			terminal.PrintInfo(fmt.Sprintf("backtracked by %d", backtrackedCount))
		}
		return true

	case input == "!a" || input == "!a exact":

		if !configObj.SaveAllSessions {
			terminal.PrintError("session search requires save_all_sessions to be enabled in config")

		}

		exact := input == "!a exact"
		session, err := chatManager.ManageSessions(terminal, exact)
		if err != nil {
			if err.Error() == "selection cancelled" {
				return true
			}
			terminal.PrintError(fmt.Sprintf("%v", err))
			return true
		}

		chatManager.RestoreSessionState(session)

		err = platformManager.Initialize()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error initializing client: %v", err))
			return true
		}

		if rl != nil {
			rl.SetPrompt(terminal.GetPrompt())
		}

		fmt.Printf("\033[91mrestored session from %s UTC\033[0m\n", time.Unix(session.Timestamp, 0).UTC().Format("2006-01-02 15:04:05"))

		for _, entry := range session.ChatHistory {
			if entry.User == state.Config.SystemPrompt {
				continue
			}

			if entry.User != "" {
				theme := terminal.GetTheme()
				fmt.Printf("%s USER \033[0m ❯ %s\n", theme.UserBox, entry.User)
			}

			if entry.Bot != "" {
				theme := terminal.GetTheme()
				fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, entry.Bot)
			}
		}

		return true

	case input == configObj.ScrapeURL:
		if fromHelp {
			fmt.Printf("\033[93m%s [url] - scrape URL(s)\033[0m\n", configObj.ScrapeURL)
			return true
		}

		historyURLs := terminal.ExtractURLsFromChatHistory(chatManager.GetChatHistory())
		messageURLs := terminal.ExtractURLsFromMessages(chatManager.GetMessages())

		seen := make(map[string]bool)
		var allURLs []string
		for _, url := range historyURLs {
			if !seen[url] {
				allURLs = append(allURLs, url)
				seen[url] = true
			}
		}
		for _, url := range messageURLs {
			if !seen[url] {
				allURLs = append(allURLs, url)
				seen[url] = true
			}
		}

		if len(allURLs) == 0 {
			terminal.PrintError("no URLs found in chat history")
			return true
		}

		selectedURLs, err := terminal.FzfMultiSelect(allURLs, "select urls to scrape (tab=multi): ")
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error selecting URLs: %v", err))
			return true
		}

		if len(selectedURLs) == 0 {
			return true
		}

		return handleScrapeURLs(selectedURLs, chatManager, terminal, state)

	case strings.HasPrefix(input, configObj.ScrapeURL+" "):
		urls := strings.Fields(strings.TrimPrefix(input, configObj.ScrapeURL+" "))
		return handleScrapeURLs(urls, chatManager, terminal, state)

	case input == configObj.WebSearch:
		if fromHelp {
			fmt.Printf("\033[93m%s [query] - web search\033[0m\n", configObj.WebSearch)
			return true
		}

		allSentences := terminal.ExtractSentencesFromChatHistory(chatManager.GetChatHistory(), chatManager.GetMessages())

		if len(allSentences) == 0 {
			terminal.PrintError("no sentences found in chat history")
			return true
		}

		selectedSentence, err := terminal.FzfSelect(allSentences, "select sentence to search: ")
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error selecting sentence: %v", err))
			return true
		}

		if selectedSentence == "" {
			return true
		}

		return handleWebSearch(selectedSentence, chatManager, terminal, state)

	case strings.HasPrefix(input, configObj.WebSearch+" "):
		query := strings.TrimPrefix(input, configObj.WebSearch+" ")
		return handleWebSearch(query, chatManager, terminal, state)

	case input == configObj.CopyToClipboard:
		err := terminal.CopyResponsesInteractive(chatManager.GetChatHistory(), chatManager.GetMessages())
		if err != nil {
			terminal.PrintError(fmt.Sprintf("%v", err))
		}
		return true

	case input == configObj.QuickCopyLatest:
		err := terminal.CopyLatestResponseToClipboard(chatManager.GetChatHistory())
		if err != nil {
			terminal.PrintError(fmt.Sprintf("%v", err))
		} else {
			terminal.PrintInfo("latest response copied to clipboard")
		}
		return true

	case input == configObj.MultiLine:
		var lines []string
		terminal.PrintInfo("multi-line mode (exit with '\\')")

		multiLineRl, err := readline.NewEx(&readline.Config{
			Prompt:      "... ",
			HistoryFile: "/dev/null",
		})
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error creating multi-line input: %v", err))
			return true
		}
		defer multiLineRl.Close()

		for {
			line, err := multiLineRl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt || err == io.EOF {
					return true
				}
				break
			}
			if line == configObj.MultiLine {
				break
			}
			lines = append(lines, line)
		}

		fullInput := strings.Join(lines, "\n")
		if strings.TrimSpace(fullInput) == "" {
			return true
		}

		chatManager.AddUserMessage(fullInput)

		ctx, animationCancel := context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "thinking")

		response, err := platformManager.SendChatRequest(chatManager.GetMessages(), chatManager.GetCurrentModel(), &state.StreamingCancel, &state.IsStreaming, animationCancel, terminal)

		animationCancel()

		if err != nil {
			if err.Error() == "request was interrupted" {
				chatManager.RemoveLastUserMessage()
				return true
			}
			terminal.PrintError(fmt.Sprintf("%v", err))
			return true
		}

		if platformManager.IsReasoningModel(chatManager.GetCurrentModel()) {
			if state.Config.IsPipedOutput {
				fmt.Printf("%s\n", response)
			} else {
				fmt.Print("\r\033[K")
				theme := terminal.GetTheme()
				fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)
			}
		}

		chatManager.AddAssistantMessage(response)
		chatManager.AddToHistory(fullInput, response)
		return true

	default:
		return false
	}
}

func handleFileLoad(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, dirPath string) bool {
	var files []string
	var err error
	var targetPath string

	if dirPath == "" {

		targetPath, _ = os.Getwd()
		files, err = terminal.GetCurrentDirFilesRecursive()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error reading current directory: %v", err))
			return true
		}
		if len(files) == 0 {
			terminal.PrintError("no files found in current directory")
			return true
		}
	} else {

		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			terminal.PrintError(fmt.Sprintf("directory does not exist: %s", dirPath))
			return true
		}
		targetPath = dirPath
		files, err = terminal.GetDirFilesRecursive(dirPath)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error reading directory %s: %v", dirPath, err))
			return true
		}
		if len(files) == 0 {
			terminal.PrintError(fmt.Sprintf("no files found in directory: %s", dirPath))
			return true
		}
	}

	if util.IsShallowLoadDir(state.Config, targetPath) {
		terminal.PrintInfo("shallow loading")
	}

	selections, err := terminal.FzfMultiSelect(files, "files: ")
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error selecting files: %v", err))
		return true
	}

	if len(selections) == 0 {
		return true
	}

	// Resolve full paths if using custom directory
	var fullPaths []string
	if dirPath != "" {
		for _, selection := range selections {
			fullPaths = append(fullPaths, filepath.Join(dirPath, selection))
		}
	} else {
		fullPaths = selections
	}

	content, err := terminal.LoadFileContent(fullPaths)
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error loading content: %v", err))
		return true
	}

	if content != "" {
		chatManager.AddUserMessage(content)
		if dirPath != "" {
			historySummary := fmt.Sprintf("Loaded from %s: %s", dirPath, strings.Join(selections, ", "))
			chatManager.AddToHistory(historySummary, "")
		} else {
			historySummary := fmt.Sprintf("Loaded: %s", strings.Join(selections, ", "))
			chatManager.AddToHistory(historySummary, "")
		}
	}

	return true
}

func handleCodeDump(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	codedump, err := terminal.CodeDump()
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error generating codedump: %v", err))
		return true
	}

	if codedump != "" {
		chatManager.AddUserMessage(codedump)
		chatManager.AddToHistory("Codedump loaded", "")
	}

	return true
}

func handleShellRecord(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	sessionContent, err := terminal.RecordShellSession()
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error recording shell session: %v", err))
		return true
	}

	if strings.TrimSpace(sessionContent) != "" {

		lines := strings.Split(sessionContent, "\n")
		var cleanedLines []string
		for _, line := range lines {
			if !strings.HasPrefix(line, "Script started on") && !strings.HasPrefix(line, "Script done on") {
				cleanedLines = append(cleanedLines, line)
			}
		}
		cleanedContent := strings.Join(cleanedLines, "\n")

		formattedContent := fmt.Sprintf("The user ran the following shell session and here is the output:\n\n---\n%s\n---", cleanedContent)

		chatManager.AddUserMessage(formattedContent)
		chatManager.AddToHistory("Shell session loaded", "")
	} else {
		terminal.PrintInfo("no activity recorded in shell session")
	}

	return true
}

func handleShellCommand(command string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	if command == "" {
		terminal.PrintError("no command specified")
		return true
	}

	cmd := exec.Command("sh", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		terminal.PrintError(fmt.Sprintf("failed to create stdout pipe: %v", err))
		return true
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		terminal.PrintError(fmt.Sprintf("failed to create stderr pipe: %v", err))
		return true
	}

	if err := cmd.Start(); err != nil {
		terminal.PrintError(fmt.Sprintf("failed to start command: %v", err))
		return true
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)

	// Capture output while streaming it live
	var outputBuffer strings.Builder

	done := make(chan bool, 2)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			outputBuffer.WriteString(line + "\n")
		}
		done <- true
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			outputBuffer.WriteString(line + "\n")
		}
		done <- true
	}()

	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	var cmdErr error
	select {
	case <-sigChan:

		if err := cmd.Process.Kill(); err != nil {
			terminal.PrintError(fmt.Sprintf("failed to kill command: %v", err))
		}
		fmt.Println("\ncommand interrupted")

		go func() {
			<-done
			<-done
		}()

		cmdErr = fmt.Errorf("command interrupted by user")

	case cmdErr = <-cmdDone:

		<-done
		<-done
	}

	outputStr := outputBuffer.String()

	var result string
	if cmdErr != nil {
		result = fmt.Sprintf("Command: %s\nError: %v\nOutput:\n%s", command, cmdErr, outputStr)
	} else {
		result = fmt.Sprintf("Command: %s\nOutput:\n%s", command, outputStr)
	}

	formattedContent := fmt.Sprintf("The user executed the following command and here is the output:\n\n---\n%s\n---", result)

	chatManager.AddUserMessage(formattedContent)
	chatManager.AddToHistory(fmt.Sprintf("!x %s", command), "Command executed and output added to context")

	return true
}

func isValidCodedumpDir(dirPath string) bool {

	if dirPath == "/" {
		return false
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func handleExportCodeBlocks(chatManager *chat.Manager, terminal *ui.Terminal) error {
	filePaths, err := chatManager.ExportCodeBlocks(terminal)
	if err != nil {
		return err
	}

	if len(filePaths) == 0 {
		terminal.PrintInfo("no code blocks found in the last response")
		return nil
	}

	for _, filePath := range filePaths {
		fmt.Println(filePath)
	}

	return nil
}

func handleExportChatInteractive(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, targetFile string) error {
	filePath, err := chatManager.ExportChatInteractive(terminal, targetFile)
	if err != nil {
		return err
	}

	if filePath != "" {
		fmt.Println(filePath)
	}

	return nil
}

func handleShowState(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) error {

	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05 MST")

	platform := chatManager.GetCurrentPlatform()
	model := chatManager.GetCurrentModel()

	chatHistory := chatManager.GetChatHistory()
	chatCount := len(chatHistory) - 1

	// Calculate total token count (including both history and messages for accuracy)
	var totalContent string
	for _, entry := range chatHistory {
		totalContent += entry.User + " " + entry.Bot + " "
	}

	for _, message := range chatManager.GetMessages() {
		totalContent += message.Content + " "
	}

	encoding := tokenizer.Cl100kBase
	enc, err := tokenizer.Get(encoding)
	if err != nil {
		return fmt.Errorf("error getting tokenizer: %v", err)
	}

	tokens, _, err := enc.Encode(totalContent)
	if err != nil {
		return fmt.Errorf("error encoding text: %v", err)
	}
	tokenCount := len(tokens)

	combinedDateTime := currentDate + " " + currentTime
	if state.Config.IsPipedOutput {
		fmt.Printf("%s %s\n", "date:", combinedDateTime)
		fmt.Printf("%s %s\n", "platform:", platform)
		fmt.Printf("%s %s\n", "model:", model)
		fmt.Printf("%s %d\n", "chats:", chatCount)
		fmt.Printf("%s %d\n", "tokens:", tokenCount)
	} else {
		fmt.Printf("\033[96m%s\033[0m \033[93m%s\033[0m\n", "date:", combinedDateTime)
		fmt.Printf("\033[96m%s\033[0m \033[95m%s\033[0m\n", "platform:", platform)
		fmt.Printf("\033[96m%s\033[0m \033[95m%s\033[0m\n", "model:", model)
		fmt.Printf("\033[96m%s\033[0m \033[92m%d\033[0m\n", "chats:", chatCount)
		fmt.Printf("\033[96m%s\033[0m \033[91m%d\033[0m\n", "tokens:", tokenCount)
	}

	return nil
}

func handleTokenCount(filePath string, model string, terminal *ui.Terminal, state *types.AppState) error {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	targetModel := model
	if targetModel == "" {
		targetModel = state.Config.CurrentModel
		if targetModel == "" {
			targetModel = state.Config.DefaultModel
		}
	}

	// Map model names to tokenizer encodings
	var encoding tokenizer.Encoding
	switch {
	case strings.Contains(strings.ToLower(targetModel), "gpt-4"):
		encoding = tokenizer.Cl100kBase
	case strings.Contains(strings.ToLower(targetModel), "gpt-3.5"):
		encoding = tokenizer.Cl100kBase
	case strings.Contains(strings.ToLower(targetModel), "gpt-2"):
		encoding = tokenizer.R50kBase
	case strings.Contains(strings.ToLower(targetModel), "claude"):
		encoding = tokenizer.Cl100kBase
	default:
		encoding = tokenizer.Cl100kBase
	}

	enc, err := tokenizer.Get(encoding)
	if err != nil {
		return fmt.Errorf("error getting tokenizer: %v", err)
	}

	tokens, _, err := enc.Encode(string(content))
	if err != nil {
		return fmt.Errorf("error encoding text: %v", err)
	}

	if state.Config.IsPipedOutput {
		fmt.Printf("%s %s\n", "file:", filePath)
		fmt.Printf("%s %s\n", "model:", targetModel)
		fmt.Printf("%s %d\n", "tokens:", len(tokens))
	} else {
		fmt.Printf("\033[96m%s\033[0m %s\n", "file:", filePath)
		fmt.Printf("\033[96m%s\033[0m \033[95m%s\033[0m\n", "model:", targetModel)
		fmt.Printf("\033[96m%s\033[0m \033[91m%d\033[0m\n", "tokens:", len(tokens))
	}

	return nil
}

// splitByDelimiters splits a string by both commas and pipes, trimming whitespace
func splitByDelimiters(input string) []string {

	parts := strings.Split(input, ",")
	var result []string

	for _, part := range parts {

		subParts := strings.Split(part, "|")
		for _, subPart := range subParts {
			trimmed := strings.TrimSpace(subPart)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}

	return result
}

// handleFlagWithPrompt sends context and prompt to AI, then displays response
// ctxContent: the loaded/scraped/searched content
// prompt: the user's query/instruction
func handleFlagWithPrompt(chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, ctxContent string, prompt string, noHistory bool) error {

	combinedMessage := ctxContent + "\n\n" + prompt

	chatManager.AddUserMessage(combinedMessage)

	// Pass animation context to SendChatRequest if needed
	var animationCancel context.CancelFunc
	if !state.Config.IsPipedOutput {
		var ctx context.Context
		ctx, animationCancel = context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "thinking")
	}

	response, err := platformManager.SendChatRequest(chatManager.GetMessages(), chatManager.GetCurrentModel(), &state.StreamingCancel, &state.IsStreaming, animationCancel, terminal)

	if animationCancel != nil {
		animationCancel()
	}

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return nil
		}
		return err
	}

	if platformManager.IsReasoningModel(chatManager.GetCurrentModel()) {
		if state.Config.IsPipedOutput {
			fmt.Printf("%s\n", response)
		} else {
			fmt.Printf("\033[92m%s\033[0m\n", response)
		}
	}

	chatManager.AddAssistantMessage(response)
	chatManager.AddToHistory(prompt, response)

	if state.Config.EnableSessionSave && !noHistory {
		if err := chatManager.SaveSessionState(); err != nil {
			terminal.PrintError(fmt.Sprintf("warning: failed to save session: %v", err))
		}
	}

	return nil
}

// handleScrapeURLs handles the !s command for scraping URLs
func handleScrapeURLs(urls []string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	if len(urls) == 0 {
		terminal.PrintError("no URLs provided")
		return true
	}

	content, err := terminal.ScrapeURLs(urls)
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error scraping URLs: %v", err))
		return true
	}

	if content != "" {
		chatManager.AddUserMessage(content)
		historySummary := fmt.Sprintf("Scraped: %s", strings.Join(urls, ", "))
		chatManager.AddToHistory(historySummary, "")
	}

	return true
}

// handleWebSearch handles the !w command for web search
func handleWebSearch(query string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	if query == "" {
		terminal.PrintError("no search query provided")
		return true
	}

	content, err := terminal.WebSearch(query)
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error searching: %v", err))
		return true
	}

	if content != "" {
		chatManager.AddUserMessage(content)
		historySummary := fmt.Sprintf("Web search: %s", query)
		chatManager.AddToHistory(historySummary, "")
	}

	return true
}

// generateUniqueCodeDumpFilename generates a unique filename for code dump with collision detection
func generateUniqueCodeDumpFilename(currentDir, content string) string {
	baseHash := chat.GenerateHashFromContent(content, 8)
	filename := fmt.Sprintf("viren_cd%s.txt", baseHash)
	fullPath := filepath.Join(currentDir, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return filename
	}

	for offset := 1; offset <= 10; offset++ {
		newHash := chat.GenerateHashFromContentWithOffset(content, 8, offset)
		filename = fmt.Sprintf("viren_cd%s.txt", newHash)
		fullPath = filepath.Join(currentDir, filename)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return filename
		}
	}

	for counter := 1; counter <= 999; counter++ {
		filename = fmt.Sprintf("viren_cd%s_%03d.txt", baseHash, counter)
		fullPath = filepath.Join(currentDir, filename)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return filename
		}
	}

	return fmt.Sprintf("viren_cd%s.txt", uuid.New().String())
}

// handleAllModels handles the !o command for selecting from all available models
func handleAllModels(chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	// Create channels for async operation
	type modelResult struct {
		models []string
		err    error
	}
	resultChan := make(chan modelResult)

	go func() {
		models, err := platformManager.FetchAllModelsAsync()
		resultChan <- modelResult{models, err}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "fetching models")

	result := <-resultChan
	cancel()

	if result.err != nil {
		terminal.PrintError(fmt.Sprintf("%v", result.err))
		return true
	}

	if len(result.models) == 0 {
		terminal.PrintError("no models found")
		return true
	}

	models := result.models

	// Create a map to store platform and model info indexed by display string
	type modelInfo struct {
		platform string
		model    string
	}
	modelMap := make(map[string]modelInfo)

	for _, m := range models {
		parts := strings.SplitN(m, "|", 2)
		if len(parts) == 2 {
			platform := parts[0]
			modelName := parts[1]
			modelMap[m] = modelInfo{platform, modelName}
		}
	}

	selectedModel, err := terminal.FzfSelect(models, "model: ")
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error selecting model: %v", err))
		return true
	}

	if selectedModel == "" {
		return true
	}

	info, exists := modelMap[selectedModel]
	if !exists {
		terminal.PrintError("invalid model selection")
		return true
	}

	platformName := info.platform
	modelName := info.model

	currentPlatform := state.Config.CurrentPlatform

	state.Config.CurrentPlatform = platformName
	state.Config.CurrentModel = modelName

	chatManager.SetCurrentPlatform(platformName)
	chatManager.SetCurrentModel(modelName)
	state.Config.CurrentPlatform = platformName
	state.Config.CurrentModel = modelName

	config.SaveConfigToFile(state.Config)

	if platformName != currentPlatform {

		state.Config.CurrentBaseURL = ""
		err := platformManager.Initialize()
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error initializing client: %v", err))

			state.Config.CurrentPlatform = currentPlatform
			return true
		}
	}

	if !state.Config.MuteNotifications {
		terminal.PrintModelSwitch(modelName)
	}

	return true
}
