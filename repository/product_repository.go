package repository

import (
	"gadget-marketplace/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Upsert(product *models.Product) error
	FindAll(category string) ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	UpdateStock(id uint, newStock int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Upsert(product *models.Product) error {
	var existing models.Product
	err := r.db.Where("external_id = ?", product.ExternalID).First(&existing).Error
	if err == nil {
		product.ID = existing.ID
		return r.db.Save(product).Error
	}
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll(category string) ([]models.Product, error) {
	var products []models.Product
	query := r.db.Order("id desc")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateStock(id uint, newStock int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", newStock).Error
}
