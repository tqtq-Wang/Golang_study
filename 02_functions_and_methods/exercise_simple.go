package main

import "fmt"

// ========== 结构体定义 ==========
type Student struct {
	ID    int
	Name  string
	Score float64
}

// ========== 值接收者方法（只读，不修改对象）==========

// 显示学生信息
func (s Student) ShowInfo() {
	fmt.Printf("学号: %d, 姓名: %s, 分数: %.2f\n", s.ID, s.Name, s.Score)
}

// 判断是否及格（>= 60分）
func (s Student) IsPassed() bool {
	return s.Score >= 60
}

// ========== 指针接收者方法（需要修改对象）==========

// 设置分数（需要验证：0-100）
func (s *Student) SetScore(score float64) error {
	// TODO: 检查分数是否在 0-100 之间
	// 如果不在范围内，返回错误：fmt.Errorf("...")
	// 如果在范围内，设置分数并返回 nil

	if score < 0 || score > 100 {
		return fmt.Errorf("分数必须在 0-100 之间")
	}
	s.Score = score
	return nil
}

// 加分（不能超过100分）
func (s *Student) AddBonus(bonus float64) error {
	// TODO: 检查加分后是否超过 100
	// 如果超过，返回错误
	// 如果不超过，增加分数并返回 nil

	newScore := s.Score + bonus
	if newScore > 100 {
		return fmt.Errorf("加分后总分 %.1f 超过100分", newScore)
	}
	s.Score = newScore
	return nil
}

// ========== 普通函数 ==========

// 计算平均分（可变参数）
func calculateAverage(students ...Student) float64 {
	// TODO: 计算所有学生的平均分
	// 提示：遍历 students，累加分数，除以人数

	if len(students) == 0 {
		return 0
	}

	total := 0.0
	for _, student := range students {
		total += student.Score
	}
	return total / float64(len(students))
}

// 找最高分学生（空切片返回错误）
func findTopStudent(students []Student) (Student, error) {
	// TODO: 如果切片为空，返回错误
	// 否则找到分数最高的学生并返回

	if len(students) == 0 {
		return Student{}, fmt.Errorf("学生列表为空")
	}

	topStudent := students[0]
	for _, student := range students[1:] {
		if student.Score > topStudent.Score {
			topStudent = student
		}
	}
	return topStudent, nil
}

// ========== 主函数 ==========

func main() {
	// defer：最后执行（即使前面出错也会执行）
	defer fmt.Println("\n程序结束")

	fmt.Println("程序开始")
	fmt.Println("===== 学生信息 =====")

	// 创建学生
	student1 := Student{ID: 1, Name: "张三", Score: 85}
	student2 := Student{ID: 2, Name: "李四", Score: 58}

	// 显示信息
	student1.ShowInfo()
	fmt.Printf("是否及格: %v\n\n", student1.IsPassed())

	student2.ShowInfo()
	fmt.Printf("是否及格: %v\n\n", student2.IsPassed())

	// ========== 修改分数 ==========
	fmt.Println("===== 修改分数 =====")

	// 给李四加分
	err := student2.AddBonus(5)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("李四加分后: %.2f\n", student2.Score)
		if student2.IsPassed() {
			fmt.Println("✓ 及格了！\n")
		}
	}

	// 尝试设置无效分数
	fmt.Println("尝试设置无效分数...")
	err = student1.SetScore(150)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println()
	}

	// ========== 统计信息 ==========
	fmt.Println("===== 统计信息 =====")

	// 计算平均分
	avg := calculateAverage(student1, student2)
	fmt.Printf("班级平均分: %.2f\n", avg)

	// 找最高分学生
	students := []Student{student1, student2}
	topStudent, err := findTopStudent(students)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("最高分学生: %s (%.2f分)\n", topStudent.Name, topStudent.Score)
	}

	// defer 会在这里执行，输出"程序结束"
}
