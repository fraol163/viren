package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/fraol163/viren/pkg/types"
)

func SaveConfigToFile(config *types.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	virenDir := filepath.Join(homeDir, ".viren")
	configPath := filepath.Join(virenDir, "config.json")

	if err := os.MkdirAll(virenDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func loadConfigFromFile() (*types.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	virenDir := filepath.Join(homeDir, ".viren")
	configPath := filepath.Join(virenDir, "config.json")

	if err := os.MkdirAll(virenDir, 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &types.Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config types.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func mergeConfigs(defaultConfig, userConfig *types.Config) *types.Config {
	if userConfig.DefaultModel != "" {
		defaultConfig.DefaultModel = userConfig.DefaultModel

		if userConfig.CurrentModel == "" {
			defaultConfig.CurrentModel = userConfig.DefaultModel
		}
	}
	if userConfig.CurrentModel != "" {
		defaultConfig.CurrentModel = userConfig.CurrentModel
	}
	if userConfig.CurrentBaseURL != "" {
		defaultConfig.CurrentBaseURL = userConfig.CurrentBaseURL
	}
	if userConfig.SystemPrompt != "" {
		defaultConfig.SystemPrompt = userConfig.SystemPrompt
	}
	if userConfig.ExitKey != "" {
		defaultConfig.ExitKey = userConfig.ExitKey
	}
	if userConfig.ModelSwitch != "" {
		defaultConfig.ModelSwitch = userConfig.ModelSwitch
	}
	if userConfig.EditorInput != "" {
		defaultConfig.EditorInput = userConfig.EditorInput
	}
	if userConfig.ClearHistory != "" {
		defaultConfig.ClearHistory = userConfig.ClearHistory
	}
	if userConfig.HelpKey != "" {
		defaultConfig.HelpKey = userConfig.HelpKey
	}
	if userConfig.ExportChat != "" {
		defaultConfig.ExportChat = userConfig.ExportChat
	}
	if userConfig.Backtrack != "" {
		defaultConfig.Backtrack = userConfig.Backtrack
	}
	if userConfig.WebSearch != "" {
		defaultConfig.WebSearch = userConfig.WebSearch
	}
	if userConfig.NumSearchResults != 0 {
		defaultConfig.NumSearchResults = userConfig.NumSearchResults
	}
	if userConfig.SearchCountry != "" {
		defaultConfig.SearchCountry = userConfig.SearchCountry
	}
	if userConfig.SearchLang != "" {
		defaultConfig.SearchLang = userConfig.SearchLang
	}
	if userConfig.ScrapeURL != "" {
		defaultConfig.ScrapeURL = userConfig.ScrapeURL
	}
	if userConfig.CopyToClipboard != "" {
		defaultConfig.CopyToClipboard = userConfig.CopyToClipboard
	}
	if userConfig.QuickCopyLatest != "" {
		defaultConfig.QuickCopyLatest = userConfig.QuickCopyLatest
	}
	if userConfig.LoadFiles != "" {
		defaultConfig.LoadFiles = userConfig.LoadFiles
	}
	if userConfig.AnswerSearch != "" {
		defaultConfig.AnswerSearch = userConfig.AnswerSearch
	}
	if userConfig.PlatformSwitch != "" {
		defaultConfig.PlatformSwitch = userConfig.PlatformSwitch
	}
	if userConfig.AllModels != "" {
		defaultConfig.AllModels = userConfig.AllModels
	}
	if userConfig.CodeDump != "" {
		defaultConfig.CodeDump = userConfig.CodeDump
	}
	if userConfig.ShellRecord != "" {
		defaultConfig.ShellRecord = userConfig.ShellRecord
	}
	if userConfig.ShellOption != "" {
		defaultConfig.ShellOption = userConfig.ShellOption
	}
	if userConfig.MultiLine != "" {
		defaultConfig.MultiLine = userConfig.MultiLine
	}
	if userConfig.PreferredEditor != "" {
		defaultConfig.PreferredEditor = userConfig.PreferredEditor
	}
	if userConfig.CurrentPlatform != "" {
		defaultConfig.CurrentPlatform = userConfig.CurrentPlatform
	}
	if userConfig.CurrentMode != "" {
		defaultConfig.CurrentMode = userConfig.CurrentMode
	}
	if userConfig.CurrentTheme != "" {
		defaultConfig.CurrentTheme = userConfig.CurrentTheme
	}
	if userConfig.CurrentPersonality != "" {
		defaultConfig.CurrentPersonality = userConfig.CurrentPersonality
	}

	if userConfig.Regenerate != "" {
		defaultConfig.Regenerate = userConfig.Regenerate
	}
	if userConfig.ExplainCode != "" {
		defaultConfig.ExplainCode = userConfig.ExplainCode
	}
	if userConfig.Summarize != "" {
		defaultConfig.Summarize = userConfig.Summarize
	}
	if userConfig.GenerateTests != "" {
		defaultConfig.GenerateTests = userConfig.GenerateTests
	}
	if userConfig.GenerateDocs != "" {
		defaultConfig.GenerateDocs = userConfig.GenerateDocs
	}
	if userConfig.OptimizeCode != "" {
		defaultConfig.OptimizeCode = userConfig.OptimizeCode
	}
	if userConfig.GitCommand != "" {
		defaultConfig.GitCommand = userConfig.GitCommand
	}
	if userConfig.CompareFiles != "" {
		defaultConfig.CompareFiles = userConfig.CompareFiles
	}
	if userConfig.TranslateCode != "" {
		defaultConfig.TranslateCode = userConfig.TranslateCode
	}
	if userConfig.FindReplace != "" {
		defaultConfig.FindReplace = userConfig.FindReplace
	}
	if userConfig.CommandReference != "" {
		defaultConfig.CommandReference = userConfig.CommandReference
	}
	// UI commands merge
	if userConfig.ModeSwitch != "" {
		defaultConfig.ModeSwitch = userConfig.ModeSwitch
	}
	if userConfig.ThemeSwitch != "" {
		defaultConfig.ThemeSwitch = userConfig.ThemeSwitch
	}
	if userConfig.PersonalitySwitch != "" {
		defaultConfig.PersonalitySwitch = userConfig.PersonalitySwitch
	}
	if userConfig.Onboarding != "" {
		defaultConfig.Onboarding = userConfig.Onboarding
	}
	// Update system merge
	defaultConfig.AutoUpdate = userConfig.AutoUpdate
	if userConfig.UpdateCommand != "" {
		defaultConfig.UpdateCommand = userConfig.UpdateCommand
	}
	if userConfig.LastUpdateCheck > 0 {
		defaultConfig.LastUpdateCheck = userConfig.LastUpdateCheck
	}

	if userConfig.DefaultModel != "" || userConfig.CurrentPlatform != "" || userConfig.SystemPrompt != "" || userConfig.ShowSearchResults {
		defaultConfig.ShowSearchResults = userConfig.ShowSearchResults
	}

	if userConfig.DefaultModel != "" || userConfig.CurrentPlatform != "" || userConfig.SystemPrompt != "" || userConfig.MuteNotifications {
		defaultConfig.MuteNotifications = userConfig.MuteNotifications
	}

	if userConfig.DefaultModel != "" || userConfig.CurrentPlatform != "" || userConfig.SystemPrompt != "" {
		defaultConfig.EnableSessionSave = userConfig.EnableSessionSave
	}
	if userConfig.SaveAllSessions {
		defaultConfig.SaveAllSessions = userConfig.SaveAllSessions
	}

	if userConfig.ShallowLoadDirs != nil {
		defaultConfig.ShallowLoadDirs = userConfig.ShallowLoadDirs
	}

	if userConfig.Platforms != nil {
		for name, platform := range userConfig.Platforms {
			defaultConfig.Platforms[name] = platform
		}
	}

	if userConfig.UserProfile.Name != "" {
		defaultConfig.UserProfile.Name = userConfig.UserProfile.Name
	}
	if userConfig.UserProfile.Role != "" {
		defaultConfig.UserProfile.Role = userConfig.UserProfile.Role
	}
	if userConfig.UserProfile.Environment != "" {
		defaultConfig.UserProfile.Environment = userConfig.UserProfile.Environment
	}
	if userConfig.UserProfile.Ambition != "" {
		defaultConfig.UserProfile.Ambition = userConfig.UserProfile.Ambition
	}
	if userConfig.UserProfile.Theme != "" {
		defaultConfig.UserProfile.Theme = userConfig.UserProfile.Theme
	}

	return defaultConfig
}

func DefaultConfig() *types.Config {

	homeDir, _ := os.UserHomeDir()

	shallowDirs := []string{
		"/",
		"/Users/",
		"/home/",
		"/usr/",
		"/var/",
		"/opt/",
		"/Library/",
		"/System/",
		"/mnt/",
		"/media/",
		"/Applications/",
		"/tmp/",
	}
	if homeDir != "" {
		shallowDirs = append(shallowDirs, homeDir)
	}

	defaultConfig := &types.Config{
		OpenAIAPIKey:	"",
		DefaultModel:	"gpt-4.1-mini",
		CurrentModel:	"gpt-4.1-mini",
		SystemPrompt:	"You are a helpful assistant powered by Viren who provides concise, clear, and accurate answers. Be brief, but ensure the response fully addresses the question without leaving out important details. Do NOT use em dashes (—) characters ever. But still, do NOT go crazy long with your response if you DON'T HAVE TO. Always return any code or file output in a Markdown code fence, with syntax ```<language or filetype>\n...``` so it can be parsed automatically. Only do this when needed, no need to do this for responses just code segments and/or when directly asked to do so from the user.",
		ExitKey:	"!q",
		ModelSwitch:	"!m",
		EditorInput:	"!t",
		ClearHistory:	"!c",
		HelpKey:	"!h",
		ExportChat:	"!e",
		Backtrack:	"!b",
		WebSearch:	"!w",
		ShowSearchResults:	true,
		NumSearchResults:	5,
		SearchCountry:	"us",
		SearchLang:	"en",
		ScrapeURL:	"!s",
		CopyToClipboard:	"!y",
		QuickCopyLatest:	"cc",
		LoadFiles:	"!l",
		AnswerSearch:	"!a",
		PlatformSwitch:	"!p",
		AllModels:	"!o",
		CodeDump:	"!d",
		ShellRecord:	"!x",
		ShellOption:	"!x",
		MultiLine:	"\\",
		PreferredEditor:	"vim",
		CurrentPlatform:	"openai",
		CurrentMode:	"standard",
		CurrentTheme:	"deepspace",
		CurrentPersonality:	"balanced",
		MuteNotifications:	false,
		EnableSessionSave:	true,
		ShallowLoadDirs:	shallowDirs,

		Regenerate:	"!r",
		ExplainCode:	"!explain",
		Summarize:	"!summarize",
		GenerateTests:	"!test",
		GenerateDocs:	"!doc",
		OptimizeCode:	"!optimize",
		GitCommand:	"!git",
		CompareFiles:	"!compare",
		TranslateCode:	"!translate",
		FindReplace:	"!f",
		CommandReference:	"!cmd",
		// UI commands with defaults
		ModeSwitch:	"!v",
		ThemeSwitch:	"!z",
		PersonalitySwitch:	"!u",
		Onboarding:	"!onboard",
		// Update system
		AutoUpdate:	true,
		UpdateCommand:	"!update",
		Platforms: map[string]types.Platform{
			"groq": {
				Name:	"groq",
				BaseURL:	types.BaseURLValue{Single: "https://api.groq.com/openai/v1"},
				EnvName:	"GROQ_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://api.groq.com/openai/v1/models",
					JSONPath:	"data.id",
				},
			},
			"openrouter": {
				Name:	"openrouter",
				BaseURL:	types.BaseURLValue{Single: "https://openrouter.ai/api/v1"},
				EnvName:	"OPENROUTER_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://openrouter.ai/api/v1/models",
					JSONPath:	"data.id",
				},
			},
			"deepseek": {
				Name:	"deepseek",
				BaseURL:	types.BaseURLValue{Single: "https://api.deepseek.com"},
				EnvName:	"DEEP_SEEK_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://api.deepseek.com/models",
					JSONPath:	"data.id",
				},
			},
			"anthropic": {
				Name:	"anthropic",
				BaseURL:	types.BaseURLValue{Single: "https://api.anthropic.com/v1/"},
				EnvName:	"ANTHROPIC_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://api.anthropic.com/v1/models",
					JSONPath:	"data.id",
				},
			},
			"xai": {
				Name:	"xai",
				BaseURL:	types.BaseURLValue{Single: "https://api.x.ai/v1"},
				EnvName:	"XAI_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://api.x.ai/v1/models",
					JSONPath:	"data.id",
				},
			},
			"ollama": {
				Name:	"ollama",
				BaseURL:	types.BaseURLValue{Single: "http://127.0.0.1:11434/v1"},
				EnvName:	"ollama",
				Models: types.PlatformModels{
					URL:	"http://127.0.0.1:11434/api/tags",
					JSONPath:	"models.name",
				},
			},
			"together": {
				Name:	"together",
				BaseURL:	types.BaseURLValue{Single: "https://api.together.xyz/v1"},
				EnvName:	"TOGETHER_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://api.together.xyz/v1/models",
					JSONPath:	"id",
				},
			},
			"google": {
				Name:	"google",
				BaseURL:	types.BaseURLValue{Single: "https://generativelanguage.googleapis.com/v1beta/openai/"},
				EnvName:	"GEMINI_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://generativelanguage.googleapis.com/v1beta/models",
					JSONPath:	"models.name",
				},
			},
			"mistral": {
				Name:	"mistral",
				BaseURL:	types.BaseURLValue{Single: "https://api.mistral.ai/v1"},
				EnvName:	"MISTRAL_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://api.mistral.ai/v1/models",
					JSONPath:	"data.id",
				},
			},
			"amazon": {
				Name:	"amazon",
				BaseURL: types.BaseURLValue{
					Multi: []string{
						"https://bedrock-runtime.us-west-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.us-east-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.us-east-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-northeast-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-northeast-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-northeast-3.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-south-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-south-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-southeast-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ap-southeast-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.ca-central-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-central-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-central-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-north-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-south-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-south-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-west-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-west-2.amazonaws.com/openai/v1",
						"https://bedrock-runtime.eu-west-3.amazonaws.com/openai/v1",
						"https://bedrock-runtime.sa-east-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.us-gov-east-1.amazonaws.com/openai/v1",
						"https://bedrock-runtime.us-gov-west-1.amazonaws.com/openai/v1",
					},
				},
				EnvName:	"AWS_BEDROCK_API_KEY",
				Models: types.PlatformModels{
					URL:	"https://bedrock.us-west-2.amazonaws.com/foundation-models",
					JSONPath:	"modelSummaries.modelId",
				},
			},
		},
	}

	if platformEnv := os.Getenv("VIREN_DEFAULT_PLATFORM"); platformEnv != "" {
		defaultConfig.CurrentPlatform = platformEnv
	}
	if modelEnv := os.Getenv("VIREN_DEFAULT_MODEL"); modelEnv != "" {
		defaultConfig.CurrentModel = modelEnv
		defaultConfig.DefaultModel = modelEnv
	}

	userConfig, err := loadConfigFromFile()
	if err == nil {
		defaultConfig = mergeConfigs(defaultConfig, userConfig)
	}

	return defaultConfig
}

func InitializeAppState() *types.AppState {
	cfg := DefaultConfig()

	state := &types.AppState{
		Config:	cfg,
		Messages: []types.ChatMessage{
			{Role: "system", Content: cfg.SystemPrompt},
		},
		ChatHistory: []types.ChatHistory{
			{Time: time.Now().Unix(), User: cfg.SystemPrompt, Bot: ""},
		},
		CurrentMode:	cfg.CurrentMode,
		CurrentTheme:	cfg.CurrentTheme,
		CurrentPersonality:	cfg.CurrentPersonality,
		RecentlyCreatedFiles:	[]string{},
		IsStreaming:	false,
		StreamingCancel:	nil,
		IsExecutingCommand:	false,
		CommandCancel:	nil,
	}

	return state
}
