# Golang 后端求职知识点总结（超详细版）

本文件对应 [golang_backend_job_plan_4weeks.md](file:///d:/Code/terraform-provider-huaweicloud/docs/golang_backend_job_plan_4weeks.md) 的学习与面试准备范围，目标是覆盖后端业务岗 80% 以上高频考点，并能直接映射到你的项目实践与简历表达。

使用方式建议：

- 不要试图一次看完。按“项目需要 + 面试高频”滚动复习。
- 每个大章至少做 1 个“可验证产出”：代码、笔记、压测截图、EXPLAIN、测试用例、复盘文档。
- 面试前 24 小时：只看“高频问答清单”和你项目的“3 个难点讲稿”。

---

## 1. Go 语言基础（面试必考）

### 1.1 基本类型、零值与类型转换

- **零值**：int 为 0，string 为 ""，bool 为 false，指针/slice/map/channel/function/interface 为 nil，struct 字段为各自零值
- **类型转换是显式的**：`int64(x)`、`[]byte(s)`、`string(b)`（注意 string/[]byte 转换会复制数据）
- **rune 与 byte**：
  - byte 是 uint8，通常表示字节
  - rune 是 int32，通常表示 Unicode code point
- **字符串不可变**：任何“修改”都会生成新字符串；对性能敏感的拼接使用 `strings.Builder` 或 `bytes.Buffer`

常问：

- 为什么 `len("中文")` 不是 2？
  - `len` 返回字节数；中文在 UTF-8 下常为 3 字节；`utf8.RuneCountInString` 统计字符数

### 1.2 指针、值传递、引用语义

- Go 的参数传递永远是“值传递”
  - 传指针：值是“地址”
  - 传 slice/map：值是“描述符”（指向底层数据结构的头部信息），因此看起来像引用
- **何时用指针接收者**：
  - 方法需要修改接收者
  - 避免大对象拷贝
  - 保持方法集合一致（接口实现一致性）

常问：

- slice 传参后修改元素是否会影响外部？
  - 修改元素会影响（共享底层数组）；append 可能扩容导致不影响（取决于是否扩容）

### 1.3 struct、方法、组合（embedding）

- **struct** 是数据载体；方法定义在类型上
- **组合优于继承**：
  - embedding 会“提升”字段/方法（promoted）
  - 发生同名冲突需要显式指定
- **tag** 常用于序列化与校验：`json:"name"`、`db:"name"`

### 1.4 interface（高频难点）

- interface 由两部分组成：
  - 动态类型（concrete type）
  - 动态值（value）
- **nil interface 坑**：
  - `var err error = (*MyErr)(nil)`：此时 err != nil，因为动态类型存在
- **类型断言**：
  - `v, ok := x.(T)`，避免 panic
- **接口设计建议**：
  - 小接口、面向使用方（consumer-driven）
  - 用行为命名：`io.Reader`、`io.Writer`

常问：

- interface 的底层结构是什么？为什么会有 nil 坑？
  - 因为 interface 包含类型信息，只要类型存在就不为 nil

### 1.5 error 处理（项目与面试都很看重）

- Go 的 error 是接口：`type error interface { Error() string }`
- **推荐风格**：
  - 早返回：if err != nil { return ... }
  - error 要携带上下文（哪个步骤/哪个参数）
- **错误封装**（Go 1.13+）：
  - `fmt.Errorf("xxx: %w", err)` 包装
  - `errors.Is/As` 做判断与解包
- **业务错误 vs 系统错误**：
  - 业务错误：参数非法、资源不存在、权限不足 -> 映射 HTTP 4xx
  - 系统错误：DB/Redis 不可用、超时 -> 映射 HTTP 5xx

项目落地建议：

- 统一错误返回结构：`code/message/request_id`，并把内部错误记录到日志里
- 不要把 DB 原始错误直接返回给前端（泄露结构与安全风险）

### 1.6 slice（几乎必问）

- slice 结构：指针（ptr）+ 长度（len）+ 容量（cap）
- **append**：
  - 不扩容：在原数组上追加
  - 扩容：分配新数组，复制旧数据
- **copy**：复制 min(len(dst), len(src)) 个元素
- **常见坑**：
  - 共享底层数组导致“意外修改”
  - 循环中引用切片元素地址（for range 变量复用的问题）

必会解释：

- `make([]T, len, cap)` 的区别，len 决定可访问元素数量，cap 决定可扩展容量

### 1.7 map（几乎必问）

- map 是哈希表，不保证遍历顺序
- **并发不安全**：并发读写会 panic 或 data race
- **读不存在 key**：返回零值；用 `v, ok := m[k]` 判断是否存在
- **map 的 value 是 struct** 时修改字段要注意：
  - `m[k].Field = ...` 可能无法编译（取到的是拷贝）
  - 常用做法：value 用指针或取出来修改再写回

并发场景应对：

- 读多写少：`sync.RWMutex`
- 特定场景：`sync.Map`（读多写少、key 稳定、并发访问）

### 1.8 defer、panic、recover（中高频）

- defer 的执行时机：函数返回前，后进先出（LIFO）
- defer 参数求值：注册 defer 时就求值（闭包变量另说）
- panic：异常中止当前 goroutine 的执行栈回溯
- recover：只能在 defer 中生效，用于捕获 panic 并恢复流程

项目落地：

- HTTP 入口统一 recover，避免服务崩溃
- recover 后要记录堆栈（至少记录 panic 信息 + request-id）

### 1.9 for range 变量复用坑（高频实战坑）

- `for _, v := range xs` 中 v 会复用同一个变量地址
- 常见 bug：起 goroutine 时捕获 v，导致全部拿到最后一个值
- 解决：
  - 在循环体内重新声明：`v := v`
  - 或使用索引访问：`for i := range xs { v := xs[i] }`

### 1.10 包管理与构建

- go mod 基础：
  - `go.mod` 声明 module 与依赖
  - `go.sum` 记录依赖校验和
- 常用命令：
  - `go test ./...`
  - `go test -race ./...`（排查 data race）
  - `go test -run TestName -v`
  - `go env`

---

## 2. 并发与运行时（GMP、channel、锁、泄漏）

### 2.1 goroutine 与 GMP 调度模型（概念要能讲）

- G：goroutine（任务）
- M：OS 线程（machine）
- P：处理器（processor，持有 runqueue）
- 核心理解：
  - goroutine 很轻量，调度由 runtime 管理
  - P 的数量默认等于 GOMAXPROCS（通常等于 CPU 核心数）

常问：

- goroutine 为什么比线程轻量？
  - 初始栈小且可增长；调度成本更低；用户态调度为主

### 2.2 channel 语义（必会）

- channel 三种操作：send、recv、close
- 缓冲/无缓冲：
  - 无缓冲：发送与接收同步握手
  - 有缓冲：缓冲未满发送不阻塞；缓冲为空接收阻塞
- close 规则：
  - 只能由发送方 close（原则上）
  - close 后：
    - 接收方还能读出缓冲区剩余数据
    - 读完后继续读得到零值，且 ok=false
  - 向已关闭 channel 发送会 panic

select 要点：

- 多个 case 同时就绪会随机选择（避免饥饿）
- `default` 会让 select 变成非阻塞

项目落地常用模式：

- worker pool：任务队列 + 多 worker 消费
- fan-out/fan-in：并行执行 + 汇总结果
- done channel / context：用于取消与退出

### 2.3 sync 包（锁、原子、并发工具）

- `sync.Mutex`：
  - 互斥锁，保护临界区
  - 避免锁粒度过大（影响吞吐）或过小（复杂且容易死锁）
- `sync.RWMutex`：
  - 读多写少场景更优
  - 注意写锁可能饥饿（看实现细节与使用方式）
- `sync.WaitGroup`：
  - 等待一组 goroutine 完成
  - 规则：Add 与 Done 成对；最好在启动 goroutine 前 Add
- `sync.Once`：
  - 保证初始化只执行一次（例如单例初始化）
- `sync.Cond`：
  - 条件变量，适合复杂同步（面试中可讲但项目一般少用）
- `sync/atomic`：
  - 原子操作，用于计数、状态位、CAS

### 2.4 数据竞争与 -race

- data race：多个 goroutine 并发访问同一内存位置且至少一个写
- 常见触发点：
  - map 并发写
  - 未加锁写共享变量
  - 结构体中包含引用类型字段的并发读写
- `go test -race` 能在运行时插桩检测数据竞争

### 2.5 context（面试与项目高频）

核心目的：

- 控制“生命周期”：取消、超时、截止时间
- 传递请求范围的元信息（request-id、trace-id 等）

常用 API：

- `context.Background()`：根 context
- `context.WithCancel(parent)`：主动取消
- `context.WithTimeout(parent, d)`：超时
- `context.WithDeadline(parent, t)`：截止时间
- `ctx.Done()`：取消信号 channel
- `ctx.Err()`：取消原因（Canceled/DeadlineExceeded）
- `context.WithValue`：
  - 只放请求范围元信息
  - 不放大对象/不当成全局变量/不传业务参数

项目落地：

- HTTP 入口生成 request-id，写入 ctx，并在日志中带上
- DB/Redis 调用使用 ctx，设置合理超时（不要无限等待）

### 2.6 goroutine 泄漏（高频排障点）

常见原因：

- 永久阻塞在 channel recv/send（没人写/没人读）
- select 没有监听取消信号（ctx.Done）
- 定时器/ticker 未停止（ticker 忘记 Stop）
- 请求已结束但后台 goroutine 还在等资源

应对原则：

- 任何后台 goroutine 都要有“退出条件”
- 统一用 context 控制生命周期
- 对外部依赖调用都加超时

---

## 3. HTTP/REST 与后端业务开发核心

### 3.1 HTTP 基础（必须能讲清）

- HTTP 请求组成：method、path、query、header、body
- 常见 header：
  - Content-Type：`application/json`
  - Authorization：`Bearer <token>`
  - Accept：期望返回格式
  - User-Agent
- 状态码（业务岗高频）：
  - 200 OK：成功
  - 201 Created：创建成功
  - 204 No Content：成功但无返回体
  - 400 Bad Request：参数错误
  - 401 Unauthorized：未登录/鉴权失败
  - 403 Forbidden：无权限
  - 404 Not Found：资源不存在
  - 409 Conflict：冲突（例如重复创建/幂等冲突）
  - 429 Too Many Requests：限流
  - 500 Internal Server Error：服务内部错误
  - 502/503：上游/服务不可用
  - 504 Gateway Timeout：网关超时

### 3.2 REST API 设计（能直接套到项目里）

资源命名：

- 用名词复数：`/users`、`/tasks`
- 子资源：`/users/{id}/tasks`

方法语义：

- GET：查询（应幂等）
- POST：创建/提交（通常非幂等；可通过幂等 key 实现幂等）
- PUT：整体更新（应幂等）
- PATCH：部分更新（通常也应幂等）
- DELETE：删除（应幂等）

统一返回结构（建议）：

- success：
  - `data`：对象/列表/分页结构
  - `request_id`
- error：
  - `code`：业务错误码（稳定、可前后端约定）
  - `message`：给用户看的信息
  - `request_id`

分页与排序（面试会追问）：

- offset 分页：
  - 参数：page/page_size 或 offset/limit
  - 缺点：深分页性能差（offset 大时扫描多）
- seek 分页（游标分页）：
  - 参数：cursor（例如 last_id/last_time）
  - 优点：性能稳定；缺点：实现更复杂、对排序有约束

过滤：

- query 参数：`/tasks?status=done&tag=work`
- 对字段范围：`created_at_gte=...`

### 3.3 鉴权（JWT 高概率问）

JWT 三段：

- Header：算法等元信息
- Payload：claims（用户信息、过期时间等）
- Signature：签名

常见 claims：

- sub：主体（用户 ID）
- exp：过期时间
- iat：签发时间
- iss：签发者

注意点（项目与面试都要提到）：

- JWT 适合无状态鉴权，但“注销/踢下线/权限变更即时生效”需要额外机制
- Token 存储：
  - 浏览器：更倾向 HttpOnly Cookie（降低 XSS 风险），但要处理 CSRF
  - App/服务端：Authorization Header 常见
- 过期与刷新：
  - access token 短期 + refresh token 长期（可选）
- 签名密钥保管：不要写死在代码里，使用环境变量/配置

### 3.4 中间件（日志、recover、鉴权、限流）

常见中间件职责：

- request-id：生成/透传，并写到响应 header
- 日志：记录 method/path/status/latency/request-id/user-id
- recover：捕获 panic，返回 500，避免进程崩溃
- 鉴权：解析 token、注入用户信息到 ctx
- 限流：按 IP/用户/接口维度控制 QPS

### 3.5 超时、重试与幂等（能讲“取舍”）

超时：

- 客户端超时：避免一直等待
- 服务端超时：下游依赖超时、DB/Redis 超时
- 原则：超时必须层层传递（context）

重试：

- 只对“可安全重试”的请求重试（通常是读请求或具备幂等保障的写请求）
- 避免雪崩：指数退避 + 抖动（概念即可）

幂等：

- 幂等定义：多次执行效果等价于一次
- 写请求幂等方案（项目加分项）：
  - 客户端传 idempotency-key
  - 服务端把 key 与结果绑定（DB 唯一键/Redis 记录）
  - 重复请求直接返回第一次结果或返回 409/200（按设计）

### 3.6 输入校验与安全（业务岗常见追问）

- 参数校验：
  - 类型、范围、长度、枚举值
  - 对字符串做 trim、对日期格式做解析校验
- SQL 注入：使用参数化查询/ORM（不要拼接 SQL）
- XSS/CSRF：主要是 Web 场景；如果用 Cookie 鉴权要考虑 CSRF
- 密码存储：
  - 不存明文
  - 使用强 hash（概念：bcrypt/argon2）；至少能说“不应该用 md5/sha1 直接 hash”

---

## 4. MySQL（表结构、索引、事务、慢 SQL）

### 4.1 表结构设计（面试很爱问）

基础原则：

- 先满足查询，再谈范式
- 核心查询要能走索引
- 字段类型选对（影响存储与索引效率）

字段类型建议：

- 主键：
  - 常见：自增 bigint / 雪花 id（概念即可）
  - 业务上如果有天然唯一键，可做唯一索引
- 字符串：
  - 统一使用 `utf8mb4`
  - 长文本用 TEXT（注意索引限制）
- 时间：
  - created_at/updated_at 常规字段
  - 业务时间单独字段（例如 due_at）
- 软删除：
  - `deleted_at` 或 `is_deleted`
  - 注意：软删除后查询要加过滤条件，否则索引设计要考虑该条件

### 4.2 索引（B+Tree）与常见索引类型

索引价值：

- 减少扫描行数（提高选择性）
- 利用有序性（范围查询、排序、group by）

高频概念：

- **聚簇索引**（InnoDB）：
  - 主键索引的叶子节点存整行数据
- **二级索引**：
  - 叶子节点存主键值，需要回表
- **覆盖索引**：
  - 查询字段都在索引上，不用回表
- **联合索引**与最左前缀：
  - 联合索引 `(a, b, c)` 能支持 `a`、`a+b`、`a+b+c` 的匹配
  - `b` 单独通常无法利用（除非发生索引下推/其他优化，但不要依赖）

什么情况下索引可能失效（面试最爱问）：

- 对索引列做函数/运算：`where date(created_at)=...`
- 隐式类型转换：字符串列用数字比较
- like 前缀通配：`like '%abc'`
- 条件范围后列无法充分利用（联合索引中范围条件之后）
- 选择性太差导致优化器放弃索引（例如性别字段）

### 4.3 EXPLAIN 必会字段（能解释你的项目查询）

重点关注：

- type：`ALL`（全表）> `index` > `range` > `ref` > `eq_ref` > `const/system`
- key：实际使用的索引
- rows：预估扫描行数
- Extra：
  - Using index：覆盖索引
  - Using where：需要回表/过滤
  - Using filesort：额外排序（可能慢）
  - Using temporary：临时表（可能慢）

复盘方法：

- 先看是否走对索引（key）
- 再看扫描行数（rows）
- 再看是否出现 filesort/temporary

### 4.4 事务与隔离级别（概念必须清楚）

ACID：

- 原子性、一致性、隔离性、持久性

隔离级别与现象：

- 读未提交：可能脏读
- 读已提交：避免脏读，可能不可重复读
- 可重复读（MySQL InnoDB 默认）：避免不可重复读；通过 MVCC + gap lock 等避免幻读（具体实现可被追问）
- 串行化：最强隔离，性能差

MVCC（能讲“是什么”即可）：

- 通过版本链 + Read View 实现快照读

锁与死锁：

- 行锁、间隙锁（gap lock）、临键锁（next-key lock）
- 死锁出现后 InnoDB 会回滚其中一个事务

项目落地建议：

- 写请求尽量短事务
- 统一按固定顺序更新资源，降低死锁概率

### 4.5 慢 SQL 排查思路（面试加分）

- 先确认是否命中索引：EXPLAIN
- 再看是否扫描行数过多：rows/实际数据分布
- 检查是否产生 filesort/temporary
- 检查查询字段是否需要回表：能否覆盖索引
- 检查分页方式：深分页是否需要 seek
- 检查连接池与并发：是否连接耗尽导致排队

---

## 5. Redis（缓存、热点、穿透/击穿/雪崩）

### 5.1 Redis 常用数据结构与场景

- string：缓存对象（序列化 JSON）、计数器（incr）
- hash：对象字段存储（适合部分更新）
- list：队列（简单场景）
- set：去重、集合关系
- zset：排行榜、按分数排序

### 5.2 缓存模式（必须能讲一种并落地）

Cache Aside（旁路缓存）最常用：

- 读：
  - 先查缓存，命中直接返回
  - 未命中查 DB，写回缓存（设置 TTL）
- 写：
  - 写 DB 成功后删除缓存（或更新缓存）

常见坑：

- “先删缓存再写 DB”会出现并发问题
- “先写 DB 再删缓存”也可能出现短暂不一致，但通常可接受

工程化建议：

- TTL 要有随机抖动（避免同一时刻大量过期）
- 缓存 key 设计要统一前缀：`app:resource:id`
- 大对象缓存要考虑体积与序列化开销

### 5.3 三大缓存问题（面试高频）

缓存穿透：

- 访问不存在的数据导致每次落到 DB
- 应对：
  - 缓存空值（短 TTL）
  - 布隆过滤器（概念即可）
  - 参数校验 + 频控

缓存击穿：

- 热点 key 过期瞬间，大量请求打到 DB
- 应对：
  - 热点 key 不过期 + 异步刷新
  - 互斥锁/单飞（singleflight）保证只有一个请求回源（概念即可）

缓存雪崩：

- 大量 key 集中过期或 Redis 故障，流量打爆 DB
- 应对：
  - TTL 加随机
  - 多级缓存/降级
  - 限流/熔断

### 5.4 分布式锁（可选加分项，能讲核心原则）

最小闭环：

- `SET key value NX PX ttl` 获取锁
- value 用随机值，释放锁前校验 value（避免误删）
- 释放锁用 Lua 脚本保证原子性（概念即可）

注意：

- 锁必须有过期时间（防死锁）
- 分布式锁不等于事务，不能保证强一致，只用于降低并发冲突

### 5.5 Redis 限流（项目加分项）

常见算法：

- 固定窗口：简单，但边界抖动明显
- 滑动窗口：更平滑
- 令牌桶：控制平均速率并允许突发

---

## 6. 工程化与项目结构（简历与面试非常吃这一块）

### 6.1 分层与职责边界（建议最少三层）

- handler/controller：
  - 解析参数、校验、组装 DTO
  - 调用 service
  - 统一返回结构
- service：
  - 业务逻辑、权限判断、事务边界
  - 组合 repo、cache、外部依赖
- repo/dao：
  - 数据持久化（SQL）
  - 尽量不掺杂业务规则

常问：

- 为什么要分层？
  - 可测试、可维护、降低耦合；面试可以结合你项目举例

### 6.2 配置管理

- 建议通过环境变量或配置文件区分：
  - dev/test/prod
  - DB/Redis 地址、超时、日志级别、JWT secret
- 原则：
  - 配置不写死
  - 敏感信息不进仓库

### 6.3 日志（必须能在项目里落地）

建议字段：

- request-id（最重要）
- method/path/status/latency
- user-id（如果鉴权后可获取）
- err（内部错误详情）

关键原则：

- 错误返回给用户“友好消息”，内部日志记录“可排障信息”
- 不要打印敏感信息（密码、token、密钥）

### 6.4 统一错误码与返回格式

建议设计：

- code：稳定枚举（例如 10001 参数错误、20001 未登录、30001 权限不足、50001 内部错误）
- message：用户可理解
- request_id：便于排障定位

面试表达要点：

- 为什么不用纯 HTTP 状态码？
  - 状态码粒度不够，业务需要稳定的细分错误码

### 6.5 测试（至少 3-5 个单测）

单测优先级：

- service 层逻辑（最值钱）
- 工具函数（校验、分页计算、错误封装）
- handler 可用 httptest 做基础测试（可选）

表驱动测试（Go 常用）：

- 用切片定义多组输入输出
- 避免重复代码，覆盖边界条件

建议最低验收：

- `go test ./...` 全绿
- `go test -race ./...` 能跑通（即使慢也建议在 CI 跑）

---

## 7. Linux 基础（第 3 周重点，面试常问排障）

### 7.1 进程与资源

- `ps`：看进程
- `top`：看 CPU/内存/负载
- `kill -TERM/-KILL`：信号与强杀
- 文件描述符：网络连接、文件、日志都依赖 FD；FD 泄漏会导致服务异常

### 7.2 网络排障

- `ss -lntp`：看监听端口与进程
- `lsof -i :port`：定位端口占用
- `curl`：验证 HTTP 接口
- `ping`/`traceroute`（概念即可）：网络可达性

### 7.3 日志定位

- `tail -f`：实时查看
- `grep`：过滤关键字（结合 request-id）
- 先定位时间窗口，再定位 request-id，再看错误堆栈与上下游依赖

---

## 8. Docker 与部署（第 3 周重点）

### 8.1 Docker 基本概念

- 镜像（image）：只读模板
- 容器（container）：镜像运行实例
- 层（layer）：Dockerfile 每条指令生成新层，影响构建缓存

### 8.2 Dockerfile 常见要点

- 使用多阶段构建（概念即可）：编译阶段 + 运行阶段，减小镜像体积
- 固定端口、用环境变量注入配置
- 不把敏感配置写进镜像

### 8.3 docker-compose 常见要点

- services：app/mysql/redis
- networks：服务互通
- volumes：MySQL/Redis 数据持久化
- env：配置注入

排障思路：

- `docker logs <container>` 看日志
- `docker exec -it <container> sh` 进入容器检查
- 检查端口映射、网络连通、依赖服务是否 ready

---

## 9. 算法与手写题（冲刺期常考题型模板）

目标：

- 会写常见题型模板
- 能清楚讲解：思路、时间复杂度、边界条件

### 9.1 时间/空间复杂度基本功

- O(1)、O(log n)、O(n)、O(n log n)、O(n^2)
- 常见优化方向：
  - 用哈希换时间
  - 双指针减少嵌套循环
  - 二分减少搜索空间
  - 单调栈/队列处理“最近更大/更小”

### 9.2 高频题型与模板清单

- 数组/字符串：
  - 双指针：去重、反转、两数之和变体
  - 滑动窗口：最长不重复子串、最小覆盖子串（至少掌握一种窗口写法）
- 哈希：
  - 两数之和、异位词分组、频次统计
- 二分：
  - 查找边界（lower_bound/upper_bound 思想）
  - 有序数组旋转查找（进阶）
- 栈/队列：
  - 括号匹配
  - 单调栈：下一个更大元素、每日温度
- 链表（看岗位要求，可选）：
  - 反转链表、找环、合并有序链表
- DFS/BFS（看岗位要求，可选）：
  - 树遍历、最短路径（BFS）

面试表达模板：

- 先讲暴力解法与复杂度
- 再讲优化思路与关键数据结构
- 最后讲边界：空输入、重复元素、溢出

---

## 10. 高频问答清单（面试前突击）

### 10.1 Go

- slice 扩容发生了什么？什么时候会影响原切片？
- map 为什么并发不安全？你项目里怎么处理并发访问？
- channel 关闭后读/写分别会怎样？如何优雅退出 worker？
- context 用在什么地方？你项目里如何设置超时与取消？
- goroutine 泄漏有哪些常见原因？你如何避免？

### 10.2 HTTP/工程

- REST 怎么设计资源与方法？状态码怎么选？
- 幂等是什么？写请求如何做幂等？
- 中间件通常做哪些事？你项目里有哪些中间件？
- 超时与重试应该放在哪里？如何避免重试放大故障？

### 10.3 MySQL/Redis

- 索引为什么快？联合索引最左前缀是什么？哪些情况索引会失效？
- 事务隔离级别与常见读现象是什么？
- 缓存穿透/击穿/雪崩分别是什么？各怎么解决？
- 你项目里的缓存策略是什么？如何处理一致性与过期？

---

## 11. 把知识点映射到“项目可讲点”（建议你最终准备 3 个）

你最终要做到：每个可讲点 1-2 分钟讲清楚，并能结合自己项目代码与数据佐证。

可选方向（任选 3 个打磨）：

- 缓存策略与一致性：为什么用 cache-aside，TTL 怎么选，如何避免击穿
- 索引与慢查询优化：查询场景是什么，索引怎么建，EXPLAIN 前后对比
- 请求链路可观测性：request-id 怎么生成与传递，日志字段怎么设计，如何定位问题
- 超时与资源保护：ctx 超时怎么设，限流/熔断怎么做（最小闭环）
- 幂等：为什么要幂等，幂等 key 怎么设计，失败与重复请求怎么处理

