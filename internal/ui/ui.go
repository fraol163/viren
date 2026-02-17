package ui

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fraol163/viren/internal/util"
	"github.com/fraol163/viren/pkg/types"
	"github.com/ledongthuc/pdf"
	"github.com/lu4p/cat"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/tealeg/xlsx/v3"
	"golang.org/x/net/html"
)

var (
	codeBlockRegex	= regexp.MustCompile("(?s)```([a-zA-Z0-9]*)\n(.*?)\n```")
	urlRegex	= regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`)
	sentenceRegex	= regexp.MustCompile(`[.!?]+\s+`)
)

type Terminal struct {
	config *types.Config
}

func NewTerminal(config *types.Config) *Terminal {
	return &Terminal{
		config: config,
	}
}

func (t *Terminal) GetTheme() types.Theme {
	if t.config.CurrentTheme == "" {
		return GetDefaultTheme()
	}
	return GetThemeByID(t.config.CurrentTheme)
}

func (t *Terminal) SetTheme(themeID string) {
	t.config.CurrentTheme = themeID
}

func (t *Terminal) ApplyTheme() {
	if t.config.IsPipedOutput {
		return
	}
	theme := t.GetTheme()

	bgColor := theme.BgColor
	switch t.config.CurrentPlatform {
	case "google":
		bgColor = "\033]11;#000000\007"
	case "openai":
		bgColor = "\033]11;#000808\007"
	case "anthropic":
		bgColor = "\033]11;#0a0a0a\007"
	case "groq":
		bgColor = "\033]11;#000008\007"
	case "deepseek":
		bgColor = "\033]11;#000a00\007"
	case "xai":
		bgColor = "\033]11;#000000\007"
	case "ollama":
		bgColor = "\033]11;#0a0a00\007"
	}

	if bgColor != "" {
		fmt.Print(bgColor)
	}
	if theme.FgColor != "" {
		fmt.Print(theme.FgColor)
	}
}

func (t *Terminal) GetPrompt() string {
	theme := t.GetTheme()
	return fmt.Sprintf("%s USER \033[0m ❯ ", theme.UserBox)
}

func (t *Terminal) ShowLogo() {
	if t.config.IsPipedOutput {
		return
	}
	theme := t.GetTheme()
	fmt.Print(theme.LogoColor)
	fmt.Println(`██╗      ██╗   ██╗██╗██████╗ ███████╗███╗   ██╗`)
	fmt.Println(`╚██╗     ██║   ██║██║██╔══██╗██╔════╝████╗  ██║`)
	fmt.Println(` ╚██╗    ██║   ██║██║██████╔╝█████╗  ██╔██╗ ██║`)
	fmt.Println(` ██╔╝    ╚██╗ ██╔╝██║██╔══██╗██╔══╝  ██║╚██╗██║`)
	fmt.Println(`██╔╝      ╚████╔╝ ██║██║  ██║███████╗██║ ╚████║`)
	fmt.Println(`╚═╝        ╚═══╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═══╝`)
	fmt.Print("\033[0m")

	fmt.Printf("\n  %s✦\033[0m \033[1mWelcome to VIREN, your advanced neural command interface.\033[0m\n", theme.LogoColor)
	fmt.Printf("\n  \033[1m%sSYSTEM CAPABILITIES:\033[0m\n", theme.MutedColor)
	fmt.Printf("  %s•\033[0m \033[1mNeural Interaction:\033[0m Instant access to elite LLMs (GPT-4, Claude 3, DeepSeek).\n", theme.LogoColor)
	fmt.Printf("  %s•\033[0m \033[1mContext Injection:\033[0m Seamless ingestion of files, dirs, and live web content.\n", theme.LogoColor)
	fmt.Printf("  %s•\033[0m \033[1mDev Suite:\033[0m Smart code export, terminal session capture, and fuzzy-search.\n", theme.LogoColor)
	fmt.Printf("  %s•\033[0m \033[1mLocal Privacy:\033[0m Zero-cloud persistence. All logic and history stay on-disk.\n", theme.LogoColor)

	fmt.Printf("\n  \033[91m!q\033[0m quit %s•\033[0m  \033[91m!h\033[0m help %s•\033[0m  \033[91m!c\033[0m clear %s•\033[0m  \033[91m!m\033[0m model %s•\033[0m  \033[91m!u\033[0m personality %s•\033[0m  \033[91m!v\033[0m mode %s•\033[0m  \033[91m!z\033[0m theme %s•\033[0m  \033[91m!p\033[0m platform\n",
		theme.BorderColor, theme.BorderColor, theme.BorderColor, theme.BorderColor, theme.BorderColor, theme.BorderColor, theme.BorderColor)

	modeStr := ""
	if t.config.CurrentMode != "" && t.config.CurrentMode != "standard" {
		modeStr = fmt.Sprintf("%s/\033[0m\033[1;93m%s\033[0m", theme.BorderColor, t.config.CurrentMode)
	}
	persStr := ""
	if t.config.CurrentPersonality != "" && t.config.CurrentPersonality != "balanced" {
		persStr = fmt.Sprintf("%s/\033[0m\033[1;35m%s\033[0m", theme.BorderColor, t.config.CurrentPersonality)
	}
	fmt.Printf("  \033[92m○\033[0m \033[1mCONNECTED\033[0m  %s[\033[0m\033[1;96m%s\033[0m%s/\033[0m\033[1;95m%s\033[0m%s%s%s]\033[0m\n", theme.BorderColor, t.config.CurrentPlatform, theme.BorderColor, t.config.CurrentModel, modeStr, persStr, theme.BorderColor)
}

func (t *Terminal) ClearTerminal() {
	if t.config.IsPipedOutput {
		return
	}

	fmt.Print("\033[H\033[2J")
}

func escapeShellArg(arg string) string {
	return "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
}

func ContainsAllOption(items []string) bool {
	for _, item := range items {
		if strings.HasPrefix(item, ">all") {
			return true
		}
	}
	return false
}

func (t *Terminal) runFzfCore(fzfArgs []string, inputText string) ([]byte, bool, error) {
	tempDir, err := util.GetTempDir()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get temp directory: %w", err)
	}

	outputFile := filepath.Join(tempDir, fmt.Sprintf("fzf-output-%d", time.Now().UnixNano()))
	defer os.Remove(outputFile)

	var escapedArgs []string
	for _, arg := range fzfArgs {
		escapedArgs = append(escapedArgs, escapeShellArg(arg))
	}

	cmdStr := fmt.Sprintf("fzf %s > %s", strings.Join(escapedArgs, " "), escapeShellArg(outputFile))

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdin = strings.NewReader(inputText)
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 130 || exitErr.ExitCode() == 1 {
			return nil, true, nil
		}
		return nil, false, fmt.Errorf("fzf failed: %w", err)
	} else if err != nil {
		return nil, false, fmt.Errorf("fzf execution failed: %w", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, true, nil
		}
		return nil, false, fmt.Errorf("failed to read fzf output: %w", err)
	}

	return content, false, nil
}

func (t *Terminal) runFzfSSHSafe(fzfArgs []string, inputText string) (string, error) {
	content, cancelled, err := t.runFzfCore(fzfArgs, inputText)
	if err != nil {
		return "", err
	}
	if cancelled || len(content) == 0 {
		return "", nil
	}
	return strings.TrimSpace(string(content)), nil
}

func (t *Terminal) runFzfSSHSafeWithQuery(fzfArgs []string, inputText string) ([]string, error) {
	content, cancelled, err := t.runFzfCore(fzfArgs, inputText)
	if err != nil {
		return nil, err
	}
	if cancelled || len(content) == 0 {
		return []string{}, nil
	}
	return strings.Split(strings.TrimRight(string(content), "\n"), "\n"), nil
}

func (t *Terminal) IsTerminal() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func (t *Terminal) ShowHelp() {
	t.ShowLogo()

	fmt.Println("\n\033[1;94m❯ COMMANDS\033[0m")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-h, --help", "Open this dashboard")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-c, --continue", "Resume last conversation")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-a, --history", "Manage sessions (Load/Delete)")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-d [dir]", "Pack directory for AI")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-p [plat]", "Switch platform")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-m [model]", "Specify model")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-l [file]", "Load file/URL context")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-w [query]", "Search the web")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-s [url]", "Scrape web content")
	fmt.Printf("  \033[93m%-15s\033[0m %s\n", "-e, --export", "Export code blocks")

	fmt.Println("\n\033[1;94m❯ FEATURES\033[0m")
	fmt.Println("  \033[38;2;0;0;0m•\033[0m \033[1mAuto-Command:\033[0m Detects and runs shell commands with confirmation.")
	fmt.Println("  \033[38;2;0;0;0m•\033[0m \033[1mSession Mgr:\033[0m Interactive deletion and loading of past sessions.")

	fmt.Println("\n\033[1;94m❯ USAGE EXAMPLES\033[0m")
	fmt.Println("  \033[38;2;0;0;0m$\033[0m viren \"How does Docker work?\"")
	fmt.Println("  \033[38;2;0;0;0m$\033[0m viren -p groq -m llama3 \"Write a Go web server\"")
	fmt.Println("  \033[38;2;0;0;0m$\033[0m cat logs.txt | viren \"Find the error in these logs\"")

	fmt.Println("\n\033[1;94m❯ CORE PLATFORMS\033[0m")
	var platforms []string
	for name := range t.config.Platforms {
		platforms = append(platforms, name)
	}
	sort.Strings(platforms)
	for i, name := range platforms {
		fmt.Printf(" \033[96m%-12s\033[0m", name)
		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
	if len(platforms)%4 != 0 {
		fmt.Println()
	}

	fmt.Println("\n\033[38;2;0;0;0m" + strings.Repeat("━", 60) + "\033[0m")
	fmt.Println(" \033[1;92mRUN 'viren' FOR INTERACTIVE CHAT\033[0m")
	fmt.Println("\033[38;2;0;0;0m" + strings.Repeat("━", 60) + "\033[0m")
}

func (t *Terminal) RecordShellSession() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	tempDir, err := util.GetTempDir()
	if err != nil {
		return "", fmt.Errorf("failed to get temp directory: %w", err)
	}

	tempFile, err := os.CreateTemp(tempDir, "viren_shell_session_*.log")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	t.PrintInfo(fmt.Sprintf("%s session started", shell))

	var cmd *exec.Cmd

	cmd = exec.Command("script", "-q", tempFile.Name())
	cmd.Env = append(os.Environ(), "SHELL="+shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {

		if _, ok := err.(*exec.ExitError); ok {
			cmd = exec.Command("script", tempFile.Name())
			cmd.Env = append(os.Environ(), "SHELL="+shell)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {

				if _, ok := err.(*exec.ExitError); !ok {
					return "", fmt.Errorf("failed to run shell session: %w", err)
				}
			}
		} else {
			return "", fmt.Errorf("failed to run shell session: %w", err)
		}
	}

	t.PrintInfo("shell session ended")

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read session recording: %w", err)
	}

	return string(content), nil
}

func (t *Terminal) ShowHelpFzf() string {
	options := t.getInteractiveHelpOptions()

	fzfArgs := []string{
		"--reverse", "--height=40%", "--border",
		"--prompt=⚡ option: ", "--multi",
	}
	inputText := strings.Join(options, "\n")

	output, err := t.runFzfSSHSafe(fzfArgs, inputText)
	if err != nil {
		t.PrintError(fmt.Sprintf("%v", err))
		return ""
	}

	if output == "" {
		return ""
	}

	selectedItems := strings.Split(output, "\n")

	if len(selectedItems) > 1 {
		for _, item := range selectedItems {
			if !strings.HasPrefix(item, ">all") {
				return t.processHelpSelection(item, options)
			}
		}
	}

	if len(selectedItems) == 1 && strings.HasPrefix(selectedItems[0], ">all") {
		for _, option := range options {
			if !strings.HasPrefix(option, ">all") && !strings.HasPrefix(option, ">state") {
				fmt.Printf("\033[93m%s\033[0m\n", option)
			}
		}
		return ""
	}

	return t.processHelpSelection(selectedItems[0], options)
}

func (t *Terminal) processHelpSelection(selected string, options []string) string {
	if strings.HasPrefix(selected, ">state") {
		return ">state"
	}

	if strings.HasPrefix(selected, ">") {
		fmt.Printf("\033[93m%s\033[0m\n", selected)
		return ""
	}

	parts := strings.Fields(selected)
	if len(parts) > 0 {
		command := parts[0]

		if strings.HasPrefix(command, "!") || command == "cc" {
			return command
		}
	}

	fmt.Printf("\033[93m%s\033[0m\n", selected)
	return ""
}

func (t *Terminal) getCommandList() []string {
	return []string{
		fmt.Sprintf("%s - exit interface", t.config.ExitKey),
		fmt.Sprintf("%s - help dashboard", t.config.HelpKey),
		fmt.Sprintf("%s - clear chat & screen", t.config.ClearHistory),
		fmt.Sprintf("%s - backtrack messages", t.config.Backtrack),
		fmt.Sprintf("%s - switch models", t.config.ModelSwitch),
		fmt.Sprintf("%s - select from all models", t.config.AllModels),
		fmt.Sprintf("%s - switch platforms", t.config.PlatformSwitch),
		"!u - select AI personality",
		"!v - select domain mode",
		"!z - change terminal theme",
		"!onboard - setup neural profile",
		fmt.Sprintf("%s - record shell session", t.config.ShellRecord),
		fmt.Sprintf("%s - generate codedump", t.config.CodeDump),
		fmt.Sprintf("%s - add to clipboard", t.config.CopyToClipboard),
		fmt.Sprintf("%s - quick copy latest", t.config.QuickCopyLatest),
		fmt.Sprintf("%s - multi-line mode ('\\\\')", t.config.MultiLine),
		fmt.Sprintf("%s [file] - export chat(s)", t.config.ExportChat),
		fmt.Sprintf("%s [buff] - text editor mode", t.config.EditorInput),
		fmt.Sprintf("%s [dir] - load files/dirs", t.config.LoadFiles),
		fmt.Sprintf("%s [url] - scrape web content", t.config.ScrapeURL),
		fmt.Sprintf("%s [query] - perform web search", t.config.WebSearch),
		fmt.Sprintf("%s [exact] - manage sessions", t.config.AnswerSearch),
		"ctrl+c - clear prompt input",
		"ctrl+d - exit completely",
	}
}

func (t *Terminal) getInteractiveHelpOptions() []string {
	options := []string{
		">all - show all help options",
		">state - show current state",
	}
	options = append(options, t.getCommandList()...)
	return options
}

func (t *Terminal) ShowLoadingAnimation(ctx context.Context, message string) {
	if t.config.IsPipedOutput {
		return
	}

	simpleMessages := []string{
		"THINKING",
		"WAITING",
		"PROCESSING",
		"WORKING",
		"ANALYZING",
	}

	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	msgIdx := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if i%10 == 0 {
				msgIdx = (msgIdx + 1) % len(simpleMessages)
			}

			fmt.Printf("\r \033[1;34m%s\033[0m \033[38;2;0;0;0m%s\033[0m ", simpleMessages[msgIdx], chars[i%len(chars)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *Terminal) FzfSelect(items []string, prompt string) (string, error) {
	fzfArgs := []string{"--reverse", "--height=40%", "--border", "--prompt=❯ " + prompt}
	inputText := strings.Join(items, "\n")

	return t.runFzfSSHSafe(fzfArgs, inputText)
}

func (t *Terminal) FzfMultiSelect(items []string, prompt string) ([]string, error) {
	fzfArgs := []string{"--reverse", "--height=40%", "--border", "--prompt=❯ " + prompt, "--multi", "--bind=tab:toggle+down"}
	inputText := strings.Join(items, "\n")

	result, err := t.runFzfSSHSafe(fzfArgs, inputText)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []string{}, nil
	}

	return strings.Split(result, "\n"), nil
}

func (t *Terminal) FzfMultiSelectExact(items []string, prompt string) ([]string, error) {
	fzfArgs := []string{"--reverse", "--height=40%", "--border", "--prompt=❯ " + prompt, "--multi", "--bind=tab:toggle+down", "--exact"}
	inputText := strings.Join(items, "\n")

	result, err := t.runFzfSSHSafe(fzfArgs, inputText)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []string{}, nil
	}

	return strings.Split(result, "\n"), nil
}

func (t *Terminal) FzfMultiSelectForCLI(items []string, prompt string) ([]string, error) {
	fzfArgs := []string{"--reverse", "--height=40%", "--border", "--prompt=❯ " + prompt, "--multi", "--bind=tab:toggle+down"}
	inputText := strings.Join(items, "\n")

	result, err := t.runFzfSSHSafe(fzfArgs, inputText)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, fmt.Errorf("user cancelled")
	}

	return strings.Split(result, "\n"), nil
}

func (t *Terminal) FzfSelectOrQuery(items []string, prompt string) (string, error) {
	fzfArgs := []string{"--reverse", "--height=40%", "--border", "--prompt=❯ " + prompt, "--print-query"}
	inputText := strings.Join(items, "\n")

	lines, err := t.runFzfSSHSafeWithQuery(fzfArgs, inputText)
	if err != nil {
		return "", err
	}

	if len(lines) == 0 {
		return "", nil
	}

	if len(lines) > 1 && lines[1] != "" {
		return lines[1], nil
	}

	if lines[0] != "" {
		return lines[0], nil
	}

	return "", nil
}

func (t *Terminal) PrintSuccess(message string) {
	if t.config.IsPipedOutput {
		fmt.Printf("%s\n", message)
	} else {
		theme := t.GetTheme()
		fmt.Printf("%s SUCCESS \033[0m %s\n", theme.SuccessBox, message)
	}
}

func (t *Terminal) PrintError(message string) {
	if t.config.IsPipedOutput {
		fmt.Fprintf(os.Stderr, "%s\n", message)
	} else {
		theme := t.GetTheme()
		fmt.Printf("%s ERROR \033[0m %s\n", theme.ErrorBox, message)
	}
}

func (t *Terminal) PrintInfo(message string) {
	if t.config.IsPipedOutput {
		return
	}
	theme := t.GetTheme()
	fmt.Printf("%s INFO \033[0m %s\n", theme.InfoBox, message)
}

func (t *Terminal) PrintModelSwitch(model string) {
	if t.config.IsPipedOutput {
		return
	}
	theme := t.GetTheme()
	modeStr := ""
	if t.config.CurrentMode != "" && t.config.CurrentMode != "standard" {
		modeStr = fmt.Sprintf("%s/\033[0m\033[1;93m%s\033[0m", theme.BorderColor, t.config.CurrentMode)
	}
	fmt.Printf("\033[92m  ○\033[0m \033[1mMODEL UPDATED\033[0m  %s[\033[0m\033[1;96m%s\033[0m%s/\033[0m\033[1;95m%s\033[0m%s%s]\033[0m\n", theme.BorderColor, t.config.CurrentPlatform, theme.BorderColor, model, modeStr, theme.BorderColor)
}

func (t *Terminal) PrintPlatformSwitch(platform, model string) {
	if t.config.IsPipedOutput {
		return
	}
	theme := t.GetTheme()
	modeStr := ""
	if t.config.CurrentMode != "" && t.config.CurrentMode != "standard" {
		modeStr = fmt.Sprintf("\033[38;2;0;0;0m/\033[0m\033[93m%s\033[0m", t.config.CurrentMode)
	}
	fmt.Printf("\033[92m  ○\033[0m \033[1mPLATFORM UPDATED\033[0m  %s[\033[0m\033[1;96m%s\033[0m%s/\033[0m\033[1;95m%s\033[0m%s%s]\033[0m\n", theme.BorderColor, platform, theme.BorderColor, model, modeStr, theme.BorderColor)
}

func (t *Terminal) LoadFileContent(selections []string) (string, error) {
	var contentBuilder strings.Builder

	for _, selection := range selections {
		if selection == "" {
			continue
		}

		if t.isURL(selection) {
			urlContent, err := t.scrapeURL(selection)
			if err != nil {
				contentBuilder.WriteString(fmt.Sprintf("Error scraping %s: %v\n", selection, err))
				continue
			}
			contentBuilder.WriteString(urlContent)
			continue
		}

		info, err := os.Stat(selection)
		if err != nil {
			continue
		}

		if info.IsDir() {
			dirContent, err := t.loadDirectoryContent(selection)
			if err != nil {
				continue
			}
			contentBuilder.WriteString(dirContent)
		} else {
			fileContent, err := t.loadTextFile(selection)
			if err != nil {
				continue
			}
			contentBuilder.WriteString(fileContent)
		}
	}

	return contentBuilder.String(), nil
}

func (t *Terminal) loadTextFile(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	var content string
	var err error

	switch ext {
	case ".pdf":
		content, err = t.loadPDF(filePath)
	case ".docx", ".odt", ".rtf":
		content, err = t.loadDOCX(filePath)
	case ".xlsx":
		content, err = t.loadXLSX(filePath)
	case ".csv":
		content, err = t.loadCSV(filePath)
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp":
		content, err = t.loadImage(filePath)
	default:

		fileContent, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return "", readErr
		}

		if !t.isTextFile(fileContent) {
			return "", fmt.Errorf("file is not a supported file type")
		}

		content = string(fileContent)
	}

	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("File: %s\n", filePath))
	result.WriteString(content)
	result.WriteString("\n\n")

	return result.String(), nil
}

func (t *Terminal) loadDirectoryContent(dirPath string) (string, error) {
	var result strings.Builder

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		fileContent, err := t.loadTextFile(path)
		if err != nil {
			return nil
		}

		result.WriteString(fileContent)
		return nil
	})

	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (t *Terminal) loadPDF(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file stats: %w", err)
	}

	reader, err := pdf.NewReader(file, stat.Size())
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var content strings.Builder
	numPages := reader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		content.WriteString(text)
		content.WriteString("\n")
	}

	return content.String(), nil
}

func (t *Terminal) loadDOCX(filePath string) (string, error) {
	text, err := cat.File(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from document: %w", err)
	}
	return text, nil
}

func (t *Terminal) loadXLSX(filePath string) (string, error) {
	workbook, err := xlsx.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open XLSX file: %w", err)
	}

	var content strings.Builder

	for _, sheet := range workbook.Sheets {
		content.WriteString(fmt.Sprintf("=== Sheet: %s ===\n", sheet.Name))

		err := sheet.ForEachRow(func(row *xlsx.Row) error {
			var rowData []string
			err := row.ForEachCell(func(cell *xlsx.Cell) error {
				text := cell.String()
				rowData = append(rowData, text)
				return nil
			})
			if err != nil {
				return err
			}

			if len(strings.TrimSpace(strings.Join(rowData, ""))) > 0 {
				content.WriteString(fmt.Sprintf("%s\n", strings.Join(rowData, " | ")))
			}
			return nil
		})

		if err != nil {
			return "", fmt.Errorf("failed to read sheet %s: %w", sheet.Name, err)
		}
		content.WriteString("\n")
	}

	return content.String(), nil
}

func (t *Terminal) loadCSV(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV file: %w", err)
	}

	var content strings.Builder

	for rowIndex, record := range records {
		content.WriteString(fmt.Sprintf("Row %d: %s\n", rowIndex+1, strings.Join(record, " | ")))
	}

	return content.String(), nil
}

func (t *Terminal) loadImage(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	var content strings.Builder
	content.WriteString(fmt.Sprintf("image analysis for: %s\n\n", filepath.Base(filePath)))

	fileInfo, err := file.Stat()
	if err == nil {
		content.WriteString(fmt.Sprintf("file size: %d bytes (%.2f KB)\n", fileInfo.Size(), float64(fileInfo.Size())/1024.0))
		content.WriteString(fmt.Sprintf("modified: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05")))
	}

	file.Seek(0, 0)

	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	content.WriteString(fmt.Sprintf("format: %s\n", strings.ToUpper(format)))
	content.WriteString(fmt.Sprintf("dimensions: %dx%d pixels\n", bounds.Dx(), bounds.Dy()))

	file.Seek(0, 0)

	exifData, err := exif.Decode(file)
	if err == nil {
		content.WriteString("\nEXIF metadata:\n")

		exifTags := []struct {
			name	string
			tag	exif.FieldName
		}{
			{"camera make", exif.Make},
			{"camera model", exif.Model},
			{"date time", exif.DateTime},
			{"date time original", exif.DateTimeOriginal},
			{"date time digitized", exif.DateTimeDigitized},
			{"software", exif.Software},
			{"artist", exif.Artist},
			{"copyright", exif.Copyright},
			{"image description", exif.ImageDescription},
			{"user comment", exif.UserComment},
			{"orientation", exif.Orientation},
			{"x resolution", exif.XResolution},
			{"y resolution", exif.YResolution},
			{"resolution unit", exif.ResolutionUnit},
			{"flash", exif.Flash},
			{"focal length", exif.FocalLength},
			{"exposure time", exif.ExposureTime},
			{"f number", exif.FNumber},
			{"iso", exif.ISOSpeedRatings},
			{"white balance", exif.WhiteBalance},
			{"gps latitude", exif.GPSLatitude},
			{"gps longitude", exif.GPSLongitude},
			{"gps altitude", exif.GPSAltitude},
		}

		for _, tagInfo := range exifTags {
			if tag, err := exifData.Get(tagInfo.tag); err == nil {
				value := strings.TrimSpace(tag.String())
				if value != "" && value != "0" && value != "0/1" {
					content.WriteString(fmt.Sprintf("  %s: %s\n", tagInfo.name, value))
				}
			}
		}

		if lat, err := exifData.Get(exif.GPSLatitude); err == nil {
			if latRef, err := exifData.Get(exif.GPSLatitudeRef); err == nil {
				if lon, err := exifData.Get(exif.GPSLongitude); err == nil {
					if lonRef, err := exifData.Get(exif.GPSLongitudeRef); err == nil {
						latDeg := convertDMSToDecimal(lat.String())
						lonDeg := convertDMSToDecimal(lon.String())
						if latRef.String() == "S" {
							latDeg = -latDeg
						}
						if lonRef.String() == "W" {
							lonDeg = -lonDeg
						}
						if latDeg != 0 || lonDeg != 0 {
							content.WriteString(fmt.Sprintf("  gps coordinates: %.6f, %.6f\n", latDeg, lonDeg))
						}
					}
				}
			}
		}
	} else {
		content.WriteString("\nno EXIF metadata found or failed to read EXIF data\n")
	}

	content.WriteString(fmt.Sprintf("\nimage properties:\n"))
	content.WriteString(fmt.Sprintf("  color mode: %T\n", img.ColorModel()))
	content.WriteString(fmt.Sprintf("  aspect ratio: %.2f:1\n", float64(bounds.Dx())/float64(bounds.Dy())))

	megapixels := float64(bounds.Dx()*bounds.Dy()) / 1000000.0
	if megapixels > 1.0 {
		content.WriteString(fmt.Sprintf("  megapixels: %.1f MP\n", megapixels))
	} else {
		content.WriteString(fmt.Sprintf("  resolution: %.0f K pixels\n", megapixels*1000))
	}

	content.WriteString("\n" + strings.Repeat("=", 50) + "\n")
	content.WriteString("text extraction (OCR):\n")
	content.WriteString(strings.Repeat("=", 50) + "\n\n")

	extractedText, err := t.extractTextFromImage(filePath)
	if err != nil {
		content.WriteString(fmt.Sprintf("OCR error: %v\n", err))
	} else if strings.TrimSpace(extractedText) == "" {
		content.WriteString("no text detected in the image.\n")
	} else {
		content.WriteString("extracted text:\n")
		content.WriteString(strings.Repeat("-", 30) + "\n")
		content.WriteString(extractedText)
		content.WriteString("\n" + strings.Repeat("-", 30) + "\n")
	}

	return content.String(), nil
}

func convertDMSToDecimal(dms string) float64 {

	parts := strings.Split(dms, ",")
	if len(parts) < 3 {
		return 0
	}

	var degrees, minutes, seconds float64

	if d := parseFraction(strings.TrimSpace(parts[0])); d >= 0 {
		degrees = d
	}
	if m := parseFraction(strings.TrimSpace(parts[1])); m >= 0 {
		minutes = m
	}
	if s := parseFraction(strings.TrimSpace(parts[2])); s >= 0 {
		seconds = s
	}

	return degrees + minutes/60.0 + seconds/3600.0
}

func parseFraction(fraction string) float64 {
	parts := strings.Split(fraction, "/")
	if len(parts) != 2 {
		if f, err := strconv.ParseFloat(fraction, 64); err == nil {
			return f
		}
		return -1
	}

	numerator, err1 := strconv.ParseFloat(parts[0], 64)
	denominator, err2 := strconv.ParseFloat(parts[1], 64)

	if err1 != nil || err2 != nil || denominator == 0 {
		return -1
	}

	return numerator / denominator
}

func (t *Terminal) isTextFile(content []byte) bool {
	if len(content) == 0 {
		return true
	}

	for _, b := range content {
		if b == 0 {
			return false
		}
	}

	printableCount := 0
	for _, b := range content {
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 || b >= 128 {
			printableCount++
		}
	}

	return float64(printableCount)/float64(len(content)) > 0.95
}

func (t *Terminal) GetCurrentDirFilesRecursive() ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %v", err)
	}
	return t.GetDirFilesRecursive(currentDir)
}

func (t *Terminal) GetDirFilesRecursive(targetDir string) ([]string, error) {
	var items []string

	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %v", err)
	}

	isShallow := util.IsShallowLoadDir(t.config, absTargetDir)

	err = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if path == targetDir {
			return nil
		}

		relPath, err := filepath.Rel(targetDir, path)
		if err != nil {
			return nil
		}

		if isShallow {
			depth := strings.Count(relPath, string(filepath.Separator))
			if d.IsDir() && depth > 0 {

				return filepath.SkipDir
			}
			if !d.IsDir() && depth > 0 {

				return nil
			}
		}

		if d.IsDir() && (filepath.Base(relPath) == ".git" || filepath.Base(relPath) == ".svn" || filepath.Base(relPath) == ".hg") {
			return filepath.SkipDir
		}

		if d.IsDir() {
			items = append(items, relPath+"/")
		} else {
			items = append(items, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	return items, nil
}

func (t *Terminal) CodeDump() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	return t.CodeDumpFromDir(pwd)
}

func (t *Terminal) CodeDumpFromDir(targetDir string) (string, error) {

	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	allFiles, err := t.discoverFiles(absDir)
	if err != nil {
		return "", fmt.Errorf("failed to discover files: %v", err)
	}

	if len(allFiles) == 0 {
		return "", fmt.Errorf("no text files found in directory")
	}

	fzfOptions := append([]string{">none"}, allFiles...)

	excludedItems, err := t.FzfMultiSelect(fzfOptions, "exclude from dump (tab=multi): ")
	if err != nil {
		return "", fmt.Errorf("failed to get exclusions: %v", err)
	}

	var filteredExclusions []string
	for _, item := range excludedItems {
		if !strings.HasPrefix(item, ">none") {
			filteredExclusions = append(filteredExclusions, item)
		}
	}
	excludedItems = filteredExclusions

	includedFiles := t.filterExcludedFiles(allFiles, excludedItems)

	if len(includedFiles) == 0 {
		return "", fmt.Errorf("no files remaining after exclusions")
	}

	return t.generateCodeDumpFromDir(includedFiles, absDir)
}

func (t *Terminal) CodeDumpFromDirForCLI(targetDir string) (string, error) {

	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	allFiles, err := t.discoverFiles(absDir)
	if err != nil {
		return "", fmt.Errorf("failed to discover files: %v", err)
	}

	if len(allFiles) == 0 {
		return "", fmt.Errorf("no text files found in directory")
	}

	fzfOptions := append([]string{">none"}, allFiles...)

	excludedItems, err := t.FzfMultiSelectForCLI(fzfOptions, "exclude from dump (tab=multi): ")
	if err != nil {
		return "", fmt.Errorf("failed to get exclusions: %v", err)
	}

	var filteredExclusions []string
	for _, item := range excludedItems {
		if !strings.HasPrefix(item, ">none") {
			filteredExclusions = append(filteredExclusions, item)
		}
	}
	excludedItems = filteredExclusions

	includedFiles := t.filterExcludedFiles(allFiles, excludedItems)

	if len(includedFiles) == 0 {
		return "", fmt.Errorf("no files remaining after exclusions")
	}

	return t.generateCodeDumpFromDir(includedFiles, absDir)
}

func (t *Terminal) discoverFiles(rootDir string) ([]string, error) {
	var allFiles []string
	var allDirs []string
	gitignorePatterns := t.loadGitignorePatterns(rootDir)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return nil
		}

		if relPath == "." {
			return nil
		}

		if t.shouldIgnore(relPath, gitignorePatterns) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			allDirs = append(allDirs, relPath+"/")
		} else {

			if t.isTextFileByPath(path) {
				allFiles = append(allFiles, relPath)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var combined []string
	combined = append(combined, allDirs...)
	combined = append(combined, allFiles...)

	return combined, nil
}

func (t *Terminal) loadGitignorePatterns(rootDir string) []string {

	patterns := []string{".git/"}

	gitignorePath := filepath.Join(rootDir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return patterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return patterns
}

func (t *Terminal) shouldIgnore(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if t.matchesPattern(path, pattern) {
			return true
		}
	}
	return false
}

func (t *Terminal) matchesPattern(path, pattern string) bool {

	pattern = strings.TrimPrefix(pattern, "/")

	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")

		return strings.HasPrefix(path, dirPattern+"/") || path == dirPattern
	}

	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}

		parts := strings.Split(path, "/")
		for i := range parts {
			partialPath := strings.Join(parts[:i+1], "/")
			if matched, _ := filepath.Match(pattern, partialPath); matched {
				return true
			}
		}
		return false
	}

	return path == pattern || strings.HasPrefix(path, pattern+"/")
}

func (t *Terminal) isTextFileByPath(filePath string) bool {

	ext := strings.ToLower(filepath.Ext(filePath))
	supportedExtensions := map[string]bool{

		".txt":	true, ".md": true, ".go": true, ".py": true, ".js": true,
		".ts":	true, ".jsx": true, ".tsx": true, ".html": true, ".css": true,
		".scss":	true, ".sass": true, ".json": true, ".xml": true, ".yaml": true,
		".yml":	true, ".toml": true, ".ini": true, ".cfg": true, ".conf": true,
		".sh":	true, ".bash": true, ".zsh": true, ".fish": true, ".ps1": true,
		".bat":	true, ".cmd": true, ".dockerfile": true, ".makefile": true,
		".c":	true, ".cpp": true, ".cc": true, ".cxx": true, ".h": true,
		".hpp":	true, ".java": true, ".kt": true, ".scala": true, ".rb": true,
		".php":	true, ".pl": true, ".pm": true, ".r": true, ".sql": true,
		".vim":	true, ".lua": true, ".rs": true, ".swift": true, ".m": true,
		".mm":	true, ".cs": true, ".vb": true, ".fs": true, ".clj": true,

		".pdf":	true, ".docx": true, ".odt": true, ".rtf": true, ".xlsx": true, ".csv": true,

		".jpg":	true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".tiff": true, ".tif": true, ".webp": true,
		".hs":	true, ".elm": true, ".ex": true, ".exs": true, ".erl": true,
		".hrl":	true, ".dart": true, ".gradle": true, ".sbt": true,
		".build":	true, ".cmake": true, ".mk": true, ".am": true, ".in": true,
		".ac":	true, ".m4": true, ".spec": true, ".desktop": true, ".service": true,
		".log":	true, ".tsv": true, ".properties": true, ".env": true,
	}

	if supportedExtensions[ext] {
		return true
	}

	if ext == "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return false
		}
		return t.isTextFile(content)
	}

	return false
}

func (t *Terminal) filterExcludedFiles(allFiles, excludedItems []string) []string {
	excludedSet := make(map[string]bool)
	var excludedDirs []string

	for _, item := range excludedItems {
		excludedSet[item] = true
		if strings.HasSuffix(item, "/") {
			excludedDirs = append(excludedDirs, strings.TrimSuffix(item, "/"))
		}
	}

	var includedFiles []string
	for _, file := range allFiles {

		if excludedSet[file] {
			continue
		}

		if strings.HasSuffix(file, "/") {
			continue
		}

		excluded := false
		for _, excludedDir := range excludedDirs {
			if strings.HasPrefix(file, excludedDir+"/") {
				excluded = true
				break
			}
		}

		if !excluded {
			includedFiles = append(includedFiles, file)
		}
	}

	return includedFiles
}

func (t *Terminal) generateCodeDumpFromDir(files []string, sourceDir string) (string, error) {
	var result strings.Builder

	result.WriteString("=== Code Dump ===\n\n")
	result.WriteString(fmt.Sprintf("generated from directory: %s\n", sourceDir))
	result.WriteString(fmt.Sprintf("total files: %d\n\n", len(files)))

	for _, file := range files {

		fullPath := filepath.Join(sourceDir, file)

		ext := strings.ToLower(filepath.Ext(file))

		supportedTypes := []string{".pdf", ".docx", ".odt", ".rtf", ".xlsx", ".csv"}
		isSpecialFile := false
		for _, supportedExt := range supportedTypes {
			if ext == supportedExt {
				isSpecialFile = true
				break
			}
		}

		var content string
		if isSpecialFile {

			fileContent, err := t.loadTextFile(fullPath)
			if err != nil {
				result.WriteString(fmt.Sprintf("=== FILE: %s ===\nError processing file: %v\n\n", file, err))
				continue
			}
			content = fileContent
		} else {

			fileBytes, err := os.ReadFile(fullPath)
			if err != nil {
				result.WriteString(fmt.Sprintf("=== FILE: %s ===\nError reading file: %v\n\n", file, err))
				continue
			}

			if !t.isTextFile(fileBytes) {
				continue
			}

			content = fmt.Sprintf("File: %s\n%s", file, string(fileBytes))
		}

		result.WriteString(fmt.Sprintf("=== FILE: %s ===\n", file))
		result.WriteString(content)
		result.WriteString("\n\n")
	}

	result.WriteString("=== END CODE DUMP ===")
	return result.String(), nil
}

func (t *Terminal) isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (t *Terminal) IsURL(str string) bool {
	return t.isURL(str)
}

func (t *Terminal) isYouTubeURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	host := strings.ToLower(u.Host)
	return strings.Contains(host, "youtube.com") ||
		strings.Contains(host, "youtu.be") ||
		strings.Contains(host, "m.youtube.com") ||
		strings.Contains(host, "www.youtube.com") ||
		strings.Contains(host, "youtube-nocookie.com")
}

func (t *Terminal) cleanURL(urlStr string) string {

	cleaned := strings.ReplaceAll(urlStr, `\?`, `?`)
	cleaned = strings.ReplaceAll(cleaned, `\=`, `=`)
	cleaned = strings.ReplaceAll(cleaned, `\&`, `&`)
	return cleaned
}

func (t *Terminal) scrapeURL(urlStr string) (string, error) {

	ctx, cancel := context.WithCancel(context.Background())
	go t.ShowLoadingAnimation(ctx, "Scraping")
	defer cancel()

	return t.scrapeURLInternal(urlStr)
}

func (t *Terminal) scrapeURLInternal(urlStr string) (string, error) {

	cleanedURL := t.cleanURL(urlStr)

	var result strings.Builder
	result.WriteString(fmt.Sprintf("=== %s ===\n\n", cleanedURL))

	var scrapeErr error
	if t.isYouTubeURL(cleanedURL) {

		content, err := t.scrapeYouTube(cleanedURL)
		if err != nil {
			scrapeErr = fmt.Errorf("failed to scrape YouTube URL: %w", err)
		} else {
			result.WriteString(content)
		}
	} else {

		content, err := t.scrapeWeb(cleanedURL)
		if err != nil {
			scrapeErr = fmt.Errorf("failed to scrape URL: %w", err)
		} else {
			result.WriteString(content)
		}
	}

	if scrapeErr != nil {
		return "", scrapeErr
	}

	result.WriteString("\n")
	return result.String(), nil
}

func (t *Terminal) scrapeWeb(urlStr string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Virenrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL: status code %d", resp.StatusCode)
	}

	return t.textContentFromHTML(resp.Body)
}

func (t *Terminal) textContentFromHTML(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var sb strings.Builder
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {

		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "nav", "header", "footer", "aside":
				return
			}
		}

		if n.Type == html.TextNode {

			trimmed := strings.TrimSpace(n.Data)
			if len(trimmed) > 0 {
				sb.WriteString(trimmed)
				sb.WriteString(" ")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "div", "h1", "h2", "h3", "h4", "h5", "h6", "li", "br", "tr":
				sb.WriteString("\n")
			}
		}
	}

	traverse(doc)

	text := sb.String()
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`(\s*\n\s*)+`).ReplaceAllString(text, "\n")

	return strings.TrimSpace(text), nil
}

func (t *Terminal) scrapeYouTube(urlStr string) (string, error) {
	var result strings.Builder

	result.WriteString("--- Metadata ---\n")
	metadataCtx, metadataCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer metadataCancel()
	metadataCmd := exec.CommandContext(metadataCtx, "yt-dlp", "-j", urlStr)
	metadataOutput, err := metadataCmd.Output()
	if err != nil {
		if metadataCtx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("failed to get YouTube metadata: command timed out")
		}
		return "", fmt.Errorf("failed to get YouTube metadata: %w", err)
	}

	metadata := string(metadataOutput)
	result.WriteString(t.parseYouTubeMetadata(metadata))

	result.WriteString("\n--- Subtitles ---\n")

	tempDir, err := util.GetTempDir()
	if err != nil {
		return "", fmt.Errorf("failed to get temp directory: %w", err)
	}

	baseName := filepath.Join(tempDir, fmt.Sprintf("yt_%d", time.Now().UnixNano()))

	subtitleCtx, subtitleCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer subtitleCancel()
	subtitleCmd := exec.CommandContext(subtitleCtx, "yt-dlp", "--quiet", "--skip-download",
		"--write-auto-subs", "--sub-lang", "en", "--sub-format", "srt",
		"-o", baseName+".%(ext)s", urlStr)

	err = subtitleCmd.Run()
	if err == nil {

		pattern := baseName + "*.srt"
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			srtContent, readErr := os.ReadFile(matches[0])
			if readErr == nil {
				result.WriteString(string(srtContent))
			}

			for _, match := range matches {
				os.Remove(match)
			}
		}
	}

	return result.String(), nil
}

func (t *Terminal) parseYouTubeMetadata(jsonStr string) string {
	var result strings.Builder

	extractField := func(field string) string {
		pattern := fmt.Sprintf(`"%s":\s*"([^"]*)"`, field)
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(jsonStr)
		if len(matches) > 1 {
			return matches[1]
		}

		pattern = fmt.Sprintf(`"%s":\s*(\d+)`, field)
		re = regexp.MustCompile(pattern)
		matches = re.FindStringSubmatch(jsonStr)
		if len(matches) > 1 {
			return matches[1]
		}

		return ""
	}

	if title := extractField("title"); title != "" {
		result.WriteString(fmt.Sprintf("title: %s\n", title))
	}
	if duration := extractField("duration"); duration != "" {
		result.WriteString(fmt.Sprintf("duration: %s seconds\n", duration))
	}
	if viewCount := extractField("view_count"); viewCount != "" {
		result.WriteString(fmt.Sprintf("view count: %s\n", viewCount))
	}
	if uploader := extractField("uploader"); uploader != "" {
		result.WriteString(fmt.Sprintf("uploader: %s\n", uploader))
	}
	if uploadDate := extractField("upload_date"); uploadDate != "" {
		result.WriteString(fmt.Sprintf("upload date: %s\n", uploadDate))
	}
	if description := extractField("description"); description != "" {

		if len(description) > 500 {
			description = description[:500] + "..."
		}
		result.WriteString(fmt.Sprintf("description: %s\n", description))
	}

	return result.String()
}

func (t *Terminal) ScrapeURLs(urls []string) (string, error) {
	var result strings.Builder

	for _, urlStr := range urls {
		if urlStr == "" {
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		go t.ShowLoadingAnimation(ctx, "Scraping")

		content, err := t.scrapeURLInternal(urlStr)

		cancel()

		if err != nil {
			result.WriteString(fmt.Sprintf("Error scraping %s: %v\n", urlStr, err))
			continue
		}

		result.WriteString(content)
	}

	return result.String(), nil
}

func (t *Terminal) WebSearch(query string) (string, error) {
	apiKey := os.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("the BRAVE_API_KEY environment variable is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go t.ShowLoadingAnimation(ctx, "Searching")
	defer cancel()

	req, err := http.NewRequest("GET", "https://api.search.brave.com/res/v1/web/search", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create search request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("count", fmt.Sprintf("%d", t.config.NumSearchResults))
	q.Add("country", t.config.SearchCountry)
	q.Add("search_lang", t.config.SearchLang)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-Subscription-Token", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read search response: %w", err)
	}

	var braveResult BraveSearchResult
	if err := json.Unmarshal(body, &braveResult); err != nil {
		return "", fmt.Errorf("failed to parse search results: %w", err)
	}

	if len(braveResult.Web.Results) == 0 {
		noResultsMsg := fmt.Sprintf("No search results found for: %s\n", query)
		if t.config.ShowSearchResults {
			fmt.Print(noResultsMsg)
		}
		return noResultsMsg, nil
	}

	formatted := t.formatBraveSearchResults(braveResult.Web.Results, query)

	if t.config.ShowSearchResults {
		fmt.Print(formatted)
	}

	return formatted, nil
}

type BraveSearchResult struct {
	Web struct {
		Results []BraveWebResult `json:"results"`
	} `json:"web"`
}

type BraveWebResult struct {
	Title		string	`json:"title"`
	URL		string	`json:"url"`
	Description	string	`json:"description"`
}

func (t *Terminal) formatBraveSearchResults(results []BraveWebResult, query string) string {
	var result strings.Builder

	if len(results) == 0 {
		result.WriteString(fmt.Sprintf("no search results found for: %s\n", query))
		return result.String()
	}

	for i, searchResult := range results {
		if searchResult.Title != "" && searchResult.URL != "" {
			if t.config.IsPipedOutput {
				result.WriteString(fmt.Sprintf("%d) %s\n", i+1, searchResult.Title))
				result.WriteString(fmt.Sprintf("%s\n", searchResult.URL))
				if searchResult.Description != "" {
					result.WriteString(fmt.Sprintf("%s\n", searchResult.Description))
				}
			} else {
				result.WriteString(fmt.Sprintf("\033[93m%d) \033[93m%s\033[0m\n", i+1, searchResult.Title))
				result.WriteString(fmt.Sprintf("\033[95m%s\033[0m\n", searchResult.URL))
				if searchResult.Description != "" {
					result.WriteString(fmt.Sprintf("\033[92m%s\033[0m\n", searchResult.Description))
				}
			}
		}
	}

	return result.String()
}

func (t *Terminal) CopyToClipboard(content string) error {
	var cmd *exec.Cmd

	if _, err := exec.LookPath("pbcopy"); err == nil {

		cmd = exec.Command("pbcopy")
	} else if _, err := exec.LookPath("xclip"); err == nil {

		cmd = exec.Command("xclip", "-selection", "clipboard")
	} else if _, err := exec.LookPath("xsel"); err == nil {

		cmd = exec.Command("xsel", "--clipboard", "--input")
	} else if _, err := exec.LookPath("wl-copy"); err == nil {

		cmd = exec.Command("wl-copy")
	} else if _, err := exec.LookPath("termux-clipboard-set"); err == nil {

		cmd = exec.Command("termux-clipboard-set")
	} else if _, err := exec.LookPath("clip"); err == nil {

		cmd = exec.Command("clip")
	} else {
		return fmt.Errorf("no clipboard utility found. Please install: pbcopy (macOS), xclip/xsel (Linux), wl-copy (Wayland), or termux-clipboard-set (Android)")
	}

	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}

func (t *Terminal) CopyResponsesInteractive(chatHistory []types.ChatHistory, messages []types.ChatMessage) error {
	if len(chatHistory) == 0 {
		return fmt.Errorf("no chat history available")
	}

	copyMode, err := t.FzfSelect([]string{"turn copy", "block copy", "manual copy", "link copy"}, "select copy mode: ")
	if err != nil {
		return fmt.Errorf("selection cancelled or failed: %v", err)
	}

	if copyMode == "turn copy" {
		return t.copyResponsesTurn(chatHistory)
	}

	if copyMode == "block copy" {
		return t.copyResponsesBlock(chatHistory)
	}

	if copyMode == "link copy" {
		return t.copyResponsesLinks(chatHistory, messages)
	}

	if copyMode == "" {
		return nil
	}

	return t.copyResponsesManual(chatHistory)
}

func (t *Terminal) CopyLatestResponseToClipboard(chatHistory []types.ChatHistory) error {
	if len(chatHistory) < 2 {
		return fmt.Errorf("no bot responses available")
	}

	latestResponse := chatHistory[len(chatHistory)-1].Bot
	if latestResponse == "" {
		return fmt.Errorf("latest response is empty")
	}

	return t.CopyToClipboard(latestResponse)
}

func (t *Terminal) copyResponsesTurn(chatHistory []types.ChatHistory) error {

	var items []string
	type replyEntry struct {
		content	string
		isUser	bool
		index	int
	}
	var entries []replyEntry

	for i := len(chatHistory) - 1; i >= 1; i-- {
		entry := chatHistory[i]

		if entry.Bot != "" {
			preview := strings.Split(entry.Bot, "\n")[0]
			if len(preview) > 70 {
				preview = preview[:70] + "..."
			}
			items = append(items, fmt.Sprintf("BOT: %s", preview))
			entries = append(entries, replyEntry{content: entry.Bot, isUser: false, index: i})
		}

		if entry.User != "" {
			preview := strings.Split(entry.User, "\n")[0]
			if len(preview) > 70 {
				preview = preview[:70] + "..."
			}
			items = append(items, fmt.Sprintf("USER: %s", preview))
			entries = append(entries, replyEntry{content: entry.User, isUser: true, index: i})
		}
	}

	if len(items) == 0 {
		return fmt.Errorf("no turns available")
	}

	fzfOptions := append([]string{">all"}, items...)

	selectedItems, err := t.FzfMultiSelect(fzfOptions, "select turns to copy (tab=multi): ")
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}

	if len(selectedItems) == 0 {
		t.PrintInfo("no turns selected")
		return nil
	}

	allSelected := ContainsAllOption(selectedItems)

	var combinedContent strings.Builder
	isSingleItem := (allSelected && len(entries) == 1) || (!allSelected && len(selectedItems) == 1)

	if allSelected {

		for i := len(entries) - 1; i >= 0; i-- {
			if i < len(entries)-1 {
				combinedContent.WriteString("\n\n")
			}

			if !isSingleItem {
				if entries[i].isUser {
					combinedContent.WriteString("USER:\n")
				} else {
					combinedContent.WriteString("BOT:\n")
				}
			}
			combinedContent.WriteString(entries[i].content)
		}
	} else {

		var selectedEntries []replyEntry
		for _, selected := range selectedItems {
			for j, item := range items {
				if item == selected {
					selectedEntries = append(selectedEntries, entries[j])
					break
				}
			}
		}

		sort.Slice(selectedEntries, func(i, j int) bool {
			if selectedEntries[i].index != selectedEntries[j].index {
				return selectedEntries[i].index < selectedEntries[j].index
			}
			return selectedEntries[i].isUser && !selectedEntries[j].isUser
		})

		for i, entry := range selectedEntries {
			if i > 0 {
				combinedContent.WriteString("\n\n")
			}

			if !isSingleItem {
				if entry.isUser {
					combinedContent.WriteString("USER:\n")
				} else {
					combinedContent.WriteString("BOT:\n")
				}
			}
			combinedContent.WriteString(entry.content)
		}
	}

	finalContent := combinedContent.String()

	err = t.CopyToClipboard(finalContent)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	count := len(selectedItems)
	if allSelected {
		count = len(entries)
	}
	turnWord := "turn"
	if count > 1 {
		turnWord = "turns"
	}
	fmt.Printf("\033[93madded %d %s to clipboard\033[0m\n", count, turnWord)
	return nil
}

func (t *Terminal) copyResponsesBlock(chatHistory []types.ChatHistory) error {

	type ExtractedBlock struct {
		Content		string
		Language	string
		Preview		string
	}
	var blocks []ExtractedBlock
	var items []string

	for i := len(chatHistory) - 1; i >= 1; i-- {
		entry := chatHistory[i]
		if entry.Bot != "" {
			matches := codeBlockRegex.FindAllStringSubmatch(entry.Bot, -1)
			for _, match := range matches {
				language := match[1]
				if language == "" {
					language = "text"
				}
				content := match[2]

				preview := strings.Split(content, "\n")[0]
				if len(preview) > 60 {
					preview = preview[:60] + "..."
				}

				displayText := fmt.Sprintf("[%s] %s", language, preview)
				items = append(items, displayText)
				blocks = append(blocks, ExtractedBlock{Content: content, Language: language, Preview: preview})
			}
		}
	}

	if len(blocks) == 0 {
		return fmt.Errorf("no code blocks found in chat history")
	}

	fzfOptions := append([]string{">all"}, items...)

	selectedItems, err := t.FzfMultiSelect(fzfOptions, "select code blocks to copy (tab=multi): ")
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}

	if len(selectedItems) == 0 {
		t.PrintInfo("no blocks selected")
		return nil
	}

	allSelected := ContainsAllOption(selectedItems)

	var combinedContent strings.Builder
	if allSelected {

		for i, block := range blocks {
			if i > 0 {
				combinedContent.WriteString("\n\n")
			}
			combinedContent.WriteString(block.Content)
		}
	} else {
		for i, selected := range selectedItems {

			for j, item := range items {
				if item == selected {
					if i > 0 {
						combinedContent.WriteString("\n\n")
					}
					combinedContent.WriteString(blocks[j].Content)
					break
				}
			}
		}
	}

	finalContent := combinedContent.String()

	err = t.CopyToClipboard(finalContent)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	count := len(selectedItems)
	if allSelected {
		count = len(blocks)
	}
	blockWord := "code block"
	if count > 1 {
		blockWord = "code blocks"
	}
	fmt.Printf("\033[93madded %d %s to clipboard\033[0m\n", count, blockWord)
	return nil
}

func (t *Terminal) copyResponsesManual(chatHistory []types.ChatHistory) error {

	var responseOptions []string
	var responseMap = make(map[string]types.ChatHistory)

	for i := len(chatHistory) - 1; i >= 1; i-- {
		entry := chatHistory[i]
		if entry.User != "" || entry.Bot != "" {

			userPreview := strings.Split(entry.User, "\n")[0]
			if len(userPreview) > 60 {
				userPreview = userPreview[:60] + "..."
			}

			timestamp := time.Unix(entry.Time, 0).Format("2006-01-02 15:04:05")
			optionText := fmt.Sprintf("%d: %s - %s", i, timestamp, userPreview)
			responseOptions = append(responseOptions, optionText)
			responseMap[optionText] = entry
		}
	}

	if len(responseOptions) == 0 {
		return fmt.Errorf("no responses found in chat history")
	}

	fzfOptions := append([]string{">all"}, responseOptions...)

	selected, err := t.FzfMultiSelect(fzfOptions, "select responses to copy (tab=multi): ")
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}

	if len(selected) == 0 {
		t.PrintInfo("no responses selected")
		return nil
	}

	allSelected := ContainsAllOption(selected)

	var combinedContent strings.Builder
	if allSelected {

		for i, option := range responseOptions {
			if entry, exists := responseMap[option]; exists {
				if i > 0 {
					combinedContent.WriteString("\n\n---\n\n")
				}
				combinedContent.WriteString(entry.Bot)
			}
		}
	} else {
		for i, selection := range selected {
			if entry, exists := responseMap[selection]; exists {
				if i > 0 {
					combinedContent.WriteString("\n\n---\n\n")
				}
				combinedContent.WriteString(entry.Bot)
			}
		}
	}

	finalContent := combinedContent.String()

	editedContent, err := t.openEditorWithContent(finalContent)
	if err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	err = t.CopyToClipboard(editedContent)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	count := len(selected)
	if allSelected {
		count = len(responseOptions)
	}
	responseWord := "response"
	if count > 1 {
		responseWord = "responses"
	}
	fmt.Printf("\033[93madded %d %s to clipboard\033[0m\n", count, responseWord)
	return nil
}

func (t *Terminal) copyResponsesLinks(chatHistory []types.ChatHistory, messages []types.ChatMessage) error {

	urls := t.ExtractURLsFromChatHistory(chatHistory)

	if len(messages) > 0 {
		messagesURLs := t.ExtractURLsFromMessages(messages)

		seen := make(map[string]bool)
		for _, url := range urls {
			seen[url] = true
		}
		for _, url := range messagesURLs {
			if !seen[url] {
				urls = append(urls, url)
				seen[url] = true
			}
		}
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs found")
	}

	fzfOptions := append([]string{">all"}, urls...)

	selectedURLs, err := t.FzfMultiSelect(fzfOptions, "select URLs to copy (tab=multi): ")
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}

	if len(selectedURLs) == 0 {
		t.PrintInfo("no URLs selected")
		return nil
	}

	allSelected := ContainsAllOption(selectedURLs)

	var finalContent string
	if allSelected {
		finalContent = strings.Join(urls, " ")
	} else {
		finalContent = strings.Join(selectedURLs, " ")
	}

	err = t.CopyToClipboard(finalContent)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	count := len(selectedURLs)
	if allSelected {
		count = len(urls)
	}
	urlWord := "URL"
	if count > 1 {
		urlWord = "URLs"
	}
	fmt.Printf("\033[93madded %d %s to clipboard\033[0m\n", count, urlWord)
	return nil
}

func (t *Terminal) openEditorWithContent(content string) (string, error) {

	tempDir, err := util.GetTempDir()
	if err != nil {
		return "", fmt.Errorf("failed to get temp directory: %w", err)
	}

	tempFile, err := os.CreateTemp(tempDir, "viren_clipboard_*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempFile.Close()

	if err := t.runEditorWithFallback(tempFile.Name()); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	editedBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited content: %w", err)
	}

	return string(editedBytes), nil
}

func (t *Terminal) runEditorWithFallback(filePath string) error {
	return RunEditorWithFallback(t.config, filePath)
}

func normalizeURL(rawURL string) string {

	cleaned := strings.TrimRight(rawURL, ".,;:!?'\"`)>]}")

	u, err := url.Parse(cleaned)
	if err != nil {
		return cleaned
	}

	u.Host = strings.ToLower(u.Host)

	if u.RawQuery == "" && strings.HasSuffix(u.Path, "/") && u.Path != "/" {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String()
}

func (t *Terminal) ExtractURLsFromText(text string) []string {
	matches := urlRegex.FindAllString(text, -1)

	seen := make(map[string]bool)
	var uniqueURLs []string

	for _, rawURL := range matches {
		normalized := normalizeURL(rawURL)

		if normalized == "" {
			continue
		}

		if seen[normalized] {
			continue
		}

		seen[normalized] = true
		uniqueURLs = append(uniqueURLs, normalized)
	}

	return uniqueURLs
}

func (t *Terminal) ExtractURLsFromChatHistory(chatHistory []types.ChatHistory) []string {
	var allURLs []string
	seen := make(map[string]bool)

	for _, entry := range chatHistory {

		userURLs := t.ExtractURLsFromText(entry.User)
		for _, url := range userURLs {
			if !seen[url] {
				allURLs = append(allURLs, url)
				seen[url] = true
			}
		}

		botURLs := t.ExtractURLsFromText(entry.Bot)
		for _, url := range botURLs {
			if !seen[url] {
				allURLs = append(allURLs, url)
				seen[url] = true
			}
		}
	}

	return allURLs
}

func (t *Terminal) ExtractURLsFromMessages(messages []types.ChatMessage) []string {
	var allURLs []string
	seen := make(map[string]bool)

	for _, message := range messages {

		urls := t.ExtractURLsFromText(message.Content)
		for _, url := range urls {
			if !seen[url] {
				allURLs = append(allURLs, url)
				seen[url] = true
			}
		}
	}

	return allURLs
}

func (t *Terminal) ExtractSentencesFromText(text string) []string {

	rawSentences := sentenceRegex.Split(text, -1)

	var sentences []string
	for _, sentence := range rawSentences {

		cleaned := strings.TrimSpace(sentence)

		if cleaned == "" {
			continue
		}

		if len(cleaned) >= 10 && len(cleaned) <= 200 {
			sentences = append(sentences, cleaned)
		}
	}

	return sentences
}

func (t *Terminal) ExtractSentencesFromChatHistory(chatHistory []types.ChatHistory, messages []types.ChatMessage) []string {
	var allSentences []string
	seen := make(map[string]bool)

	for _, entry := range chatHistory {

		userSentences := t.ExtractSentencesFromText(entry.User)
		for _, sentence := range userSentences {
			if !seen[sentence] {
				allSentences = append(allSentences, sentence)
				seen[sentence] = true
			}
		}

		botSentences := t.ExtractSentencesFromText(entry.Bot)
		for _, sentence := range botSentences {
			if !seen[sentence] {
				allSentences = append(allSentences, sentence)
				seen[sentence] = true
			}
		}
	}

	for _, message := range messages {
		sentences := t.ExtractSentencesFromText(message.Content)
		for _, sentence := range sentences {
			if !seen[sentence] {
				allSentences = append(allSentences, sentence)
				seen[sentence] = true
			}
		}
	}

	return allSentences
}
