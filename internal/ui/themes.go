package ui

import "github.com/fraol163/viren/pkg/types"

// GetThemes returns all 8 available themes based on FZX documentation
func GetThemes() []types.Theme {
	return []types.Theme{
		{
			ID:           "deepspace",
			Name:         "Deep Space (FZX)",
			LogoColor:    "\033[38;2;0;255;255m",
			UserBox:      "\033[38;2;0;0;0;48;2;0;102;204m",
			AssistantBox: "\033[38;2;0;0;0;48;2;0;204;102m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;51;51;153m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;0;255;127m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;255;51;51m",
			InfoBox:      "\033[38;2;0;0;0;48;2;0;153;255m",
			SystemBox:    "\033[38;2;0;0;0;48;2;0;204;204m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#000000\007",
			FgColor:      "\033]10;#00FFFF\007",
		},
		{
			ID:           "neonfuture",
			Name:         "Neon Future (FZX)",
			LogoColor:    "\033[38;2;255;0;255m",
			UserBox:      "\033[38;2;0;0;0;48;2;57;255;20m",
			AssistantBox: "\033[38;2;0;0;0;48;2;255;0;255m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;0;255;255m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;57;255;20m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;255;0;0m",
			InfoBox:      "\033[38;2;0;0;0;48;2;0;255;255m",
			SystemBox:    "\033[38;2;0;0;0;48;2;255;255;0m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#000000\007",
			FgColor:      "\033]10;#FF00FF\007",
		},
		{
			ID:           "retrowave",
			Name:         "Retro Wave (FZX)",
			LogoColor:    "\033[38;2;255;106;213m",
			UserBox:      "\033[38;2;0;0;0;48;2;5;213;250m",
			AssistantBox: "\033[38;2;0;0;0;48;2;255;106;213m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;255;170;0m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;0;255;153m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;255;51;102m",
			InfoBox:      "\033[38;2;0;0;0;48;2;5;213;250m",
			SystemBox:    "\033[38;2;0;0;0;48;2;153;0;255m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#050010\007",
			FgColor:      "\033]10;#FF6AD5\007",
		},
		{
			ID:           "greenglow",
			Name:         "Green Glow (FZX)",
			LogoColor:    "\033[38;2;0;255;0m",
			UserBox:      "\033[38;2;0;0;0;48;2;0;255;0m",
			AssistantBox: "\033[38;2;0;0;0;48;2;0;153;0m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;0;102;0m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;0;255;0m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;204;0;0m",
			InfoBox:      "\033[38;2;0;0;0;48;2;0;153;0m",
			SystemBox:    "\033[38;2;0;0;0;48;2;0;255;0m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#000000\007",
			FgColor:      "\033]10;#00FF00\007",
		},
		{
			ID:           "purpledream",
			Name:         "Purple Dream (FZX)",
			LogoColor:    "\033[38;2;153;51;255m",
			UserBox:      "\033[38;2;0;0;0;48;2;102;0;204m",
			AssistantBox: "\033[38;2;0;0;0;48;2;204;51;255m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;51;0;102m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;153;255;51m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;255;51;51m",
			InfoBox:      "\033[38;2;0;0;0;48;2;204;51;255m",
			SystemBox:    "\033[38;2;0;0;0;48;2;102;0;204m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#100020\007",
			FgColor:      "\033]10;#9933FF\007",
		},
		{
			ID:           "darkmode",
			Name:         "Dark Mode (FZX)",
			LogoColor:    "\033[38;2;255;255;255m",
			UserBox:      "\033[38;2;0;0;0;48;2;60;60;60m",
			AssistantBox: "\033[38;2;0;0;0;48;2;180;180;180m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;40;40;40m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;100;100;100m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;255;0;0m",
			InfoBox:      "\033[38;2;0;0;0;48;2;100;100;100m",
			SystemBox:    "\033[38;2;0;0;0;48;2;180;180;180m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#000000\007",
			FgColor:      "\033]10;#FFFFFF\007",
		},
		{
			ID:           "systemlight",
			Name:         "System Light",
			LogoColor:    "\033[38;2;0;0;0m",
			UserBox:      "\033[38;2;255;255;255;48;2;0;120;215m",
			AssistantBox: "\033[38;2;255;255;255;48;2;16;124;16m",
			ThoughtBox:   "\033[38;2;255;255;255;48;2;102;102;102m",
			SuccessBox:   "\033[38;2;255;255;255;48;2;16;124;16m",
			ErrorBox:     "\033[38;2;255;255;255;48;2;232;17;35m",
			InfoBox:      "\033[38;2;255;255;255;48;2;0;120;215m",
			SystemBox:    "\033[38;2;255;255;255;48;2;0;120;215m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#FFFFFF\007",
			FgColor:      "\033]10;#000000\007",
		},
		{
			ID:           "systemdark",
			Name:         "System Dark",
			LogoColor:    "\033[38;2;255;255;255m",
			UserBox:      "\033[38;2;0;0;0;48;2;0;120;212m",
			AssistantBox: "\033[38;2;0;0;0;48;2;16;124;16m",
			ThoughtBox:   "\033[38;2;0;0;0;48;2;50;50;50m",
			SuccessBox:   "\033[38;2;0;0;0;48;2;16;124;16m",
			ErrorBox:     "\033[38;2;0;0;0;48;2;232;17;35m",
			InfoBox:      "\033[38;2;0;0;0;48;2;0;120;212m",
			SystemBox:    "\033[38;2;0;0;0;48;2;50;50;50m",
			BorderColor:  "\033[38;2;0;0;0m",
			MutedColor:   "\033[38;2;0;0;0m",
			BgColor:      "\033]11;#000000\007",
			FgColor:      "\033]10;#FFFFFF\007",
		},
	}
}

// GetDefaultTheme returns the default Deep Space theme
func GetDefaultTheme() types.Theme {
	return GetThemes()[0]
}

// GetThemeByID finds a theme by its ID
func GetThemeByID(id string) types.Theme {
	themes := GetThemes()
	for _, t := range themes {
		if t.ID == id {
			return t
		}
	}
	return GetDefaultTheme()
}
