package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoginMiddleware provides authentication validation for login endpoints
type LoginMiddleware struct {
	allowedTokens []string
}

// NewLoginMiddleware creates a new login middleware
func NewLoginMiddleware() *LoginMiddleware {
	return &LoginMiddleware{
		allowedTokens: []string{
			"dev-token",
			"root-token",
			"vault-token",
			"test-token",
			"admin-token",
		},
	}
}

// ValidateToken validates the token format and content
func (m *LoginMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			Token string `json:"token"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.Next()
			return
		}

		// Check if token is in allowed list (for development)
		allowed := false
		for _, token := range m.allowedTokens {
			if req.Token == token {
				allowed = true
				break
			}
		}

		if !allowed && req.Token != "" {
			// Basic token validation
			if len(req.Token) < 8 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "Token too short (minimum 8 characters)",
				})
				ctx.Abort()
				return
			}

			if strings.Contains(req.Token, " ") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   "Token cannot contain spaces",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}

// LogLoginAttempt logs login attempts for security monitoring
func (m *LoginMiddleware) LogLoginAttempt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Capture request info before processing
		clientIP := ctx.ClientIP()
		userAgent := ctx.GetHeader("User-Agent")
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		// Store in context for later use
		ctx.Set("client_ip", clientIP)
		ctx.Set("user_agent", userAgent)

		// Log the attempt
		if strings.Contains(path, "/login") {
			// TODO: Implement proper logging
			_ = clientIP
			_ = userAgent
			_ = method
		}

		ctx.Next()
	}
}

// RateLimitLogin applies basic rate limiting to login attempts
func (m *LoginMiddleware) RateLimitLogin() gin.HandlerFunc {
	// Simple in-memory rate limiter
	attempts := make(map[string]int)

	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()

		// Reset count if too many attempts
		if attempts[clientIP] > 10 {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Too many login attempts. Please try again later.",
			})
			ctx.Abort()
			return
		}

		// Only increment for failed attempts
		ctx.Next()

		// Check if response was an error
		if ctx.Writer.Status() == http.StatusUnauthorized || ctx.Writer.Status() == http.StatusBadRequest {
			attempts[clientIP]++
		} else if ctx.Writer.Status() == http.StatusOK {
			// Reset on successful login
			attempts[clientIP] = 0
		}
	}
}

// ValidateContentType ensures JSON content type for login endpoints
func (m *LoginMiddleware) ValidateContentType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contentType := ctx.GetHeader("Content-Type")

		if ctx.Request.Method == "POST" && !strings.Contains(contentType, "application/json") {
			// Allow empty content type for simplicity
			if contentType != "" {
				ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
					"success": false,
					"error":   "Content-Type must be application/json",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}

// SecurityHeaders adds security headers for login endpoints
func (m *LoginMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Add security headers
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		ctx.Next()
	}
}
