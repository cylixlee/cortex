package service

import (
	"errors"

	"github.com/cylixlee/cortex/internal/auth"
	"github.com/cylixlee/cortex/internal/models"
	"github.com/cylixlee/cortex/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
	ErrFailedToHashPassword   = errors.New("failed to hash password")
	ErrFailedToGenerateToken  = errors.New("failed to generate token")
)

type UserService struct {
	userRepo     *repository.UserRepository
	tokenService *auth.TokenService
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiryHours int) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenService: auth.NewTokenService(jwtSecret, jwtExpiryHours),
	}
}

type RegisterOutput struct {
	UserID uuid.UUID
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type GetUserOutput struct {
	ID        uuid.UUID
	Email     string
	Role      models.UserRole
	CreatedAt string
}

func (s *UserService) Register(email, password string) (*RegisterOutput, error) {
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, ErrFailedToHashPassword
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterOutput{UserID: user.ID}, nil
}

func (s *UserService) Login(email, password string) (*LoginOutput, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, ErrFailedToGenerateToken
	}

	refreshToken, err := s.tokenService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, ErrFailedToGenerateToken
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", ErrUserNotFound
	}

	accessToken, err := s.tokenService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", ErrFailedToGenerateToken
	}

	return accessToken, nil
}

func (s *UserService) GetUser(userID uuid.UUID) (*GetUserOutput, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &GetUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *UserService) ListUsers(limit, offset int) ([]GetUserOutput, error) {
	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []GetUserOutput
	for _, user := range users {
		result = append(result, GetUserOutput{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return result, nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}
