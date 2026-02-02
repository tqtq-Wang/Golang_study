package main

import "fmt"

// ========== 示例1：最简单的错误处理 ==========

// 除法（可能出错）
func divide(a, b int) (int, error) {
	if b == 0 {
		// 出错了，返回错误
		return 0, fmt.Errorf("除数不能为0")
	}
	// 没错误，返回结果和 nil
	return a / b, nil
}

// ========== 示例2：方法中的错误处理 ==========

type BankAccount struct {
	Owner   string
	Balance float64
}

// 取款（可能失败）
func (acc *BankAccount) Withdraw(amount float64) error {
	// 检查1：金额必须大于0
	if amount <= 0 {
		return fmt.Errorf("取款金额必须大于0")
	}

	// 检查2：余额必须足够
	if amount > acc.Balance {
		return fmt.Errorf("余额不足：需要 %.2f，只有 %.2f", amount, acc.Balance)
	}

	// 都通过了，执行操作
	acc.Balance -= amount
	return nil // 成功返回 nil
}

// ========== 示例3：多个错误检查 ==========

type Student struct {
	Name  string
	Age   int
	Score float64
}

// 验证学生信息（可能有多个错误）
func (s *Student) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("姓名不能为空")
	}

	if s.Age < 0 || s.Age > 150 {
		return fmt.Errorf("年龄 %d 不合法", s.Age)
	}

	if s.Score < 0 || s.Score > 100 {
		return fmt.Errorf("分数 %.1f 不合法", s.Score)
	}

	return nil // 全部验证通过
}

// ========== 主函数：演示错误处理 ==========

func main() {
	fmt.Println("========== 错误处理演示 ==========\n")

	// ========== 示例1：除法 ==========
	fmt.Println("【示例1：除法运算】")

	// 正常情况
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 ÷ 2 = %d ✓\n", result)
	}

	// 除以0的情况
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("10 ÷ 0 错误: %v ✗\n\n", err)
	} else {
		fmt.Printf("10 ÷ 0 = %d\n", result)
	}

	// ========== 示例2：银行账户 ==========
	fmt.Println("【示例2：银行账户】")

	acc := BankAccount{Owner: "张三", Balance: 1000}
	fmt.Printf("初始余额: %.2f\n", acc.Balance)

	// 正常取款
	err = acc.Withdraw(300)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("取款 300 成功，剩余: %.2f ✓\n", acc.Balance)
	}

	// 取款金额超过余额
	err = acc.Withdraw(1000)
	if err != nil {
		fmt.Printf("取款失败: %v ✗\n", err)
	}

	// 取款金额无效
	err = acc.Withdraw(-100)
	if err != nil {
		fmt.Printf("取款失败: %v ✗\n\n", err)
	}

	// ========== 示例3：验证学生信息 ==========
	fmt.Println("【示例3：验证学生信息】")

	// 有效学生
	student1 := Student{Name: "李四", Age: 20, Score: 85}
	err = student1.Validate()
	if err != nil {
		fmt.Printf("验证失败: %v ✗\n", err)
	} else {
		fmt.Printf("学生 %s 信息有效 ✓\n", student1.Name)
	}

	// 无效学生（姓名为空）
	student2 := Student{Name: "", Age: 20, Score: 85}
	err = student2.Validate()
	if err != nil {
		fmt.Printf("验证失败: %v ✗\n", err)
	}

	// 无效学生（分数超范围）
	student3 := Student{Name: "王五", Age: 20, Score: 150}
	err = student3.Validate()
	if err != nil {
		fmt.Printf("验证失败: %v ✗\n", err)
	}

	fmt.Println("\n========== 演示结束 ==========")
}
