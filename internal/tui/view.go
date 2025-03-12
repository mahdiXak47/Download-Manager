package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF875F")).
			MarginLeft(2)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F87FF"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5FFF87"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3A3A"))
)

func (m Model) View() string {
	var s strings.Builder

	// Header
	s.WriteString(titleStyle.Render("Download Manager v0.1\n\n"))

	// Error message if any
	if m.ErrorMessage != "" {
		s.WriteString(errorStyle.Render(m.ErrorMessage + "\n\n"))
	}

	switch m.Menu {
	case "add":
		s.WriteString(renderAddMenu(m))
	case "list":
		s.WriteString(renderDownloadList(m))
	default:
		s.WriteString(renderMainMenu(m))
	}

	// Help
	s.WriteString("\n" + renderHelp(m))

	return s.String()
}

func renderAddMenu(m Model) string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Add Download\n\n"))
	s.WriteString("URL: " + urlStyle.Render(m.InputURL) + "█\n")
	s.WriteString("Queue: " + m.InputQueue + "\n")
	s.WriteString("\nPress Enter to add, Esc to cancel")
	return s.String()
}

func renderDownloadList(m Model) string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Downloads\n\n"))

	if len(m.Downloads) == 0 {
		s.WriteString("No downloads yet. Press 'a' to add one.\n")
		return s.String()
	}

	for i, d := range m.Downloads {
		item := fmt.Sprintf("%s [%s] %.1f%%",
			urlStyle.Render(d.URL),
			statusStyle.Render(d.Status),
			d.Progress,
		)

		if i == m.Selected {
			item = selectedStyle.Render("> " + item)
		} else {
			item = "  " + item
		}

		s.WriteString(item + "\n")
	}

	return s.String()
}

func renderMainMenu(m Model) string {
	return `Commands:
  a: Add download
  l: List downloads
  p: Pause selected
  r: Resume selected
  q: Quit
  ↑/k: Move up
  ↓/j: Move down`
}

func renderHelp(m Model) string {
	switch m.Menu {
	case "add":
		return "Enter: Save • Esc: Cancel"
	case "list":
		return "p: Pause • r: Resume • Esc: Back"
	default:
		return "q: Quit"
	}
}
