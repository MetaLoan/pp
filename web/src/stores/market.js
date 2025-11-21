import { defineStore } from 'pinia';
import api from '../api/axios';

export const useMarketStore = defineStore('market', {
  state: () => ({
    allowedSymbols: ['EURUSD', 'BTCUSDT', 'ETHUSDT'],
    selectedSymbol: 'EURUSD',
    currentPrice: 0,
    priceHistory: [], // Array of { time, value } for chart
    isConnected: false,
    socket: null,
    activeOrders: [],
    orderHistory: [],
    balance: 0,
    balanceCurrency: 'USDT',
  }),
  actions: {
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
        // Reconnect logic could go here
      };
    },
    disconnect() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
    },
    setSymbol(symbol) {
      if (this.allowedSymbols.includes(symbol)) {
        this.selectedSymbol = symbol;
        this.currentPrice = 0;
        this.priceHistory = [];
      }
    },
    async placeOrder(symbol, direction, amount, duration) {
      try {
        const response = await api.post('/trade/order', {
          symbol,
          direction,
          amount: parseFloat(amount),
          duration: parseInt(duration),
          client_price: this.currentPrice || undefined,
        });
        return response.data;
      } catch (error) {
        console.error('Order failed:', error);
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
    async fetchOrderHistory({ status = '', limit = 20 } = {}) {
      try {
        const response = await api.get('/trade/orders/history', { params: { status, limit } });
        this.orderHistory = response.data.orders;
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
    }
  },
});
