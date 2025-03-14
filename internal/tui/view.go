package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var s strings.Builder

	// Create main container with dynamic width
	container := baseStyle.Width(m.Width - 4).Margin(2)

	// Header with app name and version
	header := titleStyle.Width(m.Width - 8).Render("Download Manager v0.1")
	s.WriteString(container.Render(header))

	// Error message if any
	if m.ErrorMessage != "" {
		s.WriteString("\n" + errorStyle.Render(m.ErrorMessage))
	}

	// Content with help text integrated
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

	return s.String()
}

func renderAddMenu(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("New Download"))
	s.WriteString("\n\n")

	inputBox := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(subtle)

	s.WriteString(inputBox.Render(
		menuItemStyle.Render("URL:    "+urlStyle.Render(m.InputURL+"_")) + "\n" +
			menuItemStyle.Render("Queue:  "+urlStyle.Render(m.InputQueue)),
	))

	// Help text in the middle
	s.WriteString("\n\n" + helpStyle.Render("[ Enter ] Save   [ Esc ] Cancel"))

	return s.String()
}

func renderDownloadList(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("Downloads"))
	s.WriteString("\n\n")

	if len(m.Downloads) == 0 {
		s.WriteString(menuItemStyle.Render("No downloads yet. Press 'a' to add one."))
		s.WriteString("\n\n" + helpStyle.Render("[ a ] Add New Download"))
		return s.String()
	}

	// Help text at the top
	s.WriteString(helpStyle.Render("[ p ] Pause   [ r ] Resume   [ Esc ] Back"))
	s.WriteString("\n\n")

	for i, d := range m.Downloads {
		// URL and status
		item := fmt.Sprintf("%s    %s",
			urlStyle.Render(d.URL),
			RenderStatus(d.Status),
		)

		// Progress bar with spinner for active downloads
		progressWidth := m.Width - 50
		if d.Status == "downloading" {
			frame := spinnerFrames[int(d.Progress)%len(spinnerFrames)]
			item += "\n" + spinnerStyle.Render(frame) + " " + RenderProgressBar(progressWidth, d.Progress)
		} else {
			item += "\n  " + RenderProgressBar(progressWidth, d.Progress)
		}

		// Speed
		if d.Speed > 0 {
			item += fmt.Sprintf("  %s", formatSpeed(d.Speed))
		}

		// Selection highlight
		if i == m.Selected {
			item = selectedStyle.Render(item)
		}

		s.WriteString(menuItemStyle.Render(item) + "\n\n")
	}

	return s.String()
}

func renderMainMenu(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("Main Menu"))
	s.WriteString("\n\n")

	menuBox := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(subtle)

	// Menu items
	menuItems := menuItemStyle.Render(`
    [a]  Add download
    [l]  List downloads
    [p]  Pause selected
    [r]  Resume selected
    [q]  Quit
   [↑]  Move up
   [↓]  Move down`)

	// Help text in the middle of the menu box
	helpText := "\n\n" + helpStyle.Render("[ q ] Quit") + "\n"

	s.WriteString(menuBox.Render(menuItems + helpText))

	return s.String()
}

func formatSpeed(speed int64) string {
	if speed > 1024*1024 {
		return fmt.Sprintf("%.1f MB/s", float64(speed)/(1024*1024))
	}
	return fmt.Sprintf("%.1f KB/s", float64(speed)/1024)
}
