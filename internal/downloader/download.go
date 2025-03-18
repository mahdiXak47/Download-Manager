package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Download struct {
	URL      string  `json:"url"`
	Queue    string  `json:"queue"`
	Status   string  `json:"status"` // pending, downloading, paused, completed, error
	Progress float64 `json:"progress"`
	Speed    int64   `json:"speed"`
	Error    string  `json:"error"`

	// Control fields (not persisted to JSON)
	pauseChan  chan struct{} `json:"-"`
	resumeChan chan struct{} `json:"-"`
	mutex      sync.Mutex    `json:"-"`
}

// Initialize sets up control channels for a download
func (d *Download) Initialize() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.pauseChan == nil {
		d.pauseChan = make(chan struct{})
	}
	if d.resumeChan == nil {
		d.resumeChan = make(chan struct{})
	}
}

// Pause signals the download to pause
func (d *Download) Pause() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status == "downloading" {
		d.Status = "paused"
		d.pauseChan <- struct{}{}
	}
}

// Resume signals the download to resume
func (d *Download) Resume() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status == "paused" {
		d.Status = "downloading"
		d.resumeChan <- struct{}{}
	}
}

// StartDownload begins downloading a file with optional speed limit
func StartDownload(d *Download, speedLimit int64) error {
	// Initialize control channels
	d.Initialize()

	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Get file size
	fileSize := resp.ContentLength

	// Create output file
	fileName := filepath.Base(d.URL)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Create buffer for speed limiting
	buffer := make([]byte, 32*1024) // 32KB buffer
	downloaded := int64(0)
	start := time.Now()

	for {
		select {
		case <-d.pauseChan:
			// Wait for resume signal
			<-d.resumeChan
			// Reset speed calculation after resume
			start = time.Now()
			downloaded = int64(float64(fileSize) * d.Progress / 100)

		default:
			// Normal download operation
			if speedLimit > 0 && d.Speed > speedLimit {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			n, err := resp.Body.Read(buffer)
			if err != nil && err != io.EOF {
				return fmt.Errorf("failed to read response: %v", err)
			}

			if n == 0 {
				break
			}

			// Write to file
			if _, err := file.Write(buffer[:n]); err != nil {
				return fmt.Errorf("failed to write to file: %v", err)
			}

			// Update progress
			downloaded += int64(n)
			if fileSize > 0 {
				d.Progress = float64(downloaded) * 100 / float64(fileSize)
			}

			// Calculate speed
			elapsed := time.Since(start).Seconds()
			if elapsed > 0 {
				d.Speed = int64(float64(downloaded) / elapsed)
			}
		}
	}

	return nil
}
