package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/audit"
	"github.com/skygenesisenterprise/aether-vault/package/golang/auth"
	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/config"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type Middleware func(client.RequestOption) client.RequestOption

type AuthMiddleware struct {
	authClient  *auth.AuthClient
	token       string
	autoRefresh bool
}

func NewAuthMiddleware(authClient *auth.AuthClient, token string, autoRefresh bool) *AuthMiddleware {
	return &AuthMiddleware{
		authClient:  authClient,
		token:       token,
		autoRefresh: autoRefresh,
	}
}

func (a *AuthMiddleware) Apply(next client.RequestOption) client.RequestOption {
	return func(req *http.Request) {
		if a.token != "" {
			req.Header.Set("Authorization", "Bearer "+a.token)
		}

		if next != nil {
			next(req)
		}
	}
}

func (a *AuthMiddleware) RefreshToken(ctx context.Context) error {
	if !a.autoRefresh {
		return nil
	}

	resp, err := a.authClient.Refresh(ctx, a.token)
	if err != nil {
		return err
	}

	a.token = resp.Token
	return nil
}

type AuditMiddleware struct {
	auditClient *audit.AuditClient
	identityID  *string
}

func NewAuditMiddleware(auditClient *audit.AuditClient, identityID *string) *AuditMiddleware {
	return &AuditMiddleware{
		auditClient: auditClient,
		identityID:  identityID,
	}
}

func (a *AuditMiddleware) LogAction(ctx context.Context, action, resource, resourceID, status, message string) {
	event := &audit.AuditEvent{
		IdentityID: a.identityID,
		Action:     action,
		Resource:   resource,
		ResourceID: &resourceID,
		Status:     status,
		Message:    message,
		Timestamp:  time.Now(),
	}

	go func() {
		_, err := a.auditClient.Search(ctx, &audit.SearchAuditRequest{
			Action:     action,
			Resource:   resource,
			ResourceID: a.identityID,
			Limit:      1,
		})
		if err != nil {
			return
		}
	}()
}

type RetryMiddleware struct {
	maxRetries     int
	initialBackoff time.Duration
	maxBackoff     time.Duration
	retryableFunc  func(*http.Request, *http.Response, error) bool
	onRetryFunc    func(retry int, resp *http.Response, err error)
}

func NewRetryMiddleware(maxRetries int) *RetryMiddleware {
	return &RetryMiddleware{
		maxRetries:     maxRetries,
		initialBackoff: 100 * time.Millisecond,
		maxBackoff:     5 * time.Second,
		retryableFunc:  defaultRetryableFunc,
	}
}

func (r *RetryMiddleware) WithBackoff(initial, max time.Duration) *RetryMiddleware {
	r.initialBackoff = initial
	r.maxBackoff = max
	return r
}

func (r *RetryMiddleware) WithRetryableFunc(fn func(*http.Request, *http.Response, error) bool) *RetryMiddleware {
	r.retryableFunc = fn
	return r
}

func (r *RetryMiddleware) WithOnRetryFunc(fn func(retry int, resp *http.Response, err error)) *RetryMiddleware {
	r.onRetryFunc = fn
	return r
}

func (r *RetryMiddleware) Apply(next client.RequestOption) client.RequestOption {
	return func(req *http.Request) {
		if next != nil {
			next(req)
		}
	}
}

type MetricsMiddleware struct {
	metrics map[string]interface{}
}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: make(map[string]interface{}),
	}
}

func (m *MetricsMiddleware) RecordRequest(method, endpoint string, statusCode int, duration time.Duration) {
	key := method + ":" + endpoint

	if _, exists := m.metrics[key]; !exists {
		m.metrics[key] = map[string]interface{}{
			"count":        0,
			"total_time":   time.Duration(0),
			"status_codes": make(map[int]int),
		}
	}

	methodMetrics := m.metrics[key].(map[string]interface{})
	methodMetrics["count"] = methodMetrics["count"].(int) + 1
	methodMetrics["total_time"] = methodMetrics["total_time"].(time.Duration) + duration

	statusCodes := methodMetrics["status_codes"].(map[int]int)
	statusCodes[statusCode]++
}

func (m *MetricsMiddleware) GetMetrics() map[string]interface{} {
	return m.metrics
}

type CircuitBreakerMiddleware struct {
	maxFailures int
	resetTime   time.Duration
	failures    int
	lastFailure time.Time
	state       CircuitState
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

func NewCircuitBreakerMiddleware(maxFailures int, resetTime time.Duration) *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		maxFailures: maxFailures,
		resetTime:   resetTime,
		state:       CircuitClosed,
	}
}

func (c *CircuitBreakerMiddleware) Apply(next client.RequestOption) client.RequestOption {
	return func(req *http.Request) {
		if c.state == CircuitOpen && time.Since(c.lastFailure) < c.resetTime {
			return
		}

		if c.state == CircuitOpen && time.Since(c.lastFailure) >= c.resetTime {
			c.state = CircuitHalfOpen
			c.failures = 0
		}

		if next != nil {
			next(req)
		}
	}
}

func (c *CircuitBreakerMiddleware) RecordSuccess() {
	c.failures = 0
	c.state = CircuitClosed
}

func (c *CircuitBreakerMiddleware) RecordFailure() {
	c.failures++
	c.lastFailure = time.Now()

	if c.failures >= c.maxFailures {
		c.state = CircuitOpen
	}
}

type RateLimitMiddleware struct {
	requests     int
	window       time.Duration
	allowBurst   bool
	requestTimes []time.Time
}

func NewRateLimitMiddleware(requests int, window time.Duration, allowBurst bool) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		requests:   requests,
		window:     window,
		allowBurst: allowBurst,
	}
}

func (r *RateLimitMiddleware) Apply(next client.RequestOption) client.RequestOption {
	return func(req *http.Request) {
		now := time.Now()

		cutoff := now.Add(-r.window)
		validTimes := make([]time.Time, 0)
		for _, t := range r.requestTimes {
			if t.After(cutoff) {
				validTimes = append(validTimes, t)
			}
		}

		if len(validTimes) >= r.requests && !r.allowBurst {
			return
		}

		r.requestTimes = append(validTimes, now)

		if len(r.requestTimes) > r.requests*2 {
			r.requestTimes = r.requestTimes[len(r.requestTimes)-r.requests:]
		}

		if next != nil {
			next(req)
		}
	}
}

func defaultRetryableFunc(req *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp == nil {
		return false
	}

	status := resp.StatusCode
	return status >= 500 || status == http.StatusTooManyRequests || status == http.StatusRequestTimeout
}

type Chain struct {
	middlewares []Middleware
}

func NewChain(middlewares ...Middleware) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

func (c *Chain) Apply(option client.RequestOption) client.RequestOption {
	result := option
	for _, middleware := range c.middlewares {
		result = middleware(result)
	}
	return result
}
