# Golang Download Manager

A terminal-based Download Manager built with Golang, featuring a text-based user interface (TUI). It supports multiple queues, download management, speed limits, and scheduling.

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


