package repository

import (
	"errors"

	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
	"gorm.io/gorm"
)

type WalletRepository struct{}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{}
}

// createLedger records wallet change within an existing transaction.
func createLedger(tx *gorm.DB, userID int64, currency, typ string, amount, balance, frozen float64, ref, remark string) error {
	entry := model.WalletLedger{
		UserID:   userID,
		Currency: currency,
		Type:     typ,
		Amount:   amount,
		Balance:  balance,
		Frozen:   frozen,
		Ref:      ref,
		Remark:   remark,
	}
	return tx.Create(&entry).Error
}

func (r *WalletRepository) GetWallet(userID int64, currency string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := core.DB.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Auto create wallet for MVP if not exists
		wallet = model.Wallet{
			UserID:   userID,
			Currency: currency,
			Balance:  10000.0, // Give some demo money
		}
		core.DB.Create(&wallet)
		return &wallet, nil
	}
	return &wallet, err
}

func (r *WalletRepository) FreezeFunds(userID int64, currency string, amount float64) error {
	return core.DB.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet
		err := tx.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto create wallet for MVP if not exists (Fix for existing users)
			wallet = model.Wallet{
				UserID:   userID,
				Currency: currency,
				Balance:  10000.0,
			}
			if err := tx.Create(&wallet).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		// Optimistic locking with Version
		result := tx.Model(&wallet).Where("version = ?", wallet.Version).
			Updates(map[string]interface{}{
				"balance": gorm.Expr("balance - ?", amount),
				"frozen":  gorm.Expr("frozen + ?", amount),
				"version": gorm.Expr("version + 1"),
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("concurrent update conflict")
		}

		if err := tx.Where("id = ?", wallet.ID).First(&wallet).Error; err == nil {
			if err := createLedger(tx, userID, currency, "freeze", amount, wallet.Balance-amount, wallet.Frozen+amount, "", "freeze funds for order"); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *WalletRepository) UnfreezeAndSettle(userID int64, currency string, frozenAmount, profitAmount float64, ref string) error {
	return core.DB.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet
		if err := tx.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"frozen":  gorm.Expr("frozen - ?", frozenAmount),
			"balance": gorm.Expr("balance + ?", profitAmount), // profitAmount includes principal + profit
			"version": gorm.Expr("version + 1"),
		}

		result := tx.Model(&wallet).Where("version = ?", wallet.Version).Updates(updates)

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("concurrent update conflict")
		}

		if err := tx.Where("id = ?", wallet.ID).First(&wallet).Error; err == nil {
			typ := "settle"
			remark := "settlement"
			if profitAmount > frozenAmount {
				typ = "settle_win"
			} else if profitAmount == frozenAmount {
				typ = "settle_draw"
			} else {
				typ = "settle_lose"
			}
			if err := createLedger(tx, userID, currency, typ, profitAmount, wallet.Balance+profitAmount, wallet.Frozen-frozenAmount, ref, remark); err != nil {
				return err
			}
		}
		return nil
	})
}

// Credit increases balance with ledger.
func (r *WalletRepository) Credit(userID int64, currency string, amount float64, ref string, remark string, typ string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return core.DB.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet
		if err := tx.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error; err != nil {
			return err
		}
		result := tx.Model(&wallet).Where("version = ?", wallet.Version).
			Updates(map[string]interface{}{
				"balance": gorm.Expr("balance + ?", amount),
				"version": gorm.Expr("version + 1"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("concurrent update conflict")
		}
		if err := tx.Where("id = ?", wallet.ID).First(&wallet).Error; err == nil {
			ledgerType := typ
			if ledgerType == "" {
				ledgerType = "credit"
			}
			if err := createLedger(tx, userID, currency, ledgerType, amount, wallet.Balance+amount, wallet.Frozen, ref, remark); err != nil {
				return err
			}
		}
		return nil
	})
}

// Debit decreases balance with ledger.
func (r *WalletRepository) Debit(userID int64, currency string, amount float64, ref string, remark string, typ string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return core.DB.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet
		if err := tx.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error; err != nil {
			return err
		}
		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}
		result := tx.Model(&wallet).Where("version = ?", wallet.Version).
			Updates(map[string]interface{}{
				"balance": gorm.Expr("balance - ?", amount),
				"version": gorm.Expr("version + 1"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("concurrent update conflict")
		}
		if err := tx.Where("id = ?", wallet.ID).First(&wallet).Error; err == nil {
			ledgerType := typ
			if ledgerType == "" {
				ledgerType = "debit"
			}
			if err := createLedger(tx, userID, currency, ledgerType, -amount, wallet.Balance-amount, wallet.Frozen, ref, remark); err != nil {
				return err
			}
		}
		return nil
	})
}
