package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
)

type NetworkMiddleware struct {
	rateLimiter map[string]*RateLimiter
	mutex       sync.RWMutex
	config      *NetworkConfig
}

type NetworkConfig struct {
	MaxRequestsPerMinute int
	MaxConcurrent        int
	TimeoutSeconds       int
	BlacklistedIPs       []string
	WhitelistedIPs       []string
	AllowedProtocols     []model.ProtocolType
}

type RateLimiter struct {
	requests  []time.Time
	mutex     sync.Mutex
	maxPerMin int
}

type NetworkMetrics struct {
	TotalRequests    int64            `json:"total_requests"`
	ProtocolRequests map[string]int64 `json:"protocol_requests"`
	ErrorRate        float64          `json:"error_rate"`
	AverageLatency   int64            `json:"average_latency_ms"`
	LastActivity     time.Time        `json:"last_activity"`
}

func NewNetworkMiddleware(config *NetworkConfig) *NetworkMiddleware {
	if config == nil {
		config = &NetworkConfig{
			MaxRequestsPerMinute: 100,
			MaxConcurrent:        10,
			TimeoutSeconds:       30,
			AllowedProtocols: []model.ProtocolType{
				model.ProtocolHTTP,
				model.ProtocolHTTPS,
				model.ProtocolSSH,
				model.ProtocolFTP,
				model.ProtocolSFTP,
				model.ProtocolWebDAV,
				model.ProtocolSMB,
				model.ProtocolNFS,
				model.ProtocolRsync,
				model.ProtocolGit,
				model.ProtocolCustom,
			},
		}
	}

	return &NetworkMiddleware{
		rateLimiter: make(map[string]*RateLimiter),
		config:      config,
	}
}

func (m *NetworkMiddleware) ValidateProtocol() gin.HandlerFunc {
	return func(c *gin.Context) {
		protocol := c.GetHeader("X-Protocol-Type")
		if protocol == "" {
			protocol = c.Query("protocol")
		}

		if protocol != "" {
			protocolType := model.ProtocolType(protocol)
			if !m.isProtocolAllowed(protocolType) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":    "protocol not allowed",
					"protocol": protocol,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func (m *NetworkMiddleware) NetworkRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := m.getClientIP(c)

		m.mutex.RLock()
		limiter, exists := m.rateLimiter[clientIP]
		m.mutex.RUnlock()

		if !exists {
			m.mutex.Lock()
			limiter = &RateLimiter{
				maxPerMin: m.config.MaxRequestsPerMinute,
				requests:  make([]time.Time, 0),
			}
			m.rateLimiter[clientIP] = limiter
			m.mutex.Unlock()
		}

		limiter.mutex.Lock()
		now := time.Now()

		validRequests := make([]time.Time, 0)
		for _, req := range limiter.requests {
			if now.Sub(req) < time.Minute {
				validRequests = append(validRequests, req)
			}
		}

		if len(validRequests) >= m.config.MaxRequestsPerMinute {
			limiter.mutex.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":  "rate limit exceeded",
				"limit":  m.config.MaxRequestsPerMinute,
				"window": "1 minute",
			})
			c.Abort()
			return
		}

		limiter.requests = append(validRequests, now)
		limiter.mutex.Unlock()

		c.Next()
	}
}

func (m *NetworkMiddleware) IPWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(m.config.WhitelistedIPs) == 0 {
			c.Next()
			return
		}

		clientIP := m.getClientIP(c)
		if !m.isIPAllowed(clientIP) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "IP address not allowed",
				"ip":    clientIP,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *NetworkMiddleware) NetworkTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.config.TimeoutSeconds > 0 {
			c.Request = c.Request.WithContext(
				c.Request.Context(),
			)
		}
		c.Next()
	}
}

func (m *NetworkMiddleware) NetworkLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := m.getClientIP(c)
		userAgent := c.Request.UserAgent()
		protocol := c.GetHeader("X-Protocol-Type")

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if strings.HasPrefix(path, "/api/v1/network") {
			logData := gin.H{
				"timestamp":  start,
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"latency_ms": latency.Milliseconds(),
				"client_ip":  clientIP,
				"user_agent": userAgent,
				"protocol":   protocol,
				"request_id": c.GetString("request_id"),
			}

			if len(c.Errors) > 0 {
				logData["errors"] = c.Errors.String()
			}

			if statusCode >= 400 {
				c.Error(fmt.Errorf("network request failed"))
			}
		}
	}
}

func (m *NetworkMiddleware) ProtocolSecurity() gin.HandlerFunc {
	return func(c *gin.Context) {
		protocol := c.GetHeader("X-Protocol-Type")
		if protocol == "" {
			protocol = c.Query("protocol")
		}

		if protocol != "" {
			switch model.ProtocolType(protocol) {
			case model.ProtocolHTTP:
				if c.Request.TLS != nil && c.Request.URL.Scheme == "https" {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "HTTP protocol requested over HTTPS connection",
					})
					c.Abort()
					return
				}
			case model.ProtocolHTTPS:
				if c.Request.TLS == nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "HTTPS protocol requires TLS connection",
					})
					c.Abort()
					return
				}
			case model.ProtocolSSH, model.ProtocolSFTP:
				if c.GetHeader("X-SSH-Key") == "" && c.GetHeader("X-SSH-Password") == "" {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "SSH/SFTP protocol requires authentication credentials",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

func (m *NetworkMiddleware) ConcurrentConnectionLimit() gin.HandlerFunc {
	concurrentConnections := make(map[string]int)
	var mutex sync.Mutex

	return func(c *gin.Context) {
		clientIP := m.getClientIP(c)

		mutex.Lock()
		concurrentConnections[clientIP]++
		current := concurrentConnections[clientIP]
		mutex.Unlock()

		defer func() {
			mutex.Lock()
			concurrentConnections[clientIP]--
			if concurrentConnections[clientIP] <= 0 {
				delete(concurrentConnections, clientIP)
			}
			mutex.Unlock()
		}()

		if current > m.config.MaxConcurrent {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many concurrent connections",
				"limit": m.config.MaxConcurrent,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *NetworkMiddleware) GetMetrics() *NetworkMetrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	metrics := &NetworkMetrics{
		TotalRequests:    0,
		ProtocolRequests: make(map[string]int64),
		ErrorRate:        0.0,
		AverageLatency:   0,
		LastActivity:     time.Time{},
	}

	return metrics
}

func (m *NetworkMiddleware) getClientIP(c *gin.Context) string {
	clientIP := c.GetHeader("X-Forwarded-For")
	if clientIP == "" {
		clientIP = c.GetHeader("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	if strings.Contains(clientIP, ",") {
		ips := strings.Split(clientIP, ",")
		clientIP = strings.TrimSpace(ips[0])
	}

	return clientIP
}

func (m *NetworkMiddleware) isIPAllowed(ip string) bool {
	if len(m.config.WhitelistedIPs) == 0 {
		return true
	}

	for _, allowedIP := range m.config.WhitelistedIPs {
		if allowedIP == ip {
			return true
		}
		if strings.Contains(allowedIP, "/") {
			_, ipNet, err := net.ParseCIDR(allowedIP)
			if err == nil && ipNet.Contains(net.ParseIP(ip)) {
				return true
			}
		}
	}

	return false
}

func (m *NetworkMiddleware) isIPBlocked(ip string) bool {
	for _, blockedIP := range m.config.BlacklistedIPs {
		if blockedIP == ip {
			return true
		}
		if strings.Contains(blockedIP, "/") {
			_, ipNet, err := net.ParseCIDR(blockedIP)
			if err == nil && ipNet.Contains(net.ParseIP(ip)) {
				return true
			}
		}
	}

	return false
}

func (m *NetworkMiddleware) isProtocolAllowed(protocol model.ProtocolType) bool {
	for _, allowed := range m.config.AllowedProtocols {
		if allowed == protocol {
			return true
		}
	}

	return false
}
