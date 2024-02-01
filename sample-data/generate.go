package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type LogStats struct {
	RecordsGenerated int
	StartTime        time.Time
}

func generateSampleData() string {
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
	serviceNames := []string{"checkout", "payment", "shipping", "authentication", "inventory", "customer_support"}
	httpStatusCodes := []int{200, 201, 204, 400, 401, 403, 404, 500, 503}
	responseTime := rand.Intn(150) + 50
	userID := fmt.Sprintf("user%d", rand.Intn(9000)+1000)
	transactionID := fmt.Sprintf("tx%d", rand.Intn(9000)+1000)
	additionalInfo := "Sample log entry for testing."

	logEntry := fmt.Sprintf("%s %s %d %dms %s %s %s", timestamp,
		serviceNames[rand.Intn(len(serviceNames))], httpStatusCodes[rand.Intn(len(httpStatusCodes))],
		responseTime, userID, transactionID, additionalInfo)

	return logEntry
}

func writeToLogFile(logFileName string, wg *sync.WaitGroup, stopChan chan struct{}, statsChan chan LogStats) {
	defer wg.Done()

	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", logFileName, err)
		return
	}
	defer file.Close()

	stats := LogStats{StartTime: time.Now()}

	for {
		select {
		case <-stopChan:
			statsChan <- stats
			return
		default:
			logEntry := generateSampleData()
			// fmt.Println(logEntry)
			_, err := file.WriteString(logEntry + "\n")
			if err != nil {
				fmt.Printf("Error writing to file %s: %v\n", logFileName, err)
				return
			}
			stats.RecordsGenerated++
			time.Sleep(1 * time.Microsecond) // Adjust the sleep duration as needed
		}
	}
}

func main() {
	scriptDirectory, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	stopChan := make(chan struct{}, 2)
	statsChan := make(chan LogStats, 2)

	// Capture interrupt signals (e.g., Ctrl+C)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	for i := 1; i <= 2; i++ { // Create 2 goroutines
		// wg.Add(1)
		// go writeToLogFile(&wg, stopChan, statsChan)
		logFileName := fmt.Sprintf("%s/sample-%02d.log", scriptDirectory, i)
		wg.Add(1)
		go writeToLogFile(logFileName, &wg, stopChan, statsChan)
	}

	fmt.Println("Generating logs, press CTRL+C to interupt!")

	// Wait for termination signal
	<-signalChannel
	close(stopChan) // Signal goroutines to stop

	// Wait for goroutines to finish
	wg.Wait()

	// Close the stats channel before retrieving results
	close(statsChan)
	printStats(statsChan)

	fmt.Println("Graceful shutdown complete.")
}

func printStats(statsChan chan LogStats) {
	var generated int
	var startTime time.Time
	for stat := range statsChan {
		generated += stat.RecordsGenerated
		startTime = stat.StartTime
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("\nTotal records generated across all files: %d\n", generated)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}
