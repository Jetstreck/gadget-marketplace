package service

import (
	"errors"

	"gadget-marketplace/config"
	"gadget-marketplace/models"
	"gadget-marketplace/pkg/gemini"
	"gadget-marketplace/repository"
)

type AIService interface {
	GetRecommendation(prompt string) (*models.AIRecommendResponse, error)
}

type aiService struct {
	productRepo  repository.ProductRepository
	geminiClient *gemini.GeminiClient
}

func NewAIService(productRepo repository.ProductRepository, cfg *config.Config) AIService {
	return &aiService{
		productRepo:  productRepo,
		geminiClient: gemini.NewGeminiClient(cfg.GeminiAPIKey),
	}
}

func (s *aiService) GetRecommendation(prompt string) (*models.AIRecommendResponse, error) {
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}

	// Fetch live product catalog from database
	catalog, _ := s.productRepo.FindAll("")

	recommendation, err := s.geminiClient.GenerateGadgetRecommendation(prompt, catalog)
	if err != nil {
		return nil, err
	}

	return &models.AIRecommendResponse{
		Prompt:         prompt,
		Recommendation: recommendation,
	}, nil
}
