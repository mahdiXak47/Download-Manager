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
	RetryCount int           `json:"retry_count"`
	MaxRetries int           `json:"max_retries"`
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
	if d.MaxRetries == 0 {
		d.MaxRetries = 3
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
		d.RetryCount++
		return nil
	}
	return fmt.Errorf("download is not in error state")
}

// StartDownload begins downloading a file with optional speed limit
func StartDownload(d *Download, speedLimit int64) error {
	// Initialize control channels
	d.Initialize()
	d.Status = "downloading"

	for d.RetryCount <= d.MaxRetries {
		// Create HTTP client
		client := &http.Client{}

		// Create request
		req, err := http.NewRequest("GET", d.URL, nil)
		if err != nil {
			d.handleError(err)
			continue
		}

		// Check if the file already exists and get its size
		fileName := filepath.Base(d.URL)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			d.handleError(err)
			continue
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			d.handleError(err)
			continue
		}

		// If the file already exists, set the Range header to resume the download
		if fileInfo.Size() > 0 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fileInfo.Size()))
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			d.handleError(err)
			continue
		}
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			err := fmt.Errorf("server returned status %d", resp.StatusCode)
			d.handleError(err)
			continue
		}

		// Get file size
		fileSize := resp.ContentLength + fileInfo.Size() // Total size = downloaded + remaining

		// Seek to the end of the file to append new data
		file.Seek(fileInfo.Size(), io.SeekStart)

		// Create buffer for speed limiting
		buffer := make([]byte, 32*1024) // 32KB buffer
		downloaded := fileInfo.Size()
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
			default:
				// Normal download operation
				d.limitSpeed(buffer, speedLimit)

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
		if d.RetryCount >= d.MaxRetries {
			d.Status = "error"
			d.Error = fmt.Sprintf("download failed after %d retries", d.MaxRetries)
			return fmt.Errorf(d.Error)
		}
		d.RetryCount++
		time.Sleep(d.retryDelay)
	}

	return fmt.Errorf("download failed after %d retries", d.MaxRetries)
}

// limitSpeed implements a simple speed limiter
func (d *Download) limitSpeed(buffer []byte, speedLimit int64) {
	if speedLimit <= 0 {
		return
	}

	// Calculate the time required to read the buffer at the given speed limit
	chunkSize := int64(len(buffer))
	expectedTime := time.Duration(float64(chunkSize)/float64(speedLimit)) * time.Second

	// Sleep for the required duration
	time.Sleep(expectedTime)
}

// handleError updates the download status and error message
func (d *Download) handleError(err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Status = "error"
	d.Error = err.Error()
}
