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
  - Immediate download processing when added to queue
  - Optimized HTTP client settings for reliability
  - Proper error handling and retry logic
  - Bandwidth limiting capabilities
  - Progress tracking and speed calculation
  - File directory management
  - Multi-part downloading for large files (when servers support Range headers)

- **Multi-part Download Support**: Added parallel downloading capability:

  - Automatically splits large files (default >10MB) into multiple parts
  - Downloads parts in parallel for faster speeds
  - Intelligently falls back to single-part when servers don't support ranges
  - Configurable number of parts (default 5) and size threshold
  - Ensures parts are correctly assembled into a complete file

- **Comprehensive Logging System**: Added detailed logging functionality:

  - Records all download activities to `download-logs.log`
  - Tracks download starts, status changes, and completions
  - Logs errors with detailed reasons
  - Queue operations and system events logging
  - Allows troubleshooting of download issues

- **Enhanced Error Recovery**: Added manual retry functionality:

  - New "Try Again" option (key: 'y') for failed downloads
  - Clear visual feedback with color-coded messages
  - Limited to 3 retry attempts per download
  - Automatically processes retried downloads when queue capacity allows

- **Test Suite Implementation**: Added comprehensive testing capabilities:
  - Unit tests for logger functionality
  - Unit tests for downloader core components
  - Rate limiter testing with timing verification
  - Mock HTTP servers for controlled download testing
  - Pause/resume/cancel functionality tests
  - Multi-part download testing

## Project Architecture

```
.
├── cmd/
│   ├── main.go                # Application entry point
│   ├── main_test.go           # Main package tests
│   └── test_download/         # Standalone test download utility
│       └── main.go            # Test download implementation
├── internal/
│   ├── tui/
│   │   ├── model.go           # Core TUI model
│   │   ├── update.go          # Update logic for TUI
│   │   ├── view.go            # View rendering for TUI
│   │   ├── styles.go          # Style definitions
│   │   ├── themes.go          # Theme definitions
│   │   └── messages.go        # Message definitions
│   ├── downloader/
│   │   ├── download.go        # Download implementation
│   │   ├── ratelimiter.go     # Rate limiting functionality
│   │   ├── downloader_test.go # Downloader tests
│   │   └── ratelimiter_test.go # Rate limiter tests
│   ├── queue/
│   │   └── manager.go         # Queue management system
│   ├── logger/
│   │   ├── logger.go          # Logging system
│   │   └── logger_test.go     # Logger tests
│   └── config/
│       └── config.go          # Configuration management
├── downloads/                 # Download destination folder
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── README.md                  # Project documentation
```

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
   - Testing and quality assurance

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

# Run tests
go test ./...

# Run the test download utility (for testing download functionality independently)
go run cmd/test_download/main.go
```

## Features

- **Concurrent Downloads**: Uses Goroutines and Channels for efficient multi-threading.
- **Download Queue Management**: Set storage folder, max simultaneous downloads, bandwidth limits, and time scheduling.
- **Control Options**: Pause, Resume, Cancel, and Retry downloads.
- **Multi-part Downloading**: Supports parallel downloads for large files if the server allows `Accept-Ranges`.
- **Persistent State**: Saves queues and downloads, resuming unfinished ones on restart.
- **Keyboard Shortcuts**: Navigate tabs, control downloads, and manage queues efficiently.
- **Activity Logging**: Records all download activities and errors to a log file for troubleshooting.

## User Interface

The Download Manager features a clean, intuitive terminal user interface with tabbed navigation:

- **Tab 1**: Add new downloads - Enter URL and choose queue.
- **Tab 2**: Download List - View and manage active downloads (pause, resume, cancel).
- **Tab 3**: Queue Management - Configure and manage download queues.
- **Tab 4**: Settings & Help - Change themes and view keyboard shortcuts.

The interface supports keyboard navigation with tabs displayed at the bottom of the screen for easy access.

## Keyboard Shortcuts

- **1-4**: Switch between tabs (works globally when not in input mode)
- **↑/↓** or **j/k**: Navigate lists
- **Enter**: Confirm/Submit
- **Esc**: Cancel/Back or exit input mode
- **p**: Pause selected download
- **r**: Resume selected download
- **c**: Cancel selected download
- **y**: Try again for failed downloads (limited to 3 attempts)
- **n**: Add new queue (in Queue tab)
- **e**: Edit selected queue (in Queue tab)
- **d**: Delete selected queue (in Queue tab)
- **t**: Change theme (press when not typing in an input field)
- **q**: Quit application

## Testing

The download manager includes a comprehensive test suite:

- **Unit Tests**: Test individual components for correctness
- **Integration Tests**: Test component interaction
- **Mock Testing**: Using httptest package to simulate HTTP servers
- **Performance Testing**: Rate limiter timing verification

To run all tests:

```bash
go test ./...
```

To run tests for a specific package:

```bash
go test ./internal/logger
go test ./internal/downloader
```

## Technical Highlights

- **Concurrency**: Utilizes Goroutines and Channels.
- **Error Handling & Retry Logic**: Implements robust retry mechanism with configurable attempts and delays.
- **Networking & File I/O**: Handles HTTP requests and file writes efficiently.
- **TUI Library**: Uses `bubbletea` for terminal UI.
