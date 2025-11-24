<template>
  <div class="page-shell">
    <header class="hero">
      <div class="hero-left">
        <div class="brand-icon">
          <Activity :size="24" color="#5df7c2" />
        </div>
        <div>
          <div class="brand">PP Pro Desk</div>
          <div class="tagline">High-frequency binary options terminal</div>
        </div>
      </div>
      <div class="hero-right">
        <div class="balance-chip">
          <span class="label"><Wallet :size="14" /> Balance</span>
          <span class="value">{{ balanceDisplay }}</span>
        </div>
      </div>
    </header>

    <div class="workspace">
      <section class="chart-pane">
        <div class="chart-glass" @click="handleOutsideClick">
          <div class="chart-toolbar">
            <div class="symbol-block">
              <div class="symbol-select-wrapper">
                <button class="symbol-btn" @click.stop="toggleMenu('symbol')">
                  {{ symbolLabel(selectedSymbol) }}
                  <ChevronDown :size="16" />
                </button>
              </div>
              <div class="price-ticker" :class="{ up: isPriceUp, down: !isPriceUp }">
                {{ formattedPrice }}
                <component :is="isPriceUp ? TrendingUp : TrendingDown" :size="18" />
              </div>
            </div>

            <div class="toolbar-actions">
              <!-- Chart & Timeframe Menu -->
              <div class="tool-wrapper">
                <button class="tool-btn" :class="{ active: activeMenu === 'chart' }" @click.stop="toggleMenu('chart')">
                  <BarChart2 :size="18" />
                  <span class="tool-badge">{{ timeframeLabel(timeframe) }}</span>
                </button>
              </div>

              <!-- Indicators Menu -->
              <div class="tool-wrapper">
                <button class="tool-btn" :class="{ active: activeMenu === 'indicators' }" @click.stop="toggleMenu('indicators')">
                  <Sliders :size="18" />
                </button>
              </div>

              <!-- Drawing Menu -->
              <div class="tool-wrapper">
                <button class="tool-btn" :class="{ active: activeMenu === 'drawing' }" @click.stop="toggleMenu('drawing')">
                  <PenTool :size="18" />
                </button>
              </div>

              <!-- More & Grid -->
              <div class="tool-wrapper">
                <button class="tool-btn">
                  <MoreHorizontal :size="18" />
                </button>
              </div>
              <div class="tool-wrapper">
                <button class="tool-btn">
                  <Grid :size="18" />
                </button>
              </div>
            </div>
          </div>

          <div class="chart-wrapper">
            <div class="chart-overlay top">
              <div class="toolbar-blur">
                <button class="ghost" @click="startDrawing"><Plus :size="14" /> Line</button>
                <button class="ghost" @click="clearDrawings"><Eraser :size="14" /> Clear</button>
              </div>
            </div>
            <div class="chart-surface" ref="chartContainer"></div>
            <div class="chart-overlay bottom">
              <div class="market-stats">
                <div class="stat-item" :class="marketStats.change >= 0 ? 'up' : 'down'">
                  <span class="stat-label">24h Chg</span>
                  <span class="stat-value">
                    {{ marketStats.change >= 0 ? '+' : '' }}{{ formatWithPrecision(marketStats.change) }}
                    ({{ marketStats.changePercent.toFixed(2) }}%)
                  </span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">24h High</span>
                  <span class="stat-value">{{ formatWithPrecision(marketStats.high) }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">24h Low</span>
                  <span class="stat-value">{{ formatWithPrecision(marketStats.low) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Unified Floating Side Panel -->
          <div v-if="activeMenu" class="floating-side-panel" :class="{ 'symbol-panel': activeMenu === 'symbol' }" @click.stop>
            <div class="panel-content">
              <!-- Symbol Selector Content -->
              <div v-if="activeMenu === 'symbol'" class="symbol-menu-content">
                <!-- Search Box -->
                <div class="symbol-search-box">
                  <Search :size="16" class="search-icon" />
                  <input 
                    v-model="symbolSearchQuery" 
                    type="text" 
                    placeholder="搜索" 
                    class="symbol-search-input"
                  />
                </div>

                <!-- Category Tabs -->
                <div class="symbol-categories">
                  <button 
                    v-for="cat in symbolCategories" 
                    :key="cat.id"
                    :class="['category-btn', activeSymbolCategory === cat.id ? 'active' : '']"
                    @click="activeSymbolCategory = cat.id"
                  >
                    <component :is="cat.icon" :size="14" />
                    <span>{{ cat.label }}</span>
                  </button>
                </div>

                <!-- Symbols List -->
                <div class="symbols-list">
                  <button 
                    v-for="sym in filteredSymbols" 
                    :key="sym.symbol"
                    :class="['symbol-item', selectedSymbol === sym.symbol ? 'active' : '']"
                    @click="selectedSymbol = sym.symbol; activeMenu = null;"
                  >
                    <div class="symbol-item-left">
                      <button 
                        class="star-btn"
                        :class="{ favorited: sym.favorited }"
                        @click.stop="toggleFavorite(sym.symbol)"
                      >
                        <Star :size="14" :fill="sym.favorited ? '#ffbe3d' : 'none'" />
                      </button>
                      <div class="symbol-info">
                        <span class="symbol-name">{{ sym.display }}</span>
                        <span class="symbol-return">{{ sym.returnRate }}</span>
                      </div>
                    </div>
                    <Check v-if="selectedSymbol === sym.symbol" :size="16" class="check-icon" />
                  </button>
                </div>
              </div>

              <!-- Chart Settings Content -->
              <div v-if="activeMenu === 'chart'" class="chart-menu-content">
                <div class="menu-section">
                  <div class="menu-label">Chart Type</div>
                  <div class="type-grid">
                    <button 
                      v-for="type in ['line', 'candle', 'area']" 
                      :key="type"
                      :class="['type-btn', chartType === type ? 'active' : '']"
                      @click="chartType = type"
                    >
                      <LineChart v-if="type === 'line'" :size="20" />
                      <BarChart2 v-if="type === 'candle'" :size="20" />
                      <Activity v-if="type === 'area'" :size="20" />
                      <span>{{ type.charAt(0).toUpperCase() + type.slice(1) }}</span>
                    </button>
                    <button class="type-btn disabled" title="Coming soon">
                      <Activity :size="20" />
                      <span>Heikin Ashi</span>
                    </button>
                  </div>
                </div>
                
                <div class="menu-section">
                  <div class="menu-label">Timeframe</div>
                  <div class="tf-grid">
                    <button 
                      v-for="tf in timeframesConfig" 
                      :key="tf.value"
                      :class="['tf-btn', timeframe === tf.value ? 'active' : '']"
                      @click="timeframe = tf.value"
                    >
                      {{ tf.label }}
                    </button>
                  </div>
                </div>

                <div class="menu-section">
                  <div class="menu-label">Settings</div>
                  <div class="settings-list">
                    <label class="setting-row">
                      <span>Enable Timer</span>
                      <div class="toggle" :class="{ active: settings.timer }" @click="settings.timer = !settings.timer">
                        <div class="toggle-handle"></div>
                      </div>
                    </label>
                    <label class="setting-row">
                      <span>Enable Auto Scroll</span>
                      <div class="toggle" :class="{ active: settings.autoScroll }" @click="settings.autoScroll = !settings.autoScroll">
                        <div class="toggle-handle"></div>
                      </div>
                    </label>
                    <label class="setting-row">
                      <span>Enable Grid Snap</span>
                      <div class="toggle" :class="{ active: settings.gridSnap }" @click="settings.gridSnap = !settings.gridSnap">
                        <div class="toggle-handle"></div>
                      </div>
                    </label>
                  </div>
                  <div class="custom-color-link">Custom Candle Color</div>
                </div>
              </div>

              <!-- Indicators Content -->
              <div v-if="activeMenu === 'indicators'" class="indicators-menu-content">
                <div class="indicators-grid">
                  <div 
                    v-for="ind in indicatorsList" 
                    :key="ind" 
                    class="indicator-item"
                    @click="
                      ind === 'Moving Average' ? (showSMA = !showSMA) : 
                      ind === 'Exponential Moving Average' ? (showEMA = !showEMA) : null
                    "
                  >
                    <div class="star-icon">☆</div>
                    <span>{{ ind }}</span>
                    <div v-if="(ind === 'Moving Average' && showSMA) || (ind === 'Exponential Moving Average' && showEMA)" class="active-dot"></div>
                  </div>
                </div>
              </div>

              <!-- Drawing Tools Content -->
              <div v-if="activeMenu === 'drawing'" class="drawing-menu-content">
                <div class="tools-grid">
                  <button 
                    v-for="tool in drawingTools" 
                    :key="tool.id" 
                    class="tool-grid-item"
                    @click="startDrawing(tool.id)"
                    :title="tool.label"
                  >
                    <component :is="tool.icon" :size="20" />
                    <span class="tool-label">{{ tool.label }}</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="card-grid">
          <div class="card debug-card">
            <div class="card-head">
              <span>Live Prices</span>
            </div>
            <div class="card-body debug-body">
              <div v-if="lastTicks.length === 0" class="debug-row">No ticks yet</div>
              <div v-else class="debug-row" v-for="(t, idx) in lastTicks" :key="idx">
                <span class="debug-time">{{ t.timeLabel }}</span>
                <span class="debug-value">{{ t.valueLabel }}</span>
              </div>
            </div>
          </div>

          <div class="card positions-card">
            <div class="card-head tabs">
              <button :class="['tab', currentTab === 'active' ? 'active' : '']" @click="currentTab = 'active'">
                <Activity :size="14" /> Active ({{ activeOrders.length }})
              </button>
              <button :class="['tab', currentTab === 'recent' ? 'active' : '']" @click="currentTab = 'recent'">
                <History :size="14" /> Recent ({{ orderHistory.length }})
              </button>
            </div>
            <div class="card-body positions-body">
              <div v-if="currentTab === 'active'">
                <div v-if="activeOrders.length === 0" class="no-orders">
                  <Layers :size="48" stroke-width="1" />
                  <span>No active trades</span>
                </div>
                <div v-else class="order-item" v-for="order in activeOrders" :key="order.id">
                  <div class="order-header">
                    <span class="symbol">{{ order.asset_symbol }}</span>
                    <span :class="['direction', order.direction.toLowerCase()]">
                      <component :is="order.direction === 'CALL' ? ArrowUpRight : ArrowDownRight" :size="16" />
                      {{ order.direction }}
                    </span>
                  </div>
                  <div class="order-details">
                    <div>Entry: {{ formatPrice(order.open_price, order.asset_symbol) }}</div>
                    <div>Amount: ${{ order.amount }}</div>
                    <div class="timer"><Clock :size="12" /> {{ getTimeLeft(order) }}s</div>
                    <div class="pnl" :class="getPnlClass(order)">
                      Est. {{ getEstimatedPnl(order) }}
                    </div>
                  </div>
                </div>
              </div>
              <div v-else>
                <div class="history-actions">
                  <label>Status</label>
                  <div class="select-wrapper-small">
                    <select v-model="historyStatus">
                      <option value="">All</option>
                      <option value="won">Won</option>
                      <option value="lost">Lost</option>
                      <option value="draw">Draw</option>
                    </select>
                    <ChevronDown :size="12" class="arrow" />
                  </div>
                </div>
                <div v-if="orderHistory.length === 0" class="no-orders">
                  <History :size="48" stroke-width="1" />
                  <span>No history yet</span>
                </div>
                <div v-else class="order-item" v-for="order in orderHistory" :key="order.id">
                  <div class="order-header">
                    <span class="symbol">{{ order.asset_symbol }}</span>
                    <span class="status" :class="statusClass(order.status)">{{ order.status }}</span>
                  </div>
                  <div class="order-details">
                    <div>
                      <strong :class="['direction', order.direction.toLowerCase()]">
                        {{ order.direction }}
                      </strong>
                    </div>
                    <div>Entry: {{ formatPrice(order.open_price, order.asset_symbol) }}</div>
                    <div>Exit: {{ order.close_price ? formatPrice(order.close_price, order.asset_symbol) : '--' }}</div>
                    <div>Amount: ${{ order.amount }}</div>
                    <div>PnL: <span :class="statusClass(order.status)">{{ formatPnl(order) }}</span></div>
                    <div>{{ order.close_time ? new Date(order.close_time).toLocaleTimeString() : '--' }}</div>
                  </div>
                </div>
                <button v-if="orderHistory.length > 0 && marketStore.historyHasMore" class="load-more" @click="loadMoreHistory">
                  Load more
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <aside class="side-pane">
        <div class="card trade-card">
          <div class="card-head">
            <span><Zap :size="16" /> Trade Ticket</span>
            <span class="badge accent">Pro</span>
          </div>
          <div class="card-body">
            <div v-if="errorMsg" class="alert">{{ errorMsg }}</div>
            <div class="input-row">
              <label>Amount ($)</label>
              <div class="input-wrapper">
                <DollarSign :size="14" class="input-icon" />
                <input type="number" v-model="amount" min="1" />
              </div>
            </div>
            <div class="input-row">
              <label>Duration (s)</label>
              <div class="input-wrapper">
                <Clock :size="14" class="input-icon" />
                <select v-model="duration">
                  <option value="30">30s</option>
                  <option value="60">60s</option>
                  <option value="300">5m</option>
                </select>
              </div>
            </div>
            <div class="actions">
              <button class="btn-call" :disabled="!canTrade" @click="handleTrade('CALL')">
                <div class="btn-content">
                  <span>CALL</span>
                  <ArrowUpRight :size="20" />
                </div>
                <span class="payout">{{ (payoutRate * 100).toFixed(0) }}%</span>
              </button>
              <button class="btn-put" :disabled="!canTrade" @click="handleTrade('PUT')">
                <div class="btn-content">
                  <span>PUT</span>
                  <ArrowDownRight :size="20" />
                </div>
                <span class="payout">{{ (payoutRate * 100).toFixed(0) }}%</span>
              </button>
            </div>
            <div class="hint">Est. return: ${{ (amount * (1 + payoutRate)).toFixed(2) }}</div>
          </div>
        </div>

        <div class="card funds-card">
          <div class="card-head">
            <span><Wallet :size="16" /> Funds</span>
            <span class="badge">Wallet</span>
          </div>
          <div class="card-body funds-body">
            <div class="input-row">
              <label>Deposit (USDT)</label>
              <div class="input-group">
                <input type="number" v-model="depositAmount" min="1" />
                <button class="ghost" @click="handleDeposit"><ArrowDownRight :size="14" /> In</button>
              </div>
            </div>
            <div class="input-row">
              <label>Withdraw (USDT)</label>
              <div class="input-group">
                <input type="number" v-model="withdrawAmount" min="1" />
                <button class="ghost" @click="handleWithdraw"><ArrowUpRight :size="14" /> Out</button>
              </div>
            </div>
            <div v-if="fundsMsg" class="funds-msg">{{ fundsMsg }}</div>
            <div v-if="fundsError" class="funds-error">{{ fundsError }}</div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import { createChart, LineSeries, AreaSeries, CandlestickSeries } from 'lightweight-charts';
import { 
  Wallet, Activity, Zap, TrendingUp, TrendingDown, Clock, DollarSign, 
  ArrowUpRight, ArrowDownRight, History, Layers, Settings, Plus, Eraser, 
  RefreshCw, ChevronDown, BarChart2, LineChart, PieChart, ToggleLeft, ToggleRight,
  MoreHorizontal, PenTool, Grid, Sliders, X, Check, Star, Search,
  Minus, MoreVertical, Menu, Fan, Maximize, Square, Move, GitBranch
} from 'lucide-vue-next';
import { useMarketStore } from '../stores/market';
import api from '../api/axios';

const chartContainer = ref(null);
const marketStore = useMarketStore();

const amount = ref(10);
const duration = ref(30);
const isPriceUp = ref(true);
const errorMsg = ref('');
const currentTab = ref('active');
const historyStatus = ref('');
const depositAmount = ref(100);
const withdrawAmount = ref(50);
const fundsMsg = ref('');
const fundsError = ref('');
const chartType = ref('line'); // line | area | candle
const timeframe = ref(5); // seconds per bar
const showSMA = ref(false);
const showEMA = ref(false);
const smaPeriod = ref(10);
const timeframeOptions = [1, 5, 15, 30, 60, 300, 600]; // Keep for logic mapping

// UI State for Menus
const activeMenu = ref(null); // 'chart', 'indicators', 'drawing', 'more', 'symbol'
const toggleMenu = (menu) => {
  activeMenu.value = activeMenu.value === menu ? null : menu;
};

// Symbol Selection State
const symbolSearchQuery = ref('');
const activeSymbolCategory = ref('favorites');

// Trading Pairs Data
const tradingPairs = ref([
  // Crypto
  { symbol: 'BTCUSDT', display: 'BTC/USDT', category: 'crypto', returnRate: '+92%', favorited: true },
  { symbol: 'ETHUSDT', display: 'ETH/USDT', category: 'crypto', returnRate: '+92%', favorited: true },
  { symbol: 'BNBUSDT', display: 'BNB/USDT', category: 'crypto', returnRate: '+92%', favorited: false },
  { symbol: 'XRPUSDT', display: 'XRP/USDT', category: 'crypto', returnRate: '+92%', favorited: false },
  { symbol: 'ADAUSDT', display: 'ADA/USDT', category: 'crypto', returnRate: '+92%', favorited: false },
  { symbol: 'SOLUSDT', display: 'SOL/USDT', category: 'crypto', returnRate: '+92%', favorited: false },
  { symbol: 'DOGEUSDT', display: 'DOGE/USDT', category: 'crypto', returnRate: '+92%', favorited: false },
  { symbol: 'DOTUSDT', display: 'DOT/USDT', category: 'crypto', returnRate: '+92%', favorited: false },
  
  // Forex
  { symbol: 'EURUSD', display: 'EUR/USD', category: 'forex', returnRate: '+92%', favorited: true },
  { symbol: 'GBPUSD', display: 'GBP/USD', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'USDJPY', display: 'USD/JPY', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'AUDUSD', display: 'AUD/USD', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'USDCAD', display: 'USD/CAD', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'NZDUSD', display: 'NZD/USD', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'EURGBP', display: 'EUR/GBP', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'EURJPY', display: 'EUR/JPY', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'GBPJPY', display: 'GBP/JPY', category: 'forex', returnRate: '+92%', favorited: false },
  { symbol: 'AUDJPY', display: 'AUD/JPY', category: 'forex', returnRate: '+92%', favorited: false },
  
  // Commodities
  { symbol: 'XAUUSD', display: 'Gold', category: 'commodities', returnRate: '+92%', favorited: false },
  { symbol: 'XAGUSD', display: 'Silver', category: 'commodities', returnRate: '+92%', favorited: false },
  { symbol: 'WTIUSD', display: 'WTI Oil', category: 'commodities', returnRate: '+92%', favorited: false },
  { symbol: 'NATGAS', display: 'Natural Gas', category: 'commodities', returnRate: '+92%', favorited: false },
  
  // Indices
  { symbol: 'US500', display: 'S&P 500', category: 'indices', returnRate: '+92%', favorited: false },
  { symbol: 'US100', display: 'Nasdaq 100', category: 'indices', returnRate: '+92%', favorited: false },
  { symbol: 'US30', display: 'Dow Jones', category: 'indices', returnRate: '+92%', favorited: false },
  { symbol: 'DE40', display: 'DAX', category: 'indices', returnRate: '+92%', favorited: false },
  { symbol: 'UK100', display: 'FTSE 100', category: 'indices', returnRate: '+92%', favorited: false },
  { symbol: 'JP225', display: 'Nikkei 225', category: 'indices', returnRate: '+92%', favorited: false },
  
  // Stocks
  { symbol: 'AAPL', display: 'Apple', category: 'stocks', returnRate: '+92%', favorited: false },
  { symbol: 'MSFT', display: 'Microsoft', category: 'stocks', returnRate: '+92%', favorited: false },
  { symbol: 'GOOGL', display: 'Google', category: 'stocks', returnRate: '+92%', favorited: false },
  { symbol: 'AMZN', display: 'Amazon', category: 'stocks', returnRate: '+92%', favorited: false },
  { symbol: 'TSLA', display: 'Tesla', category: 'stocks', returnRate: '+92%', favorited: false },
  { symbol: 'META', display: 'Meta', category: 'stocks', returnRate: '+92%', favorited: false },
  { symbol: 'NVDA', display: 'NVIDIA', category: 'stocks', returnRate: '+92%', favorited: false },
]);

const symbolCategories = [
  { id: 'favorites', label: '收藏', icon: Star },
  { id: 'crypto', label: '加密货币', icon: Activity },
  { id: 'forex', label: '外汇', icon: DollarSign },
  { id: 'commodities', label: '大宗商品', icon: Layers },
  { id: 'indices', label: '指数', icon: TrendingUp },
  { id: 'stocks', label: '股票', icon: BarChart2 },
];

const filteredSymbols = computed(() => {
  let symbols = tradingPairs.value;
  
  // Filter by category
  if (activeSymbolCategory.value === 'favorites') {
    symbols = symbols.filter(s => s.favorited);
  } else {
    symbols = symbols.filter(s => s.category === activeSymbolCategory.value);
  }
  
  // Filter by search query
  if (symbolSearchQuery.value.trim()) {
    const query = symbolSearchQuery.value.toLowerCase();
    symbols = symbols.filter(s => 
      s.symbol.toLowerCase().includes(query) || 
      s.display.toLowerCase().includes(query)
    );
  }
  
  return symbols;
});

const toggleFavorite = (symbol) => {
  const pair = tradingPairs.value.find(p => p.symbol === symbol);
  if (pair) {
    pair.favorited = !pair.favorited;
  }
};

// Update allowed symbols for the store
const allowedSymbols = computed(() => tradingPairs.value.map(p => p.symbol));

// Configuration Lists
const timeframesConfig = [
  { label: 'S5', value: 5 }, { label: 'S10', value: 10 }, { label: 'S15', value: 15 }, { label: 'S30', value: 30 },
  { label: 'M1', value: 60 }, { label: 'M2', value: 120 }, { label: 'M3', value: 180 }, { label: 'M5', value: 300 },
  { label: 'M10', value: 600 }, { label: 'M15', value: 900 }, { label: 'M30', value: 1800 }, { label: 'H1', value: 3600 },
  { label: 'H4', value: 14400 }, { label: 'D1', value: 86400 }
];

const indicatorsList = [
  'Accelerator Oscillator', 'ADX', 'Aroon', 'ATR', 'Bear Power', 'Bollinger Bands', 
  'Bull Power', 'CCI', 'DeMarker', 'Envelopes', 'Exponential Moving Average', 'Fractal Chaos Bands', 'Ichimoku Kinko Hyo',
  'MACD', 'Momentum', 'Moving Average', 'OsMA', 'Parabolic SAR', 'RSI', 'Stochastic',
  'SuperTrend', 'Vortex', 'Williams %R', 'Zig Zag'
];

const drawingTools = [
  { id: 'horizontal', label: 'Horizontal Line', icon: Minus },
  { id: 'vertical', label: 'Vertical Line', icon: MoreVertical },
  { id: 'ray', label: 'Ray', icon: ArrowUpRight },
  { id: 'fib_retrace', label: 'Fibonacci Retracement', icon: Menu },
  { id: 'fib_fan', label: 'Fibonacci Fan', icon: Fan },
  { id: 'trend', label: 'Trend Line', icon: TrendingUp },
  { id: 'channel', label: 'Parallel Channel', icon: Maximize },
  { id: 'rect', label: 'Rectangle', icon: Square },
  { id: 'pitchfork', label: "Andrew's Pitchfork", icon: GitBranch },
];

// Settings State
const settings = ref({
  timer: true,
  autoScroll: true,
  gridSnap: true
});

let currentBar = null; // { time, open, high, low, close }

let chart;
let series;
let smaSeries;
let emaSeries;
let resizeObserver;
let timerInterval;
let historyInterval;
let candleInterval;
let drawings = [];
let drawingMode = false;

const interpolatedPrice = ref(0);
let animationFrameId = null;

const currentPrice = computed(() => interpolatedPrice.value || marketStore.currentPrice);
const activeOrders = computed(() => marketStore.activeOrders);
const orderHistory = computed(() => marketStore.orderHistory);
const balanceDisplay = computed(() => `${marketStore.balanceCurrency} ${marketStore.balance.toFixed(2)}`);
const selectedSymbol = computed({
  get: () => marketStore.selectedSymbol,
  set: (val) => marketStore.setSymbol(val),
});

// Market Stats Calculations
const marketStats = computed(() => {
  const history = marketStore.priceHistory;
  if (!history || history.length === 0) return { change: 0, changePercent: 0, high: 0, low: 0 };

  const current = currentPrice.value || 0;
  // Use the oldest available point as reference for "Open" in this session context
  // In a real app, this would come from a 24h ticker API
  const open = history[history.length - 1]?.value || current;
  const change = current - open;
  const changePercent = open !== 0 ? (change / open) * 100 : 0;

  // Calculate High/Low from available history
  let high = current;
  let low = current;
  for (const p of history) {
    if (p.value > high) high = p.value;
    if (p.value < low) low = p.value;
  }

  return {
    change,
    changePercent,
    high,
    low
  };
});

const candleCountdown = ref('--');

const updateCandleCountdown = () => {
  const now = Math.floor(Date.now() / 1000);
  const tf = timeframe.value;
  const nextCandleTime = (Math.floor(now / tf) + 1) * tf;
  const diff = nextCandleTime - now;
  
  const m = Math.floor(diff / 60);
  const s = diff % 60;
  candleCountdown.value = `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
};

watch(currentPrice, (newVal) => {
  if (newVal) {
    document.title = `${newVal.toFixed(pricePrecision.value)} | ${selectedSymbol.value} | PP Pro`;
  }
});

const symbolLabel = (sym) => {
  const pair = tradingPairs.value.find(p => p.symbol === sym);
  return pair ? pair.display : sym;
};

const precisionMap = {
  // Crypto
  BTCUSDT: 2,
  ETHUSDT: 2,
  BNBUSDT: 2,
  XRPUSDT: 4,
  ADAUSDT: 4,
  SOLUSDT: 2,
  DOGEUSDT: 6,
  DOTUSDT: 3,
  // Forex
  EURUSD: 5,
  GBPUSD: 5,
  USDJPY: 3,
  AUDUSD: 5,
  USDCAD: 5,
  NZDUSD: 5,
  EURGBP: 5,
  EURJPY: 3,
  GBPJPY: 3,
  AUDJPY: 3,
  // Commodities
  XAUUSD: 2,
  XAGUSD: 3,
  WTIUSD: 2,
  NATGAS: 3,
  // Indices
  US500: 2,
  US100: 2,
  US30: 2,
  DE40: 2,
  UK100: 2,
  JP225: 2,
  // Stocks
  AAPL: 2,
  MSFT: 2,
  GOOGL: 2,
  AMZN: 2,
  TSLA: 2,
  META: 2,
  NVDA: 2,
};

const pricePrecision = computed(() => {
  return precisionMap[selectedSymbol.value] || 4;
});

const formattedPrice = computed(() => {
  if (!currentPrice.value) return '--';
  return currentPrice.value.toFixed(pricePrecision.value);
});

const canTrade = computed(() => {
  return amount.value > 0 && currentPrice.value > 0;
});

const payoutRate = 0.85;

const timeframeLabel = (sec) => {
  const tf = timeframesConfig.find(t => t.value === sec);
  if (tf) return tf.label;
  
  if (sec < 60) return `S${sec}`;
  const minutes = sec / 60;
  return `M${minutes}`;
};

const getTimeLeft = (order) => {
  const closeTime = new Date(order.close_time).getTime();
  const now = Date.now();
  return Math.max(0, Math.floor((closeTime - now) / 1000));
};

const getEstimatedPnl = (order) => {
  const price = currentPrice.value;
  if (!price) return '...';
  const isWin = order.direction === 'CALL' ? price > order.open_price : price < order.open_price;
  if (isWin) return `+$${(order.amount * order.payout_rate).toFixed(2)}`;
  return `-$${order.amount}`;
};

const getPnlClass = (order) => {
  const price = currentPrice.value;
  if (!price) return '';
  const isWin = order.direction === 'CALL' ? price > order.open_price : price < order.open_price;
  return isWin ? 'win' : 'loss';
};

const statusClass = (status) => {
  if (status === 'won') return 'win';
  if (status === 'lost') return 'loss';
  if (status === 'draw') return 'draw';
  if (status === 'active' || status === 'pending') return 'pending';
  return '';
};

const formatPrice = (price, symbol) => {
  const prec = precisionMap[symbol] || 4;
  return price ? Number(price).toFixed(prec) : '--';
};

const formatWithPrecision = (value) => {
  if (value === null || value === undefined) return '--';
  return Number(value).toFixed(pricePrecision.value);
};

const formatPnl = (order) => {
  if (order.status === 'won') return `+$${order.pnl.toFixed(2)}`;
  if (order.status === 'lost') return `-$${Math.abs(order.pnl).toFixed(2)}`;
  if (order.status === 'draw') return '$0.00';
  return '--';
};

const lastTicks = computed(() => {
  const nowSec = Math.floor(Date.now() / 1000);
  const recent = marketStore.priceHistory.filter((p) => p.time >= nowSec - 180); // show only last 3 minutes
  const tail = recent.slice(-5);
  return tail
    .map((p) => ({
      timeLabel: new Date(p.time * 1000).toLocaleTimeString(),
      valueLabel: p.value?.toFixed(pricePrecision.value) ?? '--',
    }))
    .reverse();
});

const processTick = (price, time) => {
  if (!price || !time) return;
  
  const tf = timeframe.value;
  const barTime = Math.floor(time / tf) * tf;
  
  if (!currentBar || barTime > currentBar.time) {
    // New bar
    currentBar = {
      time: barTime,
      open: price,
      high: price,
      low: price,
      close: price
    };
  } else {
    // Update current bar
    currentBar.high = Math.max(currentBar.high, price);
    currentBar.low = Math.min(currentBar.low, price);
    currentBar.close = price;
  }

  if (series) {
    const item = chartType.value === 'candle' 
      ? currentBar 
      : { time: currentBar.time, value: currentBar.close };
      
    try {
      series.update(item);
    } catch (e) {
      // ignore time order errors
    }
  }
};

const resampleHistory = (history, tf) => {
  const buckets = new Map();
  for (const p of history) {
    const t = p.time;
    if (!t) continue;
    const bucketTime = Math.floor(t / tf) * tf;
    if (!buckets.has(bucketTime)) {
      buckets.set(bucketTime, { 
        time: bucketTime, 
        open: p.value, 
        high: p.value, 
        low: p.value, 
        close: p.value 
      });
    } else {
      const b = buckets.get(bucketTime);
      b.high = Math.max(b.high, p.value);
      b.low = Math.min(b.low, p.value);
      b.close = p.value;
    }
  }
  return Array.from(buckets.values()).sort((a, b) => a.time - b.time);
};

const aggregateCandles = (history, bucketSec) => {
  const buckets = new Map();
  for (const point of history) {
    const t = point.time;
    if (typeof t !== 'number' || Number.isNaN(t)) continue;
    const bucket = Math.floor(t / bucketSec) * bucketSec;
    if (!buckets.has(bucket)) {
      buckets.set(bucket, { open: point.value, high: point.value, low: point.value, close: point.value, time: bucket });
    } else {
      const b = buckets.get(bucket);
      b.high = Math.max(b.high, point.value);
      b.low = Math.min(b.low, point.value);
      b.close = point.value;
    }
  }
  return Array.from(buckets.values()).sort((a, b) => a.time - b.time);
};

const computeSMA = (data, period) => {
  const res = [];
  const vals = [];
  data.forEach((p) => {
    const val = p.close ?? p.value;
    vals.push(val);
    if (vals.length >= period) {
      const window = vals.slice(-period);
      const sum = window.reduce((a, b) => a + b, 0);
      res.push({ time: p.time, value: sum / period });
    }
  });
  return res;
};

const computeEMA = (data, period) => {
  if (period <= 1) return [];
  const res = [];
  let emaPrev = null;
  const k = 2 / (period + 1);
  data.forEach((p) => {
    const val = p.close ?? p.value;
    if (emaPrev === null) {
      emaPrev = val;
    } else {
      emaPrev = val * k + emaPrev * (1 - k);
    }
    res.push({ time: p.time, value: emaPrev });
  });
  return res.slice(period - 1);
};

const applyHistoryToSeries = (history, updatePrice = true) => {
  if (!series || history.length === 0) return;
  if (chartType.value === 'candle') return; // handled by candles fetch

  const data = resampleHistory(history, timeframe.value);
  if (data.length === 0) return;

  // Map to Line/Area format
  const lineData = data.map(d => ({ time: d.time, value: d.close }));

  if (updatePrice) {
    try {
      series.setData(lineData);
      // Initialize currentBar with the last data point
      const last = data[data.length - 1];
      currentBar = { ...last };
    } catch (e) {
      console.warn('Failed to setData on series', e);
    }

    if (lineData.length > 1) {
      const latest = lineData[lineData.length - 1];
      isPriceUp.value = latest.value >= lineData[lineData.length - 2].value;
    }
  }

  if (smaSeries) {
    const data = showSMA.value ? computeSMA(lineData, smaPeriod.value) : [];
    smaSeries.setData(data);
  }
  if (emaSeries) {
    const data = showEMA.value ? computeEMA(lineData, smaPeriod.value) : [];
    emaSeries.setData(data);
  }
};

const renderCandles = (updatePrice = true) => {
  // Only render candles when candle mode is active.
  if (!series || chartType.value !== 'candle') return;

  const candles = (marketStore.candles || []).filter(
    (c) =>
      c &&
      c.time != null &&
      c.open != null &&
      c.high != null &&
      c.low != null &&
      c.close != null
  );

  if (updatePrice) {
    series.setData(candles);

    if (candles.length > 1) {
      const last = candles[candles.length - 1];
      const prev = candles[candles.length - 2];
      isPriceUp.value = last.close >= prev.close;
    }
  }

  if (smaSeries) {
    const data = showSMA.value ? computeSMA(candles, smaPeriod.value) : [];
    smaSeries.setData(data);
  }
  if (emaSeries) {
    const data = showEMA.value ? computeEMA(candles, smaPeriod.value) : [];
    emaSeries.setData(data);
  }
};

const createSeries = () => {
  if (!chart) return;
  if (series) {
    chart.removeSeries(series);
    series = null;
  }
  if (smaSeries) {
    chart.removeSeries(smaSeries);
    smaSeries = null;
  }
  if (emaSeries) {
    chart.removeSeries(emaSeries);
    emaSeries = null;
  }
  const baseOptions = {
    priceFormat: { type: 'price', precision: pricePrecision.value, minMove: 1 / Math.pow(10, pricePrecision.value) },
  };
  if (chartType.value === 'candle') {
    series = chart.addSeries(CandlestickSeries, {});
    smaSeries = chart.addSeries(LineSeries, { color: '#ffeb3b', lineWidth: 1, ...baseOptions });
    emaSeries = chart.addSeries(LineSeries, { color: '#ff9800', lineWidth: 1, ...baseOptions });
  } else if (chartType.value === 'area') {
    series = chart.addSeries(AreaSeries, { lineColor: '#4caf50', topColor: 'rgba(76,175,80,0.3)', bottomColor: 'rgba(76,175,80,0.05)', ...baseOptions });
    smaSeries = chart.addSeries(LineSeries, { color: '#ffeb3b', lineWidth: 1, ...baseOptions });
    emaSeries = chart.addSeries(LineSeries, { color: '#ff9800', lineWidth: 1, ...baseOptions });
  } else {
    series = chart.addSeries(LineSeries, { color: '#2962FF', lineWidth: 2, ...baseOptions });
    smaSeries = chart.addSeries(LineSeries, { color: '#ffeb3b', lineWidth: 1, ...baseOptions });
    emaSeries = chart.addSeries(LineSeries, { color: '#ff9800', lineWidth: 1, ...baseOptions });
  }
};

const refreshCandles = async () => {
  if (chartType.value !== 'candle') return;
  await marketStore.fetchCandles({ symbol: selectedSymbol.value, interval: timeframe.value, limit: 300 });
  renderCandles();
};

const fetchHistory = async (reset = true) => {
  const params = {
    status: historyStatus.value,
    limit: 20,
    offset: reset ? 0 : marketStore.orderHistory.length,
    append: !reset,
  };
  await marketStore.fetchOrderHistory(params);
};

const loadMoreHistory = async () => {
  if (!marketStore.historyHasMore) return;
  await fetchHistory(false);
};

watch(
  () => marketStore.priceHistory,
  (history) => applyHistoryToSeries(history),
  { deep: true }
);

watch(
  () => selectedSymbol.value,
  () => {
    interpolatedPrice.value = 0;
    if (series) {
      createSeries();
    }
    // Update chart price formatter when symbol changes
    if (chart) {
      chart.applyOptions({
        localization: {
          priceFormatter: (price) => {
            return Number(price).toFixed(pricePrecision.value);
          },
        },
      });
    }
    marketStore.fetchActiveOrders();
    fetchHistory(true);
    refreshCandles();
  }
);

watch(
  () => historyStatus.value,
  () => {
    fetchHistory(true);
  }
);

const updateTimeScale = () => {
  if (!chart || !chartContainer.value) return;
  
  chart.timeScale().applyOptions({
    minBarSpacing: 0.5,
    barSpacing: 6,
  });
};

watch(
  () => timeframe.value,
  () => {
    marketStore.setTimeframe(timeframe.value);
    currentBar = null; // Reset current bar on timeframe change
    refreshCandles();
    if (chartType.value !== 'candle') {
      applyHistoryToSeries(marketStore.priceHistory);
    }
    restartCandleInterval();
    updateTimeScale();
  }
);

const animatePrice = () => {
  const target = marketStore.currentPrice;
  if (!target) {
    animationFrameId = null;
    return;
  }

  if (interpolatedPrice.value === 0) {
    interpolatedPrice.value = target;
  }

  const diff = target - interpolatedPrice.value;
  if (Math.abs(diff) < 0.00001) {
    interpolatedPrice.value = target;
  } else {
    // Smooth interpolation factor (adjust 0.1 for speed)
    interpolatedPrice.value += diff * 0.05;
  }

  // Update chart with interpolated price
  processTick(interpolatedPrice.value, Date.now() / 1000);

  animationFrameId = requestAnimationFrame(animatePrice);
};

watch(
  () => marketStore.currentPrice,
  (newPrice) => {
    if (newPrice && !animationFrameId) {
      animatePrice();
    }
  },
  { immediate: true }
);

watch(
  () => chartType.value,
  () => {
    createSeries();
    // render from either candles or tick history
    if (chartType.value === 'candle') {
      refreshCandles();
    } else {
      applyHistoryToSeries(marketStore.priceHistory);
    }
    updateTimeScale();
  }
);

watch(
  () => [showSMA.value, showEMA.value, smaPeriod.value],
  () => {
    if (chartType.value === 'candle') {
      renderCandles(false);
    } else {
      applyHistoryToSeries(marketStore.priceHistory, false);
    }
  }
);

onMounted(async () => {
  await nextTick();
  if (!chartContainer.value) return;

  chart = createChart(chartContainer.value, {
    width: chartContainer.value.clientWidth || 800,
    height: chartContainer.value.clientHeight || 600,
    layout: { background: { type: 'solid', color: '#1e222d' }, textColor: '#d1d4dc' },
    grid: { vertLines: { color: '#2b2b43' }, horzLines: { color: '#2b2b43' } },
    timeScale: {
      timeVisible: true,
      secondsVisible: true,
      tickMarkFormatter: (time) => {
        const date = new Date(time * 1000);
        const h = date.getHours().toString().padStart(2, '0');
        const m = date.getMinutes().toString().padStart(2, '0');
        const s = date.getSeconds().toString().padStart(2, '0');
        return `${h}:${m}:${s}`;
      },
    },
    rightPriceScale: {
      scaleMargins: {
        top: 0.2,
        bottom: 0.2,
      },
      borderVisible: false,
      minimumWidth: 80,
    },
    localization: {
      priceFormatter: (price) => {
        return Number(price).toFixed(pricePrecision.value);
      },
    },
  });

  createSeries();
  renderCandles();
  applyHistoryToSeries(marketStore.priceHistory);

  resizeObserver = new ResizeObserver((entries) => {
    if (!chartContainer.value || entries.length === 0 || entries[0].target !== chartContainer.value) return;
    const { width, height } = entries[0].contentRect;
    if (width > 0 && height > 0) chart.applyOptions({ width, height });
  });
  resizeObserver.observe(chartContainer.value);

  marketStore.connect();
  marketStore.fetchActiveOrders();
  marketStore.fetchBalance();
  refreshCandles();
  fetchHistory(true);

  timerInterval = setInterval(() => {
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
    updateCandleCountdown();
  }, 1000);

  historyInterval = setInterval(() => {
    fetchHistory(true);
  }, 5000);

  restartCandleInterval();

  chart.subscribeClick((param) => {
    if (!drawingMode) return;
    const price = currentPrice.value || (param.seriesData && param.seriesData.get(series)?.close);
    if (!price) return;
    const line = series.createPriceLine({ price, color: '#ffa500', lineWidth: 1 });
    drawings.push(line);
    drawingMode = false;
  });

  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  marketStore.disconnect();
  if (chart) chart.remove();
  if (resizeObserver && chartContainer.value) resizeObserver.unobserve(chartContainer.value);
  clearInterval(timerInterval);
  clearInterval(historyInterval);
  clearInterval(candleInterval);
  if (animationFrameId) cancelAnimationFrame(animationFrameId);
  window.removeEventListener('resize', handleResize);
  clearDrawings();
});

const handleResize = () => {
  if (chart && chartContainer.value) {
    chart.applyOptions({
      width: chartContainer.value.clientWidth,
      height: chartContainer.value.clientHeight,
    });
    updateTimeScale();
  }
};

const restartCandleInterval = () => {
  if (candleInterval) clearInterval(candleInterval);
  const ms = Math.max(timeframe.value * 1000, 1000);
  candleInterval = setInterval(() => {
    refreshCandles();
  }, ms);
};

const handleTrade = async (direction) => {
  errorMsg.value = '';
  try {
    await marketStore.placeOrder(selectedSymbol.value, direction, amount.value, duration.value);
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
  } catch (error) {
    errorMsg.value = error.response?.data?.error || error.message || 'Failed to place order';
  }
};

const handleDeposit = async () => {
  fundsMsg.value = '';
  fundsError.value = '';
  try {
    const res = await api.post('/funds/deposit', { amount: depositAmount.value, currency: 'USDT' });
    fundsMsg.value = `Deposit request submitted: #${res.data.request.id}`;
    marketStore.fetchBalance();
  } catch (err) {
    fundsError.value = err.response?.data?.error || err.message;
  }
};

const startDrawing = (toolId = 'trend') => {
  drawingMode = true;
  // In a real implementation, we would set the specific tool type here
  console.log('Starting drawing with tool:', toolId);
};

const clearDrawings = () => {
  if (drawings.length && series) {
    drawings.forEach((l) => series.removePriceLine(l));
    drawings = [];
  }
};

const handleWithdraw = async () => {
  fundsMsg.value = '';
  fundsError.value = '';
  try {
    const res = await api.post('/funds/withdraw', { amount: withdrawAmount.value, currency: 'USDT' });
    fundsMsg.value = `Withdraw request submitted: #${res.data.request.id}`;
    marketStore.fetchBalance();
  } catch (err) {
    fundsError.value = err.response?.data?.error || err.message;
  }
};

const handleOutsideClick = (e) => {
  // Close panel when clicking outside of it
  // Check if click is on chart-glass or any child that's not the panel
  if (activeMenu.value && !e.target.closest('.floating-side-panel')) {
    activeMenu.value = null;
  }
};
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  background: radial-gradient(circle at 20% 20%, rgba(108, 99, 255, 0.08), transparent 40%),
    radial-gradient(circle at 80% 10%, rgba(0, 214, 170, 0.08), transparent 40%),
    radial-gradient(circle at 50% 80%, rgba(255, 111, 97, 0.05), transparent 40%),
    #0b0e14;
  color: #e7ecf5;
  padding: 24px;
  font-family: 'Inter', 'SF Pro Display', system-ui, -apple-system, sans-serif;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(20px);
}

.hero-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.brand-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: rgba(93, 247, 194, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(93, 247, 194, 0.2);
}

.brand {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.5px;
  background: linear-gradient(90deg, #fff, #a0aab9);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.tagline {
  color: #6b7280;
  font-size: 13px;
  margin-top: 2px;
}

.hero-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.balance-chip {
  padding: 8px 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(11, 14, 20, 0.6);
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 140px;
  transition: all 0.2s;
}

.balance-chip:hover {
  border-color: rgba(255, 255, 255, 0.15);
  background: rgba(11, 14, 20, 0.8);
}

.balance-chip .label {
  font-size: 11px;
  color: #8fa1c4;
  display: flex;
  align-items: center;
  gap: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 600;
}

.balance-chip .value {
  font-weight: 700;
  font-size: 18px;
  color: #fff;
  font-variant-numeric: tabular-nums;
}

.workspace {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 20px;
  flex: 1;
}

.chart-pane {
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-height: 0;
}

.chart-glass {
  background: rgba(18, 20, 28, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  padding: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 500px;
}

.chart-toolbar {
  position: absolute;
  top: 24px;
  left: 24px;
  right: 24px;
  z-index: 20;
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  pointer-events: none;
}

.toolbar-actions {
  pointer-events: auto;
  background: rgba(12, 16, 27, 0.8);
  backdrop-filter: blur(12px);
  padding: 4px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  gap: 4px;
  align-items: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.symbol-block {
  pointer-events: auto;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 16px;
  background: rgba(12, 16, 27, 0.8);
  backdrop-filter: blur(12px);
  padding: 8px 12px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.toolbar-actions {
  align-items: center;
}

.symbol-select-wrapper {
  position: relative;
  display: inline-block;
}

.symbol-btn {
  appearance: none;
  background: transparent;
  border: none;
  color: #fff;
  padding: 4px 8px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.symbol-btn:hover {
  background: rgba(255, 255, 255, 0.08);
}

.floating-side-panel.symbol-panel {
  width: 360px;
  left: 24px;
  right: auto;
}

.symbol-search-box {
  position: relative;
  margin-bottom: 12px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7a99;
  pointer-events: none;
}

.symbol-search-input {
  width: 100%;
  padding: 10px 12px 10px 36px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  font-size: 13px;
  outline: none;
  transition: all 0.2s;
}

.symbol-search-input::placeholder {
  color: #6b7a99;
}

.symbol-search-input:focus {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
}

.symbol-categories {
  display: flex;
  gap: 6px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.category-btn {
  flex: 1;
  min-width: 80px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #8fa1c4;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.category-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
  color: #dbe3f5;
}

.category-btn.active {
  background: rgba(93, 247, 194, 0.15);
  border-color: rgba(93, 247, 194, 0.5);
  color: #5df7c2;
}

.symbol-menu-content .symbols-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 4px;
}

.symbol-menu-content .symbols-list::-webkit-scrollbar {
  width: 6px;
}

.symbol-menu-content .symbols-list::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 3px;
}

.symbol-menu-content .symbols-list::-webkit-scrollbar-thumb {
  background: rgba(93, 247, 194, 0.3);
  border-radius: 3px;
}

.symbol-menu-content .symbols-list::-webkit-scrollbar-thumb:hover {
  background: rgba(93, 247, 194, 0.5);
}

.symbol-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  color: #dbe3f5;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.symbol-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
}

.symbol-item.active {
  background: rgba(93, 247, 194, 0.1);
  border-color: rgba(93, 247, 194, 0.5);
  color: #5df7c2;
}

.symbol-item-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.star-btn {
  background: none;
  border: none;
  color: #6b7a99;
  cursor: pointer;
  padding: 2px;
  display: flex;
  align-items: center;
  transition: all 0.2s;
}

.star-btn:hover {
  color: #ffbe3d;
  transform: scale(1.1);
}

.star-btn.favorited {
  color: #ffbe3d;
}

.symbol-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
}

.symbol-name {
  font-weight: 600;
  font-size: 13px;
}

.symbol-return {
  font-size: 11px;
  color: #5df7c2;
  font-weight: 500;
}

.check-icon {
  color: #5df7c2;
}

.price-ticker {
  font-size: 24px;
  font-weight: 700;
  padding: 6px 12px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-variant-numeric: tabular-nums;
}

.price-ticker.up {
  color: #5df7c2;
  text-shadow: 0 0 20px rgba(93, 247, 194, 0.3);
}
.price-ticker.down {
  color: #ff7b7b;
  text-shadow: 0 0 20px rgba(255, 123, 123, 0.3);
}

.market-stats {
  pointer-events: auto;
  display: flex;
  gap: 16px;
  font-size: 12px;
  padding: 8px 12px;
  background: rgba(12, 16, 27, 0.8);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-label {
  font-size: 10px;
  color: #6b7a99;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-weight: 600;
}

.stat-value {
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: #dbe3f5;
}

.stat-item.up .stat-value {
  color: #5df7c2;
}
.stat-item.down .stat-value {
  color: #ff7b7b;
}

.toolbar-right {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.chip-group {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 10px;
}

.chip {
  padding: 8px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #8fa1c4;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chip:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
}

.chip.active {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.indicator-toggle {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  font-size: 12px;
  color: #8fa1c4;
}

.indicator-toggle label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.indicator-toggle input[type='checkbox'] {
  accent-color: #5f9bff;
}

.period-input {
  width: 50px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  border-radius: 6px;
  padding: 4px 8px;
  text-align: center;
  font-size: 12px;
}

.chart-wrapper {
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: #0c101b;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chart-surface {
  width: 100%;
  height: 100%;
  flex: 1;
  min-height: 400px;
}

.chart-overlay {
  position: absolute;
  left: 0;
  right: 0;
  padding: 10px;
  display: flex;
  justify-content: space-between;
  pointer-events: none;
  z-index: 10;
}
.chart-overlay.top {
  top: 0;
  justify-content: flex-start;
}
.chart-overlay.bottom {
  bottom: 0;
  justify-content: flex-start;
}

.toolbar-blur {
  pointer-events: all;
  display: flex;
  gap: 8px;
  padding: 8px 10px;
  background: rgba(12, 16, 27, 0.7);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.ghost {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #dbe3f5;
  padding: 6px 10px;
  border-radius: 8px;
  cursor: pointer;
}

.legend {
  pointer-events: all;
  padding: 6px 8px;
  background: rgba(12, 16, 27, 0.8);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  gap: 10px;
  align-items: center;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  margin-right: 6px;
}
.legend-dot.up {
  background: #5df7c2;
}
.legend-dot.down {
  background: #ff7b7b;
}

.debug-mini {
  pointer-events: all;
  padding: 6px 10px;
  background: rgba(12, 16, 27, 0.8);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  gap: 6px;
  align-items: center;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.card {
  background: rgba(18, 20, 28, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 700;
  font-size: 15px;
  color: #fff;
}

.card-head span:first-child {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.badge {
  padding: 4px 10px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.08);
  font-size: 11px;
  letter-spacing: 0.5px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  font-weight: 600;
  text-transform: uppercase;
}
.badge.accent {
  background: linear-gradient(135deg, #ff7b7b, #ffd257);
  color: #1b0f20;
  border: none;
  box-shadow: 0 4px 12px rgba(255, 123, 123, 0.3);
}

.tabs {
  gap: 8px;
  justify-content: flex-start;
}
.tab {
  background: transparent;
  border: 1px solid transparent;
  padding: 8px 12px;
  border-radius: 10px;
  cursor: pointer;
  color: #8fa1c4;
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}
.tab:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
}
.tab.active {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border-color: rgba(255, 255, 255, 0.1);
}

.positions-body {
  gap: 16px;
  max-height: 400px;
  overflow-y: auto;
  padding-right: 4px;
}

.order-item {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 12px;
  transition: all 0.2s;
}
.order-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-weight: 700;
  font-size: 14px;
}
.order-details {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  color: #8fa1c4;
  font-size: 12px;
}
.direction {
  display: flex;
  align-items: center;
  gap: 4px;
}
.direction.call {
  color: #5df7c2;
}
.direction.put {
  color: #ff7b7b;
}
.pnl.win {
  color: #5df7c2;
}
.pnl.loss {
  color: #ff7b7b;
}

.timer {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #ffd257;
}

.no-orders {
  color: #56607a;
  text-align: center;
  padding: 40px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  font-size: 14px;
}

.history-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.select-wrapper-small {
  position: relative;
  display: flex;
  align-items: center;
}

.select-wrapper-small select {
  appearance: none;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  padding: 6px 24px 6px 10px;
  border-radius: 8px;
  font-size: 12px;
  cursor: pointer;
}

.select-wrapper-small .arrow {
  position: absolute;
  right: 8px;
  pointer-events: none;
  color: #8fa1c4;
}

.side-pane {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.trade-card .actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 8px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 12px;
  color: #8fa1c4;
  pointer-events: none;
}

.trade-card input,
.trade-card select {
  width: 100%;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 12px 12px 12px 36px;
  border-radius: 12px;
  color: #fff;
  font-size: 14px;
  transition: all 0.2s;
}

.trade-card input:focus,
.trade-card select:focus {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  outline: none;
}

.trade-card .input-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.trade-card label {
  font-size: 12px;
  color: #8fa1c4;
  font-weight: 500;
}

.trade-card .btn-call,
.trade-card .btn-put {
  border: none;
  padding: 16px;
  border-radius: 16px;
  color: #0b1221;
  font-weight: 800;
  cursor: pointer;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  transition: transform 0.1s, filter 0.2s;
}

.trade-card .btn-call:active,
.trade-card .btn-put:active {
  transform: scale(0.98);
}

.trade-card .btn-call:hover,
.trade-card .btn-put:hover {
  filter: brightness(1.1);
}

.btn-content {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 16px;
}

.trade-card .btn-call {
  background: linear-gradient(135deg, #5df7c2, #5fa7ff);
}
.trade-card .btn-put {
  background: linear-gradient(135deg, #ff9a7a, #ff5f8f);
}

.payout {
  font-size: 12px;
  opacity: 0.8;
  font-weight: 600;
}

.trade-card .hint {
  color: #8fa1c4;
  font-size: 12px;
  text-align: center;
  margin-top: 8px;
}

.funds-body {
  gap: 16px;
}
.funds-body .input-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.funds-body .input-group {
  display: flex;
  gap: 8px;
}
.funds-body input {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  padding: 10px 12px;
  border-radius: 10px;
  flex: 1;
}
.funds-body .ghost {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  padding: 0 16px;
  border-radius: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.2s;
}
.funds-body .ghost:hover {
  background: rgba(255, 255, 255, 0.15);
}

.funds-msg {
  font-size: 12px;
  color: #5df7c2;
  padding: 8px;
  background: rgba(93, 247, 194, 0.1);
  border-radius: 8px;
  text-align: center;
}
.funds-error {
  font-size: 12px;
  color: #ff7b7b;
  padding: 8px;
  background: rgba(255, 123, 123, 0.1);
  border-radius: 8px;
  text-align: center;
}

.btn-call:disabled,
.btn-put:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.debug-card .card-body {
  max-height: 200px;
  overflow-y: auto;
}
.debug-body {
  gap: 8px;
  max-height: 200px;
  overflow-y: auto;
}
.debug-row {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.02);
}
.debug-time {
  color: #8fa1c4;
  font-variant-numeric: tabular-nums;
}
.debug-value {
  font-variant-numeric: tabular-nums;
  font-weight: 600;
  color: #fff;
}

.sim-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sim-reset {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #8fa1c4;
  padding: 6px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.sim-reset:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border-color: rgba(255, 255, 255, 0.2);
}

.sim-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.2s;
}
.sim-toggle:hover {
  background: rgba(255, 255, 255, 0.05);
}
.sim-label {
  font-size: 12px;
  color: #8fa1c4;
  font-weight: 600;
}
.sim-toggle.active .sim-label {
  color: #5df7c2;
}

/* Chart Menu Styles */
.chart-menu {
  width: 360px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.menu-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.menu-label {
  font-size: 12px;
  color: #6b7280;
  font-weight: 600;
  text-transform: uppercase;
}

.type-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.type-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  padding: 12px 4px;
  border-radius: 8px;
  color: #8fa1c4;
  cursor: pointer;
  transition: all 0.2s;
}

.type-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.type-btn.active {
  background: rgba(93, 247, 194, 0.1);
  border-color: rgba(93, 247, 194, 0.3);
  color: #5df7c2;
}

.type-btn span {
  font-size: 10px;
  text-align: center;
}

.tf-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 6px;
}

.tf-btn {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  padding: 6px 0;
  color: #8fa1c4;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
}

.tf-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.tf-btn.active {
  background: #5df7c2;
  color: #0b0e14;
  font-weight: 700;
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: #d1d4dc;
  cursor: pointer;
}

.toggle {
  width: 36px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  position: relative;
  transition: all 0.2s;
}

.toggle.active {
  background: #5df7c2;
}

.toggle-handle {
  width: 16px;
  height: 16px;
  background: #fff;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: all 0.2s;
}

.toggle.active .toggle-handle {
  transform: translateX(16px);
}

.custom-color-link {
  font-size: 12px;
  color: #8fa1c4;
  text-decoration: underline;
  text-align: right;
  cursor: pointer;
  margin-top: 4px;
  border-top: 1px dashed rgba(255, 255, 255, 0.1);
  padding-top: 8px;
}

/* Indicators Menu */
.indicators-menu {
  width: 600px;
  max-height: 400px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.indicators-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 4px;
}

.indicator-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  color: #d1d4dc;
  font-size: 13px;
  transition: all 0.2s;
}

.indicator-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.star-icon {
  color: #6b7280;
  font-size: 14px;
}

.indicator-item:hover .star-icon {
  color: #ffd257;
}

.active-dot {
  width: 6px;
  height: 6px;
  background: #5df7c2;
  border-radius: 50%;
  margin-left: auto;
}

/* Drawing Menu */
.drawing-menu {
  width: 240px;
  padding: 8px;
}

.drawing-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  color: #d1d4dc;
  font-size: 13px;
  transition: all 0.2s;
}

.drawing-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

/* Unified Floating Side Panel */
.floating-side-panel {
  position: absolute;
  top: 72px;
  right: 24px;
  width: 360px;
  max-height: calc(100% - 96px);
  background: rgba(18, 20, 28, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(20px);
  z-index: 50;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: slideInRight 0.2s ease-out;
}

@keyframes slideInRight {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

/* Content Specific Styles */
.chart-menu-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.indicators-menu-content .indicators-grid {
  display: grid;
  grid-template-columns: 1fr; /* Single column for better readability in side panel */
  gap: 4px;
}

.drawing-menu-content .tools-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.tool-grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 16px 8px;
  border-radius: 12px;
  color: #d1d4dc;
  cursor: pointer;
  transition: all 0.2s;
  min-height: 80px;
}

.tool-grid-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.15);
  color: #fff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.tool-grid-item:active {
  transform: translateY(0);
}

.tool-label {
  font-size: 11px;
  text-align: center;
  line-height: 1.3;
  font-weight: 500;
}

/* Immersive Toolbar Overrides */
.tool-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: transparent;
  border: 1px solid transparent;
  color: #8fa1c4;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  height: 32px;
}

.tool-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.tool-btn.active {
  background: rgba(93, 247, 194, 0.1);
  color: #5df7c2;
  border-color: rgba(93, 247, 194, 0.3);
  box-shadow: 0 0 10px rgba(93, 247, 194, 0.1);
}
</style>
