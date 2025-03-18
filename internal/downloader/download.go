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
	cancelChan chan struct{} `json:"-"`
	mutex      sync.Mutex    `json:"-"`
	retryCount int           `json:"retry_count"`
	maxRetries int           `json:"max_retries"`
	retryDelay time.Duration `json:"-"`
}

// Initialize sets up control channels for a download
func (d *Download) Initialize() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.pauseChan == nil {
		d.pauseChan = make(chan struct{}, 1)
	}
	if d.resumeChan == nil {
		d.resumeChan = make(chan struct{}, 1)
	}
	if d.cancelChan == nil {
		d.cancelChan = make(chan struct{}, 1)
	}
	if d.Status == "" {
		d.Status = "pending"
	}
	if d.maxRetries == 0 {
		d.maxRetries = 3
	}
	if d.retryDelay == 0 {
		d.retryDelay = 5 * time.Second
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

// Cancel stops the download and removes temporary files
func (d *Download) Cancel() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status != "completed" && d.Status != "cancelled" {
		d.Status = "cancelled"
		select {
		case d.cancelChan <- struct{}{}:
		default:
		}

		fileName := filepath.Base(d.URL)
		if err := os.Remove(fileName); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove file: %v", err)
		}
	}
	return nil
}

// Retry attempts to restart a failed download
func (d *Download) Retry() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status == "error" {
		d.Status = "pending"
		d.Error = ""
		d.Progress = 0
		d.Speed = 0
		d.retryCount++
		return nil
	}
	return fmt.Errorf("download is not in error state")
}

// StartDownload begins downloading a file with optional speed limit
func StartDownload(d *Download, speedLimit int64) error {
	// Initialize control channels
	d.Initialize()
	d.Status = "downloading"

	for d.retryCount <= d.maxRetries {
		// Create HTTP client
		client := &http.Client{}

		// Create request
		req, err := http.NewRequest("GET", d.URL, nil)
		if err != nil {
			d.handleError(err)
			continue
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			d.handleError(err)
			continue
		}
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("server returned status %d", resp.StatusCode)
			d.handleError(err)
			continue
		}

		// Get file size
		fileSize := resp.ContentLength

		// Create output file
		fileName := filepath.Base(d.URL)
		file, err := os.Create(fileName)
		if err != nil {
			d.handleError(err)
			continue
		}
		defer file.Close()

		// Create buffer for speed limiting
		buffer := make([]byte, 32*1024) // 32KB buffer
		downloaded := int64(0)
		start := time.Now()

		for {
			select {
			case <-d.cancelChan:
				return nil
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
					d.handleError(err)
					goto retry
				}

				if n == 0 {
					if err == io.EOF {
						d.Status = "completed"
						return nil
					}
					continue
				}

				// Write to file
				if _, err := file.Write(buffer[:n]); err != nil {
					d.handleError(err)
					goto retry
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
	retry:
		// If we get here, an error occurred during download
		if d.retryCount >= d.maxRetries {
			d.Status = "error"
			d.Error = fmt.Sprintf("download failed after %d retries", d.maxRetries)
			return fmt.Errorf(d.Error)
		}
		d.retryCount++
		time.Sleep(d.retryDelay)
	}

	return fmt.Errorf("download failed after %d retries", d.maxRetries)
}

// handleError updates the download status and error message
func (d *Download) handleError(err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Status = "error"
	d.Error = err.Error()
}
