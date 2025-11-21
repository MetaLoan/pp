package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_no"`
	UserID      int64          `gorm:"index" json:"user_id"`
	AssetSymbol string         `gorm:"type:varchar(20);not null" json:"asset_symbol"`
	Direction   string         `gorm:"type:varchar(10);not null" json:"direction"` // CALL, PUT
	Amount      float64        `gorm:"type:decimal(18,8);not null" json:"amount"`
	PayoutRate  float64        `gorm:"type:decimal(5,4);not null" json:"payout_rate"`
	OpenPrice   float64        `gorm:"type:decimal(18,8);not null" json:"open_price"`
	OpenTime    time.Time      `json:"open_time"`
	ClosePrice  float64        `gorm:"type:decimal(18,8)" json:"close_price"`
	CloseTime   time.Time      `gorm:"index" json:"close_time"`
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, active, won, lost, draw
	PnL         float64        `gorm:"type:decimal(18,8);default:0" json:"pnl"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
