package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all state updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return handleWindowSize(m, msg)
	case tea.KeyMsg:
		return handleKeyPress(m, msg)
	case TickMsg:
		return handleTick(m)
	case StartDownloadMsg:
		return handleStartDownload(m, msg)
	case DownloadProgressMsg:
		return handleProgress(m, msg)
	case ErrorMsg:
		return handleError(m, msg)
	}
	return m, nil
}

// handleWindowSize updates the terminal size
func handleWindowSize(m Model, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.UpdateSize(msg.Width, msg.Height)
	return m, nil
}

// handleKeyPress handles all keyboard input
func handleKeyPress(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle global keys first
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		if m.Menu != "main" {
			m.Menu = "main"
			m.InputMode = false
			return m, nil
		}
	case "t":
		if m.Menu == "main" {
			m.CycleTheme()
		}
	}

	// Handle menu-specific keys
	if m.InputMode {
		return handleInputMode(m, msg)
	}

	return handleNavigationMode(m, msg)
}

// handleInputMode handles keyboard input when in input mode
func handleInputMode(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.Menu == "add" && m.InputURL != "" {
			cmd := func() tea.Msg {
				return StartDownloadMsg{
					URL:   m.InputURL,
					Queue: m.InputQueue,
				}
			}
			m.InputMode = false
			m.InputURL = ""
			m.InputQueue = ""
			m.Menu = "list"
			return m, cmd
		}

	case tea.KeyBackspace:
		if len(m.InputURL) > 0 {
			m.InputURL = m.InputURL[:len(m.InputURL)-1]
		}

	case tea.KeyRunes:
		m.InputURL += string(msg.Runes)
	}

	return m, nil
}

// Handles navigation (e.g., switching menus, pausing/resuming downloads).
func handleNavigationMode(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.Menu == "list" && len(m.Downloads) > 0 {
			if m.Selected > 0 {
				m.Selected--
			} else {
				m.Selected = len(m.Downloads) - 1
			}
		}
	case tea.KeyDown:
		if m.Menu == "list" && len(m.Downloads) > 0 {
			m.Selected = (m.Selected + 1) % len(m.Downloads)
		}
	}

	switch msg.String() {
	case "a":
		m.Menu = "add"
		m.InputMode = true
	case "l":
		m.Menu = "list"
	case "p":
		if m.Menu == "list" {
			m.PauseDownload()
		}
	case "r":
		if m.Menu == "list" {
			m.ResumeDownload()
			return m, tickCmd()
		}
	case "c":
		if m.Menu == "list" {
			m.CancelDownload()
		}
	case "j":
		if m.Menu == "list" && len(m.Downloads) > 0 {
			m.Selected = (m.Selected + 1) % len(m.Downloads)
		}
	case "k":
		if m.Menu == "list" && len(m.Downloads) > 0 {
			if m.Selected > 0 {
				m.Selected--
			} else {
				m.Selected = len(m.Downloads) - 1
			}
		}
	}

	return m, nil
}

// handleStartDownload processes a new download request
func handleStartDownload(m Model, msg StartDownloadMsg) (tea.Model, tea.Cmd) {
	m.AddDownload(msg.URL, msg.Queue)
	return m, tickCmd()
}

// handleProgress updates download progress
func handleProgress(m Model, msg DownloadProgressMsg) (tea.Model, tea.Cmd) {
	for i, d := range m.Downloads {
		if d.URL == msg.URL {
			m.Downloads[i].Progress = msg.Progress
			m.Downloads[i].Speed = msg.Speed
			break
		}
	}
	return m, nil
}

// handleError displays error messages
func handleError(m Model, msg ErrorMsg) (tea.Model, tea.Cmd) {
	m.ErrorMessage = msg.Error.Error()
	return m, nil
}

// Handles periodic updates (e.g., checking for active downloads).
func handleTick(m Model) (tea.Model, tea.Cmd) {
	// Update active downloads
	hasActive := false
	for _, d := range m.Downloads {
		if d.Status == "downloading" || d.Status == "paused" {
			hasActive = true
			break
		}
	}

	if hasActive {
		return m, tickCmd()
	}
	return m, nil
}

// Schedules a periodic update
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}
