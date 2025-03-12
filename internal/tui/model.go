package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
)

type Model struct {
	downloads  []downloader.Download
	selected   int
	inputMode  bool
	inputURL   string
	inputQueue string
	menu       string // "main", "add", "list", etc.
}

func NewModel() Model {
	return Model{
		downloads: make([]downloader.Download, 0),
		menu:      "main",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
