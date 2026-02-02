package main

import "fmt"

// ==================== 示例1：接口的基本使用 ====================

// 定义接口
type Speaker interface {
	Speak() string
	GetName() string
}

// 实现1：Dog
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "汪汪汪"
}

func (d Dog) GetName() string {
	return d.Name
}

// 实现2：Cat
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "喵喵喵"
}

func (c Cat) GetName() string {
	return c.Name
}

// 多态函数
func MakeSound(s Speaker) {
	fmt.Printf("%s 说: %s\n", s.GetName(), s.Speak())
}

// ==================== 示例2：支付系统（多态） ====================

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
	fmt.Printf("[支付宝] 账号 %s 支付 %.2f 元\n", a.Account, amount)
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
	fmt.Printf("[微信支付] 账号 %s 支付 %.2f 元\n", w.Account, amount)
	return nil
}

func (w WeChatPay) GetName() string {
	return "微信支付"
}

// 统一支付处理（多态）
func ProcessPayment(pm PaymentMethod, amount float64) {
	fmt.Printf("使用 %s 进行支付\n", pm.GetName())
	if err := pm.Pay(amount); err != nil {
		fmt.Println("支付失败:", err)
	} else {
		fmt.Println("支付成功！")
	}
}

// ==================== 示例3：类型断言 ====================

func DescribeAnimal(s Speaker) {
	// 方式1：类型断言（需要判断 ok）
	if dog, ok := s.(Dog); ok {
		fmt.Printf("这是一只狗，名字是 %s\n", dog.Name)
		return
	}

	if cat, ok := s.(Cat); ok {
		fmt.Printf("这是一只猫，名字是 %s\n", cat.Name)
		return
	}

	fmt.Println("未知动物")
}

func DescribeAnimalSwitch(s Speaker) {
	// 方式2：类型开关（推荐）
	switch v := s.(type) {
	case Dog:
		fmt.Printf("类型开关: 这是狗 %s\n", v.Name)
	case Cat:
		fmt.Printf("类型开关: 这是猫 %s\n", v.Name)
	default:
		fmt.Printf("类型开关: 未知类型 %T\n", v)
	}
}

// ==================== 示例4：空接口（any） ====================

// 打印任意类型的值
func PrintAny(value any) {
	fmt.Printf("值: %v, 类型: %T\n", value, value)
}

// 处理不同类型
func ProcessValue(value any) {
	switch v := value.(type) {
	case int:
		fmt.Printf("整数: %d, 双倍: %d\n", v, v*2)
	case string:
		fmt.Printf("字符串: %s, 长度: %d\n", v, len(v))
	case []int:
		fmt.Printf("整数切片: %v, 元素个数: %d\n", v, len(v))
	case Speaker:
		fmt.Printf("会说话的动物: %s 说 %s\n", v.GetName(), v.Speak())
	default:
		fmt.Printf("未知类型: %T\n", v)
	}
}

// ==================== 示例5：组合（Embedding） ====================

// 基础类型：Person
type Person struct {
	Name string
	Age  int
}

func (p Person) SayHello() {
	fmt.Printf("你好，我是 %s，今年 %d 岁\n", p.Name, p.Age)
}

// 嵌入 Person
type Employee struct {
	Person     // 匿名字段（嵌入）
	EmployeeID string
	Department string
	Salary     float64
}

func (e Employee) Work() {
	fmt.Printf("%s 在 %s 部门工作\n", e.Name, e.Department)
}

// ==================== 示例6：方法提升 ====================

type Engine struct {
	Power int
	Brand string
}

func (e Engine) Start() {
	fmt.Printf("发动机启动！品牌: %s, 功率: %d 马力\n", e.Brand, e.Power)
}

func (e Engine) Stop() {
	fmt.Println("发动机熄火")
}

type Car struct {
	Engine // 嵌入
	Model  string
	Color  string
}

func (c Car) Drive() {
	fmt.Printf("%s 色的 %s 正在行驶\n", c.Color, c.Model)
}

// ==================== 示例7：覆盖嵌入类型的方法 ====================

type Animal struct {
	Name string
}

func (a Animal) Speak() {
	fmt.Println("动物叫声")
}

func (a Animal) Eat() {
	fmt.Printf("%s 在吃东西\n", a.Name)
}

type DogWithOverride struct {
	Animal
	Breed string
}

// 覆盖 Animal.Speak()
func (d DogWithOverride) Speak() {
	fmt.Printf("%s（%s）: 汪汪汪！\n", d.Name, d.Breed)
}

// ==================== 示例8：多重嵌入 ====================

type Flyer struct {
	Name string
}

func (f Flyer) Fly() {
	fmt.Printf("%s 在飞翔\n", f.Name)
}

type Swimmer struct {
	Name string
}

func (s Swimmer) Swim() {
	fmt.Printf("%s 在游泳\n", s.Name)
}

// 嵌入多个类型
type Duck struct {
	Flyer
	Swimmer
}

// ==================== 示例9：接口组合 ====================

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
		return "", fmt.Errorf("文件 %s 未打开", f.Name)
	}
	return f.Content, nil
}

func (f *File) Write(content string) error {
	if !f.IsOpen {
		return fmt.Errorf("文件 %s 未打开", f.Name)
	}
	f.Content = content
	return nil
}

func (f *File) Close() error {
	if !f.IsOpen {
		return fmt.Errorf("文件 %s 已关闭", f.Name)
	}
	f.IsOpen = false
	fmt.Printf("文件 %s 已关闭\n", f.Name)
	return nil
}

func (f *File) Open() error {
	f.IsOpen = true
	fmt.Printf("文件 %s 已打开\n", f.Name)
	return nil
}

// ==================== 主函数 ====================

func main() {
	fmt.Println("==================== 示例1：接口的基本使用 ====================")
	dog := Dog{Name: "旺财"}
	cat := Cat{Name: "小白"}

	MakeSound(dog)
	MakeSound(cat)

	fmt.Println("\n==================== 示例2：支付系统（多态） ====================")
	alipay := Alipay{Account: "user@example.com"}
	wechat := WeChatPay{Account: "user123"}

	ProcessPayment(alipay, 100.50)
	fmt.Println()
	ProcessPayment(wechat, 200.00)

	fmt.Println("\n==================== 示例3：类型断言 ====================")
	DescribeAnimal(dog)
	DescribeAnimal(cat)
	fmt.Println()
	DescribeAnimalSwitch(dog)
	DescribeAnimalSwitch(cat)

	fmt.Println("\n==================== 示例4：空接口（any） ====================")
	PrintAny(123)
	PrintAny("Hello Go")
	PrintAny([]int{1, 2, 3})
	PrintAny(dog)

	fmt.Println("\n--- 处理不同类型 ---")
	ProcessValue(42)
	ProcessValue("Go语言")
	ProcessValue([]int{1, 2, 3, 4, 5})
	ProcessValue(dog)

	fmt.Println("\n==================== 示例5：组合（Embedding） ====================")
	emp := Employee{
		Person:     Person{Name: "张三", Age: 30},
		EmployeeID: "E001",
		Department: "技术部",
		Salary:     10000,
	}

	// 可以直接访问嵌入类型的字段和方法
	emp.SayHello() // 调用 Person.SayHello()
	emp.Work()

	// 也可以通过类型名访问
	fmt.Printf("员工 ID: %s, 工资: %.2f\n", emp.EmployeeID, emp.Salary)
	fmt.Printf("通过 Person 访问: %s, %d 岁\n", emp.Person.Name, emp.Person.Age)

	fmt.Println("\n==================== 示例6：方法提升 ====================")
	car := Car{
		Engine: Engine{Power: 200, Brand: "丰田"},
		Model:  "凯美瑞",
		Color:  "白",
	}

	car.Start() // 自动调用 Engine.Start()（方法提升）
	car.Drive()
	car.Stop() // 自动调用 Engine.Stop()

	fmt.Println("\n==================== 示例7：覆盖嵌入类型的方法 ====================")
	dog2 := DogWithOverride{
		Animal: Animal{Name: "旺财"},
		Breed:  "柴犬",
	}

	dog2.Speak()        // 调用 DogWithOverride.Speak()（覆盖）
	dog2.Eat()          // 调用 Animal.Eat()（未覆盖）
	dog2.Animal.Speak() // 显式调用 Animal.Speak()

	fmt.Println("\n==================== 示例8：多重嵌入 ====================")
	duck := Duck{
		Flyer:   Flyer{Name: "唐老鸭"},
		Swimmer: Swimmer{Name: "唐老鸭"},
	}

	duck.Fly()  // 调用 Flyer.Fly()
	duck.Swim() // 调用 Swimmer.Swim()

	// 访问字段时需要指定类型（有歧义）
	fmt.Printf("Flyer Name: %s\n", duck.Flyer.Name)
	fmt.Printf("Swimmer Name: %s\n", duck.Swimmer.Name)

	fmt.Println("\n==================== 示例9：接口组合 ====================")
	file := &File{Name: "test.txt"}
	file.Open()

	// 使用接口
	var rwc ReadWriteCloser = file

	// 写入
	if err := rwc.Write("Hello, Go!"); err != nil {
		fmt.Println("写入失败:", err)
	}

	// 读取
	content, err := rwc.Read()
	if err != nil {
		fmt.Println("读取失败:", err)
	} else {
		fmt.Println("读取内容:", content)
	}

	// 关闭
	rwc.Close()

	// 尝试读取已关闭的文件
	_, err = rwc.Read()
	if err != nil {
		fmt.Println("错误:", err)
	}
}
