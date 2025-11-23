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
