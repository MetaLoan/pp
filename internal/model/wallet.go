package model

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64          `gorm:"uniqueIndex:idx_user_currency" json:"user_id"`
	Currency  string         `gorm:"type:varchar(10);uniqueIndex:idx_user_currency" json:"currency"`
	Balance   float64        `gorm:"type:decimal(24,8);default:0" json:"balance"`
	Frozen    float64        `gorm:"type:decimal(24,8);default:0" json:"frozen"`
	Version   int64          `gorm:"default:0" json:"version"` // Optimistic locking
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type WalletLedger struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64          `gorm:"index" json:"user_id"`
	Currency  string         `gorm:"index" json:"currency"`
	Type      string         `gorm:"type:varchar(30);index" json:"type"` // freeze, unfreeze, debit, credit, settle_win, settle_lose, settle_draw
	Amount    float64        `gorm:"type:decimal(24,8)" json:"amount"`
	Balance   float64        `gorm:"type:decimal(24,8)" json:"balance"`
	Frozen    float64        `gorm:"type:decimal(24,8)" json:"frozen"`
	Ref       string         `gorm:"type:varchar(64)" json:"ref"` // order_no or any reference
	Remark    string         `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
