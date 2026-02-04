package main

import (
	"fmt"
	"sync"
	"time"
)

// ========================================
// 版本1：最简单 - 只有1个worker
// ========================================

func simpleVersion() {
	fmt.Println("===== 版本1：最简单 =====")

	// 创建任务管道
	jobs := make(chan int, 5)

	// 启动1个工人
	go func() {
		// 从管道取任务
		for job := range jobs {
			fmt.Printf("处理任务 %d\n", job)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("工人下班了")
	}()

	// 发送3个任务
	jobs <- 1
	jobs <- 2
	jobs <- 3

	// 关闭管道（告诉工人：没活了）
	close(jobs)

	// 等待工人处理完
	time.Sleep(500 * time.Millisecond)
}

// ========================================
// 版本2：加上 WaitGroup
// ========================================

func withWaitGroup() {
	fmt.Println("\n===== 版本2：加上WaitGroup =====")

	jobs := make(chan int, 5)
	var wg sync.WaitGroup

	// 启动1个工人
	wg.Add(1) // 告诉系统：有1个任务
	go func() {
		defer wg.Done() // 下班时打卡

		for job := range jobs {
			fmt.Printf("处理任务 %d\n", job)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("工人下班了")
	}()

	// 发送任务
	jobs <- 1
	jobs <- 2
	jobs <- 3

	close(jobs)

	// 等待工人（不用猜时间了！）
	wg.Wait()
	fmt.Println("确认工人已下班")
}

// ========================================
// 版本3：多个worker
// ========================================

func multipleWorkers() {
	fmt.Println("\n===== 版本3：多个worker =====")

	jobs := make(chan int, 5)
	var wg sync.WaitGroup

	// 启动3个工人
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()

			for job := range jobs {
				fmt.Printf("工人%d 处理任务 %d\n", workerID, job)
				time.Sleep(100 * time.Millisecond)
			}
			fmt.Printf("工人%d 下班了\n", workerID)
		}()
	}

	// 发送6个任务
	for i := 1; i <= 6; i++ {
		jobs <- i
	}

	close(jobs)
	wg.Wait()
	fmt.Println("所有工人已下班")
}

// ========================================
// 版本4：加上结果收集
// ========================================

func withResultCollector() {
	fmt.Println("\n===== 版本4：加上结果收集 =====")

	jobs := make(chan int, 5)
	results := make(chan int, 10)
	done := make(chan bool)
	var wg sync.WaitGroup

	// 启动2个工人
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()

			for job := range jobs {
				fmt.Printf("工人%d 处理 %d\n", workerID, job)
				result := job * 2 // 简单处理：乘以2
				results <- result
			}
		}()
	}

	// 启动收集器
	go func() {
		sum := 0
		for result := range results {
			sum += result
			fmt.Printf("收集结果: %d (累计: %d)\n", result, sum)
		}
		fmt.Printf("总和: %d\n", sum)
		done <- true
	}()

	// 发送任务
	for i := 1; i <= 5; i++ {
		jobs <- i
	}

	// 关闭顺序（重点！）
	close(jobs)    // 1. 不发新任务
	wg.Wait()      // 2. 等工人完成
	close(results) // 3. 关闭结果管道
	<-done         // 4. 等收集器完成

	fmt.Println("全部完成")
}

// ========================================
// 版本5：加上超时控制
// ========================================

func withTimeout() {
	fmt.Println("\n===== 版本5：加上超时控制 =====")

	jobs := make(chan int, 5)
	var wg sync.WaitGroup

	// 启动1个工人
	wg.Add(1)
	go func() {
		defer wg.Done()
		for job := range jobs {
			fmt.Printf("处理任务 %d\n", job)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 发送任务（带超时）
	timeout := time.After(500 * time.Millisecond)

TaskLoop:
	for i := 1; i <= 10; i++ {
		select {
		case <-timeout:
			fmt.Println("⏰ 超时！停止发送")
			break TaskLoop
		case jobs <- i:
			fmt.Printf("发送任务 %d\n", i)
			time.Sleep(200 * time.Millisecond)
		}
	}

	close(jobs)
	wg.Wait()
	fmt.Println("完成")
}

// ========================================
// 主函数：按顺序运行所有版本
// ========================================

func main() {
	simpleVersion()

	time.Sleep(1 * time.Second)
	withWaitGroup()

	time.Sleep(1 * time.Second)
	multipleWorkers()

	time.Sleep(1 * time.Second)
	withResultCollector()

	time.Sleep(1 * time.Second)
	withTimeout()
}
