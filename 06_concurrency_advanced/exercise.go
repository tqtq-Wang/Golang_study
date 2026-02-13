package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// URL 列表（模拟要爬取的网页）
var urls = []string{
	"https://golang.org",
	"https://github.com",
	"https://google.com",
	"https://baidu.com",
	"https://zhihu.com",
	"https://bilibili.com",
	"https://taobao.com",
	"https://jd.com",
	"https://weibo.com",
	"https://douyin.com",
	"https://golang.org", // 重复！
	"https://github.com", // 重复！
	"https://example1.com",
	"https://example2.com",
	"https://example3.com",
}

type CrawlResult struct {
	URL      string
	Duration time.Duration
	Success  bool
}

func crawl(url string) CrawlResult {
	start := time.Now()

	// 模拟网络延迟（200-800ms）
	delay := time.Duration(rand.Intn(600)+200) * time.Millisecond
	time.Sleep(delay)

	success := rand.Intn(10) > 0 // 90% 成功率

	return CrawlResult{
		URL:      url,
		Duration: time.Since(start),
		Success:  success,
	}
}

func worker(id int, jobs <-chan string, results chan<- CrawlResult, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for url := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker-%d: 收到取消信号\n", id)
			return
		default:
			fmt.Printf("Worker-%d: 爬取 %s\n", id, url)
			result := crawl(url)
			results <- result
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("===== 串行爬虫（Step 1）=====")
	start := time.Now()

	for _, url := range urls {
		result := crawl(url)
		if result.Success {
			fmt.Printf("成功爬取: %s (耗时: %v)\n", result.URL, result.Duration)
		} else {
			fmt.Printf("爬取失败: %s (耗时: %v)\n", result.URL, result.Duration)
		}
	}

	fmt.Printf("总耗时: %v\n\n", time.Since(start))

	// 保存串行耗时用于后续对比
	serialDuration := time.Since(start)

	fmt.Println("===== 并发爬虫（Step 2）=====")
	start = time.Now()

	// 为并发爬虫创建带10秒超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobs := make(chan string, 10)
	results := make(chan CrawlResult, 15)
	var wg sync.WaitGroup

	// Step 5: 使用 atomic 统计成功/失败/跳过数量
	var successCount, failCount, skippedCount int64

	// 启动 5 个 worker
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg, ctx)
	}

	// 发送任务（带去重）
	go func() {
		visited := make(map[string]bool)

		for i, url := range urls {
			if visited[url] {
				fmt.Printf("⏭️  [%d] 跳过重复: %s\n", i+1, url)
				atomic.AddInt64(&skippedCount, 1)
				continue
			}
			visited[url] = true
			jobs <- url
		}
		close(jobs)
	}()

	// 收集结果
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Success {
			fmt.Printf("成功爬取: %s (耗时: %v)\n", result.URL, result.Duration)
			atomic.AddInt64(&successCount, 1)
		} else {
			fmt.Printf("爬取失败: %s (耗时: %v)\n", result.URL, result.Duration)
			atomic.AddInt64(&failCount, 1)
		}
	}

	// Step 5: 打印详细统计信息
	concurrentDuration := time.Since(start)
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 并发爬虫统计报告")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✅ 成功: %d 个\n", atomic.LoadInt64(&successCount))
	fmt.Printf("❌ 失败: %d 个\n", atomic.LoadInt64(&failCount))
	fmt.Printf("⏭️  跳过: %d 个\n", atomic.LoadInt64(&skippedCount))
	fmt.Printf("📈 总计: %d 个URL\n", len(urls))
	fmt.Printf("⏱️  串行耗时: %v\n", serialDuration)
	fmt.Printf("⚡ 并发耗时: %v\n", concurrentDuration)
	fmt.Printf("🚀 提速倍数: %.2fx\n", float64(serialDuration)/float64(concurrentDuration))
	fmt.Println(strings.Repeat("=", 50))
}
