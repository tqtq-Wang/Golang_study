# 第06节：并发编程进阶 - Context 与并发安全

> **本节目标**：掌握 Context 上下文管理、深入理解并发安全、学习竞态检测  
> **前置知识**：第05节（Goroutine、Channel、Select、WaitGroup）  
> **重要程度**：⭐⭐⭐⭐⭐ 企业级开发必备！

---

## 📖 一、Context 上下文管理

### 1.1 为什么需要 Context？

#### **场景问题**

假设你启动了一个 HTTP 请求，它内部启动了 10 个 Goroutine 去处理不同的任务：

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // 启动 10 个 goroutine 处理任务
    for i := 0; i < 10; i++ {
        go doSomething(i)
    }

    // 问题：如果用户取消了请求，这 10 个 goroutine 怎么停止？
    // 问题：如果请求超时了，怎么通知所有 goroutine？
}
```

**Java 的解决方案**：

```java
// Java 使用 Thread.interrupt() 或 ExecutorService.shutdown()
ExecutorService executor = Executors.newFixedThreadPool(10);
Future<?> future = executor.submit(task);
future.cancel(true);  // 取消任务
```

**Go 的解决方案**：使用 **Context**！

---

### 1.2 Context 的基本概念

**Context** 是 Go 的上下文管理机制，用于：

1. **取消信号传递**：通知 goroutine 停止工作
2. **超时控制**：自动取消超时的操作
3. **截止时间**：设置任务的最后期限
4. **传递请求范围的值**：（不推荐滥用）

**核心接口**：

```go
type Context interface {
    // 返回 context 的截止时间（如果有）
    Deadline() (deadline time.Time, ok bool)

    // 返回一个 channel，当 context 被取消时关闭
    Done() <-chan struct{}

    // 返回取消的原因
    Err() error

    // 返回与 context 关联的值
    Value(key interface{}) interface{}
}
```

---

### 1.3 Context 的创建

#### **1. Background 和 TODO**

```go
// Background：根 context，通常用于 main 函数、初始化、测试
ctx := context.Background()

// TODO：当不确定用什么 context 时使用（占位符）
ctx := context.TODO()
```

#### **2. WithCancel：手动取消**

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    // 监听取消信号
    <-ctx.Done()
    fmt.Println("任务被取消了:", ctx.Err())
}()

// 手动取消
time.Sleep(1 * time.Second)
cancel()  // 调用 cancel() 会关闭 Done() channel
```

**Java 对比**：

```java
// Java 需要手动检查中断标志
if (Thread.interrupted()) {
    throw new InterruptedException();
}
```

---

#### **3. WithTimeout：超时自动取消**

```go
// 3秒后自动取消
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()  // 养成好习惯：总是 defer cancel()

go func() {
    select {
    case <-time.After(5 * time.Second):
        fmt.Println("任务完成")
    case <-ctx.Done():
        fmt.Println("超时取消:", ctx.Err())  // context deadline exceeded
    }
}()

time.Sleep(4 * time.Second)
```

---

#### **4. WithDeadline：指定截止时间**

```go
// 指定绝对时间点
deadline := time.Now().Add(2 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

go func() {
    <-ctx.Done()
    fmt.Println("截止时间到:", ctx.Err())
}()

time.Sleep(3 * time.Second)
```

---

#### **5. WithValue：传递值（慎用！）**

```go
// 传递请求ID
ctx := context.WithValue(context.Background(), "requestID", "12345")

go func(ctx context.Context) {
    requestID := ctx.Value("requestID").(string)
    fmt.Println("请求ID:", requestID)
}(ctx)
```

**⚠️ 警告**：

- 不要用 Context 传递函数参数！
- 只用于传递请求范围的元数据（如请求ID、用户身份等）
- 优先使用函数参数，Context 只是补充

---

### 1.4 Context 的传递规则

#### **规则1：Context 树形结构**

```go
rootCtx := context.Background()

// 创建子 context
ctx1, cancel1 := context.WithTimeout(rootCtx, 5*time.Second)
ctx2, cancel2 := context.WithCancel(ctx1)

// 父 context 取消，所有子 context 也取消
cancel1()  // ctx1 和 ctx2 都会取消
```

**树形结构**：

```
Background
    ↓
WithTimeout(5s) ← cancel1() 取消这里
    ↓
WithCancel      ← 这里也会被取消
```

---

#### **规则2：Context 作为第一个参数**

```go
// ✅ 正确：ctx 作为第一个参数
func DoSomething(ctx context.Context, name string) error {
    // ...
}

// ❌ 错误：ctx 不在第一个位置
func DoSomething(name string, ctx context.Context) error {
    // ...
}
```

**Go 的约定**：Context 总是第一个参数，参数名通常是 `ctx`

---

### 1.5 Context 实战示例

#### **示例1：模拟 HTTP 请求处理**

```go
func HandleRequest(ctx context.Context, userID int) {
    // 设置 3 秒超时
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // 启动多个任务
    var wg sync.WaitGroup

    // 任务1：查询数据库
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := queryDatabase(ctx, userID); err != nil {
            fmt.Println("数据库查询失败:", err)
        }
    }()

    // 任务2：调用外部 API
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := callExternalAPI(ctx, userID); err != nil {
            fmt.Println("API 调用失败:", err)
        }
    }()

    // 等待所有任务完成或超时
    wg.Wait()
}

func queryDatabase(ctx context.Context, userID int) error {
    select {
    case <-time.After(2 * time.Second):
        fmt.Println("数据库查询完成")
        return nil
    case <-ctx.Done():
        return ctx.Err()  // 超时或取消
    }
}

func callExternalAPI(ctx context.Context, userID int) error {
    select {
    case <-time.After(5 * time.Second):  // 模拟慢请求
        fmt.Println("API 调用完成")
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

#### **示例2：级联取消**

```go
func main() {
    // 根 context
    rootCtx, rootCancel := context.WithCancel(context.Background())

    // 启动多层任务
    go layer1(rootCtx)

    // 3 秒后取消所有任务
    time.Sleep(3 * time.Second)
    fmt.Println("主程序：取消所有任务")
    rootCancel()  // 取消根 context

    time.Sleep(1 * time.Second)
}

func layer1(ctx context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    fmt.Println("Layer 1: 启动")
    go layer2(ctx)

    <-ctx.Done()
    fmt.Println("Layer 1: 收到取消信号")
}

func layer2(ctx context.Context) {
    fmt.Println("Layer 2: 启动")
    go layer3(ctx)

    <-ctx.Done()
    fmt.Println("Layer 2: 收到取消信号")
}

func layer3(ctx context.Context) {
    fmt.Println("Layer 3: 启动")

    <-ctx.Done()
    fmt.Println("Layer 3: 收到取消信号")
}
```

**输出**：

```
Layer 1: 启动
Layer 2: 启动
Layer 3: 启动
主程序：取消所有任务
Layer 1: 收到取消信号
Layer 2: 收到取消信号
Layer 3: 收到取消信号
```

**关键**：一个 `cancel()` 调用，所有子 context 都会收到信号！

---

## 📖 二、并发安全问题

### 2.1 什么是竞态条件（Race Condition）？

**竞态条件**：多个 goroutine 同时访问共享变量，至少有一个在写，导致结果不可预测。

#### **示例：银行账户并发扣款**

```go
type Account struct {
    balance int
}

func (a *Account) Withdraw(amount int) {
    // ❌ 不安全的实现
    if a.balance >= amount {
        time.Sleep(1 * time.Millisecond)  // 模拟处理时间
        a.balance -= amount
    }
}

func main() {
    account := &Account{balance: 100}

    // 10 个 goroutine 同时取钱
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            account.Withdraw(10)
        }()
    }

    wg.Wait()
    fmt.Println("剩余余额:", account.balance)
    // 期望：0
    // 实际：可能是 10, 20, 30... 不确定！
}
```

**问题分析**：

```
时间线：
goroutine1: 读取 balance=100 → 判断够 → 暂停
goroutine2: 读取 balance=100 → 判断够 → 暂停
goroutine3: 读取 balance=100 → 判断够 → 暂停
...
goroutine1: 恢复 → balance=90
goroutine2: 恢复 → balance=90（覆盖了！）
goroutine3: 恢复 → balance=90（又覆盖了！）
```

---

### 2.2 解决方案1：Mutex 互斥锁

**Mutex**（Mutual Exclusion）：同一时间只允许一个 goroutine 访问。

```go
import "sync"

type SafeAccount struct {
    balance int
    mu      sync.Mutex  // 互斥锁
}

func (a *SafeAccount) Withdraw(amount int) {
    a.mu.Lock()         // 加锁
    defer a.mu.Unlock() // 解锁

    if a.balance >= amount {
        time.Sleep(1 * time.Millisecond)
        a.balance -= amount
    }
}
```

**Java 对比**：

```java
// Java 使用 synchronized
public synchronized void withdraw(int amount) {
    if (balance >= amount) {
        balance -= amount;
    }
}
```

---

### 2.3 解决方案2：RWMutex 读写锁

**RWMutex**：允许多个读，但写时独占。

```go
type SafeCounter struct {
    count int
    mu    sync.RWMutex  // 读写锁
}

func (c *SafeCounter) Get() int {
    c.mu.RLock()         // 读锁（多个 goroutine 可以同时读）
    defer c.mu.RUnlock()
    return c.count
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()          // 写锁（独占）
    defer c.mu.Unlock()
    c.count++
}
```

**使用场景**：

- ✅ 读多写少：RWMutex 性能更好
- ✅ 写多读少：Mutex 更简单

---

### 2.4 解决方案3：Channel（推荐！）

**Go 的哲学**：用 Channel 代替锁！

```go
type Account struct {
    balance int
    ops     chan func(*Account)  // 操作 channel
}

func NewAccount(initial int) *Account {
    a := &Account{
        balance: initial,
        ops:     make(chan func(*Account)),
    }

    // 启动一个 goroutine 串行处理所有操作
    go func() {
        for op := range a.ops {
            op(a)  // 执行操作
        }
    }()

    return a
}

func (a *Account) Withdraw(amount int) {
    done := make(chan bool)

    a.ops <- func(a *Account) {
        if a.balance >= amount {
            a.balance -= amount
        }
        done <- true
    }

    <-done  // 等待操作完成
}

func (a *Account) GetBalance() int {
    result := make(chan int)

    a.ops <- func(a *Account) {
        result <- a.balance
    }

    return <-result
}
```

**优势**：

- ✅ 无需加锁
- ✅ 串行化访问，天然并发安全
- ✅ 符合 Go 的设计哲学

---

### 2.5 解决方案4：Atomic 原子操作

**适用场景**：简单的数值操作（加减、交换等）

```go
import "sync/atomic"

type Counter struct {
    count int64  // 必须是 int32 或 int64
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.count, 1)  // 原子操作
}

func (c *Counter) Get() int64 {
    return atomic.LoadInt64(&c.count)
}

func (c *Counter) Set(val int64) {
    atomic.StoreInt64(&c.count, val)
}
```

**Java 对比**：

```java
// Java 使用 AtomicInteger
AtomicInteger count = new AtomicInteger(0);
count.incrementAndGet();
```

---

## 📖 三、竞态检测（Race Detector）

### 3.1 什么是竞态检测？

Go 内置了**竞态检测器**，可以在运行时检测数据竞争。

#### **使用方法**

```bash
# 运行时加上 -race 参数
go run -race main.go

# 测试时检测
go test -race

# 编译带竞态检测的二进制
go build -race
```

---

### 3.2 竞态检测示例

```go
// race_example.go
package main

import (
    "fmt"
    "time"
)

func main() {
    counter := 0

    // 启动 2 个 goroutine 同时修改 counter
    go func() {
        for i := 0; i < 1000; i++ {
            counter++  // 写
        }
    }()

    go func() {
        for i := 0; i < 1000; i++ {
            counter++  // 写
        }
    }()

    time.Sleep(1 * time.Second)
    fmt.Println("Counter:", counter)
}
```

**运行**：

```bash
go run -race race_example.go
```

**输出**：

```
==================
WARNING: DATA RACE
Write at 0x00c000018090 by goroutine 7:
  main.main.func1()
      /path/to/race_example.go:13 +0x4c

Previous write at 0x00c000018090 by goroutine 8:
  main.main.func2()
      /path/to/race_example.go:19 +0x4c

Goroutine 7 (running) created at:
  main.main()
      /path/to/race_example.go:11 +0x7c
==================
Counter: 1523
```

**解释**：

- 检测到数据竞争（DATA RACE）
- 指出了冲突的代码位置
- 结果不是预期的 2000

---

### 3.3 修复竞态条件

```go
// 修复：使用 Mutex
func main() {
    counter := 0
    var mu sync.Mutex

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    }()

    go func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    }()

    wg.Wait()
    fmt.Println("Counter:", counter)  // 总是 2000
}
```

**再次运行**：

```bash
go run -race main.go
# 没有 DATA RACE 警告！
Counter: 2000
```

---

## 📖 四、并发模式回顾与扩展

### 4.1 Worker 池模式（已学）

```go
// 你已经掌握的模式
jobs := make(chan Task, 10)
results := make(chan Result, 20)

for i := 0; i < numWorkers; i++ {
    go worker(jobs, results)
}
```

---

### 4.2 Pipeline 管道模式（已学）

```go
// 数据流水线
nums := generate(1, 2, 3, 4, 5)
squared := square(nums)
filtered := filter(squared)
```

---

### 4.3 扇出扇入模式（Fan-out/Fan-in）

```go
// 扇出：一个输入，多个 worker 处理
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)

    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output

        go func() {
            for val := range input {
                output <- val * 2  // 处理
            }
            close(output)
        }()
    }

    return outputs
}

// 扇入：多个输入，合并为一个输出
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup

    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for val := range ch {
                output <- val
            }
        }(input)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}
```

---

### 4.4 超时模式（Timeout Pattern）

```go
func doWorkWithTimeout(ctx context.Context) error {
    result := make(chan error, 1)

    go func() {
        // 执行耗时任务
        time.Sleep(5 * time.Second)
        result <- nil
    }()

    select {
    case err := <-result:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## 💻 代码示例

完整示例代码请查看 `example.go`

---

## 🎯 随堂练习

### 练习要求：实现一个"并发网络爬虫"

#### **功能需求**：

1. **带 Context 的 Worker 池**：
   - 启动 5 个 worker
   - 使用 Context 实现超时控制（10秒）
   - 支持手动取消

2. **并发安全的 URL 去重**：
   - 使用 Mutex 或 Channel 实现
   - 记录已访问的 URL

3. **竞态检测**：
   - 代码必须通过 `go run -race` 检测

4. **统计信息**：
   - 已爬取 URL 数量
   - 成功/失败次数
   - 使用 atomic 或 Mutex 保证线程安全

---

### 期望输出示例

```
===== 启动爬虫系统 =====
Worker-1: 开始爬取 https://example.com
Worker-2: 开始爬取 https://golang.org
Worker-1: 完成爬取 (200ms)
Worker-3: 开始爬取 https://github.com
...

===== 10秒后超时 =====
Context 取消，停止所有任务

===== 统计信息 =====
总URL数: 20
已爬取: 15
跳过（重复）: 3
失败: 2
```

---

## 🔑 关键知识点总结

| 概念              | 作用           | vs Java         | 使用场景          |
| ----------------- | -------------- | --------------- | ----------------- |
| **Context**       | 取消信号、超时 | ThreadLocal     | HTTP 请求、长任务 |
| **Mutex**         | 互斥锁         | synchronized    | 保护共享变量      |
| **RWMutex**       | 读写锁         | ReadWriteLock   | 读多写少          |
| **Channel**       | 通信代替锁     | BlockingQueue   | Go 首选方案       |
| **Atomic**        | 原子操作       | AtomicInteger   | 简单计数器        |
| **Race Detector** | 竞态检测       | ThreadSanitizer | 开发调试          |

**Go 并发安全原则**：

1. ✅ **优先使用 Channel**：通信代替共享内存
2. ✅ **不得已用 Mutex**：保护必须共享的变量
3. ✅ **简单场景用 Atomic**：性能最好
4. ✅ **始终用 -race 测试**：发现隐藏的 bug

**Context 使用原则**：

1. ✅ **第一个参数**：func DoSomething(ctx context.Context, ...)
2. ✅ **总是 defer cancel()**：防止泄漏
3. ✅ **传递给所有子任务**：级联取消
4. ❌ **不要存储在结构体**：通过参数传递

**下一节预告**：包管理与模块化（Go Modules、包可见性、init函数）！
