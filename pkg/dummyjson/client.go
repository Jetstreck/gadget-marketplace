package dummyjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gadget-marketplace/models"
)

type DummyProduct struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Rating      float64  `json:"rating"`
	Stock       int      `json:"stock"`
	Brand       string   `json:"brand"`
	Category    string   `json:"category"`
	Thumbnail   string   `json:"thumbnail"`
	Images      []string `json:"images"`
}

type DummyResponse struct {
	Products []DummyProduct `json:"products"`
	Total    int            `json:"total"`
	Skip     int            `json:"skip"`
	Limit    int            `json:"limit"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) FetchCategoryProducts(category string) ([]models.Product, error) {
	url := fmt.Sprintf("%s/products/category/%s", c.baseURL, category)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from DummyJSON: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DummyJSON API returned status code %d", resp.StatusCode)
	}

	var apiResp DummyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode DummyJSON response: %w", err)
	}

	var products []models.Product
	for _, p := range apiResp.Products {
		// Convert USD price to IDR (approx 15,000 IDR per USD for realistic pricing)
		priceIDR := p.Price * 15000
		if priceIDR == 0 {
			priceIDR = 5000000
		}

		stockVal := p.Stock
		if stockVal <= 0 {
			stockVal = 10
		}

		products = append(products, models.Product{
			ExternalID:        p.ID,
			Name:              p.Title,
			Category:          p.Category,
			Price:             priceIDR,
			RentalCosts:       priceIDR,
			Stock:             stockVal,
			StockAvailability: stockVal,
			Brand:             p.Brand,
			Description:       p.Description,
			Thumbnail:         p.Thumbnail,
		})
	}

	return products, nil
}
