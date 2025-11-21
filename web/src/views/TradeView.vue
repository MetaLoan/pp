<template>
  <div class="trade-container">
    <div class="chart-section" ref="chartContainer"></div>

    <div class="controls-section">
      <div class="price-display">
        <div class="symbol-row">
          <label for="symbol">Symbol</label>
          <select id="symbol" v-model="selectedSymbol">
            <option v-for="sym in marketStore.allowedSymbols" :key="sym" :value="sym">
              {{ symbolLabel(sym) }}
            </option>
          </select>
        </div>
        <h3>{{ symbolLabel(selectedSymbol) }}</h3>
        <div class="current-price" :class="{ up: isPriceUp, down: !isPriceUp }">
          {{ formattedPrice }}
        </div>
      </div>

      <div class="trade-form">
        <div class="balance">
          Balance: {{ balanceDisplay }}
        </div>
        <div v-if="errorMsg" class="alert">
          {{ errorMsg }}
        </div>
        <div class="form-group">
          <label>Amount ($)</label>
          <input type="number" v-model="amount" min="1" />
        </div>
        <div class="form-group">
          <label>Duration (s)</label>
          <select v-model="duration">
            <option value="30">30s</option>
            <option value="60">60s</option>
            <option value="300">5m</option>
          </select>
        </div>

        <div class="buttons">
          <button class="btn-call" @click="handleTrade('CALL')">
            CALL (High)
            <span class="payout">85%</span>
          </button>
          <button class="btn-put" @click="handleTrade('PUT')">
            PUT (Low)
            <span class="payout">85%</span>
          </button>
        </div>
      </div>

      <div class="positions-card">
        <div class="tab-header">
          <button
            :class="['tab', currentTab === 'active' ? 'active' : '']"
            @click="currentTab = 'active'"
          >
            Active ({{ activeOrders.length }})
          </button>
          <button
            :class="['tab', currentTab === 'recent' ? 'active' : '']"
            @click="currentTab = 'recent'"
          >
            Recent ({{ orderHistory.length }})
          </button>
        </div>

        <div v-if="currentTab === 'active'" class="positions-list">
          <div v-if="activeOrders.length === 0" class="no-orders">No active trades</div>
          <div v-else class="order-item" v-for="order in activeOrders" :key="order.id">
            <div class="order-header">
              <span class="symbol">{{ order.asset_symbol }}</span>
              <span :class="['direction', order.direction.toLowerCase()]">{{ order.direction }}</span>
            </div>
            <div class="order-details">
              <div>Entry: {{ formatPrice(order.open_price, order.asset_symbol) }}</div>
              <div>Amount: ${{ order.amount }}</div>
              <div class="timer">Time Left: {{ getTimeLeft(order) }}s</div>
              <div class="pnl" :class="getPnlClass(order)">
                Est. PnL: {{ getEstimatedPnl(order) }}
              </div>
            </div>
          </div>
        </div>

        <div v-else class="positions-list recent">
          <div v-if="orderHistory.length === 0" class="no-orders">No history yet</div>
          <div v-else class="order-item" v-for="order in orderHistory" :key="order.id">
            <div class="order-header">
              <span class="symbol">{{ order.asset_symbol }}</span>
              <span class="status" :class="statusClass(order.status)">{{ order.status }}</span>
            </div>
            <div class="order-details">
              <div>Direction: <strong :class="['direction', order.direction.toLowerCase()]">{{ order.direction }}</strong></div>
              <div>Entry: {{ formatPrice(order.open_price, order.asset_symbol) }}</div>
              <div>Exit: {{ order.close_price ? formatPrice(order.close_price, order.asset_symbol) : '--' }}</div>
              <div>Amount: ${{ order.amount }}</div>
              <div>PnL: <span :class="statusClass(order.status)">{{ formatPnl(order) }}</span></div>
              <div>Closed: {{ order.close_time ? new Date(order.close_time).toLocaleTimeString() : '--' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';
import { createChart, LineSeries } from 'lightweight-charts';
import { useMarketStore } from '../stores/market';

const chartContainer = ref(null);
const marketStore = useMarketStore();

const amount = ref(10);
const duration = ref(30);
const isPriceUp = ref(true);
const errorMsg = ref('');
const currentTab = ref('active');
let chart;
let lineSeries;
let resizeObserver;
let timerInterval;
let historyInterval;

const currentPrice = computed(() => marketStore.currentPrice);
const activeOrders = computed(() => marketStore.activeOrders);
const orderHistory = computed(() => marketStore.orderHistory);
const balanceDisplay = computed(() => `${marketStore.balanceCurrency} ${marketStore.balance.toFixed(2)}`);
const selectedSymbol = computed({
  get: () => marketStore.selectedSymbol,
  set: (val) => marketStore.setSymbol(val),
});

const symbolLabel = (sym) => {
  const labels = {
    EURUSD: 'EUR/USD',
    BTCUSDT: 'BTC/USDT',
    ETHUSDT: 'ETH/USDT',
  };
  return labels[sym] || sym;
};

const pricePrecision = computed(() => {
  const precisionMap = {
    EURUSD: 5,
    BTCUSDT: 1,
    ETHUSDT: 2,
  };
  return precisionMap[selectedSymbol.value] || 4;
});

const formattedPrice = computed(() => {
  if (!currentPrice.value) return '--';
  return currentPrice.value.toFixed(pricePrecision.value);
});

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
  const precisionMap = {
    EURUSD: 5,
    BTCUSDT: 1,
    ETHUSDT: 2,
  };
  const prec = precisionMap[symbol] || 4;
  return price ? Number(price).toFixed(prec) : '--';
};

const formatPnl = (order) => {
  if (order.status === 'won') return `+$${order.pnl.toFixed(2)}`;
  if (order.status === 'lost') return `-$${Math.abs(order.pnl).toFixed(2)}`;
  if (order.status === 'draw') return '$0.00';
  return '--';
};

const applyHistoryToSeries = (history) => {
  if (!lineSeries || history.length === 0) return;
  // Ensure time is strictly increasing to avoid time scale assertion
  const deduped = [];
  let lastTime = -Infinity;
  for (const point of history) {
    if (typeof point.time !== 'number' || Number.isNaN(point.time)) continue;
    if (point.time > lastTime) {
      deduped.push(point);
      lastTime = point.time;
    } else if (point.time === lastTime) {
      deduped[deduped.length - 1] = point; // replace same-time point
    }
  }
  if (deduped.length === 0) return;
  try {
    lineSeries.setData(deduped);
  } catch (e) {
    console.warn('Failed to setData on series', e);
  }
  if (deduped.length > 1) {
    const latest = deduped[deduped.length - 1];
    isPriceUp.value = latest.value >= deduped[deduped.length - 2].value;
  }
};

watch(
  () => marketStore.priceHistory,
  (history) => applyHistoryToSeries(history),
  { deep: true }
);

watch(
  () => selectedSymbol.value,
  () => {
    if (lineSeries) {
      lineSeries.setData([]);
      lineSeries.applyOptions({
        priceFormat: { type: 'price', precision: pricePrecision.value, minMove: 1 / Math.pow(10, pricePrecision.value) },
      });
    }
    marketStore.fetchActiveOrders();
    marketStore.fetchOrderHistory({ limit: 20 });
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
    timeScale: { timeVisible: true, secondsVisible: true },
  });

  try {
    // v5 API uses addSeries with LineSeries definition
    lineSeries = chart.addSeries(LineSeries, {
      color: '#2962FF',
      lineWidth: 2,
      priceFormat: { type: 'price', precision: pricePrecision.value, minMove: 1 / Math.pow(10, pricePrecision.value) },
    });
  } catch (err) {
    console.error('Failed to create line series', err);
    return;
  }

  // If store already has points (after hot reload), render them
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
  marketStore.fetchOrderHistory({ limit: 20 });

  timerInterval = setInterval(() => {
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
  }, 1000);

  historyInterval = setInterval(() => {
    marketStore.fetchOrderHistory({ limit: 20 });
  }, 5000);

  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  marketStore.disconnect();
  if (chart) chart.remove();
  if (resizeObserver && chartContainer.value) resizeObserver.unobserve(chartContainer.value);
  clearInterval(timerInterval);
  clearInterval(historyInterval);
  window.removeEventListener('resize', handleResize);
});

const handleResize = () => {
  if (chart && chartContainer.value) {
    chart.applyOptions({
      width: chartContainer.value.clientWidth,
      height: chartContainer.value.clientHeight,
    });
  }
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
</script>

<style scoped>
.trade-container {
  display: flex;
  height: 100vh;
  background-color: #131722;
  color: white;
}

.chart-section {
  flex: 1;
  border-right: 1px solid #2b2b43;
  min-width: 0;
}

.controls-section {
  width: 300px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
}

.positions-card {
  border-top: 1px solid #2b2b43;
  padding-top: 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tab-header {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.tab {
  padding: 10px;
  background: #1e222d;
  border: 1px solid #2b2b43;
  color: #b2b5be;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
  font-weight: 600;
}

.tab.active {
  background: #2962ff;
  color: #fff;
  border-color: #2962ff;
}

.price-display {
  text-align: center;
  margin-bottom: 20px;
}

.current-price {
  font-size: 2em;
  font-weight: bold;
}

.current-price.up {
  color: #00b894;
}

.current-price.down {
  color: #d63031;
}

.trade-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.balance {
  font-size: 0.95em;
  color: #b2b5be;
  padding: 10px;
  background: #1e222d;
  border: 1px solid #2b2b43;
  border-radius: 4px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

input,
select {
  padding: 10px;
  background: #2b2b43;
  border: 1px solid #43435c;
  color: white;
  border-radius: 4px;
}

.buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
}

button {
  padding: 15px;
  border: none;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  font-size: 1.1em;
  transition: opacity 0.2s;
}

button:hover {
  opacity: 0.9;
}

.btn-call {
  background-color: #00b894;
  color: white;
}

.btn-put {
  background-color: #d63031;
  color: white;
}

.positions-list {
  margin-top: 20px;
  border-top: 1px solid #2b2b43;
  padding-top: 20px;
  flex: 1;
  min-height: 220px;
  overflow-y: auto;
}

.positions-list.recent {
  border-top: none;
  padding-top: 0;
  max-height: 240px;
  overflow-y: auto;
}

.order-item {
  background: #1e222d;
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: 1px solid #2b2b43;
}

.order-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  font-weight: bold;
}

.direction.call {
  color: #00b894;
}

.direction.put {
  color: #d63031;
}

.order-details {
  font-size: 0.9em;
  color: #b2b5be;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 5px;
}

.pnl.win {
  color: #00b894;
}

.pnl.loss {
  color: #d63031;
}

.status.win {
  color: #00b894;
}

.status.loss {
  color: #d63031;
}

.status.draw {
  color: #e5c07b;
}

.status.pending {
  color: #61dafb;
}

.symbol-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 10px;
}

.symbol-row label {
  font-size: 0.9em;
  color: #b2b5be;
}

.symbol-row select {
  flex: 1;
  background: #2b2b43;
  border: 1px solid #43435c;
  color: #fff;
  padding: 8px;
  border-radius: 4px;
}

.alert {
  padding: 10px;
  background: #2b1f1f;
  border: 1px solid #d63031;
  color: #ffb3b3;
  border-radius: 4px;
  font-size: 0.9em;
}
</style>
