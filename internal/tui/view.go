package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	// Create main container
	container := baseStyle.Width(m.Width - 4).Margin(2)

	// Header
	header := titleStyle.Render("Download Manager v0.1")
	s.WriteString(container.Render(header))

	// Error message if any
	if m.ErrorMessage != "" {
		s.WriteString("\n" + errorStyle.Render(m.ErrorMessage))
	}

	// Content
	content := ""
	switch m.Menu {
	case "add":
		content = renderAddMenu(m)
	case "list":
		content = renderDownloadList(m)
	default:
		content = renderMainMenu(m)
	}
	s.WriteString("\n" + container.Render(content))

	// Help
	help := renderHelp(m)
	s.WriteString("\n" + helpStyle.Render(help))

	return s.String()
}

func renderAddMenu(m Model) string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Add Download"))
	s.WriteString("\n\n")
	s.WriteString("URL: " + urlStyle.Render(m.InputURL+"█"))
	s.WriteString("\nQueue: " + urlStyle.Render(m.InputQueue))
	return s.String()
}

func renderDownloadList(m Model) string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Downloads"))
	s.WriteString("\n\n")

	if len(m.Downloads) == 0 {
		return s.String() + "No downloads yet. Press 'a' to add one."
	}

	for i, d := range m.Downloads {
		// URL and status
		item := fmt.Sprintf("%s %s",
			urlStyle.Render(d.URL),
			RenderStatus(d.Status),
		)

		// Progress bar
		progressWidth := m.Width - 40 // Adjust based on other content
		item += "\n" + RenderProgressBar(progressWidth, d.Progress)

		// Speed
		if d.Speed > 0 {
			item += fmt.Sprintf(" %.1f MB/s", float64(d.Speed)/(1024*1024))
		}

		// Selection highlight
		if i == m.Selected {
			item = selectedStyle.Render(item)
		}

		s.WriteString(item + "\n\n")
	}

	return s.String()
}

func renderMainMenu(m Model) string {
	menu := `
Commands:
  a: Add download
  l: List downloads
  p: Pause selected
  r: Resume selected
  q: Quit
  ↑/k: Move up
  ↓/j: Move down`

	return titleStyle.Render("Main Menu") + menu
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
