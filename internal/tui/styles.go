package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Base styles
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(CurrentTheme.Border)

	// Container style for the entire app
	containerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("87")).
			Padding(1).
			Margin(1)

	// Container style for centering content
	centerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(100).
			MarginLeft(2)

	// Header styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 1).
			MarginBottom(1).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("87"))

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
			BorderForeground(lipgloss.Color("87")).
			Padding(0, 1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.Subtle).
			Align(lipgloss.Center).
			PaddingTop(1)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.Error).
			Bold(true).
			Align(lipgloss.Center)

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
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("87"))

	// Input box style
	inputBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Spinner style
	spinnerStyle = lipgloss.NewStyle().
			Bold(true)

	// Tab styles
	tabStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(CurrentTheme.Border).
			Padding(0, 1).
			Align(lipgloss.Center)

	activeTabStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(CurrentTheme.Special).
			Foreground(CurrentTheme.Special).
			Padding(0, 1).
			Bold(true).
			Align(lipgloss.Center)

	// Header style for tables
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Align(lipgloss.Center)

	// Table cell styles
	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 2).
			MaxHeight(1)

	tableHeaderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(CurrentTheme.Highlight).
				Foreground(CurrentTheme.Special).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)

	// Table container style
	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Border).
			Padding(0, 1).
			Align(lipgloss.Center).
			Bold(true).
			Italic(true)

	// Selected item style
	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1)

	// Queue management specific styles
	queueStatsStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(CurrentTheme.Special).
			Padding(1).
			Align(lipgloss.Center).
			Width(60).
			MarginLeft(2)

	queueStatItemStyle = lipgloss.NewStyle().
				Foreground(CurrentTheme.Text).
				Bold(true).
				Italic(true)

	queueActiveItemStyle = lipgloss.NewStyle().
				Background(CurrentTheme.Primary).
				Foreground(CurrentTheme.Text).
				Bold(true)

	queueFormStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("87")).
			Padding(1).
			MarginTop(1)

	queueFormFieldStyle = lipgloss.NewStyle().
				Padding(0, 2)

	queueFormSelectedFieldStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("87")).
					Background(lipgloss.Color("237")).
					Padding(0, 2)

	// Status message styles
	successStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.Special).
			Bold(true).
			Align(lipgloss.Center)

	// Table styles
	selectedRowStyle = lipgloss.NewStyle().
				Background(CurrentTheme.Highlight).
				Foreground(CurrentTheme.Text).
				Bold(true)

	normalRowStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.Text)
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
	baseStyle = baseStyle.BorderForeground(CurrentTheme.Border)

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
	helpStyle = helpStyle.Foreground(CurrentTheme.Subtle)

	// Error style
	errorStyle = errorStyle.Foreground(CurrentTheme.Error)

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
	tabStyle = tabStyle.BorderForeground(CurrentTheme.Border)

	activeTabStyle = activeTabStyle.
		BorderForeground(CurrentTheme.Special).
		Foreground(CurrentTheme.Special)

	// Header style
	headerStyle = headerStyle.
		BorderForeground(lipgloss.Color(CurrentTheme.Subtle.Dark)).
		Foreground(lipgloss.Color(CurrentTheme.Primary.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Selected item style
	selectedItemStyle = selectedItemStyle.
		Foreground(lipgloss.Color(CurrentTheme.Highlight.Dark)).
		Background(lipgloss.Color(CurrentTheme.Background))

	// Table styles
	tableStyle = tableStyle.BorderForeground(CurrentTheme.Border)
	tableHeaderStyle = tableHeaderStyle.
		BorderForeground(CurrentTheme.Highlight).
		Foreground(CurrentTheme.Special)
	selectedRowStyle = selectedRowStyle.
		Background(CurrentTheme.Highlight).
		Foreground(CurrentTheme.Text)
	normalRowStyle = normalRowStyle.Foreground(CurrentTheme.Text)

	// Queue styles
	queueStatsStyle = queueStatsStyle.BorderForeground(CurrentTheme.Special)
	queueStatItemStyle = queueStatItemStyle.Foreground(CurrentTheme.Text)
	queueActiveItemStyle = queueActiveItemStyle.
		Background(CurrentTheme.Primary).
		Foreground(CurrentTheme.Text)

	// Status message styles
	successStyle = successStyle.Foreground(CurrentTheme.Special)
}
