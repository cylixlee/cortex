package auth

import (
	"time"

	"github.com/open-portfolios/cortex/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	secret      string
	expiryHours int
}

func NewTokenService(secret string, expiryHours int) *TokenService {
	return &TokenService{
		secret:      secret,
		expiryHours: expiryHours,
	}
}

func (s *TokenService) GenerateToken(userID uuid.UUID, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(s.expiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, _ := claims["user_id"].(string)
		userID, _ := uuid.Parse(userIDStr)
		email, _ := claims["email"].(string)
		roleStr, _ := claims["role"].(string)

		return &Claims{
			UserID: userID,
			Email:  email,
			Role:   models.UserRole(roleStr),
		}, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
