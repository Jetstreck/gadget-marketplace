package handler

import (
	"net/http"

	"gadget-marketplace/models"
	"gadget-marketplace/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload: " + err.Error()})
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email, username, and password are required"})
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email and password are required"})
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   resp.Token,
		"user":    resp.User,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User profile not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (h *UserHandler) TopUp(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req models.TopUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid top-up payload"})
	}

	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Amount must be greater than zero"})
	}

	updatedUser, err := h.userService.TopUp(userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":     "Top-up successful",
		"new_deposit": updatedUser.Deposit,
		"user":        updatedUser,
	})
}
