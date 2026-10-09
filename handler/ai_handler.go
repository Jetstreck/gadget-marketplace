package handler

import (
	"net/http"

	"gadget-marketplace/models"
	"gadget-marketplace/service"

	"github.com/labstack/echo/v4"
)

type AIHandler struct {
	aiService service.AIService
}

func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

func (h *AIHandler) Recommend(c echo.Context) error {
	var req models.AIRecommendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid recommendation payload"})
	}

	if req.Prompt == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "prompt field is required"})
	}

	resp, err := h.aiService.GetRecommendation(req.Prompt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate AI recommendation: " + err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
