package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	// Create main container with dynamic width
	mainContainer := containerStyle.Width(m.Width - 4)

	// Header with app name and version
	header := titleStyle.Width(m.Width - 8).Render("Download Manager v0.1")

	// Build the content
	var content strings.Builder
	content.WriteString(header)

	// Error message if any
	if m.ErrorMessage != "" {
		content.WriteString("\n" + errorStyle.Render(m.ErrorMessage))
	}

	// Content based on the active tab
	var tabContent string
	switch m.ActiveTab {
	case AddDownloadTab:
		tabContent = renderAddDownloadTab(m)
	case DownloadListTab:
		tabContent = renderDownloadListTab(m)
	case QueueListTab:
		tabContent = renderQueueListTab(m)
	case SettingsTab:
		tabContent = renderSettingsTab(m)
	}
	content.WriteString("\n" + tabContent)

	// Tab bar at the bottom
	tabBar := renderTabBar(m)
	content.WriteString("\n" + tabBar)

	// Wrap everything in the main container
	return mainContainer.Render(content.String())
}

// renderTabBar creates the tab bar with number keys 1-4
func renderTabBar(m Model) string {
	width := m.Width - 12 // Account for margins/paddings
	tabWidth := width / 4

	// Create tab styles based on active tab
	tab1Style := tabStyle.Copy().Width(tabWidth)
	tab2Style := tabStyle.Copy().Width(tabWidth)
	tab3Style := tabStyle.Copy().Width(tabWidth)
	tab4Style := tabStyle.Copy().Width(tabWidth)

	// Highlight active tab with more distinctive styling
	switch m.ActiveTab {
	case AddDownloadTab:
		tab1Style = activeTabStyle.Copy().Width(tabWidth)
	case DownloadListTab:
		tab2Style = activeTabStyle.Copy().Width(tabWidth)
	case QueueListTab:
		tab3Style = activeTabStyle.Copy().Width(tabWidth)
	case SettingsTab:
		tab4Style = activeTabStyle.Copy().Width(tabWidth)
	}

	// Render tabs with number key indicators
	tab1 := tab1Style.Render("1: Add Download")
	tab2 := tab2Style.Render("2: Download List")
	tab3 := tab3Style.Render("3: Queues List")
	tab4 := tab4Style.Render("4: Settings")

	return lipgloss.JoinHorizontal(lipgloss.Top, tab1, tab2, tab3, tab4)
}

func renderAddDownloadTab(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("Add New Download"))
	s.WriteString("\n\n")

	// Success or error message (if any)
	if m.AddDownloadMessage != "" {
		msgStyle := errorStyle
		if m.AddDownloadSuccess {
			msgStyle = msgStyle.Copy().
				Foreground(lipgloss.Color(CurrentTheme.Special.Dark)).
				BorderForeground(lipgloss.Color(CurrentTheme.Special.Dark))
		} else {
			msgStyle = msgStyle.Copy().
				Foreground(lipgloss.Color(CurrentTheme.Error.Dark)).
				BorderForeground(lipgloss.Color(CurrentTheme.Error.Dark))
		}
		s.WriteString(msgStyle.Render(m.AddDownloadMessage) + "\n\n")
	}

	// Queue selection first
	if m.QueueSelectionMode {
		s.WriteString(menuHeaderStyle.Render("Select Download Queue"))
		s.WriteString("\n\n")

		// Available queues
		s.WriteString(menuItemStyle.Render("Available Queues:"))
		s.WriteString("\n\n")

		// List all queues
		for i, q := range m.Config.Queues {
			itemStyle := menuItemStyle
			if i == m.QueueSelected {
				itemStyle = selectedItemStyle
			}

			activeCount := 0
			for _, d := range m.Downloads {
				if d.Queue == q.Name && d.Status == "downloading" {
					activeCount++
				}
			}

			queueInfo := fmt.Sprintf("%s (%d/%d active)", q.Name, activeCount, q.MaxConcurrent)
			s.WriteString(itemStyle.Render(queueInfo) + "\n")
		}

		// Help text
		s.WriteString("\n" + helpStyle.Render("[ ↑/↓ ] Navigate   [ Enter ] Select   [ Esc ] Cancel"))
	} else if m.URLInputMode {
		// URL input field
		s.WriteString(menuHeaderStyle.Render("Enter Download URL"))
		s.WriteString("\n\n")

		// Selected queue display
		s.WriteString(menuItemStyle.Render("Selected Queue: " + urlStyle.Render(m.InputQueue)))
		s.WriteString("\n\n")

		// URL input field
		s.WriteString(inputBoxStyle.Render(
			menuItemStyle.Render("URL: " + urlStyle.Render(m.InputURL+"_")),
		))

		// Help text for input mode
		s.WriteString("\n\n" + helpStyle.Render("[ Enter ] Start Download   [ Esc ] Back"))
	} else {
		// Initial instructions
		s.WriteString(menuItemStyle.Render("Press Enter to add a new download"))
		s.WriteString("\n\n" + helpStyle.Render("[ Enter ] Start   [ Esc ] Back"))
	}

	return s.String()
}

func renderDownloadListTab(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("Download List"))
	s.WriteString("\n\n")

	// Show retry message if any
	if m.DownloadListMessage != "" {
		msgStyle := errorStyle
		if m.DownloadListSuccess {
			// Blue color for success messages
			msgStyle = msgStyle.Copy().
				Foreground(lipgloss.Color(CurrentTheme.Special.Dark)).
				BorderForeground(lipgloss.Color(CurrentTheme.Special.Dark))
		} else {
			// Red color for error messages
			msgStyle = msgStyle.Copy().
				Foreground(lipgloss.Color(CurrentTheme.Error.Dark)).
				BorderForeground(lipgloss.Color(CurrentTheme.Error.Dark))
		}
		s.WriteString(msgStyle.Render(m.DownloadListMessage) + "\n\n")
	}

	if len(m.Downloads) == 0 {
		s.WriteString(menuItemStyle.Render("No downloads yet. Press '1' to switch to Add Download tab."))
	} else {
		// Table header
		header := lipgloss.JoinHorizontal(lipgloss.Top,
			headerStyle.Width(5).Render("#"),
			headerStyle.Width(35).Render("URL"),
			headerStyle.Width(15).Render("Status"),
			headerStyle.Width(15).Render("Progress"),
			headerStyle.Width(15).Render("Speed"),
		)
		s.WriteString(header + "\n")

		// Downloads list
		for i, d := range m.Downloads {
			// Highlight selected download
			itemStyle := menuItemStyle
			if i == m.Selected {
				itemStyle = selectedItemStyle
			}

			// Format the line
			line := lipgloss.JoinHorizontal(lipgloss.Top,
				itemStyle.Width(5).Render(fmt.Sprintf("%d", i+1)),
				itemStyle.Width(35).Render(truncateString(d.URL, 32)),
				itemStyle.Width(15).Render(d.Status),
				itemStyle.Width(15).Render(fmt.Sprintf("%.1f%%", d.Progress)),
				itemStyle.Width(15).Render(formatSpeed(d.Speed)),
			)
			s.WriteString(line + "\n")
		}
	}

	// Help text
	helpText := "[ ↑/↓ ] Navigate   [ p ] Pause   [ r ] Resume   [ c ] Cancel   [ y ] Try Again   [ a ] Add New"
	s.WriteString("\n" + helpStyle.Render(helpText))

	return s.String()
}

func renderQueueListTab(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("Queue Management"))
	s.WriteString("\n\n")

	if m.QueueFormMode {
		// Queue form
		s.WriteString(menuHeaderStyle.Render("Queue Configuration"))
		s.WriteString("\n\n")

		// Form fields
		fieldStyles := make([]lipgloss.Style, 6)
		for i := range fieldStyles {
			fieldStyles[i] = menuItemStyle
			if i == m.QueueFormField {
				fieldStyles[i] = selectedItemStyle
			}
		}

		s.WriteString(fieldStyles[0].Render("Name:           " + urlStyle.Render(m.InputQueueName)))
		s.WriteString("\n" + fieldStyles[1].Render("Path:           "+urlStyle.Render(m.InputQueuePath)))
		s.WriteString("\n" + fieldStyles[2].Render("Max Concurrent: "+urlStyle.Render(m.InputQueueConcurrent)))
		s.WriteString("\n" + fieldStyles[3].Render("Speed Limit:    "+urlStyle.Render(m.InputQueueSpeedLimit+" KB/s (0 = unlimited)")))
		s.WriteString("\n" + fieldStyles[4].Render("Start Time:     "+urlStyle.Render(m.InputQueueStartTime+" (format: HH:MM)")))
		s.WriteString("\n" + fieldStyles[5].Render("End Time:       "+urlStyle.Render(m.InputQueueEndTime+" (format: HH:MM)")))

		// Help text
		s.WriteString("\n\n" + helpStyle.Render("[ ↑/↓ ] Navigate   [ Tab ] Next Field   [ Enter ] Save   [ Esc ] Cancel"))
	} else {
		// Queue list
		if len(m.Config.Queues) == 0 {
			s.WriteString(menuItemStyle.Render("No queues configured. Press 'n' to add a queue."))
		} else {
			// Table header
			header := lipgloss.JoinHorizontal(lipgloss.Top,
				headerStyle.Width(20).Render("Name"),
				headerStyle.Width(25).Render("Path"),
				headerStyle.Width(15).Render("Max Concurrent"),
				headerStyle.Width(15).Render("Speed Limit"),
				headerStyle.Width(10).Render("Status"),
			)
			s.WriteString(header + "\n")

			// Queue list
			for i, q := range m.Config.Queues {
				// Highlight selected queue
				itemStyle := menuItemStyle
				if i == m.QueueSelected {
					itemStyle = selectedItemStyle
				}

				// Format speed limit
				speedLimit := "Unlimited"
				if q.SpeedLimit > 0 {
					speedLimit = fmt.Sprintf("%d KB/s", q.SpeedLimit)
				}

				// Format the line
				status := "Enabled"
				if !q.Enabled {
					status = "Disabled"
				}
				line := lipgloss.JoinHorizontal(lipgloss.Top,
					itemStyle.Width(20).Render(q.Name),
					itemStyle.Width(25).Render(truncateString(q.Path, 22)),
					itemStyle.Width(15).Render(fmt.Sprintf("%d", q.MaxConcurrent)),
					itemStyle.Width(15).Render(speedLimit),
					itemStyle.Width(10).Render(status),
				)
				s.WriteString(line + "\n")
			}
		}

		// Active downloads per queue
		s.WriteString("\n" + menuHeaderStyle.Render("Active Downloads Per Queue"))
		s.WriteString("\n")
		for _, q := range m.Config.Queues {
			activeCount := 0
			for _, d := range m.Downloads {
				if d.Queue == q.Name && d.Status == "downloading" {
					activeCount++
				}
			}
			s.WriteString(fmt.Sprintf("%s: %d/%d\n", q.Name, activeCount, q.MaxConcurrent))
		}

		// Help text
		s.WriteString("\n" + helpStyle.Render("[ ↑/↓ ] Navigate   [ n ] New Queue   [ e ] Edit Queue   [ d ] Delete Queue"))
	}

	return s.String()
}

func renderSettingsTab(m Model) string {
	var s strings.Builder
	s.WriteString(menuHeaderStyle.Render("Settings & Help"))
	s.WriteString("\n\n")

	// Theme settings
	s.WriteString(menuHeaderStyle.Render("Appearance"))
	s.WriteString("\n")
	s.WriteString(menuItemStyle.Render("Current Theme: " + m.CurrentTheme))
	s.WriteString("\n" + menuItemStyle.Render("Press 't' to cycle through available themes"))

	// Keyboard shortcuts
	s.WriteString("\n\n" + menuHeaderStyle.Render("Keyboard Shortcuts"))
	s.WriteString("\n")
	s.WriteString(menuItemStyle.Render("1-4:             Switch tabs"))
	s.WriteString("\n" + menuItemStyle.Render("↑/↓ or j/k:      Navigate lists"))
	s.WriteString("\n" + menuItemStyle.Render("Enter:           Confirm/Submit"))
	s.WriteString("\n" + menuItemStyle.Render("Esc:             Cancel/Back"))
	s.WriteString("\n" + menuItemStyle.Render("p:               Pause download"))
	s.WriteString("\n" + menuItemStyle.Render("r:               Resume download"))
	s.WriteString("\n" + menuItemStyle.Render("c:               Cancel download"))
	s.WriteString("\n" + menuItemStyle.Render("n:               New queue"))
	s.WriteString("\n" + menuItemStyle.Render("e:               Edit queue"))
	s.WriteString("\n" + menuItemStyle.Render("d:               Delete queue"))
	s.WriteString("\n" + menuItemStyle.Render("t:               Change theme"))
	s.WriteString("\n" + menuItemStyle.Render("q:               Quit application"))

	// About
	s.WriteString("\n\n" + menuHeaderStyle.Render("About"))
	s.WriteString("\n")
	s.WriteString(menuItemStyle.Render("Download Manager v0.1"))
	s.WriteString("\n" + menuItemStyle.Render("A terminal-based download manager with queue support"))

	return s.String()
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatSpeed(speed int64) string {
	if speed < 1024 {
		return fmt.Sprintf("%d B/s", speed)
	} else if speed < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", float64(speed)/1024)
	} else {
		return fmt.Sprintf("%.1f MB/s", float64(speed)/(1024*1024))
	}
}
