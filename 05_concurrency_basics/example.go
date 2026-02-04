package main

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 示例1：基本 Goroutine ====================

func sayHello(name string) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("%s: Hello %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func demoBasicGoroutine() {
	fmt.Println("==================== 示例1：基本 Goroutine ====================")

	// 启动 goroutine
	go sayHello("Goroutine-1")
	go sayHello("Goroutine-2")

	// 主 goroutine 也执行
	sayHello("Main")

	time.Sleep(500 * time.Millisecond) // 等待 goroutine 完成
	fmt.Println()
}

// ==================== 示例2：Goroutine 参数传递 ====================

func demoGoroutineParameters() {
	fmt.Println("==================== 示例2：Goroutine 参数传递 ====================")

	// ✅ 正确方式：传递参数
	fmt.Println("正确方式:")
	for i := 0; i < 5; i++ {
		go func(n int) {
			fmt.Printf("参数传递: %d\n", n)
		}(i)
	}
	time.Sleep(100 * time.Millisecond)

	// ❌ 错误方式：闭包捕获（可能都打印 5）
	fmt.Println("\n错误方式（闭包陷阱）:")
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Printf("闭包捕获: %d\n", i) // 可能都是 5
		}()
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// ==================== 示例3：无缓冲 Channel ====================

func demoUnbufferedChannel() {
	fmt.Println("==================== 示例3：无缓冲 Channel ====================")

	ch := make(chan int) // 无缓冲 channel

	// 发送者
	go func() {
		fmt.Println("发送者: 准备发送数据...")
		ch <- 42
		fmt.Println("发送者: 数据已发送")
	}()

	time.Sleep(500 * time.Millisecond) // 模拟延迟

	// 接收者
	fmt.Println("接收者: 准备接收数据...")
	value := <-ch
	fmt.Printf("接收者: 收到数据 %d\n\n", value)
}

// ==================== 示例4：有缓冲 Channel ====================

func demoBufferedChannel() {
	fmt.Println("==================== 示例4：有缓冲 Channel ====================")

	ch := make(chan int, 3) // 缓冲区大小为 3

	// 可以连续发送 3 次，不会阻塞
	fmt.Println("发送 3 个数据到缓冲 channel...")
	ch <- 1
	fmt.Println("已发送: 1")
	ch <- 2
	fmt.Println("已发送: 2")
	ch <- 3
	fmt.Println("已发送: 3")

	fmt.Println("\n接收数据...")
	fmt.Printf("接收: %d\n", <-ch)
	fmt.Printf("接收: %d\n", <-ch)
	fmt.Printf("接收: %d\n", <-ch)
	fmt.Println()
}

// ==================== 示例5：关闭 Channel ====================

func demoCloseChannel() {
	fmt.Println("==================== 示例5：关闭 Channel ====================")

	ch := make(chan int, 3)

	// 发送数据
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch) // 关闭 channel

	// 可以继续接收
	fmt.Println("从已关闭的 channel 接收:")
	fmt.Println(<-ch) // 1
	fmt.Println(<-ch) // 2
	fmt.Println(<-ch) // 3

	// channel 为空后，接收零值
	value, ok := <-ch
	if !ok {
		fmt.Printf("Channel 已关闭且为空，收到零值: %d\n\n", value)
	}
}

// ==================== 示例6：遍历 Channel ====================

func demoRangeChannel() {
	fmt.Println("==================== 示例6：遍历 Channel ====================")

	ch := make(chan int, 5)

	// 发送者
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
		}
		close(ch) // 必须关闭
	}()

	// 接收者使用 range
	fmt.Println("使用 range 接收:")
	for value := range ch {
		fmt.Printf("接收: %d\n", value)
	}
	fmt.Println("所有数据接收完毕\n")
}

// ==================== 示例7：Select 基础 ====================

func demoSelect() {
	fmt.Println("==================== 示例7：Select 基础 ====================")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine 1
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自 ch1"
	}()

	// Goroutine 2
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自 ch2"
	}()

	// Select 等待第一个就绪的 channel
	select {
	case msg1 := <-ch1:
		fmt.Println("收到:", msg1)
	case msg2 := <-ch2:
		fmt.Println("收到:", msg2)
	}
	fmt.Println()
}

// ==================== 示例8：Select 超时控制 ====================

func demoSelectTimeout() {
	fmt.Println("==================== 示例8：Select 超时控制 ====================")

	ch := make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- 42
	}()

	select {
	case value := <-ch:
		fmt.Printf("收到数据: %d\n", value)
	case <-time.After(1 * time.Second):
		fmt.Println("超时！1秒内未收到数据")
	}
	fmt.Println()
}

// ==================== 示例9：Select 非阻塞 ====================

func demoSelectNonBlocking() {
	fmt.Println("==================== 示例9：Select 非阻塞 ====================")

	ch := make(chan int)

	select {
	case value := <-ch:
		fmt.Printf("收到: %d\n", value)
	default:
		fmt.Println("channel 未就绪，执行默认操作")
	}
	fmt.Println()
}

// ==================== 示例10：WaitGroup ====================

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 函数结束时调用

	fmt.Printf("Worker %d: 开始工作\n", id)
	time.Sleep(time.Duration(id*100) * time.Millisecond)
	fmt.Printf("Worker %d: 完成工作\n", id)
}

func demoWaitGroup() {
	fmt.Println("==================== 示例10：WaitGroup ====================")

	var wg sync.WaitGroup

	// 启动 5 个 worker
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	fmt.Println("等待所有 worker 完成...")
	wg.Wait()
	fmt.Println("所有任务完成\n")
}

// ==================== 示例11：生产者-消费者模式 ====================

func producer(ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		ch <- i
		fmt.Printf("生产者: 生产 %d\n", i)
		time.Sleep(100 * time.Millisecond)
	}
	close(ch)
	fmt.Println("生产者: 完成生产")
}

func consumer(ch <-chan int, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for value := range ch {
		fmt.Printf("消费者 %d: 消费 %d\n", id, value)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("消费者 %d: 完成消费\n", id)
}

func demoProducerConsumer() {
	fmt.Println("==================== 示例11：生产者-消费者模式 ====================")

	ch := make(chan int, 5)
	var wg sync.WaitGroup

	// 启动生产者
	go producer(ch, 10)

	// 启动 2 个消费者
	wg.Add(2)
	go consumer(ch, 1, &wg)
	go consumer(ch, 2, &wg)

	wg.Wait()
	fmt.Println()
}

// ==================== 示例12：管道模式 ====================

// 生成数据
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// 平方
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 过滤偶数
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func demoPipeline() {
	fmt.Println("==================== 示例12：管道模式 ====================")

	// 构建管道: 生成 -> 平方 -> 过滤偶数
	nums := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(nums)
	filtered := filterEven(squared)

	// 消费结果
	fmt.Println("管道结果（平方后的偶数）:")
	for result := range filtered {
		fmt.Printf("%d ", result)
	}
	fmt.Println("\n")
}

// ==================== 示例13：Worker 池 ====================

type Job struct {
	ID   int
	Data string
}

type Result struct {
	Job Job
	Sum int
}

func workerPool(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d: 处理任务 %d\n", id, job.ID)
		time.Sleep(100 * time.Millisecond) // 模拟处理

		// 计算字符串长度作为结果
		result := Result{
			Job: job,
			Sum: len(job.Data),
		}
		results <- result
	}
}

func demoWorkerPool() {
	fmt.Println("==================== 示例13：Worker 池 ====================")

	jobs := make(chan Job, 10)
	results := make(chan Result, 10)
	var wg sync.WaitGroup

	// 启动 3 个 worker
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go workerPool(i, jobs, results, &wg)
	}

	// 发送任务
	go func() {
		for i := 1; i <= 9; i++ {
			jobs <- Job{ID: i, Data: fmt.Sprintf("Task-%d", i)}
		}
		close(jobs)
	}()

	// 关闭结果 channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	fmt.Println("\n收集结果:")
	for result := range results {
		fmt.Printf("任务 %d 完成，结果: %d\n", result.Job.ID, result.Sum)
	}
	fmt.Println()
}

// ==================== 主函数 ====================

func main() {
	demoBasicGoroutine()
	demoGoroutineParameters()
	demoUnbufferedChannel()
	demoBufferedChannel()
	demoCloseChannel()
	demoRangeChannel()
	demoSelect()
	demoSelectTimeout()
	demoSelectNonBlocking()
	demoWaitGroup()
	demoProducerConsumer()
	demoPipeline()
	demoWorkerPool()

	fmt.Println("==================== 所有示例完成 ====================")
}
