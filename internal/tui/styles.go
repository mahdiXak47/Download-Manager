package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Base styles
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1).
			MarginBottom(1)

	// Header styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 3).
			MarginBottom(1).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.DoubleBorder())

	// Content styles
	urlStyle = lipgloss.NewStyle().
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
			Padding(0, 1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1).
			MarginTop(1).
			MarginBottom(1).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder())

	// Error style
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			MarginTop(1).
			Padding(1).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder())

	// Progress bar styles
	progressBarStyle = lipgloss.NewStyle().
				MarginLeft(2).
				MarginRight(2)

	progressBarFilledStyle = lipgloss.NewStyle()

	progressBarEmptyStyle = lipgloss.NewStyle()

	// Menu styles
	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			PaddingRight(4)

	menuHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 2).
			MarginBottom(1).
			Width(30).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder())

	// Input box style
	inputBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Spinner style
	spinnerStyle = lipgloss.NewStyle().
			Bold(true)

	// Tab styles
	tabStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder(), false, true, false, false).
		Align(lipgloss.Center)

	activeTabStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.DoubleBorder(), false, true, false, false).
		BorderBottom(true).
		Bold(true).
		Italic(true).
		Align(lipgloss.Center)
	
	// Header style for tables
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)

	// Selected item style
	selectedItemStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)
)

// Spinner frames for animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

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
		return statusStyle.Copy().
			Background(CurrentTheme.Special).
			Foreground(lipgloss.Color(CurrentTheme.Background)).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Special).
			Render(" " + status + " ")
	case "paused":
		return statusStyle.Copy().
			Background(CurrentTheme.Warning).
			Foreground(lipgloss.Color(CurrentTheme.Background)).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Warning).
			Render(" " + status + " ")
	case "completed":
		return statusStyle.Copy().
			Background(CurrentTheme.Highlight).
			Foreground(lipgloss.Color(CurrentTheme.Foreground)).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Highlight).
			Render(" " + status + " ")
	case "error":
		return statusStyle.Copy().
			Background(CurrentTheme.Danger).
			Foreground(lipgloss.Color(CurrentTheme.Foreground)).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Danger).
			Render(" " + status + " ")
	case "cancelled":
		return statusStyle.Copy().
			Background(CurrentTheme.Danger).
			Foreground(lipgloss.Color(CurrentTheme.Background)).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Danger).
			Render(" " + status + " ")
	default:
		return statusStyle.Copy().
			Background(CurrentTheme.Subtle).
			Foreground(lipgloss.Color(CurrentTheme.Background)).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Subtle).
			Render(" " + status + " ")
	}
}

// UpdateStyles updates all styles based on the current theme
func UpdateStyles() {
	// Base styles
	baseStyle = baseStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Border.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Header styles
	titleStyle = titleStyle.
		Foreground(lipgloss.Color(CurrentTheme.Primary.Dark)).
		BorderForeground(lipgloss.Color(CurrentTheme.Border.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Menu styles
	menuHeaderStyle = menuHeaderStyle.
		Foreground(lipgloss.Color(CurrentTheme.Primary.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	menuItemStyle = menuItemStyle.
		Foreground(lipgloss.Color(CurrentTheme.Text.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Content styles
	urlStyle = urlStyle.
		Foreground(lipgloss.Color(CurrentTheme.Secondary.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	statusStyle = statusStyle.
		Background(lipgloss.Color(CurrentTheme.Background))

	progressStyle = progressStyle.
		Foreground(lipgloss.Color(CurrentTheme.Text.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	selectedStyle = selectedStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Highlight.Dark)).
		Foreground(lipgloss.Color(CurrentTheme.Text.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Help style
	helpStyle = helpStyle.
		Foreground(lipgloss.Color(CurrentTheme.Text.Dark)).
		BorderForeground(lipgloss.Color(CurrentTheme.Subtle.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Error style
	errorStyle = errorStyle.
		Foreground(lipgloss.Color(CurrentTheme.Error.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Input style
	inputBoxStyle = inputBoxStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Border.Dark)).
		Foreground(lipgloss.Color(CurrentTheme.Text.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Spinner style
	spinnerStyle = spinnerStyle.
		Foreground(lipgloss.Color(CurrentTheme.Highlight.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))
		
	// Tab styles
	tabStyle = tabStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Subtle.Dark)).
		Foreground(lipgloss.Color(CurrentTheme.Text.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	activeTabStyle = activeTabStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Highlight.Dark)).
		Foreground(lipgloss.Color(CurrentTheme.Highlight.Dark)).
		Background(lipgloss.Color(CurrentTheme.Subtle.Dark))
		
	// Header style
	headerStyle = headerStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Subtle.Dark)).
		Foreground(lipgloss.Color(CurrentTheme.Primary.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))
		
	// Selected item style
	selectedItemStyle = selectedItemStyle.
		Foreground(lipgloss.Color(CurrentTheme.Highlight.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))
}
