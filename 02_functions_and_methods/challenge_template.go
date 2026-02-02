package main

import "fmt"

// ========== 任务1：定义枚举 ==========

// TODO: 定义 OrderStatus 类型和枚举常量
// type OrderStatus int
// const (
//
//	Pending ...
//
// )
type OrderStatus int

const (
	Pending   OrderStatus = iota + 1 // 1: 待支付
	Paid                             // 2: 已支付
	Shipping                         // 3: 发货中
	Completed                        // 4: 已完成
	Canceled                         // 5: 已取消
)

// ========== 任务2：定义结构体 ==========

// TODO: 定义 Product 结构体
// type Product struct { ... }
type Product struct {
	ID    int     // 商品ID
	Name  string  // 商品名称
	Price float64 // 单价
	Stock int     // 库存数量
}

// TODO: 定义 OrderItem 结构体
// type OrderItem struct { ... }
type OrderItem struct {
	Product  Product // 商品信息
	Quantity int     // 购买数量
}

// TODO: 定义 Order 结构体
// type Order struct { ... }
type Order struct {
	ID     int         // 订单ID
	Items  []OrderItem // 订单项列表
	Status OrderStatus // 订单状态
}

// ========== 任务3：Product 的方法 ==========

// TODO: 实现 ShowInfo() - 值接收者
func (p Product) ShowInfo() {
	fmt.Printf("[%d] %s - ¥%.2f (库存: %d件)\n", p.ID, p.Name, p.Price, p.Stock)
}

// TODO: 实现 IsAvailable(quantity int) bool - 值接收者
func (p Product) IsAvailable(quantity int) bool {
	return p.Stock >= quantity
}

// TODO: 实现 UpdateStock(quantity int) error - 指针接收者
func (p *Product) UpdateStock(quantity int) error {
	newStock := p.Stock + quantity
	if newStock < 0 {
		return fmt.Errorf("库存不足，无法减少 %d 件", -quantity)
	}
	p.Stock = newStock
	return nil
}

// ========== 任务4：Order 的方法 ==========

// TODO: 实现 CalculateTotal() float64 - 值接收者
func (o Order) CalculateTotal() float64 {
	total := 0.0
	for _, item := range o.Items {
		total += item.Product.Price * float64(item.Quantity)
	}
	return total
}

// TODO: 实现 GetItemCount() int - 值接收者
func (o Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}

// TODO: 实现 AddItem(product Product, quantity int) error - 指针接收者
func (o *Order) AddItem(product Product, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("购买数量必须大于0")
	}
	if !product.IsAvailable(quantity) {
		return fmt.Errorf("商品 %s 库存不足", product.Name)
	}
	if o.Status != Pending {
		return fmt.Errorf("订单已支付，无法修改")
	}
	// 添加订单项
	o.Items = append(o.Items, OrderItem{Product: product, Quantity: quantity})
	return nil
}

// TODO: 实现 ChangeStatus(newStatus OrderStatus) error - 指针接收者
func (o *Order) ChangeStatus(newStatus OrderStatus) error {
	// 简单状态流转检查
	switch o.Status {
	case Pending:
		if newStatus != Paid && newStatus != Canceled {
			return fmt.Errorf("无法从 %v 变更到 %v", o.Status, newStatus)
		}
	case Paid:
		if newStatus != Shipping && newStatus != Canceled {
			return fmt.Errorf("无法从 %v 变更到 %v", o.Status, newStatus)
		}
	case Shipping:
		if newStatus != Completed {
			return fmt.Errorf("无法从 %v 变更到 %v", o.Status, newStatus)
		}
	case Completed, Canceled:
		return fmt.Errorf("订单已完成或取消，无法变更状态")
	}
	o.Status = newStatus
	return nil
}

// TODO: 实现 Cancel() error - 指针接收者
func (o *Order) Cancel() error {
	if o.Status == Completed || o.Status == Shipping {
		return fmt.Errorf("订单已发货，无法取消")
	}
	o.Status = Canceled
	return nil
}

// ========== 任务5：普通函数 ==========

// TODO: 实现 CreateOrder(id int, products ...Product) (*Order, error)
func CreateOrder(id int, products ...Product) (*Order, error) {
	if len(products) == 0 {
		return nil, fmt.Errorf("订单必须包含至少一个商品")
	}
	order := &Order{
		ID:     id,
		Items:  []OrderItem{},
		Status: Pending,
	}
	for _, product := range products {
		err := order.AddItem(product, 1) // 默认每个商品数量为1
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}

// TODO: 实现 FindMostExpensiveItem(order Order) (OrderItem, error)
func FindMostExpensiveItem(order Order) (OrderItem, error) {
	if len(order.Items) == 0 {
		return OrderItem{}, fmt.Errorf("订单为空")
	}
	expensiveItem := order.Items[0]
	for _, item := range order.Items[1:] {
		if item.Product.Price > expensiveItem.Product.Price {
			expensiveItem = item
		}
	}
	return expensiveItem, nil
}

// ========== 主函数 ==========

func main() {
	// TODO: 在这里添加 defer，输出"订单管理系统运行结束"
	defer fmt.Println("\n订单管理系统运行结束")

	fmt.Println("========================================")
	fmt.Println("订单管理系统启动")
	fmt.Println("========================================\n")

	// ========== 创建测试商品 ==========
	fmt.Println("【商品库存】")

	// TODO: 创建 3 个商品
	product1 := Product{ID: 1, Name: "笔记本电脑", Price: 5999.99, Stock: 10}
	product2 := Product{ID: 2, Name: "智能手机", Price: 3999.50, Stock: 20}
	product3 := Product{ID: 3, Name: "无线耳机", Price: 799.00, Stock: 15}

	// TODO: 显示商品信息
	product1.ShowInfo()
	product2.ShowInfo()
	product3.ShowInfo()

	// ========== 创建订单 ==========
	fmt.Println("\n【创建订单】")

	// TODO: 使用 CreateOrder 创建订单，添加商品
	order, err := CreateOrder(1001, product1, product2)
	if err != nil {
		fmt.Println("✗ 创建订单失败：", err)
		return
	}
	fmt.Printf("✓ 订单 %d 创建成功\n", order.ID)

	// TODO: 使用 AddItem 添加更多商品
	err = order.AddItem(product3, 2)
	if err != nil {
		fmt.Println("✗ 添加商品失败：", err)
	} else {
		fmt.Printf("✓ 添加商品 %s 成功\n", product3.Name)
	}

	// ========== 显示订单信息 ==========
	fmt.Println("\n【订单详情】")

	// TODO: 输出订单ID、状态、商品总件数、总金额
	fmt.Printf("订单ID: %d\n", order.ID)
	fmt.Printf("订单状态: %v\n", order.Status)
	fmt.Printf("商品总件数: %d\n", order.GetItemCount())
	fmt.Printf("订单总金额: ￥%.2f\n", order.CalculateTotal())

	// TODO: 查找最贵商品
	expensiveItem, err := FindMostExpensiveItem(*order)
	if err != nil {
		fmt.Println("✗ 查找最贵商品失败：", err)
	} else {
		fmt.Printf("最贵商品: %s ￥%.2f\n", expensiveItem.Product.Name, expensiveItem.Product.Price)
	}

	// ========== 测试订单状态流转 ==========
	fmt.Println("\n【订单流程】")

	// TODO: 测试状态变更 Pending → Paid → Shipping → Completed
	err = order.ChangeStatus(Paid)
	if err != nil {
		fmt.Println("✗ 状态变更失败：", err)
	} else {
		fmt.Printf("✓ 订单状态变更为 %v\n", order.Status)
	}

	err = order.ChangeStatus(Shipping)
	if err != nil {
		fmt.Println("✗ 状态变更失败：", err)
	} else {
		fmt.Printf("✓ 订单状态变更为 %v\n", order.Status)
	}

	err = order.ChangeStatus(Completed)
	if err != nil {
		fmt.Println("✗ 状态变更失败：", err)
	} else {
		fmt.Printf("✓ 订单状态变更为 %v\n", order.Status)
	}

	// ========== 测试错误处理 ==========
	fmt.Println("\n【错误处理测试】")

	// TODO: 尝试修改已完成订单的状态（应该失败）
	err = order.ChangeStatus(Pending)
	if err != nil {
		fmt.Printf("✗ 状态变更失败：%v\n", err)
	} else {
		fmt.Printf("✓ 订单状态变更为 %v\n", order.Status)
	}
	// TODO: 尝试取消已发货订单（应该失败）
	err = order.Cancel()
	if err != nil {
		fmt.Printf("✗ 取消订单失败：%v\n", err)
	} else {
		fmt.Printf("✓ 订单已取消，当前状态：%v\n", order.Status)
	}

	// ========== 测试库存更新 ==========
	fmt.Println("\n【库存更新】")

	// TODO: 更新库存（减少和增加）
	err = product1.UpdateStock(-2)
	if err != nil {
		fmt.Printf("✗ 更新库存失败：%v\n", err)
	} else {
		fmt.Printf("✓ 商品 %s 库存减少 2 件，当前库存：%d 件\n", product1.Name, product1.Stock)
	}

	err = product2.UpdateStock(5)
	if err != nil {
		fmt.Printf("✗ 更新库存失败：%v\n", err)
	} else {
		fmt.Printf("✓ 商品 %s 库存增加 5 件，当前库存：%d 件\n", product2.Name, product2.Stock)
	}

	// defer 会在这里执行
}
