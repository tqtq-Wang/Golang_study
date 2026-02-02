package main

import "fmt"

// ========== 包级别常量（可被其他包访问，因为首字母大写）==========
const Pi = 3.14159
const AppName = "GoLearn"

// ========== 枚举定义 ==========
type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

// 存储单位
const (
	_  = iota // 跳过 0
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

func main() {
	fmt.Println("========== 第01节：环境搭建与基础类型 ==========\n")

	// ========== 整数类型演示 ==========
	fmt.Println("【整数类型】")
	var age int = 25      // 标准声明
	height := 175         // 短变量声明（最常用）
	var weight int32 = 65 // 指定具体类型

	fmt.Printf("年龄: %d, 类型: %T\n", age, age)
	fmt.Printf("身高: %d cm, 类型: %T\n", height, height)
	fmt.Printf("体重: %d kg, 类型: %T\n", weight, weight)

	// 类型转换（必须显式）
	result := age + int(weight)
	fmt.Printf("年龄+体重 = %d\n\n", result)

	// ========== 浮点数演示 ==========
	fmt.Println("【浮点数类型】")
	var score float64 = 98.5
	price := 19.99 // 默认 float64

	fmt.Printf("考试分数: %.1f 分\n", score)
	fmt.Printf("商品价格: %.2f 元\n\n", price)

	// ========== 布尔类型演示 ==========
	fmt.Println("【布尔类型】")
	isStudent := true
	hasCar := false

	fmt.Printf("是否为学生: %v\n", isStudent)
	fmt.Printf("是否拥有汽车: %v\n", hasCar)

	if isStudent {
		fmt.Println("✓ 享受学生票优惠\n")
	}

	// ========== 字符串演示 ==========
	fmt.Println("【字符串类型】")
	name := "张三"
	greeting := "Hello, Go!"

	fmt.Println("姓名:", name)
	fmt.Println("问候语:", greeting)

	// 字符串拼接
	fullGreeting := greeting + " 欢迎 " + name
	fmt.Println("完整问候:", fullGreeting)

	// 字符串长度（字节数）
	fmt.Printf("'%s' 的字节长度: %d\n", greeting, len(greeting))
	fmt.Printf("'%s' 的字节长度: %d (中文UTF-8编码占3字节)\n\n", name, len(name))

	// ========== 字符类型（rune）演示 ==========
	fmt.Println("【字符类型（rune）】")
	var ch rune = '中'
	fmt.Printf("字符: %c\n", ch)
	fmt.Printf("Unicode 码点: U+%04X\n", ch)
	fmt.Printf("十进制值: %d\n", ch)
	fmt.Printf("类型: %T\n\n", ch)

	// ========== 零值演示 ==========
	fmt.Println("【零值概念】")
	var defaultInt int
	var defaultStr string
	var defaultBool bool
	var defaultFloat float64

	fmt.Printf("int 零值: %d\n", defaultInt)
	fmt.Printf("string 零值: '%s' (长度: %d)\n", defaultStr, len(defaultStr))
	fmt.Printf("bool 零值: %v\n", defaultBool)
	fmt.Printf("float64 零值: %.1f\n\n", defaultFloat)

	// ========== 常量演示 ==========
	fmt.Println("【常量】")
	fmt.Println("圆周率:", Pi)
	fmt.Println("应用名称:", AppName)

	const MaxUsers = 1000
	fmt.Println("最大用户数:", MaxUsers)
	// MaxUsers = 2000  // 编译错误：不能修改常量
	fmt.Println()

	// ========== 枚举（iota）演示 ==========
	fmt.Println("【枚举（iota）】")
	today := Wednesday
	fmt.Printf("今天是星期%d\n", today)

	if today == Wednesday {
		fmt.Println("✓ 今天是星期三\n")
	}

	// ========== 存储单位演示 ==========
	fmt.Println("【存储单位（iota + 位运算）】")
	fileSize := 2 * GB
	fmt.Printf("文件大小: %d 字节\n", fileSize)
	fmt.Printf("文件大小: %.2f GB\n", float64(fileSize)/float64(GB))
	fmt.Printf("文件大小: %.2f MB\n", float64(fileSize)/float64(MB))
	fmt.Printf("文件大小: %.2f KB\n\n", float64(fileSize)/float64(KB))

	// ========== 变量批量声明 ==========
	fmt.Println("【批量声明】")
	var (
		userName  string  = "李四"
		userAge   int     = 30
		userScore float64 = 95.5
		isVIP     bool    = true
	)

	fmt.Printf("用户: %s, 年龄: %d, 分数: %.1f, VIP: %v\n", userName, userAge, userScore, isVIP)

	fmt.Println("\n========== 示例程序结束 ==========")
}
