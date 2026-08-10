# Golang 求职知识点详解（含示例）

> 对应大纲：`golang-knowledge-points.md`  
> 对应计划：`golang-job-12week-plan.md`  
> 本文目标：说明「是什么、为什么、怎么用、容易错在哪、面试怎么讲」，并在关键点给出**可直接对照敲的示例**。

---

## 使用方式

- 按周学习：读说明 → 自己敲示例 → 用章末「自问」检查。
- 面试前：优先复习标 ★ 的段落与示例旁的「面试一句话」。
- 标记：★ 必会｜◆ 加分｜○ 了解。
- 示例以精简片段为主，省略 `package`/`import` 时请自行补全。

---

# 第 1 章 语言与工程基础

## 1.1 语法与类型系统 ★

### 变量、常量与 iota

Go 用 `var` / `:=` 声明变量。函数内常用 `:=`，包级变量用 `var`。常量用 `const`，编译期确定。

```go
const (
    StatusDraft = iota // 0
    StatusOpen         // 1
    StatusDone         // 2
)

// 位运算枚举
const (
    Read  = 1 << iota // 1
    Write             // 2
    Exec              // 4
)
```

**要点**：iota 只在同一个 `const` 括号块内递增。

### 切片共享底层数组（易错）★

```go
a := []int{1, 2, 3, 4}
b := a[1:3]   // b == [2 3]，与 a 共享底层
b[0] = 99
fmt.Println(a) // [1 99 3 4] —— 改 b 影响了 a

// 需要独立副本
c := append([]int(nil), a...)
// 或
d := make([]int, len(a))
copy(d, a)
```

`append` 可能触发扩容，**必须接住返回值**：

```go
s := []int{1, 2}
s = append(s, 3) // 正确
```

nil 与空切片 JSON 差异：

```go
var nilS []int          // JSON: null
empty := make([]int, 0) // JSON: []
```

### map：nil 可读不可写 ★

```go
var m map[string]int
fmt.Println(m["a"]) // 0，不 panic
// m["a"] = 1       // panic: assignment to entry in nil map

m = make(map[string]int)
m["a"] = 1 // OK

// var：仅声明变量，但值为nil
// make：声明并未变量初始化且分配内存，值为空对象
```

### 指针接收者与接口 ★

```go
type Counter struct{ n int }

func (c *Counter) Inc() { c.n++ } // 指针接收者

type Incer interface{ Inc() }

func demo() {
    var c Counter
    var i Incer = &c // 必须取地址：*Counter 才实现 Inc
    // var i Incer = c // 编译错误
    i.Inc()
}
```

### 接口 == nil 的坑 ★

```go
func returnsErr() error {
    var p *os.PathError = nil
    return p // 返回的 error 接口「动态类型」非 nil
}

err := returnsErr()
fmt.Println(err == nil) // false！，因为这个err的类型不为nil,只有当值和类型都为nil时才是nil； 面试高频

// 正确：直接 return nil
func returnsErrOK() error {
    return nil
}
```

**面试一句话**：接口值是否为 nil，要看类型和值是否都空；把 typed nil 赋给接口后，接口本身不是 nil。

### defer 参数求值时机 ★

```go
func f() {
    x := 1
    defer fmt.Println("A", x) // 注册时求值，打印 1
    defer func() {
        fmt.Println("B", x) // 闭包读变量，打印 2
    }()
    x = 2
}
// 输出顺序：B 2，然后 A 1（LIFO）
```

### recover 只在 defer 中有效

```go
func safe(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                http.Error(w, "internal error", 500)
            }
        }()
        h.ServeHTTP(w, r)
    })
}
```

---

## 1.2 错误处理 ★

```go
var ErrNotFound = errors.New("not found")

type HTTPError struct {
    Code int
    Msg  string
}

func (e *HTTPError) Error() string { return e.Msg }

func GetUser(id string) error {
    if id == "" {
        return fmt.Errorf("get user: %w", ErrNotFound)
    }
    return nil
}

func handle(err error) {
    if errors.Is(err, ErrNotFound) {
        // 404
    }
    var he *HTTPError
    if errors.As(err, &he) {
        // 用 he.Code
    }
}
```

**日志原则**：中间层只 `%w` 包装返回；在 HTTP 边界记一次日志。

---

## 1.3 context ★

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel() // 必须调用，防泄漏

    if err := queryDB(ctx); err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            http.Error(w, "timeout", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, err.Error(), 500)
    }
}

func queryDB(ctx context.Context) error {
    select {
    case <-time.After(3 * time.Second): // 模拟慢查询
        return nil
    case <-ctx.Done():
        return ctx.Err() // 被取消/超时立刻返回
    }
}
```

带 request id（仅少量元数据）：

```go
type ctxKey string
const keyReqID ctxKey = "req_id"

ctx = context.WithValue(ctx, keyReqID, "abc-123")
id, _ := ctx.Value(keyReqID).(string)
```

**易错**：`WithCancel`/`WithTimeout` 派生后忘记 `defer cancel()`。

---

## 1.4 并发模型 ★

### Worker Pool

```go
func workerPool(tasks []int, n int) {
    ch := make(chan int)
    var wg sync.WaitGroup
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for t := range ch {
                fmt.Println("work", t)
            }
        }()
    }
    for _, t := range tasks {
        ch <- t
    }
    close(ch) // 发送方关闭
    wg.Wait()
}
```

### 信号量限并发

```go
sem := make(chan struct{}, 8) // 最多 8 个并发
var wg sync.WaitGroup
for _, u := range urls {
    wg.Add(1)
    sem <- struct{}{}
    go func(u string) {
        defer wg.Done()
        defer func() { <-sem }()
        http.Get(u)
    }(u)
}
wg.Wait()
```

### select + 超时

```go
select {
case v := <-ch:
    fmt.Println(v)
case <-time.After(time.Second):
    fmt.Println("timeout")
case <-ctx.Done():
    return ctx.Err()
}
```

### Mutex 保护共享 map

```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (s *SafeMap) Get(k string) (int, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.m[k]
    return v, ok
}

func (s *SafeMap) Set(k string, v int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[k] = v
}
```

### 竞态示例（应用 -race 抓住）

```go
// 错误：无锁并发写
n := 0
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() { defer wg.Done(); n++ }()
}
wg.Wait()
// go test -race 会报 data race
```

**面试一句话**：通信选型 channel，保护已有共享状态用 mutex；无界 `go` 要限流，否则易泄漏/打爆下游。

---

## 1.5 内存与运行时（面试够用）★／◆

查看是否逃逸：

```bash
go build -gcflags="-m" .
```

**GMP 口述**：G=协程，M=线程，P=调度上下文；G 挂在 P 本地队列，空闲 P 会工作窃取。

---

## 1.6 标准库要点 ★

### 中间件链式（标准库）

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

mux := http.NewServeMux()
mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ok"))
})
http.ListenAndServe(":8080", logging(mux))
```

### JSON 与大整数

```go
type Order struct {
    ID   int64  `json:"id,string"` // 也可用 string 防精度丢
    Name string `json:"name,omitempty"`
}

// decode 到 map[string]any 时数字默认 float64，慎用
```

### 限制 Body 大小

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
```

### 表驱动测试

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"ok", 1, 2, 3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.a + tt.b; got != tt.want {
                t.Fatalf("got %d want %d", got, tt.want)
            }
        })
    }
}
```

---

### 第 1 章自问

1. 为什么对 nil map 写入会 panic，读取却不会？
2. 接口变量什么时候看起来「里面是 nil」但 `== nil` 为 false？
3. 如何用 ctx 让进行中的 HTTP handler 尽快退出？

---

# 第 2 章 Web 与业务服务

## 2.1 HTTP 与 API 设计 ★

统一响应 + 幂等键思路：

```go
type Resp struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

// POST /orders
// Header: Idempotency-Key: client-uuid
// 服务端：若 key 已存在，直接返回首次响应，不重复创建
```

游标分页（示意）：

```sql
SELECT * FROM tasks
WHERE user_id = ? AND id > ?
ORDER BY id ASC
LIMIT 20;
```

---

## 2.2 Web 框架（Gin）★

```go
r := gin.New()
r.Use(gin.Recovery(), gin.Logger())

api := r.Group("/api/v1")
api.POST("/login", login)

auth := api.Group("/")
auth.Use(JWTAuth())
auth.GET("/tasks", listTasks)

type CreateReq struct {
    Title string `json:"title" binding:"required,min=1,max=100"`
}
func create(c *gin.Context) {
    var req CreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"message": err.Error()})
        return
    }
    c.JSON(200, gin.H{"title": req.Title})
}
```

---

## 2.3 鉴权与安全 ★

### 密码哈希（bcrypt）

```go
hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
err = bcrypt.CompareHashAndPassword(hash, []byte(plain))
```

### JWT 签发与校验（示意）

```go
func signToken(userID string, secret []byte) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(15 * time.Minute).Unix(),
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return t.SignedString(secret)
}
```

### SQL 注入：参数化 vs 错误拼串

```go
// 错误
db.Query("SELECT * FROM users WHERE name = '" + name + "'")

// 正确
db.QueryContext(ctx, "SELECT id, name FROM users WHERE name = ?", name)
```

### 防越权

```go
uid := currentUserID(ctx) // 来自 token，不是来自 query
row := db.QueryRowContext(ctx,
    "SELECT id, title FROM tasks WHERE id = ? AND user_id = ?",
    taskID, uid)
```

---

## 2.4 配置与日志 ★

```go
// 环境变量优先
dsn := os.Getenv("MYSQL_DSN")

log.Info("create user",
    "request_id", reqID,
    "user_id", uid,
    "err", err,
)
```

---

## 2.5 MySQL ★

### 联合索引与最左前缀

```sql
-- 索引 (user_id, status, created_at)
EXPLAIN SELECT * FROM tasks WHERE user_id = 1 AND status = 2;
-- 通常能用到索引

EXPLAIN SELECT * FROM tasks WHERE status = 2;
-- 跳过 user_id，往往用不上该联合索引
```

### Go 事务

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer func() {
    if err != nil {
        _ = tx.Rollback()
    }
}()

if _, err = tx.ExecContext(ctx, `UPDATE accounts SET bal = bal - ? WHERE id = ?`, 100, from); err != nil {
    return err
}
if _, err = tx.ExecContext(ctx, `UPDATE accounts SET bal = bal + ? WHERE id = ?`, 100, to); err != nil {
    return err
}
err = tx.Commit()
return err
```

### 连接池

```go
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(30 * time.Minute)
```

### N+1 问题（GORM）

```go
// 坏：循环里查
for _, u := range users {
    db.Where("user_id = ?", u.ID).Find(&orders)
}

// 好：预加载
db.Preload("Orders").Find(&users)
```

---

## 2.6 Redis ★

### Cache Aside

```go
func GetUser(ctx context.Context, id string) (*User, error) {
    key := "user:" + id
    if b, err := rdb.Get(ctx, key).Bytes(); err == nil {
        var u User
        _ = json.Unmarshal(b, &u)
        return &u, nil
    }
    u, err := loadFromDB(ctx, id)
    if err != nil {
        return nil, err
    }
    if u == nil {
        // 防穿透：缓存空值短 TTL
        _ = rdb.Set(ctx, key, "null", 30*time.Second).Err()
        return nil, ErrNotFound
    }
    b, _ := json.Marshal(u)
    // 过期加抖动，防雪崩
    ttl := 10*time.Minute + time.Duration(rand.Intn(60))*time.Second
    _ = rdb.Set(ctx, key, b, ttl).Err()
    return u, nil
}

func UpdateUser(ctx context.Context, u *User) error {
    if err := saveDB(ctx, u); err != nil {
        return err
    }
    return rdb.Del(ctx, "user:"+u.ID).Err() // 先 DB 再删缓存
}
```

### 分布式锁

```go
token := uuid.NewString()
ok, err := rdb.SetNX(ctx, "lock:order:"+id, token, 10*time.Second).Result()
if !ok {
    return errors.New("busy")
}
defer unlock(ctx, "lock:order:"+id, token)

// 释放：校验 token，避免误删他人锁（Lua）
const unlockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
else
  return 0
end`
```

---

## 2.7 Docker 入门 ★

`docker-compose.yml` 示意：

```yaml
services:
  app:
    build: .
    ports: ["8080:8080"]
    environment:
      MYSQL_DSN: user:pass@tcp(mysql:3306)/app?parseTime=true
      REDIS_ADDR: redis:6379
    depends_on: [mysql, redis]
  mysql:
    image: mysql:8
    environment:
      MYSQL_ROOT_PASSWORD: pass
      MYSQL_DATABASE: app
    volumes: ["mysql_data:/var/lib/mysql"]
  redis:
    image: redis:7
volumes:
  mysql_data:
```

---

### 第 2 章自问

1. 为什么深分页慢？cursor 如何缓解？
2. 缓存与 DB 更新顺序为什么常「先 DB 再删缓存」？
3. JWT 无法立刻作废时业务上怎么补？

---

# 第 3 章 工程化与可靠性

## 3.1 超时、取消、重试 ★

```go
func callDownstream(ctx context.Context, url string) error {
    var last error
    for i := 0; i < 3; i++ {
        req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        resp, err := http.DefaultClient.Do(req)
        if err == nil && resp.StatusCode < 500 {
            resp.Body.Close()
            return nil // 4xx 多数不重试
        }
        if err == nil {
            resp.Body.Close()
            last = fmt.Errorf("status %d", resp.StatusCode)
        } else {
            last = err
        }
        // 指数退避 + 抖动
        backoff := time.Duration(1<<i) * 100 * time.Millisecond
        jitter := time.Duration(rand.Intn(50)) * time.Millisecond
        select {
        case <-time.After(backoff + jitter):
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return last
}
```

**不可重试**：非幂等 POST 已可能成功、业务校验失败（400）、未授权（401/403）。

---

## 3.2 幂等、限流 ★／◆

### DB 唯一键幂等

```sql
CREATE TABLE idempotency (
  id_key VARCHAR(64) PRIMARY KEY,
  response JSON NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```go
_, err := db.ExecContext(ctx,
    `INSERT INTO idempotency(id_key, response) VALUES(?, ?)`, key, respBody)
if isDuplicate(err) {
    // 读出首次响应返回
}
```

### 令牌桶限流（进程内示意）

```go
type limiter struct {
    ch chan struct{}
}

func newLimiter(qps int) *limiter {
    l := &limiter{ch: make(chan struct{}, qps)}
    for i := 0; i < qps; i++ {
        l.ch <- struct{}{}
    }
    go func() {
        t := time.NewTicker(time.Second / time.Duration(qps))
        for range t.C {
            select {
            case l.ch <- struct{}{}:
            default:
            }
        }
    }()
    return l
}

func (l *limiter) Allow() bool {
    select {
    case <-l.ch:
        return true
    default:
        return false
    }
}
```

---

## 3.3 优雅关闭 ★

```go
srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}

go func() {
    if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Fatal(err)
    }
}()

stop := make(chan os.Signal, 1)
signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
<-stop

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
_ = srv.Shutdown(ctx) // 停接新请求，等旧请求结束
_ = db.Close()
```

---

## 3.5 pprof 怎么开 ★／◆

```go
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
}()
```

```bash
go tool pprof http://127.0.0.1:6060/debug/pprof/profile?seconds=30
go tool pprof http://127.0.0.1:6060/debug/pprof/heap
# goroutine 泄漏看：
curl http://127.0.0.1:6060/debug/pprof/goroutine?debug=1
```

---

### 第 3 章自问

1. 什么样的错误不能重试？
2. Shutdown 和直接 `os.Exit` 差别？
3. P99 延迟和平均值哪个更重要？为什么？

---

# 第 4 章 分布式与异步

## 4.1 gRPC + Protobuf ★／◆

`user.proto`：

```protobuf
syntax = "proto3";
package user.v1;
option go_package = "example.com/app/gen/user/v1;userv1";

message GetUserRequest { string id = 1; }
message User { string id = 1; string name = 2; }
service UserService {
  rpc GetUser(GetUserRequest) returns (User);
}
```

服务端片段：

```go
func (s *Server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.User, error) {
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "id required")
    }
    u, err := s.repo.Find(ctx, req.Id)
    if errors.Is(err, ErrNotFound) {
        return nil, status.Error(codes.NotFound, "user not found")
    }
    return &userv1.User{Id: u.ID, Name: u.Name}, nil
}
```

客户端务必带超时：

```go
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()
resp, err := client.GetUser(ctx, &userv1.GetUserRequest{Id: "1"})
```

---

## 4.2 消息队列：至少一次 + 幂等 ★／◆

```go
func consume(msg Message) error {
    // 1. 幂等：处理记录唯一键
    ok, err := insertProcessed(msg.ID) // INSERT 唯一索引
    if err != nil {
        return err
    }
    if !ok {
        return nil // 已处理过，直接 ack
    }
    // 2. 业务
    if err := createOrder(msg.Payload); err != nil {
        _ = deleteProcessed(msg.ID) // 或依赖事务/状态机回滚
        return err                  // 返回错误以便重投
    }
    return nil // ack
}
```

Outbox 示意：

```sql
-- 与业务同一事务
BEGIN;
INSERT INTO orders(...);
INSERT INTO outbox(id, topic, payload, created_at) VALUES(...);
COMMIT;
-- 异步进程扫 outbox 发 MQ，成功后标记已发送
```

---

## 4.3 最终一致口述要点

跨服务不追求单机事务时：本地提交 + 消息通知 + 补偿/对账。分布式锁只护临界区，**不能代替幂等**。

---

### 第 4 章自问

1. 为什么说「至少一次 + 幂等」是实务标准答案？
2. Kafka 分区数与消费并行度关系？
3. 什么时候不该用 MQ？

---

# 第 5 章 云原生与差异化

## 5.1 多阶段 Dockerfile ★

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/app

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
```

---

## 5.2 K8s 最小清单 ★／◆

```yaml
apiVersion: apps/v1
kind: Deployment
metadata: { name: app }
spec:
  replicas: 2
  selector: { matchLabels: { app: app } }
  template:
    metadata: { labels: { app: app } }
    spec:
      containers:
        - name: app
          image: myrepo/app:1.0.0
          ports: [{ containerPort: 8080 }]
          readinessProbe:
            httpGet: { path: /healthz, port: 8080 }
            initialDelaySeconds: 3
          livenessProbe:
            httpGet: { path: /livez, port: 8080 }
            initialDelaySeconds: 10
          resources:
            requests: { cpu: "100m", memory: "128Mi" }
            limits: { cpu: "500m", memory: "256Mi" }
---
apiVersion: v1
kind: Service
metadata: { name: app }
spec:
  selector: { app: app }
  ports: [{ port: 80, targetPort: 8080 }]
```

**readiness 失败、liveness 成功**：Pod 还在，但 Service 不转发新流量。

---

## 5.3 Terraform Provider 行为对照 ◆

```text
Create: API 创建 → 等 job/状态 Active → 把 ID 写入 state
Read:   用 ID 查云侧；404 → 从 state 移除（资源已没）
Update: diff → 调更新；ForceNew 字段 → 先删后建
Delete: 调删除；云侧已 404 → 仍返回成功（幂等删）
```

异步创建伪代码思路：

```go
job, err := client.CreateCluster(ctx, req)
if err != nil { return err }
deadline := time.Now().Add(30 * time.Minute)
for time.Now().Before(deadline) {
    st, err := client.GetJob(ctx, job.ID)
    if err != nil { return retryable(err) }
    if st == "SUCCESS" { break }
    if st == "FAILED" { return fmt.Errorf("create failed") }
    time.Sleep(10 * time.Second)
}
d.SetId(clusterID)
```

---

### 第 5 章自问

1. readiness 失败但 liveness 成功时流量怎么走？
2. TF Read 返回 404 通常应如何更新 state？
3. 为什么删除 API 返回「不存在」时 Delete 仍可返回成功？

---

# 第 6 章 计算机基础（面试）

## 6.4 切片传参示意 ★

```go
func mod(s []int) {
    s[0] = 100          // 外部可见
    s = append(s, 9)    // 若扩容，外部看不到新元素
}
a := []int{1, 2}
mod(a)
fmt.Println(a) // [100 2]
```

其余 OS/网络/DB 原理：能口述进程/线程/协程、三次握手、索引与 MVCC 直觉即可。

---

# 第 7 章 算法（示例：二分边界）

```go
// 找第一个 >= target 的下标
func lowerBound(a []int, target int) int {
    l, r := 0, len(a) // 半开区间 [l,r)
    for l < r {
        m := l + (r-l)/2
        if a[m] < target {
            l = m + 1
        } else {
            r = m
        }
    }
    return l
}
```

---

# 第 8 章 项目深挖

用 STAR 讲缓存一致性、用 Provider 讲「异步 job + Read 对账 + 幂等删除」。能量化就写 QPS/P99。

---

# 第 9 章 概念对照速查

| 概念       | 一句话            |
| -------- | -------------- |
| 幂等       | 做多次等于做一次       |
| 最终一致     | 允许中间态，稍后对齐     |
| 优雅关闭     | 先摘流量再停进程       |
| 至少一次     | 可能重复，靠幂等       |
| ForceNew | 该字段变了必须重建资源    |
| 覆盖索引     | 索引已含查询列，无需回表   |
| 工作窃取     | 空闲 P 偷别的 P 的 G |
| 熔断       | 下游病了先快速失败      |

---

# 第 10 章 建议阅读顺序

1. 第 1 章 + 手敲切片/context/worker 示例
2. 第 2 章 + 把 Gin/MySQL/Redis 示例拼进主项目
3. 第 3 章 + 给项目加上 Shutdown 与超时重试
4. 第 4 章 + 写一个幂等消费者
5. 第 5 章 + 多阶段镜像 / 最小 K8s YAML
6. 第 6–7 章穿插面试与算法
7. 第 8 章写口述稿并模拟面试

---

## 附录：每日学习卡片模板

```text
日期：
知识点：
我的理解（不看文档写 5–10 行）：
代码验证（做了什么）：
易错点：
面试一句话：
仍不清楚：
```

---

## 文档关系

| 文件                                 | 作用             |
| ---------------------------------- | -------------- |
| `golang-job-12week-plan.md`        | 时间表与验收         |
| `golang-knowledge-points.md`       | 知识点清单与优先级      |
| `golang-knowledge-detailed.md`（本文） | 详解 + **关键点示例** |

已在第 1–5 章及算法/切片等高频点补充必要示例。原无示例备份见：`golang-knowledge-detailed.backup.md`。
