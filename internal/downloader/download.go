package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/mahdiXak47/Download-Manager/internal/logger"
)

// Download represents a download task with its state and control channels
type Download struct {
	URL          string    `json:"url"`
	TargetPath   string    `json:"target_path"`
	Filename     string    `json:"filename"`
	Queue        string    `json:"queue"`
	Status       string    `json:"status"` // pending, downloading, paused, completed, error, cancelled
	Progress     float64   `json:"progress"`
	Speed        int64     `json:"speed"` // bytes per second
	TotalSize    int64     `json:"total_size"`
	Downloaded   int64     `json:"downloaded"`
	Error        string    `json:"error,omitempty"`
	MaxBandwidth int64     `json:"max_bandwidth"` // in KB/s, 0 means unlimited
	StartTime    time.Time `json:"start_time,omitempty"`

	// Control fields (not persisted to JSON)
	pauseChan   chan struct{} `json:"-"`
	resumeChan  chan struct{} `json:"-"`
	cancelChan  chan struct{} `json:"-"`
	isPaused    bool          `json:"-"`
	isCancelled bool          `json:"-"`
	mutex       sync.Mutex    `json:"-"`
	retryCount  int           `json:"retry_count"`
	maxRetries  int           `json:"max_retries"`
	retryDelay  time.Duration `json:"-"`
	client      *http.Client  `json:"-"`
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
		// Log the initial pending status
		logger.LogDownloadPending(d.URL, d.Queue, "Initialized download")
	}
	if d.maxRetries == 0 {
		d.maxRetries = 3
	}
	if d.retryDelay == 0 {
		d.retryDelay = 5 * time.Second
	}
	if d.client == nil {
		// Configure HTTP client with reasonable timeouts and settings
		d.client = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConns:          10,
				IdleConnTimeout:       30 * time.Second,
				DisableCompression:    false,
			},
		}
	}
	if d.Filename == "" && d.URL != "" {
		d.Filename = filepath.Base(d.URL)
	}
	if d.TargetPath == "" && d.Filename != "" {
		d.TargetPath = d.Filename
	}
}

// Pause signals the download to pause
func (d *Download) Pause() {
	d.mutex.Lock()
	oldStatus := d.Status
	defer d.mutex.Unlock()

	if d.Status == "downloading" && !d.isPaused && !d.isCancelled {
		d.Status = "paused"
		d.isPaused = true
		// Log status change to paused
		logger.LogDownloadStatus(d.URL, oldStatus, d.Status, d.Downloaded, d.TotalSize)
		select {
		case d.pauseChan <- struct{}{}:
		default:
		}
	}
}

// Resume signals the download to resume
func (d *Download) Resume() {
	d.mutex.Lock()
	oldStatus := d.Status
	defer d.mutex.Unlock()

	if d.Status == "paused" && d.isPaused && !d.isCancelled {
		d.Status = "downloading"
		d.isPaused = false
		// Log status change to downloading (resumed)
		logger.LogDownloadStatus(d.URL, oldStatus, d.Status, d.Downloaded, d.TotalSize)
		select {
		case d.resumeChan <- struct{}{}:
		default:
		}
	}
}

// Cancel stops the download and removes temporary files
func (d *Download) Cancel() error {
	d.mutex.Lock()
	oldStatus := d.Status
	defer d.mutex.Unlock()

	if d.Status != "completed" && d.Status != "cancelled" && !d.isCancelled {
		d.Status = "cancelled"
		d.isCancelled = true
		// Log status change to cancelled
		logger.LogDownloadStatus(d.URL, oldStatus, d.Status, d.Downloaded, d.TotalSize)
		select {
		case d.cancelChan <- struct{}{}:
		default:
		}

		// Only attempt to remove the file if it was created
		if d.TargetPath != "" && d.Progress > 0 {
			if err := os.Remove(d.TargetPath); err != nil && !os.IsNotExist(err) {
				errorMsg := fmt.Sprintf("failed to remove file: %v", err)
				logger.LogDownloadError(d.URL, d.Queue, errorMsg)
				return fmt.Errorf("failed to remove file: %v", err)
			}
		}
	}
	return nil
}

// Retry attempts to restart a failed download
func (d *Download) Retry() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status == "error" {
		oldStatus := d.Status // Save the old status for logging
		d.Status = "pending"
		d.Error = ""
		d.Progress = 0
		d.Speed = 0
		d.Downloaded = 0
		d.retryCount++
		// Log status change to pending (retry)
		logger.LogDownloadPending(d.URL, d.Queue, fmt.Sprintf("Retry attempt %d of %d", d.retryCount, d.maxRetries))
		// Log the status change
		logger.LogDownloadStatus(d.URL, oldStatus, d.Status, 0, d.TotalSize)
		return nil
	}
	return fmt.Errorf("download is not in error state")
}

// GetStatus returns the current status of the download
func (d *Download) GetStatus() string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.Status
}

// GetProgress returns the current progress percentage of the download
func (d *Download) GetProgress() float64 {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.Progress
}

// GetSpeed returns the current download speed in bytes per second
func (d *Download) GetSpeed() int64 {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.Speed
}

// GetRetryCount returns the current retry count for the download
func (d *Download) GetRetryCount() int {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.retryCount
}

// Start begins downloading a file with optional progress callback
func (d *Download) Start() error {
	// Initialize control channels and fields
	d.Initialize()
	d.mutex.Lock()
	oldStatus := d.Status
	d.Status = "downloading"
	d.StartTime = time.Now()
	d.mutex.Unlock()

	// Log download start
	logger.LogDownloadStart(d.URL, d.Queue, d.MaxBandwidth)
	// Log status change
	logger.LogDownloadStatus(d.URL, oldStatus, "downloading", 0, d.TotalSize)

	// Main download loop with retry logic
	for d.retryCount <= d.maxRetries {
		err := d.performDownload()
		if err == nil {
			// Download completed successfully
			d.mutex.Lock()
			oldStatus := d.Status
			d.Status = "completed"
			d.Progress = 100.0
			d.mutex.Unlock()

			// Calculate download duration
			duration := time.Since(d.StartTime)
			// Log download completion
			logger.LogDownloadComplete(d.URL, d.TargetPath, duration, d.TotalSize)
			// Log status change
			logger.LogDownloadStatus(d.URL, oldStatus, "completed", d.TotalSize, d.TotalSize)
			return nil
		}

		// Check if download was cancelled
		d.mutex.Lock()
		if d.isCancelled {
			d.mutex.Unlock()
			logger.LogDownloadStatus(d.URL, "downloading", "cancelled", d.Downloaded, d.TotalSize)
			return fmt.Errorf("download cancelled")
		}

		// Handle error and retry if possible
		oldStatus := d.Status
		d.Status = "error"
		d.Error = err.Error()

		// Log error status
		logger.LogDownloadError(d.URL, d.Queue, err.Error())
		logger.LogDownloadStatus(d.URL, oldStatus, "error", d.Downloaded, d.TotalSize)

		// Check if we should retry
		if d.retryCount < d.maxRetries {
			d.retryCount++
			d.Status = "pending"
			retryMsg := fmt.Sprintf("Retry attempt %d of %d after error: %s",
				d.retryCount, d.maxRetries, err.Error())
			logger.LogDownloadPending(d.URL, d.Queue, retryMsg)
			logger.LogDownloadStatus(d.URL, "error", "pending", d.Downloaded, d.TotalSize)
			d.mutex.Unlock()
			time.Sleep(d.retryDelay)
			continue
		}

		d.mutex.Unlock()
		finalError := fmt.Errorf("download failed after %d retries: %v", d.maxRetries, err)
		logger.LogDownloadError(d.URL, d.Queue, finalError.Error())
		return finalError
	}

	finalError := fmt.Errorf("download failed after %d retries", d.maxRetries)
	logger.LogDownloadError(d.URL, d.Queue, finalError.Error())
	return finalError
}

// performDownload handles the actual file download process
func (d *Download) performDownload() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(d.TargetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		errorMsg := fmt.Sprintf("failed to create directory: %v", err)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Send HEAD request to get file size
	resp, err := d.client.Head(d.URL)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to send HEAD request: %v", err)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("failed to send HEAD request: %w", err)
	}
	defer resp.Body.Close()

	totalSize, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	d.mutex.Lock()
	d.TotalSize = totalSize
	d.mutex.Unlock()

	// Log file size information
	logger.LogDownloadStatus(d.URL, "downloading", "downloading", 0, totalSize)

	// Check if server supports range requests
	supportsRanges := resp.Header.Get("Accept-Ranges") == "bytes"

	// Create the GET request
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to create request: %v", err)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// If we're resuming and the server supports ranges, set the range header
	d.mutex.Lock()
	startByte := d.Downloaded
	d.mutex.Unlock()

	if startByte > 0 && supportsRanges {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))
		logger.LogDownloadStatus(d.URL, "downloading", "downloading", startByte, totalSize)
	}

	// Send the request
	resp, err = d.client.Do(req)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to send GET request: %v", err)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorMsg := fmt.Sprintf("server responded with status: %s", resp.Status)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("server responded with status: %s", resp.Status)
	}

	// Prepare file for writing
	var file *os.File
	var openMode int

	if startByte > 0 && supportsRanges {
		// Append to existing file if resuming
		openMode = os.O_WRONLY | os.O_APPEND
	} else {
		// Create new file
		openMode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		startByte = 0
	}

	file, err = os.OpenFile(d.TargetPath, openMode, 0644)
	if err != nil {
		errorMsg := fmt.Sprintf("failed to open file: %v", err)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Setup rate limiting if needed
	var limiter *RateLimiter
	if d.MaxBandwidth > 0 {
		limiter = NewRateLimiter(d.MaxBandwidth * 1024) // Convert KB/s to bytes/s
		logger.LogDownloadPending(d.URL, d.Queue, fmt.Sprintf("Applying bandwidth limit of %d KB/s", d.MaxBandwidth))
	}

	// Track progress
	buffer := make([]byte, 32*1024)
	downloaded := startByte
	startTime := time.Now()
	lastUpdateTime := startTime
	lastBytes := downloaded

	// Start the download loop
	for {
		// Check if we should pause
		select {
		case <-d.pauseChan:
			logger.LogDownloadStatus(d.URL, "downloading", "paused", downloaded, totalSize)

			// Wait for resume signal
			<-d.resumeChan

			// Reset speed calculation
			startTime = time.Now()
			lastUpdateTime = startTime
			lastBytes = downloaded

			logger.LogDownloadStatus(d.URL, "paused", "downloading", downloaded, totalSize)
			continue

		case <-d.cancelChan:
			logger.LogDownloadStatus(d.URL, "downloading", "cancelled", downloaded, totalSize)
			return fmt.Errorf("download cancelled")

		default:
			// Proceed with download
		}

		// Apply rate limiting if needed
		if limiter != nil {
			n, err := limiter.Read(resp.Body, buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				errorMsg := fmt.Sprintf("error reading from response: %v", err)
				logger.LogDownloadError(d.URL, d.Queue, errorMsg)
				return fmt.Errorf("error reading from response: %w", err)
			}

			if n == 0 {
				break
			}

			// Write to file
			if _, err := file.Write(buffer[:n]); err != nil {
				errorMsg := fmt.Sprintf("error writing to file: %v", err)
				logger.LogDownloadError(d.URL, d.Queue, errorMsg)
				return fmt.Errorf("error writing to file: %w", err)
			}

			downloaded += int64(n)
		} else {
			// No rate limiting
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				errorMsg := fmt.Sprintf("error reading from response: %v", err)
				logger.LogDownloadError(d.URL, d.Queue, errorMsg)
				return fmt.Errorf("error reading from response: %w", err)
			}

			if n == 0 {
				break
			}

			// Write to file
			if _, err := file.Write(buffer[:n]); err != nil {
				errorMsg := fmt.Sprintf("error writing to file: %v", err)
				logger.LogDownloadError(d.URL, d.Queue, errorMsg)
				return fmt.Errorf("error writing to file: %w", err)
			}

			downloaded += int64(n)
		}

		// Update progress
		if totalSize > 0 {
			d.mutex.Lock()
			d.Progress = float64(downloaded) / float64(totalSize) * 100
			d.Downloaded = downloaded
			d.mutex.Unlock()
		}

		// Calculate speed and log progress (not too often)
		now := time.Now()
		elapsed := now.Sub(lastUpdateTime)
		if elapsed >= time.Second {
			bytesPerSecond := int64(float64(downloaded-lastBytes) / elapsed.Seconds())

			d.mutex.Lock()
			d.Speed = bytesPerSecond
			d.mutex.Unlock()

			// Log progress every 10% or at least every 30 seconds
			progressPercent := float64(downloaded) / float64(totalSize) * 100
			lastProgressPercent := float64(lastBytes) / float64(totalSize) * 100
			if (int(progressPercent/10) > int(lastProgressPercent/10)) || elapsed >= 30*time.Second {
				logger.LogDownloadStatus(d.URL, "downloading", "downloading", downloaded, totalSize)
			}

			lastUpdateTime = now
			lastBytes = downloaded
		}
	}

	// Verify download completed successfully
	if totalSize > 0 && downloaded < totalSize {
		errorMsg := fmt.Sprintf("download incomplete: got %d of %d bytes", downloaded, totalSize)
		logger.LogDownloadError(d.URL, d.Queue, errorMsg)
		return fmt.Errorf("download incomplete: got %d of %d bytes", downloaded, totalSize)
	}

	// Update final download size if we didn't know it before
	if totalSize <= 0 {
		d.mutex.Lock()
		d.TotalSize = downloaded
		d.Progress = 100.0
		d.mutex.Unlock()
	}

	logger.LogDownloadStatus(d.URL, "downloading", "completed", downloaded, downloaded)
	return nil
}

// New creates a new download instance
func New(url, targetPath, queue string, maxBandwidth int64) *Download {
	download := &Download{
		URL:          url,
		TargetPath:   targetPath,
		Filename:     filepath.Base(targetPath),
		Queue:        queue,
		Status:       "pending",
		MaxBandwidth: maxBandwidth,
		maxRetries:   3,
		retryDelay:   5 * time.Second,
	}
	download.Initialize()
	return download
}

// StartDownload is a convenience function to create and start a download
func StartDownload(url, targetPath, queue string, maxBandwidth int64) (*Download, error) {
	download := New(url, targetPath, queue, maxBandwidth)
	go download.Start()
	return download, nil
}
