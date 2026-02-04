# 第05节：并发编程基础 - Goroutine 与 Channel

> **本节目标**：掌握 Go 的并发模型，理解 Goroutine 和 Channel 的使用  
> **前置知识**：第01-04节的内容（不会用到 context、sync 包的高级特性）  
> **重要程度**：⭐⭐⭐⭐⭐ Go 的核心竞争力！

---

## 📖 一、并发模型对比

### 1.1 Java 的并发模型

#### **Java 的线程模型（基于共享内存）**

```java
// Java 创建线程
Thread thread = new Thread(() -> {
    System.out.println("Hello from thread");
});
thread.start();

// 使用线程池
ExecutorService executor = Executors.newFixedThreadPool(10);
executor.submit(() -> {
    // 执行任务
});

// 共享内存 + 锁
class Counter {
    private int count = 0;

    public synchronized void increment() {
        count++;  // 需要加锁保护
    }
}
```

**Java 并发特点**：

- ✅ 基于操作系统线程（重量级）
- ✅ 共享内存通信（需要锁）
- ⚠️ 线程开销大（每个线程 1MB+ 栈空间）
- ⚠️ 上下文切换代价高
- ⚠️ 容易出现死锁、竞态条件

---

### 1.2 Go 的并发模型

#### **Go 的 Goroutine 模型（CSP - Communicating Sequential Processes）**

```go
// 创建 Goroutine（只需要 go 关键字！）
go func() {
    fmt.Println("Hello from goroutine")
}()

// 使用 Channel 通信
ch := make(chan int)

go func() {
    ch <- 42  // 发送数据
}()

value := <-ch  // 接收数据
```

**Go 并发特点**：

- ✅ **轻量级**：每个 Goroutine 只需 2KB 栈空间（可增长）
- ✅ **高并发**：可以轻松创建数十万个 Goroutine
- ✅ **通信代替共享**：使用 Channel 通信，避免锁
- ✅ **调度器优化**：Go runtime 自动调度（M:N 模型）
- ✅ **简单易用**：只需 `go` 关键字

---

### 1.3 核心理念对比

| 特性         | Java         | Go               |
| ------------ | ------------ | ---------------- |
| **并发单元** | Thread       | Goroutine        |
| **创建方式** | new Thread() | go func()        |
| **开销**     | 重（1MB+栈） | 轻（2KB栈）      |
| **数量限制** | 数千         | 数十万+          |
| **通信方式** | 共享内存+锁  | Channel          |
| **调度**     | OS 调度      | Go runtime 调度  |
| **理念**     | 共享内存     | **通信代替共享** |

**Go 的并发哲学**：

> **"Do not communicate by sharing memory; instead, share memory by communicating."**  
> **不要通过共享内存来通信，而应该通过通信来共享内存。**

---

## 📖 二、Goroutine 基础

### 2.1 创建 Goroutine

#### **基本语法**

```go
// 方式1：使用匿名函数
go func() {
    fmt.Println("Hello from goroutine")
}()

// 方式2：调用已有函数
func sayHello() {
    fmt.Println("Hello")
}

func main() {
    go sayHello()  // 启动 goroutine

    // 注意：主 goroutine 结束，所有子 goroutine 也会终止
    time.Sleep(time.Second)  // 等待 goroutine 执行
}
```

**Java 对比**：

```java
// Java 需要更多代码
new Thread(() -> {
    System.out.println("Hello from thread");
}).start();

// 或使用线程池
ExecutorService executor = Executors.newCachedThreadPool();
executor.submit(() -> {
    System.out.println("Hello");
});
```

---

### 2.2 Goroutine 的参数传递

```go
// 正确方式：传递参数
for i := 0; i < 5; i++ {
    go func(n int) {  // 参数传递
        fmt.Println(n)
    }(i)  // 传入 i
}

// ❌ 错误方式：闭包捕获（常见陷阱！）
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // 可能都打印 5！
    }()
}
```

**为什么会打印 5？**

- Goroutine 启动需要时间
- 当 Goroutine 执行时，循环可能已经结束
- 所有 Goroutine 共享同一个 `i` 变量

---

### 2.3 主 Goroutine 等待

**问题**：主 Goroutine 结束，所有子 Goroutine 都会被终止

```go
func main() {
    go func() {
        fmt.Println("This might not print")
    }()

    // main 立即结束，goroutine 可能还没执行
}
```

**解决方案**：

#### **方式1：使用 time.Sleep（不推荐）**

```go
func main() {
    go func() {
        fmt.Println("Hello")
    }()

    time.Sleep(time.Second)  // 简单但不可靠
}
```

#### **方式2：使用 Channel 同步（推荐）**

```go
func main() {
    done := make(chan bool)

    go func() {
        fmt.Println("Hello")
        done <- true  // 通知完成
    }()

    <-done  // 等待通知
}
```

#### **方式3：使用 WaitGroup（后面会讲）**

---

## 📖 三、Channel 基础

### 3.1 Channel 的概念

**Channel 是什么？**

- Go 的**并发通信机制**
- 类似于**有类型的管道**
- 用于在 Goroutine 之间传递数据

**Java 对比**：

```java
// Java 使用 BlockingQueue
BlockingQueue<Integer> queue = new LinkedBlockingQueue<>();

// 生产者
new Thread(() -> {
    queue.put(42);
}).start();

// 消费者
int value = queue.take();
```

**Go 的 Channel 更简洁**：

```go
ch := make(chan int)  // 创建 channel

// 生产者
go func() {
    ch <- 42  // 发送
}()

// 消费者
value := <-ch  // 接收
```

---

### 3.2 Channel 的创建和使用

#### **创建 Channel**

```go
// 无缓冲 channel（同步）
ch1 := make(chan int)

// 有缓冲 channel（异步）
ch2 := make(chan int, 10)  // 缓冲区大小为 10

// 只读 channel
var readOnly <-chan int = ch1

// 只写 channel
var writeOnly chan<- int = ch1
```

---

#### **发送和接收**

```go
ch := make(chan int)

// 发送数据
ch <- 42

// 接收数据
value := <-ch

// 接收并忽略
<-ch

// 接收并检查 channel 是否关闭
value, ok := <-ch
if !ok {
    fmt.Println("Channel 已关闭")
}
```

---

### 3.3 无缓冲 Channel vs 有缓冲 Channel

#### **无缓冲 Channel（同步通信）**

```go
ch := make(chan int)  // 无缓冲

// 发送操作会阻塞，直到有接收者
go func() {
    ch <- 42  // 阻塞，等待接收
    fmt.Println("发送完成")
}()

value := <-ch  // 接收，发送者被唤醒
fmt.Println(value)
```

**特点**：

- ✅ **同步通信**：发送和接收必须同时准备好
- ✅ **零容量**：不存储数据
- ✅ **强同步保证**

---

#### **有缓冲 Channel（异步通信）**

```go
ch := make(chan int, 3)  // 缓冲区大小为 3

// 可以连续发送 3 次，不会阻塞
ch <- 1
ch <- 2
ch <- 3

// 第 4 次发送会阻塞，直到有人接收
// ch <- 4  // 如果没有接收者，会阻塞

// 接收
fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
fmt.Println(<-ch)  // 3
```

**特点**：

- ✅ **异步通信**：缓冲区未满时，发送不阻塞
- ✅ **有容量**：可以存储数据
- ⚠️ **缓冲区满时阻塞**

---

### 3.4 关闭 Channel

```go
ch := make(chan int, 3)

// 发送数据
ch <- 1
ch <- 2
ch <- 3

// 关闭 channel
close(ch)

// 可以继续接收，直到 channel 为空
fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
fmt.Println(<-ch)  // 3

// channel 为空后，接收会得到零值
fmt.Println(<-ch)  // 0（int 的零值）

// 检查 channel 是否关闭
value, ok := <-ch
if !ok {
    fmt.Println("Channel 已关闭且为空")
}
```

**重要规则**：

- ✅ **发送者关闭 channel**
- ❌ **接收者不应该关闭 channel**
- ❌ **不要向已关闭的 channel 发送数据**（会 panic）
- ✅ **可以从已关闭的 channel 接收数据**

---

### 3.5 遍历 Channel

#### **使用 range 遍历**

```go
ch := make(chan int, 5)

// 发送数据
go func() {
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch)  // 必须关闭，否则 range 会永久阻塞
}()

// 遍历接收
for value := range ch {//阻塞式写法，只要ch没关闭，range 就会一直等待新数据
    fmt.Println(value)
}

fmt.Println("所有数据接收完毕")
```

**注意**：

- `range` 会持续接收，直到 channel 关闭
- 如果不关闭 channel，`range` 会永久阻塞

---

## 📖 四、Select 多路复用

### 4.1 Select 的概念

**Select** 是 Go 的多路复用机制，类似于网络编程中的 `select/epoll`

**Java 对比**：
Java 没有直接等价的语法，需要使用复杂的代码：

```java
// Java 需要手动轮询或使用复杂的异步框架
CompletableFuture<Integer> future1 = ...;
CompletableFuture<Integer> future2 = ...;
CompletableFuture.anyOf(future1, future2).thenAccept(...);
```

**Go 的 Select**：

```go
select {
case value := <-ch1:
    fmt.Println("从 ch1 接收:", value)
case value := <-ch2:
    fmt.Println("从 ch2 接收:", value)
case ch3 <- 42:
    fmt.Println("向 ch3 发送")
default:
    fmt.Println("所有 channel 都未就绪")
}
```

---

### 4.2 Select 的基本使用

#### **示例1：等待多个 Channel**

```go
ch1 := make(chan string)
ch2 := make(chan string)

go func() {
    time.Sleep(1 * time.Second)
    ch1 <- "来自 ch1"
}()

go func() {
    time.Sleep(2 * time.Second)
    ch2 <- "来自 ch2"
}()

// select 会阻塞，直到某个 case 就绪
select {
case msg1 := <-ch1:
    fmt.Println(msg1)
case msg2 := <-ch2:
    fmt.Println(msg2)
}
```

---

#### **示例2：使用 default 实现非阻塞**

```go
ch := make(chan int)

select {
case value := <-ch:
    fmt.Println("接收到:", value)
default:
    fmt.Println("channel 未就绪，继续其他工作")
}
```

---

#### **示例3：超时控制**

```go
ch := make(chan int)

go func() {
    time.Sleep(2 * time.Second)
    ch <- 42
}()

select {
case value := <-ch:
    fmt.Println("接收到:", value)
case <-time.After(1 * time.Second):
    fmt.Println("超时！")
}
```

---

### 4.3 Select 的特性

1. **随机选择**：如果多个 case 同时就绪，随机选择一个，防止饥饿，避免总是偏向执行第一个
2. **阻塞等待**：如果没有 case 就绪且无 default，当前的 Goroutine 会被**阻塞**，直到至少有一个 `case` 变为可用。
3. **非阻塞**：如果有 `default` 分支，当所有 `case` 都未准备好时，会直接执行 `default`，而不会阻塞。
4. **可以发送也可以接收**：case 可以是发送或接收操作

```go
// 随机选择示例
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

ch1 <- 1
ch2 <- 2

// 两个 case 都就绪，随机选择
select {
case <-ch1:
    fmt.Println("选择了 ch1")
case <-ch2:
    fmt.Println("选择了 ch2")
}
```

---

## 📖 五、WaitGroup 等待组

### 5.1 WaitGroup 的概念

**问题**：如何等待多个 Goroutine 完成？

**Java 对比**：

```java
// Java 使用 CountDownLatch
CountDownLatch latch = new CountDownLatch(3);

for (int i = 0; i < 3; i++) {
    new Thread(() -> {
        // 执行任务
        latch.countDown();
    }).start();
}

latch.await();  // 等待所有线程完成
```

**Go 的 WaitGroup**：

```go
import "sync"

var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)  // 计数器 +1

    go func(n int) {
        defer wg.Done()  // 计数器 -1
        fmt.Println("Goroutine", n)
    }(i)
}

wg.Wait()  // 等待计数器归零
fmt.Println("所有 goroutine 完成")
```

---

### 5.2 WaitGroup 的使用

```go
import "sync"

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // 函数结束时调用

    fmt.Printf("Worker %d 开始工作\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d 完成工作\n", id)
}

func main() {
    var wg sync.WaitGroup

    // 启动 5 个 worker
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }

    wg.Wait()  // 等待所有 worker 完成
    fmt.Println("所有任务完成")
}
```

**注意事项**：

- ✅ `Add()` 应该在 goroutine 启动**之前**调用
- ✅ `Done()` 应该在 goroutine **结束时**调用（通常用 defer）
- ✅ `Wait()` 会阻塞，直到计数器归零
- ⚠️ 传递 `*sync.WaitGroup` 指针，不要复制

---

## 📖 六、常见并发模式

### 6.1 生产者-消费者模式

```go
func producer(ch chan<- int, count int) {
    for i := 1; i <= count; i++ {
        ch <- i
        fmt.Printf("生产: %d\n", i)
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

func consumer(ch <-chan int, id int) {
    for value := range ch {
        fmt.Printf("消费者 %d 消费: %d\n", id, value)
        time.Sleep(200 * time.Millisecond)
    }
}

func main() {
    ch := make(chan int, 5)

    go producer(ch, 10)

    go consumer(ch, 1)
    go consumer(ch, 2)

    time.Sleep(3 * time.Second)
}
```

---

### 6.2 扇出扇入模式（Fan-out Fan-in）

```go
// 扇出：一个输入，多个 worker 处理
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)

    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output

        go func(out chan<- int) {
            for value := range input {
                out <- value * 2  // 处理数据
            }
            close(out)
        }(output)
    }

    return outputs
}

// 扇入：多个输入合并为一个输出
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup

    for _, input := range inputs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for value := range ch {
                output <- value
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

### 6.3 管道模式（Pipeline）

```go
// 阶段1：生成数据
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

// 阶段2：平方
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

// 阶段3：过滤偶数
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

func main() {
    // 构建管道
    nums := generate(1, 2, 3, 4, 5)
    squared := square(nums)
    filtered := filterEven(squared)

    // 消费结果
    for result := range filtered {
        fmt.Println(result)
    }
}
```

---

## 💻 代码示例

完整示例代码请查看 `example.go`

---

## 🎯 随堂练习

### 练习要求：实现一个"并发任务处理系统"

#### **功能需求**：

1. **任务生成器**：
   - 生成 20 个任务（任务 ID: 1-20）
   - 发送到 channel

2. **Worker 池**：
   - 创建 3 个 worker goroutine
   - 从 channel 接收任务并处理
   - 模拟处理时间（随机 100-500ms）

3. **结果收集器**：
   - 收集所有处理结果
   - 统计完成数量

4. **超时控制**：
   - 使用 select 实现 5 秒超时
   - 超时后停止接收新任务

5. **优雅退出**：
   - 使用 WaitGroup 等待所有 worker 完成
   - 打印统计信息

---

### 期望输出示例

```
===== 启动任务处理系统 =====
生成器: 发送任务 1
生成器: 发送任务 2
Worker-1: 开始处理任务 1
Worker-2: 开始处理任务 2
生成器: 发送任务 3
Worker-3: 开始处理任务 3
Worker-1: 完成任务 1 (耗时 234ms)
Worker-1: 开始处理任务 4
...

===== 统计信息 =====
总任务数: 20
已完成: 20
总耗时: 3.45s
平均处理时间: 287ms
```

---

## 🔑 关键知识点总结

| 概念               | 特点          | vs Java             | 使用场景        |
| ------------------ | ------------- | ------------------- | --------------- |
| **Goroutine**      | 轻量级（2KB） | Thread（1MB+）      | 高并发任务      |
| **Channel**        | 通信机制      | BlockingQueue       | 数据传递        |
| **Select**         | 多路复用      | 无直接对应          | 多 channel 选择 |
| **WaitGroup**      | 等待组        | CountDownLatch      | 等待多任务      |
| **无缓冲 Channel** | 同步通信      | SynchronousQueue    | 强同步          |
| **有缓冲 Channel** | 异步通信      | LinkedBlockingQueue | 削峰填谷        |

**Go 并发哲学**：

- ✅ **轻量级并发**：Goroutine 开销极小
- ✅ **通信代替共享**：Channel 避免锁
- ✅ **简单易用**：只需 `go` 关键字
- ✅ **强大的抽象**：Select、Pipeline 等模式

**下一节预告**：错误处理深度剖析、包管理、测试！
