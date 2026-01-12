package middleware

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SNMPMiddleware struct {
	maxRequestsPerMinute int
	requests             map[string][]time.Time
}

type SNMPConfig struct {
	MaxRequestsPerMinute int
	AllowedNetworks      []string
	TimeoutSeconds       int
}

func NewSNMPMiddleware(config *SNMPConfig) *SNMPMiddleware {
	if config == nil {
		config = &SNMPConfig{
			MaxRequestsPerMinute: 30,
			AllowedNetworks:      []string{"127.0.0.1/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
			TimeoutSeconds:       30,
		}
	}

	return &SNMPMiddleware{
		maxRequestsPerMinute: config.MaxRequestsPerMinute,
		requests:             make(map[string][]time.Time),
	}
}

func (m *SNMPMiddleware) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		now := time.Now()

		if m.requests[clientIP] == nil {
			m.requests[clientIP] = []time.Time{}
		}

		var validRequests []time.Time
		for _, reqTime := range m.requests[clientIP] {
			if now.Sub(reqTime) < time.Minute {
				validRequests = append(validRequests, reqTime)
			}
		}
		m.requests[clientIP] = validRequests

		if len(m.requests[clientIP]) >= m.maxRequestsPerMinute {
			ctx.JSON(429, gin.H{
				"error": "Rate limit exceeded for SNMP requests",
				"limit": m.maxRequestsPerMinute,
			})
			ctx.Abort()
			return
		}

		m.requests[clientIP] = append(m.requests[clientIP], now)
		ctx.Next()
	}
}

func (m *SNMPMiddleware) ValidateTarget() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var target struct {
			Target string `json:"target"`
		}

		if err := ctx.ShouldBindJSON(&target); err != nil {
			ctx.Next()
			return
		}

		if target.Target == "" {
			ctx.Next()
			return
		}

		if net.ParseIP(target.Target) == nil {
			ips, err := net.LookupIP(target.Target)
			if err != nil || len(ips) == 0 {
				ctx.JSON(400, gin.H{
					"error": "Invalid target hostname or IP address",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}

func (m *SNMPMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Next()
	}
}

func (m *SNMPMiddleware) LogSNMPRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		statusCode := ctx.Writer.Status()

		if strings.Contains(path, "/snmp") {
			gin.DefaultWriter.Write([]byte(
				fmt.Sprintf("[%s] SNMP %s %s from %s - %d (%v)\n",
					time.Now().Format("2006-01-02 15:04:05"),
					method,
					path,
					clientIP,
					statusCode,
					duration,
				),
			))
		}
	}
}
