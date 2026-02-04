package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID int
}

type Result struct {
	TaskID   int
	WorkerID int
	Duration time.Duration
}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) { //只能读 tasks， 只能写 results
	defer wg.Done() // 确保在函数结束时调用 Done

	for task := range tasks {
		fmt.Printf("Worker-%d: 开始处理任务 %d\n", id, task.ID)

		// 模拟处理时间
		start := time.Now()
		processTime := time.Duration(rand.Intn(400)+100) * time.Millisecond
		time.Sleep(processTime)
		duration := time.Since(start)

		results <- Result{
			TaskID:   task.ID,
			WorkerID: id,
			Duration: duration,
		}
		fmt.Printf("Worker-%d: 完成任务 %d，耗时 %v\n", id, task.ID, duration)
	}
}

func collector(results <-chan Result, done chan<- struct{}) {
	count := 0
	var totalDuration time.Duration
	for result := range results {
		fmt.Printf("Collector: 任务 %d 由 Worker-%d 完成，耗时 %v\n", result.TaskID, result.WorkerID, result.Duration)
		count++
		totalDuration += result.Duration
	}

	// 打印最终统计
	fmt.Println("\n===== 统计信息 =====")
	fmt.Println("总完成任务数:", count)
	if count > 0 {
		fmt.Printf("平均处理时间: %v\n", totalDuration/time.Duration(count))
	}
	done <- struct{}{}
}

func main() {
	// 设置随机种子，确保每次运行的随机数不同
	rand.Seed(time.Now().UnixNano())

	tasks := make(chan Task, 10)
	results := make(chan Result, 20)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// 启动 worker
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// 启动 collector
	go collector(results, done)

	// 生成任务（5秒超时）
	fmt.Println("===== 启动任务处理系统 (限时 5 秒) =====")

	timer := time.NewTimer(5 * time.Second) // 设置 5 秒定时器

	//循环生产20个任务
	go func() {
		for i := 1; i <= 20; i++ {
			select {
			case <-timer.C:
				fmt.Println("任务生成超时，停止生成新任务。")
				close(tasks)
				return
			default:
				tasks <- Task{ID: i}
				fmt.Printf("主程序: 生成任务 %d\n", i)
				time.Sleep(200 * time.Millisecond) // 模拟任务生成间隔
			}
		}
		close(tasks) // 关闭任务通道，表示没有更多任务
	}()

	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(results) // 关闭结果通道，表示没有更多结果
	}()

	// 等待 collector 完成
	<-done
	fmt.Println("系统安全退出")
}
