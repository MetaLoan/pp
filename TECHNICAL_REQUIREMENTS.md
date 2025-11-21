# 技术需求文档 (Technical Requirements Document)

## 1. 项目概述
本项目旨在构建一个高性能、高可用的在线交易平台（参考 Pocket Option 功能），支持多种资产（外汇、加密货币、股票等）的交易。系统分为前端用户界面和后端服务，后端采用 Go 语言技术栈以确保高并发处理能力和低延迟。

## 2. 系统架构
采用前后端分离架构，后端提供 RESTful API 和 WebSocket 服务，前端通过 HTTP 和 WebSocket 与后端交互。

### 2.1 架构图示 (概念)
```mermaid
graph TD
    Client[Web/Mobile Client] -->|HTTP/REST| Gateway[API Gateway / Nginx]
    Client -->|WebSocket| Gateway
    Gateway -->|HTTP| AuthSvc[认证服务 (Go)]
    Gateway -->|HTTP| TradeSvc[交易核心服务 (Go)]
    Gateway -->|HTTP| WalletSvc[钱包服务 (Go)]
    Gateway -->|WS| MarketSvc[行情推送服务 (Go)]
    
    TradeSvc --> DB[(PostgreSQL)]
    WalletSvc --> DB
    MarketSvc --> Redis[(Redis Pub/Sub)]
    MarketSvc --> TSDB[(InfluxDB/TimescaleDB)]
```

## 3. 后端技术栈 (Go)

### 3.1 核心语言与框架
*   **编程语言**: Go (Golang) 1.21+
*   **Web 框架**: Gin 或 Echo (高性能 HTTP 框架)
*   **微服务框架 (可选)**: go-zero 或 Kratos (如果项目规模较大)
*   **RPC 框架**: gRPC (用于内部服务间通讯)

### 3.2 数据存储
*   **关系型数据库**: PostgreSQL 15+ (用于存储用户、订单、资金流水等核心事务数据)
*   **缓存数据库**: Redis 7+ (用于 Session 管理、热点数据缓存、行情分发、排行榜)
*   **时序数据库**: InfluxDB 或 TimescaleDB (用于存储历史行情数据、K线数据)

### 3.3 实时通讯
*   **WebSocket**: `gorilla/websocket` 或 `melody` (用于实时推送行情、订单状态变更)
*   **消息队列**: Kafka 或 RabbitMQ (用于削峰填谷，异步处理撮合、结算、邮件发送等)

### 3.4 工具与库
*   **ORM**: GORM 或 sqlx
*   **配置管理**: Viper
*   **日志**: Zap 或 Logrus
*   **API 文档**: Swagger (Gin-Swagger)

## 4. 前端技术栈 (建议)

### 4.1 Web 端
*   **框架**: React 18+ 或 Vue 3
*   **构建工具**: Vite
*   **状态管理**: Zustand (React) 或 Pinia (Vue)
*   **UI 组件库**: Ant Design 或 Tailwind CSS
*   **图表库**: **TradingView Lightweight Charts** (专为金融交易设计，性能好) 或 Highcharts
*   **网络请求**: Axios (HTTP), Native WebSocket API

### 4.2 移动端 (可选)
*   **方案**: Flutter (跨平台) 或 React Native

## 5. 功能模块需求

### 5.1 用户系统 (User System)
*   **注册/登录**: 邮箱/手机号注册，支持 OAuth (Google, Facebook)。
*   **安全认证**: JWT (JSON Web Token) 双 Token 机制 (Access/Refresh)。
*   **KYC 认证**: 实名认证流程，证件上传与审核。
*   **账户等级**: 根据入金量或交易量划分等级 (Stranger -> Guru)。

### 5.2 行情系统 (Market Data System)
*   **数据源接入**: 对接外部行情 API (如 Binance, Alpha Vantage 等)。
*   **数据处理**: 清洗、聚合、生成 K 线 (1m, 5m, 1h, 1d)。
*   **实时推送**: 通过 WebSocket 向前端广播最新价格 (Tick 数据)。

### 5.3 交易系统 (Trading System)
*   **资产管理**: 支持外汇、加密货币、股票等多种资产配置。
*   **下单接口**: 支持市价单、限价单。
*   **二元期权逻辑**: 
    *   看涨/看跌 (Call/Put)。
    *   到期时间设定 (30s, 60s, 5m 等)。
    *   收益率计算与结算。
*   **撮合/结算引擎**: 高性能内存撮合，到期自动结算盈亏。

### 5.4 钱包与支付 (Wallet & Payment)
*   **多币种钱包**: 支持 USD, BTC, ETH 等余额管理。
*   **充值**: 集成第三方支付网关 (Stripe, PayPal) 及 区块链支付节点。
*   **提现**: 提现申请、审核、自动/人工打款。
*   **资金流水**: 记录每一笔资金变动，保证账目平配。

### 5.5 社交与活动 (Social & Gamification)
*   **跟单系统**: 允许用户订阅“大神”交易员，自动复制订单。
*   **排行榜**: 基于收益率、胜率的实时排行榜。
*   **锦标赛**: 定期交易比赛逻辑。

## 6. 非功能性需求

### 6.1 性能
*   **低延迟**: 行情推送延迟 < 100ms。
*   **高并发**: 支持万级以上用户同时在线，千级 TPS (每秒交易数)。

### 6.2 安全
*   **数据加密**: 敏感数据 (密码、密钥) 加密存储。
*   **传输安全**: 全站 HTTPS / WSS。
*   **防护**: 防 SQL 注入，XSS，CSRF，DDoS 防护 (Cloudflare)。

### 6.3 可靠性
*   **数据备份**: 数据库定时冷备 + 热备。
*   **容灾**: 服务多节点部署，无单点故障。

## 7. 开发与部署流程
*   **版本控制**: Git
*   **容器化**: Docker
*   **编排**: Kubernetes (K8s) 或 Docker Compose (小型部署)
*   **CI/CD**: GitHub Actions 或 GitLab CI
