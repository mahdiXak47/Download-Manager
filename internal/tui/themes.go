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
		Subtle:     lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
		Special:    lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"},
		Danger:     lipgloss.AdaptiveColor{Light: "#FF5D62", Dark: "#FF6B70"},
		Warning:    lipgloss.AdaptiveColor{Light: "#F2B155", Dark: "#F9C97C"},
		Info:       lipgloss.AdaptiveColor{Light: "#2D9CDB", Dark: "#4DA8DA"},
		Primary:    lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"},
		Text:       lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#EEEEEE"},
		Border:     lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"},
		Error:      lipgloss.AdaptiveColor{Light: "#FF5D62", Dark: "#FF6B70"},
		Background: "#2F2F2F",
		Foreground: "#FFFFFF",
	}

	OceanTheme = Theme{
		Name:       "ocean",
		Subtle:     lipgloss.AdaptiveColor{Light: "#A7C4C2", Dark: "#2C3E50"},
		Highlight:  lipgloss.AdaptiveColor{Light: "#6C5CE7", Dark: "#A29BFE"},
		Special:    lipgloss.AdaptiveColor{Light: "#00B894", Dark: "#55EFC4"},
		Danger:     lipgloss.AdaptiveColor{Light: "#D63031", Dark: "#FF7675"},
		Warning:    lipgloss.AdaptiveColor{Light: "#FDCB6E", Dark: "#FFA502"},
		Info:       lipgloss.AdaptiveColor{Light: "#0984E3", Dark: "#74B9FF"},
		Primary:    lipgloss.AdaptiveColor{Light: "#6C5CE7", Dark: "#A29BFE"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#00B894", Dark: "#55EFC4"},
		Text:       lipgloss.AdaptiveColor{Light: "#F5F0E1", Dark: "#ECF0F1"},
		Border:     lipgloss.AdaptiveColor{Light: "#A7C4C2", Dark: "#2C3E50"},
		Error:      lipgloss.AdaptiveColor{Light: "#D63031", Dark: "#FF7675"},
		Background: "#1E3D59",
		Foreground: "#F5F0E1",
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
)

// CurrentTheme holds the active theme
var CurrentTheme = ModernTheme
