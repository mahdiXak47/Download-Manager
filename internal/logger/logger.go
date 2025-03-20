package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile *os.File
	mu      sync.Mutex
	logger  *Logger
)

// Logger represents a simple logging utility
type Logger struct {
	filePath string
	enabled  bool
}

// Initialize sets up the logger with the specified file path
func Initialize(filePath string) error {
	mu.Lock()
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		mu.Unlock()
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		mu.Unlock()
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Close existing log file if open
	if logFile != nil {
		logFile.Close()
	}

	logFile = file
	logger = &Logger{
		filePath: filePath,
		enabled:  true,
	}
	
	// Unlock the mutex before logging the initialization event
	mu.Unlock()

	// Log initialization without mutex contention
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, "SYSTEM", "Logger initialized")
	
	// Write directly to the file
	if _, err := logFile.WriteString(logLine); err != nil {
		return fmt.Errorf("failed to write initialization log: %w", err)
	}

	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if logger == nil {
		// Create default logger to download-logs.log in current directory
		err := Initialize("download-logs.log")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		}
	}
	return logger
}

// logDownloadEvent logs an event related to downloads
func logDownloadEvent(eventType, message string) error {
	if logFile == nil {
		return fmt.Errorf("logger not initialized")
	}

	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, eventType, message)

	_, err := logFile.WriteString(logLine)
	if err != nil {
		return fmt.Errorf("failed to write to log: %w", err)
	}

	return nil
}

// LogDownloadStart logs when a download starts
func LogDownloadStart(url, queue string, maxBandwidth int64) {
	message := fmt.Sprintf("Download started - URL: %s, Queue: %s, Bandwidth Limit: %d KB/s", 
		url, queue, maxBandwidth)
	logDownloadEvent("START", message)
}

// LogDownloadStatus logs status changes for downloads
func LogDownloadStatus(url, oldStatus, newStatus string, downloadedBytes, totalBytes int64) {
	var message string
	if totalBytes > 0 {
		progress := float64(downloadedBytes) / float64(totalBytes) * 100
		message = fmt.Sprintf("Status changed for %s: %s -> %s (Progress: %.2f%%, Downloaded: %d/%d bytes)", 
			url, oldStatus, newStatus, progress, downloadedBytes, totalBytes)
	} else {
		message = fmt.Sprintf("Status changed for %s: %s -> %s", url, oldStatus, newStatus)
	}
	logDownloadEvent("STATUS", message)
}

// LogDownloadError logs download errors
func LogDownloadError(url, queue, errorMsg string) {
	message := fmt.Sprintf("Error for download %s in queue %s: %s", url, queue, errorMsg)
	logDownloadEvent("ERROR", message)
}

// LogDownloadPending logs when a download becomes pending
func LogDownloadPending(url, queue, reason string) {
	message := fmt.Sprintf("Download pending - URL: %s, Queue: %s, Reason: %s", url, queue, reason)
	logDownloadEvent("PENDING", message)
}

// LogDownloadComplete logs when a download completes
func LogDownloadComplete(url, targetPath string, duration time.Duration, size int64) {
	speedMBps := float64(size) / (1024 * 1024 * duration.Seconds())
	message := fmt.Sprintf("Download complete - URL: %s, Path: %s, Duration: %s, Size: %d bytes, Avg Speed: %.2f MB/s", 
		url, targetPath, duration.String(), size, speedMBps)
	logDownloadEvent("COMPLETE", message)
}

// LogDownloadEvent logs a general download-related event
func LogDownloadEvent(eventType, message string) error {
	return logDownloadEvent(eventType, message)
}

// Close closes the log file
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		return logFile.Close()
	}
	return nil
} 