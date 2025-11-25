# 🎯 PP Pro Desk - M30 信号批处理脚本系统

## 📦 已创建文件清单

### 1️⃣ 核心脚本

```
web/scripts/
├── generate-signals.js          ← Node.js核心脚本 (跨平台)
├── generate-signals.sh          ← Bash脚本 (macOS/Linux)
├── generate-signals.bat         ← 批处理脚本 (Windows)
└── README.md                    ← 完整使用文档
```

### 2️⃣ 快速开始指南

```
SIGNALS_BATCH_QUICK_START.md     ← 1分钟上手指南
SETUP_SCRIPTS_SUMMARY.md         ← 本文档
```

### 3️⃣ 生成的数据目录

```
web/generated/                   ← 所有生成的M30信号数据
├── *.json                       ← JSON格式数据
├── *.js                         ← Vue ref格式（推荐）
├── *.csv                        ← CSV格式（Excel兼容）
└── *.sql                        ← SQL INSERT语句
```

---

## 🚀 快速命令速查

### macOS/Linux

```bash
# 生成20条M30信号（JSON）
bash web/scripts/generate-signals.sh generate

# 生成50条M30信号（Vue格式）
bash web/scripts/generate-signals.sh generate -c 50 -f vue

# 显示最新数据预览
bash web/scripts/generate-signals.sh display

# 复制Vue代码到剪贴板（macOS）
bash web/scripts/generate-signals.sh copy

# 交互式生成
bash web/scripts/generate-signals.sh insert
```

### Windows

```batch
REM 生成20条M30信号
scripts\generate-signals.bat generate

REM 生成50条Vue格式
scripts\generate-signals.bat generate -c 50 -f vue

REM 显示最新数据
scripts\generate-signals.bat display

REM 交互式生成
scripts\generate-signals.bat insert
```

### Node.js（跨平台）

```bash
cd web

# JSON格式
node scripts/generate-signals.js 20 json

# Vue格式（推荐）
node scripts/generate-signals.js 50 vue

# CSV格式
node scripts/generate-signals.js 30 csv

# SQL格式
node scripts/generate-signals.js 100 sql

# JavaScript格式
node scripts/generate-signals.js 25 js
```

---

## 📊 脚本功能特性

### 支持的参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-c, --count` | 生成的信号数量 | `-c 50` 或 `50` |
| `-f, --format` | 输出格式 | `-f json` 或 `json` |

### 支持的输出格式

| 格式 | 后缀 | 用途 | 推荐场景 |
|------|------|------|---------|
| `json` | `.json` | 通用JSON格式 | API测试、通用数据 |
| `vue` | `.js` | Vue ref() 格式 | **直接粘贴到组件** ⭐ |
| `csv` | `.csv` | 逗号分隔值 | Excel/Sheets导入 |
| `sql` | `.sql` | SQL INSERT语句 | 数据库导入 |
| `js` | `.js` | JavaScript常量 | 通用JavaScript项目 |

---

## 💡 实际使用场景

### 🎯 场景1：快速UI测试

```bash
# 1. 生成20条信号（Vue格式）
bash web/scripts/generate-signals.sh generate -c 20 -f vue

# 2. 复制到剪贴板
bash web/scripts/generate-signals.sh copy

# 3. 粘贴到 src/views/TradeView.vue：
#    找到 const signalFeed = ref(...);
#    替换整个数组

# 4. 保存文件，自动热更新
# ✅ 20条M30信号出现在界面
```

### 📊 场景2：数据分析（Excel导出）

```bash
# 1. 生成100条CSV格式
bash web/scripts/generate-signals.sh generate -c 100 -f csv

# 2. 打开 web/generated/signals-*.csv 文件
# 3. 用Excel/Google Sheets打开
# 4. 进行分析、图表制作等
```

### 🗄️ 场景3：后端测试（数据库导入）

```bash
# 1. 生成200条SQL格式
bash web/scripts/generate-signals.sh generate -c 200 -f sql

# 2. 复制SQL语句
# 3. 在数据库客户端执行
# 4. 测试后端查询、排序、分页等功能
```

### 🤖 场景4：自动化集成

在 `package.json` 中添加脚本：

```json
{
  "scripts": {
    "generate:signals": "node web/scripts/generate-signals.js 20 vue",
    "generate:signals:large": "node web/scripts/generate-signals.js 100 json",
    "generate:signals:db": "node web/scripts/generate-signals.js 200 sql"
  }
}
```

运行：
```bash
npm run generate:signals
npm run generate:signals:large
npm run generate:signals:db
```

---

## 📈 生成数据统计

### 每条信号包含字段

```javascript
{
  title: "BTCUSDT 突破",           // 标的+信号标题
  metric: "RSI 指标强势",           // 技术指标
  confidence: 0.82,                 // 信心度 (0.6-0.95)
  action: "CALL" | "PUT",          // 交易方向
  timing: "m30",                    // 时间框架 (固定m30)
  symbol: "BTCUSDT",               // 交易对代码
  timeframe: "M30",                // 时间框架 (固定M30)
  strength: 1 | 2,                 // 信号强度
  amount: 75,                      // 建议金额(25-175)
  duration: 2400,                  // 订单时长(秒) (1800-3600)
  followers: 1450,                 // 跟随者数(100-2000)
  copied: 5,                       // 已复制次数(0-15)
  createdAt: timestamp,            // 创建时间
  validity: 2200000,               // 有效期(ms)
  expiryTime: timestamp,           // 过期时间
  isNew: true | false              // 新信号标记
}
```

### 数据范围

| 字段 | 最小值 | 最大值 | 说明 |
|------|--------|--------|------|
| confidence | 0.60 | 0.95 | 信心度 |
| amount | 25 | 175 | 下单金额 |
| duration | 1800s | 3600s | 30分钟-60分钟 |
| followers | 100 | 2000 | 跟随者数 |
| copied | 0 | 15 | 复制次数 |
| validity | 1800000ms | 3600000ms | 30分钟-60分钟有效 |

---

## ✅ 验证清单

运行以下命令验证环境是否正确配置：

```bash
# ✅ 检查Node.js版本
node -v                    # 应显示 v14+ 版本

# ✅ 检查脚本是否存在
ls -la web/scripts/generate-signals.*

# ✅ 检查脚本权限（Linux/Mac）
ls -la web/scripts/generate-signals.sh

# ✅ 测试生成（生成5条）
bash web/scripts/generate-signals.sh generate -c 5

# ✅ 检查输出目录
ls -la web/generated/

# ✅ 检查是否成功生成
cat web/generated/signals-*.json | head -20
```

---

## 🔧 高级用法

### 修改生成参数

编辑 `web/scripts/generate-signals.js` 中的以下部分自定义数据：

```javascript
// 修改交易对列表
const TRADING_PAIRS = [ ... ];

// 修改信号标题
const SIGNAL_TITLES = [ ... ];

// 修改技术指标
const METRICS = [ ... ];

// 修改数值范围
amount: randomInt(50, 200),        // 改为 50-200
duration: randomInt(3600, 7200),   // 改为 1h-2h
followers: randomInt(1000, 5000),  // 改为 1000-5000
```

### 在CI/CD中使用

```yaml
# GitHub Actions 示例
- name: Generate test signals
  run: |
    cd web
    node scripts/generate-signals.js 100 json > test-signals.json
    
- name: Upload artifacts
  uses: actions/upload-artifact@v2
  with:
    name: test-signals
    path: web/generated/signals-*.json
```

---

## 📚 文档位置

| 文档 | 内容 | 位置 |
|------|------|------|
| **快速开始** | 1分钟上手指南 | `SIGNALS_BATCH_QUICK_START.md` |
| **完整使用** | 详细功能和选项 | `web/scripts/README.md` |
| **本文档** | 脚本系统总结 | `SETUP_SCRIPTS_SUMMARY.md` |

---

## 🆘 故障排除

### 常见错误及解决方案

```bash
# ❌ 错误：permission denied
chmod +x web/scripts/generate-signals.sh
chmod +x web/scripts/generate-signals.js

# ❌ 错误：Node.js not found
# 安装Node.js: https://nodejs.org/

# ❌ 错误：找不到文件
cd /Users/jack/Desktop/GITHUB/PP
ls web/scripts/

# ❌ 错误：生成的数据已过期
# 检查系统时间: date
```

---

## 📞 支持

脚本相关问题：

1. 查看 `web/scripts/README.md` 完整文档
2. 运行 `bash web/scripts/generate-signals.sh -h` 查看帮助
3. 检查 `web/generated/` 中的生成文件

---

## 📋 版本信息

- **创建日期**: 2025-11-24
- **版本**: 1.0.0
- **支持系统**: macOS, Linux, Windows
- **依赖**: Node.js v14+

---

**🎉 现在您已拥有完整的M30信号批处理系统！**

随时使用以下命令生成测试数据：

```bash
bash web/scripts/generate-signals.sh generate -c [数量] -f [格式]
```

祝您测试愉快！🚀

