package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== 示例1：Context 基本使用 ====================

func demoContextBasic() {
	fmt.Println("==================== 示例1：Context 基本使用 ====================")

	// WithCancel：手动取消
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		fmt.Println("收到取消信号:", ctx.Err())
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("发送取消信号")
	cancel()

	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

// ==================== 示例2：Context 超时控制 ====================

func demoContextTimeout() {
	fmt.Println("==================== 示例2：Context 超时控制 ====================")

	// 2秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go longRunningTask(ctx, "任务1")
	go longRunningTask(ctx, "任务2")

	time.Sleep(3 * time.Second)
	fmt.Println()
}

func longRunningTask(ctx context.Context, name string) {
	fmt.Printf("%s: 开始执行\n", name)

	select {
	case <-time.After(5 * time.Second):
		fmt.Printf("%s: 正常完成\n", name)
	case <-ctx.Done():
		fmt.Printf("%s: 被取消 (%v)\n", name, ctx.Err())
	}
}

// ==================== 示例3：Context 级联取消 ====================

func demoContextCascade() {
	fmt.Println("==================== 示例3：Context 级联取消 ====================")

	// 根 context
	rootCtx, rootCancel := context.WithCancel(context.Background())

	// 启动多层任务
	go level1(rootCtx)

	// 2秒后取消
	time.Sleep(2 * time.Second)
	fmt.Println(">>> 主程序：取消根 context")
	rootCancel()

	time.Sleep(1 * time.Second)
	fmt.Println()
}

func level1(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("Level 1: 启动")
	go level2(ctx)

	<-ctx.Done()
	fmt.Println("Level 1: 收到取消信号")
}

func level2(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("Level 2: 启动")
	go level3(ctx)

	<-ctx.Done()
	fmt.Println("Level 2: 收到取消信号")
}

func level3(ctx context.Context) {
	fmt.Println("Level 3: 启动")

	<-ctx.Done()
	fmt.Println("Level 3: 收到取消信号")
}

// ==================== 示例4：竞态条件演示 ====================

func demoRaceCondition() {
	fmt.Println("==================== 示例4：竞态条件演示 ====================")
	fmt.Println("提示：运行 'go run -race example.go' 可以检测数据竞争")

	counter := 0
	var wg sync.WaitGroup

	// 启动 10 个 goroutine 同时修改 counter
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter++ // ❌ 不安全！
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d (期望: 1000, 实际可能不是)\n", counter)
	fmt.Println()
}

// ==================== 示例5：使用 Mutex 修复竞态 ====================

func demoMutex() {
	fmt.Println("==================== 示例5：使用 Mutex 修复竞态 ====================")

	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d (总是 1000)\n", counter)
	fmt.Println()
}

// ==================== 示例6：RWMutex 读写锁 ====================

type SafeCounter struct {
	count int
	mu    sync.RWMutex
}

func (c *SafeCounter) Get() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func demoRWMutex() {
	fmt.Println("==================== 示例6：RWMutex 读写锁 ====================")

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 5 个写 goroutine
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				counter.Increment()
				fmt.Printf("Writer %d: 写入\n", id)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// 10 个读 goroutine（可以并发读）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				val := counter.Get()
				fmt.Printf("Reader %d: 读取 %d\n", id, val)
				time.Sleep(5 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("最终值: %d\n", counter.Get())
	fmt.Println()
}

// ==================== 示例7：使用 Channel 实现并发安全 ====================

type ChannelCounter struct {
	ops chan func(int) int
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		ops: make(chan func(int) int),
	}

	// 启动一个 goroutine 串行处理所有操作
	go func() {
		count := 0
		for op := range c.ops {
			count = op(count)
		}
	}()

	return c
}

func (c *ChannelCounter) Increment() {
	c.ops <- func(count int) int {
		return count + 1
	}
}

func (c *ChannelCounter) Get() int {
	result := make(chan int)
	c.ops <- func(count int) int {
		result <- count
		return count
	}
	return <-result
}

func demoChannelSafety() {
	fmt.Println("==================== 示例7：使用 Channel 实现并发安全 ====================")

	counter := NewChannelCounter()
	var wg sync.WaitGroup

	// 10 个 goroutine 同时增加
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d (总是 1000)\n", counter.Get())
	fmt.Println()
}

// ==================== 示例8：Atomic 原子操作 ====================

func demoAtomic() {
	fmt.Println("==================== 示例8：Atomic 原子操作 ====================")

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d (总是 1000)\n", atomic.LoadInt64(&counter))
	fmt.Println()
}

// ==================== 示例9：Context 实战 - HTTP 请求模拟 ====================

func demoContextPractical() {
	fmt.Println("==================== 示例9：Context 实战 - HTTP 请求模拟 ====================")

	// 模拟处理 HTTP 请求
	ctx := context.Background()
	handleRequest(ctx, 123)

	fmt.Println()
}

func handleRequest(ctx context.Context, userID int) {
	// 设置 3 秒超时
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// 任务1：查询数据库
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := queryDatabase(ctx, userID); err != nil {
			fmt.Printf("❌ 数据库查询失败: %v\n", err)
		} else {
			fmt.Println("✅ 数据库查询成功")
		}
	}()

	// 任务2：调用外部 API
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := callExternalAPI(ctx, userID); err != nil {
			fmt.Printf("❌ API 调用失败: %v\n", err)
		} else {
			fmt.Println("✅ API 调用成功")
		}
	}()

	// 任务3：发送通知
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := sendNotification(ctx, userID); err != nil {
			fmt.Printf("❌ 发送通知失败: %v\n", err)
		} else {
			fmt.Println("✅ 发送通知成功")
		}
	}()

	wg.Wait()
}

func queryDatabase(ctx context.Context, userID int) error {
	select {
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func callExternalAPI(ctx context.Context, userID int) error {
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func sendNotification(ctx context.Context, userID int) error {
	select {
	case <-time.After(5 * time.Second): // 模拟慢请求
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ==================== 示例10：并发安全的 Map ====================

type SafeMap struct {
	data map[string]int
	mu   sync.RWMutex
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

func (m *SafeMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

func demoSafeMap() {
	fmt.Println("==================== 示例10：并发安全的 Map ====================")

	safeMap := NewSafeMap()
	var wg sync.WaitGroup

	// 10 个 goroutine 写入
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			safeMap.Set(key, id*10)
			fmt.Printf("写入: %s = %d\n", key, id*10)
		}(i)
	}

	wg.Wait()

	// 读取所有值
	fmt.Println("\n读取所有值:")
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		if val, ok := safeMap.Get(key); ok {
			fmt.Printf("%s = %d\n", key, val)
		}
	}

	fmt.Println()
}

// ==================== 主函数 ====================

func main() {
	demoContextBasic()

	time.Sleep(500 * time.Millisecond)
	demoContextTimeout()

	time.Sleep(500 * time.Millisecond)
	demoContextCascade()

	time.Sleep(500 * time.Millisecond)
	demoRaceCondition()

	time.Sleep(500 * time.Millisecond)
	demoMutex()

	time.Sleep(500 * time.Millisecond)
	demoRWMutex()

	time.Sleep(500 * time.Millisecond)
	demoChannelSafety()

	time.Sleep(500 * time.Millisecond)
	demoAtomic()

	time.Sleep(500 * time.Millisecond)
	demoContextPractical()

	time.Sleep(500 * time.Millisecond)
	demoSafeMap()

	fmt.Println("==================== 所有示例完成 ====================")
	fmt.Println("\n提示：运行 'go run -race example.go' 可以检测数据竞争！")
}
