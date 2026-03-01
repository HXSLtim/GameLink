package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:8081/api/v1/healthz", "URL to benchmark")
	concurrency := flag.Int("c", 100, "Number of concurrent requests")
	duration := flag.Int("d", 10, "Duration in seconds")
	flag.Parse()

	fmt.Printf("Benchmarking %s\n", *url)
	fmt.Printf("Concurrency: %d, Duration: %ds\n\n", *concurrency, *duration)

	var (
		totalRequests int64
		successCount  int64
		errorCount    int64
		totalLatency  int64
		minLatency    int64 = 1<<63 - 1
		maxLatency    int64
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

	// 启动工作协程
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
					resp, err := client.Get(*url)
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
						// fmt.Printf("Error: %v\n", err)
						continue
					}

					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()

					if resp.StatusCode == 200 {
						atomic.AddInt64(&successCount, 1)
					} else if resp.StatusCode == 429 {
						// 限流，不计入错误
						atomic.AddInt64(&errorCount, 1)
					} else {
						atomic.AddInt64(&errorCount, 1)
					}
				}
			}
		}()
	}

	// 等待指定时间
	time.Sleep(time.Duration(*duration) * time.Second)
	close(done)
	wg.Wait()

	// 计算结果
	total := atomic.LoadInt64(&totalRequests)
	success := atomic.LoadInt64(&successCount)
	errors := atomic.LoadInt64(&errorCount)
	avgLatency := float64(atomic.LoadInt64(&totalLatency)) / float64(total) / 1000 // ms
	minLat := float64(atomic.LoadInt64(&minLatency)) / 1000
	maxLat := float64(atomic.LoadInt64(&maxLatency)) / 1000
	qps := float64(total) / float64(*duration)

	fmt.Println("========== Results ==========")
	fmt.Printf("Total Requests:  %d\n", total)
	fmt.Printf("Successful:      %d (%.2f%%)\n", success, float64(success)/float64(total)*100)
	fmt.Printf("Errors:          %d (%.2f%%)\n", errors, float64(errors)/float64(total)*100)
	fmt.Printf("QPS:             %.2f req/s\n", qps)
	fmt.Printf("Avg Latency:     %.2f ms\n", avgLatency)
	fmt.Printf("Min Latency:     %.2f ms\n", minLat)
	fmt.Printf("Max Latency:     %.2f ms\n", maxLat)
}
