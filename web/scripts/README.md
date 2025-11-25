# M30 信号数据批处理脚本 - 完整使用指南

## 📋 目录

- [快速开始](#快速开始)
- [脚本介绍](#脚本介绍)
- [使用方法](#使用方法)
- [实际案例](#实际案例)
- [数据格式](#数据格式)
- [故障排除](#故障排除)

---

## 快速开始

### macOS/Linux

```bash
# 生成20条M30信号 (JSON格式)
bash web/scripts/generate-signals.sh generate

# 生成50条M30信号 (Vue格式，可直接粘贴)
bash web/scripts/generate-signals.sh generate -c 50 -f vue

# 显示最新生成的数据
bash web/scripts/generate-signals.sh display

# 复制Vue代码到剪贴板
bash web/scripts/generate-signals.sh copy
```

### Windows

```batch
REM 生成20条M30信号 (JSON格式)
scripts\generate-signals.bat generate

REM 生成50条M30信号 (Vue格式)
scripts\generate-signals.bat generate -c 50 -f vue

REM 显示最新生成的数据
scripts\generate-signals.bat display

REM 交互式生成
scripts\generate-signals.bat insert
```

### Node.js (跨平台)

```bash
cd web

# 生成20条JSON信号
node scripts/generate-signals.js 20 json

# 生成30条CSV信号
node scripts/generate-signals.js 30 csv

# 生成50条SQL INSERT语句
node scripts/generate-signals.js 50 sql
```

---

## 脚本介绍

### 三个脚本的区别

| 脚本 | 系统 | 功能 | 推荐场景 |
|------|------|------|---------|
| `generate-signals.js` | 跨平台 | Node.js核心脚本，生成数据 | 集成到CI/CD、自动化 |
| `generate-signals.sh` | macOS/Linux | Bash包装脚本，功能完整 | 日常开发测试 |
| `generate-signals.bat` | Windows | 批处理脚本，功能完整 | Windows用户 |

### 生成的文件存储位置

所有生成的文件都保存在 `web/generated/` 目录：

```
web/generated/
├── signals-20-2025-11-24T20-30-45-123.json   # JSON格式
├── signals-20-2025-11-24T20-31-12-456.js     # Vue格式（推荐）
├── signals-20-2025-11-24T20-31-45-789.csv    # CSV格式
└── signals-20-2025-11-24T20-32-10-012.sql    # SQL格式
```

---

## 使用方法

### 方法1：Shell脚本 (macOS/Linux)

#### 基本生成

```bash
bash web/scripts/generate-signals.sh generate
```

**输出示例：**
```
✅ Node.js 已就绪 (v18.16.0)

📊 生成 20 条M30信号数据 (json格式)

✅ 数据已保存到: /path/to/web/generated/signals-20-2025-11-24T20-30-45-123.json

📋 统计信息:
   - 信号数量: 20
   - 时间框架: M30 (全部)
   - 交易对数: 18
   - CALL数: 11
   - PUT数: 9
```

#### 自定义数量和格式

```bash
# 生成50条Vue格式信号
bash web/scripts/generate-signals.sh generate -c 50 -f vue

# 生成30条CSV格式信号
bash web/scripts/generate-signals.sh generate -c 30 -f csv

# 生成100条SQL INSERT语句
bash web/scripts/generate-signals.sh generate -c 100 -f sql
```

#### 显示最新数据

```bash
bash web/scripts/generate-signals.sh display
```

**输出示例：**
```
ℹ️  显示最新数据: signals-20-2025-11-24T20-30-45-123.json

[
  {
    "title": "BTCUSDT 突破",
    "metric": "RSI 指标强势",
    "confidence": 0.82,
    "action": "CALL",
    "timing": "m30",
    "symbol": "BTCUSDT",
    "timeframe": "M30",
    "strength": 2,
    "amount": 75,
    "duration": 2400,
    "followers": 1450,
    "copied": 5,
    "createdAt": 1700815845123,
    "validity": 2200000,
    "expiryTime": 1700818045123,
    "isNew": true
  },
  ...
]

ℹ️  完整文件位置: /path/to/web/generated/signals-20-2025-11-24T20-30-45-123.json
```

#### 复制到剪贴板 (macOS only)

```bash
bash web/scripts/generate-signals.sh copy
```

**输出：**
```
✅ Vue 代码已复制到剪贴板
ℹ️  可直接粘贴到 TradeView.vue 的 signalFeed 变量
```

然后可以直接粘贴到 `TradeView.vue` 中：

```vue
const signalFeed = ref([
  { title: 'BTCUSDT 突破', metric: '...', ... },
  { title: 'EURUSD 反弹', metric: '...', ... },
  // 20条信号...
]);
```

#### 交互式生成

```bash
bash web/scripts/generate-signals.sh insert
```

**交互过程：**
```
ℹ️  交互式信号数据插入工具

请输入要生成的信号数量 (默认: 20): 50
选择输出格式:
  1. JSON (推荐用于测试)
  2. CSV (用于电子表格)
  3. SQL (用于数据库)
  4. JavaScript
  5. Vue (推荐用于组件)
请选择 [1-5] (默认: 5): 5

ℹ️  正在生成 50 条 M30 信号数据 (vue 格式)...
```

---

### 方法2：批处理脚本 (Windows)

#### 基本生成

```batch
scripts\generate-signals.bat generate
```

#### 自定义数量和格式

```batch
REM 生成50条Vue格式
scripts\generate-signals.bat generate -c 50 -f vue

REM 生成30条CSV格式
scripts\generate-signals.bat generate -c 30 -f csv
```

#### 显示最新数据

```batch
scripts\generate-signals.bat display
```

#### 交互式生成

```batch
scripts\generate-signals.bat insert
```

---

### 方法3：直接使用Node.js脚本

#### 基础用法

```bash
cd web
node scripts/generate-signals.js [数量] [格式]
```

#### 具体示例

```bash
# 生成20条JSON格式信号
node scripts/generate-signals.js 20 json

# 生成50条CSV格式
node scripts/generate-signals.js 50 csv

# 生成100条SQL INSERT
node scripts/generate-signals.js 100 sql

# 生成Vue格式（推荐）
node scripts/generate-signals.js 30 vue
```

---

## 实际案例

### 案例1：快速测试跟随信号功能

```bash
# 1. 生成20条M30信号
bash web/scripts/generate-signals.sh generate -c 20 -f vue

# 2. 复制到剪贴板
bash web/scripts/generate-signals.sh copy

# 3. 粘贴到 TradeView.vue 中的 signalFeed
# 4. 保存文件，自动热更新
# 5. 页面中出现20条新信号，可进行跟随测试
```

### 案例2：导出到电子表格进行分析

```bash
# 1. 生成50条CSV格式信号
bash web/scripts/generate-signals.sh generate -c 50 -f csv

# 2. 打开 web/generated/ 目录中的 .csv 文件
# 3. 用 Excel/Google Sheets 打开进行分析
```

### 案例3：导入数据库进行后端测试

```bash
# 1. 生成100条SQL格式信号
bash web/scripts/generate-signals.sh generate -c 100 -f sql

# 2. 复制SQL语句到数据库客户端执行
# 3. 测试后端查询、排序、过滤功能
```

### 案例4：自动化CI/CD集成

```bash
# 在CI/CD脚本中调用
node web/scripts/generate-signals.js 50 json > test-signals.json

# 或在package.json中添加脚本
{
  "scripts": {
    "generate:signals": "node web/scripts/generate-signals.js 20 vue",
    "generate:signals:large": "node web/scripts/generate-signals.js 100 json"
  }
}

# 然后运行
npm run generate:signals
```

---

## 数据格式

### JSON 格式

```json
[
  {
    "title": "BTCUSDT 突破",
    "metric": "RSI 指标强势",
    "confidence": 0.82,
    "action": "CALL",
    "timing": "m30",
    "symbol": "BTCUSDT",
    "timeframe": "M30",
    "strength": 2,
    "amount": 75,
    "duration": 2400,
    "followers": 1450,
    "copied": 5,
    "createdAt": 1700815845123,
    "validity": 2200000,
    "expiryTime": 1700818045123,
    "isNew": true
  }
]
```

### CSV 格式

```csv
title,metric,confidence,action,timing,symbol,timeframe,strength,amount,duration,followers,copied,createdAt,validity,expiryTime,isNew
"BTCUSDT 突破","RSI 指标强势",0.82,"CALL","m30","BTCUSDT","M30",2,75,2400,1450,5,1700815845123,2200000,1700818045123,false
```

### SQL 格式

```sql
INSERT INTO signals (title, metric, confidence, action, timing, symbol, timeframe, strength, amount, duration, followers, copied, createdAt, validity, expiryTime, isNew) VALUES ('BTCUSDT 突破', 'RSI 指标强势', 0.82, 'CALL', 'm30', 'BTCUSDT', 'M30', 2, 75, 2400, 1450, 5, 1700815845123, 2200000, 1700818045123, 1);
```

### Vue/JS 格式

```javascript
const signalFeed = ref([
  {
    "title": "BTCUSDT 突破",
    "metric": "RSI 指标强势",
    "confidence": 0.82,
    "action": "CALL",
    "timing": "m30",
    "symbol": "BTCUSDT",
    "timeframe": "M30",
    "strength": 2,
    "amount": 75,
    "duration": 2400,
    "followers": 1450,
    "copied": 5,
    "createdAt": 1700815845123,
    "validity": 2200000,
    "expiryTime": 1700818045123,
    "isNew": true
  }
]);
```

---

## 数据字段说明

| 字段名 | 类型 | 说明 | 范围/示例 |
|--------|------|------|----------|
| `title` | String | 信号标题 | "BTCUSDT 突破" |
| `metric` | String | 技术指标 | "RSI 指标强势", "MACD 金叉" |
| `confidence` | Number | 信心度 | 0.6 - 0.95 |
| `action` | String | 交易方向 | "CALL" 或 "PUT" |
| `timing` | String | 时间框架标签 | "m30" (M30专用) |
| `symbol` | String | 交易对 | "BTCUSDT", "EURUSD" |
| `timeframe` | String | 时间框架 | "M30" (M30专用) |
| `strength` | Number | 信号强度 | 1 或 2 |
| `amount` | Number | 建议下单金额 | 25 - 175 |
| `duration` | Number | 订单时长(秒) | 1800 - 3600 (30m-60m) |
| `followers` | Number | 跟随者数 | 100 - 2000 |
| `copied` | Number | 已被复制次数 | 0 - 15 |
| `createdAt` | Number | 创建时间戳(ms) | Unix timestamp |
| `validity` | Number | 有效期(ms) | 1800000 - 3600000 |
| `expiryTime` | Number | 过期时间戳(ms) | Unix timestamp |
| `isNew` | Boolean | 是否新信号 | true / false |

---

## 故障排除

### 问题1：Node.js 未找到

**错误信息：**
```
❌ 未找到 Node.js，请先安装
```

**解决方案：**
```bash
# 安装 Node.js
# macOS
brew install node

# Ubuntu/Debian
sudo apt-get install nodejs npm

# 验证安装
node -v
npm -v
```

### 问题2：权限不足

**错误信息：**
```
bash: ./generate-signals.sh: Permission denied
```

**解决方案：**
```bash
# 添加执行权限
chmod +x web/scripts/generate-signals.sh
chmod +x web/scripts/generate-signals.js

# 再次运行
bash web/scripts/generate-signals.sh generate
```

### 问题3：找不到生成的文件

**错误信息：**
```
❌ 未找到生成的信号数据
```

**解决方案：**
```bash
# 首先生成数据
bash web/scripts/generate-signals.sh generate

# 检查web/generated目录是否存在
ls -la web/generated/

# 查看最新的文件
ls -lt web/generated/ | head -10
```

### 问题4：时间戳错误

**症状：** 所有信号都已过期

**解决方案：**
脚本使用系统当前时间自动生成。检查系统时间是否正确：

```bash
date

# 如果时间错误，调整系统时间
# macOS
date -s "2025-11-24 20:30:00"

# Linux
sudo date -s "2025-11-24 20:30:00"
```

### 问题5：路径问题

**错误信息：**
```
❌ 文件不存在: /path/to/file
```

**解决方案：**
```bash
# 确保在项目根目录运行
cd /Users/jack/Desktop/GITHUB/PP

# 或使用绝对路径
bash /Users/jack/Desktop/GITHUB/PP/web/scripts/generate-signals.sh generate
```

---

## 高级用法

### 自定义修改脚本

编辑 `web/scripts/generate-signals.js` 修改生成参数：

```javascript
// 修改信号标题
const SIGNAL_TITLES = [
  '自定义标题1',
  '自定义标题2',
  // ...
];

// 修改技术指标
const METRICS = [
  '自定义指标1',
  '自定义指标2',
  // ...
];

// 修改交易对（可选）
const TRADING_PAIRS = [
  { symbol: 'CUSTOM1', name: '自定义交易对1' },
  // ...
];

// 修改数据范围
// 如修改信号数量范围：
amount: randomInt(50, 200),  // 改为 50-200
duration: randomInt(3600, 7200),  // 改为 1h-2h
```

### 集成到npm脚本

编辑 `web/package.json`：

```json
{
  "scripts": {
    "generate:signals": "node scripts/generate-signals.js 20 vue",
    "generate:signals:large": "node scripts/generate-signals.js 100 json",
    "generate:signals:csv": "node scripts/generate-signals.js 50 csv",
    "generate:signals:sql": "node scripts/generate-signals.js 100 sql"
  }
}
```

然后运行：
```bash
npm run generate:signals
npm run generate:signals:large
```

---

## 常见问题 (FAQ)

**Q: 为什么信号的有效期会过期？**
A: 脚本基于系统当前时间生成。若需要长期有效的信号，修改脚本中的 `validity` 参数。

**Q: 能否生成其他时间框架的信号？**
A: 可以。修改脚本中的 `timeframe` 和 `timing` 字段，或创建新的生成函数。

**Q: 生成的数据是否真实？**
A: 生成的是**模拟测试数据**，所有字段都是随机生成。不代表真实行情。

**Q: 能否导入到后端数据库？**
A: 可以。使用 SQL 格式输出，然后在数据库客户端执行。

**Q: 脚本支持哪些操作系统？**
A: Node.js脚本跨平台支持。Bash脚本支持 macOS/Linux。Batch脚本支持 Windows。

---

## 技术支持

如遇到问题，请检查：

1. ✅ Node.js 是否安装（`node -v`）
2. ✅ 文件权限是否正确（`ls -la`）
3. ✅ 系统时间是否正确（`date`）
4. ✅ web/generated 目录是否存在
5. ✅ 运行路径是否正确

更多帮助：

```bash
# 查看脚本帮助
bash web/scripts/generate-signals.sh -h
node web/scripts/generate-signals.js --help
scripts\generate-signals.bat -h
```

---

**Last Updated:** 2025-11-24
**Version:** 1.0.0
