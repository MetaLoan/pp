package repository

import (
	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
)

type WalletLedgerRepository struct{}

func NewWalletLedgerRepository() *WalletLedgerRepository {
	return &WalletLedgerRepository{}
}

func (r *WalletLedgerRepository) Create(entry *model.WalletLedger) error {
	return core.DB.Create(entry).Error
}
