# Golang 求职详细知识点大纲

> 对应文档：`golang-job-12week-plan.md`  
> 用途：按阶段查缺补漏；面试前按章节过一遍「会讲 / 会写 / 会排障」。  
> 标记说明：
> - ★ 必会（多数岗位）
> - ◆ 加分（中高级 / 云原生）
> - ○ 了解即可（先混个脸熟）

---

## 目录

1. [语言与工程基础](#1-语言与工程基础第-1–2-周)
2. [Web 与业务服务](#2-web-与业务服务第-3–4-周)
3. [工程化与可靠性](#3-工程化与可靠性第-5–6-周)
4. [分布式与异步](#4-分布式与异步第-7–8-周)
5. [云原生与差异化](#5-云原生与差异化第-9–10-周)
6. [面试计算机基础](#6-面试计算机基础第-11–12-周贯穿)
7. [算法与数据结构](#7-算法与数据结构贯穿)
8. [项目与简历话术知识点](#8-项目与简历话术知识点)

---

## 1. 语言与工程基础（第 1–2 周）

### 1.1 语法与类型系统 ★

- 变量、常量、iota、类型推断
- 基本类型：整型、浮点、布尔、字符串、rune/byte
- 复合类型：数组、切片、map、结构体
- 切片底层：数组指针、len、cap；扩容规则；共享底层数组的坑
- map：无序、非并发安全、`make` 与 nil map 区别
- 指针：值拷贝 vs 指针传递；何时用指针接收者
- 结构体嵌套、匿名字段、标签（`json:"xxx"`）
- 方法集：值接收者 / 指针接收者对接口实现的影响
- 接口：隐式实现、空接口、类型断言、类型 switch
- 泛型基础：类型参数、约束（`comparable` 等）、常见用法边界
- `defer`：执行顺序、参数求值时机、与返回值/named return 的关系
- `panic` / `recover`：适用场景；为何业务错误不应用 panic

### 1.2 错误处理 ★

- `error` 接口与自定义错误
- `errors.New` / `fmt.Errorf`
- `errors.Is` / `errors.As` / `errors.Unwrap`
- 错误包装：`%w` 与错误链
- 哨兵错误 vs 错误类型
- 何时返回错误、何时记录日志（避免重复打日志）
- 业务错误码设计（与 HTTP 状态码映射）

### 1.3 context ★

- `context.Background` / `TODO` / `WithCancel` / `WithTimeout` / `WithDeadline` / `WithValue`
- 超时与取消如何向下传递
- 在 HTTP、DB、RPC 中如何贯穿 context
- `WithValue` 的使用边界（不放可选参数、不放大对象）
- 常见坑：context 泄漏、在已取消后继续干活、忽略 `ctx.Done()`

### 1.4 并发模型 ★

- goroutine 启动与生命周期
- channel：无缓冲 / 有缓冲；发送接收语义；关闭规则；从已关闭 channel 读
- `select`：多路复用、default、超时写法
- 常见模式：
  - worker pool
  - fan-in / fan-out
  - 信号通知（done channel）
  - 限流（buffered channel 当信号量）
- `sync` 包：
  - `Mutex` / `RWMutex`
  - `WaitGroup`
  - `Once`
  - `Map`（适用场景与代价）
  - `Cond`（了解）
  - `Pool`（对象复用思路）
- `atomic`：计数器、标志位；与 mutex 的选型
- 竞态：什么是 data race；`go test -race` 怎么用
- 死锁、活锁、饥饿（能举例）
- goroutine 泄漏常见原因与排查思路

### 1.5 内存与运行时（面试够用）★ / ◆

- 逃逸分析：什么情况会分配到堆（概念级）
- 栈 vs 堆（口述即可）
- GC 粗略：三色标记、STW 概念、如何减少分配
- GMP 调度模型：G/M/P 职责、工作窃取（面试常问）
- channel 实现要点（队列 + sudog，概念级）

### 1.6 标准库必会 ★

#### net/http

- `http.Handler` / `HandlerFunc` / `ServeMux`
- 路由、方法匹配、路径参数（标准库限制与自实现思路）
- 中间件链式写法
- `Request` / `ResponseWriter` 常用字段
- 客户端：`http.Client`、超时、Transport
- Cookie、Header、Query、Form、JSON Body 读写

#### encoding/json

- 序列化 / 反序列化
- 忽略字段、omitempty、自定义时间格式
- `json.RawMessage`、流式 Decoder/Encoder
- 数字精度问题（`json.Number`）

#### io / bufio / bytes / strings

- `Reader` / `Writer` / `Closer`
- `io.Copy`、`io.ReadAll`、`LimitReader`
- 缓冲读写适用场景

#### time

- `Duration`、`Time`、时区、定时器、`Ticker`
- 单调时钟与墙钟差异（概念）

#### testing

- 表驱动测试
- `t.Run` 子测试
- `Helper`、`Cleanup`
- benchmark：`testing.B`
- 测试覆盖率：`go test -cover`
- 测试中的并发与 `-race`

### 1.7 工程与工具链 ★

- `go mod`：init、tidy、vendor（是否需要）、replace
- 版本选择与语义化版本
- 项目布局常见做法：`cmd/`、`internal/`、`pkg/`（理解即可，不教条）
- `go fmt` / `go vet` / `staticcheck`（了解）
- 构建：`go build`、交叉编译、`-ldflags`
- 调试：日志、Delve 基础（○）

---

## 2. Web 与业务服务（第 3–4 周）

### 2.1 HTTP 与 API 设计 ★

- REST 常见约定：资源命名、幂等语义（GET/PUT/DELETE）
- 状态码选用：2xx / 4xx / 5xx
- 统一响应结构：code / message / data
- 分页：offset/limit vs cursor
- 过滤、排序、字段选择
- 幂等键（Idempotency-Key）思路
- API 版本管理（URL / Header）

### 2.2 Web 框架（Gin 为例，会一个即可）★

- 路由分组、路径参数、查询参数
- 中间件：日志、恢复、鉴权、CORS
- 参数绑定与校验（binding / validator）
- 文件上传下载
- 优雅退出与框架集成方式
- 框架 vs 标准库：取舍

### 2.3 鉴权与安全基础 ★

- 认证 vs 授权
- Session / Cookie vs JWT / Token
- JWT：结构（header/payload/signature）、过期、刷新令牌思路
- 密码存储：哈希 + salt（bcrypt/argon2 概念）
- 常见安全问题：
  - SQL 注入
  - XSS / CSRF（后端视角）
  - 越权（水平/垂直）
  - 敏感信息进日志
- HTTPS、TLS 基本概念
- 密钥与配置不进仓库

### 2.4 配置、日志、依赖注入 ★

- 配置来源：文件、环境变量、命令行
- 配置热更新（○）
- 结构化日志：字段、级别、trace id
- 日志库选型思路（slog / zap / logrus）
- 依赖组装：手工注入 vs wire（了解）

### 2.5 MySQL / SQL ★

#### SQL 基础

- DDL / DML / 事务语句
- JOIN、聚合、子查询、UNION
- 索引类型：主键、唯一、普通、联合索引
- 最左前缀原则
- 覆盖索引、回表
- EXPLAIN 关键字段：type、key、rows、Extra
- 慢查询排查步骤

#### 事务与隔离

- ACID
- 隔离级别：读未提交 / 读已提交 / 可重复读 / 串行化
- 脏读、不可重复读、幻读
- InnoDB 行锁、间隙锁（概念）
- 事务中长事务的危害

#### 表设计

- 范式与反范式取舍
- 主键选型（自增 / 雪花）
- 软删除、时间字段、状态机字段
- 字符集与排序规则（utf8mb4）

#### Go 访问层

- `database/sql`：连接池参数（`SetMaxOpenConns` 等）
- prepared statement
- 事务写法：`Begin` / `Commit` / `Rollback`
- GORM（或同类 ORM）：
  - CRUD、关联、钩子
  - 预加载 N+1 问题
  - 迁移（○）
- SQL 与 ORM 如何选型

### 2.6 Redis ★

- 常用数据结构：String、Hash、List、Set、ZSet
- 过期策略、内存淘汰粗略概念
- 缓存模式：Cache Aside（旁路缓存）
- 缓存问题：
  - 穿透
  - 击穿
  - 雪崩
  - 数据不一致
- 分布式锁：SET NX EX + 持有者校验；红锁仅了解
- 计数器、限流简单实现
- Pipeline / 事务 / Lua（○～◆）
- Go 客户端基本用法与连接池

### 2.7 Docker 入门（项目可运行）★

- 镜像、容器、仓库
- Dockerfile 常用指令
- `docker compose`：多服务编排（app + mysql + redis）
- 卷、网络、端口映射
- 多阶段构建（第 9–10 周加深）

---

## 3. 工程化与可靠性（第 5–6 周）

### 3.1 超时、取消、重试 ★

- 各级超时：客户端 / 网关 / 服务 / DB
- 重试条件：只对幂等或可安全重试的错误重试
- 退避：固定、指数退避、抖动（jitter）
- 重试风暴与舱壁（bulkhead）概念
- context 取消后的资源清理

### 3.2 幂等、限流、熔断、降级 ★ / ◆

- 幂等实现：唯一键、状态机、Token
- 限流算法：固定窗口、滑动窗口、漏桶、令牌桶（会讲会写一种）
- 熔断：关闭 / 半开 / 打开
- 降级：返回缓存、默认值、裁剪功能
- 负载保护：排队、拒绝、过载

### 3.3 服务生命周期 ★

- 优雅启动：依赖就绪再接流量
- 优雅关闭：停监听 → 等在途请求 → 关连接
- 健康检查：liveness vs readiness
- 信号处理：`SIGINT` / `SIGTERM`

### 3.4 测试策略 ★

- 单元测试：纯逻辑、mock 边界
- 表驱动 + 黄金文件（○）
- 集成测试：testcontainers / compose（○～◆）
- 接口契约测试思路（○）
- 哪些必须测：鉴权、幂等、超时、边界输入

### 3.5 性能与排查 ★ / ◆

- `pprof`：CPU、heap、goroutine、block、mutex
- `go tool pprof` / 火焰图基本阅读
- benchmark 与实战优化的关系
- 常见瓶颈：锁竞争、盲目序列化、N+1 查询、无界并发
- 指标：QPS、延迟分位（P50/P95/P99）、错误率、饱和度

### 3.6 可观测性入门 ★ / ◆

- 日志、指标、链路追踪三者关系
- trace id / span 概念
- OpenTelemetry 粗略了解（○～◆）
- Prometheus 指标类型：Counter、Gauge、Histogram（◆）

---

## 4. 分布式与异步（第 7–8 周）

### 4.1 gRPC 与 Protobuf ★ / ◆

- IDL：message、service、rpc
- 生成代码流程
- 四种 RPC：Unary / Server stream / Client stream / Bidi（先会 Unary）
- 错误模型：status code
- metadata、拦截器（interceptor）
- 超时、取消与 context
- 兼容性：字段号、不要复用字段号
- HTTP/JSON 网关（○）
- 对比 REST：性能、契约、浏览器友好度

### 4.2 消息队列 ★ / ◆

任选 Kafka 或 RabbitMQ 深挖一种，另一种了解差异即可。

#### 共性知识点

- 为什么用 MQ：解耦、削峰、异步
- 投递语义：至多一次 / 至少一次 / 恰好一次（现实中如何近似）
- 生产者确认、消费者确认 / ack
- 重复消费与幂等
- 顺序消息：分区键 / 单队列限制
- 死信队列、重试队列
- 消息堆积排查
- 事务消息 / 本地消息表 / Outbox 模式（◆）

#### Kafka 要点（若选 Kafka）

- Topic、Partition、Consumer Group
- offset 管理
- ISR、副本（概念）
- 再均衡（rebalance）影响

#### RabbitMQ 要点（若选 RabbitMQ）

- Exchange、Queue、Binding、Routing Key
- 工作队列、发布订阅
- ACK / NACK / requeue
- 延迟消息实现思路（○）

### 4.3 分布式常见问题 ★ / ◆

- 服务间调用失败模式
- 分布式事务：2PC（了解）、Saga、TCC（概念对比）
- 一致性：强一致 vs 最终一致
- CAP / BASE（面试口述）
- 分布式 ID：UUID、雪花、号段
- 分布式锁应用边界（别当银弹）
- 会话与粘滞、无状态服务设计

### 4.4 配置与服务发现（了解）○ / ◆

- 注册发现：Consul / Nacos / etcd 概念
- 配置中心用途
- etcd 在 K8s 中的角色（为后面铺垫）

---

## 5. 云原生与差异化（第 9–10 周）

### 5.1 Docker 进阶 ★

- 多阶段构建、减小镜像
- 基础镜像选型（distroless/alpine 取舍）
- 非 root 用户运行
- 健康检查指令
- 构建缓存与 `.dockerignore`
- 镜像安全扫面概念（○）

### 5.2 Kubernetes 基础 ★ / ◆

- Pod、ReplicaSet、Deployment、Service
- ConfigMap、Secret
- Ingress 概念
- 探针：liveness / readiness / startup
- 资源 request / limit
- 滚动更新与回滚
- 命名空间、标签、选择器
- 常见排障：CrashLoopBackOff、ImagePullBackOff、未就绪
- ServiceAccount / RBAC（○～◆）

### 5.3 IaC / Terraform Provider（差异化加分）◆

结合现有 Provider 开发经验重点整理：

- Resource 生命周期：Create / Read / Update / Delete / Import
- 状态（state）与云上真实资源的对齐
- 幂等：重复 apply 的结果
- 重试：限流、暂时性错误、最终一致性窗口
- 部分失败与回滚策略
- schema / 参数校验 / ForceNew
- 数据源（Data Source）与资源差异
- API 分页、异步任务（job）等待
- 鉴权：AK/SK、项目、区域（region）
- 可测试性：acceptance test 思路（○）

### 5.4 Operator / 控制器（可选）○ / ◆

- 控制回路：reconcile
- CRD 概念
- 最终一致与 requeue
- client-go / controller-runtime 粗略地图（不必深挖完）

### 5.5 Linux 与网络排障基础 ★

- 常用命令：`ps`、`top`、`lsof`、`netstat`/`ss`、`curl`、`dig`/`nslookup`
- 文件描述符、连接数
- DNS 解析失败排查
- 抓包概念（tcpdump，○）
- 日志定位：应用日志 + 系统日志

---

## 6. 面试计算机基础（第 11–12 周，贯穿）

### 6.1 操作系统 ★

- 进程 vs 线程 vs 协程
- 用户态 / 内核态
- 上下文切换
- 内存管理：虚拟内存、缺页（概念）
- IO：阻塞 / 非阻塞、同步 / 异步、多路复用（select/poll/epoll 概念）
- 零拷贝粗略了解（○）

### 6.2 计算机网络 ★

- OSI / TCP-IP 分层（能对应常见协议）
- TCP：三次握手、四次挥手、状态机关键状态
- 可靠传输：序号、ACK、重传、滑动窗口、拥塞控制（口述）
- UDP 对比
- HTTP/1.1 vs HTTP/2 粗略差异
- HTTPS：证书、握手概览、对称/非对称加密角色
- WebSocket 概念（○）
- 常见头：Host、Content-Type、Authorization、Connection
- 代理、反向代理、负载均衡类型（四层/七层）

### 6.3 数据库原理补充 ★

- B+ Tree 索引结构直觉
- 聚簇索引 vs 非聚簇
- MVCC 粗略
- 主从复制与读写分离（○～◆）
- 分库分表思路与代价（◆）

### 6.4 Go 面试高频专题 ★

- 切片扩容与复制
- map 非线程安全
- defer 经典题
- interface 底层（类型信息 + 数据指针，概念）
- channel 关闭与广播
- 调度：syscall、网络轮询器（概念）
- GC 会不会停顿、如何优化分配
- 为什么有 goroutine 仍可能阻塞进程（线程耗尽等）
- 如何定位泄漏：pprof goroutine

---

## 7. 算法与数据结构（贯穿）

### 7.1 必刷题型 ★

- 数组 / 双指针 / 滑动窗口
- 哈希表
- 链表（反转、环、合并）
- 栈与队列
- 二叉树遍历、层序、BST 基础
- 二分查找
- 排序：快排 / 归并思路；稳定性概念
- 堆 / TopK（◆）
- 回溯与 DFS/BFS 基础（◆）
- 动态规划入门：爬楼梯、背包感知（◆，量力）

### 7.2 复杂度 ★

- 时间 / 空间复杂度口述
- 均摊分析直觉（切片扩容）

### 7.3 每周节奏建议

- 每周 4–6 题，优先「会讲思路」而不是刷数量
- 同一题型连续练到能无提示写出来

---

## 8. 项目与简历话术知识点

### 8.1 主项目建议覆盖的能力清单 ★

- 用户注册登录与鉴权
- 业务 CRUD + 分页过滤
- MySQL 事务或状态流转
- Redis 缓存至少一个读路径
- 统一错误码与请求日志（含 request id）
- 超时与优雅关闭
- Docker Compose 一键启动
- （进阶）MQ 异步链路
- （进阶）gRPC 内部调用
- （进阶）K8s 部署清单

### 8.2 面试项目深挖准备 ★

每个项目至少准备：

1. **背景与目标**：解决什么问题
2. **架构图**：客户端 → API → DB/Cache/MQ
3. **你负责的模块**与技术选型原因
4. **一个难点**：现象 → 分析 → 方案 → 结果
5. **一个故障或风险**：如何发现与止血
6. **数据指标**：QPS、延迟、错误率、成本（能量化就量化）
7. **如果重做**：会改什么

### 8.3 Provider / 云 API 经历可讲点 ◆

- 一次 apply 背后的 API 调用链
- 如何处理异步创建（轮询 job / 等状态）
- 限流与重试策略
- Read 对账如何避免漂移
- 删除保护、ForceNew 字段变更
- 多 region / 项目隔离踩坑

---

## 9. 按优先级的学习顺序（总览）

```text
P0 必会（先找工作）
  Go 语法并发 + context + 测试
  HTTP/API + MySQL + Redis
  Docker Compose 可运行项目
  操作系统/网络/DB 八股基础
  算法热题

P1 强烈加分
  超时重试幂等限流
  gRPC
  一种 MQ
  pprof 与基本排障
  K8s 基础部署

P2 差异化
  Terraform Provider / IaC 深挖
  可观测（OTel/Prometheus）
  Operator / 云原生控制面
  分库分表、复杂分布式事务
```

---

## 10. 自检表（可打印）

### 语言

- [ ] 能手写带中间件的标准库 HTTP 服务
- [ ] 能解释 channel 与 mutex 选型
- [ ] 会用 context 做超时取消
- [ ] 会写表驱动测试并跑 `-race`

### 业务服务

- [ ] 会做 JWT（或 Session）鉴权
- [ ] 会写 SQL 并看懂基础 EXPLAIN
- [ ] 会用 Redis 做缓存并讲清穿透/击穿/雪崩
- [ ] 项目能 `docker compose up`

### 可靠性

- [ ] 服务能优雅关闭
- [ ] 关键路径有超时
- [ ] 能说明一个幂等方案
- [ ] 用过 pprof 看过至少一种 profile

### 分布式

- [ ] 写过 gRPC Unary
- [ ] 写过 MQ 生产消费并处理重复消费
- [ ] 能画同步 + 异步架构图

### 云原生 / 差异化

- [ ] 能解释 Deployment + Service
- [ ] 能讲清一个 Resource 的 CRUD 与重试
- [ ] （可选）本地 K8s 跑通过项目

### 面试

- [ ] 项目 3 分钟 / 10 分钟两套讲法
- [ ] 两个难点故事
- [ ] 网络 + MySQL + Go 调度常见题能展开

---

## 附录 A：建议最小工具集

- Go 工具链、Git、Docker Desktop / Docker Engine
- Postman 或 curl
- MySQL 客户端（或 GUI）
- Redis CLI
- IDE：GoLand 或 VS Code + Go 插件
- （可选）kind / minikube、protoc

## 附录 B：文档维护建议

- 每学完一小节，在对应 checkbox 打勾
- 不会讲的知识点单独建「薄弱清单」
- 面试前只复习：薄弱清单 + 第 8 章项目话术 + 第 6 章高频题
