package ratelimit

import "github.com/gin-gonic/gin"

type Limiter interface {
	Limit(key string) gin.HandlerFunc
}
