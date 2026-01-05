package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	clients map[string]*ClientLimiter
	mutex   sync.RWMutex
	rate    int
	window  time.Duration
}

type ClientLimiter struct {
	requests  int
	lastReset time.Time
	mutex     sync.Mutex
}

func NewRateLimitMiddleware(requests int, window time.Duration) *RateLimitMiddleware {
	limiter := &RateLimitMiddleware{
		clients: make(map[string]*ClientLimiter),
		rate:    requests,
		window:  window,
	}

	go limiter.cleanup()

	return limiter
}

func (m *RateLimitMiddleware) Limit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()

		m.mutex.RLock()
		limiter, exists := m.clients[clientIP]
		m.mutex.RUnlock()

		if !exists {
			m.mutex.Lock()
			limiter = &ClientLimiter{
				requests:  0,
				lastReset: time.Now(),
			}
			m.clients[clientIP] = limiter
			m.mutex.Unlock()
		}

		limiter.mutex.Lock()

		if time.Since(limiter.lastReset) >= m.window {
			limiter.requests = 0
			limiter.lastReset = time.Now()
		}

		if limiter.requests >= m.rate {
			limiter.mutex.Unlock()

			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "VAULT_RATE_LIMIT_EXCEEDED",
					"message": "Rate limit exceeded",
				},
			})
			ctx.Abort()
			return
		}

		limiter.requests++
		limiter.mutex.Unlock()

		ctx.Next()
	}
}

func (m *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(m.window)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		for ip, limiter := range m.clients {
			limiter.mutex.Lock()
			if time.Since(limiter.lastReset) >= m.window*2 {
				delete(m.clients, ip)
			}
			limiter.mutex.Unlock()
		}
		m.mutex.Unlock()
	}
}
