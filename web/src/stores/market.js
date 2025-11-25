import { defineStore } from 'pinia';
import api from '../api/axios';

export const useMarketStore = defineStore('market', {
  state: () => ({
    // 扩展支持所有交易对
    allowedSymbols: [
      'EURUSD', 'GBPUSD', 'USDJPY', 'AUDUSD', 'USDCAD', 'NZDUSD', 'EURGBP', 'EURJPY', 'GBPJPY', 'AUDJPY',
      'BTCUSDT', 'ETHUSDT', 'BNBUSDT', 'XRPUSDT', 'ADAUSDT', 'SOLUSDT', 'DOGEUSDT', 'DOTUSDT',
      'XAUUSD', 'XAGUSD', 'WTIUSD', 'NATGAS',
      'US500', 'US100', 'US30', 'DE40', 'UK100', 'JP225',
      'AAPL', 'MSFT', 'GOOGL', 'AMZN', 'TSLA', 'META', 'NVDA', 'NFLX'
    ],
    // 为每个交易对存储其初始价格（用于模拟）
    symbolPrices: {
      'EURUSD': 1.0850, 'GBPUSD': 1.2680, 'USDJPY': 149.50, 'AUDUSD': 0.6750, 'USDCAD': 1.3650, 
      'NZDUSD': 0.6150, 'EURGBP': 0.8560, 'EURJPY': 161.50, 'GBPJPY': 187.20, 'AUDJPY': 100.80,
      'BTCUSDT': 62450, 'ETHUSDT': 2850, 'BNBUSDT': 610, 'XRPUSDT': 0.52, 'ADAUSDT': 0.98, 
      'SOLUSDT': 138, 'DOGEUSDT': 0.32, 'DOTUSDT': 8.50,
      'XAUUSD': 2285, 'XAGUSD': 27.50, 'WTIUSD': 78.50, 'NATGAS': 2.95,
      'US500': 5810, 'US100': 18950, 'US30': 38620, 'DE40': 17850, 'UK100': 8120, 'JP225': 33200,
      'AAPL': 189.50, 'MSFT': 415, 'GOOGL': 185, 'AMZN': 195, 'TSLA': 245, 'META': 480, 'NVDA': 875, 'NFLX': 680
    },
    selectedSymbol: 'EURUSD',
    currentPrice: 0,
    priceHistory: [], // Array of { time, value } for chart
    currentTickTime: null,
    isConnected: false,
    socket: null,
    mockInterval: null, // 本地模拟价格定时器
    activeOrders: [],
    orderHistory: [],
    historyHasMore: true,
    balance: 0,
    balanceCurrency: 'USDT',
    candles: [],
    timeframeSec: 5,
  }),
  actions: {
    // 为指定的交易对生成模拟价格
    generateMockPrice(symbol) {
      const basePrice = this.symbolPrices[symbol] || 100;
      // 每次随机波动 0.5% - 1% 范围内
      const volatility = (Math.random() - 0.5) * 0.02;
      return basePrice * (1 + volatility);
    },
    
    connect() {
      if (this.socket) return;

      // Use current host to avoid "localhost" issues on LAN/mobile testing
      const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
      const wsHost = import.meta.env.VITE_WS_HOST || window.location.hostname;
      const wsPort = import.meta.env.VITE_WS_PORT || '8080';
      const wsUrl = `${wsProtocol}://${wsHost}:${wsPort}/ws`;

      this.socket = new WebSocket(wsUrl);

      this.socket.onopen = () => {
        console.log('WebSocket connected', wsUrl);
        this.isConnected = true;
        // Clear old ticks so we don't show stale prices after a reconnect/restart.
        this.priceHistory = [];
        this.currentPrice = 0;
      };

      this.socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          // data format: { symbol: "EURUSD", price: 1.0503, timestamp: 1716288000000 }

          if (data.symbol !== this.selectedSymbol) {
            // Ignore other symbols for now; could be extended to per-symbol cache
            return;
          }

          this.currentPrice = data.price;

          // Convert timestamp to seconds for lightweight-charts
          // Ensure integer seconds to avoid sub-second issues with standard scale
          const time = Math.floor(data.timestamp / 1000);

          // expose latest tick time for consumers (seconds)
          this.currentTickTime = time;

          // Avoid duplicates if we get multiple updates in same second
          if (this.priceHistory.length > 0) {
            const lastTime = this.priceHistory[this.priceHistory.length - 1].time;
            if (time <= lastTime) {
                // Update the last candle instead of pushing new one?
                // For LineSeries, we usually just want the latest value for that time.
                // If time is same, we update the last point.
                this.priceHistory[this.priceHistory.length - 1] = { time, value: data.price };
            } else {
                this.priceHistory.push({ time, value: data.price });
            }
          } else {
             this.priceHistory.push({ time, value: data.price });
          }

          // Keep history size manageable
          if (this.priceHistory.length > 1000) {
            this.priceHistory.shift();
          }
        } catch (e) {
          console.error('Error parsing WS message:', e);
        }
      };

      this.socket.onerror = (err) => {
        console.error('WebSocket error', err);
      };

      this.socket.onclose = () => {
        console.log('WebSocket disconnected', wsUrl);
        this.isConnected = false;
        this.socket = null;
        
        // Auto-reconnect after 3 seconds
        setTimeout(() => {
          console.log('Attempting to reconnect...');
          this.connect();
        }, 3000);
        // Start local mock as a fallback when WS is disconnected
        if (!this.mockInterval) {
          // small delay before starting mock to avoid racing with reconnect attempts
          setTimeout(() => {
            if (!this.socket && !this.mockInterval) this.startMockPriceGenerator();
          }, 1500);
        }
      };
      
      // 本地模拟器不在连接成功时立即启动，避免与服务器实时数据冲突。
      // 仅在连接关闭（或无法连接）时作为回退启动。
    },
    
    // 启动本地模拟价格生成器（每500ms生成一个新价格）
    startMockPriceGenerator() {
      if (this.mockInterval) {
        clearInterval(this.mockInterval);
      }
      
      this.mockInterval = setInterval(() => {
        const mockPrice = this.generateMockPrice(this.selectedSymbol);
        this.currentPrice = mockPrice;
        
        const time = Math.floor(Date.now() / 1000);
        
        if (this.priceHistory.length > 0) {
          const lastTime = this.priceHistory[this.priceHistory.length - 1].time;
          if (time <= lastTime) {
            this.priceHistory[this.priceHistory.length - 1] = { time, value: mockPrice };
          } else {
            this.priceHistory.push({ time, value: mockPrice });
          }
        } else {
          this.priceHistory.push({ time, value: mockPrice });
        }
        
        if (this.priceHistory.length > 1000) {
          this.priceHistory.shift();
        }
      }, 500); // 每500ms更新一次价格
    },
    
    disconnect() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
      if (this.mockInterval) {
        clearInterval(this.mockInterval);
        this.mockInterval = null;
      }
    },
    setSymbol(symbol) {
      if (this.allowedSymbols.includes(symbol)) {
        this.selectedSymbol = symbol;
        this.currentPrice = 0;
        this.priceHistory = [];
        this.orderHistory = [];
        this.historyHasMore = true;
        this.candles = [];
      }
    },
    setTimeframe(sec) {
      if (sec > 0) this.timeframeSec = sec;
    },
    async placeOrder(symbol, direction, amount, duration) {
      try {
        const orderData = {
          symbol,
          direction: direction.toUpperCase(), // 确保方向是大写（后端要求 CALL 或 PUT）
          amount: parseFloat(amount),
          duration: parseInt(duration),
          client_price: this.currentPrice || undefined,
        };
        console.log('Sending order:', orderData);
        const response = await api.post('/trade/order', orderData);
        return response.data;
      } catch (error) {
        console.error('Order failed:', error);
        console.error('Response:', error.response?.data);
        throw error;
      }
    },
    async fetchActiveOrders() {
      try {
        const response = await api.get('/trade/orders/active');
        this.activeOrders = response.data.orders;
      } catch (error) {
        console.error('Failed to fetch orders:', error);
      }
    },
    async fetchOrderHistory({ status = '', limit = 20, offset = 0, append = false } = {}) {
      try {
        const response = await api.get('/trade/orders/history', { params: { status, limit, offset } });
        const data = response.data.orders || [];
        this.historyHasMore = data.length === limit;
        if (append && offset > 0) {
          this.orderHistory = [...this.orderHistory, ...data];
        } else {
          this.orderHistory = data;
        }
      } catch (error) {
        console.error('Failed to fetch order history:', error);
      }
    },
    async fetchBalance(currency = 'USDT') {
      try {
        const response = await api.get('/wallet/balance', { params: { currency } });
        this.balance = response.data.balance;
        this.balanceCurrency = response.data.currency;
      } catch (error) {
        console.error('Failed to fetch balance:', error);
      }
    },
    async fetchCandles({ symbol, interval = 5, limit = 200 } = {}) {
      try {
        const response = await api.get('/market/candles', {
          params: {
            symbol: symbol || this.selectedSymbol,
            interval,
            limit,
          },
        });
        this.candles = response.data.candles || [];
      } catch (error) {
        console.error('Failed to fetch candles:', error);
      }
    }
  },
});
