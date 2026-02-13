# Golang 学习路线 - 从 Java 到 Go

> **学习者背景**：有 Java 开发经验的程序员  
> **学习目标**：达到企业级 Go 开发入门水平，掌握 Gin 框架  
> **学习时长**：几天速通  
> **教学方式**：对比教学法（Go vs Java）+ 理论 + 实战

---

## 📚 学习大纲

### **阶段一：Go 语言基础（1-2天）**

#### 第01节：环境搭建与基础类型

- Go 开发环境配置
- 包管理（go mod vs Maven/Gradle）
- 基本数据类型对比（int, string, bool 等）
- 变量声明方式（`:=` vs Java 显式声明）
- 常量与枚举（const/iota vs Java enum）

#### 第02节：函数与方法

- 函数定义（多返回值特性）
- 参数传递（值传递 vs Java 引用传递）
- 方法接收者（Receiver vs Java this）
- defer 延迟执行（vs Java try-finally）
- 错误处理（error vs Java Exception）

#### 第03节：复合数据类型

- 数组与切片（Slice vs Java ArrayList）
- Map（vs Java HashMap）
- 结构体（Struct vs Java Class）
- 指针基础（vs Java 引用）

#### 第04节：面向对象思想在 Go 中的实现

- 接口（Interface）：隐式实现 vs Java 显式实现
- 组合优于继承（Embedding vs Java 继承）
- 多态的实现方式
- 空接口（interface{} 和 any vs Java Object）

---

### **阶段二：Go 语言进阶（1-2天）**

#### 第05节：并发编程基础

- Goroutine vs Java Thread
- Channel：Go 的核心并发工具（vs Java BlockingQueue）
- select 多路复用（vs Java Selector）
- sync 包：Mutex、WaitGroup、Once 等

#### 第06节：并发编程进阶

- Context 上下文管理（vs Java ThreadLocal）
- 并发模式：生产者消费者、工作池等
- 并发安全问题与解决方案
- 竞态检测（go run -race）

#### 第07节：包管理与模块化

- 包的可见性（大小写 vs Java public/private）
- init 函数与包初始化顺序
- go mod 依赖管理详解
- 内部包（internal package）

#### 第08节：错误处理与测试

- error 接口深入（vs Java Exception 体系）
- panic/recover 机制（vs Java try-catch）
- 单元测试（testing 包 vs JUnit）
- 基准测试（Benchmark）
- 表格驱动测试

---

### **阶段三：实战与 Web 开发（1-2天）**

#### 第09节：标准库常用包

- io 包：Reader/Writer 接口（vs Java InputStream/OutputStream）
- json 序列化（encoding/json vs Jackson/Gson）
- http 包：标准库 Web 服务器
- time、strings、fmt 等常用工具包

#### 第10节：Gin 框架入门

- Gin vs Spring Boot 对比
- 路由与路由组（vs Spring MVC @RequestMapping）
- 中间件机制（vs Spring Interceptor/Filter）
- 请求参数绑定与验证

#### 第11节：Gin 框架进阶

- RESTful API 设计
- 文件上传与下载
- JWT 认证实现（vs Spring Security）
- 日志与错误处理中间件

#### 第12节：数据库操作

- database/sql 标准库（vs JDBC）
- GORM 入门（vs MyBatis/Hibernate）
- 数据库连接池管理
- 事务处理

---

### **阶段四：企业级开发实践（可选扩展）**

#### 第13节：项目结构与最佳实践

- 标准项目布局（project-layout）
- 依赖注入模式（vs Spring IoC）
- 配置管理（Viper vs Spring Config）
- 优雅关闭（Graceful Shutdown）

#### 第14节：综合实战项目

- 构建一个完整的 RESTful API 服务
- 集成 Redis 缓存
- 集成消息队列（可选）
- 部署与容器化

---

## 📝 学习进度记录

- [x] 第01节：环境搭建与基础类型
- [x] 第02节：函数与方法
- [x] 第03节：复合数据类型 ⬅️ 当前
- [x] 第04节：面向对象思想在 Go 中的实现
- [x] 第05节：并发编程基础
- [ ] 第06节：并发编程进阶
- [ ] 第07节：包管理与模块化
- [ ] 第08节：错误处理与测试
- [ ] 第09节：标准库常用包
- [ ] 第10节：Gin 框架入门
- [ ] 第11节：Gin 框架进阶
- [ ] 第12节：数据库操作
- [ ] 第13节：项目结构与最佳实践
- [ ] 第14节：综合实战项目

---

## 🎯 当前学习状态

**当前节次**：第06节 - 并发编程进阶  
**完成进度**：5/14  
**最后更新**：2026-02-05

---

## 💡 学习建议

1. **循序渐进**：不要跳过基础部分，Go 的并发模型是核心亮点
2. **动手实践**：每个练习都要自己敲代码，不要复制粘贴
3. **对比思考**：时刻对比 Java，理解设计哲学的差异
4. **看官方文档**：Go 官方文档非常优秀，是最佳参考资料
5. **写地道 Go**：避免用 Java 思维写 Go 代码

---

## 📖 推荐资源

- [Go 官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Gin 官方文档](https://gin-gonic.com/)
