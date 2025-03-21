package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Simple logger that writes to a specified file
type SimpleLogger struct {
	file *os.File
}

// Initialize creates a new logger
func NewLogger(logPath string) (*SimpleLogger, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &SimpleLogger{file: file}, nil
}

// Log writes a message to the log file
func (l *SimpleLogger) Log(event, message string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, event, message)
	
	_, err := l.file.WriteString(logLine)
	if err != nil {
		return fmt.Errorf("failed to write to log: %v", err)
	}
	
	// Also print to console
	fmt.Print(logLine)
	
	return nil
}

// Close closes the log file
func (l *SimpleLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// DownloadWithLogging downloads a file with logging
func DownloadWithLogging(url, targetPath, logPath string) error {
	// Create logger
	logger, err := NewLogger(logPath)
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}
	defer logger.Close()
	
	logger.Log("START", fmt.Sprintf("Starting download of %s to %s", url, targetPath))
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Log("ERROR", fmt.Sprintf("Failed to create directory: %v", err))
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Get file information
	logger.Log("INFO", "Sending HEAD request")
	resp, err := http.Head(url)
	if err != nil {
		logger.Log("ERROR", fmt.Sprintf("Failed to send HEAD request: %v", err))
		return fmt.Errorf("failed to send HEAD request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check if the server supports ranges
	supportsRanges := resp.Header.Get("Accept-Ranges") == "bytes"
	logger.Log("INFO", fmt.Sprintf("Server supports ranges: %v", supportsRanges))
	
	// Get file size
	totalSize := int64(0)
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		totalSize, _ = strconv.ParseInt(contentLength, 10, 64)
		logger.Log("INFO", fmt.Sprintf("File size: %d bytes", totalSize))
	}
	
	// Create GET request
	logger.Log("INFO", "Sending GET request")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Log("ERROR", fmt.Sprintf("Failed to create request: %v", err))
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	// Send request
	client := &http.Client{
		Timeout: 0, // No timeout for downloads
	}
	response, err := client.Do(request)
	if err != nil {
		logger.Log("ERROR", fmt.Sprintf("Failed to send GET request: %v", err))
		return fmt.Errorf("failed to send GET request: %v", err)
	}
	defer response.Body.Close()
	
	// Check response status
	logger.Log("INFO", fmt.Sprintf("Response status: %s", response.Status))
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusPartialContent {
		logger.Log("ERROR", fmt.Sprintf("Failed to download file: status code %d", response.StatusCode))
		return fmt.Errorf("failed to download file: status code %d", response.StatusCode)
	}
	
	// Create the output file
	logger.Log("INFO", fmt.Sprintf("Creating output file: %s", targetPath))
	file, err := os.Create(targetPath)
	if err != nil {
		logger.Log("ERROR", fmt.Sprintf("Failed to create file: %v", err))
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Prepare for download
	buffer := make([]byte, 32*1024)
	downloaded := int64(0)
	startTime := time.Now()
	lastLogTime := startTime
	
	// Start download
	logger.Log("INFO", "Starting file download")
	for {
		// Read a chunk
		n, err := response.Body.Read(buffer)
		if err != nil && err != io.EOF {
			logger.Log("ERROR", fmt.Sprintf("Error reading from response: %v", err))
			return fmt.Errorf("error reading from response: %v", err)
		}
		
		if n == 0 {
			break
		}
		
		// Write chunk to file
		if _, err := file.Write(buffer[:n]); err != nil {
			logger.Log("ERROR", fmt.Sprintf("Error writing to file: %v", err))
			return fmt.Errorf("error writing to file: %v", err)
		}
		
		// Update progress
		downloaded += int64(n)
		now := time.Now()
		
		// Log progress every second
		if now.Sub(lastLogTime) >= time.Second {
			if totalSize > 0 {
				progress := float64(downloaded) / float64(totalSize) * 100
				speed := float64(downloaded) / now.Sub(startTime).Seconds() / 1024
				logger.Log("PROGRESS", fmt.Sprintf("Progress: %.2f%% (%d/%d bytes), Speed: %.2f KB/s", 
					progress, downloaded, totalSize, speed))
			} else {
				speed := float64(downloaded) / now.Sub(startTime).Seconds() / 1024
				logger.Log("PROGRESS", fmt.Sprintf("Downloaded: %d bytes, Speed: %.2f KB/s", 
					downloaded, speed))
			}
			lastLogTime = now
		}
		
		if err == io.EOF {
			break
		}
	}
	
	// Calculate download stats
	duration := time.Since(startTime)
	speed := float64(downloaded) / duration.Seconds() / 1024
	
	logger.Log("COMPLETE", fmt.Sprintf("Download completed in %s, Total: %d bytes, Avg Speed: %.2f KB/s", 
		duration, downloaded, speed))
	
	return nil
}

// TestDownload provides a simple test function to download a file
func TestDownload() {
	url := "https://dl.musichi.ir/1403/11/26/Amin%20Tijay%20-%20Chikar%20Kardi%20Baam.mp3"
	downloadDir := "./downloads"
	filename := filepath.Base(url)
	targetPath := filepath.Join(downloadDir, filename)
	logPath := "./download_test.log"
	
	fmt.Printf("Downloading %s to %s\n", url, targetPath)
	fmt.Printf("Logs will be written to %s\n", logPath)
	
	err := DownloadWithLogging(url, targetPath, logPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Download completed successfully!")
	}
}

// This file can be used to test downloading functionality
// Run with: go run cmd/test_download/main.go
func main() {
	TestDownload()
} 