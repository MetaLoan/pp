package repository

import (
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
