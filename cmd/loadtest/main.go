package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "Base URL of the API")
	totalRequests := flag.Int("n", 1000, "Total number of requests per endpoint")
	concurrency := flag.Int("c", 100, "Number of concurrent workers")
	flag.Parse()

	endpoints := []struct {
		name     string
		path     string
		isDetail bool
		maxID    int
	}{
		{"Get All Products", "/products", false, 0},
		{"Get Product By ID", "/products/%d", true, 40}, // Assumption: 40 products from seeders
		{"Get All Categories", "/categories", false, 0},
		{"Get Category By ID", "/categories/%d", true, 8}, // Assumption: 8 categories from seeders
	}

	fmt.Printf("Starting load test on %s\n", *baseURL)
	fmt.Printf("Requests per endpoint: %d\n", *totalRequests)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	for _, ep := range endpoints {
		runTest(*baseURL, ep.name, ep.path, ep.isDetail, ep.maxID, *totalRequests, *concurrency)
	}
}

func runTest(baseURL, name, pathTmpl string, isDetail bool, maxID, totalReqs, concurrency int) {
	fmt.Printf("Testing %s...\n", name)

	var wg sync.WaitGroup
	// Handle case where totalReqs < concurrency
	if totalReqs < concurrency {
		concurrency = totalReqs
	}
	requestsPerWorker := totalReqs / concurrency
	remainder := totalReqs % concurrency

	start := time.Now()
	var successCount, failCount int64
	var mu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		// Distribute remainder requests to the first few workers
		tasks := requestsPerWorker
		if i < remainder {
			tasks++
		}

		go func(count int) {
			defer wg.Done()
			localSuccess := 0
			localFail := 0

			for j := 0; j < count; j++ {
				url := baseURL + pathTmpl
				if isDetail {
					// Use random ID between 1 and maxID
					id := rand.Intn(maxID) + 1
					url = fmt.Sprintf(baseURL+pathTmpl, id)
				}

				resp, err := http.Get(url)
				if err != nil {
					localFail++
					continue
				}
				// Typically 200 OK is expected
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					localSuccess++
				} else {
					localFail++
					// Optional: Print first few errors to debug
					// fmt.Printf("Error: %d %s\n", resp.StatusCode, resp.Status)
				}
				resp.Body.Close()
			}

			mu.Lock()
			successCount += int64(localSuccess)
			failCount += int64(localFail)
			mu.Unlock()
		}(tasks)
	}

	wg.Wait()
	duration := time.Since(start)

	total := successCount + failCount
	var rps float64
	if duration.Seconds() > 0 {
		rps = float64(total) / duration.Seconds()
	}

	fmt.Printf("  Duration: %v\n", duration)
	fmt.Printf("  Total Requests: %d\n", total)
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Failed: %d\n", failCount)
	fmt.Printf("  RPS: %.2f\n", rps)
	fmt.Println("------------------------------------------------")
}
