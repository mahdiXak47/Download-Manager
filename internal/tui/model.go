package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
)

// Model represents the application state
type Model struct {
	// Core state
	Menu      string // "main", "add", "list"
	InputMode bool   // Whether we're capturing input
	Selected  int    // Currently selected download

	// Input fields
	InputURL   string
	InputQueue string

	// Data
	Downloads []downloader.Download

	// UI State
	ErrorMessage string
	Width        int
	Height       int
}

// NewModel creates and initializes a new model
func NewModel() Model {
	return Model{
		Menu:      "main",
		Downloads: make([]downloader.Download, 0),
		Selected:  0,
		Width:     80, // default terminal width
		Height:    24, // default terminal height
	}
}

// Init runs any initial IO
func (m Model) Init() tea.Cmd {
	return nil
}

// HandleInput processes text input when in input mode
func (m *Model) HandleInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.InputMode = false
		m.Menu = "main"
		return *m, nil

	case tea.KeyBackspace:
		if m.InputMode {
			if len(m.InputURL) > 0 {
				m.InputURL = m.InputURL[:len(m.InputURL)-1]
			}
		}
		return *m, nil

	case tea.KeyRunes:
		if m.InputMode {
			m.InputURL += string(msg.Runes)
		}
		return *m, nil
	}

	return *m, nil
}

// UpdateSize updates the model's terminal size
func (m *Model) UpdateSize(width, height int) {
	m.Width = width
	m.Height = height
}

// AddDownload adds a new download to the model
func (m *Model) AddDownload(url, queue string) {
	download := downloader.Download{
		URL:      url,
		Queue:    queue,
		Status:   "pending",
		Progress: 0,
	}
	m.Downloads = append(m.Downloads, download)
}

// PauseDownload pauses the selected download
func (m *Model) PauseDownload() {
	if m.Selected >= 0 && m.Selected < len(m.Downloads) {
		m.Downloads[m.Selected].Status = "paused"
	}
}

// ResumeDownload resumes the selected download
func (m *Model) ResumeDownload() {
	if m.Selected >= 0 && m.Selected < len(m.Downloads) {
		m.Downloads[m.Selected].Status = "downloading"
	}
}
