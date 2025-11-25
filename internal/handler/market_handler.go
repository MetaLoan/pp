package handler

import (
	"net/http"
	"strconv"

	"github.com/MetaLoan/pp/internal/core"
	"github.com/gin-gonic/gin"
)

type MarketHandler struct{}

func NewMarketHandler() *MarketHandler {
	return &MarketHandler{}
}

// GetPrice returns the latest cached price for a symbol (debug endpoint)
func (h *MarketHandler) GetPrice(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(400, gin.H{"error": "symbol required"})
		return
	}
	// core.GlobalMarket holds lastPrice map
	if d, ok := core.GetLastPrice(symbol); ok {
		c.JSON(200, gin.H{"symbol": symbol, "price": d.Price, "timestamp": d.Timestamp})
		return
	}
	c.JSON(404, gin.H{"error": "no price for symbol"})
}

func (h *MarketHandler) GetCandles(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol required"})
		return
	}
	intervalSec := 5
	if v := c.Query("interval"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			intervalSec = parsed
		}
	}
	limit := 200
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	candles := core.GlobalMarket.GetCandles(symbol, intervalSec, limit)
	c.JSON(http.StatusOK, gin.H{"candles": candles})
}

// GetHistory returns raw server-side ticks for a symbol (debug)
func (h *MarketHandler) GetHistory(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol required"})
		return
	}
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	hist := core.GetPriceHistory(symbol, limit)
	// convert to simple JSON array
	out := make([]map[string]interface{}, 0, len(hist))
	for _, p := range hist {
		out = append(out, map[string]interface{}{"symbol": p.Symbol, "price": p.Price, "timestamp": p.Timestamp})
	}
	c.JSON(http.StatusOK, gin.H{"history": out})
}
