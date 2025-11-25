import { BollingerBands, EMA, RSI, SMA } from 'technicalindicators';

const mapSeries = (values, candles, offset = 0) => {
  return values
    .map((value, index) => {
      const candle = candles[index + offset];
      const numeric = Number(value);
      if (!candle || !Number.isFinite(numeric)) {
        return null;
      }
      return { time: candle.time, value: numeric };
    })
    .filter(Boolean);
};

const computeSMA = (candles, settings) => {
  const values = candles.map((candle) => candle.close);
  const result = SMA.calculate({ period: settings.period, values });
  return {
    sma: mapSeries(result, candles, settings.period - 1),
  };
};

const computeEMA = (candles, settings) => {
  const values = candles.map((candle) => candle.close);
  const result = EMA.calculate({ period: settings.period, values });
  return {
    ema: mapSeries(result, candles, settings.period - 1),
  };
};

const computeRSI = (candles, settings) => {
  const values = candles.map((candle) => candle.close);
  const result = RSI.calculate({ period: settings.period, values });
  return {
    rsi: mapSeries(result, candles, settings.period),
  };
};

const computeBollingerBands = (candles, settings) => {
  const values = candles.map((candle) => candle.close);
  const result = BollingerBands.calculate({
    period: settings.period,
    stdDev: settings.stdDev,
    values,
  });

  const upper = [];
  const middle = [];
  const lower = [];

  result.forEach((entry, index) => {
    const candle = candles[index + settings.period - 1];
    const upperVal = Number(entry.upper);
    const middleVal = Number(entry.middle);
    const lowerVal = Number(entry.lower);
    if (!candle || !Number.isFinite(upperVal) || !Number.isFinite(middleVal) || !Number.isFinite(lowerVal)) {
      return;
    }
    upper.push({ time: candle.time, value: upperVal });
    middle.push({ time: candle.time, value: middleVal });
    lower.push({ time: candle.time, value: lowerVal });
  });

  return { upper, basis: middle, lower };
};

export const indicatorsCatalog = [
  {
    id: 'sma',
    label: 'Simple Moving Average',
    description: '平滑价格、识别趋势方向的经典指标',
    supported: true,
    settingsSchema: [
      { key: 'period', label: '周期', type: 'number', min: 2, max: 200, step: 1, default: 14 },
    ],
    seriesConfigs: [
      { id: 'sma', color: '#fef08a', lineWidth: 2, priceScaleId: 'right' },
    ],
    compute: computeSMA,
  },
  {
    id: 'ema',
    label: 'Exponential Moving Average',
    description: '对最新价格更敏感，适合捕捉加速行情',
    supported: true,
    settingsSchema: [
      { key: 'period', label: '周期', type: 'number', min: 2, max: 200, step: 1, default: 21 },
    ],
    seriesConfigs: [
      { id: 'ema', color: '#fb923c', lineWidth: 2, priceScaleId: 'right' },
    ],
    compute: computeEMA,
  },
  {
    id: 'bollinger',
    label: 'Bollinger Bands',
    description: '波动率通道，用于判断趋势扩张与收缩',
    supported: true,
    settingsSchema: [
      { key: 'period', label: '周期', type: 'number', min: 5, max: 200, step: 1, default: 20 },
      { key: 'stdDev', label: '标准差', type: 'number', min: 1, max: 5, step: 0.5, default: 2 },
    ],
    seriesConfigs: [
      { id: 'upper', color: '#38bdf8', lineWidth: 1, priceScaleId: 'right' },
      { id: 'basis', color: '#fde047', lineWidth: 1, priceScaleId: 'right' },
      { id: 'lower', color: '#f472b6', lineWidth: 1, priceScaleId: 'right' },
    ],
    compute: computeBollingerBands,
  },
  {
    id: 'rsi',
    label: 'RSI',
    description: '衡量超买/超卖强弱的经典振荡指标',
    supported: true,
    settingsSchema: [
      { key: 'period', label: '周期', type: 'number', min: 2, max: 50, step: 1, default: 14 },
    ],
    seriesConfigs: [
      {
        id: 'rsi',
        color: '#22d3ee',
        lineWidth: 2,
        priceScaleId: 'indicator-rsi',
        priceFormat: { type: 'price', precision: 0, minMove: 1 },
        scaleMargins: { top: 0.8, bottom: 0 },
      },
    ],
    compute: computeRSI,
  },
  {
    id: 'macd',
    label: 'MACD',
    description: '即将上线，敬请期待',
    supported: false,
  },
  {
    id: 'stochastic',
    label: 'Stochastic',
    description: '即将上线，敬请期待',
    supported: false,
  },
];

export const getIndicatorDefaults = (indicator) => {
  if (!indicator?.settingsSchema) return {};
  return indicator.settingsSchema.reduce((acc, field) => {
    acc[field.key] = field.default;
    return acc;
  }, {});
};

export const sanitizeCandles = (candles = []) => {
  return candles
    .map((candle) => ({
      time: Number(candle.time),
      open: Number(candle.open),
      high: Number(candle.high),
      low: Number(candle.low),
      close: Number(candle.close),
    }))
    .filter(
      (candle) =>
        Number.isFinite(candle.time) &&
        Number.isFinite(candle.open) &&
        Number.isFinite(candle.high) &&
        Number.isFinite(candle.low) &&
        Number.isFinite(candle.close)
    );
};
// Technical Indicators Calculation Utilities

/**
 * Calculate Simple Moving Average (SMA)
 * @param {Array} data - Array of price objects { time, value } or numbers
 * @param {number} period - The period for SMA
 * @returns {Array} - Array of { time, value }
 */
export function calculateSMA(data, period) {
  const result = [];
  for (let i = 0; i < data.length; i++) {
    if (i < period - 1) continue;
    
    let sum = 0;
    for (let j = 0; j < period; j++) {
      sum += data[i - j].close || data[i - j].value || data[i - j];
    }
    
    result.push({
      time: data[i].time,
      value: sum / period
    });
  }
  return result;
}

/**
 * Calculate Exponential Moving Average (EMA)
 * @param {Array} data - Array of price objects { time, value }
 * @param {number} period - The period for EMA
 * @returns {Array} - Array of { time, value }
 */
export function calculateEMA(data, period) {
  const result = [];
  const k = 2 / (period + 1);
  
  let ema = data[0].close || data[0].value || data[0];
  
  for (let i = 0; i < data.length; i++) {
    const price = data[i].close || data[i].value || data[i];
    if (i === 0) {
      ema = price; // Start with first price (or SMA of first N periods ideally)
    } else {
      ema = price * k + ema * (1 - k);
    }
    
    if (i >= period - 1) {
      result.push({
        time: data[i].time,
        value: ema
      });
    }
  }
  return result;
}

/**
 * Calculate Bollinger Bands
 * @param {Array} data - Array of price objects { time, value }
 * @param {number} period - SMA period (default 20)
 * @param {number} stdDevMultiplier - Standard Deviation Multiplier (default 2)
 * @returns {Object} - { upper: [], middle: [], lower: [] }
 */
export function calculateBollingerBands(data, period = 20, stdDevMultiplier = 2) {
  const upper = [];
  const middle = [];
  const lower = [];
  
  for (let i = 0; i < data.length; i++) {
    if (i < period - 1) continue;
    
    // Calculate SMA (Middle Band)
    let sum = 0;
    for (let j = 0; j < period; j++) {
      sum += data[i - j].close || data[i - j].value || data[i - j];
    }
    const sma = sum / period;
    
    // Calculate Standard Deviation
    let sumSqDiff = 0;
    for (let j = 0; j < period; j++) {
      const val = data[i - j].close || data[i - j].value || data[i - j];
      sumSqDiff += Math.pow(val - sma, 2);
    }
    const stdDev = Math.sqrt(sumSqDiff / period);
    
    const time = data[i].time;
    
    middle.push({ time, value: sma });
    upper.push({ time, value: sma + (stdDev * stdDevMultiplier) });
    lower.push({ time, value: sma - (stdDev * stdDevMultiplier) });
  }
  
  return { upper, middle, lower };
}

/**
 * Calculate Average True Range (ATR)
 * @param {Array} data - Array of candles { time, high, low, close }
 * @param {number} period - Default 14
 * @returns {Array} - Array of { time, value }
 */
export function calculateATR(data, period = 14) {
  if (data.length < 2) return [];
  
  const trs = []; // True Ranges
  // First TR is High - Low
  trs.push(data[0].high - data[0].low);
  
  for (let i = 1; i < data.length; i++) {
    const high = data[i].high;
    const low = data[i].low;
    const prevClose = data[i-1].close;
    
    const tr = Math.max(
      high - low,
      Math.abs(high - prevClose),
      Math.abs(low - prevClose)
    );
    trs.push(tr);
  }
  
  const result = [];
  let atr = 0;
  
  // First ATR is simple average of first 'period' TRs
  for (let i = 0; i < period; i++) {
    atr += trs[i];
  }
  atr /= period;
  
  result.push({ time: data[period-1].time, value: atr });
  
  // Subsequent ATRs are smoothed
  for (let i = period; i < data.length; i++) {
    atr = ((atr * (period - 1)) + trs[i]) / period;
    result.push({ time: data[i].time, value: atr });
  }
  
  return result;
}

/**
 * Calculate RSI (Relative Strength Index)
 * @param {Array} data - Array of prices
 * @param {number} period - Default 14
 * @returns {Array} 
 */
export function calculateRSI(data, period = 14) {
  const result = [];
  if (data.length <= period) return result;

  let gains = 0;
  let losses = 0;

  // Calculate initial average gain/loss
  for (let i = 1; i <= period; i++) {
    const change = (data[i].close || data[i].value) - (data[i - 1].close || data[i - 1].value);
    if (change > 0) gains += change;
    else losses -= change;
  }

  let avgGain = gains / period;
  let avgLoss = losses / period;

  for (let i = period + 1; i < data.length; i++) {
    const change = (data[i].close || data[i].value) - (data[i - 1].close || data[i - 1].value);
    let gain = change > 0 ? change : 0;
    let loss = change < 0 ? -change : 0;

    avgGain = (avgGain * (period - 1) + gain) / period;
    avgLoss = (avgLoss * (period - 1) + loss) / period;

    let rs = avgLoss === 0 ? 100 : avgGain / avgLoss;
    let rsi = 100 - (100 / (1 + rs));

    result.push({
      time: data[i].time,
      value: rsi
    });
  }

  return result;
}


