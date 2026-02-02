# 第02节：函数与方法

> **本节目标**：掌握 Go 的函数特性（多返回值、defer），理解方法接收者，学会错误处理模式

---

## 📖 一、核心概念讲解

### 1.1 函数定义：多返回值特性

#### **Java 的函数定义（单返回值）**

```java
// Java 只能返回一个值
public int divide(int a, int b) {
    return a / b;
}

// 如果要返回多个值，必须用对象包装
public class DivResult {
    int quotient;
    int remainder;
}

public DivResult divideWithRemainder(int a, int b) {
    DivResult result = new DivResult();
    result.quotient = a / b;
    result.remainder = a % b;
    return result;
}
```

#### **Go 的函数定义（多返回值）**

```go
// Go 可以直接返回多个值（无需包装类）
func divide(a, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

// 使用方式
q, r := divide(10, 3)  // q=3, r=1

// 如果只需要部分返回值，用 _ 忽略
q, _ := divide(10, 3)  // 只要商，忽略余数
```

#### **关键差异**：

| 特性           | Go                     | Java                 |
| -------------- | ---------------------- | -------------------- |
| **返回值数量** | 支持多返回值           | 只能返回1个          |
| **错误处理**   | `return result, error` | 抛出异常或返回特殊值 |
| **命名返回值** | 支持（自动初始化零值） | 不支持               |
| **语法简洁性** | 极简                   | 需要包装类           |

---

### 1.2 命名返回值（Named Return Values）

**Go 特有特性**：可以给返回值命名，函数内直接使用，return 时自动返回

```go
// 命名返回值（相当于声明了变量）
func divide(a, b int) (quotient int, remainder int) {
    quotient = a / b    // 直接赋值，无需声明
    remainder = a % b
    return              // 裸返回（naked return），自动返回命名变量
}

// 等价于
func divide(a, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}
```

**优点**：

- 文档化：返回值名称即说明
- 简化代码：零值自动初始化
- 适合复杂函数：多处 return 时减少重复

**缺点**：

- 裸返回（`return` 不写变量）可能降低可读性
- 企业开发建议：**简短函数用裸返回，复杂函数显式返回**

---

### 1.3 参数传递：值传递 vs Java 的引用传递

#### **Java 的参数传递**

```java
public class PassByValue {
    // 基本类型：值传递
    public void modifyInt(int x) {
        x = 100;  // 不影响外部
    }

    // 引用类型：引用的值传递（可修改对象内容）
    public void modifyList(List<String> list) {
        list.add("new");  // 影响外部（修改对象）
        list = new ArrayList<>();  // 不影响外部（修改引用）
    }
}
```

#### **Go 的参数传递**

**重要**：Go 中**所有参数都是值传递**（包括指针、切片、map）

```go
// 基本类型：值传递（拷贝）
func modifyInt(x int) {
    x = 100  // 不影响外部
}

// 结构体：值传递（拷贝整个结构体）
type Person struct {
    Name string
    Age  int
}

func modifyPerson(p Person) {
    p.Age = 30  // 不影响外部（修改的是副本）
}

// 指针：传递指针的值（可以修改指向的数据）
func modifyPersonPtr(p *Person) {
    p.Age = 30  // 影响外部（通过指针修改原对象）
}

// 切片/Map：传递头部结构的值（可修改底层数据）
func modifySlice(s []int) {
    s[0] = 999  // 影响外部（修改底层数组）
    s = append(s, 100)  // 不一定影响外部（可能扩容重新分配）
}
```

#### **对比总结**：

| 类型            | Java         | Go                                       |
| --------------- | ------------ | ---------------------------------------- |
| **基本类型**    | 值传递       | 值传递                                   |
| **对象/结构体** | 引用的值传递 | 值传递（拷贝整个结构体）                 |
| **指针**        | 无显式指针   | 值传递（传递指针的副本，但指向同一地址） |
| **数组**        | 引用类型     | 值类型（拷贝整个数组）                   |
| **切片/Map**    | 类似引用     | 传递头部结构（可修改底层数据）           |

**Go 中高效传递大结构体的方式**：传递指针 `*Person` 而不是值 `Person`

---

### 1.4 方法接收者（Method Receiver）

#### **Java 的方法（绑定到类）**

```java
public class Person {
    private String name;
    private int age;

    // 实例方法（隐含 this）
    public void sayHello() {
        System.out.println("Hello, I'm " + this.name);
    }

    // 静态方法（无 this）
    public static void info() {
        System.out.println("Person class");
    }
}
```

#### **Go 的方法（绑定到类型）**

**Go 没有类**，但可以为任何类型（结构体、基本类型）定义方法

```go
type Person struct {
    Name string
    Age  int
}

// 值接收者（类似 Java 的值传递）
func (p Person) SayHello() {
    fmt.Printf("Hello, I'm %s\n", p.Name)
}

// 指针接收者（类似 Java 的实例方法）
func (p *Person) SetAge(age int) {
    p.Age = age  // 可以修改原对象
}

// 使用
p := Person{Name: "张三", Age: 25}
p.SayHello()   // 值接收者
p.SetAge(30)   // 指针接收者（Go 自动取地址）
```

#### **值接收者 vs 指针接收者**

| 接收者类型     | 语法          | 能否修改对象        | 拷贝开销             | 使用场景         |
| -------------- | ------------- | ------------------- | -------------------- | ---------------- |
| **值接收者**   | `(p Person)`  | ❌ 不能（操作副本） | 高（拷贝整个结构体） | 小对象、只读方法 |
| **指针接收者** | `(p *Person)` | ✅ 可以             | 低（只拷贝指针）     | 大对象、需要修改 |

**企业开发建议**：

1. **需要修改对象** → 必须用指针接收者
2. **大结构体（>几百字节）** → 用指针接收者（避免拷贝）
3. **实现接口时** → 保持一致性（全用值或全用指针）
4. **不确定时** → 默认用指针接收者

#### **Go 的自动取地址与解引用**

```go
p := Person{Name: "张三"}
p.SetAge(30)  // Go 自动转换为 (&p).SetAge(30)

ptr := &Person{Name: "李四"}
ptr.SayHello()  // Go 自动转换为 (*ptr).SayHello()
```

**Java 无此特性**，Go 的语法糖让指针使用更便捷。

---

### 1.5 defer 延迟执行

#### **Java 的资源管理（try-finally）**

```java
public void readFile() {
    FileReader reader = null;
    try {
        reader = new FileReader("file.txt");
        // 读取文件...
    } catch (IOException e) {
        e.printStackTrace();
    } finally {
        if (reader != null) {
            try {
                reader.close();  // 保证关闭
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}

// Java 7+ try-with-resources
try (FileReader reader = new FileReader("file.txt")) {
    // 自动关闭
}
```

#### **Go 的 defer 机制**

**defer**：延迟执行语句，函数返回前执行（类似 finally，但更灵活）

```go
func readFile() error {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer file.Close()  // 函数返回前自动执行

    // 读取文件...
    data := make([]byte, 100)
    file.Read(data)

    return nil  // file.Close() 会在这之前执行
}
```

#### **defer 的执行顺序（LIFO：后进先出）**

```go
func example() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    fmt.Println("start")
}
// 输出：
// start
// 3
// 2
// 1
```

**类比**：defer 类似一个栈，先 defer 的后执行（像洗碗一样，最上面的碗最后放，最先拿）

#### **defer 的典型使用场景**

```go
// 1. 资源清理
func process() {
    mu.Lock()
    defer mu.Unlock()  // 保证释放锁

    // 业务逻辑...
}

// 2. 追踪函数执行
func trace(name string) func() {
    fmt.Println("enter:", name)
    return func() {
        fmt.Println("exit:", name)
    }
}

func business() {
    defer trace("business")()  // 进入时打印，退出时打印
    // 业务逻辑...
}

// 3. 修改命名返回值
func calculate() (result int) {
    defer func() {
        result += 10  // 可以修改返回值
    }()
    result = 5
    return  // 实际返回 15
}
```

#### **defer 的陷阱（闭包变量）**

```go
// ❌ 错误示例
for i := 0; i < 3; i++ {
    defer fmt.Println(i)  // 输出：2 2 2（延迟求值）
}

// ✅ 正确示例
for i := 0; i < 3; i++ {
    i := i  // 创建新变量
    defer fmt.Println(i)  // 输出：2 1 0
}

// 或者
for i := 0; i < 3; i++ {
    defer func(n int) {
        fmt.Println(n)
    }(i)  // 立即求值
}
```

---

### 1.6 错误处理：error vs Java Exception

#### **Java 的异常机制**

```java
public int divide(int a, int b) throws ArithmeticException {
    if (b == 0) {
        throw new ArithmeticException("division by zero");
    }
    return a / b;
}

// 使用
try {
    int result = divide(10, 0);
} catch (ArithmeticException e) {
    System.out.println("Error: " + e.getMessage());
}
```

#### **Go 的错误处理**

**Go 没有异常机制**，用返回值传递错误（显式错误处理）

```go
// error 是内置接口
type error interface {
    Error() string
}

// 函数返回 error
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")  // 返回错误
    }
    return a / b, nil  // nil 表示无错误
}

// 使用（必须检查错误）
result, err := divide(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

#### **错误处理模式对比**

| 特性         | Go (error)           | Java (Exception)        |
| ------------ | -------------------- | ----------------------- |
| **错误类型** | 值（返回值）         | 对象（抛出）            |
| **控制流**   | 显式检查             | 隐式跳转（try-catch）   |
| **性能**     | 高（无栈展开）       | 低（栈展开开销）        |
| **强制处理** | 否（编译器不强制）   | 是（checked exception） |
| **代码冗余** | 较多 `if err != nil` | 较少，集中处理          |

#### **Go 的错误处理惯例**

```go
// 1. 立即检查错误
f, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("open file: %w", err)  // %w 包装错误
}
defer f.Close()

// 2. 错误包装（error wrapping）
if err := doSomething(); err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}

// 3. 错误链追踪
err := errors.New("original error")
err = fmt.Errorf("layer 2: %w", err)
err = fmt.Errorf("layer 3: %w", err)

// 检查错误类型
if errors.Is(err, os.ErrNotExist) {
    // 处理文件不存在
}
```

#### **自定义错误类型**

```go
// Java
public class CustomException extends Exception {
    private int errorCode;

    public CustomException(String message, int code) {
        super(message);
        this.errorCode = code;
    }
}

// Go
type CustomError struct {
    Code    int
    Message string
}

func (e *CustomError) Error() string {
    return fmt.Sprintf("code %d: %s", e.Code, e.Message)
}

// 使用
func validate() error {
    return &CustomError{Code: 400, Message: "invalid input"}
}
```

---

### 1.7 panic 和 recover（类似 Java 的异常）

**Go 也有类似异常的机制**，但只用于**不可恢复的错误**

```go
// panic：类似 throw
func mustOpen(filename string) *os.File {
    f, err := os.Open(filename)
    if err != nil {
        panic(err)  // 程序崩溃（除非 recover）
    }
    return f
}

// recover：类似 catch（只能在 defer 中使用）
func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    panic("something went wrong")  // 会被 recover 捕获
    fmt.Println("This won't print")
}
```

#### **panic/recover vs Java try/catch**

| 特性         | Go panic/recover   | Java try/catch |
| ------------ | ------------------ | -------------- |
| **使用场景** | 严重错误（程序级） | 常规错误处理   |
| **推荐度**   | 不推荐（用 error） | 推荐           |
| **性能**     | 类似异常（栈展开） | 栈展开         |
| **习惯**     | 只用于不可恢复错误 | 常规控制流     |

**Go 哲学**：

- 普通错误 → 用 `error` 返回值
- 程序bug → 用 `panic`（如数组越界、空指针）
- 库代码 → 永远不要 `panic`，交给调用者处理

---

## 💻 二、代码示例

### 示例1：函数与多返回值

```go
package main

import (
	"fmt"
	"math"
)

// ========== 多返回值 ==========
// 计算圆的面积和周长
func calculateCircle(radius float64) (area float64, circumference float64) {
	area = math.Pi * radius * radius
	circumference = 2 * math.Pi * radius
	return  // 裸返回（返回命名变量）
}

// 除法运算（返回商、余数、错误）
func divideWithError(a, b int) (int, int, error) {
	if b == 0 {
		return 0, 0, fmt.Errorf("division by zero")
	}
	return a / b, a % b, nil
}

// ========== 可变参数（类似 Java 的 varargs）==========
func sum(numbers ...int) int {  // ...int 表示可变参数
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

func main() {
	// 使用多返回值
	area, circumference := calculateCircle(5.0)
	fmt.Printf("半径5的圆 - 面积: %.2f, 周长: %.2f\n", area, circumference)

	// 忽略部分返回值
	area, _ = calculateCircle(10.0)
	fmt.Printf("半径10的圆 - 面积: %.2f\n", area)

	// 错误处理
	quotient, remainder, err := divideWithError(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 ÷ 3 = %d 余 %d\n", quotient, remainder)
	}

	// 除以0的情况
	_, _, err = divideWithError(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 可变参数
	fmt.Println("sum(1,2,3) =", sum(1, 2, 3))
	fmt.Println("sum(1,2,3,4,5) =", sum(1, 2, 3, 4, 5))

	// 传递切片给可变参数
	nums := []int{10, 20, 30}
	fmt.Println("sum(nums...) =", sum(nums...))  // ... 展开切片
}
```

### 示例2：方法接收者

```go
package main

import "fmt"

// ========== 定义结构体 ==========
type BankAccount struct {
	Owner   string
	Balance float64
}

// ========== 值接收者（只读方法）==========
func (acc BankAccount) ShowBalance() {
	fmt.Printf("账户: %s, 余额: %.2f 元\n", acc.Owner, acc.Balance)
}

// 尝试修改（无效）
func (acc BankAccount) TryDeposit(amount float64) {
	acc.Balance += amount  // 修改的是副本，不影响原对象
	fmt.Printf("（值接收者）存入后余额: %.2f\n", acc.Balance)
}

// ========== 指针接收者（可修改方法）==========
func (acc *BankAccount) Deposit(amount float64) {
	acc.Balance += amount
	fmt.Printf("存入 %.2f 元，当前余额: %.2f 元\n", amount, acc.Balance)
}

func (acc *BankAccount) Withdraw(amount float64) error {
	if amount > acc.Balance {
		return fmt.Errorf("余额不足：需要 %.2f，只有 %.2f", amount, acc.Balance)
	}
	acc.Balance -= amount
	fmt.Printf("取出 %.2f 元，当前余额: %.2f 元\n", amount, acc.Balance)
	return nil
}

func main() {
	// 创建账户
	acc := BankAccount{Owner: "张三", Balance: 1000.0}

	// 值接收者（只读）
	acc.ShowBalance()

	// 尝试修改（无效）
	fmt.Println("\n测试值接收者:")
	acc.TryDeposit(500)
	acc.ShowBalance()  // 余额不变，仍然是 1000

	// 指针接收者（可修改）
	fmt.Println("\n测试指针接收者:")
	acc.Deposit(500)
	acc.ShowBalance()  // 余额变为 1500

	// 取款
	err := acc.Withdraw(800)
	if err != nil {
		fmt.Println("Error:", err)
	}
	acc.ShowBalance()  // 余额变为 700

	// 取款失败
	err = acc.Withdraw(1000)
	if err != nil {
		fmt.Println("Error:", err)
	}
	acc.ShowBalance()  // 余额不变，仍然是 700
}
```

### 示例3：defer 延迟执行

```go
package main

import "fmt"

// ========== defer 基础 ==========
func deferBasic() {
	fmt.Println("开始")
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")
	fmt.Println("结束")
	// 输出顺序：开始 → 结束 → defer 3 → defer 2 → defer 1
}

// ========== defer 修改返回值 ==========
func deferModifyReturn() (result int) {
	defer func() {
		result += 10  // 修改命名返回值
	}()
	return 5  // 实际返回 15
}

// ========== defer 捕获 panic ==========
func deferRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到 panic:", r)
		}
	}()

	fmt.Println("正常执行")
	panic("发生错误！")  // 触发 panic
	fmt.Println("这行不会执行")
}

// ========== defer 在循环中的陷阱 ==========
func deferInLoop() {
	fmt.Println("错误示例:")
	for i := 0; i < 3; i++ {
		defer fmt.Println("defer:", i)  // 全部延迟到函数结束
	}
	// 输出：defer: 2, defer: 1, defer: 0
}

func deferInLoopFixed() {
	fmt.Println("\n正确示例:")
	for i := 0; i < 3; i++ {
		func(n int) {
			defer fmt.Println("defer:", n)
		}(i)  // 立即执行匿名函数
	}
	// 输出：defer: 0, defer: 1, defer: 2
}

func main() {
	// 1. 基础示例
	fmt.Println("=== defer 基础 ===")
	deferBasic()

	// 2. 修改返回值
	fmt.Println("\n=== defer 修改返回值 ===")
	result := deferModifyReturn()
	fmt.Println("返回值:", result)

	// 3. 捕获 panic
	fmt.Println("\n=== defer 捕获 panic ===")
	deferRecover()
	fmt.Println("程序继续执行")

	// 4. 循环中的陷阱
	fmt.Println("\n=== defer 在循环中 ===")
	deferInLoop()
	deferInLoopFixed()
}
```

---

## 🎯 三、随堂练习

### 练习要求：

编写一个 **学生成绩管理系统**，实现以下功能：

#### 1. 定义结构体

```go
type Student struct {
    ID    int
    Name  string
    Score float64
}
```

#### 2. 实现以下方法

**值接收者方法**：

- `ShowInfo()`：显示学生信息（只读）
- `IsPassed()` `bool`：判断是否及格（>= 60分）

**指针接收者方法**：

- `SetScore(score float64)` `error`：设置分数（0-100，否则返回错误）
- `AddBonus(bonus float64)` `error`：加分（不能超过100分）

#### 3. 实现函数

- `calculateAverage(students ...Student)` `float64`：计算平均分（可变参数）
- `findTopStudent(students []Student)` `(Student, error)`：找最高分学生（空切片返回错误）

#### 4. 使用 defer

在 `main` 函数中用 `defer` 输出"程序结束"（最后执行）

#### 5. 错误处理

所有可能出错的地方都要检查并处理错误

---

### 期望输出示例

```
程序开始
===== 学生信息 =====
学号: 1, 姓名: 张三, 分数: 85.00
是否及格: true

学号: 2, 姓名: 李四, 分数: 58.00
是否及格: false

===== 修改分数 =====
李四加分5分后: 63.00
✓ 及格了！

尝试设置无效分数...
Error: 分数必须在 0-100 之间

===== 统计信息 =====
班级平均分: 74.00
最高分学生: 张三 (85.00分)

程序结束
```

---

### 提示

- 使用 `fmt.Errorf()` 创建错误
- 检查错误时用 `if err != nil`
- defer 放在 main 函数开头
- 指针接收者用 `*Student`
- 可变参数用 `...Student`

---

## 📝 提交方式

完成后，将代码保存为 `e:\Golang_study\02_functions_and_methods\exercise.go`，运行后把代码和结果发给我！

---

## 🔑 关键知识点总结

| 知识点       | Go 特点          | vs Java                      |
| ------------ | ---------------- | ---------------------------- |
| **返回值**   | 支持多返回值     | 只能返回1个                  |
| **参数传递** | 全部值传递       | 基本类型值传递，对象引用传递 |
| **方法**     | 可为任何类型定义 | 只能在类中定义               |
| **接收者**   | 值/指针接收者    | this（隐式引用）             |
| **defer**    | 延迟执行（LIFO） | finally（必须配try）         |
| **错误**     | 返回值error      | 异常throw/catch              |
| **panic**    | 不可恢复错误     | 类似RuntimeException         |

**下一节预告**：复合数据类型（Slice、Map、Struct、指针）！
