package main

import (
	"fmt"
	"slices"
)

type Student struct {
	ID    int
	Name  string
	Age   int
	Class string
}

type ScoreManager struct {
	Students []Student       // 切片：学生列表
	Scores   map[int]float64 // Map：学号 -> 分数
}

func (sm *ScoreManager) AddStudent(s Student) {
	sm.Students = append(sm.Students, s)
}

func (sm *ScoreManager) SetScore(studentID int, score float64) error {
	for _, s := range sm.Students {
		if s.ID == studentID {
			sm.Scores[studentID] = score
			return nil
		}
	}
	return fmt.Errorf("学生ID %d 不存在", studentID)
}

func (sm *ScoreManager) GetScore(studentID int) (float64, error) {
	score, exists := sm.Scores[studentID]
	if !exists {
		return 0, fmt.Errorf("学生ID %d 的成绩不存在", studentID)
	}
	return score, nil
}

func (sm *ScoreManager) GetAverageScore() float64 {
	if len(sm.Scores) == 0 {
		return 0.0
	}
	total := 0.0
	for _, score := range sm.Scores {
		total += score
	}
	return total / float64(len(sm.Scores))
}

func (sm *ScoreManager) GetTopStudents(n int) []Student {
	if n <= 0 || len(sm.Students) == 0 {
		return nil
	}

	if n > len(sm.Students) {
		n = len(sm.Students)
	}

	// 复制所有学生（不是只复制 n 个）
	allStudents := make([]Student, len(sm.Students))
	copy(allStudents, sm.Students)

	// 按成绩降序排序
	slices.SortFunc(allStudents, func(a, b Student) int {
		scoreA := sm.Scores[a.ID]
		scoreB := sm.Scores[b.ID]
		if scoreA > scoreB {
			return -1
		} else if scoreA < scoreB {
			return 1
		} else {
			return 0
		}
	})

	// 返回前 n 个
	return allStudents[:n]
}

func FindStudentByID(students []Student, id int) (*Student, error) {
	for i := range students {
		if students[i].ID == id {
			return &students[i], nil
		}
	}
	return nil, fmt.Errorf("学生ID %d 不存在", id)
}

// 深拷贝学生切片
func CopyStudents(src []Student) []Student {
	dst := make([]Student, len(src))
	copy(dst, src)
	return dst
}

// ===== 添加学生 =====
// ✓ 添加学生: 张三 (ID: 1001)
// ✓ 添加学生: 李四 (ID: 1002)
// ✓ 添加学生: 王五 (ID: 1003)

// ===== 设置成绩 =====
// ✓ 张三的成绩: 95.50
// ✓ 李四的成绩: 88.00
// ✓ 王五的成绩: 92.50

// ===== 查询成绩 =====
// 学生 1001 的成绩: 95.50
// 学生 9999 的成绩: 学生不存在

// ===== 统计信息 =====
// 班级平均分: 92.00

// ===== Top 2 学生 =====
// 1. 张三 - 95.50分
// 2. 王五 - 92.50分

// ===== 测试切片拷贝 =====
// 原切片: [张三 李四 王五]
// 复制后修改不影响原切片: [张三 李四 王五]
func main() {
	sm := ScoreManager{
		Students: []Student{},
		Scores:   make(map[int]float64),
	}

	fmt.Println("===== 添加学生 =====")
	students := []Student{
		{ID: 1001, Name: "张三", Age: 20, Class: "一班"},
		{ID: 1002, Name: "李四", Age: 21, Class: "一班"},
		{ID: 1003, Name: "王五", Age: 19, Class: "一班"},
	}
	for _, s := range students {
		sm.AddStudent(s)
		fmt.Printf("✓ 添加学生: %s (ID: %d)\n", s.Name, s.ID)
	}

	fmt.Println("\n===== 设置成绩 =====")
	scores := map[int]float64{
		1001: 95.5,
		1002: 88.0,
		1003: 92.5,
	}
	for id, score := range scores {
		err := sm.SetScore(id, score)
		if err != nil {
			fmt.Printf("✗ 设置学生 %d 成绩失败: %v\n", id, err)
		} else {
			student, _ := FindStudentByID(sm.Students, id)
			fmt.Printf("✓ %s的成绩: %.2f\n", student.Name, score)
		}
	}

	fmt.Println("\n===== 查询成绩 =====")
	testIDs := []int{1001, 9999}
	for _, id := range testIDs {
		score, err := sm.GetScore(id)
		if err != nil {
			fmt.Printf("学生 %d 的成绩: %v\n", id, err)
		} else {
			fmt.Printf("学生 %d 的成绩: %.2f\n", id, score)
		}
	}

	fmt.Println("\n===== 统计信息 =====")
	avgScore := sm.GetAverageScore()
	fmt.Printf("班级平均分: %.2f\n", avgScore)

	fmt.Println("\n===== Top 2 学生 =====")
	topStudents := sm.GetTopStudents(2)
	for i, s := range topStudents {
		score, _ := sm.GetScore(s.ID)
		fmt.Printf("%d. %s - %.2f分\n", i+1, s.Name, score)
	}

	fmt.Println("\n===== 测试切片拷贝 =====")
	original := []Student{
		{Name: "张三"},
		{Name: "李四"},
		{Name: "王五"},
	}
	copied := CopyStudents(original)
	copied[0].Name = "赵六" // 修改复制的切片

	fmt.Print("原切片: [")
	for i, s := range original {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(s.Name)
	}
	fmt.Println("]")

	fmt.Print("复制后修改: [")
	for i, s := range copied {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(s.Name)
	}
	fmt.Println("]")
	fmt.Println("✓ 复制后修改不影响原切片")
}
