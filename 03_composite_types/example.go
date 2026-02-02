package main

import (
	"fmt"
	"sort"
)

// ========== 一、数组演示 ==========

func arrayDemo() {
	fmt.Println("========== 数组（Array）==========\n")

	// 数组声明与初始化
	var arr1 [5]int                      // 零值：[0 0 0 0 0]
	arr2 := [3]int{1, 2, 3}              // 完整初始化
	arr3 := [5]int{1, 2}                 // 部分初始化：[1 2 0 0 0]
	arr4 := [...]int{10, 20, 30, 40, 50} // 自动推断长度

	fmt.Println("零值数组:", arr1)
	fmt.Println("完整初始化:", arr2)
	fmt.Println("部分初始化:", arr3)
	fmt.Println("自动推断:", arr4)

	// 数组是值类型（重要！）
	fmt.Println("\n【数组是值类型】")
	a := [3]int{1, 2, 3}
	b := a                               // b 是 a 的副本
	b[0] = 999                           // 修改 b
	fmt.Printf("a = %v, b = %v\n", a, b) // a 不受影响

	// 数组遍历
	fmt.Println("\n【数组遍历】")
	scores := [5]int{85, 92, 78, 95, 88}

	// 方式1：传统 for
	for i := 0; i < len(scores); i++ {
		fmt.Printf("scores[%d] = %d\n", i, scores[i])
	}

	// 方式2：range（推荐）
	fmt.Println("\n使用 range:")
	for index, value := range scores {
		fmt.Printf("索引:%d, 值:%d\n", index, value)
	}
}

// ========== 二、切片演示 ==========

func sliceDemo() {
	fmt.Println("\n========== 切片（Slice）==========\n")

	// 切片创建
	fmt.Println("【切片创建】")

	// 方式1：字面量
	slice1 := []int{1, 2, 3, 4, 5}
	fmt.Println("字面量:", slice1)

	// 方式2：make
	slice2 := make([]int, 3, 5) // 长度3，容量5
	fmt.Printf("make创建: %v, 长度:%d, 容量:%d\n", slice2, len(slice2), cap(slice2))

	// 方式3：从数组切片
	arr := [5]int{10, 20, 30, 40, 50}
	slice3 := arr[1:4] // [20 30 40]
	fmt.Println("从数组切片:", slice3)

	// append 添加元素
	fmt.Println("\n【append 添加元素】")
	var slice []int // nil 切片
	fmt.Printf("初始: %v (nil? %v)\n", slice, slice == nil)

	slice = append(slice, 1)
	slice = append(slice, 2, 3, 4)
	fmt.Println("append后:", slice)

	// 合并切片
	slice4 := []int{5, 6, 7}
	slice = append(slice, slice4...) // ... 展开切片
	fmt.Println("合并切片:", slice)

	// 切片截取
	fmt.Println("\n【切片截取】")
	s := []int{0, 1, 2, 3, 4, 5}
	fmt.Println("原切片:", s)
	fmt.Println("s[1:4]:", s[1:4]) // [1 2 3]
	fmt.Println("s[:3]:", s[:3])   // [0 1 2]
	fmt.Println("s[2:]:", s[2:])   // [2 3 4 5]
	fmt.Println("s[:]:", s[:])     // [0 1 2 3 4 5]

	// 切片的陷阱
	fmt.Println("\n【切片陷阱：共享底层数组】")
	arr2 := [5]int{1, 2, 3, 4, 5}
	s1 := arr2[1:4] // [2 3 4]
	s2 := arr2[2:5] // [3 4 5]

	fmt.Printf("原数组: %v\n", arr2)
	fmt.Printf("s1: %v\n", s1)
	fmt.Printf("s2: %v\n", s2)

	s1[1] = 999 // 修改 s1
	fmt.Printf("修改s1后 - 原数组: %v\n", arr2)
	fmt.Printf("修改s1后 - s2: %v (也受影响)\n", s2)

	// 切片复制（避免共享）
	fmt.Println("\n【切片复制】")
	src := []int{10, 20, 30}
	dst := make([]int, len(src))
	copy(dst, src) // 复制
	dst[0] = 999

	fmt.Printf("src: %v (不受影响)\n", src)
	fmt.Printf("dst: %v\n", dst)
}

// ========== 三、Map 演示 ==========

func mapDemo() {
	fmt.Println("\n========== Map（映射）==========\n")

	// Map 创建
	fmt.Println("【Map 创建】")

	// 方式1：make
	m1 := make(map[string]int)
	m1["张三"] = 25
	m1["李四"] = 30
	fmt.Println("make创建:", m1)

	// 方式2：字面量
	m2 := map[string]int{
		"Go":     10,
		"Java":   20,
		"Python": 15,
	}
	fmt.Println("字面量:", m2)

	// Map 操作
	fmt.Println("\n【Map 操作】")
	scores := make(map[string]float64)

	// 添加/修改
	scores["张三"] = 95.5
	scores["李四"] = 88.0
	scores["王五"] = 92.5
	fmt.Println("添加后:", scores)

	// 获取值（推荐方式）
	score, exists := scores["张三"]
	if exists {
		fmt.Printf("张三的成绩: %.1f\n", score)
	}

	score, exists = scores["赵六"]
	if !exists {
		fmt.Println("赵六不存在，返回零值:", score) // 0.0
	}

	// 删除
	delete(scores, "李四")
	fmt.Println("删除李四后:", scores)

	// 遍历
	fmt.Println("\n【Map 遍历】")
	fmt.Println("遍历键值对:")
	for name, score := range scores {
		fmt.Printf("  %s: %.1f\n", name, score)
	}

	// 只遍历键
	fmt.Println("只遍历键:")
	for name := range scores {
		fmt.Printf("  %s\n", name)
	}

	// Map 是引用类型
	fmt.Println("\n【Map 是引用类型】")
	map1 := map[string]int{"a": 1}
	map2 := map1 // map2 和 map1 指向同一个 map
	map2["b"] = 2

	fmt.Printf("map1: %v (受影响)\n", map1)
	fmt.Printf("map2: %v\n", map2)
}

// ========== 四、指针演示 ==========

type Person struct {
	Name string
	Age  int
}

func pointerDemo() {
	fmt.Println("\n========== 指针（Pointer）==========\n")

	// 基本类型指针
	fmt.Println("【基本类型指针】")
	x := 10
	ptr := &x // 取地址
	fmt.Printf("x的值: %d\n", x)
	fmt.Printf("x的地址: %p\n", ptr)
	fmt.Printf("指针指向的值: %d\n", *ptr) // 解引用

	*ptr = 20 // 通过指针修改
	fmt.Printf("修改后x的值: %d\n", x)

	// 结构体指针
	fmt.Println("\n【结构体指针】")
	p1 := Person{Name: "张三", Age: 25}
	p2 := p1 // 值拷贝
	p2.Age = 30

	fmt.Printf("p1: %+v (不受影响)\n", p1)
	fmt.Printf("p2: %+v\n", p2)

	// 使用指针
	fmt.Println("\n使用指针:")
	ptrPerson := &p1   // 指向 p1
	ptrPerson.Age = 35 // Go 自动解引用
	fmt.Printf("p1: %+v (受影响)\n", p1)

	// new 函数
	fmt.Println("\n【new 函数】")
	ptrInt := new(int) // 返回 *int，指向零值
	*ptrInt = 100
	fmt.Printf("ptrInt指向的值: %d\n", *ptrInt)

	ptrPerson2 := new(Person) // 返回 *Person
	ptrPerson2.Name = "李四"
	ptrPerson2.Age = 28
	fmt.Printf("ptrPerson2: %+v\n", *ptrPerson2)

	// 值传递 vs 指针传递
	fmt.Println("\n【值传递 vs 指针传递】")
	person := Person{Name: "王五", Age: 20}

	modifyByValue(person)
	fmt.Printf("值传递后: %+v (未改变)\n", person)

	modifyByPointer(&person)
	fmt.Printf("指针传递后: %+v (已改变)\n", person)
}

// 值传递：不能修改原对象
func modifyByValue(p Person) {
	p.Age = 99
	fmt.Println("  函数内修改:", p.Age)
}

// 指针传递：可以修改原对象
func modifyByPointer(p *Person) {
	p.Age = 99
	fmt.Println("  函数内修改:", p.Age)
}

// ========== 五、综合示例：学生管理 ==========

type Student struct {
	ID   int
	Name string
	Age  int
}

// 使用指针避免大结构体拷贝
func updateStudent(s *Student, newAge int) {
	s.Age = newAge
}

// 查找学生（返回指针）
func findStudent(students []Student, id int) *Student {
	for i := range students {
		if students[i].ID == id {
			return &students[i] // 返回指针
		}
	}
	return nil
}

func comprehensiveDemo() {
	fmt.Println("\n========== 综合示例：学生管理 ==========\n")

	// 切片存储学生
	students := []Student{
		{ID: 1001, Name: "张三", Age: 20},
		{ID: 1002, Name: "李四", Age: 21},
		{ID: 1003, Name: "王五", Age: 19},
	}

	// Map 存储成绩
	scores := map[int]float64{
		1001: 95.5,
		1002: 88.0,
		1003: 92.5,
	}

	fmt.Println("【学生列表】")
	for _, s := range students {
		score := scores[s.ID]
		fmt.Printf("ID:%d, 姓名:%s, 年龄:%d, 成绩:%.1f\n", s.ID, s.Name, s.Age, score)
	}

	// 查找学生（指针）
	fmt.Println("\n【查找学生】")
	student := findStudent(students, 1002)
	if student != nil {
		fmt.Printf("找到: %+v\n", *student)

		// 修改学生信息
		updateStudent(student, 22)
		fmt.Printf("修改后: %+v\n", students[1]) // 原切片也改变了
	}

	// 计算平均分
	fmt.Println("\n【统计信息】")
	total := 0.0
	for _, score := range scores {
		total += score
	}
	avg := total / float64(len(scores))
	fmt.Printf("平均分: %.2f\n", avg)

	// 排序（按成绩降序）
	fmt.Println("\n【成绩排序】")
	// 创建临时切片用于排序
	type StudentScore struct {
		Student Student
		Score   float64
	}

	var studentScores []StudentScore
	for _, s := range students {
		studentScores = append(studentScores, StudentScore{
			Student: s,
			Score:   scores[s.ID],
		})
	}

	// 排序
	sort.Slice(studentScores, func(i, j int) bool {
		return studentScores[i].Score > studentScores[j].Score // 降序
	})

	fmt.Println("排名:")
	for i, ss := range studentScores {
		fmt.Printf("%d. %s - %.1f分\n", i+1, ss.Student.Name, ss.Score)
	}
}

// ========== 主函数 ==========

func main() {
	fmt.Println("========== 第03节：复合数据类型 ==========")

	arrayDemo()         // 数组
	sliceDemo()         // 切片
	mapDemo()           // Map
	pointerDemo()       // 指针
	comprehensiveDemo() // 综合示例

	fmt.Println("\n========== 示例程序结束 ==========")
}
