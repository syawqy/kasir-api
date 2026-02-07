package main

import (
	"bytes"
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

type Product struct {
	ID    int `json:"id"`
	Stock int `json:"stock"`
}

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "Base URL of the API")
	totalRequests := flag.Int("n", 1000, "Total number of requests per endpoint")
	concurrency := flag.Int("c", 100, "Number of concurrent workers")
	debug := flag.Bool("debug", false, "Enable verbose error logging")
	target := flag.String("target", "", "Filter endpoints by name (case-insensitive substring)")
	flag.Parse()

	// 1. Fetch valid IDs dynamically
	products, err := fetchProducts(*baseURL + "/products")
	if err != nil {
		fmt.Printf("Warning: Failed to fetch products: %v. Using default range 1-40.\n", err)
		// Fallback
		for i := 1; i <= 40; i++ {
			products = append(products, Product{ID: i, Stock: 100})
		}
	} else {
		fmt.Printf("Fetched %d valid products\n", len(products))
	}

	productIDs := make([]int, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
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
		method   string
		isDetail bool
		isPost   bool
		ids      []int
	}{
		{"Get All Products", "/products", "GET", false, false, nil},
		{"Get Product By ID", "/products/%d", "GET", true, false, productIDs},
		{"Get All Categories", "/categories", "GET", false, false, nil},
		{"Get Category By ID", "/categories/%d", "GET", true, false, categoryIDs},
		{"Search Products", "/products?name=%s", "GET", false, false, nil},
		{"Create Transaction", "/checkout", "POST", false, true, productIDs},
	}

	fmt.Printf("Starting load test on %s\n", *baseURL)
	fmt.Printf("Requests per endpoint: %d\n", *totalRequests)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	for _, ep := range endpoints {
		if *target != "" && !strings.Contains(strings.ToLower(ep.name), strings.ToLower(*target)) {
			continue
		}
		if ep.isPost {
			runPostTest(*baseURL, ep.name, ep.path, ep.ids, *totalRequests, *concurrency, *debug)
		} else {
			runTest(*baseURL, ep.name, ep.path, ep.method, ep.isDetail, ep.ids, *totalRequests, *concurrency, *debug)
		}
	}
}

func fetchProducts(url string) ([]Product, error) {
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

	var products []Product
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, err
	}

	return products, nil
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

	// Try unmarshal as wrapper
	var result struct {
		Data []struct {
			ID int `json:"id"`
		} `json:"data"`
	}

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

func runTest(baseURL, name, pathTmpl, method string, isDetail bool, ids []int, totalReqs, concurrency int, debug bool) {
	fmt.Printf("Testing %s (%s)...\n", name, method)

	var wg sync.WaitGroup
	if totalReqs < concurrency {
		concurrency = totalReqs
	}
	requestsPerWorker := totalReqs / concurrency
	remainder := totalReqs % concurrency

	start := time.Now()
	var successCount, failCount int64
	var mu sync.Mutex

	// Search keywords for variety
	searchKeywords := []string{"Laptop", "Mouse", "Keyboard", "Coffee", "Tea", "Noodles", "Shirt", "Jeans"}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

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
				if isDetail && len(ids) > 0 {
					id := ids[rand.Intn(len(ids))]
					url = fmt.Sprintf(baseURL+pathTmpl, id)
				} else if strings.Contains(pathTmpl, "name=") {
					// Search endpoint - pick random keyword
					keyword := searchKeywords[rand.Intn(len(searchKeywords))]
					url = fmt.Sprintf(baseURL+pathTmpl, keyword)
				}

				req, err := http.NewRequest(method, url, nil)
				if err != nil {
					localFail++
					continue
				}

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					localFail++
					continue
				}

				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					localSuccess++
				} else {
					localFail++
					if debug {
						fmt.Printf("[%s] Request Failed: %s | Status: %d\n", name, url, resp.StatusCode)
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

func runPostTest(baseURL, name, path string, ids []int, totalReqs, concurrency int, debug bool) {
	fmt.Printf("Testing %s (POST)...\n", name)

	// First, restock all products to ensure sufficient inventory
	fmt.Printf("Restocking products...\n")
	products, err := restockProducts(baseURL)
	if err != nil {
		fmt.Printf("Error restocking products: %v\n", err)
		return
	}

	if len(products) == 0 {
		fmt.Printf("No products available for testing\n")
		return
	}

	fmt.Printf("Restocked %d products with 1000 units each\n", len(products))

	// Create a stock tracker to avoid exceeding available stock
	stockTracker := make(map[int]int)
	for _, p := range products {
		stockTracker[p.ID] = p.Stock
	}

	var wg sync.WaitGroup
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

		tasks := requestsPerWorker
		if i < remainder {
			tasks++
		}

		go func(count int) {
			defer wg.Done()
			localSuccess := 0
			localFail := 0

			for j := 0; j < count; j++ {
				// Create random checkout payload
				numItems := rand.Intn(3) + 1 // 1-3 items per transaction
				items := make([]map[string]interface{}, 0, numItems)

				// Track products used in this transaction to avoid duplicates
				usedProducts := make(map[int]bool)

				for k := 0; k < numItems && len(items) < numItems; k++ {
					// Find a product with available stock
					attempts := 0
					for attempts < 10 {
						productIdx := rand.Intn(len(products))
						product := products[productIdx]

						// Skip if already used in this transaction
						if usedProducts[product.ID] {
							attempts++
							continue
						}

						mu.Lock()
						availableStock := stockTracker[product.ID]
						mu.Unlock()

						if availableStock >= 1 {
							// Calculate safe quantity (max 3 or available stock)
							maxQty := 3
							if availableStock < maxQty {
								maxQty = availableStock
							}
							quantity := rand.Intn(maxQty) + 1

							mu.Lock()
							stockTracker[product.ID] -= quantity
							mu.Unlock()

							items = append(items, map[string]interface{}{
								"product_id": product.ID,
								"quantity":   quantity,
							})
							usedProducts[product.ID] = true
							break
						}
						attempts++
					}
				}

				if len(items) == 0 {
					// No products with stock available, skip this request
					continue
				}

				payload := map[string]interface{}{
					"items": items,
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					localFail++
					continue
				}

				url := baseURL + path
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
				if err != nil {
					localFail++
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					localFail++
					continue
				}

				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					localSuccess++
				} else {
					localFail++
					if debug {
						body, _ := io.ReadAll(resp.Body)
						fmt.Printf("[%s] Request Failed: %s | Status: %d | Body: %s\n", name, url, resp.StatusCode, string(body))
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

func restockProducts(baseURL string) ([]Product, error) {
	// Fetch current products
	products, err := fetchProducts(baseURL + "/products")
	if err != nil {
		return nil, err
	}

	// Update each product's stock to 1000
	for _, product := range products {
		payload := map[string]interface{}{
			"name":        fmt.Sprintf("Product-%d", product.ID),
			"price":       1000,
			"stock":       1000,
			"category_id": 1,
		}

		jsonBody, err := json.Marshal(payload)
		if err != nil {
			continue
		}

		url := fmt.Sprintf("%s/products/%d", baseURL, product.ID)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
	}

	// Fetch updated products
	return fetchProducts(baseURL + "/products")
}
