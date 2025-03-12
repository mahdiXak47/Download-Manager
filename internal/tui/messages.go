package tui

// Custom messages for our application
type StartDownloadMsg struct {
	URL   string
	Queue string
}

type TickMsg struct{}

type DownloadProgressMsg struct {
	URL      string
	Progress float64
	Speed    int64
}

type ErrorMsg struct {
	Error error
}
