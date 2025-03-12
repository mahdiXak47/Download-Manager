package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			if !m.InputMode {
				m.Menu = "add"
				m.InputMode = true
			}
		case "l":
			m.Menu = "list"
			m.InputMode = false
		case "p":
			if m.Selected >= 0 && m.Selected < len(m.Downloads) {
				m.Downloads[m.Selected].Status = "paused"
			}
		case "r":
			if m.Selected >= 0 && m.Selected < len(m.Downloads) {
				m.Downloads[m.Selected].Status = "downloading"
			}
		case "enter":
			if m.InputMode && m.Menu == "add" {
				m.Downloads = append(m.Downloads, downloader.Download{
					URL:      m.InputURL,
					Queue:    m.InputQueue,
					Status:   "pending",
					Progress: 0,
				})
				m.InputMode = false
				m.InputURL = ""
				m.Menu = "main"
			}
		}
	}
	return m, nil
}
