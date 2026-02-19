package types

import "encoding/json"

type BaseURLValue struct {
	Single	string
	Multi	[]string
}

func (b *BaseURLValue) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		b.Single = str
		b.Multi = nil
		return nil
	}

	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		b.Single = ""
		b.Multi = arr
		return nil
	}

	return json.Unmarshal(data, &str)
}

func (b BaseURLValue) MarshalJSON() ([]byte, error) {
	if b.IsMulti() {
		return json.Marshal(b.Multi)
	}
	if b.Single != "" {
		return json.Marshal(b.Single)
	}
	return json.Marshal("")
}

func (b *BaseURLValue) IsMulti() bool {
	return len(b.Multi) > 0
}

func (b *BaseURLValue) GetURLs() []string {
	if b.IsMulti() {
		return b.Multi
	}
	if b.Single != "" {
		return []string{b.Single}
	}
	return []string{}
}

type ChatMessage struct {
	Role	string		`json:"role"`
	Content	string		`json:"content"`
}

type ChatHistory struct {
	Time	int64		`json:"time"`
	User	string		`json:"user"`
	Bot	string		`json:"bot"`
	Platform	string		`json:"platform"`
	Model	string		`json:"model"`
}

type Platform struct {
	Name	string		`json:"name"`
	BaseURL	BaseURLValue		`json:"base_url"`
	EnvName	string		`json:"env_name"`
	Models	PlatformModels		`json:"models"`
	Headers	map[string]string		`json:"headers"`
}

type PlatformModels struct {
	URL	string		`json:"url"`
	JSONPath	string		`json:"json_name_path"`
	Headers	map[string]string		`json:"headers"`
}

type UserProfile struct {
	Name	string		`json:"name"`
	Role	string		`json:"role"`
	Environment	string		`json:"environment"`
	Ambition	string		`json:"ambition"`
	Theme	string		`json:"theme"`
}

type Config struct {
	UserProfile	UserProfile		`json:"user_profile,omitempty"`
	OpenAIAPIKey	string		`json:"openai_api_key,omitempty"`
	DefaultModel	string		`json:"default_model"`
	CurrentModel	string		`json:"current_model"`
	CurrentBaseURL	string		`json:"current_base_url"`
	SystemPrompt	string		`json:"system_prompt"`
	ExitKey	string		`json:"exit_key"`
	ModelSwitch	string		`json:"model_switch"`
	EditorInput	string		`json:"editor_input"`
	ClearHistory	string		`json:"clear_history"`
	HelpKey	string		`json:"help_key"`
	ExportChat	string		`json:"export_chat"`
	Backtrack	string		`json:"backtrack"`
	WebSearch	string		`json:"web_search"`
	ShowSearchResults	bool		`json:"show_search_results"`
	NumSearchResults	int		`json:"num_search_results"`
	SearchCountry	string		`json:"search_country"`
	SearchLang	string		`json:"search_lang"`
	ScrapeURL	string		`json:"scrape_url"`
	CopyToClipboard	string		`json:"copy_to_clipboard"`
	QuickCopyLatest	string		`json:"quick_copy_latest"`
	LoadFiles	string		`json:"load_files"`
	AnswerSearch	string		`json:"answer_search"`
	PlatformSwitch	string		`json:"platform_switch"`
	CodeDump	string		`json:"code_dump"`
	ShellRecord	string		`json:"shell_record"`
	ShellOption	string		`json:"shell_option"`
	MultiLine	string		`json:"multi_line"`
	PreferredEditor	string		`json:"preferred_editor"`
	CurrentPlatform	string		`json:"current_platform"`
	CurrentMode	string		`json:"current_mode"`
	CurrentTheme	string		`json:"current_theme"`
	CurrentPersonality	string		`json:"current_personality"`
	AllModels	string		`json:"all_models,omitempty"`
	MuteNotifications	bool		`json:"mute_notifications,omitempty"`
	EnableSessionSave	bool		`json:"enable_session_save"`
	SaveAllSessions	bool		`json:"save_all_sessions,omitempty"`
	ShallowLoadDirs	[]string		`json:"shallow_load_dirs,omitempty"`
	IsPipedOutput	bool		`json:"-"`
	Platforms	map[string]Platform		`json:"platforms,omitempty"`
	// New commands
	Regenerate	string		`json:"regenerate,omitempty"`
	ExplainCode	string		`json:"explain_code,omitempty"`
	Summarize	string		`json:"summarize,omitempty"`
	GenerateTests	string		`json:"generate_tests,omitempty"`
	GenerateDocs	string		`json:"generate_docs,omitempty"`
	OptimizeCode	string		`json:"optimize_code,omitempty"`
	GitCommand	string		`json:"git_command,omitempty"`
	CompareFiles	string		`json:"compare_files,omitempty"`
	TranslateCode	string		`json:"translate_code,omitempty"`
	FindReplace	string		`json:"find_replace,omitempty"`
	CommandReference	string		`json:"command_reference,omitempty"`
	// UI commands (configurable)
	ModeSwitch	string		`json:"mode_switch,omitempty"`
	ThemeSwitch	string		`json:"theme_switch,omitempty"`
	PersonalitySwitch	string		`json:"personality_switch,omitempty"`
	Onboarding	string		`json:"onboarding,omitempty"`
	// Update system
	AutoUpdate		bool		`json:"auto_update,omitempty"`
	UpdateCommand	string		`json:"update_command,omitempty"`
	LastUpdateCheck	int64		`json:"last_update_check,omitempty"`
}

type ExportEntry struct {
	Platform	string		`json:"platform"`
	ModelName	string		`json:"model_name"`
	UserPrompt	string		`json:"user_prompt"`
	BotResponse	string		`json:"bot_response"`
	Timestamp	int64		`json:"timestamp"`
}

type ChatExport struct {
	ExportedAt	int64		`json:"exported_at"`
	Entries	[]ExportEntry		`json:"entries"`
}

type SessionFile struct {
	Timestamp	int64		`json:"timestamp"`
	Platform	string		`json:"platform"`
	Model	string		`json:"model"`
	Mode	string		`json:"mode"`
	Personality	string		`json:"personality"`
	SystemPrompt	string		`json:"system_prompt"`
	BaseURL	string		`json:"base_url"`
	ChatHistory	[]ChatHistory		`json:"messages"`
}

type AppState struct {
	Config	*Config
	Messages	[]ChatMessage
	ChatHistory	[]ChatHistory
	CurrentMode	string
	CurrentTheme	string
	CurrentPersonality	string
	RecentlyCreatedFiles	[]string
	IsStreaming	bool
	StreamingCancel	func()
	IsExecutingCommand	bool
	CommandCancel	func()
	SessionStartTime	int64
}

type Theme struct {
	ID	string
	Name	string
	LogoColor	string
	UserBox	string
	AssistantBox	string
	ThoughtBox	string
	SuccessBox	string
	ErrorBox	string
	InfoBox	string
	SystemBox	string
	BorderColor	string
	MutedColor	string
	BgColor	string
	FgColor	string
}
