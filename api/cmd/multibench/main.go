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

// BenchmarkResult 压测结果
type BenchmarkResult struct {
	Endpoint      string
	TotalRequests int64
	SuccessCount  int64
	ErrorCount    int64
	RateLimited   int64
	TotalQPS      float64
	SuccessQPS    float64
	AvgLatency    float64
	MinLatency    float64
	MaxLatency    float64
	Duration      int
}

// APIClient 带认证的 API 客户端
type APIClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewAPIClient 创建 API 客户端
func NewAPIClient(baseURL string, concurrency int) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        concurrency * 2,
				MaxIdleConnsPerHost: concurrency * 2,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// Login 登录获取 token
func (c *APIClient) Login(email, password string) error {
	loginReq := map[string]string{"username": email, "password": password}
	body, _ := json.Marshal(loginReq)

	resp, err := c.httpClient.Post(c.baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var loginResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return fmt.Errorf("decode error: %v, body: %s", err, string(respBody))
	}

	if !loginResp.Success || loginResp.Data.Token == "" {
		return fmt.Errorf("login failed, response: %s", string(respBody))
	}

	c.token = loginResp.Data.Token
	return nil
}

// Get 发送 GET 请求
func (c *APIClient) Get(endpoint string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", c.baseURL+endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

// Post 发送 POST 请求
func (c *APIClient) Post(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewReader(jsonBody))
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

// BenchmarkEndpoint 压测单个接口
func BenchmarkEndpoint(client *APIClient, endpoint, method string, body interface{}, concurrency, duration int) BenchmarkResult {
	var (
		totalRequests int64
		successCount  int64
		errorCount    int64
		rateLimited   int64
		totalLatency  int64
		minLatency    int64 = 1<<63 - 1
		maxLatency    int64
	)

	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					start := time.Now()
					var resp *http.Response
					var err error

					if method == "GET" {
						resp, err = client.Get(endpoint)
					} else {
						resp, err = client.Post(endpoint, body)
					}

					latency := time.Since(start).Microseconds()
					atomic.AddInt64(&totalRequests, 1)
					atomic.AddInt64(&totalLatency, latency)

					// 更新最小延迟
					for {
						old := atomic.LoadInt64(&minLatency)
						if latency >= old || atomic.CompareAndSwapInt64(&minLatency, old, latency) {
							break
						}
					}

					// 更新最大延迟
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

					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						atomic.AddInt64(&successCount, 1)
					} else if resp.StatusCode == 429 {
						atomic.AddInt64(&rateLimited, 1)
					} else {
						atomic.AddInt64(&errorCount, 1)
					}
				}
			}
		}()
	}

	time.Sleep(time.Duration(duration) * time.Second)
	close(done)
	wg.Wait()

	total := atomic.LoadInt64(&totalRequests)
	success := atomic.LoadInt64(&successCount)
	errors := atomic.LoadInt64(&errorCount)
	limited := atomic.LoadInt64(&rateLimited)
	avgLatency := float64(atomic.LoadInt64(&totalLatency)) / float64(total) / 1000
	minLat := float64(atomic.LoadInt64(&minLatency)) / 1000
	maxLat := float64(atomic.LoadInt64(&maxLatency)) / 1000

	return BenchmarkResult{
		Endpoint:      endpoint,
		TotalRequests: total,
		SuccessCount:  success,
		ErrorCount:    errors,
		RateLimited:   limited,
		TotalQPS:      float64(total) / float64(duration),
		SuccessQPS:    float64(success) / float64(duration),
		AvgLatency:    avgLatency,
		MinLatency:    minLat,
		MaxLatency:    maxLat,
		Duration:      duration,
	}
}

// PrintResult 打印压测结果
func PrintResult(r BenchmarkResult) {
	fmt.Printf("\n========== %s ==========\n", r.Endpoint)
	fmt.Printf("Duration:        %ds\n", r.Duration)
	fmt.Printf("Total Requests:  %d\n", r.TotalRequests)
	fmt.Printf("Successful:      %d (%.2f%%)\n", r.SuccessCount, float64(r.SuccessCount)/float64(r.TotalRequests)*100)
	fmt.Printf("Rate Limited:    %d (%.2f%%)\n", r.RateLimited, float64(r.RateLimited)/float64(r.TotalRequests)*100)
	fmt.Printf("Errors:          %d (%.2f%%)\n", r.ErrorCount, float64(r.ErrorCount)/float64(r.TotalRequests)*100)
	fmt.Printf("Total QPS:       %.2f req/s\n", r.TotalQPS)
	fmt.Printf("Success QPS:     %.2f req/s\n", r.SuccessQPS)
	fmt.Printf("Avg Latency:     %.2f ms\n", r.AvgLatency)
	fmt.Printf("Min Latency:     %.2f ms\n", r.MinLatency)
	fmt.Printf("Max Latency:     %.2f ms\n", r.MaxLatency)
}

func main() {
	baseURL := flag.String("base", "http://localhost:8082", "Base URL")
	email := flag.String("email", "admin@gameLink.com", "Admin email")
	password := flag.String("password", "Admin2025@Pass#", "Admin password")
	concurrency := flag.Int("c", 50, "Number of concurrent requests")
	duration := flag.Int("d", 10, "Duration in seconds per endpoint")
	testAll := flag.Bool("all", false, "Test all endpoints")
	endpoint := flag.String("endpoint", "", "Single endpoint to test")
	flag.Parse()

	client := NewAPIClient(*baseURL, *concurrency)

	// 登录
	fmt.Println("Logging in...")
	if err := client.Login(*email, *password); err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	fmt.Printf("Login successful!\n")

	// 定义测试接口
	endpoints := []struct {
		name     string
		endpoint string
		method   string
		body     interface{}
	}{
		{"健康检查", "/api/v1/healthz", "GET", nil},
		{"游戏列表", "/api/v1/admin/games", "GET", nil},
		{"用户列表", "/api/v1/admin/users", "GET", nil},
		{"订单列表", "/api/v1/admin/orders", "GET", nil},
		{"角色列表", "/api/v1/admin/roles", "GET", nil},
		{"权限列表", "/api/v1/admin/permissions", "GET", nil},
		{"聊天群组", "/api/v1/user/chat/groups", "GET", nil},
		{"陪玩师列表", "/api/v1/user/players", "GET", nil},
		{"通知列表", "/api/v1/user/notifications", "GET", nil},
	}

	if *testAll {
		fmt.Printf("\n开始全量压测，并发数: %d，每接口持续: %ds\n", *concurrency, *duration)
		fmt.Println("============================================")

		var results []BenchmarkResult
		for _, ep := range endpoints {
			fmt.Printf("\n正在测试: %s (%s)...\n", ep.name, ep.endpoint)
			result := BenchmarkEndpoint(client, ep.endpoint, ep.method, ep.body, *concurrency, *duration)
			results = append(results, result)
			PrintResult(result)
		}

		// 汇总报告
		fmt.Println("\n\n============ 汇总报告 ============")
		fmt.Printf("%-30s %10s %10s %10s %10s\n", "接口", "总QPS", "成功QPS", "平均延迟", "成功率")
		fmt.Println("--------------------------------------------------------------------------------")
		for i, r := range results {
			successRate := float64(r.SuccessCount) / float64(r.TotalRequests) * 100
			fmt.Printf("%-30s %10.2f %10.2f %8.2fms %8.2f%%\n",
				endpoints[i].name, r.TotalQPS, r.SuccessQPS, r.AvgLatency, successRate)
		}
	} else if *endpoint != "" {
		fmt.Printf("\n测试单个接口: %s\n", *endpoint)
		result := BenchmarkEndpoint(client, *endpoint, "GET", nil, *concurrency, *duration)
		PrintResult(result)
	} else {
		fmt.Println("\n使用方法:")
		fmt.Println("  -all          测试所有接口")
		fmt.Println("  -endpoint     测试单个接口")
		fmt.Println("  -c            并发数 (默认 50)")
		fmt.Println("  -d            持续时间秒 (默认 10)")
		fmt.Println("\n示例:")
		fmt.Println("  go run multi_benchmark.go -all")
		fmt.Println("  go run multi_benchmark.go -endpoint /api/v1/admin/games -c 100 -d 30")
	}
}
