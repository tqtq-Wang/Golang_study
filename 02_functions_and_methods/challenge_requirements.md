# 综合练习：电商订单管理系统

> **难度**：⭐⭐⭐⭐ (中等偏上)  
> **预计时间**：30-40分钟  
> **涵盖知识点**：枚举、结构体、值/指针接收者、错误处理、多返回值、可变参数、defer

---

## 📦 需求说明

实现一个简单的**电商订单管理系统**，支持添加商品、计算总价、修改订单状态等功能。

---

## 📋 任务要求

### 任务1：定义枚举（使用 iota）

定义 **订单状态枚举** `OrderStatus`，包含：

- `Pending` = 1（待支付）
- `Paid` = 2（已支付）
- `Shipping` = 3（配送中）
- `Completed` = 4（已完成）
- `Cancelled` = 5（已取消）

**要求**：

- 定义类型别名 `type OrderStatus int`
- 使用 `iota` 从 1 开始
- 首字母大写（可导出）

---

### 任务2：定义结构体

#### 2.1 商品结构体 `Product`

```go
type Product struct {
    ID    int       // 商品ID
    Name  string    // 商品名称
    Price float64   // 单价
    Stock int       // 库存数量
}
```

#### 2.2 订单项结构体 `OrderItem`

```go
type OrderItem struct {
    Product  Product  // 商品信息
    Quantity int      // 购买数量
}
```

#### 2.3 订单结构体 `Order`

```go
type Order struct {
    ID     int          // 订单ID
    Items  []OrderItem  // 订单项列表
    Status OrderStatus  // 订单状态
}
```

---

### 任务3：为 Product 实现方法

#### 3.1 值接收者方法（只读）

**`ShowInfo()`**：显示商品信息

- 输出格式：`[商品ID] 商品名称 - ¥单价 (库存: N件)`

**`IsAvailable(quantity int) bool`**：判断库存是否充足

- 参数：需要购买的数量
- 返回：库存是否 >= 需要数量

---

#### 3.2 指针接收者方法（修改）

**`UpdateStock(quantity int) error`**：更新库存

- 如果 `quantity` 为正数：增加库存
- 如果 `quantity` 为负数：减少库存
- 如果库存不足（减少后 < 0），返回错误
- 成功返回 `nil`

**示例**：

```go
product := Product{ID: 1, Name: "手机", Price: 3999, Stock: 10}
err := product.UpdateStock(-3)  // 卖出3件，库存变为7
err = product.UpdateStock(5)    // 进货5件，库存变为12
```

---

### 任务4：为 Order 实现方法

#### 4.1 值接收者方法

**`CalculateTotal() float64`**：计算订单总金额

- 遍历所有订单项，计算 `单价 × 数量` 的总和
- 返回总金额

**`GetItemCount() int`**：获取订单中商品的总件数

- 返回所有订单项的数量之和

---

#### 4.2 指针接收者方法

**`AddItem(product Product, quantity int) error`**：添加商品到订单

- 检查数量是否 > 0，否则返回错误
- 检查商品库存是否充足（调用 `product.IsAvailable()`），否则返回错误
- 检查订单状态是否为 `Pending`，否则返回错误："订单已支付，无法修改"
- 成功则添加到 `Items` 切片，返回 `nil`

**`ChangeStatus(newStatus OrderStatus) error`**：修改订单状态

- 状态转换规则（必须按顺序）：
  - `Pending` → 只能变为 `Paid` 或 `Cancelled`
  - `Paid` → 只能变为 `Shipping` 或 `Cancelled`
  - `Shipping` → 只能变为 `Completed`
  - `Completed` 和 `Cancelled` → 不能再修改
- 如果不符合规则，返回错误
- 成功则修改状态，返回 `nil`

**`Cancel() error`**：取消订单

- 如果订单已经是 `Completed` 或 `Shipping`，返回错误："订单已发货，无法取消"
- 否则修改状态为 `Cancelled`，返回 `nil`

---

### 任务5：实现普通函数

**`CreateOrder(id int, products ...Product) (*Order, error)`**：创建订单（可变参数）

- 参数：订单ID 和 任意数量的商品（每个商品数量默认为1）
- 如果没有传入商品，返回错误："订单至少需要一个商品"
- 创建订单对象，初始状态为 `Pending`
- 将每个商品添加到订单（数量为1）
- 返回订单指针和 `nil`

**`FindMostExpensiveItem(order Order) (OrderItem, error)`**：查找订单中最贵的商品

- 如果订单没有商品，返回错误："订单为空"
- 返回单价最高的订单项

---

### 任务6：在 main 函数中使用 defer

在 `main` 函数开头使用 `defer`，输出：

```
========================================
订单管理系统运行结束
========================================
```

---

## 🎯 期望输出示例

```
========================================
订单管理系统启动
========================================

【商品库存】
[1] iPhone 15 - ¥5999.00 (库存: 50件)
[2] AirPods Pro - ¥1999.00 (库存: 100件)
[3] MacBook Pro - ¥12999.00 (库存: 20件)

【创建订单】
✓ 订单 #1001 创建成功
✓ 添加商品: iPhone 15 x1
✓ 添加商品: AirPods Pro x2

【订单详情】
订单ID: 1001
订单状态: 1 (待支付)
商品总件数: 3
订单总金额: ¥9997.00
最贵商品: iPhone 15 (¥5999.00)

【订单流程】
✓ 订单状态变更: Pending → Paid
✓ 订单状态变更: Paid → Shipping
✓ 订单状态变更: Shipping → Completed

【错误处理测试】
✗ 错误: 订单已完成，无法修改状态
✗ 错误: 订单已发货，无法取消

【库存更新】
卖出后库存: iPhone 15 - 49件
进货后库存: iPhone 15 - 69件

========================================
订单管理系统运行结束
========================================
```

---

## 📝 提示

1. **错误消息要清晰**：用 `fmt.Errorf()` 包含具体信息
2. **检查顺序很重要**：先检查参数，再检查状态，最后执行操作
3. **指针接收者选择**：
   - 需要修改对象 → 用 `*Order`, `*Product`
   - 只读操作 → 用 `Order`, `Product`
4. **零值处理**：空切片长度为 0，可以用 `len(slice) == 0` 判断
5. **枚举比较**：直接用 `==` 比较即可

---

## ✅ 提交要求

完成后，将代码保存为：

```
e:\Golang_study\02_functions_and_methods\challenge.go
```

运行命令：

```bash
cd e:\Golang_study\02_functions_and_methods
go run challenge.go
```

把完整代码和运行结果发给我，我会进行详细 Review！

---

## 🌟 挑战加分项（可选）

如果完成基础任务后还有余力，可以尝试：

1. **为 OrderStatus 实现 String() 方法**，自动转换为中文描述

   ```go
   func (s OrderStatus) String() string {
       // 返回 "待支付", "已支付" 等
   }
   ```

2. **实现批量添加商品**：

   ```go
   func (o *Order) AddItems(items []OrderItem) error
   ```

3. **实现订单摘要方法**：
   ```go
   func (o Order) Summary() string
   // 返回格式化的订单摘要字符串
   ```

加油！这道题会让你对 Go 的理解更上一层楼！💪
