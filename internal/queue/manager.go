package queue

import (
	"sync"
	"time"

	"github.com/mahdiXak47/Download-Manager/internal/config"
	"github.com/mahdiXak47/Download-Manager/internal/downloader"
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
		ticker:     time.NewTicker(time.Minute),
	}

	// Initialize existing downloads
	for i := range cfg.Downloads {
		d := &cfg.Downloads[i]
		m.downloads[d.URL] = d
		if d.Status == "downloading" {
			m.activeJobs[d.Queue]++
		}
	}

	return m
}

// Start begins the queue manager's operation
func (m *Manager) Start() {
	go m.run()
}

// Stop stops the queue manager
func (m *Manager) Stop() {
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
		d.Pause()
		m.activeJobs[d.Queue]--

		// Save state
		if err := config.SaveConfig(m.config); err != nil {
			// Handle error
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
			return
		}

		if !queueCfg.IsTimeAllowed() {
			return
		}

		if m.activeJobs[d.Queue] >= queueCfg.MaxConcurrent {
			return
		}

		// Resume the download
		d.Resume()
		m.activeJobs[d.Queue]++

		// Save state
		if err := config.SaveConfig(m.config); err != nil {
			// Handle error
		}
	}
}

// processQueues checks each queue and starts eligible downloads
func (m *Manager) processQueues() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, queueCfg := range m.config.Queues {
		if !queueCfg.IsTimeAllowed() {
			continue
		}

		activeCount := m.activeJobs[queueCfg.Name]
		if activeCount >= queueCfg.MaxConcurrent {
			continue
		}

		// Find pending downloads for this queue
		for i := range m.config.Downloads {
			download := &m.config.Downloads[i]
			if download.Queue == queueCfg.Name && download.Status == "pending" {
				m.startDownload(download, &queueCfg)
				m.downloads[download.URL] = download

				activeCount++
				if activeCount >= queueCfg.MaxConcurrent {
					break
				}
			}
		}
	}
}

// startDownload begins a new download
func (m *Manager) startDownload(d *downloader.Download, q *config.QueueConfig) {
	d.Status = "downloading"
	m.activeJobs[q.Name]++

	go func() {
		// Start the actual download
		err := downloader.StartDownload(d, q.SpeedLimit)

		m.mutex.Lock()
		defer m.mutex.Unlock()

		// Update download status
		if err != nil {
			d.Status = "error"
			d.Error = err.Error()
		} else {
			d.Status = "completed"
		}

		// Decrease active job count
		m.activeJobs[q.Name]--

		// Save the updated state
		if err := config.SaveConfig(m.config); err != nil {
			// Handle error
		}
	}()
}
