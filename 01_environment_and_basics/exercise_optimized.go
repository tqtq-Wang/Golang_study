/*
 * 第01节练习 - 优化版本
 * 展示地道的 Go 编程风格
 */
package main

import "fmt"

// ========== 常量定义 ==========
const (
	CompanyName = "阿里巴巴" // 首字母大写表示可导出（public）
	MinAge      = 18
)

// ========== 员工状态枚举 ==========
// 定义类型别名，提供类型安全性
type EmployeeStatus int

const (
	Probation EmployeeStatus = iota + 1 // 1: 试用期
	Regular                             // 2: 正式员工
	Resigned                            // 3: 已离职
)

// 为枚举类型添加 String() 方法（类似 Java 的 toString）
func (s EmployeeStatus) String() string {
	switch s {
	case Probation:
		return "试用期"
	case Regular:
		return "正式员工"
	case Resigned:
		return "已离职"
	default:
		return "未知状态"
	}
}

func main() {
	// ========== 个人信息变量 ==========
	name := "李四"
	age := 25
	var height float64 = 1.75 // 显式声明类型（更清晰）
	employed := true          // Go 推荐简短命名
	status := Regular         // 当前员工状态

	// ========== 输出个人信息 ==========
	fmt.Println("===== 个人信息 =====")
	fmt.Printf("姓名: %s\n", name)
	fmt.Printf("年龄: %d\n", age)
	fmt.Printf("身高: %.2f 米\n", height) // 保留2位小数更精确
	fmt.Printf("是否在职: %t\n", employed)
	fmt.Printf("公司: %s\n", CompanyName)
	fmt.Println()

	// ========== 年龄检查 ==========
	fmt.Println("===== 年龄检查 =====")
	if age >= MinAge {
		fmt.Printf("✓ 年龄符合入职要求（>= %d岁）\n", MinAge)
	} else {
		fmt.Printf("✗ 年龄不符合入职要求（< %d岁）\n", MinAge)
	}
	fmt.Println()

	// ========== 员工状态 ==========
	fmt.Println("===== 员工状态 =====")
	fmt.Printf("当前状态: %d\n", status)

	// 方法1：使用 switch（传统方式）
	switch status {
	case Probation:
		fmt.Println("状态说明: 试用期")
	case Regular:
		fmt.Println("状态说明: 正式员工")
	case Resigned:
		fmt.Println("状态说明: 已离职")
	default:
		fmt.Println("状态说明: 未知状态")
	}

	// 方法2：使用 String() 方法（更优雅）
	// fmt.Printf("状态说明: %s\n", status)  // 自动调用 String() 方法
}
