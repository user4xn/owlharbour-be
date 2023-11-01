package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requestsPerSecond int
	mu                sync.Mutex
	lastAccess        map[string]time.Time
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		lastAccess:        make(map[string]time.Time),
	}
}

func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		r.mu.Lock()
		defer r.mu.Unlock()

		if lastAccess, ok := r.lastAccess[clientIP]; ok && time.Since(lastAccess).Seconds() < 1.0/float64(r.requestsPerSecond) {
			c.JSON(429, gin.H{"message": "Too Many Requests"})
			c.Abort()
			return
		}

		r.lastAccess[clientIP] = time.Now()
		c.Next()
	}
}
