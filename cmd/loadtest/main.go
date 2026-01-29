package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	baseURL := flag.String("url", "https://kasir-api-xfuadi7395-bmt39rrk.leapcell.dev", "Base URL of the API")
	totalRequests := flag.Int("n", 1000, "Total number of requests per endpoint")
	concurrency := flag.Int("c", 100, "Number of concurrent workers")
	debug := flag.Bool("debug", false, "Enable verbose error logging")
	target := flag.String("target", "", "Filter endpoints by name (case-insensitive substring)")
	flag.Parse()

	// 1. Fetch valid IDs dynamically
	productIDs, err := fetchIDs(*baseURL + "/products")
	if err != nil {
		fmt.Printf("Warning: Failed to fetch products: %v. Using default range 1-40.\n", err)
		// Fallback
		for i := 1; i <= 40; i++ {
			productIDs = append(productIDs, i)
		}
	} else {
		fmt.Printf("Fetched %d valid product IDs\n", len(productIDs))
	}

	categoryIDs, err := fetchIDs(*baseURL + "/categories")
	if err != nil {
		fmt.Printf("Warning: Failed to fetch categories: %v. Using default range 1-8.\n", err)
		// Fallback
		for i := 1; i <= 8; i++ {
			categoryIDs = append(categoryIDs, i)
		}
	} else {
		fmt.Printf("Fetched %d valid category IDs\n", len(categoryIDs))
	}

	endpoints := []struct {
		name     string
		path     string
		isDetail bool
		ids      []int
	}{
		{"Get All Products", "/products", false, nil},
		{"Get Product By ID", "/products/%d", true, productIDs},
		{"Get All Categories", "/categories", false, nil},
		{"Get Category By ID", "/categories/%d", true, categoryIDs},
	}

	fmt.Printf("Starting load test on %s\n", *baseURL)
	fmt.Printf("Requests per endpoint: %d\n", *totalRequests)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	for _, ep := range endpoints {
		if *target != "" && !strings.Contains(strings.ToLower(ep.name), strings.ToLower(*target)) {
			continue
		}
		runTest(*baseURL, ep.name, ep.path, ep.isDetail, ep.ids, *totalRequests, *concurrency, *debug)
	}
}

func fetchIDs(url string) ([]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Assuming response structure: { "data": [ { "id": 1, ... }, ... ] }
	// or array [ { "id": 1, ... } ] depending on API
	// Let's assume generic structure wrapper or array of structs with ID
	var result struct {
		Data []struct {
			ID int `json:"id"`
		} `json:"data"` // Check if wrapped in data
	}

	// Try unmarshal as wrapper
	if err := json.Unmarshal(body, &result); err == nil && len(result.Data) > 0 {
		ids := make([]int, len(result.Data))
		for i, item := range result.Data {
			ids[i] = item.ID
		}
		return ids, nil
	}

	// Try unmarshal as array
	var arrayResult []struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(body, &arrayResult); err == nil {
		ids := make([]int, len(arrayResult))
		for i, item := range arrayResult {
			ids[i] = item.ID
		}
		return ids, nil
	}

	return nil, fmt.Errorf("could not parse IDs from response")
}

func runTest(baseURL, name, pathTmpl string, isDetail bool, ids []int, totalReqs, concurrency int, debug bool) {
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
					// Use random ID from valid list
					if len(ids) > 0 {
						id := ids[rand.Intn(len(ids))]
						url = fmt.Sprintf(baseURL+pathTmpl, id)
					} else {
						// Fallback if no IDs (shouldn't happen with default)
						url = fmt.Sprintf(baseURL+pathTmpl, 1)
					}
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
					if debug {
						fmt.Printf("[%s] Request Failed: %s | Status: %d | Error: %v\n", name, url, resp.StatusCode, resp.Status)
					}
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
