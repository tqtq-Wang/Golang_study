# 第04节：面向对象思想在 Go 中的实现

> **本节目标**：理解 Go 的面向对象设计，掌握接口、组合、多态的使用  
> **前置知识**：第01-03节的内容（不会用到并发、包管理等后续知识）

---

## 📖 一、Go 的面向对象哲学

### 1.1 Go vs Java 的面向对象对比

#### **Java 的面向对象三大特性**

1. **封装**：private/public/protected
2. **继承**：extends、super
3. **多态**：接口/抽象类、向上转型

#### **Go 的面向对象特性**

1. **封装**：大小写控制可见性
2. **组合**：嵌入（Embedding）代替继承
3. **多态**：接口（Interface）隐式实现

**核心理念**：

- ❌ **Go 没有类**（class）
- ❌ **Go 没有继承**（extends）
- ✅ **Go 有结构体+方法**
- ✅ **Go 有接口**（隐式实现）
- ✅ **组合优于继承**

---

## 📖 二、接口（Interface）

### 2.1 接口的基本概念

#### **Java 的接口（显式实现）**

```java
// 定义接口
public interface Speaker {
    void speak();
    String getName();
}

// 实现接口（必须显式声明 implements）
public class Dog implements Speaker {
    private String name;

    @Override
    public void speak() {
        System.out.println("汪汪汪！");
    }

    @Override
    public String getName() {
        return this.name;
    }
}
```

#### **Go 的接口（隐式实现）**

```go
// 定义接口
type Speaker interface {
    Speak()
    GetName() string
}

// 定义结构体
type Dog struct {
    Name string
}

// 实现方法（无需声明实现接口）
func (d Dog) Speak() {
    fmt.Println("汪汪汪！")
}

func (d Dog) GetName() string {
    return d.Name
}

// Dog 自动实现了 Speaker 接口（编译器推断）
```

**关键差异**：
| 特性 | Go | Java |
|------|-----|------|
| **实现方式** | 隐式（自动） | 显式（implements） |
| **耦合度** | 低（接口独立） | 高（类必须知道接口） |
| **灵活性** | 高（事后定义接口） | 低（必须提前声明） |
| **接口位置** | 可以在任何包 | 必须在实现前 |

---

### 2.2 接口的定义与实现

#### **接口定义规则**

```go
// 接口命名：通常以 -er 结尾
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 小接口更好（单一职责）
type Closer interface {
    Close() error
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}
```

**Go 的接口设计哲学**：

- ✅ **小接口**：1-3 个方法最好
- ✅ **单一职责**：一个接口做一件事
- ❌ **避免大接口**：不要像 Java 那样定义几十个方法

---

#### **接口实现示例**

```go
// 定义接口
type Animal interface {
    Speak() string
    Move() string
}

// 实现1：Dog
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "汪汪汪"
}

func (d Dog) Move() string {
    return "跑"
}

// 实现2：Bird
type Bird struct {
    Name string
}

func (b Bird) Speak() string {
    return "叽叽喳喳"
}

func (b Bird) Move() string {
    return "飞"
}

// 使用接口（多态）
func MakeSound(a Animal) {
    fmt.Printf("%s: %s\n", a.Speak(), a.Move())
}

func main() {
    dog := Dog{Name: "旺财"}
    bird := Bird{Name: "小鸟"}

    MakeSound(dog)   // 汪汪汪: 跑
    MakeSound(bird)  // 叽叽喳喳: 飞
}
```

---

### 2.3 接口的多态性

#### **多态示例：支付系统**

```go
// 支付接口
type PaymentMethod interface {
    Pay(amount float64) error
    GetName() string
}

// 实现1：支付宝
type Alipay struct {
    Account string
}

func (a Alipay) Pay(amount float64) error {
    fmt.Printf("支付宝支付: %.2f 元\n", amount)
    return nil
}

func (a Alipay) GetName() string {
    return "支付宝"
}

// 实现2：微信支付
type WeChatPay struct {
    Account string
}

func (w WeChatPay) Pay(amount float64) error {
    fmt.Printf("微信支付: %.2f 元\n", amount)
    return nil
}

func (w WeChatPay) GetName() string {
    return "微信支付"
}

// 统一支付接口（多态）
func ProcessPayment(pm PaymentMethod, amount float64) error {
    fmt.Printf("使用 %s 支付\n", pm.GetName())
    return pm.Pay(amount)
}

func main() {
    alipay := Alipay{Account: "user@example.com"}
    wechat := WeChatPay{Account: "user123"}

    ProcessPayment(alipay, 100.50)  // 支付宝支付
    ProcessPayment(wechat, 200.00)  // 微信支付
}
```

---

### 2.4 接口的类型断言

**类型断言**：判断接口变量的实际类型

```go
// 方式1：简单断言
var a Animal = Dog{Name: "旺财"}
dog, ok := a.(Dog)  // 类型断言
if ok {
    fmt.Println("是 Dog 类型:", dog.Name)
}

// 方式2：类型开关（推荐）
func DescribeAnimal(a Animal) {
    switch v := a.(type) {
    case Dog:
        fmt.Printf("这是狗: %s\n", v.Name)
    case Bird:
        fmt.Printf("这是鸟: %s\n", v.Name)
    default:
        fmt.Println("未知动物")
    }
}
```

**Java 对比**：

```java
// Java 使用 instanceof
if (animal instanceof Dog) {
    Dog dog = (Dog) animal;
    System.out.println("这是狗: " + dog.getName());
}
```

---

### 2.5 空接口（interface{} 和 any）

#### **空接口的概念**

**空接口**：没有任何方法的接口，可以表示**任意类型**

```go
// Go 1.18 之前
var x interface{}  // 可以赋值为任意类型

// Go 1.18+ 推荐（等价）
var x any  // any 是 interface{} 的别名
```

**类比 Java**：

```java
Object obj = "Hello";     // Object 是所有类的父类
obj = 123;
obj = new ArrayList<>();
```

#### **空接口的使用**

```go
// 存储任意类型
func PrintAny(value any) {
    fmt.Printf("值: %v, 类型: %T\n", value, value)
}

func main() {
    PrintAny(123)              // 值: 123, 类型: int
    PrintAny("Hello")          // 值: Hello, 类型: string
    PrintAny([]int{1, 2, 3})   // 值: [1 2 3], 类型: []int
}

// 切片存储不同类型
values := []any{1, "hello", 3.14, true}
for _, v := range values {
    fmt.Printf("%v ", v)
}
```

#### **空接口的类型断言**

```go
func ProcessValue(value any) {
    switch v := value.(type) {
    case int:
        fmt.Printf("整数: %d\n", v)
    case string:
        fmt.Printf("字符串: %s\n", v)
    case []int:
        fmt.Printf("整数切片: %v\n", v)
    default:
        fmt.Printf("未知类型: %T\n", v)
    }
}

func main() {
    ProcessValue(42)
    ProcessValue("Go语言")
    ProcessValue([]int{1, 2, 3})
}
```

---

## 📖 三、组合（Embedding）

### 3.1 Java 的继承 vs Go 的组合

#### **Java 的继承（extends）**

```java
// 父类
public class Animal {
    protected String name;

    public void eat() {
        System.out.println(name + " 在吃东西");
    }
}

// 子类继承父类
public class Dog extends Animal {
    public Dog(String name) {
        this.name = name;
    }

    public void bark() {
        System.out.println(name + " 在叫");
    }
}

// 使用
Dog dog = new Dog("旺财");
dog.eat();   // 继承自 Animal
dog.bark();  // Dog 自己的方法
```

#### **Go 的组合（Embedding）**

**Go 没有继承**，使用**嵌入**（组合）实现类似功能：

```go
// 基础类型
type Animal struct {
    Name string
}

func (a Animal) Eat() {
    fmt.Printf("%s 在吃东西\n", a.Name)
}

// 嵌入 Animal（组合）
type Dog struct {
    Animal  // 匿名字段（嵌入）
    Breed string
}

func (d Dog) Bark() {
    fmt.Printf("%s 在叫\n", d.Name)  // 可以直接访问 Animal 的字段
}

func main() {
    dog := Dog{
        Animal: Animal{Name: "旺财"},
        Breed:  "柴犬",
    }

    dog.Eat()   // 自动调用 Animal.Eat()
    dog.Bark()  // Dog 自己的方法
    fmt.Println(dog.Name)  // 直接访问嵌入字段
}
```

---

### 3.2 嵌入的语法

#### **匿名字段（Anonymous Field）**

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) SayHello() {
    fmt.Printf("你好，我是 %s\n", p.Name)
}

// 嵌入 Person
type Employee struct {
    Person      // 匿名字段（字段名就是类型名）
    EmployeeID string
    Salary     float64
}

func main() {
    emp := Employee{
        Person:     Person{Name: "张三", Age: 30},
        EmployeeID: "E001",
        Salary:     10000,
    }

    // 可以直接访问 Person 的字段和方法
    fmt.Println(emp.Name)      // 张三
    fmt.Println(emp.Age)       // 30
    emp.SayHello()             // 你好，我是 张三

    // 也可以通过类型名访问
    fmt.Println(emp.Person.Name)  // 张三
}
```

---

### 3.3 方法提升（Method Promotion）

**嵌入类型的方法会自动提升到外层类型**

```go
type Engine struct {
    Power int
}

func (e Engine) Start() {
    fmt.Printf("发动机启动，功率: %d\n", e.Power)
}

type Car struct {
    Engine  // 嵌入
    Brand string
}

func main() {
    car := Car{
        Engine: Engine{Power: 200},
        Brand:  "丰田",
    }

    car.Start()  // 自动调用 Engine.Start()（方法提升）
    // 等价于 car.Engine.Start()
}
```

---

### 3.4 覆盖嵌入类型的方法

**外层类型可以覆盖嵌入类型的方法**

```go
type Animal struct {
    Name string
}

func (a Animal) Speak() {
    fmt.Println("动物叫声")
}

type Dog struct {
    Animal
}

// 覆盖 Animal.Speak()
func (d Dog) Speak() {
    fmt.Printf("%s: 汪汪汪\n", d.Name)
}

func main() {
    dog := Dog{Animal: Animal{Name: "旺财"}}

    dog.Speak()         // 汪汪汪（调用 Dog.Speak）
    dog.Animal.Speak()  // 动物叫声（调用 Animal.Speak）
}
```

**Java 对比**：

```java
// Java 使用 @Override
@Override
public void speak() {
    System.out.println("汪汪汪");
}

// 调用父类方法
super.speak();
```

---

### 3.5 多重嵌入

**可以嵌入多个类型（类似多继承，但更安全）**

```go
type Flyer interface {
    Fly()
}

type Swimmer interface {
    Swim()
}

type Bird struct {
    Name string
}

func (b Bird) Fly() {
    fmt.Printf("%s 在飞\n", b.Name)
}

type Fish struct {
    Name string
}

func (f Fish) Swim() {
    fmt.Printf("%s 在游泳\n", f.Name)
}

// 嵌入多个类型
type Duck struct {
    Bird  // 可以飞
    Fish  // 可以游泳
}

func main() {
    duck := Duck{
        Bird: Bird{Name: "唐老鸭"},
        Fish: Fish{Name: "唐老鸭"},
    }

    duck.Fly()   // 唐老鸭 在飞
    duck.Swim()  // 唐老鸭 在游泳

    // 注意：访问 Name 字段时需要指定类型（有歧义）
    fmt.Println(duck.Bird.Name)  // 唐老鸭
    fmt.Println(duck.Fish.Name)  // 唐老鸭
}
```

---

## 📖 四、接口与组合的结合

### 4.1 经典设计模式：策略模式

```go
// 定义策略接口
type SortStrategy interface {
    Sort(data []int) []int
}

// 策略1：冒泡排序
type BubbleSort struct{}

func (b BubbleSort) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)

    for i := 0; i < len(result)-1; i++ {
        for j := 0; j < len(result)-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    return result
}

// 策略2：快速排序（简化版）
type QuickSort struct{}

func (q QuickSort) Sort(data []int) []int {
    // 简化实现...
    return data
}

// 上下文（使用策略）
type Sorter struct {
    Strategy SortStrategy
}

func (s Sorter) DoSort(data []int) []int {
    return s.Strategy.Sort(data)
}

func main() {
    data := []int{5, 2, 8, 1, 9}

    // 使用冒泡排序
    sorter := Sorter{Strategy: BubbleSort{}}
    result := sorter.DoSort(data)
    fmt.Println("冒泡排序:", result)

    // 切换为快速排序
    sorter.Strategy = QuickSort{}
    result = sorter.DoSort(data)
    fmt.Println("快速排序:", result)
}
```

---

### 4.2 接口组合示例：文件系统

```go
// 定义小接口
type Reader interface {
    Read() (string, error)
}

type Writer interface {
    Write(content string) error
}

type Closer interface {
    Close() error
}

// 组合接口
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// 实现
type File struct {
    Name    string
    Content string
    IsOpen  bool
}

func (f *File) Read() (string, error) {
    if !f.IsOpen {
        return "", fmt.Errorf("文件未打开")
    }
    return f.Content, nil
}

func (f *File) Write(content string) error {
    if !f.IsOpen {
        return fmt.Errorf("文件未打开")
    }
    f.Content = content
    return nil
}

func (f *File) Close() error {
    if !f.IsOpen {
        return fmt.Errorf("文件已关闭")
    }
    f.IsOpen = false
    return nil
}

func (f *File) Open() error {
    f.IsOpen = true
    return nil
}

// 使用
func ProcessFile(rwc ReadWriteCloser) error {
    defer func() {
        if err := rwc.Close(); err != nil {
            fmt.Println("关闭失败:", err)
        }
    }()

    content, err := rwc.Read()
    if err != nil {
        return err
    }

    fmt.Println("读取内容:", content)
    return rwc.Write("新内容")
}
```

---

## 💻 二、代码示例

完整示例代码请查看 `example.go`

---

## 🎯 三、随堂练习

### 练习要求：实现一个"图形计算系统"

#### **功能需求**：

1. 定义 `Shape` 接口，包含方法：
   - `Area() float64` - 计算面积
   - `Perimeter() float64` - 计算周长
   - `GetName() string` - 获取图形名称

2. 实现三种图形：
   - **Circle**（圆形）：半径
   - **Rectangle**（矩形）：长、宽
   - **Triangle**（三角形）：三条边

3. 定义 `ColoredShape` 结构体，嵌入任意 Shape 并添加颜色属性

4. 实现函数：
   - `PrintShapeInfo(s Shape)` - 打印图形信息
   - `CompareArea(s1, s2 Shape) string` - 比较面积
   - `TotalArea(shapes ...Shape) float64` - 计算总面积

5. 使用空接口 `any` 实现：
   - `Describe(value any)` - 描述任意类型的值

---

### 期望输出示例

```
===== 图形信息 =====
图形: 圆形
面积: 78.54
周长: 31.42

图形: 矩形
面积: 20.00
周长: 18.00

===== 面积比较 =====
圆形 的面积大于 矩形

===== 总面积 =====
所有图形总面积: 98.54

===== 带颜色的图形 =====
红色的圆形
面积: 78.54

===== 类型断言 =====
这是一个圆形，半径: 5.00

===== 空接口演示 =====
整数: 42
字符串: Hello Go
图形: 圆形 (面积: 78.54)
```

---

## 📝 提交方式

完成后，将代码保存为 `e:\Golang_study\04_oop_in_go\exercise.go`

---

## 🔑 关键知识点总结

| 概念       | Go 特点           | vs Java         | 使用建议             |
| ---------- | ----------------- | --------------- | -------------------- |
| **接口**   | 隐式实现          | 显式 implements | 小接口，单一职责     |
| **继承**   | 无，用组合        | extends         | 组合优于继承         |
| **多态**   | 接口实现          | 接口/抽象类     | 面向接口编程         |
| **空接口** | any / interface{} | Object          | 类型断言获取实际类型 |
| **嵌入**   | 匿名字段          | 继承            | 方法自动提升         |

**Go 哲学**：

- ✅ **小接口**：1-3 个方法
- ✅ **隐式实现**：低耦合
- ✅ **组合优于继承**：灵活可扩展
- ✅ **面向接口编程**：依赖抽象而非具体

**下一节预告**：并发编程基础（Goroutine、Channel）！
