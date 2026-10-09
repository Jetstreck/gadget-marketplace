package service_test

import (
	"testing"

	"gadget-marketplace/config"
	"gadget-marketplace/models"
	"gadget-marketplace/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	if user.ID == 0 {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateDeposit(userID uint, amount float64) (*models.User, error) {
	args := m.Called(userID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpirationHours: 24}
	svc := service.NewUserService(mockRepo, nil, cfg)

	req := &models.RegisterRequest{
		Email:    "john@example.com",
		Username: "john_doe",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", req.Email).Return(nil, nil)
	mockRepo.On("FindByUsername", req.Username).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := svc.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "john_doe", user.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserService_TopUp_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "testsecret", JWTExpirationHours: 24}
	svc := service.NewUserService(mockRepo, nil, cfg)

	expectedUser := &models.User{
		ID:      1,
		Email:   "john@example.com",
		Deposit: 500000,
	}

	mockRepo.On("UpdateDeposit", uint(1), 500000.0).Return(expectedUser, nil)

	user, err := svc.TopUp(1, 500000.0)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 500000.0, user.Deposit)
	mockRepo.AssertExpectations(t)
}
