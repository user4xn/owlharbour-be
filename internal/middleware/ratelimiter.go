package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requestsPerSecond int
	windowSize        time.Duration
	mu                sync.Mutex
	requestQueue      map[string][]time.Time
}

func NewRateLimiter(requestsPerSecond int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		windowSize:        windowSize,
		requestQueue:      make(map[string][]time.Time),
	}
}

func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		r.mu.Lock()
		defer r.mu.Unlock()

		now := time.Now()
		queue, ok := r.requestQueue[clientIP]
		if !ok {
			r.requestQueue[clientIP] = []time.Time{now}
		} else {
			// Remove timestamps that are older than the sliding window
			for len(queue) > 0 && now.Sub(queue[0]) > r.windowSize {
				queue = queue[1:]
			}
			queue = append(queue, now)
			r.requestQueue[clientIP] = queue
		}

		if len(queue) > r.requestsPerSecond {
			c.JSON(429, gin.H{"message": "Too Many Requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
