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
  - Support for various protocols (HTTP, HTTPS, FTP)

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

## Project Architecture

```
.
├── cmd/                    # Application entry points
│   └── download-manager/   # Main application
├── internal/              # Private application code
│   ├── tui/              # Terminal UI components
│   ├── downloader/       # Download management
│   ├── queue/            # Queue management
│   ├── scheduler/        # Time scheduling
│   └── config/           # Configuration management
├── pkg/                  # Public libraries
│   ├── protocol/         # Download protocols
│   └── utils/            # Utility functions
├── configs/              # Configuration files
└── docs/                 # Documentation
```

## Technical Stack

- Go 1.21 or higher
- Key Libraries (Planned):
  - [bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
  - [lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
  - [cobra](https://github.com/spf13/cobra) - CLI commands
  - [viper](https://github.com/spf13/viper) - Configuration management

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
```

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
- **Tab 1**: Add new downloads.
- **Tab 2**: View & manage active downloads.
- **Tab 3**: Configure and manage download queues.

## Technical Highlights
- **Concurrency**: Utilizes Goroutines and Channels.
- **Error Handling & Retry Logic**.
- **Networking & File I/O**: Handles HTTP GET requests and file writes efficiently.
- **TUI Library**: Uses `tview` for terminal UI.