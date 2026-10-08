package service

import (
	"errors"

	"gadget-marketplace/models"
	"gadget-marketplace/repository"
)

type OrderService interface {
	Checkout(userID uint, req *models.CheckoutRequest) (*models.CheckoutResponse, error)
	GetUserOrders(userID uint) ([]models.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) Checkout(userID uint, req *models.CheckoutRequest) (*models.CheckoutResponse, error) {
	if req.ProductID == 0 {
		return nil, errors.New("product_id is required")
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	order, updatedUser, err := s.orderRepo.CreateTransaction(userID, req.ProductID, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &models.CheckoutResponse{
		Message:    "Gadget purchase checkout successful",
		Order:      *order,
		NewDeposit: updatedUser.Deposit,
	}, nil
}

func (s *orderService) GetUserOrders(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}
