package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Create a temporary directory for the log file
	tempDir, err := os.MkdirTemp("", "logger-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up the log file
	logPath := filepath.Join(tempDir, "test.log")
	
	// Initialize the logger
	err = Initialize(logPath)
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Close()

	// Test various logging functions
	tests := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{
			name:     "LogDownloadStart",
			logFunc:  func() { LogDownloadStart("http://example.com/file.zip", "default", 100) },
			contains: "Download started",
		},
		{
			name:     "LogDownloadStatus",
			logFunc:  func() { LogDownloadStatus("http://example.com/file.zip", "pending", "downloading", 0, 1000) },
			contains: "Status changed",
		},
		{
			name:     "LogDownloadError",
			logFunc:  func() { LogDownloadError("http://example.com/file.zip", "default", "network error") },
			contains: "Error for download",
		},
		{
			name:     "LogDownloadPending",
			logFunc:  func() { LogDownloadPending("http://example.com/file.zip", "default", "waiting for queue") },
			contains: "Download pending",
		},
		{
			name:     "LogDownloadComplete",
			logFunc:  func() { LogDownloadComplete("http://example.com/file.zip", "/downloads/file.zip", 5*time.Second, 1024) },
			contains: "Download complete",
		},
		{
			name:     "LogDownloadEvent",
			logFunc:  func() { LogDownloadEvent("CUSTOM", "custom event message") },
			contains: "custom event message",
		},
	}

	// Execute each test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc()
			
			// Read the log file to verify the message was written
			content, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatalf("Failed to read log file: %v", err)
			}
			
			if !strings.Contains(string(content), tt.contains) {
				t.Errorf("Log file does not contain expected text. Got: %s, Want substring: %s", string(content), tt.contains)
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Test getting the global logger instance
	logger := GetLogger()
	if logger == nil {
		t.Error("GetLogger() returned nil")
	}
} 