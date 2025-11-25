#!/usr/bin/env node

/**
 * 信号数据生成脚本
 * 用途：为PP Pro Desk生成M30时间框架的测试信号数据
 * 用法：node scripts/generate-signals.js [数量] [输出格式]
 * 
 * 示例：
 *   node scripts/generate-signals.js 20        # 生成20条M30信号，输出为JSON
 *   node scripts/generate-signals.js 50 json   # 生成50条M30信号，输出为JSON
 *   node scripts/generate-signals.js 30 csv    # 生成30条M30信号，输出为CSV
 *   node scripts/generate-signals.js 10 sql    # 生成10条M30信号，输出为SQL INSERT语句
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// ==================== 配置 ====================

const TRADING_PAIRS = [
  // 外汇对
  { symbol: 'EURUSD', name: '欧美' },
  { symbol: 'GBPUSD', name: '英美' },
  { symbol: 'USDJPY', name: '美日' },
  { symbol: 'AUDUSD', name: '澳美' },
  { symbol: 'USDCAD', name: '美加' },
  { symbol: 'NZDUSD', name: '纽美' },
  { symbol: 'EURGBP', name: '欧英' },
  { symbol: 'EURJPY', name: '欧日' },
  { symbol: 'GBPJPY', name: '英日' },
  { symbol: 'AUDJPY', name: '澳日' },
  
  // 加密货币
  { symbol: 'BTCUSDT', name: '比特币' },
  { symbol: 'ETHUSDT', name: '以太坊' },
  { symbol: 'BNBUSDT', name: '币安币' },
  { symbol: 'XRPUSDT', name: '瑞波币' },
  { symbol: 'ADAUSDT', name: '卡尔达诺' },
  { symbol: 'SOLUSDT', name: '索拉纳' },
  { symbol: 'DOGEUSDT', name: '狗狗币' },
  { symbol: 'DOTUSDT', name: '波卡' },
  
  // 贵金属和能源
  { symbol: 'XAUUSD', name: '现货黄金' },
  { symbol: 'XAGUSD', name: '现货白银' },
  { symbol: 'WTIUSD', name: '美原油' },
  { symbol: 'NATGAS', name: '天然气' },
  
  // 指数
  { symbol: 'US500', name: '标普500' },
  { symbol: 'US100', name: '纳斯达克100' },
  { symbol: 'US30', name: '道琼斯30' },
  { symbol: 'DE40', name: '德指DAX' },
  { symbol: 'UK100', name: '英国富时' },
  { symbol: 'JP225', name: '日经225' },
  
  // 个股
  { symbol: 'AAPL', name: '苹果' },
  { symbol: 'MSFT', name: '微软' },
  { symbol: 'GOOGL', name: '谷歌' },
  { symbol: 'AMZN', name: '亚马逊' },
  { symbol: 'TSLA', name: '特斯拉' },
  { symbol: 'META', name: '元宇宙' },
  { symbol: 'NVDA', name: '英伟达' },
  { symbol: 'NFLX', name: '奈飞' },
];

const SIGNAL_TITLES = [
  '突破', '反弹', '拐点', '加速', '回调', '强势', '弱势', '盘整', '加仓', '获利',
  '冲高', '探底', '缩量', '放量', '修复', '衰竭', '启动', '转折', '蓄势', '狂欢'
];

const METRICS = [
  'RSI 指标强势',
  'MACD 金叉',
  'MA 均线突破',
  'Stoch 信号',
  'CCI 极值',
  'ATR 突破',
  'Volume 突增',
  '布林突破',
  'Trend 确认',
  '支撑反弹',
  '阻力突破',
  '角度强势',
  '分时强力',
  '级别共振',
  '多空转变',
  '确认有效',
  '信心强势',
  '破位启动',
  '连续突破',
  '黄金位置'
];

const ACTIONS = ['CALL', 'PUT'];

// ==================== 函数 ====================

/**
 * 生成随机整数
 */
function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

/**
 * 生成随机浮点数
 */
function randomFloat(min, max, decimals = 2) {
  return parseFloat((Math.random() * (max - min) + min).toFixed(decimals));
}

/**
 * 获取随机数组元素
 */
function getRandomItem(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

/**
 * 生成单个信号
 */
function generateSignal(index, baseTime) {
  const pair = getRandomItem(TRADING_PAIRS);
  const titleIndex = index % SIGNAL_TITLES.length;
  const metricIndex = (index + Math.floor(Math.random() * METRICS.length)) % METRICS.length;
  
  const createdAt = baseTime - index * 60000; // 每条间隔1分钟
  const validity = 1800000 + randomInt(0, 1800000); // 30m - 60m
  
  return {
    title: `${pair.symbol} ${SIGNAL_TITLES[titleIndex]}`,
    metric: METRICS[metricIndex],
    confidence: randomFloat(0.6, 0.95, 2),
    action: getRandomItem(ACTIONS),
    timing: 'm30',
    symbol: pair.symbol,
    timeframe: 'M30',
    strength: Math.random() > 0.5 ? 2 : 1,
    amount: randomInt(25, 175),
    duration: randomInt(1800, 3600), // 30m - 60m (seconds)
    followers: randomInt(100, 2000),
    copied: randomInt(0, 15),
    createdAt: createdAt,
    validity: validity,
    expiryTime: createdAt + validity,
    isNew: index === 0,
  };
}

/**
 * 生成多个信号
 */
function generateSignals(count = 20) {
  const now = Date.now();
  const signals = [];
  
  for (let i = 0; i < count; i++) {
    signals.push(generateSignal(i, now));
  }
  
  return signals;
}

/**
 * 转换为JSON格式
 */
function toJSON(signals) {
  return JSON.stringify(signals, null, 2);
}

/**
 * 转换为CSV格式
 */
function toCSV(signals) {
  const headers = [
    'title', 'metric', 'confidence', 'action', 'timing', 'symbol', 'timeframe',
    'strength', 'amount', 'duration', 'followers', 'copied', 'createdAt', 'validity', 'expiryTime', 'isNew'
  ];
  
  const rows = [headers.join(',')];
  
  signals.forEach(signal => {
    const values = headers.map(header => {
      const value = signal[header];
      if (typeof value === 'string') {
        return `"${value}"`;
      }
      return value;
    });
    rows.push(values.join(','));
  });
  
  return rows.join('\n');
}

/**
 * 转换为SQL INSERT语句
 */
function toSQL(signals, tableName = 'signals') {
  const statements = signals.map((signal, index) => {
    const values = [
      `'${signal.title}'`,
      `'${signal.metric}'`,
      signal.confidence,
      `'${signal.action}'`,
      `'${signal.timing}'`,
      `'${signal.symbol}'`,
      `'${signal.timeframe}'`,
      signal.strength,
      signal.amount,
      signal.duration,
      signal.followers,
      signal.copied,
      signal.createdAt,
      signal.validity,
      signal.expiryTime,
      signal.isNew ? 1 : 0,
    ];
    
    return `INSERT INTO ${tableName} (title, metric, confidence, action, timing, symbol, timeframe, strength, amount, duration, followers, copied, createdAt, validity, expiryTime, isNew) VALUES (${values.join(', ')});`;
  });
  
  return statements.join('\n');
}

/**
 * 转换为JavaScript代码片段
 */
function toJavaScript(signals) {
  return `const testSignals = ${JSON.stringify(signals, null, 2)};`;
}

/**
 * 转换为Vue ref格式（可直接粘贴到组件中）
 */
function toVueRef(signals) {
  return `const signalFeed = ref(${JSON.stringify(signals, null, 2)});`;
}

/**
 * 主函数
 */
function main() {
  const args = process.argv.slice(2);
  const count = parseInt(args[0]) || 20;
  const format = (args[1] || 'json').toLowerCase();
  
  console.log(`\n📊 生成 ${count} 条M30信号数据 (${format}格式)\n`);
  
  const signals = generateSignals(count);
  let output = '';
  
  switch (format) {
    case 'json':
      output = toJSON(signals);
      console.log(output);
      break;
      
    case 'csv':
      output = toCSV(signals);
      console.log(output);
      break;
      
    case 'sql':
      output = toSQL(signals);
      console.log(output);
      break;
      
    case 'js':
      output = toJavaScript(signals);
      console.log(output);
      break;
      
    case 'vue':
      output = toVueRef(signals);
      console.log(output);
      break;
      
    default:
      console.error(`❌ 未知格式: ${format}`);
      console.log('支持的格式: json, csv, sql, js, vue');
      process.exit(1);
  }
  
  // 保存到文件
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  const fileExtMap = {
    json: 'json',
    csv: 'csv',
    sql: 'sql',
    js: 'js',
    vue: 'js'
  };
  
  const filename = `signals-${count}-${timestamp}.${fileExtMap[format]}`;
  const filepath = path.join(__dirname, '..', 'generated', filename);
  
  // 创建generated目录
  const generatedDir = path.join(__dirname, '..', 'generated');
  if (!fs.existsSync(generatedDir)) {
    fs.mkdirSync(generatedDir, { recursive: true });
  }
  
  fs.writeFileSync(filepath, output);
  console.log(`\n✅ 数据已保存到: ${filepath}\n`);
  console.log(`📋 统计信息:`);
  console.log(`   - 信号数量: ${count}`);
  console.log(`   - 时间框架: M30 (全部)`);
  console.log(`   - 交易对数: ${new Set(signals.map(s => s.symbol)).size}`);
  console.log(`   - CALL数: ${signals.filter(s => s.action === 'CALL').length}`);
  console.log(`   - PUT数: ${signals.filter(s => s.action === 'PUT').length}`);
  console.log();
}

// ==================== 执行 ====================

main();

export { generateSignals, toJSON, toCSV, toSQL, toJavaScript, toVueRef };
