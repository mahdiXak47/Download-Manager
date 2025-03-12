package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	// Base styles
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	// Header styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Background(subtle).
			Padding(0, 1).
			MarginBottom(1)

	// Content styles
	urlStyle = lipgloss.NewStyle().
			Foreground(special).
			PaddingLeft(2)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1)

	progressStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Background(subtle).
			Padding(0, 1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			MarginTop(1)

	// Progress bar styles
	progressBarStyle = lipgloss.NewStyle().
				Foreground(special).
				MarginLeft(2)

	progressBarFilledStyle = lipgloss.NewStyle().
				Background(highlight)

	progressBarEmptyStyle = lipgloss.NewStyle().
				Background(subtle)
)

// RenderProgressBar creates a styled progress bar
func RenderProgressBar(width int, percent float64) string {
	w := float64(width)
	filled := int(w * percent / 100)
	empty := width - filled

	bar := progressBarFilledStyle.Render(strings.Repeat("█", filled))
	bar += progressBarEmptyStyle.Render(strings.Repeat("░", empty))
	return progressBarStyle.Render(bar)
}

// RenderStatus returns a styled status indicator
func RenderStatus(status string) string {
	switch status {
	case "downloading":
		return statusStyle.
			Background(special).
			Render(status)
	case "paused":
		return statusStyle.
			Background(subtle).
			Render(status)
	case "completed":
		return statusStyle.
			Background(highlight).
			Render(status)
	default:
		return statusStyle.Render(status)
	}
}
