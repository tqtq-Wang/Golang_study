/*
 * @Date: 2026-01-28 23:14:21
 * @LastEditors: tqtqtq
 * @LastEditTime: 2026-01-28 23:26:35
 * @FilePath: \Golang_study\01_environment_and_basics\exercise.go
 */
package main

import "fmt"

const (
	companyName = "阿里巴巴"
	minAge      = 18
)

const (
	probation = iota + 1
	regular
	resigned
)

func main() {
	name := "李四"
	age := 25
	height := 1.75
	isCurrentlyEmployed := true
	fmt.Printf("===== 个人信息 =====\n")
	fmt.Printf("姓名: %s\n", name)
	fmt.Printf("年龄: %d\n", age)
	fmt.Printf("身高: %.1f 米\n", height)
	fmt.Printf("是否在职: %t\n", isCurrentlyEmployed)
	fmt.Printf("公司名称: %s\n", companyName)
	fmt.Println()

	fmt.Printf("===== 年龄检查 =====\n")
	if age >= minAge {
		fmt.Printf("年龄符合入职要求（>= 18岁）\n")
	} else {
		fmt.Printf("年龄不符合入职要求（< 18岁）\n")
	}
	fmt.Println()

	fmt.Printf("===== 员工状态 =====\n")
	fmt.Printf("当前状态: ")
	switch probation {
	case 1:
		fmt.Printf("试用期\n")
	case 2:
		fmt.Printf("正式员工\n")
	case 3:
		fmt.Printf("已离职\n")
	default:
		fmt.Printf("未知状态\n")
	}
}
