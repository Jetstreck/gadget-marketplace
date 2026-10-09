package service

import (
	"errors"
	"log"

	"gadget-marketplace/models"
	"gadget-marketplace/pkg/email"
	"gadget-marketplace/repository"
)

type OrderService interface {
	Checkout(userID uint, req *models.CheckoutRequest) (*models.CheckoutResponse, error)
	GetUserOrders(userID uint) ([]models.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	emailSvc email.EmailService
}

func NewOrderService(orderRepo repository.OrderRepository, emailSvc email.EmailService) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		emailSvc: emailSvc,
	}
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

	if s.emailSvc != nil {
		go func() {
			if err := s.emailSvc.SendOrderInvoiceEmail(updatedUser.Email, updatedUser.Username, order.Product.Name, order.Quantity, order.TotalPrice, updatedUser.Deposit); err != nil {
				log.Printf("Non-blocking notice: invoice email failed (%v)", err)
			}
		}()
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
