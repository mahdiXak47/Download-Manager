package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	downloads  []Download
	selected   int
	inputMode  bool
	inputURL   string
	inputQueue string
	menu       string // "main", "add", "list", etc.
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.downloads = append(m.downloads, Download{
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

func (m model) View() string {
	s := "Download Manager v0.1\n\n"

	switch m.menu {
	case "add":
		s += "Add Download\n"
		s += "URL: " + m.inputURL + "\n"
		s += "Queue: " + m.inputQueue + "\n"
		s += "\nPress Enter to add, Esc to cancel"

	case "list":
		s += "Downloads:\n"
		for i, d := range m.downloads {
			cursor := " "
			if i == m.selected {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s [%s] %.1f%%\n", cursor, d.URL, d.Status, d.Progress)
		}

	default:
		s += "Commands:\n"
		s += "a: Add download\n"
		s += "l: List downloads\n"
		s += "p: Pause selected\n"
		s += "r: Resume selected\n"
		s += "q: Quit\n"
	}

	return s
}

func main() {
	p := tea.NewProgram(model{
		downloads: make([]Download, 0),
		menu:      "main",
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
