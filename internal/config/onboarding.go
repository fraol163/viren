package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fraol163/viren/internal/ui"
	"github.com/fraol163/viren/pkg/types"
)

// RunOnboarding runs the first-time setup wizard
func RunOnboarding(terminal *ui.Terminal, cfg *types.Config) error {
	terminal.ClearTerminal()
	terminal.ShowLogo()

	fmt.Println("\n\033[1;96mWELCOME TO VIREN. INITIALIZING NEURAL PROFILE...\033[0m")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(" \033[1;37;40m IDENTITY \033[0m What should I call you? ❯ ")
	name, _ := reader.ReadString('\n')
	cfg.UserProfile.Name = strings.TrimSpace(name)

	roles := []string{
		"Backend Developer", "Frontend Developer", "Full Stack Developer",
		"DevOps Engineer", "Mobile Developer (Expo)", "Mobile Developer (Flutter)",
		"Data Scientist", "Security Researcher", "AI Engineer", "System Architect",
	}
	role, err := terminal.FzfSelect(roles, "Select your primary role: ")
	if err != nil {
		return err
	}
	if role == "" {
		role = "Developer"
	}
	cfg.UserProfile.Role = role
	fmt.Printf(" \033[1;37;40m ROLE \033[0m %s\n", role)

	editors := []string{"VS Code", "Vim", "Neovim", "Zed", "IntelliJ", "Sublime Text", "Emacs", "Nano"}
	editor, err := terminal.FzfSelect(editors, "Preferred Editor: ")
	if err != nil {
		return err
	}
	if editor == "" {
		editor = "VS Code"
	}
	cfg.UserProfile.Environment = editor
	fmt.Printf(" \033[1;37;40m ENV \033[0m %s\n", editor)

	fmt.Print(" \033[1;37;40m AMBITION \033[0m What is your primary goal? ❯ ")
	ambition, _ := reader.ReadString('\n')
	cfg.UserProfile.Ambition = strings.TrimSpace(ambition)

	themes := ui.GetThemes()
	var themeNames []string
	for _, t := range themes {
		themeNames = append(themeNames, fmt.Sprintf("%s (%s)", t.Name, t.ID))
	}
	themeSelection, err := terminal.FzfSelect(themeNames, "Select Interface Theme: ")
	if err == nil && themeSelection != "" {
		parts := strings.Split(themeSelection, "(")
		if len(parts) >= 2 {
			themeID := strings.TrimSuffix(parts[len(parts)-1], ")")
			cfg.CurrentTheme = themeID
			cfg.UserProfile.Theme = themeID
			terminal.SetTheme(themeID)
		}
	}

	personalities := []string{"Analytical", "Creative", "Focused", "Empathetic", "Playful", "Balanced", "Rick Sanchez"}
	pSelection, err := terminal.FzfSelect(personalities, "Select AI Personality: ")
	if err == nil && pSelection != "" {
		pID := strings.ToLower(strings.Split(pSelection, " ")[0])
		if pSelection == "Rick Sanchez" {
			pID = "rick"
		}
		cfg.CurrentPersonality = pID
	}

	cfg.SystemPrompt = "VIREN SYSTEM ARCHITECTURE: Integrated Intelligence Layer\n\nYou must perfectly blend the following Domain Mode and Personality Mandate into a single, cohesive persona.\n\n--- DOMAIN MODE ---\nBalanced conversation with professional tone.\n\n"

	pMap := map[string]string{
		"analytical": "Logical, systematic, and data-driven.",
		"creative":   "Artistic, imaginative, and metaphor-heavy.",
		"focused":    "Goal-oriented, concise, and productivity-focused.",
		"empathetic": "Emotionally intelligent and supportive.",
		"playful":    "Fun, energetic, and full of wit.",
		"balanced":   "Versatile and adaptive.",
		"rick":       "Rick Sanchez from C-137. Cynical, scientifically brilliant, and sarcastic.",
	}
	if mandate, ok := pMap[cfg.CurrentPersonality]; ok {
		cfg.SystemPrompt += "--- PERSONALITY MANDATE ---\n" + mandate + "\n\n"
	}

	if cfg.UserProfile.Name != "" {
		cfg.SystemPrompt += fmt.Sprintf("--- USER NEURAL PROFILE ---\n- Identity: %s\n- Role: %s\n- Environment: %s\n- Ambition: %s",
			cfg.UserProfile.Name, cfg.UserProfile.Role, cfg.UserProfile.Environment, cfg.UserProfile.Ambition)
	}

	currentTheme := terminal.GetTheme()
	fmt.Printf("\n%s SUCCESS \033[0m Profile Established. Welcome aboard.\n\n", currentTheme.SuccessBox)

	return SaveConfig(cfg)
}

// SaveConfig writes the configuration to disk
func SaveConfig(cfg *types.Config) error {
	return SaveConfigToFile(cfg)
}
