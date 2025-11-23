package handler

import (
	"net/http"
	"strconv"

	"github.com/MetaLoan/pp/internal/model"
	"github.com/MetaLoan/pp/internal/service"
	"github.com/gin-gonic/gin"
)

type FundsHandler struct {
	fundsService *service.FundsService
}

func NewFundsHandler(fundsService *service.FundsService) *FundsHandler {
	return &FundsHandler{fundsService: fundsService}
}

type fundsRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required"`
}

func (h *FundsHandler) Deposit(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := user.(*model.User).ID

	var req fundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tr, err := h.fundsService.RequestDeposit(userID, req.Currency, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"request": tr})
}

func (h *FundsHandler) Withdraw(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := user.(*model.User).ID

	var req fundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tr, err := h.fundsService.RequestWithdraw(userID, req.Currency, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"request": tr})
}

// Review is a placeholder; in real world should be protected by admin role.
func (h *FundsHandler) Review(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := user.(*model.User).ID

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	approve := c.Query("approve") == "true"
	remark := c.DefaultQuery("remark", "")
	reviewer := strconv.FormatInt(userID, 10)

	tr, err := h.fundsService.Review(userID, id, approve, reviewer, remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"request": tr})
}
