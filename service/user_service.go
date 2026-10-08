package service

import (
	"errors"

	"gadget-marketplace/config"
	"gadget-marketplace/models"
	"gadget-marketplace/pkg/utils"
	"gadget-marketplace/repository"
)

type UserService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	GetProfile(userID uint) (*models.User, error)
	TopUp(userID uint, amount float64) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *userService) Register(req *models.RegisterRequest) (*models.User, error) {
	existingEmail, _ := s.userRepo.FindByEmail(req.Email)
	if existingEmail != nil {
		return nil, errors.New("email is already registered")
	}

	existingUsername, _ := s.userRepo.FindByUsername(req.Username)
	if existingUsername != nil {
		return nil, errors.New("username is already taken")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Deposit:  0,
		Role:     "customer",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpirationHours)
	if err != nil {
		return nil, errors.New("failed to generate auth token")
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) GetProfile(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) TopUp(userID uint, amount float64) (*models.User, error) {
	if amount <= 0 {
		return nil, errors.New("top-up amount must be greater than zero")
	}
	return s.userRepo.UpdateDeposit(userID, amount)
}
