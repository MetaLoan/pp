# 核心算法规范 (Core Algorithms Specification)

## 1. 动态赔率引擎 (Dynamic Payout Engine)

本算法旨在通过数学模型动态调整赔率，以平衡平台的多空风险敞口，确保平台在长期运行中的数学期望为正，同时保持对用户的吸引力。

### 1.1 基础公式
最终赔率 $P_{final}$ 由基础赔率 $P_{base}$ 经过风险系数修正得到。

$$ P_{final} = P_{base} \times (1 - K_{vol} \times V) \times K_{skew} $$

其中：
*   $P_{base}$: 资产的基础赔率 (如 0.85)。
*   $V$: 当前市场波动率 (Volatility)，通常取过去 1 分钟价格标准差。
*   $K_{vol}$: 波动率惩罚系数。波动越大，预测难度变大（或被操纵风险变大），赔率微调降低。
*   $K_{skew}$: 多空倾斜系数 (Long/Short Skew Adjustment)。

### 1.2 多空倾斜系数 ($K_{skew}$) 计算
旨在平衡资金池。如果买涨资金过多，降低涨的赔率，提高跌的赔率。

设 $Pool_{call}$ 为当前未结算的看涨资金池，$Pool_{put}$ 为看跌资金池。
计算倾斜比率 $R = \frac{Pool_{call}}{Pool_{call} + Pool_{put}}$

对于看涨 (Call) 订单：
$$ K_{skew\_call} = \begin{cases} 1.0 & \text{if } R \le 0.5 \\ 1.0 - \alpha \times (R - 0.5) & \text{if } R > 0.5 \end{cases} $$

对于看跌 (Put) 订单：
$$ K_{skew\_put} = \begin{cases} 1.0 & \text{if } R \ge 0.5 \\ 1.0 - \alpha \times (0.5 - R) & \text{if } R < 0.5 \end{cases} $$

*   $\alpha$: 敏感度系数 (如 0.2)。当一方资金占比达到 100% 时，赔率会显著降低，抑制继续单边下注。

---

## 2. 反延迟套利算法 (Anti-Latency Arbitrage)

### 2.1 订单时间戳校验
防止用户利用过期的价格信息下单。

1.  客户端发送订单时包含 `client_ts` (毫秒级时间戳)。
2.  服务器接收时记录 `server_ts`。
3.  计算延迟 $\Delta t = server\_ts - client\_ts$。
4.  **规则**:
    *   如果 $\Delta t > 200ms$ (可配置阈值)，拒绝订单 (Error: "Price expired")。
    *   如果 $\Delta t < 0$，标记为异常 (时钟不同步或攻击)。

### 2.2 滑点控制 (Slippage Control)
防止在价格剧烈变动瞬间成交。

1.  用户下单请求包含 `request_price` (用户看到的价格)。
2.  服务器撮合时获取 `current_price` (最新聚合价格)。
3.  计算偏差 $Diff = \frac{|current\_price - request\_price|}{request\_price}$。
4.  **规则**:
    *   如果 $Diff > 0.0005$ (0.05%)，触发 **Re-quote** (重新报价)，询问用户是否接受新价格，或直接拒绝。

---

## 3. 价格聚合算法 (Price Aggregation)

确保结算价格公允，剔除单一交易所的插针数据。

### 3.1 加权平均
假设接入 $N$ 个数据源 $S_1, S_2, ..., S_N$，权重为 $W_1, W_2, ..., W_N$。

$$ P_{index} = \frac{\sum_{i=1}^{N} P_i \times W_i}{\sum_{i=1}^{N} W_i} $$

### 3.2 异常剔除 (Sanity Check)
在计算平均值前，先进行预处理：
1.  计算所有源的中位数 $P_{median}$。
2.  对于每个源 $P_i$，如果 $|P_i - P_{median}| / P_{median} > 0.5\%$，则将其权重 $W_i$ 临时置为 0。
3.  如果有效源数量 $< 2$，暂停交易 (Circuit Breaker)，防止错误报价。

---

## 4. 游戏化升级算法 (Level Up Logic)

基于 XP (经验值) 的非线性升级曲线，增加用户粘性。

$$ XP_{required} = Base \times (Level)^{1.5} $$

*   $Base$: 基础经验值 (如 1000)。
*   用户每交易 1 USD，获得 1 XP。
*   **奖励**:
    *   Level 1 -> 2: 赠送 100 积分。
    *   Level 5: 解锁 "高级图表" 功能。
    *   Level 10: 赔率加成 +1%。
