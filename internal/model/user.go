package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID         string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255)" json:"-"` // 不返回给前端
	Nickname     string         `gorm:"type:varchar(50)" json:"nickname"`
	Role         string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	RiskLevel    string         `gorm:"type:varchar(20);default:'normal'" json:"risk_level"`
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
