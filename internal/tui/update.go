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
			if !m.inputMode {
				m.menu = "add"
				m.inputMode = true
			}
		case "l":
			m.menu = "list"
			m.inputMode = false
		case "p":
			if m.selected >= 0 && m.selected < len(m.downloads) {
				m.downloads[m.selected].Status = "paused"
			}
		case "r":
			if m.selected >= 0 && m.selected < len(m.downloads) {
				m.downloads[m.selected].Status = "downloading"
			}
		case "enter":
			if m.inputMode && m.menu == "add" {
				m.downloads = append(m.downloads, downloader.Download{
					URL:      m.inputURL,
					Queue:    m.inputQueue,
					Status:   "pending",
					Progress: 0,
				})
				m.inputMode = false
				m.inputURL = ""
				m.menu = "main"
			}
		}
	}
	return m, nil
}
