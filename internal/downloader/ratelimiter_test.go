package downloader

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	t.Run("Test Token Generation", func(t *testing.T) {
		// Create a rate limiter with 1000 bytes per second
		limiter := NewRateLimiter(1000)
		defer limiter.Stop()
		
		// Wait for tokens to be generated
		time.Sleep(100 * time.Millisecond)
		
		// Take some tokens and measure the time
		start := time.Now()
		limiter.GetToken(500) // Take half the tokens
		
		// This should be very quick as tokens are already available
		firstDuration := time.Since(start)
		if firstDuration > 100*time.Millisecond {
			t.Errorf("Initial token acquisition too slow: %v", firstDuration)
		}
		
		// Now take more tokens than available and measure time
		start = time.Now()
		limiter.GetToken(600) // This should take additional time
		
		// This should take some time to generate additional tokens
		secondDuration := time.Since(start)
		if secondDuration < 50*time.Millisecond {
			t.Errorf("Rate limiting not working, second acquisition too fast: %v", secondDuration)
		}
	})
	
	t.Run("Test Rate Limited Reader Small", func(t *testing.T) {
		// Create a smaller test to be less sensitive to timing issues
		// Create a rate limiter with 1KB per second
		limiter := NewRateLimiter(1024)
		defer limiter.Stop()
		
		// Create a test reader with 2KB of data
		data := make([]byte, 2*1024)
		reader := bytes.NewReader(data)
		
		// Read data with rate limiting
		buffer := make([]byte, 512) // 0.5KB buffer
		start := time.Now()
		
		bytesRead := 0
		for bytesRead < len(data) {
			n, err := limiter.Read(reader, buffer)
			if err != nil && err != io.EOF {
				t.Fatalf("Error reading data: %v", err)
			}
			
			bytesRead += n
			
			if err == io.EOF {
				break
			}
		}
		
		duration := time.Since(start)
		
		// Should take approximately 2 seconds to read 2KB at 1KB/s
		// Allow for some flexibility in timing
		// This just verifies it's not instant or extremely slow
		if duration < 500*time.Millisecond {
			t.Errorf("Rate limiter read too fast: %v, expected slow read for 2KB at 1KB/s", duration)
		}
		
		if duration > 5*time.Second {
			t.Errorf("Rate limiter read too slow: %v, expected approximately 2s for 2KB at 1KB/s", duration)
		}
		
		if bytesRead != len(data) {
			t.Errorf("Did not read expected amount of data: got %d bytes, want %d bytes", bytesRead, len(data))
		}
	})
} 