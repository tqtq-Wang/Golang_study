# 错误处理极简入门（第02节专用）

> **说明**：这里只讲最基础的错误处理，让你能完成本节练习。第08节会深入讲解完整的错误处理机制。

---

## 1️⃣ error 是什么？

**简单理解**：`error` 就是一个特殊的返回值，用来告诉调用者"出错了"。

### Java vs Go

```java
// Java：用异常表示错误
public int divide(int a, int b) {
    if (b == 0) {
        throw new ArithmeticException("除数不能为0");  // 抛出异常
    }
    return a / b;
}

// 调用时必须捕获
try {
    int result = divide(10, 0);
} catch (Exception e) {
    System.out.println("出错了: " + e.getMessage());
}
```

```go
// Go：用返回值表示错误（不抛异常）
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("除数不能为0")  // 返回错误
    }
    return a / b, nil  // nil 表示没错误
}

// 调用时检查返回值
result, err := divide(10, 0)
if err != nil {  // 如果有错误
    fmt.Println("出错了:", err)
} else {  // 如果没错误
    fmt.Println("结果:", result)
}
```

---

## 2️⃣ 如何创建错误？（两种方式）

### 方式1：errors.New（固定消息）

```go
import "errors"

func checkAge(age int) error {
    if age < 18 {
        return errors.New("年龄必须大于18岁")  // 创建错误
    }
    return nil  // 没错误，返回 nil
}
```

### 方式2：fmt.Errorf（格式化消息）

```go
import "fmt"

func checkAge(age int) error {
    if age < 18 {
        // 可以插入变量
        return fmt.Errorf("年龄 %d 不符合要求，必须 >= 18", age)
    }
    return nil
}
```

**推荐**：用 `fmt.Errorf`，可以包含具体的错误信息。

---

## 3️⃣ 如何检查错误？（固定套路）

### 标准模式（90%的情况都这样写）

```go
result, err := someFunction()
if err != nil {
    // 有错误，处理错误
    fmt.Println("Error:", err)
    return  // 或者 return err 传递给上层
}
// 没错误，继续使用 result
fmt.Println("成功:", result)
```

### 完整示例

```go
func main() {
    // 调用可能出错的函数
    score, err := getScore("张三")

    // 检查错误（必须检查！）
    if err != nil {
        fmt.Println("获取成绩失败:", err)
        return  // 直接退出
    }

    // 没错误，继续执行
    fmt.Println("张三的成绩:", score)
}

func getScore(name string) (float64, error) {
    if name == "" {
        return 0, fmt.Errorf("姓名不能为空")
    }
    // 假设从数据库查询...
    return 85.5, nil  // 返回成绩和 nil（表示没错误）
}
```

---

## 4️⃣ 在方法中返回错误

### 指针接收者 + 错误返回

```go
type Student struct {
    Name  string
    Score float64
}

// 设置分数（可能失败，返回 error）
func (s *Student) SetScore(score float64) error {
    if score < 0 || score > 100 {
        // 返回错误
        return fmt.Errorf("分数 %.1f 无效，必须在 0-100 之间", score)
    }
    // 没错误，修改对象
    s.Score = score
    return nil  // 返回 nil 表示成功
}

// 使用
func main() {
    student := Student{Name: "李四", Score: 0}

    // 尝试设置有效分数
    err := student.SetScore(95)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("设置成功，分数:", student.Score)
    }

    // 尝试设置无效分数
    err = student.SetScore(150)
    if err != nil {
        fmt.Println("Error:", err)  // 会输出错误
    }
}
```

---

## 5️⃣ 常见错误处理模式

### 模式1：立即返回（最常见）

```go
func process() error {
    err := step1()
    if err != nil {
        return err  // 直接返回错误
    }

    err = step2()
    if err != nil {
        return err
    }

    return nil  // 全部成功
}
```

### 模式2：打印并继续

```go
func main() {
    err := doSomething()
    if err != nil {
        fmt.Println("Warning:", err)  // 只是警告
        // 继续执行其他逻辑...
    }
}
```

### 模式3：尝试多个操作，记录所有错误

```go
func validateStudent(s Student) []error {
    var errors []error

    if s.Name == "" {
        errors = append(errors, fmt.Errorf("姓名不能为空"))
    }

    if s.Score < 0 || s.Score > 100 {
        errors = append(errors, fmt.Errorf("分数无效"))
    }

    return errors  // 返回所有错误
}
```

---

## 6️⃣ 本节练习只需要知道这些

### 你需要会的：

1. ✅ 用 `fmt.Errorf("消息")` 创建错误
2. ✅ 函数返回值加 `error`：`func xxx() (结果, error)`
3. ✅ 检查错误：`if err != nil { ... }`
4. ✅ 成功时返回 `nil`：`return result, nil`
5. ✅ 失败时返回错误：`return 零值, fmt.Errorf("...")`

### 你暂时不需要知道的（第08节再学）：

- ❌ 错误包装（`%w`）
- ❌ errors.Is / errors.As
- ❌ 自定义错误类型
- ❌ panic / recover（这是特殊机制，不是常规错误处理）

---

## 7️⃣ 快速参考：返回值零值

当返回错误时，第一个返回值需要返回"零值"：

```go
func getAge() (int, error) {
    return 0, fmt.Errorf("错误")  // int 的零值是 0
}

func getName() (string, error) {
    return "", fmt.Errorf("错误")  // string 的零值是 ""
}

func getScore() (float64, error) {
    return 0.0, fmt.Errorf("错误")  // float64 的零值是 0.0
}

func getStudent() (*Student, error) {
    return nil, fmt.Errorf("错误")  // 指针的零值是 nil
}
```

---

## 8️⃣ 完整示例：学生管理

```go
package main

import "fmt"

type Student struct {
    Name  string
    Score float64
}

// 设置分数（带错误检查）
func (s *Student) SetScore(score float64) error {
    if score < 0 || score > 100 {
        return fmt.Errorf("分数必须在 0-100 之间")
    }
    s.Score = score
    return nil
}

// 判断是否及格（不会出错，不需要返回 error）
func (s Student) IsPassed() bool {
    return s.Score >= 60
}

func main() {
    student := Student{Name: "张三"}

    // 设置有效分数
    err := student.SetScore(85)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("%s 的分数: %.1f\n", student.Name, student.Score)

    // 设置无效分数
    err = student.SetScore(150)
    if err != nil {
        fmt.Println("Error:", err)  // 会输出：分数必须在 0-100 之间
    }

    // 判断是否及格
    if student.IsPassed() {
        fmt.Println("及格了！")
    }
}
```

---

## ✅ 总结：错误处理三步走

```go
// 1. 函数声明：返回值加 error
func doSomething(input int) (result string, err error)

// 2. 函数内部：出错时返回错误
if input < 0 {
    return "", fmt.Errorf("输入不能为负数")
}
return "success", nil

// 3. 调用时：检查错误
result, err := doSomething(-1)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Success:", result)
```

---

## 🎯 现在可以做练习了！

只需要用：

- `fmt.Errorf("消息")` 创建错误
- `if err != nil` 检查错误
- 成功返回 `nil`，失败返回错误

第08节会深入讲解更高级的错误处理技巧！
