package downloader

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

// MockDownload is a simplified version of Download for testing
type MockDownload struct {
	URL         string
	TargetPath  string
	Filename    string
	Queue       string
	Status      string
	Progress    float64
	TotalSize   int64
	Downloaded  int64
	Error       string
	MaxBandwidth int64
	maxRetries  int
	retryCount  int
	client      *http.Client
	mutex       sync.Mutex
}

// Initialize sets up the mock download
func (d *MockDownload) Initialize() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	if d.Status == "" {
		d.Status = "pending"
	}
	if d.maxRetries == 0 {
		d.maxRetries = 3
	}
	if d.client == nil {
		d.client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
}

// Start simulates starting a download
func (d *MockDownload) Start() error {
	d.Initialize()
	d.mutex.Lock()
	d.Status = "downloading"
	d.mutex.Unlock()
	
	// Simulate successful download
	time.Sleep(100 * time.Millisecond)
	
	d.mutex.Lock()
	d.Status = "completed"
	d.Progress = 100.0
	d.mutex.Unlock()
	
	return nil
}

// Pause simulates pausing a download
func (d *MockDownload) Pause() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	if d.Status == "downloading" {
		d.Status = "paused"
	}
}

// Resume simulates resuming a download
func (d *MockDownload) Resume() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	if d.Status == "paused" {
		d.Status = "downloading"
	}
}

// Cancel simulates cancelling a download
func (d *MockDownload) Cancel() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	d.Status = "cancelled"
	return nil
}

// Retry simulates retrying a download
func (d *MockDownload) Retry() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	if d.Status == "error" {
		d.Status = "pending"
		d.Error = ""
		d.Progress = 0
		d.Downloaded = 0
		d.retryCount++
		return nil
	}
	return fmt.Errorf("download is not in error state")
}

// GetStatus returns the current status
func (d *MockDownload) GetStatus() string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.Status
}

func TestMockDownload(t *testing.T) {
	t.Run("Basic Download", func(t *testing.T) {
		download := &MockDownload{
			URL:        "http://example.com/file.txt",
			TargetPath: "test.txt",
			Filename:   "test.txt",
			Queue:      "default",
			Status:     "pending",
		}
		
		err := download.Start()
		if err != nil {
			t.Fatalf("Download failed: %v", err)
		}
		
		if download.Status != "completed" {
			t.Errorf("Download status = %s, want completed", download.Status)
		}
	})
	
	t.Run("Pause and Resume Download", func(t *testing.T) {
		download := &MockDownload{
			URL:        "http://example.com/file.txt",
			TargetPath: "test.txt",
			Filename:   "test.txt",
			Queue:      "default",
			Status:     "downloading",
		}
		
		download.Pause()
		
		if download.Status != "paused" {
			t.Errorf("Download status after pause = %s, want paused", download.Status)
		}
		
		download.Resume()
		
		if download.Status != "downloading" {
			t.Errorf("Download status after resume = %s, want downloading", download.Status)
		}
	})
	
	t.Run("Cancel Download", func(t *testing.T) {
		download := &MockDownload{
			URL:        "http://example.com/file.txt",
			TargetPath: "test.txt",
			Filename:   "test.txt",
			Queue:      "default",
			Status:     "downloading",
		}
		
		err := download.Cancel()
		if err != nil {
			t.Errorf("Failed to cancel download: %v", err)
		}
		
		if download.Status != "cancelled" {
			t.Errorf("Download status after cancel = %s, want cancelled", download.Status)
		}
	})
	
	t.Run("Retry Download", func(t *testing.T) {
		download := &MockDownload{
			URL:        "http://example.com/file.txt",
			TargetPath: "test.txt",
			Filename:   "test.txt",
			Queue:      "default",
			Status:     "error",
			Error:      "test error",
		}
		
		err := download.Retry()
		if err != nil {
			t.Errorf("Failed to set download for retry: %v", err)
		}
		
		if download.Status != "pending" {
			t.Errorf("Download status after retry = %s, want pending", download.Status)
		}
		
		if download.Error != "" {
			t.Errorf("Error message still present after retry: %s", download.Error)
		}
	})
} 