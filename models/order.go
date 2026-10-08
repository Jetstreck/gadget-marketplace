package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID  uint           `gorm:"not null" json:"product_id"`
	Product    Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity   int            `gorm:"not null;default:1" json:"quantity"`
	TotalPrice float64        `gorm:"not null" json:"total_price"`
	Status     string         `gorm:"default:'SUCCESS'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CheckoutRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CheckoutResponse struct {
	Message    string  `json:"message"`
	Order      Order   `json:"order"`
	NewDeposit float64 `json:"new_deposit"`
}
