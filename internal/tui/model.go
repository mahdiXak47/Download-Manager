package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/config"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
	"github.com/mahdiXak47/Download-Manager/internal/queue"
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
	Downloads    []downloader.Download
	Config       *config.Config
	QueueManager *queue.Manager
	ErrorMessage string

	// UI State
	Width  int
	Height int
}

// NewModel creates and initializes a new model
func NewModel() Model {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return Model{
			Menu:         "list",
			Downloads:    make([]downloader.Download, 0),
			Selected:     0,
			Width:        80,
			Height:       24,
			ErrorMessage: "Failed to load config: " + err.Error(),
		}
	}

	// Create queue manager
	queueManager := queue.NewManager(cfg)
	queueManager.Start()

	return Model{
		Menu:         "list",
		Downloads:    cfg.Downloads,
		Config:       cfg,
		QueueManager: queueManager,
		Selected:     0,
		Width:        80,
		Height:       24,
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
	if queue == "" {
		queue = m.Config.DefaultQueue
	}

	download := downloader.Download{
		URL:      url,
		Queue:    queue,
		Status:   "pending",
		Progress: 0,
	}
	m.Downloads = append(m.Downloads, download)

	// Update config with new download
	if m.Config != nil {
		m.Config.Downloads = m.Downloads
		if err := config.SaveConfig(m.Config); err != nil {
			m.ErrorMessage = "Failed to save config: " + err.Error()
		}
	}
}

// PauseDownload pauses the selected download
func (m *Model) PauseDownload() {
	if m.Selected >= 0 && m.Selected < len(m.Downloads) {
		download := &m.Downloads[m.Selected]
		if download.Status == "downloading" {
			m.QueueManager.PauseDownload(download.URL)
		}
	}
}

// ResumeDownload resumes the selected download
func (m *Model) ResumeDownload() {
	if m.Selected >= 0 && m.Selected < len(m.Downloads) {
		download := &m.Downloads[m.Selected]
		if download.Status == "paused" {
			m.QueueManager.ResumeDownload(download.URL)
		}
	}
}
