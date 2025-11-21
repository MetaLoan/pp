package handler

import (
	"net/http"

	"github.com/MetaLoan/pp/internal/model"
	"github.com/MetaLoan/pp/internal/service"
	"github.com/gin-gonic/gin"
)

type TradeHandler struct {
	tradeService *service.TradeService
}

func NewTradeHandler() *TradeHandler {
	return &TradeHandler{
		tradeService: service.NewTradeService(),
	}
}

type PlaceOrderRequest struct {
	Symbol    string  `json:"symbol" binding:"required"`
	Direction string  `json:"direction" binding:"required,oneof=CALL PUT"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Duration  int     `json:"duration" binding:"required,min=30"` // seconds
	ClientPrice float64 `json:"client_price"` // optional: price user saw on UI
}

func (h *TradeHandler) PlaceOrder(c *gin.Context) {
	// Get User from Context (set by Auth Middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := user.(*model.User).ID

	var req PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.tradeService.PlaceOrder(userID, req.Symbol, req.Direction, req.Amount, req.Duration, req.ClientPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order placed", "order": order})
}

func (h *TradeHandler) GetActiveOrders(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := user.(*model.User).ID

	orders, err := h.tradeService.GetActiveOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *TradeHandler) GetBalance(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := user.(*model.User).ID

	currency := c.DefaultQuery("currency", "USDT")
	wallet, err := h.tradeService.GetWallet(userID, currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"currency": wallet.Currency,
		"balance":  wallet.Balance,
		"frozen":   wallet.Frozen,
	})
}
