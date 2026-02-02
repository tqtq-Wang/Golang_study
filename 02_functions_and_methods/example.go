package main

import (
	"fmt"
	"math"
)

// ========== 示例1：多返回值 ==========

// 计算圆的面积和周长（命名返回值）
func calculateCircle(radius float64) (area float64, circumference float64) {
	area = math.Pi * radius * radius
	circumference = 2 * math.Pi * radius
	return // 裸返回
}

// 除法运算（返回商、余数、错误）
func divide(a, b int) (quotient int, remainder int, err error) {
	if b == 0 {
		err = fmt.Errorf("division by zero")
		return // 返回零值和错误
	}
	quotient = a / b
	remainder = a % b
	return
}

// 可变参数函数
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// ========== 示例2：方法接收者 ==========

type BankAccount struct {
	Owner   string
	Balance float64
}

// 值接收者（只读）
func (acc BankAccount) ShowInfo() {
	fmt.Printf("【账户信息】持有人: %s, 余额: ¥%.2f\n", acc.Owner, acc.Balance)
}

// 值接收者（判断）
func (acc BankAccount) IsRich() bool {
	return acc.Balance >= 10000
}

// 指针接收者（修改）
func (acc *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("存款金额必须大于0")
	}
	acc.Balance += amount
	fmt.Printf("✓ 存入 ¥%.2f，当前余额: ¥%.2f\n", amount, acc.Balance)
	return nil
}

// 指针接收者（修改）
func (acc *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("取款金额必须大于0")
	}
	if amount > acc.Balance {
		return fmt.Errorf("余额不足：需要 ¥%.2f，只有 ¥%.2f", amount, acc.Balance)
	}
	acc.Balance -= amount
	fmt.Printf("✓ 取出 ¥%.2f，当前余额: ¥%.2f\n", amount, acc.Balance)
	return nil
}

// ========== 示例3：defer 延迟执行 ==========

func deferDemo() {
	defer fmt.Println("defer 1 (最后执行)")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3 (最先执行)")
	fmt.Println("正常执行")
}

// defer 修改返回值
func calculate() (result int) {
	defer func() {
		result += 100 // 修改命名返回值
		fmt.Println("defer 中修改返回值")
	}()
	result = 50
	return // 实际返回 150
}

// defer 捕获 panic
func safeDivide(a, b int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到 panic: %v\n", r)
			result = 0 // 设置默认返回值
		}
	}()

	result = a / b // 可能触发 panic
	return
}

// ========== 示例4：错误处理 ==========

// 自定义错误类型
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("字段 '%s' 验证失败: %s", e.Field, e.Message)
}

// 验证年龄
func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "不能为负数"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "不能超过150"}
	}
	return nil
}

// ========== 主函数 ==========

func main() {
	// defer 在函数最开始声明（最后执行）
	defer fmt.Println("\n========== 程序结束 ==========")

	fmt.Println("========== 第02节：函数与方法 ==========\n")

	// ========== 多返回值演示 ==========
	fmt.Println("【多返回值】")
	area, circum := calculateCircle(5.0)
	fmt.Printf("半径5的圆 - 面积: %.2f, 周长: %.2f\n", area, circum)

	// 忽略部分返回值
	area, _ = calculateCircle(10.0)
	fmt.Printf("半径10的圆 - 面积: %.2f\n\n", area)

	// ========== 错误处理演示 ==========
	fmt.Println("【错误处理】")
	q, r, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 ÷ 3 = %d 余 %d\n", q, r)
	}

	// 除以0
	_, _, err = divide(10, 0)
	if err != nil {
		fmt.Printf("✗ Error: %v\n\n", err)
	}

	// ========== 可变参数演示 ==========
	fmt.Println("【可变参数】")
	fmt.Println("sum(1,2,3) =", sum(1, 2, 3))
	fmt.Println("sum(1,2,3,4,5) =", sum(1, 2, 3, 4, 5))

	nums := []int{10, 20, 30, 40, 50}
	fmt.Printf("sum(%v) = %d\n\n", nums, sum(nums...))

	// ========== 方法接收者演示 ==========
	fmt.Println("【方法接收者】")
	acc := BankAccount{Owner: "张三", Balance: 5000}
	acc.ShowInfo()

	// 值接收者判断
	if acc.IsRich() {
		fmt.Println("✓ 是富人")
	} else {
		fmt.Println("✗ 还不够富")
	}

	// 指针接收者修改
	err = acc.Deposit(6000)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if acc.IsRich() {
		fmt.Println("✓ 现在是富人了")
	}

	// 取款
	err = acc.Withdraw(5000)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 取款失败
	err = acc.Withdraw(10000)
	if err != nil {
		fmt.Printf("✗ Error: %v\n\n", err)
	}

	// ========== defer 演示 ==========
	fmt.Println("【defer 延迟执行】")
	deferDemo()
	fmt.Println()

	fmt.Println("【defer 修改返回值】")
	result := calculate()
	fmt.Printf("最终返回值: %d\n\n", result)

	fmt.Println("【defer 捕获 panic】")
	result = safeDivide(10, 2)
	fmt.Printf("10 / 2 = %d\n", result)

	result = safeDivide(10, 0)
	fmt.Printf("10 / 0 = %d (默认值)\n\n", result)

	// ========== 自定义错误演示 ==========
	fmt.Println("【自定义错误】")
	err = validateAge(25)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✓ 年龄25：验证通过")
	}

	err = validateAge(-5)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	}

	err = validateAge(200)
	if err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	}

	// 最后会执行 defer 输出"程序结束"
}
