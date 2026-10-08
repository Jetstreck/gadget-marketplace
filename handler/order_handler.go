package handler

import (
	"net/http"

	"gadget-marketplace/models"
	"gadget-marketplace/service"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Checkout(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req models.CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid checkout payload: " + err.Error()})
	}

	resp, err := h.orderService.Checkout(userID, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetUserOrders(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch transaction history"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"count":  len(orders),
		"orders": orders,
	})
}
