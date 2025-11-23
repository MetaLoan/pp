package repository

import (
	"time"

	"github.com/MetaLoan/pp/internal/core"
	"github.com/MetaLoan/pp/internal/model"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	return core.DB.Create(order).Error
}

func (r *OrderRepository) GetOrderByID(id int64) (*model.Order, error) {
	var order model.Order
	err := core.DB.First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) UpdateOrder(order *model.Order) error {
	return core.DB.Save(order).Error
}

func (r *OrderRepository) GetOrdersByUserIDAndStatus(userID int64, status string) ([]model.Order, error) {
	var orders []model.Order
	err := core.DB.Where("user_id = ? AND status = ?", userID, status).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) CountOpenOrders(userID int64) (int64, error) {
	var count int64
	err := core.DB.Model(&model.Order{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Count(&count).Error
	return count, err
}

func (r *OrderRepository) GetOrdersByUserID(userID int64, limit, offset int, status string) ([]model.Order, error) {
	var orders []model.Order
	tx := core.DB.Where("user_id = ?", userID)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if offset > 0 {
		tx = tx.Offset(offset)
	}
	err := tx.Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) SumPnLSince(userID int64, from time.Time) (float64, error) {
	var total float64
	err := core.DB.Model(&model.Order{}).
		Select("COALESCE(SUM(pn_l),0)").
		Where("user_id = ? AND close_time >= ? AND status IN ?", userID, from, []string{"won", "lost", "draw"}).
		Scan(&total).Error
	return total, err
}

// GetDueOrders returns active orders whose close_time has passed.
func (r *OrderRepository) GetDueOrders(limit int) ([]model.Order, error) {
	var orders []model.Order
	tx := core.DB.Where("status = ? AND close_time <= ?", "active", time.Now())
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	err := tx.Order("close_time asc").Find(&orders).Error
	return orders, err
}
