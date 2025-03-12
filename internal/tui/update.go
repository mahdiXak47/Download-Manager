package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// First check if we're in input mode
		if m.InputMode {
			model, cmd := m.HandleInput(msg)
			return model, cmd
		}

		// If not in input mode, handle navigation
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.Menu = "add"
			m.InputMode = true
		case "l":
			m.Menu = "list"
		case "p":
			m.PauseDownload()
		case "r":
			m.ResumeDownload()
		case "j", "down":
			if m.Selected < len(m.Downloads)-1 {
				m.Selected++
			}
		case "k", "up":
			if m.Selected > 0 {
				m.Selected--
			}
		case "enter":
			if m.InputMode && m.Menu == "add" {
				m.AddDownload(m.InputURL, m.InputQueue)
				m.InputMode = false
				m.InputURL = ""
				m.Menu = "main"
			}
		}
	}
	return m, nil
}
