package auth

import (
	"github.com/cylixlee/cortex/internal/models"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID       `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
}
