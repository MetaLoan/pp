# 📊 趋势实时刷新机制 - 功能说明

## 功能概述

实现了一个动态趋势刷新系统，根据信号的复制人数和有效期自动更新[标的+时间框架]组合的趋势强度和方向。

---

## 🎯 核心逻辑

### 1. 趋势判断规则

**当前选中时间框架下，每个交易对的趋势由以下信号决定：**

```
对于 [标的 + 时间框架] 的组合：
  ├─ 筛选所有有效信号（expiryTime > now）
  ├─ 找出复制人数最多的有效信号
  ├─ 该信号的方向 → 趋势方向 (CALL↑ 或 PUT↓)
  └─ 该信号的强度 → 趋势强度 (1条或2条箭头)
```

### 2. 自动刷新触发器

#### ✅ 触发器A：信号过期
**场景：最多人复制的信号到期时**
- 定时器每1秒检查一次
- 自动移除过期信号 (expiryTime <= now)
- 自动选择新的最多人复制的信号
- 趋势立即更新

#### ✅ 触发器B：信号被超越
**场景：新信号获得更多复制人数时**
- 用户跟随某个信号时，其 `copied` 计数 +1
- 若该信号复制人数超越原最多人复制的信号
- 趋势立即切换到新的主导信号

#### ✅ 触发器C：响应式更新
- 所有趋势计算都基于实时的 `signalFeed.value`
- Vue 的响应式系统自动监控变化
- `getTrendForPair()` 实时返回最新趋势

---

## 📝 实现细节

### 文件修改位置
- `web/src/views/TradeView.vue`

### 关键代码变更

#### 1. getTrendForPair 函数 (第975行)

**改动：从基于followers改为基于copied**

```javascript
const getTrendForPair = (symbol) => {
  // 过滤该标的+时间框架的所有信号
  const matchingSignals = signalFeed.value.filter(s => 
    s.symbol === symbol && s.timeframe === selectedTimeframe.value
  );
  
  // 过滤有效信号（未过期）
  const now = Date.now();
  const validSignals = matchingSignals.filter(s => s.expiryTime > now);
  
  // 找复制人数最多的有效信号
  const maxCopied = Math.max(...validSignals.map(s => s.copied));
  const topSignals = validSignals.filter(s => s.copied === maxCopied);
  
  // 返回最多人复制的信号的方向和强度
  return {
    direction: selectedSignal.action,
    strength: selectedSignal.strength
  };
};
```

#### 2. 趋势刷新定时器 (onMounted, 第1898行)

```javascript
// 每秒检查一次信号过期情况
trendRefreshInterval = setInterval(() => {
  const now = Date.now();
  
  // 检查是否有信号过期
  let hasExpiredSignals = signalFeed.value.some(s => s.expiryTime <= now);
  
  if (hasExpiredSignals) {
    // 自动移除过期信号
    signalFeed.value = signalFeed.value.filter(s => s.expiryTime > now);
  }
}, 1000);
```

#### 3. 信号跟随时的检查 (handleSignalTrade, 第2063行)

```javascript
// 跟随信号时，复制人数+1
signal.copied += 1;

// 检查该信号是否成为此组合的最多人复制信号
const sameGroupSignals = signalFeed.value.filter(s => 
  s.symbol === signal.symbol && 
  s.timeframe === signal.timeframe && 
  s.expiryTime > Date.now()
);

if (signal.copied === Math.max(...sameGroupSignals.map(s => s.copied))) {
  // 该信号成为主导信号，趋势自动更新
  console.log(`✓ 信号"${signal.title}"成为主导信号`);
}
```

---

## 🔄 实时刷新流程

### 流程1：信号过期导致的趋势刷新

```
T0: 信号A有5个人复制 → 显示趋势A
T10: 定时器检测 → 发现信号A过期 → 移除信号A
T10: getTrendForPair自动重算 → 找信号B(3个人复制) → 显示趋势B
```

### 流程2：新信号超越导致的趋势刷新

```
T0: 信号A有5个人复制 → 显示趋势A
T5: 用户跟随信号B(原来2个人) → B变为3个人
T5: 检查后发现A仍是最多 → 趋势保持为A
T7: 用户又跟随信号B → B变为4个人
T7: 检查后发现A仍是最多 → 趋势保持为A
T9: 用户再跟随信号B → B变为5个人
T9: 检查后发现A和B并列 → 趋势切换为B（先到先得）
T11: 用户再跟随信号B → B变为6个人
T11: getTrendForPair自动更新 → 显示趋势B（6 > 5）
```

---

## 📊 显示效果

### Markets TAB 中的趋势显示

```
BTCUSDT        ↑↑ (双箭头 = strength:2)  ← 基于M30时间框架最多人复制的信号
EURUSD         ↓  (单箭头 = strength:1)  ← 基于M30时间框架最多人复制的信号
XAUUSD         -(无趋势)                  ← 没有有效信号或已过期
```

### 实时更新场景

```
【初始状态】
信号1: BTCUSDT/M30, CALL, copied=5 ─→ 显示 ↑↑
信号2: BTCUSDT/M30, PUT,  copied=2

【5秒后 - 信号2被跟随3次】
信号1: BTCUSDT/M30, CALL, copied=5
信号2: BTCUSDT/M30, PUT,  copied=5 (并列)
─→ 显示切换为信号2的方向 (PUT/↓)

【10秒后 - 信号1过期】
信号1: 已过期 → 自动移除
信号2: BTCUSDT/M30, PUT,  copied=5
─→ 显示维持为信号2 (PUT/↓)
```

---

## ⚙️ 技术细节

### 变量

| 变量 | 类型 | 说明 |
|------|------|------|
| `signalFeed` | ref(Array) | 信号列表，包含所有有效和已过期的信号 |
| `trendRefreshInterval` | number | 趋势刷新定时器ID |
| `selectedTimeframe` | ref(String) | 当前选中的时间框架 |

### 关键字段

**信号对象中相关字段：**

| 字段 | 说明 | 变化时机 |
|------|------|---------|
| `copied` | 已被复制人数 | 用户跟随信号时+1 |
| `expiryTime` | 过期时间戳 | 固定值（创建时确定） |
| `strength` | 信号强度(1\|2) | 固定值（创建时确定） |
| `action` | 方向(CALL\|PUT) | 固定值（创建时确定） |

### 响应式链路

```
signalFeed.value 变化
  ↓
getTrendForPair() 重新计算
  ↓
Markets TAB中的趋势显示更新
  ↓
用户看到实时刷新的趋势
```

---

## 🧪 测试场景

### 测试1：信号过期后自动刷新

1. ✅ 生成20条M30信号数据
2. ✅ 观察某个交易对的趋势显示
3. ⏳ 等待最多人复制的信号过期（30-60秒）
4. ✅ 验证趋势自动刷新为次多人复制的信号

### 测试2：跟随信号导致趋势切换

1. ✅ 生成20条M30信号（signal1:CALL/5人, signal2:PUT/3人）
2. ✅ 确认当前显示趋势为signal1(↑)
3. ✅ 跟随signal2两次（3→5人）
4. ✅ 验证趋势是否自动切换为signal2(↓)

### 测试3：多个标的独立刷新

1. ✅ 选择M30时间框架
2. ✅ 查看BTCUSDT和EURUSD的趋势
3. ✅ 跟随BTCUSDT的信号（导致趋势切换）
4. ✅ 验证EURUSD的趋势保持不变（独立更新）

### 测试4：跨时间框架隔离

1. ✅ 查看M30时间框架下BTCUSDT的趋势
2. ✅ 切换到M1时间框架
3. ✅ 验证BTCUSDT的趋势因为M1没有信号而消失
4. ✅ 切换回M30，趋势恢复

---

## 🔍 监控和调试

### 浏览器控制台日志

```javascript
// 当信号成为主导信号时
✓ 信号"BTCUSDT 突破"成为BTCUSDT/M30的主导信号

// 当信号过期时（自动移除，控制台无特殊输出）
```

### 查看当前趋势

在浏览器开发者工具中运行：

```javascript
// 查看所有有效信号
signalFeed.value.filter(s => s.expiryTime > Date.now())

// 查看BTCUSDT/M30的最多人复制信号
const signals = signalFeed.value.filter(s => 
  s.symbol === 'BTCUSDT' && 
  s.timeframe === 'M30' && 
  s.expiryTime > Date.now()
);
const maxCopied = Math.max(...signals.map(s => s.copied));
signals.filter(s => s.copied === maxCopied);
```

---

## 📌 重要注意事项

1. **时间框架隔离**：不同时间框架的趋势是完全独立的
   - M30的BTCUSDT趋势不受M1的BTCUSDT影响
   - 切换时间框架会导致显示的趋势改变

2. **有效期检查**：只有 `expiryTime > now` 的信号才算有效
   - 过期信号自动被移除
   - 无有效信号时不显示趋势

3. **复制人数优先级**：基于 `copied` 不是 `followers`
   - `followers`: 信号创建时的初始值（不变）
   - `copied`: 用户跟随时实时增加（影响趋势）

4. **实时性**：最多1秒的延迟
   - 信号过期检查：每1秒一次
   - 复制人数变化：立即响应
   - 总延迟：< 1.5秒

---

## 📋 更新清单

- ✅ 修改 `getTrendForPair()` 基于 `copied` 不是 `followers`
- ✅ 过滤有效信号（检查 `expiryTime`）
- ✅ 添加趋势刷新定时器（每秒检查过期）
- ✅ 自动移除过期信号
- ✅ 信号跟随时检查是否成为主导信号
- ✅ 在 `onUnmounted` 中清理定时器
- ✅ 控制台日志输出

---

**版本**: 1.0.0  
**更新日期**: 2025-11-24
