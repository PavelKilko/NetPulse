package services

import (
	"github.com/PavelKilko/NetPulse/monitor/repository"
	"log"
	"net/http"
	"sync"
	"time"
)

// Mutex to manage concurrent access to monitoring tasks
var mu sync.Mutex

// Map to track active monitoring routines
var monitoringTasks = make(map[uint]chan bool)

// StartMonitoring initiates a new monitoring Go routine for the specified URL
func StartMonitoring(urlID uint, url string) {
	mu.Lock()
	defer mu.Unlock()

	// If a monitoring task already exists for this URL, don't start a new one
	if _, exists := monitoringTasks[urlID]; exists {
		log.Printf("Monitoring already running for URL ID %d", urlID)
		return
	}

	// Create a new control channel to stop the monitoring routine
	stopChan := make(chan bool)
	monitoringTasks[urlID] = stopChan

	go func() {
		for {
			select {
			case <-stopChan:
				log.Printf("Stopping monitoring for URL ID %d", urlID)
				return
			default:
				// Perform the monitoring task
				start := time.Now()
				resp, err := http.Get(url)
				responseTime := time.Since(start).Milliseconds()
				statusCode := 0

				if err != nil {
					log.Printf("Failed to reach URL %s: %s", url, err)
				} else {
					statusCode = resp.StatusCode
					resp.Body.Close()
				}

				// Store the result in MongoDB
				timestamp := time.Now()
				repository.StoreMonitoringResult(urlID, int(responseTime), statusCode, timestamp)

				// Sleep for the monitoring interval (e.g., 1 minute)
				time.Sleep(1 * time.Minute)
			}
		}
	}()
}

// StopMonitoring stops the monitoring Go routine for the specified URL
func StopMonitoring(urlID uint) {
	mu.Lock()
	defer mu.Unlock()

	// Check if a monitoring task exists for this URL ID
	if stopChan, exists := monitoringTasks[urlID]; exists {
		// Send a signal to stop the Go routine
		close(stopChan)
		delete(monitoringTasks, urlID)
		log.Printf("Monitoring stopped for URL ID %d", urlID)
	} else {
		log.Printf("No monitoring task found for URL ID %d", urlID)
	}
}
