package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/config"
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
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		if m.QueueFormMode {
			m.QueueFormMode = false
			return m, nil
		}
		if m.InputMode {
			m.InputMode = false
			return m, nil
		}
	case tea.KeyF1:
		// Switch to Add Download tab
		m.ActiveTab = AddDownloadTab
		m.Menu = "add"
		return m, nil
	case tea.KeyF2:
		// Switch to Download List tab
		m.ActiveTab = DownloadListTab
		m.Menu = "list"
		return m, nil
	case tea.KeyF3:
		// Switch to Queue List tab
		m.ActiveTab = QueueListTab
		m.Menu = "queues"
		return m, nil
	case tea.KeyF4:
		// Switch to Settings tab
		m.ActiveTab = SettingsTab
		m.Menu = "settings"
		return m, nil
	}

	// Handle rune keys (when not in input mode)
	if !m.InputMode && !m.QueueFormMode {
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "t":
			// Allow theme change from any tab when not in input mode
			m.CycleTheme()
			return m, nil
		}
	}

	// Handle tab-specific keys
	switch m.ActiveTab {
	case AddDownloadTab:
		return handleAddDownloadTab(m, msg)
	case DownloadListTab:
		return handleDownloadListTab(m, msg)
	case QueueListTab:
		return handleQueueListTab(m, msg)
	case SettingsTab:
		return handleSettingsTab(m, msg)
	}

	return m, nil
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

// handleAddDownloadTab handles keys for the Add Download tab
func handleAddDownloadTab(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.InputMode {
		// Handle input mode separately
		m, cmd := m.HandleInput(msg)
		return m, cmd
	}

	switch msg.String() {
	case "enter":
		m.InputMode = true
		m.InputURL = ""
		return m, nil
	}

	return m, nil
}

// handleDownloadListTab handles keys for the Download List tab
func handleDownloadListTab(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down", "j":
		if m.Selected < len(m.Downloads)-1 {
			m.Selected++
		}
	case "p":
		m.PauseDownload()
	case "r":
		m.ResumeDownload()
	case "c":
		m.CancelDownload()
	case "a":
		// Switch to Add Download tab
		m.ActiveTab = AddDownloadTab
		m.Menu = "add"
	}

	return m, nil
}

// handleQueueListTab handles keys for the Queue List tab
func handleQueueListTab(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.QueueFormMode {
		// Handle queue form input
		switch msg.String() {
		case "up", "shift+tab":
			if m.QueueFormField > 0 {
				m.QueueFormField--
			}
		case "down", "tab":
			if m.QueueFormField < 5 { // 6 fields total (0-5)
				m.QueueFormField++
			}
		case "enter":
			if m.QueueFormField < 5 {
				// Move to next field
				m.QueueFormField++
			} else {
				// Submit form
				m.SaveQueueForm()
				m.QueueFormMode = false
			}
		default:
			// Handle text input
			m, cmd := m.HandleInput(msg)
			return m, cmd
		}
		return m, nil
	}

	// Not in form mode
	switch msg.String() {
	case "up", "k":
		if m.QueueSelected > 0 {
			m.QueueSelected--
		}
	case "down", "j":
		if m.QueueSelected < len(m.Config.Queues)-1 {
			m.QueueSelected++
		}
	case "n":
		// New queue
		m.QueueFormMode = true
		m.QueueFormField = 0
		m.InputQueueName = ""
		m.InputQueuePath = ""
		m.InputQueueConcurrent = "3"
		m.InputQueueSpeedLimit = "0"
		m.InputQueueStartTime = "00:00"
		m.InputQueueEndTime = "23:59"
	case "e":
		// Edit queue
		if m.QueueSelected >= 0 && m.QueueSelected < len(m.Config.Queues) {
			q := m.Config.Queues[m.QueueSelected]
			m.QueueFormMode = true
			m.QueueFormField = 0
			m.InputQueueName = q.Name
			m.InputQueuePath = q.Path
			m.InputQueueConcurrent = fmt.Sprintf("%d", q.MaxConcurrent)
			m.InputQueueSpeedLimit = fmt.Sprintf("%d", q.SpeedLimit)
			m.InputQueueStartTime = q.StartTime
			m.InputQueueEndTime = q.EndTime
		}
	case "d":
		// Delete queue
		if m.QueueSelected >= 0 && m.QueueSelected < len(m.Config.Queues) {
			// Don't delete default queue
			if m.Config.Queues[m.QueueSelected].Name != m.Config.DefaultQueue {
				m.Config.Queues = append(m.Config.Queues[:m.QueueSelected], m.Config.Queues[m.QueueSelected+1:]...)
				if m.QueueSelected >= len(m.Config.Queues) {
					m.QueueSelected = len(m.Config.Queues) - 1
				}
				config.SaveConfig(m.Config)
			}
		}
	}

	return m, nil
}

// handleSettingsTab handles keys for the Settings tab
func handleSettingsTab(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// No need to handle 't' here as it's handled globally
	return m, nil
}
