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
)

// Download represents a download task with its state and control channels
type Download struct {
	URL         string  `json:"url"`
	TargetPath  string  `json:"target_path"`
	Filename    string  `json:"filename"`
	Queue       string  `json:"queue"`
	Status      string  `json:"status"` // pending, downloading, paused, completed, error, cancelled
	Progress    float64 `json:"progress"`
	Speed       int64   `json:"speed"` // bytes per second
	TotalSize   int64   `json:"total_size"`
	Downloaded  int64   `json:"downloaded"`
	Error       string  `json:"error,omitempty"`
	MaxBandwidth int64  `json:"max_bandwidth"` // in KB/s, 0 means unlimited
	StartTime   time.Time `json:"start_time,omitempty"`
	
	// Control fields (not persisted to JSON)
	pauseChan  chan struct{} `json:"-"`
	resumeChan chan struct{} `json:"-"`
	cancelChan chan struct{} `json:"-"`
	isPaused   bool          `json:"-"`
	isCancelled bool         `json:"-"`
	mutex      sync.Mutex    `json:"-"`
	retryCount int           `json:"retry_count"`
	maxRetries int           `json:"max_retries"`
	retryDelay time.Duration `json:"-"`
	client     *http.Client  `json:"-"`
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
	if d.client == nil {
		d.client = &http.Client{
			Timeout: 30 * time.Second,
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
	defer d.mutex.Unlock()

	if d.Status == "downloading" && !d.isPaused && !d.isCancelled {
		d.Status = "paused"
		d.isPaused = true
		select {
		case d.pauseChan <- struct{}{}:
		default:
		}
	}
}

// Resume signals the download to resume
func (d *Download) Resume() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status == "paused" && d.isPaused && !d.isCancelled {
		d.Status = "downloading"
		d.isPaused = false
		select {
		case d.resumeChan <- struct{}{}:
		default:
		}
	}
}

// Cancel stops the download and removes temporary files
func (d *Download) Cancel() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Status != "completed" && d.Status != "cancelled" && !d.isCancelled {
		d.Status = "cancelled"
		d.isCancelled = true
		select {
		case d.cancelChan <- struct{}{}:
		default:
		}

		// Only attempt to remove the file if it was created
		if d.TargetPath != "" && d.Progress > 0 {
			if err := os.Remove(d.TargetPath); err != nil && !os.IsNotExist(err) {
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
		d.Status = "pending"
		d.Error = ""
		d.Progress = 0
		d.Speed = 0
		d.Downloaded = 0
		d.retryCount++
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

// Start begins downloading a file with optional progress callback
func (d *Download) Start() error {
	// Initialize control channels and fields
	d.Initialize()
	d.mutex.Lock()
	d.Status = "downloading"
	d.StartTime = time.Now()
	d.mutex.Unlock()

	// Main download loop with retry logic
	for d.retryCount <= d.maxRetries {
		err := d.performDownload()
		if err == nil {
			// Download completed successfully
			d.mutex.Lock()
			d.Status = "completed"
			d.Progress = 100.0
			d.mutex.Unlock()
			return nil
		}

		// Check if download was cancelled
		d.mutex.Lock()
		if d.isCancelled {
			d.mutex.Unlock()
			return fmt.Errorf("download cancelled")
		}

		// Handle error and retry if possible
		d.Status = "error"
		d.Error = err.Error()
		
		// Check if we should retry
		if d.retryCount < d.maxRetries {
			d.retryCount++
			d.Status = "pending"
			d.mutex.Unlock()
			time.Sleep(d.retryDelay)
			continue
		}
		
		d.mutex.Unlock()
		return fmt.Errorf("download failed after %d retries: %v", d.maxRetries, err)
	}

	return fmt.Errorf("download failed after %d retries", d.maxRetries)
}

// performDownload handles the actual file download process
func (d *Download) performDownload() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(d.TargetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Send HEAD request to get file size
	resp, err := d.client.Head(d.URL)
	if err != nil {
		return fmt.Errorf("failed to send HEAD request: %w", err)
	}
	defer resp.Body.Close()

	totalSize, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	d.mutex.Lock()
	d.TotalSize = totalSize
	d.mutex.Unlock()
	
	// Check if server supports range requests
	supportsRanges := resp.Header.Get("Accept-Ranges") == "bytes"
	
	// Create the GET request
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// If we're resuming and the server supports ranges, set the range header
	d.mutex.Lock()
	startByte := d.Downloaded
	d.mutex.Unlock()
	
	if startByte > 0 && supportsRanges {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))
	} else {
		startByte = 0 // Reset to 0 if we can't resume
	}
	
	// Send the GET request
	resp, err = d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()
	
	// Create or open the output file
	var file *os.File
	if startByte > 0 {
		file, err = os.OpenFile(d.TargetPath, os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		file, err = os.Create(d.TargetPath)
	}
	
	if err != nil {
		return fmt.Errorf("failed to create/open file: %w", err)
	}
	defer file.Close()
	
	// Create a buffer for reading
	buffer := make([]byte, 32*1024) // 32KB buffer
	
	// Set up progress tracking
	downloaded := startByte
	lastUpdateTime := time.Now()
	lastDownloaded := downloaded
	
	// Main download loop
	for {
		// Check if download should be cancelled
		d.mutex.Lock()
		if d.isCancelled {
			d.mutex.Unlock()
			return fmt.Errorf("download cancelled")
		}
		
		// Check if download should be paused
		if d.isPaused {
			d.mutex.Unlock()
			select {
			case <-d.resumeChan:
				// Reset time measurement after resume
				lastUpdateTime = time.Now()
				lastDownloaded = downloaded
				continue
			case <-d.cancelChan:
				return fmt.Errorf("download cancelled")
			case <-time.After(100 * time.Millisecond):
				// Regularly check for cancel/resume signals
				continue
			}
		}
		d.mutex.Unlock()
		
		// Apply bandwidth limiting if needed
		d.mutex.Lock()
		maxBandwidth := d.MaxBandwidth
		d.mutex.Unlock()
		
		if maxBandwidth > 0 {
			// Calculate current speed and limit if necessary
			elapsed := time.Since(lastUpdateTime).Seconds()
			if elapsed > 0 {
				currentSpeed := float64(downloaded-lastDownloaded) / elapsed
				maxSpeed := float64(maxBandwidth * 1024) // Convert KB/s to B/s
				
				if currentSpeed > maxSpeed {
					sleepTime := time.Duration((float64(downloaded-lastDownloaded)/maxSpeed - elapsed) * float64(time.Second))
					if sleepTime > 0 {
						time.Sleep(sleepTime)
					}
				}
			}
		}
		
		// Read from response body
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			// Write to file
			if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to file: %w", writeErr)
			}
			
			// Update progress
			downloaded += int64(n)
			
			// Calculate speed and update progress
			elapsed := time.Since(lastUpdateTime).Seconds()
			if elapsed >= 0.5 { // Update stats every 0.5 seconds
				speed := int64(float64(downloaded-lastDownloaded) / elapsed)
				progress := 0.0
				if totalSize > 0 {
					progress = float64(downloaded) * 100.0 / float64(totalSize)
				}
				
				// Update download stats
				d.mutex.Lock()
				d.Speed = speed
				d.Progress = progress
				d.Downloaded = downloaded
				d.mutex.Unlock()
				
				// Reset for next measurement
				lastUpdateTime = time.Now()
				lastDownloaded = downloaded
			}
		}
		
		// Handle end of file or error
		if err != nil {
			if err == io.EOF {
				// Update final progress
				d.mutex.Lock()
				d.Downloaded = downloaded
				if totalSize > 0 {
					d.Progress = 100.0
				} else {
					d.Progress = 100.0 // Assume complete if size was unknown
				}
				d.Speed = 0
				d.mutex.Unlock()
				return nil
			}
			return fmt.Errorf("error reading response: %w", err)
		}
	}
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
