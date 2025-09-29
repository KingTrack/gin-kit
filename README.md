# Gin-Kit

> 一个基于 Gin 框架的企业级微服务开发脚手架，提供开箱即用的中间件生态和可观测性支持

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.25-blue.svg)](https://golang.org/)
[![Gin Version](https://img.shields.io/badge/Gin-v1.10.1-green.svg)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

## ✨ 核心特性

- 🚀 **快速启动**: 零配置启动，5分钟搭建生产级HTTP服务
- 🔧 **模块化设计**: 插件化架构，按需加载组件
- 📊 **可观测性**: 内置日志、指标、链路追踪三大支柱
- 🔌 **中间件生态**: 丰富的中间件支持，覆盖常见业务场景
- 🏗️ **生产就绪**: 性能优化、错误处理、优雅关闭等企业级特性
- 📖 **配置驱动**: 支持多种配置格式和热加载
- 🔗 **外部集成**: 无缝对接 MySQL、Redis、Consul、Nacos 等
- 📈 **监控集成**: 原生支持 Prometheus、Jaeger、Zipkin、夜莺等

## 📖 目录

- [核心特性](#✨-核心特性)
- [设计哲学](#🎯-设计哲学)
- [架构设计](#🏗️-架构设计)
- [核心组件](#🔧-核心组件)
- [快速开始](#🚀-快速开始)
- [配置详解](#⚙️-配置详解)
- [中间件使用](#🔌-中间件使用)
- [监控集成](#📊-监控集成)
- [最佳实践](#💡-最佳实践)
- [API文档](#📚-api文档)
- [贡献指南](#🤝-贡献指南)
- [更新日志](#📝-更新日志)

## 🎯 设计哲学

### 🧩 模块化优先
Gin-Kit 采用高度模块化的设计思想，每个功能组件都是独立的模块，可以单独配置、使用和扩展：
- **松耦合**: 模块之间依赖最小化，便于维护和测试
- **高内聚**: 每个模块专注于单一职责
- **易扩展**: 新功能可以通过插件方式无缝集成

### 🔌 插件化架构
框架基于插件化架构设计，核心功能通过统一的注册机制进行管理：
- **注册表模式**: 所有组件通过 Registry 进行统一管理
- **工厂模式**: 动态创建和配置各种组件
- **依赖注入**: 通过 Runtime Engine 提供全局访问点

### ⚙️ 配置驱动
所有组件行为都通过配置文件驱动，支持：
- **多格式配置**: TOML、JSON、YAML
- **命名空间隔离**: 不同环境、不同服务的配置完全隔离
- **热加载**: 配置变更无需重启服务（开发中）

### 🛠️ 生产就绪
框架从设计之初就考虑生产环境需求：
- **可观测性**: 内置日志、链路追踪、指标监控
- **高性能**: 对象池化、连接复用、异步处理
- **容错性**: 优雅降级、熔断保护、错误恢复

### 👨‍💻 开发者友好
提供开发者最佳体验：
- **零配置启动**: 合理的默认配置
- **类型安全**: 强类型接口设计
- **清晰的错误信息**: 详细的错误上下文

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              Gin-Kit 架构图                                    │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                                用户层 (User Layer)                              │
├─────────────────────────┬─────────────────────────┬─────────────────────────────┤
│    Business Application │      HTTP Server        │       Client APIs           │
│    (业务应用代码)        │     (Gin HTTP 服务)     │    (客户端API调用)           │
└─────────────────────────┴─────────────────────────┴─────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            中间件层 (Middleware Layer)                          │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬───────────┤
│   Recovery  │   Context   │   Logger    │   Tracer    │   Metric    │Response   │
│  异常恢复    │  上下文管理  │  访问日志    │  链路追踪    │  指标收集    │Capture   │
│  中间件      │   中间件     │   中间件     │   中间件     │   中间件     │响应捕获   │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┴───────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          运行时引擎 (Runtime Engine)                            │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                        全局单例管理器                                    │   │
│  │                    - 组件生命周期控制                                     │   │
│  │                    - 统一访问入口                                        │   │
│  │                    - 组件间协调                                          │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            注册表层 (Registry Layer)                            │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬───────────┤
│  Context    │   Logger    │   Tracer    │   Metric    │   MySQL     │  Redis    │
│  Registry   │  Registry   │  Registry   │  Registry   │  Registry   │ Registry  │
│ 上下文注册表 │  日志注册表  │ 追踪注册表   │ 指标注册表   │ 数据库注册表 │缓存注册表  │
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┴───────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          内部组件层 (Internal Layer)                            │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬───────────┤
│ Context     │    Zap      │OpenTelemetry│ go-metrics  │    GORM     │ go-redis  │
│   Pool      │   Logger    │/OpenTracing │+ Prometheus │   MySQL     │   Redis   │
│ 对象池管理   │  高性能日志  │  分布式追踪  │   指标系统   │   数据库ORM  │  缓存客户端│
└─────────────┴─────────────┴─────────────┴─────────────┴─────────────┴───────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            配置层 (Configuration Layer)                         │
├─────────────────────────┬─────────────────────────┬─────────────────────────────┤
│   Configuration Manager │    Namespace Manager    │     Plugin System           │
│      配置管理器          │     命名空间管理器       │       插件系统              │
│                         │                         │                             │
│  ┌───────────────────┐  │  ┌───────────────────┐  │  ┌───────────────────────┐  │
│  │ TOML/JSON/YAML    │  │  │  多租户配置隔离    │  │  │   Source Plugin       │  │
│  │   配置解析器       │  │  │  环境配置分离      │  │  │   Decoder Plugin      │  │
│  └───────────────────┘  │  └───────────────────┘  │  └───────────────────────┘  │
└─────────────────────────┴─────────────────────────┴─────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            外部系统 (External Systems)                          │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────────┬───────────┤
│   Jaeger    │   Zipkin    │ SkyWalking  │ Prometheus  │    夜莺      │   MySQL   │
│  链路追踪    │  链路追踪    │  链路追踪    │   指标监控   │   指标监控    │  数据库    │
├─────────────┼─────────────┼─────────────┼─────────────┼─────────────┼───────────┤
│            Redis            │         PushGateway         │      其他外部服务      │
│           缓存服务           │          推送网关            │                       │
└─────────────────────────────┴─────────────────────────────┴─────────────────────┘
```

### 📊 架构层次说明

#### 1. 用户层 (User Layer)
- **Business Application**: 业务应用代码
- **HTTP Server**: Gin HTTP 服务器
- **Client APIs**: 各种客户端API

#### 2. 中间件层 (Middleware Layer)
- **Recovery**: 异常恢复中间件
- **Context**: 请求上下文管理
- **Logger**: 访问日志记录
- **Tracer**: 分布式链路追踪
- **Metric**: 指标收集

#### 3. 运行时引擎 (Runtime Engine)
- 全局单例，管理所有组件的生命周期
- 提供统一的访问入口
- 负责组件间的协调

#### 4. 注册表层 (Registry Layer)
- 各种资源的注册表和管理器
- 实现组件的创建、配置和销毁
- 提供组件实例的获取接口

## 🔧 核心组件

### 🔧 核心模块

#### 1. 🚀 Runtime Engine (运行时引擎)
**路径**: [`kit/runtime/`](kit/runtime/) & [`kit/engine/`](kit/engine/)

运行时引擎是整个框架的核心，负责：
- 全局组件管理和生命周期控制
- 配置加载和命名空间管理
- 各种 Registry 的初始化和协调

```go
// 获取全局引擎实例
engine := runtime.Get()

// 访问各种注册表
logger := engine.LoggerRegistry().AppLogger()
db := engine.MySQLRegistry().GetDB(ctx, "main")
redis := engine.RedisRegistry().GetClient(ctx, "cache")
```

#### 2. 🔄 Context Management (上下文管理)
**路径**: [`kit/internal/context/`](kit/internal/context/)

高性能的请求上下文管理：
- 对象池化减少GC压力
- 自动的生命周期管理
- 线程安全的访问控制

#### 3. 📝 Logger System (日志系统)
**路径**: [`kit/internal/logger/`](kit/internal/logger/) & [`kit/client/logger/`](kit/client/logger/)

基于 Zap 的高性能日志系统：
- 多种日志类型：访问日志、应用日志、错误日志等
- 结构化日志输出
- 日志轮转和归档

#### 4. 📊 Metrics Collection (指标收集)
**路径**: [`kit/internal/metric/`](kit/internal/metric/) & [`kit/client/metric/`](kit/client/metric/)

基于 go-metrics 的指标系统：
- 多种指标类型：Counter、Gauge、Timer、Histogram、Meter
- Prometheus 格式输出
- 支持夜莺（n9e）监控平台
- P50/P95/P99 分位数统计

#### 5. 🔗 HTTP Middleware Chain (HTTP中间件链)
**路径**: [`kit/httpserver/internal/middleware/`](kit/httpserver/internal/middleware/)

丰富的中间件支持：
- **Recovery**: 异常恢复和错误处理
- **Context**: 请求上下文管理
- **Logger**: 结构化访问日志
- **Tracer**: 分布式追踪集成
- **Metric**: 请求指标收集
- **ResponseCapture**: 响应内容捕获

## 🚀 快速开始

### 1. 环境要求

- **Go**: >= 1.25.0
- **操作系统**: Linux、macOS、Windows
- **内存**: 建议 >= 512MB
- **外部依赖**: 可选，支持 MySQL、Redis、Consul、Nacos、ETCD 等
- **监控系统**: 可选，支持 Prometheus、Jaeger、Zipkin、SkyWalking、夜莺等

### 2. 安装

```bash
# 克隆项目
git clone git.inke.cn/nvwa/httpserver/gin-kit.git
cd gin-kit

# 安装依赖
go mod tidy
```

### 3. 创建最小示例

#### 3.1 创建配置文件 `config.toml`

```toml
[httpserver]
service_name = "gin-kit-demo"
port = 8080
read_timeout_sec = 30
write_timeout_sec = 30
idle_timeout_sec = 60

[logger]
level = "info"
log_dir = "./logs"
max_size = 100
max_backups = 10
max_age = 30

[metric]
enabled = true
service_name = "gin-kit-demo"
backend_name = "prometheus"

[metric.prometheus]
path = "/metrics"

[tracer]
service_name = "gin-kit-demo"
enabled = false
proto = "OpenTelemetry"
backend_name = "jaeger"
report_url = "http://localhost:14268/api/traces"
```

#### 3.2 创建主程序 `main.go`

```go
package main

import (
    "log"

    "github.com/KingTrack/gin-kit/kit/conf"
    "github.com/KingTrack/gin-kit/kit/engine"
    "github.com/KingTrack/gin-kit/kit/httpserver"
    contextmiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/context"
    loggermiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/logger"
    metricmiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/metric"
    recovermiddleware "github.com/KingTrack/gin-kit/kit/httpserver/internal/middleware/recover"
    "github.com/KingTrack/gin-kit/kit/plugin/decoder"
    "github.com/KingTrack/gin-kit/kit/plugin/source"
    "github.com/KingTrack/gin-kit/kit/runtime"
    "github.com/gin-gonic/gin"
)

func main() {
    // 1. 创建引擎
    e := engine.New("./")
    
    // 2. 初始化配置
    namespace := &conf.Namespace{
        RootPath: "./config.toml",
        Source:   &source.File{},
        Decoder:  &decoder.Toml{},
    }
    
    if err := e.Init(engine.WithNamespace(namespace)); err != nil {
        log.Fatal("Engine initialization failed:", err)
    }
    
    // 3. 设置全局运行时
    runtime.Set(e)
    
    // 4. 创建 HTTP 服务器
    server := httpserver.New(
        httpserver.WithRecovery(&recovermiddleware.Middleware{}),
        httpserver.WithContext(&contextmiddleware.Middleware{}),
        httpserver.WithLogger(&loggermiddleware.Middleware{}),
        httpserver.WithMetric(&metricmiddleware.Middleware{}),
    )
    
    // 5. 添加路由
    setupRoutes(server)
    
    // 6. 启动服务器
    log.Println("Starting httpserver...")
    if err := server.Run(); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}

func setupRoutes(server *httpserver.Server) {
    // 健康检查
    server.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
    })
    
    // API 路由组
    api := server.Group("/api/v1")
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
        api.GET("/users/:id", getUserByID)
    }
}

func getUsers(c *gin.Context) {
    users := []map[string]interface{}{
        {"id": 1, "name": "张三", "email": "zhangsan@example.com"},
        {"id": 2, "name": "李四", "email": "lisi@example.com"},
    }
    c.JSON(200, gin.H{"data": users, "count": len(users)})
}

func createUser(c *gin.Context) {
    c.JSON(201, gin.H{"message": "用户创建成功"})
}
```

#### 3.3 运行示例

```bash
# 运行程序
go run main.go

# 输出信息
# Server is running on :8080

# 测试接口
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"王五","email":"wangwu@example.com"}'

# 查看监控指标（如果启用了 Prometheus）
curl http://localhost:8080/metrics
```

## ⚙️ 配置详解

### 📝 完整配置示例

```toml
# 全局命名空间配置
namespace = "production"
hostname = "app-server-01"

# HTTP 服务器配置
[httpserver]
service_name = "gin-kit-demo"          # 服务名称
port = 8080                            # 监听端口
read_timeout_sec = 30                  # 读取超时（秒）
write_timeout_sec = 30                 # 写入超时（秒）
idle_timeout_sec = 60                  # 空闲超时（秒）
print_request_body_size_kb = 4         # 打印请求体大小限制（KB）
print_response_body_size_kb = 4        # 打印响应体大小限制（KB）
close_request_body = false             # 是否关闭请求体记录
close_response_body = false            # 是否关闭响应体记录

# 日志配置
[logger]
level = "info"                         # 日志级别：debug, info, warn, error
log_dir = "./logs"                     # 日志目录
max_size = 100                         # 单个日志文件最大大小（MB）
max_backups = 10                       # 保留的历史日志文件数量
max_age = 30                           # 日志文件保留天数
compress = true                        # 是否压缩历史日志

# 指标配置
[metric]
service_name = "gin-kit-demo"          # 服务名称
backend_name = "prometheus"            # 后端类型：prometheus, n9e
endpoint = "localhost"                 # 端点地址

# Prometheus 配置
[metric.prometheus]
path = "/metrics"                      # 指标暴露路径

# 夜莺（n9e）配置
[metric.n9e]
url = "http://n9e.example.com"         # 夜莺服务地址
token = "your_token_here"              # 认证令牌
interval_sec = 60                      # 上报间隔（秒）
step_sec = 15                          # 采集步长（秒）

# 链路追踪配置
[tracer]
service_name = "gin-kit-demo"          # 服务名称
enabled = true                         # 是否启用追踪
proto = "OpenTelemetry"                # 协议：OpenTelemetry, OpenTracing
backend_name = "jaeger"                # 后端：jaeger, zipkin, skywalking
report_url = "http://localhost:14268/api/traces"  # 上报地址

# 数据中心配置（服务发现）
[datacenter]
registry_type = "consul"               # 注册中心类型：consul, nacos, etcd

[datacenter.consul]
address = "localhost:8500"             # Consul 地址

[datacenter.nacos]
addresses = ["localhost:8848"]         # Nacos 地址列表
namespace_id = "public"               # 命名空间 ID
group = "DEFAULT_GROUP"               # 分组

[datacenter.etcd]
endpoints = ["localhost:2379"]        # ETCD 端点列表

# MySQL 数据库配置（支持多个数据库）
[[mysql]]
name = "main"                          # 数据库名称
dsn = "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
max_open_conns = 100                   # 最大打开连接数
max_idle_conns = 10                    # 最大空闲连接数
conn_max_lifetime_minutes = 60         # 连接最大生存时间（分钟）

[[mysql]]
name = "analytics"                     # 分析数据库
dsn = "user:password@tcp(analytics-db:3306)/analytics?charset=utf8mb4&parseTime=True&loc=Local"
max_open_conns = 50
max_idle_conns = 5
conn_max_lifetime_minutes = 30

# Redis 缓存配置（支持多个 Redis 实例）
[[redis]]
name = "cache"                         # Redis 实例名称
addr = "localhost:6379"               # Redis 地址
password = ""                         # 密码
db = 0                                # 数据库编号
pool_size = 10                        # 连接池大小
min_idle_conns = 5                    # 最小空闲连接数
max_conn_age_minutes = 30             # 连接最大存活时间（分钟）
pool_timeout_seconds = 5              # 池超时时间（秒）
idle_timeout_minutes = 5              # 空闲超时时间（分钟）
idle_check_frequency_minutes = 1      # 空闲检查频率（分钟）

[[redis]]
name = "session"                      # 会话 Redis
addr = "localhost:6380"
password = "session_password"
db = 1
pool_size = 20
```

### 🔧 按需配置

您可以根据实际需求选择性配置各个模块：

```toml
# 最小化配置 - 只启用基础功能
[httpserver]
service_name = "my-service"
port = 8080

[logger]
level = "info"
```

```toml
# 生产环境配置 - 启用完整监控
[httpserver]
service_name = "my-service"
port = 8080

[logger]
level = "warn"
log_dir = "/var/log/my-service"

[metric]
service_name = "my-service"
backend_name = "prometheus"

[tracer]
service_name = "my-service"
enabled = true
proto = "OpenTelemetry"
backend_name = "jaeger"
report_url = "http://jaeger:14268/api/traces"
```

## 🔌 中间件使用

### 📊 指标收集中间件

```go
// 创建自定义指标
metricclient := client.New("business")

// 计数器 - 统计事件发生次数
metricclient.IncCounter("user_login_total", map[string]string{
    "method": "password",
    "status": "success",
})

// 计时器 - 统计操作耗时和 QPS
stopTimer := metricclient.StartTimer("database_query_duration", map[string]string{
    "table": "users",
    "operation": "select",
})
defer stopTimer()

// 仪表盘 - 记录实时数值
metricclient.UpdateGauge("active_connections", 150, map[string]string{
    "server": "web-01",
})

// 直方图 - 统计数值分布
metricclient.UpdateHistogram("request_size_bytes", float64(requestSize), map[string]string{
    "endpoint": "/api/users",
})
```

### 📝 日志中间件

```go
loggerclient := client.New("business")
logger := loggerclient.Logger()

// 结构化日志记录
logger.Info("用户登录",
    zap.String("user_id", "12345"),
    zap.String("ip", "192.168.1.100"),
    zap.Duration("duration", time.Since(start)),
)

logger.Error("数据库连接失败",
    zap.Error(err),
    zap.String("database", "users"),
    zap.Int("retry_count", 3),
)

// 访问日志自动记录，包含：
// - 请求方法、路径、状态码
// - 请求耗时、客户端 IP
// - 链路追踪 ID
// - 请求和响应体（可配置大小限制）
```

### 🔍 链路追踪中间件

```go
// OpenTelemetry 示例
func businessHandler(c *gin.Context) {
    ctx := c.Request.Context()
    tracer := otel.Tracer("business-service")
    
    // 创建子 Span
    ctx, span := tracer.Start(ctx, "process-user-request")
    defer span.End()
    
    // 添加标签
    span.SetAttributes(
        attribute.String("user.id", "12345"),
        attribute.String("operation", "get_profile"),
    )
    
    // 业务逻辑
    result, err := processUser(ctx, "12345")
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        c.JSON(500, gin.H{"error": "处理失败"})
        return
    }
    
    c.JSON(200, result)
}
```

### 🗃️ 数据库中间件

```go
// MySQL 客户端使用
mysqlClient := client.New("main") // 对应配置中的 name

// 主库操作
db := mysqlClient.Master(ctx)
var user User
if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
    return err
}

// 从库操作（如果配置了读写分离）
db = mysqlClient.Slave(ctx)
var users []User
if err := db.Where("status = ?", "active").Find(&users).Error; err != nil {
    return err
}

// Redis 客户端使用
redisClient := client.New("cache") // 对应配置中的 name
cacheClient := redisClient.GetClient(ctx)

// 设置缓存
if err := cacheClient.Set(ctx, "user:12345", userData, time.Hour).Err(); err != nil {
    return err
}

// 获取缓存
val, err := cacheClient.Get(ctx, "user:12345").Result()
if err == redis.Nil {
    // 缓存不存在
} else if err != nil {
    return err
}
```

## 🤝 如何为本项目提交代码

### 1. 开发环境设置

```bash
# Fork 项目到你的账号，然后克隆
git clone git.inke.cn/your-username/gin-kit.git
cd gin-kit

# 添加上游仓库
git remote add upstream git.inke.cn/nvwa/httpserver/gin-kit.git

# 安装开发工具
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. 代码规范

#### 2.1 Go 代码规范

```bash
# 格式化代码
goimports -w .

# 代码检查
golangci-lint run

# 运行测试
go test ./...
```

#### 2.2 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]
```

**Type 类型：**
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建或工具相关

**示例：**
```bash
git commit -m "feat(metric): add go-metrics integration with prometheus export"
git commit -m "fix(logger): resolve memory leak in log rotation"
git commit -m "docs(readme): update quick start guide"
```

### 3. 分支管理规范

#### 🌳 分支类型

我们使用 **Git Flow** 分支模型，主要包含以下分支类型：

| 分支类型 | 命名规范 | 用途 | 生命周期 |
|---------|----------|------|----------|
| `main` | `main` | 主分支，保持稳定可发布状态 | 永久 |
| `develop` | `develop` | 开发分支，集成最新功能 | 永久 |
| `feature` | `feature/{功能描述}` | 功能开发分支 | 临时 |
| `hotfix` | `hotfix/{版本号}-{问题描述}` | 紧急修复分支 | 临时 |
| `release` | `release/{版本号}` | 发布准备分支 | 临时 |
| `bugfix` | `bugfix/{问题描述}` | Bug 修复分支 | 临时 |

#### 🎯 功能分支规范（Feature Branch）

**命名格式：**
```
feature/{类型}-{简短描述}
```

**分支类型前缀：**
- `feature/add-` : 新增功能
- `feature/update-` : 功能更新
- `feature/remove-` : 功能移除
- `feature/refactor-` : 代码重构
- `feature/optimize-` : 性能优化

**示例：**
```bash
# 新增功能分支
feature/add-rate-limiting-middleware
feature/add-redis-cluster-support
feature/add-grpc-server

# 功能更新分支
feature/update-logger-format
feature/update-metric-labels

# 重构分支
feature/refactor-engine-initialization
feature/refactor-registry-pattern

# 优化分支
feature/optimize-connection-pool
feature/optimize-memory-usage
```

#### 🔧 创建和管理功能分支

**1. 创建新功能分支**
```bash
# 从 develop 分支创建功能分支
git checkout develop
git pull upstream develop
git checkout -b feature/add-rate-limiting-middleware
```

**2. 功能开发过程**
```bash
# 定期同步 develop 分支的更新
git checkout develop
git pull upstream develop
git checkout feature/add-rate-limiting-middleware
git merge develop  # 或使用 rebase: git rebase develop

# 提交代码
git add .
git commit -m "feat(middleware): add rate limiting middleware"

# 推送到远程仓库
git push origin feature/add-rate-limiting-middleware
```

**3. 完成功能开发**
```bash
# 确保分支是最新的
git checkout develop
git pull upstream develop
git checkout feature/add-rate-limiting-middleware
git rebase develop  # 保持提交历史整洁

# 推送最终版本
git push origin feature/add-rate-limiting-middleware --force-with-lease
```

#### 🐛 Bug 修复分支规范

**命名格式：**
```
bugfix/{问题类型}-{简短描述}
```

**问题类型：**
- `memory-leak` : 内存泄漏
- `deadlock` : 死锁问题
- `crash` : 程序崩溃
- `data-loss` : 数据丢失
- `security` : 安全问题
- `performance` : 性能问题

**示例：**
```bash
bugfix/memory-leak-in-context-pool
bugfix/deadlock-in-mysql-registry
bugfix/crash-on-invalid-config
```

#### 🚨 热修复分支规范

**命名格式：**
```
hotfix/{版本号}-{问题描述}
```

**示例：**
```bash
hotfix/v1.2.1-critical-memory-leak
hotfix/v1.2.1-security-vulnerability
```

**热修复流程：**
```bash
# 从 main 分支创建热修复分支
git checkout main
git pull upstream main
git checkout -b hotfix/v1.2.1-critical-memory-leak

# 进行修复
# ...

# 提交修复
git commit -m "fix(engine): resolve critical memory leak in context pool"

# 合并到 main 和 develop
git checkout main
git merge hotfix/v1.2.1-critical-memory-leak
git tag v1.2.1

git checkout develop
git merge hotfix/v1.2.1-critical-memory-leak

# 删除热修复分支
git branch -d hotfix/v1.2.1-critical-memory-leak
```

#### 📦 发布分支规范

**命名格式：**
```
release/{版本号}
```

**示例：**
```bash
release/v1.3.0
release/v2.0.0-beta.1
```

**发布流程：**
```bash
# 从 develop 创建发布分支
git checkout develop
git pull upstream develop
git checkout -b release/v1.3.0

# 更新版本号、文档等
# 只允许 bug 修复，不允许新功能

# 完成发布准备后合并到 main
git checkout main
git merge release/v1.3.0
git tag v1.3.0

# 合并回 develop
git checkout develop
git merge release/v1.3.0

# 删除发布分支
git branch -d release/v1.3.0
```

#### ⚡ 分支操作最佳实践

**1. 分支命名规则**
- 使用小写字母和连字符
- 描述要简短但有意义
- 避免使用特殊字符
- 包含工作类型和简短描述

**2. 提交频率**
```bash
# 👍 推荐：小而频繁的提交
git commit -m "feat(middleware): add rate limiter interface"
git commit -m "feat(middleware): implement token bucket algorithm"
git commit -m "feat(middleware): add rate limiter tests"
git commit -m "docs(middleware): add rate limiter documentation"

# 👎 不推荐：大而少的提交
git commit -m "feat(middleware): add complete rate limiting functionality"
```

**3. 分支同步**
```bash
# 定期同步上游更新（建议每天至少一次）
git checkout develop
git pull upstream develop
git checkout feature/your-branch
git rebase develop  # 保持提交历史整洁
```

**4. 分支清理**
```bash
# 功能合并后删除本地分支
git branch -d feature/add-rate-limiting-middleware

# 删除远程分支
git push origin --delete feature/add-rate-limiting-middleware

# 清理已合并的分支
git branch --merged | grep -v "\*\|main\|develop" | xargs -n 1 git branch -d
```

### 4. Pull Request 流程

**1. 准备 PR**
```bash
# 确保分支是最新的
git checkout develop
git pull upstream develop
git checkout feature/add-rate-limiting-middleware
git rebase develop

# 运行完整测试
go test ./...
go vet ./...
golangci-lint run

# 推送到远程
git push origin feature/add-rate-limiting-middleware --force-with-lease
```

**2. 创建 PR**
- 填写详细的 PR 模板
- 关联相关的 Issue
- 添加适当的标签
- 请求代码审查

**3. PR 要求**
- 标题遵循 Conventional Commits 规范
- 包含功能说明和测试说明
- 所有 CI 检查通过
- 至少一个维护者审批

**PR 模板示例：**
```markdown
## 📝 变更说明
简要描述本次变更的内容

## 🎯 变更类型
- [ ] Bug 修复
- [x] 新功能
- [ ] 重大变更
- [ ] 文档更新

## 🧪 测试
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 手动测试完成

## 📋 检查清单
- [x] 代码遵循项目规范
- [x] 添加了必要的测试
- [x] 更新了相关文档
- [x] 运行了 linter 检查

## 📸 截图（如适用）

## 🔗 相关链接
Closes #123
```

### 5. 代码审查要求

- 至少需要 1 个 maintainer 审批
- 所有 CI 检查必须通过
- 测试覆盖率不低于 80%
- 新功能必须包含文档

---

## 📄 许可证

本项目采用 MIT 许可证。详情请查看 [LICENSE](LICENSE) 文件。

## 🙏 贡献者

感谢所有为本项目做出贡献的开发者！

## 📞 联系我们

- 项目地址：[github.com/KingTrack/gin-kit](github.com/KingTrack/gin-kit)
- 问题反馈：[Issues](github.com/KingTrack/gin-kit/issues)
- 功能建议：[Discussions](github.com/KingTrack/gin-kit/discussions)
```