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

		return nil
	})
}

func (r *WalletRepository) UnfreezeAndSettle(userID int64, currency string, frozenAmount, profitAmount float64) error {
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

		return nil
	})
}
