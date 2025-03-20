package tui

import (
	// "fmt"
	// "net/http"
	// "strings"
	"path/filepath"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/config"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
	"github.com/mahdiXak47/Download-Manager/internal/queue"
)

// TabID represents different tabs in the application
type TabID int

const (
	AddDownloadTab TabID = iota
	DownloadListTab
	QueueListTab
	SettingsTab
)

// Model represents the application state
type Model struct {
	// Core state
	ActiveTab  TabID  // Current active tab
	Menu       string // "main", "add", "list"
	InputMode  bool   // Whether we're capturing input
	Selected   int    // Currently selected download
	QueueSelected int // Currently selected queue

	// Input fields
	InputURL   string
	InputQueue string

	// Input fields for queue form
	InputQueueName    string
	InputQueuePath    string
	InputQueueConcurrent string
	InputQueueSpeedLimit string
	InputQueueStartTime  string
	InputQueueEndTime    string
	QueueFormMode     bool  // Whether we're in queue form mode
	QueueFormField    int   // Current field in queue form

	// Data
	Downloads    []downloader.Download
	Config       *config.Config
	QueueManager *queue.Manager
	ErrorMessage string

	// UI State
	Width  int
	Height int

	CurrentTheme string
}

// NewModel creates and initializes a new model
func NewModel() Model {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return Model{
			ActiveTab:   DownloadListTab,
			Menu:        "list",
			Downloads:   make([]downloader.Download, 0),
			Selected:    0,
			Width:       80,
			Height:      24,
			ErrorMessage: "Failed to load config: " + err.Error(),
		}
	}

	// Create queue manager
	queueManager := queue.NewManager(cfg)
	queueManager.Start()

	return Model{
		ActiveTab:    DownloadListTab,
		Menu:         "list",
		Downloads:    cfg.Downloads,
		Config:       cfg,
		QueueManager: queueManager,
		Selected:     0,
		QueueSelected: 0,
		Width:        80,
		Height:       24,
		CurrentTheme: "modern", // Default theme
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
		// Cancel input and clear fields
		if m.QueueFormMode {
			m.QueueFormMode = false
			m.InputQueueName = ""
			m.InputQueuePath = ""
			m.InputQueueConcurrent = ""
			m.InputQueueSpeedLimit = ""
			m.InputQueueStartTime = ""
			m.InputQueueEndTime = ""
			m.QueueFormField = 0
		} else {
			m.InputMode = false
			m.InputURL = ""
			m.InputQueue = ""
		}
		return *m, nil

	case tea.KeyBackspace:
		if m.QueueFormMode {
			switch m.QueueFormField {
			case 0:
				if len(m.InputQueueName) > 0 {
					m.InputQueueName = m.InputQueueName[:len(m.InputQueueName)-1]
				}
			case 1:
				if len(m.InputQueuePath) > 0 {
					m.InputQueuePath = m.InputQueuePath[:len(m.InputQueuePath)-1]
				}
			case 2:
				if len(m.InputQueueConcurrent) > 0 {
					m.InputQueueConcurrent = m.InputQueueConcurrent[:len(m.InputQueueConcurrent)-1]
				}
			case 3:
				if len(m.InputQueueSpeedLimit) > 0 {
					m.InputQueueSpeedLimit = m.InputQueueSpeedLimit[:len(m.InputQueueSpeedLimit)-1]
				}
			case 4:
				if len(m.InputQueueStartTime) > 0 {
					m.InputQueueStartTime = m.InputQueueStartTime[:len(m.InputQueueStartTime)-1]
				}
			case 5:
				if len(m.InputQueueEndTime) > 0 {
					m.InputQueueEndTime = m.InputQueueEndTime[:len(m.InputQueueEndTime)-1]
				}
			}
		} else if m.InputMode {
			if len(m.InputURL) > 0 {
				m.InputURL = m.InputURL[:len(m.InputURL)-1]
			}
		}
		return *m, nil

	case tea.KeyRunes:
		if m.QueueFormMode {
			switch m.QueueFormField {
			case 0:
				m.InputQueueName += string(msg.Runes)
			case 1:
				m.InputQueuePath += string(msg.Runes)
			case 2:
				m.InputQueueConcurrent += string(msg.Runes)
			case 3:
				m.InputQueueSpeedLimit += string(msg.Runes)
			case 4:
				m.InputQueueStartTime += string(msg.Runes)
			case 5:
				m.InputQueueEndTime += string(msg.Runes)
			}
		} else if m.InputMode {
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

	// Get the queue configuration to set bandwidth limit
	var maxBandwidth int64 = 0
	for _, q := range m.Config.Queues {
		if q.Name == queue {
			maxBandwidth = q.SpeedLimit
			break
		}
	}

	// Get a proper target path from the URL
	filename := filepath.Base(url)
	queuePath := m.Config.SavePath // Default to the global SavePath
	
	// Find the queue configuration
	for _, q := range m.Config.Queues {
		if q.Name == queue {
			// If the queue has a configured Path, use it
			if q.Path != "" {
				queuePath = q.Path
			}
			break
		}
	}
	
	var targetPath string
	if queuePath != "" {
		targetPath = filepath.Join(queuePath, filename)
	} else {
		// Use current directory if no path specified
		targetPath = filename
	}

	// Create and initialize download object
	download := downloader.New(url, targetPath, queue, maxBandwidth)
	m.Downloads = append(m.Downloads, *download)

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

// CancelDownload removes the selected download from the queue and downloads list
func (m *Model) CancelDownload() {
	if m.Selected >= 0 && m.Selected < len(m.Downloads) {
		download := m.Downloads[m.Selected]

		// Cancel the download if it's active
		if download.Status == "downloading" || download.Status == "paused" {
			download.Cancel()
		}

		// Remove from queue manager
		m.QueueManager.RemoveDownload(download.URL)

		// Remove from downloads list
		m.Downloads = append(m.Downloads[:m.Selected], m.Downloads[m.Selected+1:]...)

		// Adjust selection if needed
		if m.Selected >= len(m.Downloads) {
			m.Selected = len(m.Downloads) - 1
		}

		// Update config
		if m.Config != nil {
			m.Config.Downloads = m.Downloads
			if err := config.SaveConfig(m.Config); err != nil {
				m.ErrorMessage = "Failed to save config: " + err.Error()
			}
		}
	}
}

// CycleTheme switches to the next available theme
func (m *Model) CycleTheme() {
	themes := map[string]Theme{
		"modern":    ModernTheme,
		"ocean":     OceanTheme,
		"solarized": SolarizedTheme,
		"nord":      NordTheme,
	}

	// Get ordered theme names
	themeNames := []string{"modern", "ocean", "solarized", "nord"}
	currentIndex := 0
	for i, name := range themeNames {
		if name == m.CurrentTheme {
			currentIndex = i
			break
		}
	}

	// Switch to next theme
	nextIndex := (currentIndex + 1) % len(themeNames)
	m.CurrentTheme = themeNames[nextIndex]
	CurrentTheme = themes[m.CurrentTheme]
	UpdateStyles()
}

// SaveQueueForm saves the current queue form values to the config
func (m *Model) SaveQueueForm() error {
	// Convert string inputs to appropriate types
	maxConcurrent := 3 // Default
	if m.InputQueueConcurrent != "" {
		if val, err := strconv.Atoi(m.InputQueueConcurrent); err == nil && val > 0 {
			maxConcurrent = val
		}
	}

	var speedLimit int64 = 0 // Default - unlimited
	if m.InputQueueSpeedLimit != "" {
		if val, err := strconv.ParseInt(m.InputQueueSpeedLimit, 10, 64); err == nil && val >= 0 {
			speedLimit = val
		}
	}

	// Validate time formats
	startTime := "00:00" // Default
	if m.InputQueueStartTime != "" {
		if _, err := time.Parse("15:04", m.InputQueueStartTime); err == nil {
			startTime = m.InputQueueStartTime
		}
	}

	endTime := "23:59" // Default
	if m.InputQueueEndTime != "" {
		if _, err := time.Parse("15:04", m.InputQueueEndTime); err == nil {
			endTime = m.InputQueueEndTime
		}
	}

	// Create the queue config
	queue := config.QueueConfig{
		Name:          m.InputQueueName,
		Path:          m.InputQueuePath,
		MaxConcurrent: maxConcurrent,
		SpeedLimit:    speedLimit,
		StartTime:     startTime,
		EndTime:       endTime,
		Enabled:       true,
	}

	// Check if we're editing an existing queue or creating a new one
	found := false
	for i, q := range m.Config.Queues {
		if q.Name == m.InputQueueName {
			// Update existing queue
			m.Config.Queues[i] = queue
			found = true
			break
		}
	}

	if !found {
		// Add new queue
		m.Config.Queues = append(m.Config.Queues, queue)
	}

	// Save config
	return config.SaveConfig(m.Config)
}
