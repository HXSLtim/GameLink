package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func main() {
	baseURL := flag.String("base", "http://localhost:8081", "Base URL")
	email := flag.String("email", "admin@gameLink.com", "Admin email")
	password := flag.String("password", "Admin2025@Pass#", "Admin password")
	endpoint := flag.String("endpoint", "/api/v1/admin/games", "Endpoint to test")
	concurrency := flag.Int("c", 50, "Number of concurrent requests")
	duration := flag.Int("d", 10, "Duration in seconds")
	flag.Parse()

	// 1. 先登录获取 token
	fmt.Println("Logging in...")
	token, err := login(*baseURL, *email, *password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	fmt.Printf("Login successful, token: %s...\n\n", token[:20])

	// 2. 压测业务接口
	url := *baseURL + *endpoint
	fmt.Printf("Benchmarking %s\n", url)
	fmt.Printf("Concurrency: %d, Duration: %ds\n\n", *concurrency, *duration)

	var (
		totalRequests int64
		successCount  int64
		errorCount    int64
		totalLatency  int64
		minLatency    int64 = 1<<63 - 1
		maxLatency    int64
		status429     int64
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        *concurrency * 2,
			MaxIdleConnsPerHost: *concurrency * 2,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					start := time.Now()

					req, _ := http.NewRequest("GET", url, nil)
					req.Header.Set("Authorization", "Bearer "+token)
					req.Header.Set("Content-Type", "application/json")

					resp, err := client.Do(req)
					latency := time.Since(start).Microseconds()

					atomic.AddInt64(&totalRequests, 1)
					atomic.AddInt64(&totalLatency, latency)

					for {
						old := atomic.LoadInt64(&minLatency)
						if latency >= old || atomic.CompareAndSwapInt64(&minLatency, old, latency) {
							break
						}
					}

					for {
						old := atomic.LoadInt64(&maxLatency)
						if latency <= old || atomic.CompareAndSwapInt64(&maxLatency, old, latency) {
							break
						}
					}

					if err != nil {
						atomic.AddInt64(&errorCount, 1)
						continue
					}

					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()

					if resp.StatusCode == 200 {
						atomic.AddInt64(&successCount, 1)
					} else if resp.StatusCode == 429 {
						atomic.AddInt64(&status429, 1)
					} else {
						atomic.AddInt64(&errorCount, 1)
					}
				}
			}
		}()
	}

	time.Sleep(time.Duration(*duration) * time.Second)
	close(done)
	wg.Wait()

	total := atomic.LoadInt64(&totalRequests)
	success := atomic.LoadInt64(&successCount)
	errors := atomic.LoadInt64(&errorCount)
	rateLimit := atomic.LoadInt64(&status429)
	avgLatency := float64(atomic.LoadInt64(&totalLatency)) / float64(total) / 1000
	minLat := float64(atomic.LoadInt64(&minLatency)) / 1000
	maxLat := float64(atomic.LoadInt64(&maxLatency)) / 1000
	qps := float64(total) / float64(*duration)
	successQPS := float64(success) / float64(*duration)

	fmt.Println("========== Results ==========")
	fmt.Printf("Total Requests:  %d\n", total)
	fmt.Printf("Successful:      %d (%.2f%%)\n", success, float64(success)/float64(total)*100)
	fmt.Printf("Rate Limited:    %d (%.2f%%)\n", rateLimit, float64(rateLimit)/float64(total)*100)
	fmt.Printf("Errors:          %d (%.2f%%)\n", errors, float64(errors)/float64(total)*100)
	fmt.Printf("Total QPS:       %.2f req/s\n", qps)
	fmt.Printf("Success QPS:     %.2f req/s\n", successQPS)
	fmt.Printf("Avg Latency:     %.2f ms\n", avgLatency)
	fmt.Printf("Min Latency:     %.2f ms\n", minLat)
	fmt.Printf("Max Latency:     %.2f ms\n", maxLat)
}

func login(baseURL, email, password string) (string, error) {
	loginReq := LoginRequest{Username: email, Password: password}
	body, _ := json.Marshal(loginReq)

	resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体用于调试
	respBody, _ := io.ReadAll(resp.Body)

	var loginResp LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return "", fmt.Errorf("decode error: %v, body: %s", err, string(respBody))
	}

	if !loginResp.Success || loginResp.Data.Token == "" {
		return "", fmt.Errorf("login failed, response: %s", string(respBody))
	}

	return loginResp.Data.Token, nil
}
