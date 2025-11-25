# 📊 M30 信号数据批处理脚本 - 快速开始

## 🚀 一分钟快速开始

### macOS/Linux

```bash
# 生成20条M30信号数据
bash web/scripts/generate-signals.sh generate

# 生成50条Vue格式（可直接粘贴到组件）
bash web/scripts/generate-signals.sh generate -c 50 -f vue

# 显示最新生成的数据
bash web/scripts/generate-signals.sh display

# 复制到剪贴板（macOS only）
bash web/scripts/generate-signals.sh copy
```

### Windows

```batch
REM 生成20条M30信号数据
scripts\generate-signals.bat generate

REM 生成50条Vue格式
scripts\generate-signals.bat generate -c 50 -f vue

REM 显示最新生成的数据
scripts\generate-signals.bat display

REM 交互式生成
scripts\generate-signals.bat insert
```

### Node.js (所有平台)

```bash
cd web

# 生成20条JSON格式
node scripts/generate-signals.js 20 json

# 生成50条Vue格式
node scripts/generate-signals.js 50 vue

# 生成100条CSV格式
node scripts/generate-signals.js 100 csv

# 生成SQL INSERT语句
node scripts/generate-signals.js 100 sql
```

---

## 📁 输出位置

所有生成的文件保存在：`web/generated/`

```
web/generated/
├── signals-20-2025-11-24T13-40-27-545Z.json    ← JSON格式
├── signals-50-2025-11-24T13-40-39-782Z.js      ← Vue格式（推荐）
├── signals-100-2025-11-24T13-40-45-121Z.csv    ← CSV格式
└── signals-100-2025-11-24T13-40-51-749Z.sql    ← SQL格式
```

---

## ✨ 使用场景

### 场景1️⃣：快速测试UI功能

```bash
# 1. 生成20条M30信号
bash web/scripts/generate-signals.sh generate -c 20 -f vue

# 2. 复制Vue代码到剪贴板
bash web/scripts/generate-signals.sh copy

# 3. 粘贴到 TradeView.vue 的 signalFeed 变量
# 4. 保存，自动热更新
# ✅ 信号列表刷新，可测试跟随、下单功能
```

### 场景2️⃣：分析数据（导出Excel）

```bash
# 1. 生成50条CSV格式
bash web/scripts/generate-signals.sh generate -c 50 -f csv

# 2. 打开 web/generated/signals-*.csv
# 3. 用Excel/Google Sheets打开分析
```

### 场景3️⃣：数据库测试

```bash
# 1. 生成100条SQL格式
bash web/scripts/generate-signals.sh generate -c 100 -f sql

# 2. 复制SQL到数据库客户端执行
```

---

## 📊 生成数据特性

✅ **时间框架**：全部为 M30 (30分钟)  
✅ **交易对**：随机选择 38 个交易对  
✅ **信号标题**：20种不同标题  
✅ **技术指标**：20种技术指标组合  
✅ **信心度**：60%-95% 随机分布  
✅ **有效期**：30-60分钟有效  
✅ **数量**：支持1-1000条自定义  

---

## 🎯 支持的输出格式

| 格式 | 用途 | 命令 |
|------|------|------|
| **json** | 通用数据格式 | `node scripts/generate-signals.js 20 json` |
| **vue** | Vue组件代码 | `node scripts/generate-signals.js 20 vue` |
| **csv** | Excel/表格导入 | `node scripts/generate-signals.js 20 csv` |
| **sql** | 数据库导入 | `node scripts/generate-signals.js 20 sql` |
| **js** | JavaScript代码 | `node scripts/generate-signals.js 20 js` |

---

## ❓ 常见问题

**Q: 如何自定义信号数量？**  
A: 在命令中指定 `-c` 或 `--count` 参数：
```bash
bash web/scripts/generate-signals.sh generate -c 100   # 生成100条
```

**Q: 为什么我的信号已过期？**  
A: 脚本基于系统当前时间生成。检查系统时间是否正确。

**Q: 能否生成其他时间框架的信号？**  
A: 本脚本专门生成 M30 数据。要生成其他时间框架，需要修改脚本源代码。

**Q: 生成的数据是否真实？**  
A: 否。这是纯模拟数据用于测试，所有字段随机生成。

---

## 📋 脚本功能对比

| 功能 | Bash脚本 | 批处理脚本 | Node.js脚本 |
|------|---------|---------|-----------|
| 生成信号 | ✅ | ✅ | ✅ |
| 显示预览 | ✅ | ✅ | ✅ |
| 复制剪贴板 | ✅ | ❌ | ❌ |
| 交互式模式 | ✅ | ✅ | ❌ |
| 跨平台 | macOS/Linux | Windows | ✅ |

---

## 🔧 故障排除

```bash
# 添加执行权限 (Linux/Mac)
chmod +x web/scripts/generate-signals.sh
chmod +x web/scripts/generate-signals.js

# 检查Node.js
node -v

# 查看生成的文件
ls -la web/generated/

# 查看最新文件内容
cat web/generated/signals-*.json | head -50
```

---

## 📚 完整文档

更多高级用法、自定义配置、故障排除等详细内容，请查看：

📖 **`web/scripts/README.md`** - 完整使用指南

---

**最后更新:** 2025-11-24  
**版本:** 1.0.0  
**支持系统:** macOS, Linux, Windows
