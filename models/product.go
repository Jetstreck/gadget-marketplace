package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	ExternalID        int            `gorm:"index" json:"external_id,omitempty"`
	Name              string         `gorm:"not null" json:"name"`
	Category          string         `gorm:"not null" json:"category"`
	Price             float64        `gorm:"not null" json:"price"`
	RentalCosts       float64        `gorm:"-" json:"rental_costs"`
	Stock             int            `gorm:"not null;default:0" json:"stock"`
	StockAvailability int            `gorm:"-" json:"stock_availability"`
	Brand             string         `json:"brand"`
	Description       string         `json:"description"`
	Thumbnail         string         `json:"thumbnail"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	p.RentalCosts = p.Price
	p.StockAvailability = p.Stock
	return nil
}

type SyncProductsResponse struct {
	SyncedCount int       `json:"synced_count"`
	Message     string    `json:"message"`
	Products    []Product `json:"products"`
}
