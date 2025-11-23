package repository

import (
	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
)

type SettlementLogRepository struct{}

func NewSettlementLogRepository() *SettlementLogRepository {
	return &SettlementLogRepository{}
}

func (r *SettlementLogRepository) Create(log *model.SettlementLog) error {
	return core.DB.Create(log).Error
}
