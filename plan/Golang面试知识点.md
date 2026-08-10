# Golang 面试知识点与答案（详细版）

> 每个考点附带可直接口述的参考答案。建议先遮住答案自测，再对照补漏。

---

## 一、语言基础

### 1. Go 有哪些基本类型？零值是什么？

**答案：**

基本类型包括：`bool`、整型（`int/int8~int64`、`uint...`）、浮点（`float32/64`）、复数、`byte`（`uint8` 别名）、`rune`（`int32` 别名，表示 Unicode 码点）、`string`。

零值：

- `bool` → `false`
- 数值 → `0`
- `string` → `""`
- 指针 / slice / map / chan / func / interface → `nil`

Go 强调零值可用：如空 `slice` 可 `append`，`Mutex` 零值可直接用。

---

### 2. 哪些类型可以是 nil？

**答案：**

只有引用类型相关的：`pointer`、`slice`、`map`、`channel`、`function`、`interface`。

`struct`、数组、基本类型不能是 `nil`。

注意：`nil slice` 和 `empty slice`（`[]int{}`）都 `len=0`，但 `nil slice` 与 `nil` 比较为 true，空切片不是。

---

### 3. var 和 := 有什么区别？

**答案：**

- `var` 可在包级和函数内使用，可只声明不赋值（得零值），也可指定类型。
- `:=` 只能在函数内，必须初始化，至少有一个变量是新声明的。
- `:=` 容易造成变量遮蔽（shadowing），同名新变量盖住外层变量，是常见坑。

---

### 4. make 和 new 的区别？

**答案：**

- `new(T)`：分配一块内存，初始化为 T 的零值，返回 `*T`。可用于任何类型。
- `make`：只用于 `slice`、`map`、`chan`，完成内部结构初始化，返回类型本身（不是指针）。

例如：`make([]int, 0, 10)` 得到可用切片；`new([]int)` 得到指向 nil 切片的指针，一般不如 `make` 直观。

---

### 5. 数组和切片的区别？切片底层是什么？

**答案：**

- 数组：长度固定，长度是类型一部分；`[3]int` 与 `[4]int` 不同类型；赋值/传参会拷贝整个数组。
- 切片：动态长度，传参拷贝的是 slice header（指针、len、cap），可能共享底层数组。

切片结构大致为：

```text
type slice struct {
    ptr unsafe.Pointer // 指向底层数组
    len int
    cap int
}
```

---

### 6. 切片扩容规则是什么？append 有什么坑？

**答案：**

大致规则（不同版本细节略有差异）：

- 容量较小时代码中常见“约 2 倍”增长；较大时增长因子会降低。
- 扩容会分配新底层数组，并把旧元素拷过去。

常见坑：

1. **共享底层数组**：`b := a[1:3]` 后改 `b` 可能影响 `a`。
2. **append 是否扩容不确定**：没扩容时改新切片会影响原切片；扩容后则不会。
3. **内存泄漏**：从大切片截一小段长期持有，可能拖住整块底层数组；可用 `copy` 到新切片断开引用。

口述要点：`append` 后若返回的切片与原切片底层不同，说明发生了扩容。

---

### 7. nil slice 和 empty slice 有什么区别？

**答案：**

```go
var s1 []int      // nil slice，s1 == nil 为 true
s2 := []int{}     // empty slice，s2 == nil 为 false
s3 := make([]int, 0)
```

`len`/`cap` 都是 0，`json.Marshal` 时：nil slice 常编码为 `null`，空切片为 `[]`（具体行为面试可提一句）。功能上多数场景可互换，序列化时要注意。

---

### 8. map 是线程安全的吗？为什么？

**答案：**

不是。并发读写（至少一个写）会直接 **fatal**（`concurrent map writes` / `read and write`）。

原因：map 底层是哈希表（桶数组 + overflow），扩容、搬迁过程中没有内置锁保护；为了性能默认不加锁。

解决：

- `sync.Mutex` / `RWMutex` 包一层
- 或使用 `sync.Map`（适合读多写少、key 相对稳定、条目相互独立的场景）

---

### 9. 为什么 map 不能取元素地址？遍历顺序为什么随机？

**答案：**

- 不能取地址：map 扩容时元素会搬迁，地址会变，取地址不安全。
- 遍历随机：Go 故意随机化起始位置，避免依赖遍历顺序写出脆弱代码。

---

### 10. string、[]byte、rune 的关系？len(string) 是什么？

**答案：**

- `string`：只读字节序列，通常是 UTF-8。
- `[]byte`：可变字节切片；与 string 互转一般会发生拷贝。
- `rune`：Unicode 码点（`int32`）。

`len(s)` 返回 **字节数**，不是字符数。统计字符用 `utf8.RuneCountInString` 或 `range` 按 rune 遍历。

中文等多字节字符：一个汉字通常占 3 个字节。

---

### 11. 值接收者 vs 指针接收者怎么选？

**答案：**

优先考虑：

1. **方法需要修改接收者** → 必须用指针。
2. **接收者是大结构体** → 指针避免拷贝。
3. **需要保持一致性** → 同一类型方法接收者尽量统一。
4. **包含 sync.Mutex 等不可拷贝字段** → 必须指针。
5. **小的不可变类型（如 time.Time 风格）** → 可用值接收者。

补充：

- 值方法，指针/值都能调；指针方法，值变量也可调（编译器取址），但 **接口中如果方法集是指针接收者，则只有指针类型实现该接口**。

---

### 12. 接口是什么？如何实现接口？

**答案：**

接口定义一组方法签名。Go 是 **隐式实现**：类型拥有接口要求的全部方法即实现，无需 `implements` 关键字。

接口值内部可理解为 `(动态类型, 动态值)`。面向接口编程便于解耦和单测 mock。

小接口原则：如 `io.Reader` 只有一个 `Read`，组合优于大而全接口。

---

### 13. 经典题：为什么接口不等于 nil？（nil 接口陷阱）

**答案：**

```go
func main() {
    var p *int = nil
    var i interface{} = p
    fmt.Println(i == nil) // false
}
```

因为接口的 `type` 是 `*int`，`value` 是 nil；只有 type 和 value 都 nil 时，接口才是 nil。

这是面试超高频坑，尤其出现在返回 `error` 时返回了 `(*MyError)(nil)`。

---

### 14. defer 的执行顺序与参数求值？

**答案：**

- 多个 `defer`：**LIFO**（后进先出）。
- `defer` 注册时 **参数立刻求值**；函数体延迟执行。
- 可修改命名返回值。
- 常用于解锁、关文件、捕获 panic。

```go
func f() (x int) {
    defer func() { x++ }()
    return 1 // 实际返回 2
}
```

---

### 15. panic 和 recover 怎么用？error 怎么处理？

**答案：**

- 业务错误用 `error` 返回，不要用 panic。
- `panic` 用于真正的程序异常（如不可能发生的不变量破坏）。
- `recover` 只能在 `defer` 中有效，常用于隔离 goroutine 崩溃，避免拖垮进程。

现代错误处理：

- `fmt.Errorf("...: %w", err)` 包装
- `errors.Is` 判断哨兵错误
- `errors.As` 提取具体错误类型

---

### 16. init 函数有什么注意点？

**答案：**

- 每个包可有多个 `init`，在 `main` 之前执行。
- 执行顺序：导入依赖 → 包级变量初始化 → `init` → `main`。
- 不要在 `init` 里做复杂、易失败、难测试的逻辑；更推荐显式 `NewXXX()` 初始化。

---

### 17. Go 泛型什么时候用？

**答案：**

适用：通用数据结构、通用算法、减少重复且类型安全。

慎用：仅为“看起来高级”、复杂约束导致可读性下降、接口已经足够表达行为时。

约束常用 `any`、`comparable` 或自定义 interface constraint。

---

## 二、并发编程

### 18. goroutine 和线程的区别？

**答案：**

- OS 线程：内核调度，栈通常较大（MB 级），创建成本高。
- goroutine：用户态，由 Go runtime 调度，初始栈很小（约 2KB 级，可伸缩），成千上万个很常见。
- 多个 goroutine 复用少量 OS 线程（GMP 模型）。

---

### 19. 什么是 GMP？请简述调度过程

**答案：**

- **G**：goroutine
- **M**：worker thread（OS 线程）
- **P**：processor，持有本地运行队列，数量默认约等于 CPU 核数（`GOMAXPROCS`）

调度要点：

1. M 必须绑定 P 才能执行 G。
2. 每个 P 有本地队列；本地没有就去全局队列，或从其他 P **work stealing**。
3. G 发生阻塞型系统调用时，M 可能卡住，P 会与 M 分离，去找其他 M 继续跑别的 G。
4. 现代 Go 支持基于信号的异步抢占，长耗时循环也能被调度。

口述金句：P 是“逻辑处理器”，G 是任务，M 是执行者。

---

### 20. 有缓冲 channel 和无缓冲 channel 区别？

**答案：**

- **无缓冲**：发送和接收必须同时就绪，否则阻塞；强调同步握手。
- **有缓冲**：缓冲区未满发送不阻塞，未空接收不阻塞；用于解耦生产消费速率。

选择：

- 传递信号/同步事件 → 无缓冲更清晰
- 吞吐解耦 → 有缓冲（容量要有依据，不能拍脑袋无限放大）

---

### 21. channel 关闭规则有哪些？

**答案：**

1. 关闭已关闭 channel → panic
2. 向已关闭 channel 发送 → panic
3. 从已关闭 channel 接收：立即返回零值，`ok=false`
4. 关闭 nil channel → panic
5. 对 nil channel 发送/接收会永久阻塞

惯例：**发送方关闭** channel，接收方不要关（除非有明确所有权设计）。

---

### 22. select 的作用和特点？

**答案：**

`select` 等待多个 channel 操作，哪个就绪执行哪个；多个同时就绪时 **随机选一个**，避免饥饿偏见。

- 带 `default`：非阻塞尝试
- 常用于：超时、多路复用、退出信号

```go
select {
case v := <-ch:
    // 处理
case <-ctx.Done():
    return ctx.Err()
case <-time.After(time.Second):
    // 超时
}
```

---

### 23. 如何避免 goroutine 泄漏？

**答案：**

常见泄漏原因：

- 永远阻塞在 channel 收/发
- 没有退出条件的 `for {}`
- HTTP/DB 请求未设超时
- `time.After` 在循环中滥用（旧定时器未回收）

对策：

- 所有后台 goroutine 必须能响应 `context` 取消
- 用 `WaitGroup` / errgroup 管理生命周期
- 设置超时
- 用 pprof goroutine 剖面排查

---

### 24. Mutex 和 RWMutex、channel 怎么选？

**答案：**

**追求数据结构维度的状态保护，选锁（Mutex/RWMutex）；**

**追求并发实体（goroutine）之间的协作、通信和控制流，选通道（channel）。**

---

- **场景属于“数据共享与保护”还是“逻辑编排”？**
  + **数据保护**：需要安全地读写一个 map、slice、struct 或计数器 ➡️ **选 Mutex / RWMutex**。
  + **逻辑编排**：需要通知协作、传递数据、任务分发、生命周期控制 ➡️ **选 channel**。
- **如果确定选锁，看“读写比例”：**
  + **写少读多**（读操作占比超过 80%~90%，且每次读操作有一定耗时） ➡️ **选 RWMutex**。
  + **写操作频繁**，或者读写几乎对半 ➡️ **选 Mutex**（因为 RWMutex 维护读锁计数有额外开销，写多时性能反而不如标准 Mutex）。

**1. Mutex（互斥锁）**

- **本质**：最极端的保护，同一时间**只允许一个** goroutine 访问代码块。
- **优势**：性能极高，内存占用极小。
- **最佳场景**：
  + 保护底层的结构体字段、结构体方法、全局计数器。
  + 写操作非常密集的共享资源。
  + 替代不安全的 `map`（配合 Mutex 封装成并发安全的 map）。

**2. RWMutex（读写锁）**

- **本质**：读写分离。**允许多个读** goroutine 同时访问，但**只允许一个写** goroutine 独占。
- **注意点**：如果读操作非常快（例如只是读取一个整型变量），读锁的加锁/解锁开销甚至可能超过它带来的并发收益。
- **最佳场景**：
  + **配置信息对象**：服务启动时加载一次，运行过程中几乎全是在读取配置，偶尔（如动态配置刷新）才会写一次。
  + **本地缓存/KV 数据库**：绝大多数请求是查询（Get），极少数是写入（Set）。

**3. channel（通道）**

- **本质**：Go 官方推崇的哲学的体现——*“不要通过共享内存来通信，而要通过通信来共享内存”*。它控制的是**数据的流动和控制权**。
- **优势**：天生具备阻塞、唤醒机制，完美支持 CSP 并发模型。
- **最佳场景**：
  + **数据传递**：生产者-消费者模型（如：日志收集器、任务队列、Worker Pool）。
  + **信号通知**：利用 `close(ch)` 实现一键通知所有 goroutine 退出（配合 `context` 底层）。
  + **生命周期编排**：等待多个并发任务结束、超时控制（配合 `time.After`）。

---

---

### 25. WaitGroup 使用注意点？

**答案：**

- `Add` 必须在 `Wait` 之前，且最好在启动 goroutine 前 `Add`
- 每个任务 `Done` 一次，次数要匹配
- `Add` 负数导致计数 `<0` 会 panic
- `WaitGroup` 不能被拷贝（传指针）

---

### 26. sync.Map 适用场景？

**答案：**

适合：

- key 相对稳定
- 读写密集且不同 key 竞争少
- 场景如缓存、只增不改的注册表

不适合：

- 所有请求都频繁读写同一小集合，且需要强一致复杂操作（往往普通 map+Mutex 更可控）

---

### 27. Context 是什么？WithValue 能干什么不能干什么？

**答案：**

`context`: 高并发编程中**最核心**的并发控制工具之一。它主要用来在 goroutine 之间**传递取消信号、超时通知、截止时间以及请求范围内的元数据（Key-Value）**。

`context` 用于在调用链中传递：

- 取消信号
- 截止时间/超时
- 请求级元数据（trace id 等）

API：

- `WithCancel` / `WithTimeout` / `WithDeadline` / `WithValue`

`WithValue`：

- 只放请求范围数据（如 request id、user id）
- **不要**拿它传可选参数、塞一大坨业务依赖，否则隐式耦合难维护

约定：函数参数列表里 `ctx` 放第一个。

- **不要把 Context 放入结构体中**：应该显式地作为函数的**第一个参数**传入，变量名通常固定叫 `ctx`。
- **不要传递 nil**：如果不确定用什么，传入 `context.TODO()`，绝对不能传 `nil`。
- **记得调用 cancel()**：通过 `WithTimeout` 或 `WithCancel` 生成的 `cancel` 函数，务必在函数结束前执行（通常用 `defer cancel()`），否则会导致子上下文无法及时释放，引发**内存泄漏**。
- **WithValue 只能传请求生命周期的元数据**：不要用它来传递业务函数里的可选参数。

我们通常以 `context.Background()`（作为根节点）开始，使用以下四个函数来衍生（Derive）出子 Context，形成一棵 **Context 树**。当父节点被取消时，其所有的子节点都会被自动取消。

**.** `context.WithCancel(parent)` **— 手动取消**

- **作用**：返回一个子 context 和一个 `cancel` 取消函数。调用 `cancel()` 时，会通知所有监听该 context 的 goroutine 停止工作。
- **场景**：用户关闭了连接，或者某个主任务失败了，需要立刻通知所有并发的子协程停止计算。

**2.** `context.WithTimeout(parent, timeout)` **— 超时取消**

- **作用**：指定一个持续时间（如 5 秒）。超时后自动触发取消信号。
- **场景**：微服务 HTTP/RPC 请求，限制接口必须在 2 秒内返回，否则超时报错。

**3.** `context.WithDeadline(parent, time)` **— 截止时间取消**

- **作用**：与 `WithTimeout` 类似，但它接收的是一个具体的时间点（如 2026-12-31 23:59:59）。
- **场景**：需要任务在特定绝对时间前必须结束。

**4.** `context.WithValue(parent, key, val)` **— 传递元数据**

- **作用**：在 Context 树中携带一些只读的、与当前请求生命周期绑定的全局数据。
- **场景**：传递全链路追踪的 `TraceID`、用户身份 `UserID`、内部鉴权的 Token 等。**（注意：绝对不要用它来传递可选的业务参数！）**

在 Go 中，启动一个 goroutine 非常简单，但**关闭它却很难**。

- **痛点**：假设用户发起一个 HTTP 请求，后台启动了 A goroutine，A 又启动了 B，B 又启动了 C。如果用户中途**取消了请求**（比如关闭了浏览器），如果没有机制通知 A、B、C 停止工作，这些 goroutine 就会在后台继续运行，造成**内存泄漏**和 **CPU 资源浪费**。
- **解决方案**：`context` 提供了一颗“树状拓扑图”，当根部的 `context` 被取消时，所有派生出的子 `context` 都会同时收到取消信号，从而优雅地通知所有相关的 goroutine 退出。

---

---

### 28. 如何优雅停止一组 goroutine？

**答案：**

推荐模式：

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

var wg sync.WaitGroup
for i := 0; i < n; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            select {
            case <-ctx.Done():
                return
            case job := <-jobs:
                // handle job
            }
        }
    }()
}

// 出错或收尾
cancel()
wg.Wait()
```

要点：取消可传播 + 等待收尾完成。

---

### 29. 什么是 data race？怎么查？

**答案：**

多个 goroutine 同时访问同一变量，且至少一个写，没有同步，就是 data race。

检测：

```bash
go test -race ./...
go run -race main.go
```

修复：加锁、用 channel 串行化、用 atomic、避免共享。

---

### 30. worker pool 怎么实现？（口述/手写）

**答案思路：**

1. 任务 channel + 结果 channel
2. 启动固定数量 worker
3. 每个 worker `for job := range jobs`
4. 用 `WaitGroup` 或 `errgroup` 等待
5. 用 context 支持取消

控制并发数的本质：限制同时运行的 goroutine 数量（也可用有缓冲 channel 当信号量）。

---

## 三、内存、GC 与性能

### 31. 什么是逃逸分析？哪些情况容易逃逸？

**答案：**

编译器决定变量分配在栈还是堆。分配到堆称为逃逸。

常见逃逸：

- 返回局部变量指针
- 变量大小运行时才知道
- 赋给 interface
- 被闭包引用且生命周期超过当前栈帧
- 切片/映射过大等

查看：

```bash
go build -gcflags="-m" .
```

逃逸会增加堆分配和 GC 压力，但不代表一定差；正确性优先。

---

### 32. Go GC 大概怎么工作？如何减少 GC 压力？

**答案：**

Go 使用并发标记清除思路，配合写屏障；目标是降低 STW 停顿。

面试口述：

- 标记可达对象（三色抽象）
- 清理不可达对象
- 与用户程序并发执行，短 STW

减压力：

- 减少分配（预分配 slice、复用 buffer）
- `sync.Pool` 复用临时对象
- 避免频繁小对象/字符串拼接
- 控制 goroutine 与缓存膨胀

---

### 33. 如何排查 CPU / 内存过高？

**答案：**

1. 开 pprof：`net/http/pprof` 或 `runtime/pprof`
2. CPU：`go tool pprof` 看热点函数
3. 内存：heap profile 看分配与残留
4. goroutine profile 看是否泄漏/阻塞
5. 配合指标：GC 频率、分配速率、goroutine 数

常用：

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

### 34. 字符串拼接怎么做才高效？

**答案：**

循环里避免 `s = s + x`（反复分配）。

用：

- `strings.Builder`
- `bytes.Buffer`
- `strings.Join`

Builder 内部可预增长，减少拷贝。

---

## 四、标准库与工程

### 35. HTTP Client 为什么必须设超时？

**答案：**

`http.DefaultClient` 没有超时，对端卡住会导致 goroutine 一直挂起，最终资源耗尽。

正确做法：自定义 `http.Client`，设置 `Timeout`，或更细的 `Transport`/`context` 超时。

---

### 36. HTTP Server 如何优雅关闭？

**答案：**

```go
server := &http.Server{Addr: ":8080", Handler: mux}

go server.ListenAndServe()

// 收到 SIGINT/SIGTERM
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
server.Shutdown(ctx) // 停接新连接，等旧请求结束
```

同时关闭 DB、刷新日志、停止后台任务。

---

### 37. Go 里时间格式化为什么是 2006-01-02？

**答案：**

Go 使用固定参考时间：`Mon Jan 2 15:04:05 MST 2006`（记忆：1 2 3 4 5 6 7）。

所以布局要写成这个参考时刻的样子，而不是 `YYYY-MM-DD` 这种其他语言常见占位符。

---

### 38. 表驱动测试是什么？

**答案：**

把多个 case 放在切片里循环执行，结构清晰、易扩展：

```go
tests := []struct{
    name string
    in   int
    want int
}{
    {"ok", 1, 2},
    {"zero", 0, 1},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        if got := Inc(tt.in); got != tt.want {
            t.Fatalf("got %d want %d", got, tt.want)
        }
    })
}
```

---

## 五、Web / Gin / API

### 39. Gin 中间件是什么？如何实现鉴权中间件？

**答案：**

中间件是请求前后的拦截逻辑（日志、鉴权、限流、recover）。

思路：

1. 从 Header 取 token
2. 校验 JWT / Session
3. 失败则 `Abort` 并返回 401
4. 成功把 user 写入 context，调用 `c.Next()`

注意：中间件里要用 `c.Request.Context()` 往下传取消信号。

---

### 40. Session 和 JWT 对比？

**答案：**

|     | Session         | JWT             |
| --- | --------------- | --------------- |
| 状态  | 服务端存会话          | 通常无状态，信息在 token |
| 注销  | 易（删会话）          | 需黑名单/短过期+刷新     |
| 扩展  | 多机需共享存储         | 天然易水平扩展         |
| 安全  | 防伪造靠 session id | 靠签名；注意泄露与篡改     |

实践：访问令牌短过期 + 刷新令牌；敏感系统可仍用服务端会话。

---

### 41. 如何设计统一 API 响应与错误码？

**答案：**

建议固定结构，例如：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

- `code=0` 成功，非 0 业务错误
- HTTP 状态码表达传输层结果（401/403/404/500）
- 业务错误与内部 error 分离，避免把内部细节回给客户端
- 日志记录完整 error chain

---

### 42. 常见 Web 安全注意点？

**答案：**

- SQL 注入：参数化查询 / ORM 占位符
- XSS：输出编码；API 场景主要防存储型脏数据
- CSRF：Cookie 会话场景要防；纯 JWT Header 风险不同
- 密码哈希：bcrypt/argon2
- 限流防爆破
- 最小权限、密钥不入库明文

---

## 六、数据库与 Redis

### 43. 事务 ACID 是什么？

**答案：**

- **A 原子性**：要么全成要么全不成
- **C 一致性**：约束始终成立
- **I 隔离性**：并发事务互不干扰程度由隔离级别决定
- **D 持久性**：提交后不丢

---

### 44. MySQL 隔离级别与脏读/不可重复读/幻读？

**答案：**

| 级别 | 脏读 | 不可重复读 | 幻读 |
| -------------- | --- | ----- | -------------------------- |
| 读未提交 | 可能 | 可能 | 可能 |
| 读已提交 | 否 | 可能 | 可能 |
| 可重复读（MySQL 默认） | 否 | 否 | 理论上可能（InnoDB RR + 间隙锁大幅缓解） |
| 串行化 | 否 | 否 | 否 |

- 脏读：读到未提交数据  
- 不可重复读：同一行两次读结果变了（更新导致）  
- 幻读：同一范围两次读行数变了（插入/删除导致）

---

### 45. 索引为什么快？什么是最左前缀？

**答案：**

InnoDB 常用 B+Tree 索引，能减少磁盘 IO，支持排序范围查询。

联合索引 `(a,b,c)` 最左前缀：可用 `a`、`a,b`、`a,b,c`；一般不能跳过 `a` 只靠 `b`。

注意：

- 区分度低的列不适合单独索引
- 避免在索引列上函数运算导致失效
- 覆盖索引可减少回表

---

### 46. 数据库连接池怎么配？（Go）

**答案：**

`database/sql` 自带池：

- `SetMaxOpenConns`：最大打开连接
- `SetMaxIdleConns`：最大空闲
- `SetConnMaxLifetime`：连接最大存活
- `SetConnMaxIdleTime`：空闲多久关闭

原则：

- 不超过 DB 承受能力
- 结合实例数：每实例 open 数 × 实例数 < DB max_connections
- 用 `context` 给查询设超时

---

### 47. GORM 的 N+1 问题是什么？

**答案：**

查列表后再循环查关联，导致 1 次主查询 + N 次关联查询。

解决：预加载 `Preload`、`Joins`，或一次查出后自行组装；避免在循环里查库。

---

### 48. Redis 缓存穿透、击穿、雪崩？

**答案：**

| 问题  | 含义                   | 对策                    |
| --- | -------------------- | --------------------- |
| 穿透  | 查不存在的数据，每次打到 DB      | 布隆过滤器、缓存空值（短 TTL）     |
| 击穿  | 热点 key 过期瞬间高并发打 DB   | 互斥重建、逻辑过期、永不过期+异步刷新   |
| 雪崩  | 大量 key 同时过期或 Redis 挂 | TTL 加随机、多级缓存、限流降级、高可用 |

---

### 49. Redis 分布式锁怎么做？注意什么？

**答案：**

基本：`SET key value NX EX seconds`

注意：

1. value 用唯一 ID，释放时校验，避免误删别人的锁
2. 业务超时要续期（看门狗）或合理设过期
3. 不要假设绝对可靠；Redis 主从切换可能丢锁
4. 解锁建议用 Lua 脚本保证原子性

红锁（Redlock）有争议，面试能说出权衡即可。

---

### 50. Cache Aside（旁路缓存）流程？

**答案：**

读：

1. 读缓存
2. 命中则返回
3. 未命中读 DB，回写缓存

写：

1. 写 DB
2. 删除缓存（推荐删而不是更新，避免并发脏数据更复杂）

并发下仍可能短暂不一致，需结合业务容忍度。

---

## 七、微服务与分布式（加分）

### 51. REST 和 gRPC 怎么选？

**答案：**

- REST/JSON：易调试、浏览器友好、对外 API 常见
- gRPC/Protobuf：性能更好、强契约、适合内部服务、支持流式

很多公司对外 REST，对内 gRPC。

---

### 52. 消息队列解决什么问题？如何保证消息不丢/不重复？

**答案：**

解决：异步解耦、削峰填谷、最终一致。

不丢（至少一次）：

- 生产者确认、持久化、消费者成功后再 ack

不重复（业务幂等）：

- 唯一键去重表、状态机、幂等 token

恰好一次很难纯中间件保证，通常是“至少一次 + 幂等”。

---

### 53. CAP / BASE？

**答案：**

- CAP：分布式系统在分区时，一致性 C 与可用性 A 难以同时完美满足。
- BASE：基本可用、软状态、最终一致；很多互联网系统用最终一致换可用性与性能。

---

### 54. 常见限流算法？

**答案：**

- 固定窗口：实现简单，有临界突刺
- 滑动窗口：更平滑
- 漏桶：匀速流出，平滑流量
- 令牌桶：允许一定突发

单机可用内存计数；分布式常用 Redis + Lua。

---

### 55. 分布式 ID 怎么生成？

**答案：**

常见：

- UUID：简单但无序、较长
- 雪花算法：趋势递增，含时间/机器/序列
- 号段模式：DB 分配号段，本地消化，性能好

注意时钟回拨问题（雪花）。

---

## 八、网络与 OS

### 56. TCP 三次握手、四次挥手？TIME_WAIT 是什么？

**答案：**

三次握手：SYN → SYN+ACK → ACK，确认双方收发能力，同步初始序号。

四次挥手：主动方 FIN → 对端 ACK → 对端 FIN → 主动方 ACK。因为 TCP 全双工，关闭两个方向通常分两次。

TIME_WAIT：主动关闭方等 2MSL，确保最后 ACK 丢失时可重传，并让旧包消失，避免影响新连接。

---

### 57. HTTP 与 HTTPS 区别？HTTP/2 有什么特点？

**答案：**

- HTTPS = HTTP + TLS，加密与身份认证，防窃听篡改。
- HTTP/2：二进制分帧、多路复用、头部压缩、服务器推送（部分特性实际使用有限）。

---

### 58. 进程、线程、协程区别？

**答案：**

- 进程：资源分配单位，隔离强，切换成本高
- 线程：CPU 调度单位，共享进程内存
- 协程（goroutine）：用户态轻量任务，由运行时调度，成本更低

---

## 九、设计与项目

### 59. Go 项目如何分层？

**答案：**

常见：

```text
Handler（接口层）
  → Service（业务层）
    → Repository（数据层）
```

依赖方向向内：上层依赖下层接口，便于替换实现和单测。  
`internal` 包防止被外部错误引用。

---

### 60. 介绍项目用什么结构？

**答案（STAR）：**

1. **背景**：业务是什么，量级如何
2. **任务**：你负责哪块
3. **行动**：技术选型、关键设计、难点怎么解（并发、缓存、一致性、性能）
4. **结果**：QPS、延迟、故障率、成本或效率提升

准备至少 2 个亮点：一个性能优化，一个线上问题排查。

---

### 61. 短链系统怎么设计？（简化）

**答案要点：**

- 生成短码：哈希截断 / 发号器转 62 进制
- 存储：短码 → 长链（Redis + DB）
- 跳转：302/301，注意缓存
- 防冲突、自定义短码、过期、统计点击
- 高并发读多写少：缓存优先

---

### 62. 秒杀怎么设计？（简化）

**答案要点：**

- 前端限流 + 答疑问答/验证码
- 网关/接口限流
- 库存预扣在 Redis，异步落库
- 独立库存服务，避免打爆主库
- 幂等下单、防超卖（Lua 扣减或分段库存）
- 热点 key 优化

---

## 十、高频题速答卡（临考速记）

| 题 | 一句话答案 |
| ----------- | ---------------------------------- |
| make vs new | make 初始化 slice/map/chan；new 返回零值指针 |
| slice 传参 | 拷贝 header，可能改共享底层数组 |
| map 并发 | 不安全，要加锁或 sync.Map |
| 接口 nil | type/value 都 nil 才是 nil |
| defer 顺序 | LIFO，参数先求值 |
| GMP | G 任务，M 线程，P 逻辑处理器 |
| 无缓冲 channel | 同步握手 |
| context | 取消、超时、请求元数据 |
| 优雅停服 | 信号 + Shutdown(ctx) |
| 缓存穿透 | 布隆/空值缓存 |
| 索引最左前缀 | 联合索引从左边连续使用 |
| 至少一次+幂等 | MQ 常见可靠方案 |

---

## 十一、并发手写题参考答案

### 63. 控制并发数为 N

```go
func DoAll(ctx context.Context, jobs []Job, n int) error {
    sem := make(chan struct{}, n) // 信号量
    var wg sync.WaitGroup
    errCh := make(chan error, 1)

    for _, job := range jobs {
        j := job
        wg.Add(1)
        go func() {
            defer wg.Done()
            select {
            case sem <- struct{}{}:
            case <-ctx.Done():
                return
            }
            defer func() { <-sem }()

            if err := j.Run(ctx); err != nil {
                select {
                case errCh <- err:
                default:
                }
            }
        }()
    }

    wg.Wait()
    select {
    case err := <-errCh:
        return err
    default:
        return ctx.Err()
    }
}
```

### 64. 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

ch := make(chan result, 1)
go func() { ch <- doWork() }()

select {
case r := <-ch:
    // 成功
case <-ctx.Done():
    // 超时
}
```

---

## 十二、复习建议

### 按层级抓重点

| 层级    | 必会带答案的章节                  |
| ----- | ------------------------- |
| 初级/校招 | 一、二（到 context）、五、六基础、十速答卡 |
| 中级    | 全员 + 三、七、手写并发             |
| 高级    | 再补源码细节、分布式权衡、系统题          |

### 自测方法

1. 遮住答案，口述 1～2 分钟
2. 对照缺什么补什么
3. 把答案里的例子改成你项目里的真实例子（面试官更爱听）

---

> 文件路径：`Golang面试知识点.md`  
> 建议：先刷「十、高频题速答卡」，再按章节深挖；结合你的 Gin 项目把鉴权、DB、缓存、优雅退出写成可演示代码，答案会更有说服力。
