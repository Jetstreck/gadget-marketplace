package service

import (
	"fmt"
	"log"

	"gadget-marketplace/config"
	"gadget-marketplace/models"
	"gadget-marketplace/pkg/dummyjson"
	"gadget-marketplace/repository"
)

type ProductService interface {
	SyncFromDummyJSON() ([]models.Product, error)
	GetAllProducts(category string) ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
}

type productService struct {
	productRepo     repository.ProductRepository
	dummyJSONClient *dummyjson.Client
	cfg             *config.Config
}

func NewProductService(productRepo repository.ProductRepository, cfg *config.Config) ProductService {
	return &productService{
		productRepo:     productRepo,
		dummyJSONClient: dummyjson.NewClient(cfg.DummyJSONURL),
		cfg:             cfg,
	}
}

func (s *productService) SyncFromDummyJSON() ([]models.Product, error) {
	categories := []string{"smartphones", "laptops"}
	var syncedProducts []models.Product

	for _, cat := range categories {
		prods, err := s.dummyJSONClient.FetchCategoryProducts(cat)
		if err != nil {
			log.Printf("Warning: failed to fetch category %s from DummyJSON: %v", cat, err)
			continue
		}

		for _, p := range prods {
			productCopy := p
			if err := s.productRepo.Upsert(&productCopy); err != nil {
				log.Printf("Warning: failed to upsert product %s: %v", p.Name, err)
			} else {
				syncedProducts = append(syncedProducts, productCopy)
			}
		}
	}

	if len(syncedProducts) == 0 {
		return nil, fmt.Errorf("no products synced from DummyJSON")
	}

	return syncedProducts, nil
}

func (s *productService) GetAllProducts(category string) ([]models.Product, error) {
	products, err := s.productRepo.FindAll(category)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		log.Println("Database product catalog is empty. Auto-syncing from 3rd Party DummyJSON API...")
		_, _ = s.SyncFromDummyJSON()
		products, err = s.productRepo.FindAll(category)
	}

	return products, err
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}
