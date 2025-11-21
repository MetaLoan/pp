# 数据库设计文档 (Database Schema Design)

## 1. 概述
本系统采用 PostgreSQL 作为核心关系型数据库，Redis 作为缓存和即时状态存储，InfluxDB 用于存储时序数据（行情、K线）。

## 2. PostgreSQL 核心表结构

### 2.1 用户与认证 (Users & Auth)
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE, -- 对外暴露的 ID
    email VARCHAR(255) UNIQUE,
    telegram_id BIGINT UNIQUE, -- Telegram 登录支持
    password_hash VARCHAR(255),
    nickname VARCHAR(50),
    avatar_url TEXT,
    role VARCHAR(20) DEFAULT 'user', -- user, admin, agent
    risk_level VARCHAR(20) DEFAULT 'normal', -- normal, high_risk, whale
    status VARCHAR(20) DEFAULT 'active', -- active, suspended, banned
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_kyc (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    country VARCHAR(50),
    id_number VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 2.2 资产与钱包 (Assets & Wallets)
```sql
CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL UNIQUE, -- e.g., EURUSD, BTCUSDT
    type VARCHAR(20) NOT NULL, -- forex, crypto, stock
    is_active BOOLEAN DEFAULT true,
    base_payout DECIMAL(5, 2) DEFAULT 0.85, -- 基础赔率，如 0.85 代表 85%
    min_trade_amount DECIMAL(18, 8) DEFAULT 1.0,
    max_trade_amount DECIMAL(18, 8) DEFAULT 1000.0
);

CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    currency VARCHAR(10) NOT NULL, -- USDT, USD, GAME_POINT
    balance DECIMAL(24, 8) DEFAULT 0,
    frozen DECIMAL(24, 8) DEFAULT 0, -- 下单冻结金额
    version BIGINT DEFAULT 0, -- 乐观锁版本号
    UNIQUE(user_id, currency)
);
```

### 2.3 交易订单 (Orders)
```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(64) UNIQUE NOT NULL,
    user_id BIGINT REFERENCES users(id),
    asset_symbol VARCHAR(20) NOT NULL,
    direction VARCHAR(10) NOT NULL, -- CALL (涨), PUT (跌)
    amount DECIMAL(18, 8) NOT NULL,
    payout_rate DECIMAL(5, 4) NOT NULL, -- 下单时的赔率，如 0.8500
    
    open_price DECIMAL(18, 8) NOT NULL, -- 开仓价格
    open_time TIMESTAMPTZ NOT NULL,
    
    close_price DECIMAL(18, 8), -- 结算价格 (空为未结算)
    close_time TIMESTAMPTZ NOT NULL, -- 预计结算时间
    
    status VARCHAR(20) DEFAULT 'pending', -- pending, active, won, lost, draw
    pnl DECIMAL(18, 8) DEFAULT 0, -- 盈亏金额 (+85 或 -100)
    
    device_info JSONB, -- 风控用：IP, UA, DeviceID
    
    -- 挂单与道具支持
    type VARCHAR(20) DEFAULT 'market', -- market, pending_price, pending_time
    trigger_price DECIMAL(18, 8), -- 挂单触发价格
    trigger_time TIMESTAMPTZ, -- 挂单触发时间
    used_item_id BIGINT, -- 关联使用的道具ID (如免死金牌)
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
CREATE INDEX idx_orders_close_time ON orders(close_time) WHERE status = 'active';
```

### 2.4 资金流水 (Transactions)
```sql
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    wallet_id BIGINT REFERENCES wallets(id),
    type VARCHAR(30) NOT NULL, -- deposit, withdraw, trade_lock, trade_settle, bonus
    amount DECIMAL(24, 8) NOT NULL, -- 变动金额 (+/-)
    balance_after DECIMAL(24, 8) NOT NULL,
    reference_id VARCHAR(64), -- 关联订单号或支付流水号
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 2.5 游戏化与任务 (Gamification)
```sql
CREATE TABLE user_levels (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    current_xp BIGINT DEFAULT 0,
    level INT DEFAULT 1,
    next_level_xp BIGINT DEFAULT 1000
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    condition_type VARCHAR(50), -- trade_count, win_streak
    condition_value JSONB, -- { "count": 5 }
    reward_type VARCHAR(20), -- point, bonus
    reward_amount DECIMAL(18, 8)
);

CREATE TABLE tournaments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100),
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    prize_pool DECIMAL(18, 8),
    currency VARCHAR(10), -- USDT, USD
    status VARCHAR(20), -- upcoming, active, completed
    entry_fee DECIMAL(18, 8) DEFAULT 0,
    min_level INT DEFAULT 1 -- 参与门槛
);

CREATE TABLE tournament_participants (
    tournament_id INT REFERENCES tournaments(id),
    user_id BIGINT REFERENCES users(id),
    initial_balance DECIMAL(18, 8), -- 比赛初始资金
    current_balance DECIMAL(18, 8), -- 当前比赛资金
    rank INT,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tournament_id, user_id)
);
```

### 2.6 道具与商城 (Market & Items)
```sql
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50), -- e.g., "Risk Free Trade $10"
    type VARCHAR(30), -- risk_free, bonus, cashback, payout_boost
    effect_config JSONB, -- { "max_amount": 10, "refund_rate": 1.0 }
    price_gems INT, -- 购买价格 (宝石)
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE user_items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    item_id INT REFERENCES items(id),
    is_used BOOLEAN DEFAULT false,
    acquired_at TIMESTAMPTZ DEFAULT NOW(),
    used_at TIMESTAMPTZ,
    order_id BIGINT -- 关联使用该道具的订单
);
```

### 2.7 交易信号 (Signals)
```sql
CREATE TABLE trading_signals (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20),
    direction VARCHAR(10), -- CALL, PUT
    timeframe VARCHAR(10), -- M1, M5
    strength DECIMAL(3, 2), -- 信号强度 0.0 - 1.0
    source VARCHAR(50), -- AI_MODEL_A, INDICATOR_RSI
    generated_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);
```

## 3. Redis 数据结构设计

### 3.1 实时行情 (Market Data)
*   **Key**: `market:ticker:{symbol}`
*   **Value**: JSON `{ "price": 1.0523, "ts": 1715000000 }`
*   **TTL**: 无 (实时更新)

### 3.2 订单结算队列 (Settlement Queue)
*   **Type**: ZSET (Sorted Set)
*   **Key**: `orders:active`
*   **Score**: `close_time` (Unix Timestamp)
*   **Member**: `order_id`
*   *用途*: 轮询 ZRANGE 检查到期订单，实现毫秒级结算触发。

### 3.3 用户风控状态 (Risk State)
*   **Key**: `risk:user:{user_id}`
*   **Value**: Hash `{ "consecutive_wins": 5, "total_profit_today": 500.0 }`
*   *用途*: 快速读取用户当前状态，决定是否触发风控限制。

## 4. InfluxDB 时序数据

### 4.1 历史价格 (Price History)
*   **Measurement**: `price_ticks`
*   **Tags**: `symbol`, `source` (e.g., binance)
*   **Fields**: `price` (float)
*   *用途*: 生成 K 线图，以及结算时的价格回溯校验。

## 5. 运营与管理系统 (Operations & Admin)

### 5.1 代理与返佣 (Affiliate System)
```sql
CREATE TABLE affiliates (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    code VARCHAR(20) UNIQUE NOT NULL, -- 推广码
    commission_model VARCHAR(20) DEFAULT 'revshare', -- revshare (盈亏), turnover (流水), cpa (人头)
    commission_rate DECIMAL(5, 4) DEFAULT 0.0, -- 返佣比例 (e.g. 0.50 for 50%)
    parent_affiliate_id BIGINT, -- 上级代理 (支持多级)
    total_commission DECIMAL(18, 8) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE referrals (
    user_id BIGINT PRIMARY KEY REFERENCES users(id), -- 被邀请人
    affiliate_id BIGINT REFERENCES affiliates(user_id), -- 邀请人
    campaign_id VARCHAR(50), -- 推广活动追踪
    joined_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE commission_logs (
    id BIGSERIAL PRIMARY KEY,
    affiliate_id BIGINT REFERENCES affiliates(user_id),
    source_user_id BIGINT REFERENCES users(id), -- 贡献佣金的用户
    order_id BIGINT, -- 关联订单 (如果是按交易返佣)
    amount DECIMAL(18, 8) NOT NULL,
    type VARCHAR(20), -- trade_rebate, loss_share
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 5.2 系统配置与管理 (System Config & Audit)
```sql
CREATE TABLE system_configs (
    key VARCHAR(100) PRIMARY KEY, -- e.g., "risk.global_payout_limit", "payment.usdt_address"
    value TEXT,
    description VARCHAR(255),
    group_name VARCHAR(50), -- risk, payment, ui
    is_public BOOLEAN DEFAULT false, -- 是否下发给前端
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    password_hash VARCHAR(255),
    role VARCHAR(20), -- super_admin, finance, support
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE admin_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    admin_id INT REFERENCES admin_users(id),
    action VARCHAR(50), -- "approve_withdraw", "ban_user", "update_config"
    target_id VARCHAR(50), -- 被操作对象的ID
    details JSONB, -- 修改前后的值
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```
