package queue

import (
	"fmt"
	"sync"
	"time"

	"github.com/mahdiXak47/Download-Manager/internal/config"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
	"github.com/mahdiXak47/Download-Manager/internal/logger"
)

type Manager struct {
	config     *config.Config
	activeJobs map[string]int                  // queue name -> active download count
	downloads  map[string]*downloader.Download // URL -> Download for quick lookup
	mutex      sync.Mutex
	ticker     *time.Ticker
}

func NewManager(cfg *config.Config) *Manager {
	m := &Manager{
		config:     cfg,
		activeJobs: make(map[string]int),
		downloads:  make(map[string]*downloader.Download),
		ticker:     time.NewTicker(10 * time.Second),
	}

	// Initialize existing downloads
	for i := range cfg.Downloads {
		d := &cfg.Downloads[i]
		m.downloads[d.URL] = d
		if d.Status == "downloading" {
			m.activeJobs[d.Queue]++
		}
	}

	logger.LogDownloadEvent("SYSTEM", fmt.Sprintf("Queue Manager initialized with %d downloads", len(cfg.Downloads)))
	return m
}

// Start begins the queue manager's operation
func (m *Manager) Start() {
	logger.LogDownloadEvent("SYSTEM", "Queue Manager started")
	go m.run()
}

// Stop stops the queue manager
func (m *Manager) Stop() {
	logger.LogDownloadEvent("SYSTEM", "Queue Manager stopped")
	m.ticker.Stop()
}

// run is the main loop that processes downloads
func (m *Manager) run() {
	for range m.ticker.C {
		m.processQueues()
	}
}

// PauseDownload pauses a specific download
func (m *Manager) PauseDownload(url string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if d, exists := m.downloads[url]; exists && d.Status == "downloading" {
		logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Pausing download %s in queue %s", url, d.Queue))
		d.Pause()
		m.activeJobs[d.Queue]--

		// Save state
		if err := config.SaveConfig(m.config); err != nil {
			logger.LogDownloadError(url, d.Queue, fmt.Sprintf("Failed to save config when pausing: %v", err))
		}
	}
}

// ResumeDownload resumes a specific download
func (m *Manager) ResumeDownload(url string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if d, exists := m.downloads[url]; exists && d.Status == "paused" {
		// Check if we can resume based on queue limits
		queueCfg := m.config.GetQueue(d.Queue)
		if queueCfg == nil {
			logger.LogDownloadError(url, d.Queue, "Cannot resume: queue configuration not found")
			return
		}

		if !queueCfg.IsTimeAllowed() {
			logger.LogDownloadPending(url, d.Queue, fmt.Sprintf("Cannot resume: outside allowed time window (%s-%s)",
				queueCfg.StartTime, queueCfg.EndTime))
			return
		}

		if m.activeJobs[d.Queue] >= queueCfg.MaxConcurrent {
			logger.LogDownloadPending(url, d.Queue, fmt.Sprintf("Cannot resume: queue at maximum capacity (%d downloads)",
				queueCfg.MaxConcurrent))
			return
		}

		// Resume the download
		logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Resuming download %s in queue %s", url, d.Queue))
		d.Resume()
		m.activeJobs[d.Queue]++

		// Save state
		if err := config.SaveConfig(m.config); err != nil {
			logger.LogDownloadError(url, d.Queue, fmt.Sprintf("Failed to save config when resuming: %v", err))
		}
	}
}

// processQueues checks each queue and starts eligible downloads
func (m *Manager) processQueues() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	logger.LogDownloadEvent("SYSTEM", "Processing queues")

	for _, queueCfg := range m.config.Queues {
		if !queueCfg.IsTimeAllowed() {
			logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Queue %s: Outside allowed time window (%s-%s)",
				queueCfg.Name, queueCfg.StartTime, queueCfg.EndTime))
			continue
		}

		activeCount := m.activeJobs[queueCfg.Name]
		if activeCount >= queueCfg.MaxConcurrent {
			logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Queue %s: At maximum capacity (%d/%d downloads)",
				queueCfg.Name, activeCount, queueCfg.MaxConcurrent))
			continue
		}

		// Find pending downloads for this queue
		pendingCount := 0
		startedCount := 0
		for i := range m.config.Downloads {
			download := &m.config.Downloads[i]
			if download.Queue == queueCfg.Name && download.Status == "pending" {
				pendingCount++
				m.startDownload(download, &queueCfg)
				m.downloads[download.URL] = download

				activeCount++
				startedCount++
				if activeCount >= queueCfg.MaxConcurrent {
					break
				}
			}
		}

		logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Queue %s: Started %d of %d pending downloads (%d/%d active)",
			queueCfg.Name, startedCount, pendingCount, activeCount, queueCfg.MaxConcurrent))
	}
}

// startDownload begins a new download
func (m *Manager) startDownload(d *downloader.Download, q *config.QueueConfig) {
	d.Status = "downloading"
	m.activeJobs[q.Name]++

	logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Starting download %s in queue %s", d.URL, q.Name))

	go func() {
		// Start the actual download
		err := d.Start()

		m.mutex.Lock()
		defer m.mutex.Unlock()

		// Update download status
		if err != nil && d.Status != "cancelled" {
			d.Status = "error"
			d.Error = err.Error()
			logger.LogDownloadError(d.URL, q.Name, fmt.Sprintf("Download failed: %v", err))
		} else if d.Status != "cancelled" {
			d.Status = "completed"
			logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Download %s completed in queue %s", d.URL, q.Name))
		}

		// Decrease active job count
		m.activeJobs[q.Name]--
		logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Queue %s: Active downloads decreased to %d/%d",
			q.Name, m.activeJobs[q.Name], q.MaxConcurrent))

		// Save the updated state
		if err := config.SaveConfig(m.config); err != nil {
			logger.LogDownloadError(d.URL, q.Name, fmt.Sprintf("Failed to save config after download: %v", err))
		}
	}()
}

// RemoveDownload removes a download from the queue
func (m *Manager) RemoveDownload(url string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Find the download first to log its details
	var queueName string
	for _, d := range m.config.Downloads {
		if d.URL == url {
			queueName = d.Queue
			break
		}
	}

	// Remove from active downloads
	delete(m.downloads, url)

	// Remove from config downloads
	for i, d := range m.config.Downloads {
		if d.URL == url {
			m.config.Downloads = append(m.config.Downloads[:i], m.config.Downloads[i+1:]...)
			logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Removed download %s from queue %s", url, queueName))
			break
		}
	}

	// Update active jobs count if needed
	if d, exists := m.downloads[url]; exists && d.Status == "downloading" {
		m.activeJobs[d.Queue]--
		logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Queue %s: Active downloads decreased to %d",
			d.Queue, m.activeJobs[d.Queue]))
	}
}

// ProcessDownload processes a specific download (used for retrying downloads)
func (m *Manager) ProcessDownload(url string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if d, exists := m.downloads[url]; exists && d.Status == "pending" {
		// Find the queue configuration
		var queueCfg *config.QueueConfig
		for _, q := range m.config.Queues {
			if q.Name == d.Queue {
				queueCfg = &q
				break
			}
		}

		if queueCfg == nil {
			logger.LogDownloadError(url, d.Queue, "Cannot process: queue configuration not found")
			return
		}

		if !queueCfg.IsTimeAllowed() {
			logger.LogDownloadPending(url, d.Queue, fmt.Sprintf("Cannot process: outside allowed time window (%s-%s)",
				queueCfg.StartTime, queueCfg.EndTime))
			return
		}

		// Check if we can start the download based on queue limits
		if m.activeJobs[d.Queue] >= queueCfg.MaxConcurrent {
			// Queue is at capacity, leave as pending
			logger.LogDownloadPending(url, d.Queue, fmt.Sprintf("Cannot process: queue at maximum capacity (%d downloads)",
				queueCfg.MaxConcurrent))
			return
		}

		// Process the download
		m.startDownload(d, queueCfg)

		logger.LogDownloadEvent("QUEUE", fmt.Sprintf("Processing download %s in queue %s", url, d.Queue))
	}
}

// ProcessAllQueues immediately processes all queues (used when a new download is added)
func (m *Manager) ProcessAllQueues() {
	go func() {
		// Small delay to allow the UI to update
		//time.Sleep(100 * time.Millisecond)
		m.processQueues()
	}()
}
