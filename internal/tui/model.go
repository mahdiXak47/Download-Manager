package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
)


type Model struct {
    Menu        string
    InputMode   bool   //A boolean to track whether the app is in input mode
    InputURL    string //Stores the URL entered by the user in the TUI.
    InputQueue  string //Stores the queue name entered by the user in the TUI.
    Downloads   []downloader.Download //A slice of downloader.Download structs to track ongoing or completed downloads
    Selected    int //The index of the currently selected download in the Downloads slice.
    ErrorMessage string // A string to display error messages or feedback to the user in the TUI.
}

func NewModel() Model {
	return Model{
		Downloads: make([]downloader.Download, 0),
		Menu:      "main",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
