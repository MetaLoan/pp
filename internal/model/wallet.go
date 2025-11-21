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
