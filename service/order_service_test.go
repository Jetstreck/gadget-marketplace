package service_test

import (
	"testing"

	"gadget-marketplace/models"
	"gadget-marketplace/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateTransaction(userID uint, productID uint, quantity int) (*models.Order, *models.User, error) {
	args := m.Called(userID, productID, quantity)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*models.Order), args.Get(1).(*models.User), args.Error(2)
}

func (m *MockOrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Order), args.Error(1)
}

func TestOrderService_Checkout_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	svc := service.NewOrderService(mockRepo, nil)

	req := &models.CheckoutRequest{
		ProductID: 1,
		Quantity:  1,
	}

	expectedOrder := &models.Order{
		ID:         1,
		UserID:     10,
		ProductID:  1,
		Quantity:   1,
		TotalPrice: 15000000,
		Status:     "SUCCESS",
	}

	expectedUser := &models.User{
		ID:      10,
		Deposit: 5000000,
	}

	mockRepo.On("CreateTransaction", uint(10), uint(1), 1).Return(expectedOrder, expectedUser, nil)

	resp, err := svc.Checkout(10, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Gadget purchase checkout successful", resp.Message)
	assert.Equal(t, 5000000.0, resp.NewDeposit)
	mockRepo.AssertExpectations(t)
}
