package main

type queue struct {
	URL               string
	FileName          string
	SavedAddress      string
	MaxNumOfDownloads int
	MaxBandwidth      int
	ActiveDuration    int
	MaxNumOfRetry     int
	Status            string //pending, downloading, completed, failed
	Progress          int
}

var downloadQueue = make(chan queue, 100)

func addQueue(
	url string, fileName string, savedAddress string, maxNumOfDownloads int,
	maxBandwidth int, activeDuration int, maxNumOfRetry int) {

	task := queue{
		URL:               url,
		FileName:          fileName,
		SavedAddress:      savedAddress,
		MaxNumOfDownloads: maxNumOfDownloads,
		MaxBandwidth:      maxBandwidth,
		ActiveDuration:    activeDuration,
		MaxNumOfRetry:     maxNumOfRetry,
		Status:            "pending",
		Progress:          0,
	}
	downloadQueue <- task
}

func processQueue() {
	for task := range downloadQueue {
		task.Status = "downloading"

	}
}
