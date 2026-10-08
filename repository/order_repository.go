package repository

import (
	"errors"
	"fmt"

	"gadget-marketplace/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateTransaction(userID uint, productID uint, quantity int) (*models.Order, *models.User, error)
	FindByUserID(userID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateTransaction(userID uint, productID uint, quantity int) (*models.Order, *models.User, error) {
	var order models.Order
	var user models.User
	var product models.Product

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("user not found")
		}

		if err := tx.First(&product, productID).Error; err != nil {
			return errors.New("product not found")
		}

		if product.Stock < quantity {
			return fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, quantity)
		}

		totalPrice := product.Price * float64(quantity)
		if user.Deposit < totalPrice {
			return fmt.Errorf("insufficient deposit balance: your balance is Rp %.2f, total required is Rp %.2f", user.Deposit, totalPrice)
		}

		user.Deposit -= totalPrice
		if err := tx.Save(&user).Error; err != nil {
			return fmt.Errorf("failed to update user deposit balance: %w", err)
		}

		product.Stock -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		order = models.Order{
			UserID:     userID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: totalPrice,
			Status:     "SUCCESS",
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		order.Product = product
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &order, &user, nil
}

func (r *orderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Product").Preload("User").Where("user_id = ?", userID).Order("id desc").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
