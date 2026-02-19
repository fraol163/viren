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

	buildTime	= "unknown"
	gitCommit	= "unknown"
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
		vFlag	= flag.Bool("v", false, "Show version")
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

func runInteractiveMode(chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, noHistory bool) {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:	terminal.GetPrompt(),
		InterruptPrompt:	"",
		EOFPrompt:	"exit",
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
				Prompt:	"... ",
				HistoryFile:	"/dev/null",
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
		return handleFileLoad(chatManager, terminal, state, platformManager, "")

	case strings.HasPrefix(input, configObj.LoadFiles+" "):
		dirPath := strings.TrimSpace(strings.TrimPrefix(input, configObj.LoadFiles+" "))
		return handleFileLoad(chatManager, terminal, state, platformManager, dirPath)

	case input == configObj.CodeDump:
		return handleCodeDump(chatManager, terminal, state, platformManager)

	case input == configObj.ShellRecord:
		if fromHelp {
			fmt.Printf("\033[93m%s - record shell session\033[0m\n", configObj.ShellRecord)
			return true
		}
		return handleShellRecord(chatManager, terminal, state, platformManager)

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
		terminal.PrintError(fmt.Sprintf("unknown argument: %s. Use 'buff' for buffer mode.", arg))
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

	case input == configObj.AnswerSearch || input == configObj.AnswerSearch+" exact":
		if !configObj.SaveAllSessions {
			terminal.PrintError("session search requires save_all_sessions to be enabled in config")
			return true
		}

		exact := input == configObj.AnswerSearch+" exact"
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

		return handleScrapeURLs(selectedURLs, chatManager, terminal, state, platformManager)

	case strings.HasPrefix(input, configObj.ScrapeURL+" "):
		urls := strings.Fields(strings.TrimPrefix(input, configObj.ScrapeURL+" "))
		return handleScrapeURLs(urls, chatManager, terminal, state, platformManager)

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

		return handleWebSearch(selectedSentence, chatManager, terminal, state, platformManager)

	case strings.HasPrefix(input, configObj.WebSearch+" "):
		query := strings.TrimPrefix(input, configObj.WebSearch+" ")
		return handleWebSearch(query, chatManager, terminal, state, platformManager)

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
			Prompt:	"... ",
			HistoryFile:	"/dev/null",
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

	case input == configObj.Regenerate:
		return handleRegenerate(chatManager, terminal, state, platformManager)

	case input == configObj.ExplainCode:
		return handleExplainCode(chatManager, terminal, state, platformManager)

	case input == configObj.Summarize:
		return handleSummarize(chatManager, terminal, state, platformManager)

	case input == configObj.GenerateTests:
		return handleGenerateTests(chatManager, terminal, state, platformManager)

	case input == configObj.GenerateDocs:
		return handleGenerateDocs(chatManager, terminal, state, platformManager)

	case input == configObj.OptimizeCode:
		return handleOptimizeCode(chatManager, terminal, state, platformManager)

	case strings.HasPrefix(input, configObj.GitCommand+" "):
		command := strings.TrimSpace(strings.TrimPrefix(input, configObj.GitCommand+" "))
		return handleGitCommand(command, chatManager, terminal, state, platformManager)

	case strings.HasPrefix(input, configObj.CompareFiles+" "):
		files := strings.Fields(strings.TrimPrefix(input, configObj.CompareFiles+" "))
		return handleCompareFiles(files, chatManager, terminal, state, platformManager)

	case strings.HasPrefix(input, configObj.TranslateCode+" "):
		targetLang := strings.TrimSpace(strings.TrimPrefix(input, configObj.TranslateCode+" "))
		return handleTranslateCode(targetLang, chatManager, terminal, state, platformManager)

	case strings.HasPrefix(input, configObj.FindReplace+" "):

		arg := strings.TrimSpace(strings.TrimPrefix(input, configObj.FindReplace+" "))
		if len(arg) >= 3 && arg[0] == '/' {
			parts := strings.Split(arg[1:], "/")
			if len(parts) >= 2 {
				return handleFindReplace(parts[0], parts[1], chatManager, terminal, state)
			}
		}
		terminal.PrintError("invalid format. Use: !f /old/new/")
		return true

	case input == configObj.CommandReference:
		return handleCommandReference(terminal)

	default:
		return false
	}
}

func handleFileLoad(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager, dirPath string) bool {
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

		terminal.PrintInfo(fmt.Sprintf("loaded %d file(s)", len(selections)))

		preview := content
		if len(preview) > 1500 {
			preview = preview[:1500] + "\n\n... (content truncated)"
		}

		theme := terminal.GetTheme()
		fmt.Printf("\n%s LOADED CONTENT \033[0m\n", theme.AssistantBox)
		fmt.Printf("%s\n", strings.Repeat("─", 60))
		fmt.Printf("%s\n", preview)
		fmt.Printf("%s\n\n", strings.Repeat("─", 60))

		if dirPath != "" {
			historySummary := fmt.Sprintf("Loaded from %s: %s", dirPath, strings.Join(selections, ", "))
			chatManager.AddToHistory(historySummary, "")
		} else {
			historySummary := fmt.Sprintf("Loaded: %s", strings.Join(selections, ", "))
			chatManager.AddToHistory(historySummary, "")
		}

		formattedContent := fmt.Sprintf("The user loaded the following files:\n\n---\n%s\n---", content)
		chatManager.AddUserMessage(formattedContent)

		terminal.PrintInfo("analyzing loaded content...")

		ctx, animationCancel := context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "analyzing")

		response, err := platformManager.SendChatRequest(
			chatManager.GetMessages(),
			chatManager.GetCurrentModel(),
			&state.StreamingCancel,
			&state.IsStreaming,
			animationCancel,
			terminal,
		)

		animationCancel()

		if err != nil {
			if err.Error() == "request was interrupted" {
				chatManager.RemoveLastUserMessage()
				return true
			}
			terminal.PrintError(fmt.Sprintf("error analyzing content: %v", err))
			return true
		}

		if !state.Config.IsPipedOutput {
			fmt.Println()
		}
		theme = terminal.GetTheme()
		fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

		chatManager.AddAssistantMessage(response)
	}

	return true
}

func handleCodeDump(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	codedump, err := terminal.CodeDump()
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error generating codedump: %v", err))
		return true
	}

	if codedump != "" {

		preview := codedump
		if len(preview) > 1500 {
			preview = preview[:1500] + "\n\n... (content truncated)"
		}

		theme := terminal.GetTheme()
		fmt.Printf("\n%s CODEDUMP CONTENT \033[0m\n", theme.AssistantBox)
		fmt.Printf("%s\n", strings.Repeat("─", 60))
		fmt.Printf("%s\n", preview)
		fmt.Printf("%s\n\n", strings.Repeat("─", 60))

		chatManager.AddToHistory("Codedump loaded", "")
		formattedContent := fmt.Sprintf("The user dumped the following codebase:\n\n---\n%s\n---", codedump)
		chatManager.AddUserMessage(formattedContent)

		terminal.PrintInfo("analyzing codebase...")

		ctx, animationCancel := context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "analyzing")

		response, err := platformManager.SendChatRequest(
			chatManager.GetMessages(),
			chatManager.GetCurrentModel(),
			&state.StreamingCancel,
			&state.IsStreaming,
			animationCancel,
			terminal,
		)

		animationCancel()

		if err != nil {
			if err.Error() == "request was interrupted" {
				chatManager.RemoveLastUserMessage()
				return true
			}
			terminal.PrintError(fmt.Sprintf("error analyzing codebase: %v", err))
			return true
		}

		if !state.Config.IsPipedOutput {
			fmt.Println()
		}
		theme = terminal.GetTheme()
		fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

		chatManager.AddAssistantMessage(response)
	}

	return true
}

func handleShellRecord(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
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

		preview := cleanedContent
		if len(preview) > 1500 {
			preview = preview[:1500] + "\n\n... (content truncated)"
		}

		theme := terminal.GetTheme()
		fmt.Printf("\n%s SHELL SESSION OUTPUT \033[0m\n", theme.AssistantBox)
		fmt.Printf("%s\n", strings.Repeat("─", 60))
		fmt.Printf("%s\n", preview)
		fmt.Printf("%s\n\n", strings.Repeat("─", 60))

		formattedContent := fmt.Sprintf("The user ran the following shell session and here is the output:\n\n---\n%s\n---", cleanedContent)
		chatManager.AddUserMessage(formattedContent)
		chatManager.AddToHistory("Shell session loaded", "")

		terminal.PrintInfo("analyzing shell session output...")

		ctx, animationCancel := context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "analyzing")

		response, err := platformManager.SendChatRequest(
			chatManager.GetMessages(),
			chatManager.GetCurrentModel(),
			&state.StreamingCancel,
			&state.IsStreaming,
			animationCancel,
			terminal,
		)

		animationCancel()

		if err != nil {
			if err.Error() == "request was interrupted" {
				chatManager.RemoveLastUserMessage()
				return true
			}
			terminal.PrintError(fmt.Sprintf("error analyzing shell output: %v", err))
			return true
		}

		if !state.Config.IsPipedOutput {
			fmt.Println()
		}
		theme = terminal.GetTheme()
		fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

		chatManager.AddAssistantMessage(response)
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

func handleRegenerate(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	history := chatManager.GetChatHistory()
	if len(history) < 2 {
		terminal.PrintError("no previous conversation to regenerate")
		return true
	}

	// Find the last user message
	var lastUserMsg string
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].User != "" && history[i].User != state.Config.SystemPrompt {
			lastUserMsg = history[i].User
			break
		}
	}

	if lastUserMsg == "" {
		terminal.PrintError("no user message found to regenerate")
		return true
	}

	if len(history) > 1 && history[len(history)-1].Bot != "" {
		chatManager.RemoveLastUserMessage()
	}

	terminal.PrintInfo("regenerating response...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "regenerating")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			return true
		}
		terminal.PrintError(fmt.Sprintf("error regenerating: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)
	chatManager.AddToHistory(lastUserMsg, response)

	return true
}

func handleExplainCode(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	messages := chatManager.GetMessages()
	if len(messages) < 2 {
		terminal.PrintError("no code to explain")
		return true
	}

	lastUserMsg := messages[len(messages)-1]
	if lastUserMsg.Role != "user" {
		terminal.PrintError("no code found to explain")
		return true
	}

	prompt := fmt.Sprintf("Please explain the following code step by step in detail:\n\n%s", lastUserMsg.Content)

	chatManager.AddUserMessage(prompt)
	terminal.PrintInfo("explaining code...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "explaining")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error explaining code: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleSummarize(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	messages := chatManager.GetMessages()
	if len(messages) < 2 {
		terminal.PrintError("no content to summarize")
		return true
	}

	lastUserMsg := messages[len(messages)-1]
	if lastUserMsg.Role != "user" {
		terminal.PrintError("no content found to summarize")
		return true
	}

	prompt := fmt.Sprintf("Please summarize the following content, highlighting the key points:\n\n%s", lastUserMsg.Content)

	chatManager.AddUserMessage(prompt)
	terminal.PrintInfo("summarizing content...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "summarizing")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error summarizing: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleGenerateTests(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	code := extractLastCodeBlock(chatManager)
	if code == "" {
		terminal.PrintError("no code found to generate tests for")
		return true
	}

	prompt := fmt.Sprintf("Generate comprehensive unit tests for the following code. Use best practices and cover edge cases:\n\n```%s", code)

	chatManager.AddUserMessage(prompt)
	terminal.PrintInfo("generating tests...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "generating tests")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error generating tests: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleGenerateDocs(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	code := extractLastCodeBlock(chatManager)
	if code == "" {
		terminal.PrintError("no code found to generate documentation for")
		return true
	}

	prompt := fmt.Sprintf("Generate comprehensive documentation for the following code. Include function descriptions, parameters, return values, and usage examples:\n\n```%s", code)

	chatManager.AddUserMessage(prompt)
	terminal.PrintInfo("generating documentation...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "generating docs")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error generating docs: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleOptimizeCode(chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	code := extractLastCodeBlock(chatManager)
	if code == "" {
		terminal.PrintError("no code found to optimize")
		return true
	}

	prompt := fmt.Sprintf("Optimize the following code for performance, readability, and best practices. Explain the improvements:\n\n```%s", code)

	chatManager.AddUserMessage(prompt)
	terminal.PrintInfo("optimizing code...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "optimizing")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error optimizing code: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleGitCommand(command string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	if command == "" {
		terminal.PrintError("no git command specified")
		return true
	}

	cmd := exec.Command("git", strings.Fields(command)...)
	output, err := cmd.CombinedOutput()

	var result string
	if err != nil {
		result = fmt.Sprintf("Git command '%s' failed:\nError: %v\nOutput:\n%s", command, err, string(output))
	} else {
		result = fmt.Sprintf("Git command '%s' succeeded:\n%s", command, string(output))
	}

	theme := terminal.GetTheme()
	fmt.Printf("\n%s GIT OUTPUT \033[0m\n", theme.AssistantBox)
	fmt.Printf("%s\n", strings.Repeat("─", 60))
	fmt.Printf("%s\n", result)
	fmt.Printf("%s\n\n", strings.Repeat("─", 60))

	formattedContent := fmt.Sprintf("The user ran the following git command: %s\n\n---\n%s\n---", command, result)
	chatManager.AddUserMessage(formattedContent)
	chatManager.AddToHistory(fmt.Sprintf("!git %s", command), "")

	if strings.Contains(command, "diff") || strings.Contains(command, "log") || strings.Contains(command, "status") {
		terminal.PrintInfo("analyzing git output...")

		ctx, animationCancel := context.WithCancel(context.Background())
		go terminal.ShowLoadingAnimation(ctx, "analyzing")

		response, err := platformManager.SendChatRequest(
			chatManager.GetMessages(),
			chatManager.GetCurrentModel(),
			&state.StreamingCancel,
			&state.IsStreaming,
			animationCancel,
			terminal,
		)

		animationCancel()

		if err != nil {
			if err.Error() == "request was interrupted" {
				chatManager.RemoveLastUserMessage()
				return true
			}
			terminal.PrintError(fmt.Sprintf("error analyzing git output: %v", err))
			return true
		}

		if !state.Config.IsPipedOutput {
			fmt.Println()
		}
		theme = terminal.GetTheme()
		fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

		chatManager.AddAssistantMessage(response)
	}

	return true
}

func handleCompareFiles(files []string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	if len(files) < 2 {
		terminal.PrintError("please provide at least 2 files to compare")
		return true
	}

	var contents []string
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			terminal.PrintError(fmt.Sprintf("error reading file %s: %v", file, err))
			return true
		}
		contents = append(contents, string(content))
	}

	var prompt strings.Builder
	prompt.WriteString("Compare the following files and highlight the key differences, similarities, and improvements:\n\n")
	for i, content := range contents {
		prompt.WriteString(fmt.Sprintf("=== File %d: %s ===\n%s\n\n", i+1, files[i], content))
	}

	chatManager.AddUserMessage(prompt.String())
	terminal.PrintInfo(fmt.Sprintf("comparing %d files...", len(files)))

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "comparing")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error comparing files: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleTranslateCode(targetLang string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	code := extractLastCodeBlock(chatManager)
	if code == "" {
		terminal.PrintError("no code found to translate")
		return true
	}

	if targetLang == "" {
		terminal.PrintError("please specify target language (e.g., !translate python)")
		return true
	}

	prompt := fmt.Sprintf("Translate the following code to %s. Maintain the same functionality and follow best practices for the target language:\n\n```%s", targetLang, code)

	chatManager.AddUserMessage(prompt)
	terminal.PrintInfo(fmt.Sprintf("translating code to %s...", targetLang))

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "translating")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error translating code: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme := terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleFindReplace(find, replace string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState) bool {
	if find == "" {
		terminal.PrintError("please specify text to find (e.g., !f /old/new/)")
		return true
	}

	code := extractLastCodeBlock(chatManager)
	if code == "" {
		terminal.PrintError("no code found to perform find/replace")
		return true
	}

	newCode := strings.ReplaceAll(code, find, replace)

	theme := terminal.GetTheme()
	fmt.Printf("\n%s FIND/REPLACE RESULT \033[0m\n", theme.AssistantBox)
	fmt.Printf("%s\n", strings.Repeat("─", 60))
	fmt.Printf("Replaced '%s' with '%s'\n\n", find, replace)
	fmt.Printf("Result:\n```%s\n```\n", newCode)
	fmt.Printf("%s\n\n", strings.Repeat("─", 60))

	terminal.PrintInfo("find/replace completed. Use !e to export code blocks.")

	return true
}

func handleCommandReference(terminal *ui.Terminal) bool {
	commands := getCommandReference()

	theme := terminal.GetTheme()
	fmt.Printf("\n%s COMMAND REFERENCE \033[0m\n", theme.AssistantBox)
	fmt.Printf("%s\n", strings.Repeat("═", 70))

	for _, cmd := range commands {
		fmt.Printf("\n\033[1;96m%s\033[0m - %s\n", cmd.Command, cmd.Description)
		fmt.Printf("  \033[93mUsage:\033[0m %s\n", cmd.Usage)
		if cmd.Example != "" {
			fmt.Printf("  \033[92mExample:\033[0m %s\n", cmd.Example)
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("═", 70))
	fmt.Printf("\033[96mTotal commands:\033[0m %d\n", len(commands))

	return true
}

type CommandInfo struct {
	Command	string
	Description	string
	Usage	string
	Example	string
}

func getCommandReference() []CommandInfo {
	return []CommandInfo{

		{"!q", "Quit Viren", "!q", ""},
		{"!h", "Show help menu", "!h", ""},
		{"!c", "Clear chat history and screen", "!c", ""},
		{"!m", "Switch AI model", "!m [model]", "!m gpt-4"},
		{"!p", "Switch AI platform", "!p [platform]", "!p anthropic"},
		{"!u", "Change AI personality", "!u", ""},
		{"!v", "Change domain mode", "!v", ""},
		{"!z", "Change theme", "!z", ""},
		{"!e", "Export chat or code blocks", "!e [filename]", "!e output.txt"},
		{"!b", "Backtrack chat history", "!b", ""},
		{"!a", "Manage/load sessions", "!a [exact]", "!a exact"},
		{"!y", "Copy response to clipboard", "!y", ""},
		{"cc", "Quick copy latest response", "cc", ""},
		{"!d", "Dump codebase for analysis", "!d [dir]", "!d ./src"},
		{"!x", "Record shell session or run command", "!x [command]", "!x ls -la"},
		{"!l", "Load files into context", "!l [dir]", "!l ./config"},
		{"!s", "Scrape URL content", "!s [url]", "!s https://example.com"},
		{"!w", "Web search", "!w [query]", "!w Go programming"},
		{"!t", "Open text editor", "!t [buff]", "!t buff"},
		{"\\", "Multi-line input mode", "\\", ""},

		{"!r", "Regenerate last AI response", "!r", ""},
		{"!explain", "Explain code in detail", "!explain", ""},
		{"!summarize", "Summarize content", "!summarize", ""},
		{"!test", "Generate unit tests for code", "!test", ""},
		{"!doc", "Generate documentation for code", "!doc", ""},
		{"!optimize", "Optimize code for performance", "!optimize", ""},
		{"!git", "Run git commands with AI analysis", "!git [command]", "!git diff"},
		{"!compare", "Compare multiple files", "!compare file1 file2", "!compare main.go main_old.go"},
		{"!translate", "Translate code to another language", "!translate [language]", "!translate python"},
		{"!f", "Find and replace in code", "!f /old/new/", "!f /foo/bar/"},
		{"!cmd", "Show this command reference", "!cmd", ""},
	}
}

func extractLastCodeBlock(chatManager *chat.Manager) string {
	messages := chatManager.GetMessages()
	if len(messages) == 0 {
		return ""
	}

	lastMsg := messages[len(messages)-1]
	content := lastMsg.Content

	codeBlockRegex := regexp.MustCompile("(?s)```([a-zA-Z0-9]*)\n(.*?)\n```")
	matches := codeBlockRegex.FindAllStringSubmatch(content, -1)

	if len(matches) > 0 {
		lastMatch := matches[len(matches)-1]
		return lastMatch[2]
	}

	return content
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

func handleFlagWithPrompt(chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState, ctxContent string, prompt string, noHistory bool) error {

	combinedMessage := ctxContent + "\n\n" + prompt

	chatManager.AddUserMessage(combinedMessage)

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

func handleScrapeURLs(urls []string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	if len(urls) == 0 {
		terminal.PrintError("no URLs provided")
		return true
	}

	terminal.PrintInfo(fmt.Sprintf("scraping %d URL(s)...", len(urls)))

	content, err := terminal.ScrapeURLs(urls)
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error scraping URLs: %v", err))
		return true
	}

	if content == "" || strings.TrimSpace(content) == "" {
		terminal.PrintError("no content could be extracted from the provided URLs")
		return true
	}

	theme := terminal.GetTheme()
	fmt.Printf("\n%s SCRAPED CONTENT \033[0m\n", theme.AssistantBox)
	fmt.Printf("%s\n", strings.Repeat("─", 60))

	preview := content
	if len(preview) > 2000 {
		preview = preview[:2000] + "\n\n... (content truncated, full content sent to AI)"
	}
	fmt.Printf("%s\n", preview)
	fmt.Printf("%s\n\n", strings.Repeat("─", 60))

	formattedContent := fmt.Sprintf("The user scraped the following content from URL(s): %s\n\n---\n%s\n---", strings.Join(urls, ", "), content)
	chatManager.AddUserMessage(formattedContent)
	historySummary := fmt.Sprintf("Scraped: %s", strings.Join(urls, ", "))
	chatManager.AddToHistory(historySummary, "")

	terminal.PrintInfo("analyzing scraped content...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "analyzing")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error analyzing content: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme = terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

func handleWebSearch(query string, chatManager *chat.Manager, terminal *ui.Terminal, state *types.AppState, platformManager *platform.Manager) bool {
	if query == "" {
		terminal.PrintError("no search query provided")
		return true
	}

	terminal.PrintInfo(fmt.Sprintf("searching web for: %s", query))

	content, err := terminal.WebSearch(query)
	if err != nil {
		terminal.PrintError(fmt.Sprintf("error searching: %v", err))
		return true
	}

	if content == "" || strings.TrimSpace(content) == "" {
		terminal.PrintError("no search results found")
		return true
	}

	theme := terminal.GetTheme()
	fmt.Printf("\n%s SEARCH RESULTS \033[0m\n", theme.AssistantBox)
	fmt.Printf("%s\n", strings.Repeat("─", 60))
	fmt.Printf("%s\n", content)
	fmt.Printf("%s\n\n", strings.Repeat("─", 60))

	formattedContent := fmt.Sprintf("The user performed a web search for: %s\n\n---\n%s\n---", query, content)
	chatManager.AddUserMessage(formattedContent)
	historySummary := fmt.Sprintf("Web search: %s", query)
	chatManager.AddToHistory(historySummary, "")

	terminal.PrintInfo("analyzing search results...")

	ctx, animationCancel := context.WithCancel(context.Background())
	go terminal.ShowLoadingAnimation(ctx, "analyzing")

	response, err := platformManager.SendChatRequest(
		chatManager.GetMessages(),
		chatManager.GetCurrentModel(),
		&state.StreamingCancel,
		&state.IsStreaming,
		animationCancel,
		terminal,
	)

	animationCancel()

	if err != nil {
		if err.Error() == "request was interrupted" {
			chatManager.RemoveLastUserMessage()
			return true
		}
		terminal.PrintError(fmt.Sprintf("error analyzing search results: %v", err))
		return true
	}

	if !state.Config.IsPipedOutput {
		fmt.Println()
	}
	theme = terminal.GetTheme()
	fmt.Printf("%s ASSISTANT \033[0m ❯ \033[92m%s\033[0m\n", theme.AssistantBox, response)

	chatManager.AddAssistantMessage(response)

	return true
}

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

func handleAllModels(chatManager *chat.Manager, platformManager *platform.Manager, terminal *ui.Terminal, state *types.AppState) bool {

	type modelResult struct {
		models	[]string
		err	error
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

	type modelInfo struct {
		platform	string
		model	string
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
