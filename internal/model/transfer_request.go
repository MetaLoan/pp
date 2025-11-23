package model

import "time"

type TransferRequest struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"index" json:"user_id"`
	Type       string    `gorm:"type:varchar(20);index" json:"type"` // deposit, withdraw
	Currency   string    `gorm:"type:varchar(10);index" json:"currency"`
	Amount     float64   `gorm:"type:decimal(24,8)" json:"amount"`
	Status     string    `gorm:"type:varchar(20);index" json:"status"` // pending, approved, rejected
	Remark     string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ReviewedAt time.Time `json:"reviewed_at"`
	Reviewer   string    `gorm:"type:varchar(50)" json:"reviewer"`
}
