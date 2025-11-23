package repository

import (
	"time"

	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
)

type TransferRepository struct{}

func NewTransferRepository() *TransferRepository {
	return &TransferRepository{}
}

func (r *TransferRepository) Create(req *model.TransferRequest) error {
	return core.DB.Create(req).Error
}

func (r *TransferRepository) Update(req *model.TransferRequest) error {
	return core.DB.Save(req).Error
}

func (r *TransferRepository) GetPendingByID(id int64) (*model.TransferRequest, error) {
	var req model.TransferRequest
	err := core.DB.Where("id = ? AND status = ?", id, "pending").First(&req).Error
	return &req, err
}

func (r *TransferRepository) GetByUser(userID int64, limit int) ([]model.TransferRequest, error) {
	var reqs []model.TransferRequest
	tx := core.DB.Where("user_id = ?", userID).Order("created_at desc")
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	err := tx.Find(&reqs).Error
	return reqs, err
}

func (r *TransferRepository) SumApprovedAmountSince(userID int64, typ string, since time.Time) (float64, error) {
	var total float64
	err := core.DB.Model(&model.TransferRequest{}).
		Select("COALESCE(SUM(amount),0)").
		Where("user_id = ? AND type = ? AND status = ? AND reviewed_at >= ?", userID, typ, "approved", since).
		Scan(&total).Error
	return total, err
}

func (r *TransferRepository) LastRequestTime(userID int64, typ string) (time.Time, error) {
	var t time.Time
	err := core.DB.Model(&model.TransferRequest{}).
		Select("COALESCE(MAX(created_at), 'epoch')").
		Where("user_id = ? AND type = ?", userID, typ).
		Scan(&t).Error
	return t, err
}
