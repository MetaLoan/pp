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
                      <option value="30">30s</option>
                      <option value="60">60s</option>
                      <option value="300">5m</option>
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
                <div class="control-group">
                  <span class="control-label">周期</span>
                  <div class="filter-buttons">
                    <button :class="['filter-btn', { active: signalFilterTiming === 'all' }]" @click="signalFilterTiming = 'all'" style="width: 40px">
                      全
                    </button>
                    <button v-for="t in ['1m', '2m', '3m', '4m', '5m']" :key="t" :class="['filter-btn', { active: signalFilterTiming === t }]" @click="signalFilterTiming = t" style="width: 40px">
                      {{ t }}
                    </button>
                  </div>
                </div>
              </div>

              <div class="signal-list">
                <!-- Signal Table Header -->
                <div class="signal-header-row">
                  <div class="signal-cell symbol-cell"></div>
                  <div class="signal-cell trend-cell"></div>
                  <div class="signal-cell copies-cell"></div>
                  <div class="signal-cell expiry-cell"></div>
                  <div class="signal-cell time-cell"></div>
                  <div class="signal-cell action-cell"></div>
                </div>

                <!-- Signal Rows -->
                <div v-for="sig in filteredSignals" :key="sig.title" :class="['signal-row', { 'signal-expired': !getSignalValidity(sig).isValid }]">
                  <!-- 交易标的 -->
                  <div class="signal-cell symbol-cell">
                    <span class="symbol-badge">{{ sig.symbol }}</span>
                  </div>

                  <!-- 趋势强度 -->
                  <div class="signal-cell trend-cell">
                    <div :class="['trend-indicator', { up: sig.confidence >= 0.5, down: sig.confidence < 0.5 }]">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                        <polyline v-if="sig.confidence >= 0.5" points="23 6 13.5 15.5 8.5 10.5 1 18"></polyline>
                        <polyline v-else points="23 18 13.5 8.5 8.5 13.5 1 6"></polyline>
                      </svg>
                      <span class="trend-arrows">{{ sig.confidence >= 0.8 ? '↑↑↑' : sig.confidence >= 0.5 ? '↑↑' : sig.confidence >= 0.3 ? '↓' : '↓↓' }}</span>
                    </div>
                  </div>

                  <!-- 已复制数量 -->
                  <div class="signal-cell copies-cell">
                    <span class="copies-badge">{{ sig.copied }}</span>
                  </div>

                  <!-- 到期时间 -->
                  <div class="signal-cell expiry-cell">
                    <div v-if="getSignalValidity(sig).isValid" class="expiry-timer">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"></circle>
                        <polyline points="12 6 12 12 16 14"></polyline>
                      </svg>
                      <span>{{ Math.ceil(getSignalValidity(sig).remaining / 1000) }}s</span>
                    </div>
                    <span v-else class="expiry-badge expired">已过期</span>
                  </div>

                  <!-- 产生时间（相对时间） -->
                  <div class="signal-cell time-cell">
                    <span class="time-text">{{ formatSignalTime(sig.createdAt) }}</span>
                  </div>

                  <!-- 跟随按钮 (放在最后) -->
                  <div class="signal-cell action-cell">
                    <button 
                      v-if="getSignalValidity(sig).isValid"
                      :class="['btn-follow', sig.action === 'CALL' ? 'btn-call' : 'btn-put']" 
                      @click="handleSignalTrade(sig)"
                      :title="`跟随 ${sig.action} 信号`">
                      {{ sig.action }}
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
            
            <!-- Tab 2: Market Trends by Timeframe -->
            <div v-show="activeSignalTab === 'markets'" class="market-trends-view">
              <div class="timeframe-selector">
                <button v-for="tf in ['1m', '2m', '3m', '4m', '5m']" :key="tf" 
                  :class="['tf-btn', { active: signalFilterTiming === tf || (signalFilterTiming === 'all' && tf === '1m') }]"
                  @click="signalFilterTiming = tf">
                  {{ tf }}
                </button>
              </div>
              
              <div class="market-pairs-grid">
                <div v-for="marketData in marketTrends" :key="marketData.symbol" class="market-pair-card">
                  <div class="pair-header">
                    <span class="pair-name">{{ marketData.display }}</span>
                    <button class="pair-switch-btn" @click="selectedSymbol = marketData.symbol; activeRightModule = 'orders'" title="切换到此标的">
                      →
                    </button>
                  </div>
                  <div class="trend-for-timeframe">
                    <div v-if="signalFilterTiming === 'all' || signalFilterTiming === '1m'" class="trend-item">
                      <span class="trend-label">1m</span>
                      <span :class="['trend-strength', marketData.trends['1m'].direction > 0 ? 'up' : 'down']">
                        {{ getTrendArrows(marketData.trends['1m'].strength, marketData.trends['1m'].direction) }}
                      </span>
                    </div>
                    <div v-else class="trend-item">
                      <span class="trend-label">{{ signalFilterTiming }}</span>
                      <span :class="['trend-strength', marketData.trends[signalFilterTiming].direction > 0 ? 'up' : 'down']">
                        {{ getTrendArrows(marketData.trends[signalFilterTiming].strength, marketData.trends[signalFilterTiming].direction) }}
                      </span>
                    </div>
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
const timeframe = ref(5); // seconds per bar
const showSMA = ref(false);
const showEMA = ref(false);
const smaPeriod = ref(10);
const timeframeOptions = [1, 5, 15, 30, 60, 300, 600]; // Keep for logic mapping

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

const signalFeed = ref([
  { title: 'EUR/USD 突破', metric: '动量 +1.2σ', confidence: 0.82, action: 'CALL', timing: '5m', symbol: 'EURUSD', amount: 50, duration: 300, copied: 1240, createdAt: Date.now() - 180000, validity: 600000, expiryTime: Date.now() - 180000 + 600000 },
  { title: 'BTC/USDT 回踩', metric: 'RSI 34 · 趋势向上', confidence: 0.74, action: 'CALL', timing: '3m', symbol: 'BTCUSDT', amount: 100, duration: 180, copied: 856, createdAt: Date.now() - 120000, validity: 180000, expiryTime: Date.now() - 120000 + 180000 },
  { title: 'XAU/USD 拐点', metric: '布林中轨反弹', confidence: 0.68, action: 'PUT', timing: '1m', symbol: 'XAUUSD', amount: 25, duration: 60, copied: 543, createdAt: Date.now() - 60000, validity: 60000, expiryTime: Date.now() - 60000 + 60000 },
]);

// Signal module tabs
const activeSignalTab = ref('signals'); // 'signals' | 'markets'

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
const filteredSignals = computed(() => {
  let signals = signalFeed.value;

  // 按交易方向过滤
  if (signalFilterAction.value !== 'all') {
    signals = signals.filter(s => s.action === signalFilterAction.value);
  }

  // 按标的过滤（仅显示当前选中标的的信号）
  signals = signals.filter(s => s.symbol === selectedSymbol.value);

  // 按时间框架过滤
  if (signalFilterTiming.value !== 'all') {
    signals = signals.filter(s => s.timing === signalFilterTiming.value);
  }

  // 按条件排序
  const sorted = [...signals];
  switch (signalSortBy.value) {
    case 'confidence':
      // 从高到低排序
      sorted.sort((a, b) => b.confidence - a.confidence);
      break;
    case 'timing':
      // 时间框架排序(数字小的在前)
      const timingOrder = { '1m': 1, '2m': 2, '3m': 3, '4m': 4, '5m': 5 };
      sorted.sort((a, b) => {
        const aVal = timingOrder[a.timing] || 999;
        const bVal = timingOrder[b.timing] || 999;
        return aVal - bVal;
      });
      break;
    case 'new':
      // 新信号优先
      sorted.sort((a, b) => (b.isNew ? 1 : 0) - (a.isNew ? 1 : 0));
      break;
  }

  return sorted;
});

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
  const arrowCount = strength >= 0.5 ? 2 : 1;
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
    currentTime.value = Date.now(); // 更新当前时间用于倒计时
  }, 1000);

  historyInterval = setInterval(() => {
    fetchHistory(true);
  }, 5000);

  restartCandleInterval();

  // 信号推送间隔 - 模拟实时信号
  signalInterval = setInterval(() => {
    pushNewSignal();
  }, 12000); // 每12秒推送一条新信号

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
    
    // 刷新数据
    marketStore.fetchActiveOrders();
    marketStore.fetchBalance();
    
    // 显示成功反馈并切换到订单面板
    errorMsg.value = `✓ 信号交易执行: ${signal.action} ${signal.title} (${orderDuration}s)`;
    activeRightModule.value = 'orders';
    
    // 3秒后清除提示
    setTimeout(() => {
      errorMsg.value = '';
    }, 3000);
  } catch (error) {
    errorMsg.value = `⚠ 下单失败: ${error.response?.data?.error || error.message || 'Unknown error'}`;
  }
};


const pushNewSignal = () => {
  // 新信号数据池
  const now = Date.now();
  const signalPool = [
    { title: 'GBP/USD 上破', metric: 'RSI 65 · 突破阻力', confidence: 0.79, action: 'CALL', timing: '2m', symbol: 'GBPUSD', amount: 40, duration: 120, copied: 0, createdAt: now, validity: 120000, expiryTime: now + 120000 },
    { title: 'ETH/USDT 下行', metric: '布林下轨测试', confidence: 0.75, action: 'PUT', timing: '3m', symbol: 'ETHUSDT', amount: 75, duration: 180, copied: 0, createdAt: now, validity: 180000, expiryTime: now + 180000 },
    { title: 'Gold 反弹', metric: '动量 +0.8σ', confidence: 0.71, action: 'CALL', timing: '5m', symbol: 'XAUUSD', amount: 30, duration: 300, copied: 0, createdAt: now, validity: 300000, expiryTime: now + 300000 },
    { title: 'US100 均线黄金交叉', metric: 'MA20 穿 MA50', confidence: 0.86, action: 'CALL', timing: '4m', symbol: 'US100', amount: 60, duration: 240, copied: 0, createdAt: now, validity: 240000, expiryTime: now + 240000 },
    { title: 'Oil 下跌', metric: 'MACD 负值扩大', confidence: 0.73, action: 'PUT', timing: '2m', symbol: 'WTIUSD', amount: 45, duration: 120, copied: 0, createdAt: now, validity: 120000, expiryTime: now + 120000 },
    { title: 'USDJPY 震荡', metric: 'RSI 50 · 中性信号', confidence: 0.68, action: 'CALL', timing: '1m', symbol: 'USDJPY', amount: 20, duration: 60, copied: 0, createdAt: now, validity: 60000, expiryTime: now + 60000 },
    { title: 'S&P500 新高', metric: 'ADX 强势向上', confidence: 0.81, action: 'CALL', timing: '5m', symbol: 'US500', amount: 80, duration: 300, copied: 0, createdAt: now, validity: 300000, expiryTime: now + 300000 },
  ];

  // 随机选择一个新信号
  const newSignal = signalPool[Math.floor(Math.random() * signalPool.length)];
  
  // 检查是否已存在相同的信号(避免重复)
  const exists = signalFeed.value.some(sig => sig.title === newSignal.title);
  if (exists) {
    return;
  }

  // 在列表头部插入新信号
  signalFeed.value.unshift({
    ...newSignal,
    isNew: true, // 标记为新信号，用于动画
  });

  // 移除 isNew 标记以停止动画(500ms后)
  setTimeout(() => {
    const idx = signalFeed.value.findIndex(s => s.title === newSignal.title);
    if (idx !== -1) {
      signalFeed.value[idx].isNew = false;
    }
  }, 500);

  // 限制信号列表不超过20条
  if (signalFeed.value.length > 20) {
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
  padding: 10px 12px;
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
  width: 100%;
  border: 1px solid transparent;
  background: transparent;
  color: #8fa1c4;
  padding: 6px 6px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  flex-direction: column;
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
  grid-template-columns: 1.2fr 1.2fr 1fr 1fr 1.3fr 1.3fr;
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

.copies-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 4px 10px;
  border-radius: 4px;
  background: rgba(93, 247, 194, 0.1);
  color: #8fa1c4;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid rgba(93, 247, 194, 0.2);
}

/* Action Cell */
.action-cell {
  justify-content: flex-end;
}

.btn-follow {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px 14px;
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
  color: #5df7c2;
  border-color: #5df7c2;
}

.btn-follow.btn-call:hover:not(:disabled) {
  background: rgba(93, 247, 194, 0.15);
}

.btn-follow.btn-put {
  color: #ff7b7b;
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

/* Expiry Cell */
.expiry-cell {
  justify-content: center;
}

.expiry-timer {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.03);
  color: #8fa1c4;
  font-size: 12px;
  font-weight: 500;
}

.expiry-timer svg {
  color: #5df7c2;
}

.expiry-badge {
  display: inline-flex;
  padding: 4px 10px;
  border-radius: 4px;
  background: rgba(255, 123, 123, 0.08);
  color: #ff7b7b;
  font-size: 11px;
  font-weight: 600;
}

.expiry-badge.expired {
  opacity: 0.7;
}

/* Time Cell */
.time-cell {
  justify-content: flex-end;
  color: #8fa1c4;
  font-size: 12px;
}

.time-text {
  font-weight: 500;
}

/* Signal List Container */
.signal-list {
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(93, 247, 194, 0.1);
  overflow: hidden;
}

/* Signal Table Header */
.signal-header-row {
  display: grid;
  grid-template-columns: 1.2fr 1.2fr 1fr 1fr 1.3fr 1.3fr;
  align-items: center;
  padding: 12px 16px;
  gap: 12px;
  background: linear-gradient(90deg, rgba(93, 247, 194, 0.08), rgba(93, 247, 194, 0.03));
  border-bottom: 1px solid rgba(93, 247, 194, 0.15);
  font-size: 12px;
  font-weight: 700;
  color: #8fa1c4;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.header-label {
  display: flex;
  align-items: center;
  gap: 4px;
  opacity: 0.8;
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
