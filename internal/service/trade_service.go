package service

import (
	"errors"
	"math"
	"time"

	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
	"github.com/MetaLoan/pp/internal/repository"
	"github.com/google/uuid"
)

type TradeService struct {
	orderRepo  *repository.OrderRepository
	walletRepo *repository.WalletRepository
}

func NewTradeService() *TradeService {
	return &TradeService{
		orderRepo:  repository.NewOrderRepository(),
		walletRepo: repository.NewWalletRepository(),
	}
}

func (s *TradeService) PlaceOrder(userID int64, symbol, direction string, amount float64, durationSeconds int, clientPrice float64) (*model.Order, error) {
	// 1. Get Current Price (server view)
	serverPrice, err := core.GlobalMarket.GetCurrentPrice(symbol)
	if err != nil {
		return nil, errors.New("failed to get market price")
	}
	openPrice := serverPrice
	if clientPrice > 0 && math.Abs(clientPrice-serverPrice) <= 0.001 { // within 0.1 pip, trust client view
		openPrice = clientPrice
	}

	// 2. Check Balance & Freeze Funds
	// Assuming USDT for now
	currency := "USDT"
	if err := s.walletRepo.FreezeFunds(userID, currency, amount); err != nil {
		return nil, err
	}

	// 3. Create Order
	payoutRate := 0.85 // Fixed 85% payout for MVP
	now := time.Now()
	closeTime := now.Add(time.Duration(durationSeconds) * time.Second)

	order := &model.Order{
		OrderNo:     uuid.New().String(),
		UserID:      userID,
		AssetSymbol: symbol,
		Direction:   direction,
		Amount:      amount,
		PayoutRate:  payoutRate,
		OpenPrice:   openPrice,
		OpenTime:    now,
		CloseTime:   closeTime,
		Status:      "active",
	}

	if err := s.orderRepo.CreateOrder(order); err != nil {
		// TODO: Rollback frozen funds if order creation fails
		return nil, err
	}

	// 4. Schedule Settlement (Simplified: using goroutine)
	go func() {
		time.Sleep(time.Duration(durationSeconds) * time.Second)
		s.SettleOrder(order.ID)
	}()

	return order, nil
}

func (s *TradeService) SettleOrder(orderID int64) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != "active" {
		return nil
	}

	// Get Closing Price
	closePrice, _ := core.GlobalMarket.GetCurrentPrice(order.AssetSymbol)
	order.ClosePrice = closePrice

	// Determine Win/Loss
	isWin := false
	if order.Direction == "CALL" && closePrice > order.OpenPrice {
		isWin = true
	} else if order.Direction == "PUT" && closePrice < order.OpenPrice {
		isWin = true
	}

	// Calculate PnL and Settlement Amount
	settlementAmount := 0.0
	if isWin {
		profit := order.Amount * order.PayoutRate
		order.PnL = profit
		order.Status = "won"
		settlementAmount = order.Amount + profit
	} else {
		order.PnL = -order.Amount
		order.Status = "lost"
		settlementAmount = 0.0
	}

	// Update Wallet
	// We always unfreeze the principal (order.Amount).
	// If won, we add principal + profit back to balance.
	// If lost, we add 0 back to balance (principal is consumed).
	// Wait, my UnfreezeAndSettle logic was: balance = balance + profitAmount.
	// So if won: profitAmount = principal + profit.
	// If lost: profitAmount = 0.
	if err := s.walletRepo.UnfreezeAndSettle(order.UserID, "USDT", order.Amount, settlementAmount); err != nil {
		// Log critical error: money stuck in frozen
		return err
	}

	// Update Order
	return s.orderRepo.UpdateOrder(order)
}

func (s *TradeService) GetActiveOrders(userID int64) ([]model.Order, error) {
	return s.orderRepo.GetOrdersByUserIDAndStatus(userID, "active")
}

func (s *TradeService) GetWallet(userID int64, currency string) (*model.Wallet, error) {
	return s.walletRepo.GetWallet(userID, currency)
}
