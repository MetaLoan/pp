package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/MetaLoan/pp/internal/config"
	"github.com/MetaLoan/pp/internal/model"
	"github.com/MetaLoan/pp/internal/repository"
)

type FundsService struct {
	transferRepo *repository.TransferRepository
	walletRepo   *repository.WalletRepository
}

func NewFundsService() *FundsService {
	return &FundsService{
		transferRepo: repository.NewTransferRepository(),
		walletRepo:   repository.NewWalletRepository(),
	}
}

func (s *FundsService) RequestDeposit(userID int64, currency string, amount float64) (*model.TransferRequest, error) {
	if err := validateAmount(amount, config.GlobalConfig.Funds.Deposit.MinAmount, config.GlobalConfig.Funds.Deposit.MaxAmount); err != nil {
		return nil, err
	}
	if err := s.checkLimits(userID, "deposit", amount); err != nil {
		return nil, err
	}
	req := &model.TransferRequest{
		UserID:   userID,
		Type:     "deposit",
		Currency: currency,
		Amount:   amount,
		Status:   "pending",
	}
	if err := s.transferRepo.Create(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *FundsService) RequestWithdraw(userID int64, currency string, amount float64) (*model.TransferRequest, error) {
	if err := validateAmount(amount, config.GlobalConfig.Funds.Withdraw.MinAmount, config.GlobalConfig.Funds.Withdraw.MaxAmount); err != nil {
		return nil, err
	}
	if err := s.checkLimits(userID, "withdraw", amount); err != nil {
		return nil, err
	}
	// Check balance
	wallet, err := s.walletRepo.GetWallet(userID, currency)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}
	req := &model.TransferRequest{
		UserID:   userID,
		Type:     "withdraw",
		Currency: currency,
		Amount:   amount,
		Status:   "pending",
	}
	if err := s.transferRepo.Create(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *FundsService) Review(userID int64, id int64, approve bool, reviewer string, remark string) (*model.TransferRequest, error) {
	req, err := s.transferRepo.GetPendingByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if approve {
		if req.Type == "deposit" {
			if err := s.walletRepo.Credit(req.UserID, req.Currency, req.Amount, buildRef(req), "deposit approved", "deposit"); err != nil {
				return nil, err
			}
		} else if req.Type == "withdraw" {
			withdrawCfg := config.GlobalConfig.Funds.Withdraw
			fee := req.Amount * withdrawCfg.FeeRate
			if fee < withdrawCfg.MinFee {
				fee = withdrawCfg.MinFee
			}
			total := req.Amount + fee
			if err := s.walletRepo.Debit(req.UserID, req.Currency, req.Amount, buildRef(req), "withdraw approved", "withdraw"); err != nil {
				return nil, err
			}
			if fee > 0 {
				if err := s.walletRepo.Debit(req.UserID, req.Currency, fee, buildRef(req), "withdraw fee", "withdraw_fee"); err != nil {
					return nil, err
				}
			}
			remark = fmt.Sprintf("withdraw approved, total %.2f (fee %.2f)", total, fee)
		}
		req.Status = "approved"
	} else {
		req.Status = "rejected"
	}
	req.ReviewedAt = now
	req.Reviewer = reviewer
	req.Remark = remark
	if err := s.transferRepo.Update(req); err != nil {
		return nil, err
	}
	return req, nil
}

func validateAmount(amount, min, max float64) error {
	if min > 0 && amount < min {
		return errors.New("amount below minimum")
	}
	if max > 0 && amount > max {
		return errors.New("amount exceeds maximum")
	}
	return nil
}

func buildRef(req *model.TransferRequest) string {
	return fmt.Sprintf("transfer:%d", req.ID)
}

func (s *FundsService) checkLimits(userID int64, typ string, amount float64) error {
	var cfg struct {
		DailyMax float64
		Cooldown string
	}
	switch typ {
	case "deposit":
		cfg.DailyMax = config.GlobalConfig.Funds.Deposit.DailyMax
		cfg.Cooldown = config.GlobalConfig.Funds.Deposit.Cooldown
	case "withdraw":
		cfg.DailyMax = config.GlobalConfig.Funds.Withdraw.DailyMax
		cfg.Cooldown = config.GlobalConfig.Funds.Withdraw.Cooldown
	default:
		return errors.New("invalid type")
	}

	// Cooldown
	if cfg.Cooldown != "" {
		if cd, err := time.ParseDuration(cfg.Cooldown); err == nil && cd > 0 {
			last, err := s.transferRepo.LastRequestTime(userID, typ)
			if err == nil && !last.IsZero() && time.Since(last) < cd {
				return errors.New("cooldown in effect, try later")
			}
		}
	}

	// Daily limit (approved amounts)
	if cfg.DailyMax > 0 {
		startOfDay := time.Now().Truncate(24 * time.Hour)
		total, err := s.transferRepo.SumApprovedAmountSince(userID, typ, startOfDay)
		if err != nil {
			return err
		}
		if total+amount > cfg.DailyMax {
			return errors.New("daily limit exceeded")
		}
	}

	return nil
}
