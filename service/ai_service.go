package service

import (
	"errors"

	"gadget-marketplace/config"
	"gadget-marketplace/models"
	"gadget-marketplace/pkg/gemini"
)

type AIService interface {
	GetRecommendation(prompt string) (*models.AIRecommendResponse, error)
}

type aiService struct {
	geminiClient *gemini.GeminiClient
}

func NewAIService(cfg *config.Config) AIService {
	return &aiService{
		geminiClient: gemini.NewGeminiClient(cfg.GeminiAPIKey),
	}
}

func (s *aiService) GetRecommendation(prompt string) (*models.AIRecommendResponse, error) {
	if prompt == "" {
		return nil, errors.New("prompt is required")
	}

	recommendation, err := s.geminiClient.GenerateGadgetRecommendation(prompt)
	if err != nil {
		return nil, err
	}

	return &models.AIRecommendResponse{
		Prompt:         prompt,
		Recommendation: recommendation,
	}, nil
}
