package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type loginAttempt struct {
	count     int
	resetTime time.Time
}

type LoginLimiter struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempt
	limit    int
	window   time.Duration
}

var loginLimiter = &LoginLimiter{
	attempts: make(map[string]*loginAttempt),
	limit:    5,
	window:   time.Minute * 15,
}

func (r *LoginLimiter) isBlocked(email string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	attempt, exists := r.attempts[email]

	if !exists || now.After(attempt.resetTime) {
		r.attempts[email] = &loginAttempt{
			count:     1,
			resetTime: now.Add(r.window),
		}
		return false
	}

	if attempt.count >= r.limit {
		return true
	}

	attempt.count++
	return false
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() != "/api/v1/auth/login" {
			c.Next()
			return
		}

		var req struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.Next()
			return
		}

		if loginLimiter.isBlocked(req.Email) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many login attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
