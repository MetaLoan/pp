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
      
      <nav class="main-nav">
        <button 
          :class="['nav-item', { active: activeNavItem === 'trade' }]"
          @click="activeNavItem = 'trade'"
        >
          <TrendingUp :size="20" />
          <span>Trade</span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'finance' }]"
          @click="activeNavItem = 'finance'"
        >
          <DollarSign :size="20" />
          <span>Finance</span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'profile', 'has-badge': true }]"
          @click="activeNavItem = 'profile'"
        >
          <User :size="20" />
          <span>Profile</span>
          <span class="nav-badge"></span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'market' }]"
          @click="activeNavItem = 'market'"
        >
          <ShoppingCart :size="20" />
          <span>Market</span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'achievements', 'has-badge': true }]"
          @click="activeNavItem = 'achievements'"
        >
          <Gem :size="20" />
          <span>Achievements</span>
          <span class="nav-badge"></span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'contests' }]"
          @click="activeNavItem = 'contests'"
        >
          <Trophy :size="20" />
          <span>Contests</span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'chat', 'has-count': true }]"
          @click="activeNavItem = 'chat'"
        >
          <MessageCircle :size="20" />
          <span>Chat</span>
          <span class="nav-count">5</span>
        </button>
        
        <button 
          :class="['nav-item', { active: activeNavItem === 'help', 'has-badge': true }]"
          @click="activeNavItem = 'help'"
        >
          <HelpCircle :size="20" />
          <span>Help</span>
          <span class="nav-badge"></span>
        </button>
      </nav>
      
      <div class="hero-right">
        <div class="balance-chip">
          <div class="balance-info">
            <span class="label"><Wallet :size="14" /> Balance</span>
            <span class="value">{{ balanceDisplay }}</span>
          </div>
          <button class="deposit-btn" @click="showDepositModal = true">
            <ArrowDownRight :size="16" />
            Deposit
          </button>
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
                  <LineChart v-if="chartType === 'line'" :size="18" />
                  <BarChart2 v-else-if="chartType === 'candle'" :size="18" />
                  <Activity v-else-if="chartType === 'area'" :size="18" />
                  <BarChart2 v-else :size="18" />
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
                      @click="setTimeframeByValue(tf.value)"
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
                <div class="drawing-actions">
                  <button class="action-btn primary" @click="startDrawing('trend')">
                    <Plus :size="16" />
                    <span>Start Drawing</span>
                  </button>
                  <button class="action-btn secondary" @click="clearDrawings">
                    <Eraser :size="16" />
                    <span>Clear All</span>
                  </button>
                </div>
                <div class="menu-divider"></div>
                <div class="menu-label">Drawing Tools</div>
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

        <!-- Trade Ticket - Floating Bottom Module -->
        <div class="trade-ticket-floating">
          <div class="ticket-content">
            <div v-if="errorMsg" class="alert">{{ errorMsg }}</div>
            
            <div class="dock-layout">
              <div class="ticket-inputs">
                <div class="input-group-compact">
                  <div class="input-label">
                    <DollarSign :size="16" class="label-icon" />
                    <span>Amount</span>
                  </div>
                  <div class="input-field">
                    <input type="number" v-model="amount" min="1" placeholder="10" />
                    <span class="input-unit">USDT</span>
                  </div>
                </div>
                
                <div class="dock-divider"></div>
                
                <div class="input-group-compact">
                  <div class="input-label">
                    <Clock :size="16" class="label-icon" />
                    <span>Duration</span>
                  </div>
                  <div class="input-field">
                    <select v-model="duration">
                      <option value="60">1m</option>
                      <option value="180">3m</option>
                      <option value="300">5m</option>
                      <option value="1800">30m</option>
                      <option value="3600">1h</option>
                      <option value="14400">4h</option>
                    </select>
                  </div>
                </div>
                
                <div class="dock-divider"></div>
                
                <div class="input-group-compact">
                  <div class="input-label">
                    <TrendingUp :size="16" class="label-icon" />
                    <span>Est. Return</span>
                  </div>
                  <div class="input-field input-field-readonly">
                    <div class="return-display">${{ (amount * (1 + payoutRate)).toFixed(2) }}</div>
                  </div>
                </div>
              </div>
              
              <div class="dock-divider"></div>
              
              <div class="ticket-actions">
                <div class="action-group-compact">
                  <div class="input-label">
                    <Zap :size="16" class="label-icon" />
                    <span>Direction</span>
                  </div>
                  <div class="action-buttons">
                    <button class="btn-call" :disabled="!canTrade" @click="handleTrade('CALL')">
                      <ArrowUpRight :size="16" />
                      <span>CALL</span>
                    </button>
                    <button class="btn-put" :disabled="!canTrade" @click="handleTrade('PUT')">
                      <ArrowDownRight :size="16" />
                      <span>PUT</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="ticket-hint">Est. return: <span class="highlight">${{ (amount * (1 + payoutRate)).toFixed(2) }}</span></div>
          </div>
        </div>
      </section>

      <aside class="side-pane">
        <div class="dock-rail">
          
          <button 
            v-for="dock in rightDockItems" 
            :key="dock.id" 
            :class="['dock-btn', { active: activeRightModule === dock.id }]"
            @click="activeRightModule = dock.id"
          >
            <!-- Order Active订单数量角标 -->
            <div v-if="dock.id === 'orders' && activeOrders.length > 0" class="dock-badge dock-badge-active">{{ activeOrders.length }}</div>
            <component :is="dock.icon" :size="18" />
            <div class="dock-text">
              <span class="dock-label">{{ dock.label }}</span>
            </div>
          </button>
        </div>

        <div class="side-panel-body">
          <!-- Positions Card with Tabs -->
          <div v-if="activeRightModule === 'orders'" class="card positions-card">
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

          <!-- Signal Center -->
          <div v-else-if="activeRightModule === 'signals'" class="card signal-card">
            <div class="card-head">
              <span><Antenna :size="14" /> 交易信号</span>
              <!-- Signal Tabs (moved to header) -->
              <div class="signal-tabs-header">
                <button :class="['signal-tab-header-btn', { active: activeSignalTab === 'signals' }]" @click="activeSignalTab = 'signals'">
                  交易信号
                </button>
                <button :class="['signal-tab-header-btn', { active: activeSignalTab === 'markets' }]" @click="activeSignalTab = 'markets'">
                  切换标的
                </button>
              </div>
            </div>
            
            <!-- Tab 1: Trading Signals (current symbol only) -->
            <div v-show="activeSignalTab === 'signals'">
              <!-- 信号控制栏 -->
              <div class="signal-controls">
                <div class="control-group">
                  <span class="control-label">排序</span>
                  <select v-model="signalSortBy" class="control-select">
                    <option value="new">最新优先</option>
                    <option value="confidence">趋势强度 ↓</option>
                    <option value="timing">时间框架</option>
                  </select>
                </div>
                <div class="control-group">
                  <span class="control-label">方向</span>
                  <div class="filter-buttons">
                    <button :class="['filter-btn', { active: signalFilterAction === 'all' }]" @click="signalFilterAction = 'all'">
                      全部
                    </button>
                    <button :class="['filter-btn', 'filter-call', { active: signalFilterAction === 'CALL' }]" @click="signalFilterAction = 'CALL'">
                      CALL
                    </button>
                    <button :class="['filter-btn', 'filter-put', { active: signalFilterAction === 'PUT' }]" @click="signalFilterAction = 'PUT'">
                      PUT
                    </button>
                  </div>
                </div>
              </div>

              <div class="signal-list">
                <!-- Signal Rows (Header Removed) -->
                <div v-for="sig in filteredSignals" :key="sig.title" :class="['signal-row', { 'signal-expired': !getSignalValidity(sig).isValid }]">
                  <!-- 交易标的 + 复制人数角标 -->
                  <div class="signal-cell symbol-cell">
                    <div class="symbol-wrapper">
                      <div class="copies-badge-small">
                        <Users :size="10" /> {{ sig.followers }}
                      </div>
                      <span class="symbol-text">{{ sig.symbol }}</span>
                    </div>
                  </div>

                  <!-- 趋势强度 -->
                  <div class="signal-cell trend-cell">
                    <span :class="['trend-arrows', sig.action === 'CALL' ? 'up' : 'down']">
                      {{ sig.action === 'CALL' ? '↑' : '↓' }}{{ sig.strength > 1 ? (sig.action === 'CALL' ? '↑' : '↓') : '' }}
                    </span>
                  </div>

                  <!-- 产生时间 -->
                  <div class="signal-cell time-cell">
                    <span class="time-text">{{ formatSignalTime(sig.createdAt) }}</span>
                  </div>

                  <!-- 跟随按钮 (包含倒计时) -->
                  <div class="signal-cell action-cell">
                    <button 
                      v-if="getSignalValidity(sig).isValid"
                      :class="['btn-follow', sig.action === 'CALL' ? 'btn-call' : 'btn-put']" 
                      @click="handleSignalTrade(sig)"
                      :title="`跟随 ${sig.action} 信号`">
                      <span class="btn-action-text">{{ sig.action }}</span>
                      <span class="btn-timer-text">
                        <Clock :size="10" />
                        {{ formatDuration(getSignalValidity(sig).remaining) }}
                      </span>
                    </button>
                    <span v-else class="btn-expired">已过期</span>
                  </div>
                </div>
                <div v-if="filteredSignals.length === 0" class="signal-empty">
                  <Antenna :size="32" />
                  <span>暂无符合条件的信号</span>
                </div>
              </div>
            </div>
            
            <!-- Tab 2: Switch Symbol by Category and Timeframe -->
            <div v-show="activeSignalTab === 'markets'" class="market-trends-view">
              <!-- Timeframe Selector (6 options) -->
              <div class="section-header">时间框架</div>
              <div class="timeframe-grid">
                <button v-for="tf in ['M1', 'M3', 'M5', 'M30', 'H1', 'H4']" :key="tf" 
                  :class="['timeframe-btn', { active: selectedTimeframe === tf }]"
                  @click="setTimeframeByLabel(tf)">
                  {{ tf }}
                </button>
              </div>
              
              <!-- Trading Pairs List (with dynamic trends) -->
              <div class="section-header">交易标的</div>
              <div class="pairs-list-container">
                <div v-for="pair in tradingPairs" :key="pair.symbol" 
                  :class="['pair-list-item', { active: selectedSymbol === pair.symbol }]"
                  @click="selectedSymbol = pair.symbol">
                  <div class="pair-main">
                    <span class="pair-symbol">{{ pair.symbol }}</span>
                    <span class="pair-display">{{ pair.display }}</span>
                  </div>
                  <div v-if="getTrendForPair(pair.symbol)" :class="['pair-trend', getTrendForPair(pair.symbol).direction === 'CALL' ? 'trend-up' : 'trend-down']">
                    <span class="trend-arrow">{{ getTrendForPair(pair.symbol).direction === 'CALL' ? '↑' : '↓' }}{{ getTrendForPair(pair.symbol).strength > 1 ? getTrendForPair(pair.symbol).direction === 'CALL' ? '↑' : '↓' : '' }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Social Trading -->
          <div v-else-if="activeRightModule === 'social'" class="card social-card">
            <div class="card-head">
              <span><Users :size="14" /> 社交交易</span>
              <span class="badge">Top Movers</span>
            </div>
            <div class="social-list">
              <div v-for="trader in socialLeaders" :key="trader.name" class="social-row">
                <div class="avatar">{{ trader.initials }}</div>
                <div class="social-meta">
                  <div class="social-name">{{ trader.name }}</div>
                  <div class="social-sub">{{ trader.region }} · 胜率 {{ trader.winRate }}%</div>
                </div>
                <div class="social-stats">
                  <span class="pill pill-green">ROI {{ trader.roi }}%</span>
                  <span class="pill pill-soft">{{ trader.copiers }} 跟单</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Quick Trade -->
          <div v-else-if="activeRightModule === 'quick'" class="card quick-card">
            <div class="card-head">
              <span><Target :size="14" /> 快速交易</span>
              <span class="badge">One-tap</span>
            </div>
            <div class="quick-grid">
              <button 
                v-for="preset in quickSetups" 
                :key="preset.label" 
                class="quick-chip"
                @click="() => { amount.value = preset.amount; duration.value = preset.duration; activeRightModule.value = 'orders'; }"
              >
                <span class="quick-label">{{ preset.label }}</span>
                <span class="quick-sub">{{ preset.duration }}s · ${{ preset.amount }}</span>
                <span class="pill" :class="preset.dir === 'CALL' ? 'pill-green' : 'pill-red'">{{ preset.dir }}</span>
              </button>
            </div>
            <div class="quick-footer">
              <span>当前回报率</span>
              <div class="roi-meter">
                <div class="roi-fill" :style="{ width: '85%' }"></div>
                <span class="roi-label">85%</span>
              </div>
            </div>
          </div>

          <!-- Pending / Trigger Orders -->
          <div v-else-if="activeRightModule === 'pending'" class="card pending-card">
            <div class="card-head">
              <span><Hourglass :size="14" /> 挂单交易</span>
              <span class="badge">计划</span>
            </div>
            <div class="pending-list">
              <div v-for="item in pendingOrders" :key="item.id" class="pending-row">
                <div class="pending-meta">
                  <div class="pending-name">{{ item.symbol }}</div>
                  <div class="pending-sub">触发价 {{ formatPrice(item.trigger, item.symbol) }}</div>
                </div>
                <div class="pending-actions">
                  <span class="pill pill-soft">{{ item.window }}</span>
                  <span class="pill" :class="item.dir === 'CALL' ? 'pill-green' : 'pill-red'">{{ item.dir }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Shortcuts -->
          <div v-else-if="activeRightModule === 'shortcuts'" class="card shortcuts-card">
            <div class="card-head">
              <span><Keyboard :size="14" /> 快捷键</span>
              <span class="badge accent">Pro</span>
            </div>
            <div class="shortcut-grid">
              <div v-for="hotkey in shortcutList" :key="hotkey.keys" class="shortcut-row">
                <span class="keycap">{{ hotkey.keys }}</span>
                <div class="shortcut-meta">
                  <div class="shortcut-name">{{ hotkey.label }}</div>
                  <div class="shortcut-sub">{{ hotkey.detail }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Signal Detail Modal -->
        <!-- Deposit Modal -->
        <div v-if="showDepositModal" class="modal-overlay" @click="showDepositModal = false">
          <div class="modal-card" @click.stop>
            <div class="modal-header">
              <h3><Wallet :size="20" /> Deposit Funds</h3>
              <button class="modal-close" @click="showDepositModal = false">
                <X :size="20" />
              </button>
            </div>
            <div class="modal-body">
              <div class="deposit-input-group">
                <label class="deposit-label">Amount</label>
                <div class="deposit-input-wrapper">
                  <span class="currency-prefix">USDT</span>
                  <input 
                    type="number" 
                    v-model="depositAmount" 
                    min="1" 
                    placeholder="0.00"
                    class="deposit-input"
                  />
                </div>
              </div>
              <div class="quick-amounts">
                <button v-for="amt in [10, 50, 100, 500]" :key="amt" class="quick-btn" @click="depositAmount = amt">
                  ${{ amt }}
                </button>
              </div>
              <div v-if="fundsMsg" class="funds-msg">{{ fundsMsg }}</div>
              <div v-if="fundsError" class="funds-error">{{ fundsError }}</div>
            </div>
            <div class="modal-footer">
              <button class="btn-secondary" @click="showDepositModal = false">Cancel</button>
              <button class="btn-primary" @click="handleDeposit">
                <ArrowDownRight :size="16" />
                Confirm Deposit
              </button>
            </div>
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
  Minus, MoreVertical, Menu, Fan, Maximize, Square, Move, GitBranch,
  ShoppingCart, Award, Trophy, MessageCircle, HelpCircle, User, Gem,
  Antenna, Users, Target, Hourglass, Keyboard
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
const activeNavItem = ref('trade');
const depositAmount = ref(100);
const withdrawAmount = ref(50);
const fundsMsg = ref('');
const fundsError = ref('');
const showDepositModal = ref(false);
const signalSortBy = ref('confidence'); // 'confidence' | 'timing' | 'new'
const signalFilterAction = ref('all'); // 'all' | 'CALL' | 'PUT'
const signalFilterTiming = ref('all'); // 'all' | '1m' | '2m' | '3m' | '4m' | '5m'
const chartType = ref('line'); // line | area | candle
const timeframe = ref(60); // seconds per bar (default M1)
const showSMA = ref(false);
const showEMA = ref(false);
const smaPeriod = ref(10);
const timeframeOptions = [1, 5, 15, 30, 60, 300, 600]; // Keep for logic mapping

// 时间框架标签到秒数的映射
const timeframeMap = {
  'M1': 60,
  'M3': 180,
  'M5': 300,
  'M30': 1800,
  'H1': 3600,
  'H4': 14400
};

// 秒数到标签的反向映射
const reverseTimeframeMap = {
  60: 'M1',
  180: 'M3',
  300: 'M5',
  1800: 'M30',
  3600: 'H1',
  14400: 'H4'
};

// 根据时间框架标签更新图表级别
const setTimeframeByLabel = (label) => {
  selectedTimeframe.value = label;
  timeframe.value = timeframeMap[label] || 60;
  // 保存到localStorage以便页面刷新时恢复
  localStorage.setItem('tradeView_timeframe', String(timeframeMap[label] || 60));
};

// 根据秒数值设置时间框架，只在预设列表内时更新selectedTimeframe
const setTimeframeByValue = (seconds) => {
  timeframe.value = seconds;
  // 保存到localStorage以便页面刷新时恢复
  localStorage.setItem('tradeView_timeframe', String(seconds));
  // 如果这个秒数值对应预设的6个时间框架之一，则更新selectedTimeframe
  if (reverseTimeframeMap[seconds]) {
    selectedTimeframe.value = reverseTimeframeMap[seconds];
  } else {
    // 如果不在预设列表内，则清空selectedTimeframe（使Markets TAB中的按钮都不激活）
    selectedTimeframe.value = null;
  }
};

// Right dock modules
const activeRightModule = ref('orders');
const rightDockItems = [
  { id: 'orders', label: 'Orders', icon: History },
  { id: 'signals', label: 'Signals', icon: Antenna },
  { id: 'social', label: 'Social', icon: Users },
  { id: 'quick', label: 'Quick Trade', icon: Target },
  { id: 'pending', label: 'Pending', icon: Hourglass },
  { id: 'shortcuts', label: 'Shortcuts', icon: Keyboard },
];

// 生成20条M30时间框架的测试信号数据
const generateM30TestSignals = () => {
  const signals = [];
  const now = Date.now();
  const actions = ['CALL', 'PUT'];
  
  // 交易对列表
  const tradingPairsList = [
    { symbol: 'EURUSD' },
    { symbol: 'GBPUSD' },
    { symbol: 'USDJPY' },
    { symbol: 'AUDUSD' },
    { symbol: 'USDCAD' },
    { symbol: 'NZDUSD' },
    { symbol: 'EURGBP' },
    { symbol: 'EURJPY' },
    { symbol: 'GBPJPY' },
    { symbol: 'AUDJPY' },
    { symbol: 'BTCUSDT' },
    { symbol: 'ETHUSDT' },
    { symbol: 'BNBUSDT' },
    { symbol: 'XRPUSDT' },
    { symbol: 'ADAUSDT' },
    { symbol: 'SOLUSDT' },
    { symbol: 'DOGEUSDT' },
    { symbol: 'DOTUSDT' },
    { symbol: 'XAUUSD' },
    { symbol: 'XAGUSD' },
    { symbol: 'WTIUSD' },
    { symbol: 'NATGAS' },
    { symbol: 'US500' },
    { symbol: 'US100' },
    { symbol: 'US30' },
    { symbol: 'DE40' },
    { symbol: 'UK100' },
    { symbol: 'JP225' },
    { symbol: 'AAPL' },
    { symbol: 'MSFT' },
    { symbol: 'GOOGL' },
    { symbol: 'AMZN' },
    { symbol: 'TSLA' },
    { symbol: 'META' },
    { symbol: 'NVDA' },
    { symbol: 'NFLX' },
  ];
  
  const signalTitles = [
    '突破', '反弹', '拐点', '加速', '回调', '强势', '弱势', '盘整', '加仓', '获利',
    '冲高', '探底', '缩量', '放量', '修复', '衰竭', '启动', '转折', '蓄势', '狂欢'
  ];
  
  const metrics = [
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
  
  // 为每个交易对和时间框架生成20条信号
  for (let i = 0; i < 20; i++) {
    const randomSymbol = tradingPairsList[Math.floor(Math.random() * tradingPairsList.length)];
    const titleIndex = i % signalTitles.length;
    const metricIndex = (i + Math.floor(Math.random() * metrics.length)) % metrics.length;
    
    const signal = {
      title: `${randomSymbol.symbol} ${signalTitles[titleIndex]}`,
      metric: metrics[metricIndex],
      confidence: 0.6 + Math.random() * 0.35, // 0.6 - 0.95
      action: actions[Math.floor(Math.random() * actions.length)],
      timing: 'm30',
      symbol: randomSymbol.symbol,
      timeframe: 'M30',
      strength: Math.random() > 0.5 ? 2 : 1,
      amount: Math.floor(Math.random() * 150) + 25, // 25 - 175
      duration: 1800 + Math.floor(Math.random() * 1800), // 30m - 60m (1800s - 3600s)
      followers: Math.floor(Math.random() * 1900) + 100,
      copied: Math.floor(Math.random() * 15),
      createdAt: now - i * 60000, // 每条间隔1分钟
      validity: 1800000 + Math.floor(Math.random() * 1800000), // 30m - 60m
      expiryTime: now - i * 60000 + (1800000 + Math.floor(Math.random() * 1800000)),
      isNew: i === 0, // 第一条标记为新
    };
    
    signals.push(signal);
  }
  
  return signals;
};

const signalFeed = ref(generateM30TestSignals());

// Signal module tabs
const activeSignalTab = ref('signals'); // 'signals' | 'markets'

// 时间框架选择
const selectedTimeframe = ref('M1');

// 根据标的和时间框架获取趋势强度和方向
const getTrendForPair = (symbol) => {
  // 找出该 [symbol + selectedTimeframe] 组合下所有信号
  const matchingSignals = signalFeed.value.filter(s => 
    s.symbol === symbol && s.timeframe === selectedTimeframe.value
  );
  
  if (matchingSignals.length === 0) {
    return null;
  }
  
  // 过滤有效信号（未过期）
  const now = Date.now();
  const validSignals = matchingSignals.filter(s => s.expiryTime > now);
  
  if (validSignals.length === 0) {
    return null;
  }
  
  // 找复制人数最多的有效信号
  const maxCopied = Math.max(...validSignals.map(s => s.copied));
  const topSignals = validSignals.filter(s => s.copied === maxCopied);
  
  // 如果并列，选第一个（最新的）
  const selectedSignal = topSignals[0];
  
  return {
    direction: selectedSignal.action, // 'CALL' or 'PUT'
    strength: selectedSignal.strength // 1 or 2
  };
};

// 信号列表过滤 - 按选中的标的过滤（不按时间框架过滤）
const filteredSignals = computed(() => {
  return signalFeed.value
    .filter(s => s.symbol === selectedSymbol.value)
    .filter(s => signalFilterAction.value === 'all' || s.action === signalFilterAction.value)
    .sort((a, b) => {
      if (signalSortBy.value === 'confidence') {
        return b.confidence - a.confidence;
      } else if (signalSortBy.value === 'followers') {
        return b.followers - a.followers;
      } else {
        return b.createdAt - a.createdAt;
      }
    });
});

const socialLeaders = ref([
  { name: 'NovaQuant', initials: 'NQ', region: '新加坡', winRate: 78, roi: 34, copiers: '2.1k' },
  { name: 'AlphaWave', initials: 'AW', region: '伦敦', winRate: 74, roi: 28, copiers: '1.4k' },
  { name: 'TokyoEdge', initials: 'TE', region: '东京', winRate: 71, roi: 22, copiers: '940' },
]);

const quickSetups = ref([
  { label: '顺势突破', duration: 60, amount: 25, dir: 'CALL' },
  { label: '回撤进场', duration: 180, amount: 50, dir: 'PUT' },
  { label: '脉冲抢先', duration: 30, amount: 15, dir: 'CALL' },
]);

const pendingOrders = ref([
  { id: 1, symbol: 'EURUSD', trigger: 1.0825, dir: 'CALL', window: '60s' },
  { id: 2, symbol: 'BTCUSDT', trigger: 62850, dir: 'PUT', window: '5m' },
]);

const shortcutList = [
  { keys: 'A', label: '快速下单 CALL', detail: '使用当前金额和时长' },
  { keys: 'S', label: '快速下单 PUT', detail: '使用当前金额和时长' },
  { keys: 'H', label: '显示/隐藏信号', detail: '切换到交易信号面板' },
  { keys: '⌘ + /', label: '查看所有快捷键', detail: '打开帮助弹窗' },
];

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

// 信号过滤和排序
// 时间框架市场数据（用于切换标的Tab）
const marketTrends = computed(() => {
  // 为所有交易对和时间框架生成趋势数据
  return tradingPairs.value.map(pair => {
    const timeframes = ['1m', '2m', '3m', '4m', '5m'];
    const trends = {};
    timeframes.forEach(tf => {
      const key = `${pair.symbol}-${tf}`;
      trends[tf] = {
        symbol: pair.symbol,
        display: pair.display,
        timeframe: tf,
        strength: getRandomTrendStrength(key),
        direction: getRandomTrendDirection(key),
      };
    });
    return { symbol: pair.symbol, display: pair.display, trends };
  });
});

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
let signalInterval;
let trendRefreshInterval; // 趋势刷新定时器
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

const currentTime = ref(Date.now()); // 用于计算相对时间

// 格式化信号生成时间（相对值）
const formatSignalTime = (createdAt) => {
  const diff = currentTime.value - createdAt;
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes}分钟前`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}小时前`;
  const days = Math.floor(hours / 24);
  return `${days}天前`;
};

// 格式化倒计时 (h:m:s)
const formatDuration = (ms) => {
  if (ms <= 0) return '00:00';
  const seconds = Math.floor((ms / 1000) % 60);
  const minutes = Math.floor((ms / (1000 * 60)) % 60);
  const hours = Math.floor((ms / (1000 * 60 * 60)) % 24);
  
  const s = seconds.toString().padStart(2, '0');
  const m = minutes.toString().padStart(2, '0');
  
  if (hours > 0) {
    const h = hours.toString().padStart(2, '0');
    return `${h}:${m}:${s}`;
  }
  return `${m}:${s}`;
};

// 计算信号有效性和倒计时
const getSignalValidity = (signal) => {
  const elapsed = currentTime.value - signal.createdAt;
  const remaining = signal.validity - elapsed;
  if (remaining <= 0) return { isValid: false, remaining: 0, percent: 0 };
  const percent = Math.max(0, (remaining / signal.validity) * 100);
  return { isValid: true, remaining, percent };
};

// 生成随机趋势强度（用于演示）
const trendCache = new Map();
const getRandomTrendStrength = (symbol = '') => {
  // 缓存趋势强度避免闪烁
  // 返回0-1，表示1-2个箭头（0.5及以上 = 2个箭头，以下 = 1个箭头）
  if (!trendCache.has(`strength-${symbol}`)) {
    trendCache.set(`strength-${symbol}`, Math.random());
  }
  return trendCache.get(`strength-${symbol}`);
};

// 生成随机趋势方向（上升1 或下降-1）
const getRandomTrendDirection = (symbol = '') => {
  if (!trendCache.has(`direction-${symbol}`)) {
    trendCache.set(`direction-${symbol}`, Math.random() > 0.5 ? 1 : -1);
  }
  return trendCache.get(`direction-${symbol}`);
};

// 获取箭头显示（1-2个）
const getTrendArrows = (strength, direction) => {
  // 0.8以上为强趋势（2个箭头），否则为普通趋势（1个箭头）
  const arrowCount = strength >= 0.8 ? 2 : 1;
  const arrow = direction > 0 ? '↑' : '↓';
  return arrow.repeat(arrowCount);
};

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
    // 更新market store中的selectedSymbol，这会清空旧的价格历史
    marketStore.setSymbol(selectedSymbol.value);
    
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
  
  // 从localStorage恢复时间框架选择
  const savedTimeframe = localStorage.getItem('tradeView_timeframe');
  if (savedTimeframe) {
    const timeframeValue = parseInt(savedTimeframe, 10);
    setTimeframeByValue(timeframeValue);
  }
  
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
    currentTime.value = Date.now(); // 更新当前时间用于倒计时
  }, 1000);

  historyInterval = setInterval(() => {
    fetchHistory(true);
  }, 5000);

  restartCandleInterval();

  // 信号推送间隔 - 模拟实时信号
  signalInterval = setInterval(() => {
    pushNewSignal();
  }, 30000); // 每30秒推送一条新信号

  // 趋势刷新间隔 - 监控信号过期和复制人数变化，实时更新趋势
  // 当最多人复制的信号过期或被新信号超越时，自动重新计算趋势
  trendRefreshInterval = setInterval(() => {
    const now = Date.now();
    
    // 检查是否有信号过期
    let hasExpiredSignals = signalFeed.value.some(s => s.expiryTime <= now);
    
    if (hasExpiredSignals) {
      // 移除过期信号
      signalFeed.value = signalFeed.value.filter(s => s.expiryTime > now);
      
      // 强制重新渲染所有对应的趋势
      // 因为计算属性会自动响应signalFeed的变化，所以无需手动刷新
    }
  }, 1000); // 每秒检查一次

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
  clearInterval(signalInterval);
  clearInterval(trendRefreshInterval);
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
    // 检查 duration 是否最少 5 秒
    if (duration.value < 5) {
      errorMsg.value = `⚠ 订单时长过短（${duration.value}s），最少需要5秒`;
      return;
    }
    
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
    setTimeout(() => {
      showDepositModal.value = false;
      fundsMsg.value = '';
    }, 2000);
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

const handleSignalTrade = async (signal) => {
  errorMsg.value = '';
  try {
    // 检查信号是否已过期
    const now = Date.now();
    if (now >= signal.expiryTime) {
      errorMsg.value = `⚠ 信号已过期，无法执行交易`;
      return;
    }
    
    // 计算剩余有效期（毫秒转秒）
    const remainingTime = signal.expiryTime - now;
    const remainingSeconds = Math.ceil(remainingTime / 1000);
    
    // 检查交割时间是否 >= 5秒（过期前5秒不允许跟随）
    if (remainingSeconds < 5) {
      errorMsg.value = `⚠ 信号即将过期（${remainingSeconds}s），不允许执行（需至少5秒有效期）`;
      return;
    }
    
    // 使用剩余时间作为订单时长，而不是signal.duration
    const orderDuration = remainingSeconds;
    
    // 更新交易参数
    amount.value = signal.amount;
    duration.value = orderDuration; // 设置为剩余时间
    selectedSymbol.value = signal.symbol;
    
    // 执行下单（使用计算出的orderDuration）
    await marketStore.placeOrder(signal.symbol, signal.action, signal.amount, orderDuration);
    
    // 递增 copied 计数（跟随 = 复制）
    signal.copied += 1;
    
    // 检查该信号是否成为此[标的+时间框架]组合的最多人复制信号
    // 如果是，趋势会自动刷新（因为getTrendForPair是响应式计算）
    const sameGroupSignals = signalFeed.value.filter(s => 
      s.symbol === signal.symbol && 
      s.timeframe === signal.timeframe && 
      s.expiryTime > Date.now() // 只考虑有效信号
    );
    if (sameGroupSignals.length > 0) {
      const maxCopied = Math.max(...sameGroupSignals.map(s => s.copied));
      if (signal.copied === maxCopied) {
        // 该信号现在是最多人复制的，趋势已自动更新
        console.log(`✓ 信号"${signal.title}"成为${signal.symbol}/${signal.timeframe}的主导信号`);
      }
    }
    
    // 刷新数据
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
    
    // 显示成功反馈
    errorMsg.value = `✓ 信号交易执行: ${signal.action} ${signal.title} (${orderDuration}s)`;
    
    // 3秒后清除提示
    setTimeout(() => {
      errorMsg.value = '';
    }, 3000);
  } catch (error) {
    const errorText = error.response?.data?.error || error.message || 'Unknown error';
    
    // 处理特定的业务错误
    if (errorText.includes('too many open orders')) {
      errorMsg.value = `⚠ 开仓订单过多，请先关闭一些订单后再试`;
    } else if (errorText.includes('insufficient balance')) {
      errorMsg.value = `⚠ 余额不足，请充值后再试`;
    } else if (errorText.includes('invalid amount')) {
      errorMsg.value = `⚠ 订单金额无效`;
    } else {
      errorMsg.value = `⚠ 下单失败: ${errorText}`;
    }
  }
};


const pushNewSignal = () => {
  // 为每个 [标的 + 时间框架] 组合生成一个随机信号
  const now = Date.now();
  const timeframes = ['M1', 'M3', 'M5', 'M30', 'H1', 'H4'];
  const actions = ['CALL', 'PUT'];
  
  // 随机选择一个标的
  const randomSymbol = tradingPairs.value[Math.floor(Math.random() * tradingPairs.value.length)];
  
  // 随机选择一个时间框架
  const randomTimeframe = timeframes[Math.floor(Math.random() * timeframes.length)];
  
  // 随机选择方向
  const randomAction = actions[Math.floor(Math.random() * actions.length)];
  
  // 随机选择强度 (1 或 2)
  const randomStrength = Math.random() > 0.5 ? 2 : 1;
  
  // 随机生成跟踪人数 (100-2000)
  const randomFollowers = Math.floor(Math.random() * 1900) + 100;
  
  // 随机设置有效期 (60s - 600s)
  const validity = (Math.floor(Math.random() * 540) + 60) * 1000;
  
  // 生成信号标题和指标
  const signalTitles = [
    `${randomSymbol.symbol} 突破`,
    `${randomSymbol.symbol} 反弹`,
    `${randomSymbol.symbol} 拐点`,
    `${randomSymbol.symbol} 加速`,
    `${randomSymbol.symbol} 回调`,
  ];
  const title = signalTitles[Math.floor(Math.random() * signalTitles.length)];
  
  const metrics = [
    'RSI 指标强势',
    'MACD 金叉',
    'MA 均线突破',
    'Stoch 信号',
    'CCI 极值',
    'ATR 突破',
    'Volume 突增',
  ];
  const metric = metrics[Math.floor(Math.random() * metrics.length)];
  
  const newSignal = {
    title,
    metric,
    confidence: 0.6 + Math.random() * 0.3, // 0.6 - 0.9
    action: randomAction,
    timing: randomTimeframe.toLowerCase(),
    symbol: randomSymbol.symbol,
    timeframe: randomTimeframe,
    strength: randomStrength,
    amount: Math.floor(Math.random() * 150) + 25, // 25 - 175
    duration: 60 * (Math.floor(Math.random() * 9) + 1), // 60s - 540s
    followers: randomFollowers,
    createdAt: now,
    validity: validity,
    expiryTime: now + validity,
    isNew: true,
  };
  
  // 在列表头部插入新信号
  signalFeed.value.unshift(newSignal);
  
  // 移除 isNew 标记以停止动画(500ms后)
  setTimeout(() => {
    const idx = signalFeed.value.findIndex(s => s.title === newSignal.title && s.createdAt === newSignal.createdAt);
    if (idx !== -1) {
      signalFeed.value[idx].isNew = false;
    }
  }, 500);
  
  // 限制信号列表不超过50条
  if (signalFeed.value.length > 50) {
    signalFeed.value.pop();
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
  padding-bottom: 140px;
  font-family: 'Inter', 'SF Pro Display', system-ui, -apple-system, sans-serif;
  display: flex;
  flex-direction: column;
  gap: 20px;
  --dock-width: clamp(68px, 7vw, 86px);
  --dock-offset: clamp(10px, 1.6vw, 18px);
}

/* 全局滚动条美化 */
/* Chrome / Safari / Edge 滚动条 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
}

::-webkit-scrollbar-thumb {
  background: rgba(93, 247, 194, 0.25);
  border-radius: 4px;
  transition: background 0.2s;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(93, 247, 194, 0.45);
}

::-webkit-scrollbar-corner {
  background: rgba(255, 255, 255, 0.02);
}

/* Firefox 滚动条 */
* {
  scrollbar-width: thin;
  scrollbar-color: rgba(93, 247, 194, 0.25) rgba(255, 255, 255, 0.02);
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
  gap: 32px;
}

/* Main Navigation */
.main-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
  justify-content: center;
}

.nav-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 16px;
  background: transparent;
  border: none;
  border-radius: 12px;
  color: #8fa1c4;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.nav-item.active {
  background: rgba(93, 247, 194, 0.1);
  color: #5df7c2;
}

.nav-item.active:hover {
  background: rgba(93, 247, 194, 0.15);
}

/* Navigation Badge (dot) */
.nav-badge {
  position: absolute;
  top: 8px;
  right: 12px;
  width: 6px;
  height: 6px;
  background: #ef4444;
  border-radius: 50%;
  border: 1.5px solid rgba(11, 14, 20, 0.9);
  animation: pulse-badge 2s ease-in-out infinite;
}

@keyframes pulse-badge {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7);
  }
  50% {
    box-shadow: 0 0 0 4px rgba(239, 68, 68, 0);
  }
}

/* Navigation Count */
.nav-count {
  position: absolute;
  top: 6px;
  right: 10px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: #fff;
  font-size: 10px;
  font-weight: 700;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1.5px solid rgba(11, 14, 20, 0.9);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.4);
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
  padding: 12px 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(11, 14, 20, 0.6);
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 16px;
  min-width: 280px;
  transition: all 0.2s;
}

.balance-chip:hover {
  border-color: rgba(255, 255, 255, 0.15);
  background: rgba(11, 14, 20, 0.8);
}

.balance-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
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

.deposit-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #5df7c2 0%, #3dffb5 100%);
  color: #0a0e14;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(93, 247, 194, 0.3);
}

.deposit-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(93, 247, 194, 0.4);
}

.deposit-btn:active {
  transform: translateY(0);
}

.workspace {
  display: grid;
  grid-template-columns: 65fr 35fr;
  gap: 8px;
  flex: 1;
  align-items: start;
}

.chart-pane {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 0;
  position: relative;
}

/* Trade Ticket Floating Module - Mac Dock Style */
.trade-ticket-floating {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  width: auto;
  max-width: calc(100vw - 48px);
  overflow: visible;
  filter: drop-shadow(0 8px 32px rgba(0, 0, 0, 0.6));
}

.ticket-content {
  position: relative;
  background: rgba(18, 20, 28, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 24px;
  padding: 16px 24px;
  box-shadow: none;
  backdrop-filter: blur(40px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: visible;
}

.ticket-content:hover {
  background: rgba(18, 20, 28, 0.9);
  transform: translateY(-2px);
}

.dock-layout {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: nowrap;
  overflow-x: auto;
  overflow-y: visible;
  scrollbar-width: none;
  position: relative;
  z-index: 1;
  filter: drop-shadow(0 0 0 transparent);
}

.dock-layout::-webkit-scrollbar {
  display: none;
}

.dock-divider {
  width: 1px;
  height: 32px;
  background: linear-gradient(to bottom, 
    rgba(255, 255, 255, 0),
    rgba(255, 255, 255, 0.15),
    rgba(255, 255, 255, 0)
  );
}

.ticket-header {
  display: none;
}

.ticket-inputs {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 0;
}

.input-group-compact {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 0 0 auto;
  width: 150px;
}

.input-group-compact .input-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #8fa1c4;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.input-group-compact .label-icon {
  color: #5df7c2;
  opacity: 0.8;
}

.input-group-compact .input-field {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1.5px solid rgba(255, 255, 255, 0.12);
  border-radius: 14px;
  padding: 2px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  width: 100%;
  height: 48px;
  box-sizing: border-box;
}

.input-group-compact .input-field-readonly {
  background: rgba(93, 247, 194, 0.08);
  border-color: rgba(93, 247, 194, 0.25);
  cursor: default;
}

.input-group-compact .input-field-readonly:hover {
  background: rgba(93, 247, 194, 0.1);
  border-color: rgba(93, 247, 194, 0.3);
}

.input-group-compact .return-display {
  flex: 1;
  padding: 10px 12px;
  color: #5df7c2;
  font-size: 15px;
  font-weight: 800;
  text-align: center;
  width: 100%;
  box-sizing: border-box;
}

.input-group-compact .input-field:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
}

.input-group-compact .input-field:focus-within {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(93, 247, 194, 0.5);
  box-shadow: 0 0 0 3px rgba(93, 247, 194, 0.1);
}

.input-group-compact input,
.input-group-compact select {
  flex: 1;
  background: transparent;
  border: none;
  padding: 10px 12px 10px 36px;
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  width: 100%;
  box-sizing: border-box;
  outline: none;
}

.input-group-compact .input-unit {
  padding: 6px 12px;
  background: rgba(93, 247, 194, 0.15);
  color: #5df7c2;
  font-size: 11px;
  font-weight: 700;
  border-radius: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-right: 4px;
}

.input-group-compact select {
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  padding-right: 32px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%235df7c2' d='M6 9L1 4h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
}

.input-group-compact select option {
  background: #12141c;
  color: #fff;
  padding: 8px;
}

.ticket-actions {
  display: flex;
  flex: 0 0 auto;
}

.action-group-compact {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 320px;
}

.action-group-compact .input-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #8fa1c4;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.action-group-compact .label-icon {
  color: #5df7c2;
  opacity: 0.8;
}

.action-buttons {
  display: flex;
  gap: 10px;
  width: 100%;
  height: 48px;
  box-sizing: border-box;
}

.ticket-actions .btn-call,
.ticket-actions .btn-put {
  flex: 1;
  border: none;
  border-radius: 12px;
  font-weight: 700;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  white-space: nowrap;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  height: 100%;
  padding: 0;
  min-width: 0;
  box-sizing: border-box;
  position: relative;
}

.ticket-actions .btn-call {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: #ffffff;
  box-shadow: inset 0 0 12px rgba(16, 185, 129, 0.3);
}

.ticket-actions .btn-call:hover:not(:disabled) {
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
  transform: translateY(-2px);
  box-shadow: inset 0 0 20px rgba(16, 185, 129, 0.5), inset 0 0 8px rgba(255, 255, 255, 0.2);
}

.ticket-actions .btn-put {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: #ffffff;
  box-shadow: inset 0 0 12px rgba(239, 68, 68, 0.3);
}

.ticket-actions .btn-put:hover:not(:disabled) {
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
  transform: translateY(-2px);
  box-shadow: inset 0 0 20px rgba(239, 68, 68, 0.5), inset 0 0 8px rgba(255, 255, 255, 0.2);
}

.ticket-actions .btn-call:disabled,
.ticket-actions .btn-put:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.ticket-hint {
  display: none;
}

.ticket-content .alert {
  padding: 10px 12px;
  background: rgba(255, 123, 123, 0.1);
  border: 1px solid rgba(255, 123, 123, 0.3);
  border-radius: 12px;
  color: #ff7b7b;
  font-size: 12px;
  margin-bottom: 16px;
  font-weight: 500;
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
  align-items: center;
  pointer-events: none;
}

.toolbar-actions {
  pointer-events: auto;
  display: flex;
  gap: 8px;
  align-items: center;
  background: rgba(12, 16, 27, 0.8);
  backdrop-filter: blur(12px);
  padding: 4px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.symbol-block {
  pointer-events: auto;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  background: rgba(12, 16, 27, 0.8);
  backdrop-filter: blur(12px);
  padding: 4px 8px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  height: 42px;
  box-sizing: border-box;
}

.quick-symbols {
  pointer-events: auto;
  display: flex;
  gap: 4px;
  align-items: center;
  background: rgba(12, 16, 27, 0.8);
  backdrop-filter: blur(12px);
  padding: 4px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  height: 42px;
  box-sizing: border-box;
  overflow-x: auto;
}

.quick-symbol-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #8fa1c4;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
  min-width: 32px;
  justify-content: center;
}

.quick-symbol-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
  color: #fff;
}

.quick-symbol-btn.active {
  background: rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  border-color: rgba(93, 247, 194, 0.5);
}

.symbol-code {
  font-weight: 700;
  font-size: 10px;
}

.arrow {
  font-size: 12px;
  font-weight: 700;
}

.arrow-up {
  color: #5df7c2;
}

.arrow-down {
  color: #ff7b7b;
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
  top: 88px;
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
  box-sizing: border-box;
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
  max-height: 350px;
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
  font-size: 18px;
  font-weight: 700;
  padding: 0 8px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-variant-numeric: tabular-nums;
  height: 100%;
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
  top: 50px;
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
  padding: 8px 12px;
  background: rgba(12, 16, 27, 0.9);
  border-radius: 12px;
  border: 1px solid rgba(93, 247, 194, 0.2);
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(93, 247, 194, 0.1);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.toolbar-blur:hover {
  background: rgba(12, 16, 27, 0.95);
  border-color: rgba(93, 247, 194, 0.4);
  box-shadow: 0 12px 32px rgba(93, 247, 194, 0.15), 0 0 0 1px rgba(93, 247, 194, 0.2);
}

.ghost {
  background: linear-gradient(135deg, rgba(93, 247, 194, 0.08), rgba(61, 255, 181, 0.04));
  border: 1px solid rgba(93, 247, 194, 0.25);
  color: #5df7c2;
  padding: 8px 14px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.ghost:hover {
  background: linear-gradient(135deg, rgba(93, 247, 194, 0.15), rgba(61, 255, 181, 0.1));
  border-color: rgba(93, 247, 194, 0.5);
  box-shadow: 0 0 16px rgba(93, 247, 194, 0.3);
  transform: translateY(-1px);
}

.ghost:active {
  transform: translateY(0);
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
  letter-spacing: 0.6px;
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
  gap: 14px;
  position: relative;
}

/* Reserve space so fixed dock does not block cards or chart */
@media (min-width: 1025px) {
  .workspace {
    padding-right: calc(var(--dock-width) + var(--dock-offset) + 2px);
    box-sizing: border-box;
  }
}

@media (max-width: 1024px) {
  .workspace {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}

@media (max-width: 768px) {
  .chart-glass {
    min-height: 350px;
  }
  
  .trade-ticket-floating {
    bottom: 12px;
    max-width: calc(100vw - 24px);
  }
  
  .ticket-content {
    padding: 12px 16px;
  }
}

.dock-rail {
  display: flex;
  flex-direction: column;
  gap: 4px;
  background: linear-gradient(180deg, rgba(24, 29, 41, 0.9), rgba(18, 20, 28, 0.9));
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 8px 6px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.28);
  position: fixed;
  top: 150px;
  right: var(--dock-offset);
  transform: translateY(0);
  z-index: 90;
  width: var(--dock-width);
  backdrop-filter: blur(14px);
  overflow: visible;
}

.dock-title {
  font-size: 11px;
  color: #6b7280;
  letter-spacing: 0.6px;
  padding: 2px 8px;
  text-transform: uppercase;
  text-align: center;
}

.dock-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  width: auto;
  min-width: 100%;
  border: 1px solid transparent;
  background: transparent;
  color: #8fa1c4;
  padding: 6px 6px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  flex-direction: column;
  position: relative;
  overflow: visible;
}

.dock-badge {
  position: absolute;
  top: 0px;
  right: 0px;
  min-width: 20px;
  height: 20px;
  padding: 0 4px;
  background: linear-gradient(135deg, #5df7c2 0%, #3dffb5 100%);
  color: #0b0e14;
  font-size: 9px;
  font-weight: 700;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(11, 14, 20, 0.95);
  box-shadow: 0 2px 8px rgba(93, 247, 194, 0.4);
  pointer-events: none;
  line-height: 1;
}

.dock-badge-active {
  background: linear-gradient(135deg, #5df7c2 0%, #3dffb5 100%);
  box-shadow: 0 2px 12px rgba(93, 247, 194, 0.5);
}

.dock-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}

.dock-btn.active {
  background: rgba(93, 247, 194, 0.12);
  border-color: rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  box-shadow: 0 0 0 1px rgba(93, 247, 194, 0.1);
}

.dock-text {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  text-align: center;
}

.dock-label {
  font-size: 10px;
  font-weight: 700;
}

.side-panel-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
  width: 100%;
}

@media (max-width: 1200px) {
  .dock-rail {
    right: 8px;
  }
}

@media (max-width: 1024px) {
  .dock-rail {
    position: static;
    transform: none;
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
  }
  .dock-btn {
    flex: 1;
    flex-direction: row;
    justify-content: center;
  }
  .dock-text {
    align-items: flex-start;
  }
}

.signal-list,
.social-list,
.pending-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* Signal Tabs Header (in card-head) */
.signal-tabs-header {
  display: flex;
  gap: 8px;
  margin-left: auto;
}

.signal-tab-header-btn {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #8fa1c4;
  border-radius: 6px;
  cursor: pointer;
  font-size: 11px;
  font-weight: 600;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.signal-tab-header-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.2);
  color: #fff;
}

.signal-tab-header-btn.active {
  background: rgba(93, 247, 194, 0.15);
  border-color: rgba(93, 247, 194, 0.4);
  color: #5df7c2;
}

/* Signal Tabs (keep for backwards compatibility if needed) */
.signal-tabs {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(255, 255, 255, 0.02);
}

.signal-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #8fa1c4;
  border-radius: 8px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.signal-tab:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.2);
  color: #fff;
}

.signal-tab.active {
  background: rgba(93, 247, 194, 0.15);
  border-color: rgba(93, 247, 194, 0.4);
  color: #5df7c2;
}

/* Market Trends View */
.market-trends-view {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
  max-height: 500px;
  overflow-y: auto;
}

.timeframe-selector {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.tf-btn {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #8fa1c4;
  border-radius: 6px;
  cursor: pointer;
  font-size: 11px;
  font-weight: 600;
  transition: all 0.2s;
}

.tf-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.2);
}

.tf-btn.active {
  background: rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  border-color: rgba(93, 247, 194, 0.4);
}

.market-trends-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0;
}

.section-header {
  font-size: 12px;
  font-weight: 700;
  color: #8fa1c4;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 0 8px 0;
  border-bottom: 1px solid rgba(93, 247, 194, 0.1);
}

/* 时间框架网格 (参考图2: 3x4 布局) */
.timeframe-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.timeframe-btn {
  padding: 10px 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #8fa1c4;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}

.timeframe-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
  color: #fff;
}

.timeframe-btn.active {
  background: rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  border-color: rgba(93, 247, 194, 0.4);
  box-shadow: 0 0 8px rgba(93, 247, 194, 0.2);
}

/* 交易标的列表 */
.pairs-list-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 360px;
  overflow-y: auto;
  padding-right: 4px;
}

/* 美化滚动条 - Chrome/Safari/Edge */
.pairs-list-container::-webkit-scrollbar {
  width: 8px;
}

.pairs-list-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 4px;
}

.pairs-list-container::-webkit-scrollbar-thumb {
  background: rgba(93, 247, 194, 0.25);
  border-radius: 4px;
  transition: background 0.2s;
}

.pairs-list-container::-webkit-scrollbar-thumb:hover {
  background: rgba(93, 247, 194, 0.45);
}

/* Firefox 滚动条 */
.pairs-list-container {
  scrollbar-width: thin;
  scrollbar-color: rgba(93, 247, 194, 0.25) rgba(255, 255, 255, 0.02);
}

.pair-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 10px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.pair-list-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
  transform: translateX(4px);
}

.pair-list-item.active {
  background: rgba(93, 247, 194, 0.15);
  border-color: rgba(93, 247, 194, 0.5);
  box-shadow: inset 0 0 12px rgba(93, 247, 194, 0.1), 0 0 8px rgba(93, 247, 194, 0.15);
}

.pair-main {
  display: flex;
  align-items: baseline;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.pair-symbol {
  font-weight: 700;
  font-size: 13px;
  color: #fff;
  white-space: nowrap;
}

.pair-display {
  font-size: 11px;
  color: #8fa1c4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pair-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
  font-weight: 600;
  font-size: 12px;
  white-space: nowrap;
  flex-shrink: 0;
}

.pair-trend.trend-up {
  color: #5df7c2;
}

.pair-trend.trend-up .trend-arrow {
  color: #5df7c2;
  font-weight: 800;
}

.pair-trend.trend-down {
  color: #ff7b7b;
}

.pair-trend.trend-down .trend-arrow {
  color: #ff7b7b;
  font-weight: 800;
}

.trend-arrow {
  font-size: 14px;
  line-height: 1;
}

.market-pairs-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.market-pair-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  transition: all 0.2s;
}

.market-pair-card:hover {
  background: rgba(93, 247, 194, 0.1);
  border-color: rgba(93, 247, 194, 0.3);
  box-shadow: 0 4px 12px rgba(93, 247, 194, 0.1);
}

.pair-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.pair-name {
  font-size: 12px;
  font-weight: 700;
  color: #fff;
}

.pair-switch-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  background: rgba(93, 247, 194, 0.15);
  border: 1px solid rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pair-switch-btn:hover {
  background: rgba(93, 247, 194, 0.25);
  border-color: rgba(93, 247, 194, 0.4);
  transform: scale(1.05);
}

.trend-for-timeframe {
  display: flex;
  gap: 10px;
}

.trend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
}

.trend-label {
  font-size: 10px;
  color: #6b7a99;
  font-weight: 600;
  min-width: 20px;
  text-align: center;
}

.trend-strength {
  font-size: 16px;
  font-weight: 700;
  flex: 1;
}

.trend-strength.up {
  color: #5df7c2;
}

.trend-strength.down {
  color: #ff7b7b;
}

/* Signal Table Row - Compact Professional Design */
.signal-row,
.pending-row,
.social-row {
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr 1.5fr;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.01);
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  padding: 12px 16px;
  gap: 12px;
  transition: background 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: default;
}

.signal-row:last-child {
  border-bottom: none;
}

.signal-row:hover {
  background: rgba(93, 247, 194, 0.04);
}

.signal-row.signal-expired {
  opacity: 0.5;
  background: rgba(255, 255, 255, 0.01);
}

.signal-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

/* Symbol Cell */
.symbol-cell {
  justify-content: flex-start;
}

.symbol-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(93, 247, 194, 0.15), rgba(93, 247, 194, 0.05));
  border: 1px solid rgba(93, 247, 194, 0.3);
  border-radius: 6px;
  padding: 6px 10px;
  font-weight: 600;
  color: #5df7c2;
  font-size: 12px;
  letter-spacing: 0.5px;
}

/* Trend Cell */
.trend-cell {
  justify-content: center;
}

.trend-indicator {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.03);
  transition: all 0.2s;
}

.trend-indicator.up {
  color: #5df7c2;
  background: rgba(93, 247, 194, 0.08);
}

.trend-indicator.down {
  color: #ff7b7b;
  background: rgba(255, 123, 123, 0.08);
}

.trend-indicator svg {
  flex-shrink: 0;
}

.trend-arrows {
  font-weight: 700;
  font-size: 13px;
  letter-spacing: 0px;
}

/* Copies Cell */
.copies-cell {
  justify-content: center;
}

/* Action Cell */
.action-cell {
  justify-content: flex-end;
}

.btn-follow {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4px 12px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border: 1.5px solid currentColor;
  height: 38px;
  min-width: 80px;
}

.btn-action-text {
  font-size: 12px;
  line-height: 1.2;
  color: #ffffff;
}

.btn-timer-text {
  font-size: 10px;
  color: #ffffff;
  opacity: 1;
  font-variant-numeric: tabular-nums;
  line-height: 1;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 3px;
  margin-top: 2px;
}

.btn-follow:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.btn-follow:active:not(:disabled) {
  transform: translateY(0);
}

.btn-follow:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-follow.btn-call {
  color: #fff;
  border-color: #5df7c2;
}

.btn-follow.btn-call:hover:not(:disabled) {
  background: rgba(93, 247, 194, 0.15);
}

.btn-follow.btn-put {
  color: #fff;
  border-color: #ff7b7b;
}

.btn-follow.btn-put:hover:not(:disabled) {
  background: rgba(255, 123, 123, 0.15);
}

.btn-expired {
  display: inline-flex;
  padding: 6px 12px;
  border-radius: 6px;
  background: rgba(255, 123, 123, 0.08);
  color: #ff7b7b;
  font-size: 12px;
  font-weight: 600;
}

/* Time Cell */
.time-cell {
  justify-content: center;
  color: #8fa1c4;
  font-size: 12px;
}

.time-text {
  font-weight: 500;
}

/* Symbol Wrapper & Badge */
.symbol-wrapper {
  display: flex;
  align-items: center;
}

.copies-badge-small {
  display: flex;
  align-items: center;
  gap: 2px;
  background: rgba(93, 247, 194, 0.1);
  border: 1px solid rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  font-size: 9px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
  line-height: 1;
  white-space: nowrap;
  margin-right: 8px;
}

/* Trend Arrows */
.trend-arrows {
  font-weight: 800;
  font-size: 14px;
  letter-spacing: -2px;
}
.trend-arrows.up { color: #5df7c2; }
.trend-arrows.down { color: #ff7b7b; }

/* Signal List Container */
.signal-list {
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(93, 247, 194, 0.1);
  overflow: hidden;
}

/* Symbol Text */
.symbol-text {
  font-weight: 700;
  color: #fff;
  font-size: 13px;
  letter-spacing: 0.5px;
}

/* Signal Controls */
.signal-controls {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background: rgba(93, 247, 194, 0.05);
  border: 1px solid rgba(93, 247, 194, 0.1);
  border-radius: 12px;
  flex-wrap: wrap;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.control-label {
  font-size: 11px;
  font-weight: 700;
  color: #8fa1c4;
  text-transform: uppercase;
  letter-spacing: 0.4px;
}

.control-select {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.control-select:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(93, 247, 194, 0.3);
}

.control-select:focus {
  outline: none;
  background: rgba(93, 247, 194, 0.1);
  border-color: rgba(93, 247, 194, 0.5);
}

.filter-buttons {
  display: flex;
  gap: 4px;
}

.filter-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #8fa1c4;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

.filter-btn.active {
  background: rgba(93, 247, 194, 0.2);
  color: #5df7c2;
  border-color: rgba(93, 247, 194, 0.4);
}

.filter-btn.filter-call {
  border-color: rgba(16, 185, 129, 0.2);
}

.filter-btn.filter-call.active {
  background: rgba(16, 185, 129, 0.2);
  color: #5df7c2;
  border-color: rgba(16, 185, 129, 0.4);
}

.filter-btn.filter-put {
  border-color: rgba(239, 68, 68, 0.2);
}

.filter-btn.filter-put.active {
  background: rgba(239, 68, 68, 0.2);
  color: #ff7b7b;
  border-color: rgba(239, 68, 68, 0.4);
}

.signal-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px 20px;
  color: #6b7a99;
  text-align: center;
}

/* Removed: Old Signal Row Content CSS - replaced with modern table design */

/* Removed old signal CSS - replaced with new table design */

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: linear-gradient(135deg, #5df7c2, #5fa7ff);
  display: grid;
  place-items: center;
  color: #0b0e14;
  font-weight: 800;
  letter-spacing: 0.5px;
}

.quick-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.quick-chip {
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
  border-radius: 12px;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  cursor: pointer;
  color: #e7ecf5;
  transition: all 0.2s;
}

.quick-chip:hover {
  border-color: rgba(93, 247, 194, 0.3);
  background: rgba(93, 247, 194, 0.06);
  color: #fff;
}

.quick-label {
  font-weight: 700;
}

.quick-sub {
  color: #8fa1c4;
  font-size: 12px;
}

.quick-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
  color: #8fa1c4;
  font-size: 12px;
}

.roi-meter {
  flex: 1;
  height: 8px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
  position: relative;
  overflow: hidden;
}

.roi-fill {
  height: 100%;
  background: linear-gradient(90deg, #5df7c2, #3dffb5);
}

.roi-label {
  position: absolute;
  top: -18px;
  right: 0;
  font-size: 11px;
  color: #5df7c2;
  font-weight: 700;
}

.shortcut-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.shortcut-row {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  padding: 10px 12px;
}

.shortcut-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.shortcut-name {
  font-weight: 700;
  color: #fff;
}

.shortcut-sub {
  color: #8fa1c4;
  font-size: 12px;
}

.keycap {
  min-width: 42px;
  padding: 6px 10px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  text-align: center;
  font-weight: 700;
  color: #e7ecf5;
  box-shadow: inset 0 -2px 0 rgba(0, 0, 0, 0.2);
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
  padding: 16px;
  display: flex;
  flex-direction: column;
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

.drawing-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.action-btn.primary {
  background: linear-gradient(135deg, #5df7c2 0%, #3dffb5 100%);
  color: #0a0e14;
  box-shadow: 0 4px 12px rgba(93, 247, 194, 0.3);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(93, 247, 194, 0.4);
}

.action-btn.secondary {
  background: linear-gradient(135deg, rgba(93, 247, 194, 0.15), rgba(61, 255, 181, 0.08));
  border: 1px solid rgba(93, 247, 194, 0.3);
  color: #5df7c2;
}

.action-btn.secondary:hover {
  background: linear-gradient(135deg, rgba(93, 247, 194, 0.25), rgba(61, 255, 181, 0.15));
  border-color: rgba(93, 247, 194, 0.5);
  box-shadow: 0 0 12px rgba(93, 247, 194, 0.2);
}

.menu-divider {
  height: 1px;
  background: linear-gradient(to right,
    rgba(255, 255, 255, 0),
    rgba(93, 247, 194, 0.2),
    rgba(255, 255, 255, 0)
  );
  margin: 16px 0;
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

/* Deposit Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-card {
  background: rgba(18, 20, 28, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  width: 90%;
  max-width: 480px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from { 
    opacity: 0;
    transform: translateY(20px);
  }
  to { 
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.modal-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
}

.modal-close {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #8fa1c4;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.modal-body {
  padding: 24px;
}

.deposit-input-group {
  margin-bottom: 20px;
}

.deposit-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #8fa1c4;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.deposit-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.03);
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.deposit-input-wrapper:focus-within {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(93, 247, 194, 0.5);
  box-shadow: 0 0 0 4px rgba(93, 247, 194, 0.1);
}

.currency-prefix {
  padding: 16px 0 16px 20px;
  font-size: 16px;
  font-weight: 700;
  color: #5df7c2;
  background: rgba(93, 247, 194, 0.1);
  border-right: 2px solid rgba(93, 247, 194, 0.2);
  min-width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.deposit-input {
  flex: 1;
  background: transparent;
  border: none;
  padding: 16px 20px;
  color: #fff;
  font-size: 24px;
  font-weight: 600;
  outline: none;
  font-variant-numeric: tabular-nums;
}

.deposit-input::placeholder {
  color: rgba(255, 255, 255, 0.3);
  font-weight: 500;
}

.deposit-input::-webkit-outer-spin-button,
.deposit-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.deposit-input[type=number] {
  appearance: textfield;
  -moz-appearance: textfield;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.quick-amounts {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  flex-wrap: wrap;
}

.quick-btn {
  flex: 1;
  min-width: 70px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.03);
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: #8fa1c4;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  overflow: hidden;
}

.quick-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(93, 247, 194, 0.1) 0%, rgba(61, 255, 181, 0.1) 100%);
  opacity: 0;
  transition: opacity 0.2s;
}

.quick-btn:hover {
  background: rgba(93, 247, 194, 0.05);
  border-color: rgba(93, 247, 194, 0.4);
  color: #5df7c2;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(93, 247, 194, 0.2);
}

.quick-btn:hover::before {
  opacity: 1;
}

.quick-btn:active {
  transform: translateY(0);
}

.btn-primary {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #5df7c2 0%, #3dffb5 100%);
  color: #0a0e14;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(93, 247, 194, 0.4);
}

.btn-secondary {
  flex: 1;
  padding: 12px 24px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #8fa1c4;
  border-radius: 10px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-call {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
}

.btn-call:hover {
  box-shadow: 0 8px 20px rgba(16, 185, 129, 0.3) !important;
}

.btn-put {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%) !important;
}

.btn-put:hover {
  box-shadow: 0 8px 20px rgba(239, 68, 68, 0.3) !important;
}
</style>
