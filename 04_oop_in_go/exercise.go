package main

import "fmt"

type Shape interface {
	Area() float64
	Perimeter() float64
	GetName() string
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.Radius
}

func (c Circle) GetName() string {
	return "圆形"
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) GetName() string {
	return "矩形"
}

type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

func (t Triangle) Area() float64 {
	s := (t.SideA + t.SideB + t.SideC) / 2
	return sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC))
}

func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t Triangle) GetName() string {
	return "三角形"
}

func sqrt(value float64) float64 {
	z := value
	for i := 0; i < 10; i++ {
		z -= (z*z - value) / (2 * z)
	}
	return z
}

type ColoredShape struct {
	Shape
	Color string
}

func PrintShapeInfo(s Shape) {
	fmt.Printf("形状: %s\n", s.GetName())
	fmt.Printf("面积: %.2f\n", s.Area())
	fmt.Printf("周长: %.2f\n", s.Perimeter())
	fmt.Println()
}

func CompareAreas(s1, s2 Shape) string {
	area1 := s1.Area()
	area2 := s2.Area()

	if area1 > area2 {
		return fmt.Sprintf("%s 的面积 大于 %s 的面积", s1.GetName(), s2.GetName())
	} else if area1 < area2 {
		return fmt.Sprintf("%s 的面积 小于 %s 的面积", s1.GetName(), s2.GetName())
	} else {
		return fmt.Sprintf("%s 的面积 等于 %s 的面积", s1.GetName(), s2.GetName())
	}
}

func TotalArea(shapes ...Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// 空接口any实现描述任意类型的值
func Describe(value any) {
	switch v := value.(type) {
	case int:
		fmt.Printf("整数: %d\n", v)
	case string:
		fmt.Printf("字符串: %s\n", v)
	case Shape:
		fmt.Printf("图形: %s (面积: %.2f)\n", v.GetName(), v.Area())
	default:
		fmt.Printf("未知类型: %T, 值: %v\n", v, v)
	}
}

// ===== 图形信息 =====
// 图形: 圆形
// 面积: 78.54
// 周长: 31.42

// 图形: 矩形
// 面积: 20.00
// 周长: 18.00

// ===== 面积比较 =====
// 圆形 的面积大于 矩形

// ===== 总面积 =====
// 所有图形总面积: 98.54

// ===== 带颜色的图形 =====
// 红色的圆形
// 面积: 78.54

// ===== 类型断言 =====
// 这是一个圆形，半径: 5.00

// ===== 空接口演示 =====
// 整数: 42
// 字符串: Hello Go
// 图形: 圆形 (面积: 78.54)

func main() {
	// 创建图形实例
	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 4, Height: 5}
	triangle := Triangle{SideA: 3, SideB: 4, SideC: 5}

	// 打印图形信息
	fmt.Println("===== 图形信息 =====")
	PrintShapeInfo(circle)
	PrintShapeInfo(rectangle)
	PrintShapeInfo(triangle)

	// 比较面积
	fmt.Println("===== 面积比较 =====")
	fmt.Println(CompareAreas(circle, rectangle))
	fmt.Println()

	// 计算总面积
	fmt.Println("===== 总面积 =====")
	total := TotalArea(circle, rectangle, triangle)
	fmt.Printf("所有图形总面积: %.2f\n\n", total)

	// 带颜色的图形
	fmt.Println("===== 带颜色的图形 =====")
	coloredCircle := ColoredShape{
		Shape: circle,
		Color: "红色",
	}
	fmt.Printf("%s的%s\n", coloredCircle.Color, coloredCircle.GetName())
	fmt.Printf("面积: %.2f\n\n", coloredCircle.Area())

	// 类型断言示例
	fmt.Println("===== 类型断言 =====")
	var s Shape = circle
	if c, ok := s.(Circle); ok {
		fmt.Printf("这是一个%s，半径: %.2f\n\n", s.GetName(), c.Radius)
	} else {
		fmt.Println("未知类型\n")
	}

	// 空接口示例
	fmt.Println("===== 空接口演示 =====")
	Describe(42)
	Describe("Hello Go")
	Describe(circle)
}
