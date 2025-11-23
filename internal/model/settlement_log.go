package model

import "time"

type SettlementLog struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID     int64     `gorm:"index" json:"order_id"`
	OrderNo     string    `gorm:"type:varchar(64);index" json:"order_no"`
	UserID      int64     `gorm:"index" json:"user_id"`
	AssetSymbol string    `gorm:"type:varchar(20)" json:"asset_symbol"`
	Direction   string    `gorm:"type:varchar(10)" json:"direction"`
	Amount      float64   `gorm:"type:decimal(18,8)" json:"amount"`
	PayoutRate  float64   `gorm:"type:decimal(5,4)" json:"payout_rate"`
	OpenPrice   float64   `gorm:"type:decimal(18,8)" json:"open_price"`
	ClosePrice  float64   `gorm:"type:decimal(18,8)" json:"close_price"`
	Result      string    `gorm:"type:varchar(10)" json:"result"` // won, lost, draw
	PnL         float64   `gorm:"type:decimal(18,8)" json:"pnl"`
	SettledAt   time.Time `json:"settled_at"`
	CreatedAt   time.Time `json:"created_at"`
}
