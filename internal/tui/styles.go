package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette - red tones
	subtle    = lipgloss.AdaptiveColor{Light: "#FFD1D1", Dark: "#8B0000"}
	highlight = lipgloss.AdaptiveColor{Light: "#FF4D4D", Dark: "#FF6B6B"}
	special   = lipgloss.AdaptiveColor{Light: "#FF7676", Dark: "#FF8989"}
	accent    = lipgloss.AdaptiveColor{Light: "#FF3333", Dark: "#FF4444"}
	muted     = lipgloss.AdaptiveColor{Light: "#FFB3B3", Dark: "#CC0000"}

	// Base styles with subtle gradient border
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(highlight).
			Padding(1, 2)

	// Header styles with more prominence
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(highlight).
			Padding(1, 3).
			MarginBottom(1).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(highlight)

	// Content styles
	urlStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			PaddingLeft(2)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1)

	progressStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1).
			Background(lipgloss.Color("#FFF0F0")).
			Foreground(lipgloss.Color("#000000"))

	// Help style with better visibility and centered positioning
	helpStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#FFF5F5")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	// Error style with better contrast
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			MarginTop(1).
			Padding(1, 2).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000"))

	// Progress bar styles with rounded corners
	progressBarStyle = lipgloss.NewStyle().
				MarginLeft(2).
				MarginRight(2)

	progressBarFilledStyle = lipgloss.NewStyle().
				Background(accent).
				Foreground(lipgloss.Color("#FFFFFF"))

	progressBarEmptyStyle = lipgloss.NewStyle().
				Background(subtle).
				Foreground(lipgloss.Color("#FFFFFF"))

	// Menu styles with better hierarchy
	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingRight(4)

	menuHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Background(lipgloss.Color("#FFF0F0")).
			PaddingLeft(2).
			PaddingRight(2).
			PaddingTop(1).
			PaddingBottom(1).
			MarginBottom(1).
			Width(30).
			Align(lipgloss.Center)

	// Spinner style with accent color
	spinnerStyle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true)
)

// Spinner frames for animation
var spinnerFrames = []string{"◐", "◓", "◑", "◒"}

// RenderProgressBar creates a styled progress bar
func RenderProgressBar(width int, percent float64) string {
	w := float64(width)
	filled := int(w * percent / 100)
	empty := width - filled

	bar := progressBarFilledStyle.Render(strings.Repeat("━", filled))
	bar += progressBarEmptyStyle.Render(strings.Repeat("─", empty))
	return fmt.Sprintf(" %s %.1f%%", bar, percent)
}

// RenderStatus returns a styled status indicator
func RenderStatus(status string) string {
	switch status {
	case "downloading":
		return statusStyle.
			Background(special).
			Foreground(lipgloss.Color("#FFFFFF")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(special).
			Render(" " + status + " ")
	case "paused":
		return statusStyle.
			Background(subtle).
			Foreground(lipgloss.Color("#000000")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Render(" " + status + " ")
	case "completed":
		return statusStyle.
			Background(highlight).
			Foreground(lipgloss.Color("#FFFFFF")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Render(" " + status + " ")
	default:
		return statusStyle.
			Background(muted).
			Foreground(lipgloss.Color("#FFFFFF")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(muted).
			Render(" " + status + " ")
	}
}
