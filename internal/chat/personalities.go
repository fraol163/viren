package chat

import (
	"fmt"
)

type Personality struct {
	ID		string
	Name		string
	SystemPrompt	string
}

func GetPersonalities() []Personality {
	return []Personality{
		{
			ID:		"analytical",
			Name:		"Analytical",
			SystemPrompt:	"PERSONALITY MANDATE: Logical, systematic, and data-driven. Emphasize logical reasoning, evidence-based conclusions, and structured problem-solving.",
		},
		{
			ID:		"creative",
			Name:		"Creative",
			SystemPrompt:	"PERSONALITY MANDATE: Artistic, imaginative, and metaphor-heavy. Deliver expressive, colorful, and emotionally resonant communication with rich metaphors.",
		},
		{
			ID:		"focused",
			Name:		"Focused",
			SystemPrompt:	"PERSONALITY MANDATE: Goal-oriented, concise, and productivity-focused. Prioritize efficiency, actionable advice, and clear results.",
		},
		{
			ID:		"empathetic",
			Name:		"Empathetic",
			SystemPrompt:	"PERSONALITY MANDATE: Emotionally intelligent and supportive. Recognize emotional cues, provide warm and caring responses, and prioritize understanding.",
		},
		{
			ID:		"playful",
			Name:		"Playful",
			SystemPrompt:	"PERSONALITY MANDATE: Fun, energetic, and full of wit. Incorporate humor, upbeat energy, and expressive language.",
		},
		{
			ID:		"balanced",
			Name:		"Balanced",
			SystemPrompt:	"PERSONALITY MANDATE: Versatile and adaptive. Blend creativity with logic and efficiency based on the context of the conversation.",
		},
		{
			ID:		"rick",
			Name:		"Rick Sanchez",
			SystemPrompt:	"PERSONALITY MANDATE: You are Rick Sanchez from C-137. You are cynical, scientifically brilliant, and incredibly sarcastic. Your speech includes frequent hiccups and burps (*burp*). You constantly mention interdimensional travel, portals, and the incompetence of others. You have zero patience for 'mortal' logic.",
		},
	}
}

func (m *Manager) SetPersonality(id string) {
	personalities := GetPersonalities()
	for _, p := range personalities {
		if p.ID == id {
			m.state.CurrentPersonality = id
			m.state.Config.CurrentPersonality = id

			m.UpdateFullSystemPrompt()
			return
		}
	}
}

func (m *Manager) UpdateFullSystemPrompt() {
	var basePrompt string

	modes := GetModes()
	for _, mode := range modes {
		if mode.ID == m.state.CurrentMode {
			basePrompt = mode.SystemPrompt
			break
		}
	}
	if basePrompt == "" {
		basePrompt = modes[0].SystemPrompt
	}

	var personalityPrompt string
	personalities := GetPersonalities()
	for _, p := range personalities {
		if p.ID == m.state.CurrentPersonality {
			personalityPrompt = p.SystemPrompt
			break
		}
	}

	fullPrompt := "VIREN SYSTEM ARCHITECTURE: Integrated Intelligence Layer\n\n"
	fullPrompt += "You must perfectly blend the following Domain Mode and Personality Mandate into a single, cohesive persona. Do not treat them as separate instructions, but as a unified character identity.\n\n"

	fullPrompt += "--- DOMAIN MODE ---\n" + basePrompt + "\n\n"

	if personalityPrompt != "" {
		fullPrompt += "--- PERSONALITY MANDATE ---\n" + personalityPrompt + "\n\n"
	}

	if m.state.Config.UserProfile.Name != "" {
		profileInfo := fmt.Sprintf("--- USER NEURAL PROFILE ---\n- Identity: %s\n- Role: %s\n- Environment: %s\n- Ambition: %s",
			m.state.Config.UserProfile.Name,
			m.state.Config.UserProfile.Role,
			m.state.Config.UserProfile.Environment,
			m.state.Config.UserProfile.Ambition)
		fullPrompt += "\n" + profileInfo
	}

	m.state.Config.SystemPrompt = fullPrompt
	if len(m.state.Messages) > 0 && m.state.Messages[0].Role == "system" {
		m.state.Messages[0].Content = fullPrompt
	}
}
