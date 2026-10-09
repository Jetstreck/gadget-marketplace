package models

type AIRecommendRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type AIRecommendResponse struct {
	Prompt         string `json:"prompt"`
	Recommendation string `json:"recommendation"`
}
