package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mahdiXak47/Download-Manager/internal/logger"
	"github.com/mahdiXak47/Download-Manager/internal/tui"
)

func main() {
	logFile := "download-logs.log"
	cwd, err := os.Getwd()
	if err == nil {
    	logFile = filepath.Join(cwd, "download-logs.log")
	}


	if err := logger.Initialize(logFile); err != nil {
		fmt.Printf("Warning: Could not initialize logger: %v\n", err)
	}

	p := tea.NewProgram(tui.NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	logger.Close()
}
