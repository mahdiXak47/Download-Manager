# Download Manager

A sophisticated download manager application written in Go with Terminal User Interface (TUI).

## Features

- Terminal User Interface (TUI)

  - Interactive command-line interface
  - Real-time download progress visualization
  - Queue management interface
  - System status dashboard

- Queue Management

  - Multiple download queues support
  - Queue prioritization
  - Queue-specific settings
  - Pause/Resume entire queues

- Download Management

  - Concurrent downloads
  - Progress tracking
  - Pause/Resume individual downloads
  - Download retry mechanism
  - Support for HTTP/HTTPS protocols
  - Download speed limiting

- Speed Control

  - Global speed limit
  - Per-queue speed limits
  - Per-download speed limits
  - Bandwidth allocation

- Scheduling
  - Time-based download scheduling
  - Queue-based scheduling
  - Priority-based scheduling
  - Bandwidth scheduling

## Recent Updates

- **Improved Download Engine**: The download functionality has been fully implemented with:
  - Robust pause/resume support using HTTP Range headers
  - Proper error handling and retry logic
  - Bandwidth limiting capabilities
  - Progress tracking and speed calculation
  - File directory management

## Project Architecture

````
.
├── cmd/
│   └── download-manager/
│       └── main.go
├── internal/
│   ├── tui/
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   ├── components/
│   │   │   ├── progress.go
│   │   │   ├── queue.go
│   │   │   ├── status.go
│   │   │   └── tabs.go
│   │   ├── constants.go
│   │   └── messages.go
│   ├── downloader/
│   │   └── download.go     # Consolidated download implementation
│   ├── queue/
│   │   └── manager.go      # Queue management system
│   └── config/
│       └── config.go       # Configuration management
├── pkg/
│   ├── protocol/
│   └── utils/
├── configs/
├── docs/
│   ├── architecture.md
│   ├── api.md
│   ├── user-guide.md
│   └── dev-guide.md
├── logs/
│   └── download.log
├── scripts/
│   ├── install.sh
│   ├── setup.sh
│   └── test.sh
├── tests/
│   ├── downloader/
│   ├── queue/
│   ├── scheduler/
│   └── tui/
├── Makefile
├── Dockerfile
├── go.mod
├── go.sum
├── .gitignore
├── CHANGELOG.md
├── CONTRIBUTING.md
└── README.md

## Technical Stack

- Go 1.21 or higher
- Key Libraries:
  - [bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
  - [lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## Development Phases

1. **Foundation Phase**

   - Project structure setup
   - Basic CLI framework
   - Configuration management

2. **Core Components Phase**

   - Download engine implementation
   - Queue management system
   - Speed control mechanism

3. **TUI Phase**

   - Basic TUI implementation
   - Progress visualization
   - Interactive controls

4. **Advanced Features Phase**

   - Scheduling system
   - Advanced queue management
   - Protocol handlers

5. **Polish Phase**
   - Error handling
   - Logging system
   - Performance optimization

## Installation

```bash
# Clone the repository
git clone [your-repo-url]

# Change to project directory
cd download-manager

# Build the application
go build
````

## Usage

```bash
# Run the application
./download-manager
```

## Features

- **Concurrent Downloads**: Uses Goroutines and Channels for efficient multi-threading.
- **Download Queue Management**: Set storage folder, max simultaneous downloads, bandwidth limits, and time scheduling.
- **Control Options**: Pause, Resume, Cancel, and Retry downloads.
- **Multi-part Downloading**: Supports parallel downloads for large files if the server allows `Accept-Ranges`.
- **Persistent State**: Saves queues and downloads, resuming unfinished ones on restart.
- **Keyboard Shortcuts**: Navigate tabs, control downloads, and manage queues efficiently.

## User Interface

The Download Manager features a clean, intuitive terminal user interface with tabbed navigation:

- **Tab 1 (F1)**: Add new downloads - Enter URL and choose queue.
- **Tab 2 (F2)**: Download List - View and manage active downloads (pause, resume, cancel).
- **Tab 3 (F3)**: Queue Management - Configure and manage download queues.
- **Tab 4 (F4)**: Settings & Help - Change themes and view keyboard shortcuts.

The interface supports keyboard navigation with tabs displayed at the bottom of the screen for easy access.

## Keyboard Shortcuts

- **F1-F4**: Switch between tabs (function keys work globally)
- **↑/↓** or **j/k**: Navigate lists
- **Enter**: Confirm/Submit
- **Esc**: Cancel/Back or exit input mode
- **p**: Pause selected download
- **r**: Resume selected download
- **c**: Cancel selected download
- **n**: Add new queue (in Queue tab)
- **e**: Edit selected queue (in Queue tab)
- **d**: Delete selected queue (in Queue tab)
- **t**: Change theme (press when not typing in an input field)
- **q**: Quit application

## Technical Highlights

- **Concurrency**: Utilizes Goroutines and Channels.
- **Error Handling & Retry Logic**: Implements robust retry mechanism with configurable attempts and delays.
- **Networking & File I/O**: Handles HTTP requests and file writes efficiently.
- **TUI Library**: Uses `bubbletea` for terminal UI.
