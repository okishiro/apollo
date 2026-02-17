package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// --- CONFIGURATION ---
	totalRequests := 10000
	concurrency := 10     // How many "workers" run at once
	targetSeconds := 10.0 // Duration to spread the load

	// Calculate delay to hit exactly 100 requests per 100 seconds (1 req/sec)
	// Rate = totalRequests / targetSeconds
	delay := time.Duration(float64(time.Second) / (float64(totalRequests) / targetSeconds))
	// ---------------------

	var wg sync.WaitGroup
	jobs := make(chan int, totalRequests)

	// 1. Start Workers
	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{Timeout: 5 * time.Second}

			for i := range jobs {
				sendRequest(client, i)
				time.Sleep(delay) // Throttle the speed
			}
		}(w)
	}

	// 2. Feed Jobs (Movie names 1 to 10000)
	start := time.Now()
	for j := 1; j <= totalRequests; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("\nDone! Took %v for %d requests\n", time.Since(start), totalRequests)
}

func sendRequest(client *http.Client, id int) {
	url := "http://localhost:8123/movie/bleh"
	jsonBody := []byte(fmt.Sprintf(`{
		"name": "%d",
		"date": "2024-02-01T00:00:00Z",
		"comment": "Stress test movie"
	}`, id))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request %d: Error -> %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request %d: Status %d\n", id, resp.StatusCode)
	} else if id%100 == 0 {
		fmt.Printf("Reached request %d...\n", id)
	}
}
