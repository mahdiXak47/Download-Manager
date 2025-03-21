package tui

import "github.com/charmbracelet/lipgloss"

// Theme represents a color scheme
type Theme struct {
	Name       string
	Subtle     lipgloss.AdaptiveColor
	Highlight  lipgloss.AdaptiveColor
	Special    lipgloss.AdaptiveColor
	Danger     lipgloss.AdaptiveColor
	Warning    lipgloss.AdaptiveColor
	Info       lipgloss.AdaptiveColor
	Primary    lipgloss.AdaptiveColor
	Secondary  lipgloss.AdaptiveColor
	Text       lipgloss.AdaptiveColor
	Border     lipgloss.AdaptiveColor
	Error      lipgloss.AdaptiveColor
	Background string
	Foreground string
}

// Available themes
var (
	ModernTheme = Theme{
		Name:       "modern",
		Subtle:     lipgloss.AdaptiveColor{Light: "#6C7A89", Dark: "#4A5568"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
		Special:    lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF5D62", Dark: "#FF6B70"},
		Warning:    lipgloss.AdaptiveColor{Light: "#F2B155", Dark: "#F9C97C"},
		Info:       lipgloss.AdaptiveColor{Light: "#2D9CDB", Dark: "#4DA8DA"},
		Primary:    lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"},
		Text:       lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#EEEEEE"},
		Border:     lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF5D62", Dark: "#FF6B70"},
		Background: "#1E1E2E",
		Foreground: "#FFFFFF",
	}

	OceanTheme = Theme{
		Name:       "ocean",
		Subtle:     lipgloss.AdaptiveColor{Light: "#A7C4C2", Dark: "#2C3E50"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#00B4D8", Dark: "#90E0EF"},
		Special:    lipgloss.AdaptiveColor{Light: "#48CAE4", Dark: "#00B4D8"},
		Danger:     lipgloss.AdaptiveColor{Light: "#D63031", Dark: "#FF7675"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FFB703", Dark: "#FB8500"},
		Info:       lipgloss.AdaptiveColor{Light: "#0096C7", Dark: "#48CAE4"},
		Primary:    lipgloss.AdaptiveColor{Light: "#023E8A", Dark: "#0077B6"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#48CAE4", Dark: "#90E0EF"},
		Text:       lipgloss.AdaptiveColor{Light: "#CAF0F8", Dark: "#CAF0F8"},
		Border:     lipgloss.AdaptiveColor{Light: "#0077B6", Dark: "#48CAE4"},
		Error:      lipgloss.AdaptiveColor{Light: "#D63031", Dark: "#FF7675"},
		Background: "#03045E",
		Foreground: "#CAF0F8",
	}

	SolarizedTheme = Theme{
		Name:       "solarized",
		Subtle:     lipgloss.AdaptiveColor{Light: "#93A1A1", Dark: "#586E75"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#268BD2", Dark: "#839496"},
		Special:    lipgloss.AdaptiveColor{Light: "#859900", Dark: "#B58900"},
		Danger:     lipgloss.AdaptiveColor{Light: "#DC322F", Dark: "#CB4B16"},
		Warning:    lipgloss.AdaptiveColor{Light: "#B58900", Dark: "#CB4B16"},
		Info:       lipgloss.AdaptiveColor{Light: "#2AA198", Dark: "#6C71C4"},
		Primary:    lipgloss.AdaptiveColor{Light: "#268BD2", Dark: "#839496"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#859900", Dark: "#B58900"},
		Text:       lipgloss.AdaptiveColor{Light: "#FDF6E3", Dark: "#EEE8D5"},
		Border:     lipgloss.AdaptiveColor{Light: "#93A1A1", Dark: "#586E75"},
		Error:      lipgloss.AdaptiveColor{Light: "#DC322F", Dark: "#CB4B16"},
		Background: "#002B36",
		Foreground: "#FDF6E3",
	}

	NordTheme = Theme{
		Name:       "nord",
		Subtle:     lipgloss.AdaptiveColor{Light: "#D8DEE9", Dark: "#4C566A"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#88C0D0", Dark: "#81A1C1"},
		Special:    lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#8FBCBB"},
		Danger:     lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#D08770"},
		Warning:    lipgloss.AdaptiveColor{Light: "#EBCB8B", Dark: "#D08770"},
		Info:       lipgloss.AdaptiveColor{Light: "#5E81AC", Dark: "#88C0D0"},
		Primary:    lipgloss.AdaptiveColor{Light: "#88C0D0", Dark: "#81A1C1"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#A3BE8C", Dark: "#8FBCBB"},
		Text:       lipgloss.AdaptiveColor{Light: "#ECEFF4", Dark: "#D8DEE9"},
		Border:     lipgloss.AdaptiveColor{Light: "#D8DEE9", Dark: "#4C566A"},
		Error:      lipgloss.AdaptiveColor{Light: "#BF616A", Dark: "#D08770"},
		Background: "#2E3440",
		Foreground: "#ECEFF4",
	}

	// New Themes
	SynthwaveTheme = Theme{
		Name:       "synthwave",
		Subtle:     lipgloss.AdaptiveColor{Light: "#2B213A", Dark: "#2B213A"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#FF00FF", Dark: "#FF00FF"},
		Special:    lipgloss.AdaptiveColor{Light: "#00FFF5", Dark: "#00FFF5"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF3366", Dark: "#FF3366"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FFD700", Dark: "#FFD700"},
		Info:       lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#00FFFF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#FF00FF", Dark: "#FF00FF"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#00FFF5", Dark: "#00FFF5"},
		Text:       lipgloss.AdaptiveColor{Light: "#F0F0FF", Dark: "#F0F0FF"},
		Border:     lipgloss.AdaptiveColor{Light: "#FF00FF", Dark: "#FF00FF"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF3366", Dark: "#FF3366"},
		Background: "#241B2F",
		Foreground: "#F0F0FF",
	}

	DraculaTheme = Theme{
		Name:       "dracula",
		Subtle:     lipgloss.AdaptiveColor{Light: "#44475A", Dark: "#44475A"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#BD93F9", Dark: "#BD93F9"},
		Special:    lipgloss.AdaptiveColor{Light: "#50FA7B", Dark: "#50FA7B"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF5555", Dark: "#FF5555"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FFB86C", Dark: "#FFB86C"},
		Info:       lipgloss.AdaptiveColor{Light: "#8BE9FD", Dark: "#8BE9FD"},
		Primary:    lipgloss.AdaptiveColor{Light: "#FF79C6", Dark: "#FF79C6"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#50FA7B", Dark: "#50FA7B"},
		Text:       lipgloss.AdaptiveColor{Light: "#F8F8F2", Dark: "#F8F8F2"},
		Border:     lipgloss.AdaptiveColor{Light: "#BD93F9", Dark: "#BD93F9"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF5555", Dark: "#FF5555"},
		Background: "#282A36",
		Foreground: "#F8F8F2",
	}

	CyberpunkTheme = Theme{
		Name:       "cyberpunk",
		Subtle:     lipgloss.AdaptiveColor{Light: "#2B2B2B", Dark: "#2B2B2B"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#F9E900", Dark: "#F9E900"},
		Special:    lipgloss.AdaptiveColor{Light: "#00FF9F", Dark: "#00FF9F"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF003C", Dark: "#FF003C"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FF9B00", Dark: "#FF9B00"},
		Info:       lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#00FFFF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#FF003C", Dark: "#FF003C"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#00FF9F", Dark: "#00FF9F"},
		Text:       lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
		Border:     lipgloss.AdaptiveColor{Light: "#F9E900", Dark: "#F9E900"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF003C", Dark: "#FF003C"},
		Background: "#0D0208",
		Foreground: "#FFFFFF",
	}

	RetroTheme = Theme{
		Name:       "retro",
		Subtle:     lipgloss.AdaptiveColor{Light: "#5C4B37", Dark: "#5C4B37"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#FFB000", Dark: "#FFB000"},
		Special:    lipgloss.AdaptiveColor{Light: "#8BC34A", Dark: "#8BC34A"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF6B6B"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FFD700", Dark: "#FFD700"},
		Info:       lipgloss.AdaptiveColor{Light: "#87CEEB", Dark: "#87CEEB"},
		Primary:    lipgloss.AdaptiveColor{Light: "#FFB000", Dark: "#FFB000"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#8BC34A", Dark: "#8BC34A"},
		Text:       lipgloss.AdaptiveColor{Light: "#F4D03F", Dark: "#F4D03F"},
		Border:     lipgloss.AdaptiveColor{Light: "#FFB000", Dark: "#FFB000"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF6B6B"},
		Background: "#2C231D",
		Foreground: "#F4D03F",
	}

	NeonTheme = Theme{
		Name:       "neon",
		Subtle:     lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#1A1A1A"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#FF10F0", Dark: "#FF10F0"},
		Special:    lipgloss.AdaptiveColor{Light: "#39FF14", Dark: "#39FF14"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF2400", Dark: "#FF2400"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FFFF00", Dark: "#FFFF00"},
		Info:       lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#00FFFF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#FF10F0", Dark: "#FF10F0"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#39FF14", Dark: "#39FF14"},
		Text:       lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
		Border:     lipgloss.AdaptiveColor{Light: "#FF10F0", Dark: "#FF10F0"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF2400", Dark: "#FF2400"},
		Background: "#000000",
		Foreground: "#FFFFFF",
	}

	// New Aurora Theme
	AuroraTheme = Theme{
		Name:       "aurora",
		Subtle:     lipgloss.AdaptiveColor{Light: "#4B6584", Dark: "#4B6584"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#7F00FF", Dark: "#7F00FF"},
		Special:    lipgloss.AdaptiveColor{Light: "#00FF87", Dark: "#00FF87"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF3366", Dark: "#FF3366"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FFA07A", Dark: "#FFA07A"},
		Info:       lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#00FFFF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#7F00FF", Dark: "#7F00FF"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#00FF87", Dark: "#00FF87"},
		Text:       lipgloss.AdaptiveColor{Light: "#E0FFFF", Dark: "#E0FFFF"},
		Border:     lipgloss.AdaptiveColor{Light: "#7F00FF", Dark: "#7F00FF"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF3366", Dark: "#FF3366"},
		Background: "#1A1B26",
		Foreground: "#E0FFFF",
	}
)

// CurrentTheme holds the active theme
var CurrentTheme = ModernTheme
