<template>
  <div class="trade-container">
    <div class="chart-section" ref="chartContainer"></div>

    <div class="controls-section">
      <div class="price-display">
        <h3>EUR/USD</h3>
        <div class="current-price" :class="{ up: isPriceUp, down: !isPriceUp }">
          {{ currentPrice ? currentPrice.toFixed(5) : '--' }}
        </div>
      </div>

      <div class="trade-form">
        <div class="balance">
          Balance: {{ balanceDisplay }}
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

      <div class="positions-list">
        <h3>Active Positions</h3>
        <div v-if="activeOrders.length === 0" class="no-orders">No active trades</div>
        <div v-else class="order-item" v-for="order in activeOrders" :key="order.id">
          <div class="order-header">
            <span class="symbol">{{ order.asset_symbol }}</span>
            <span :class="['direction', order.direction.toLowerCase()]">{{ order.direction }}</span>
          </div>
          <div class="order-details">
            <div>Entry: {{ order.open_price.toFixed(5) }}</div>
            <div>Amount: ${{ order.amount }}</div>
            <div class="timer">Time Left: {{ getTimeLeft(order) }}s</div>
            <div class="pnl" :class="getPnlClass(order)">
              Est. PnL: {{ getEstimatedPnl(order) }}
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
let chart;
let lineSeries;
let resizeObserver;
let timerInterval;

const currentPrice = computed(() => marketStore.currentPrice);
const activeOrders = computed(() => marketStore.activeOrders);
const balanceDisplay = computed(() => `${marketStore.balanceCurrency} ${marketStore.balance.toFixed(2)}`);

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
      priceFormat: { type: 'price', precision: 5, minMove: 0.00001 },
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

  timerInterval = setInterval(() => {
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
  }, 1000);

  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  marketStore.disconnect();
  if (chart) chart.remove();
  if (resizeObserver && chartContainer.value) resizeObserver.unobserve(chartContainer.value);
  clearInterval(timerInterval);
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
  try {
    await marketStore.placeOrder('EURUSD', direction, amount.value, duration.value);
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
  } catch (error) {
    alert('Failed to place order: ' + (error.response?.data?.error || error.message));
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
</style>
