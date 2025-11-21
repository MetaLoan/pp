package service

import (
	"errors"

	"github.com/MetaLoan/pp/internal/model"
	"github.com/MetaLoan/pp/internal/repository"
	"github.com/MetaLoan/pp/pkg/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	walletRepo *repository.WalletRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:   repository.NewUserRepository(),
		walletRepo: repository.NewWalletRepository(),
	}
}

func (s *AuthService) Register(email, password, nickname string) (*model.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetUserByEmail(email)
	if existingUser != nil && existingUser.ID != 0 {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UUID:         uuid.New().String(),
		Email:        email,
		PasswordHash: hashedPassword,
		Nickname:     nickname,
		Role:         "user",
		RiskLevel:    "normal",
		Status:       "active",
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// Create initial wallet with demo funds
	// Note: In a real app, we wouldn't give free money, but this is for MVP/Demo
	_, err = s.walletRepo.GetWallet(user.ID, "USDT")
	if err != nil {
		// Log error but don't fail registration? Or fail?
		// For now, let's just log it or ignore, as GetWallet auto-creates it.
		// Actually, GetWallet auto-creates. So calling it is enough.
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.UUID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
