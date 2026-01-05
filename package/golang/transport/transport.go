package transport

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	RetryableFunc  func(*http.Request, *http.Response, error) bool
	OnRetryFunc    func(retry int, resp *http.Response, err error)
}

type Transport struct {
	Transport   http.RoundTripper
	RetryConfig *RetryConfig
	Debug       bool
	Timeout     time.Duration
	TLSConfig   *tls.Config
}

func NewTransport(opts ...Option) *Transport {
	t := &Transport{
		Transport:   http.DefaultTransport,
		RetryConfig: DefaultRetryConfig(),
		Timeout:     30 * time.Second,
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.TLSConfig != nil {
		if transport, ok := t.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = t.TLSConfig
		}
	}

	return t
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     5 * time.Second,
		RetryableFunc:  DefaultRetryableFunc,
	}
}

func DefaultRetryableFunc(req *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			return netErr.Timeout() || netErr.Temporary()
		}
		return false
	}

	if resp == nil {
		return false
	}

	status := resp.StatusCode
	return status >= 500 || status == http.StatusTooManyRequests || status == http.StatusRequestTimeout
}

type Option func(*Transport)

func WithTransport(transport http.RoundTripper) Option {
	return func(t *Transport) {
		t.Transport = transport
	}
}

func WithRetryConfig(config *RetryConfig) Option {
	return func(t *Transport) {
		t.RetryConfig = config
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(t *Transport) {
		t.Timeout = timeout
	}
}

func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(t *Transport) {
		t.TLSConfig = tlsConfig
	}
}

func WithDebug(debug bool) Option {
	return func(t *Transport) {
		t.Debug = debug
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if t.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt <= t.RetryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := calculateBackoff(attempt-1, t.RetryConfig.InitialBackoff, t.RetryConfig.MaxBackoff)

			if t.Debug {
				t.logRetry(attempt, lastResp, lastErr, backoff)
			}

			if t.RetryConfig.OnRetryFunc != nil {
				t.RetryConfig.OnRetryFunc(attempt, lastResp, lastErr)
			}

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := t.Transport.RoundTrip(req)
		if err != nil {
			lastErr = err
			lastResp = nil

			if t.RetryConfig.RetryableFunc == nil || !t.RetryConfig.RetryableFunc(req, nil, err) {
				return nil, err
			}
			continue
		}

		if t.RetryConfig.RetryableFunc == nil || !t.RetryConfig.RetryableFunc(req, resp, nil) {
			return resp, nil
		}

		lastResp = resp
		lastErr = nil

		if resp.Body != nil {
			resp.Body.Close()
		}
	}

	return lastResp, lastErr
}

func calculateBackoff(attempt int, initialBackoff, maxBackoff time.Duration) time.Duration {
	backoff := time.Duration(float64(initialBackoff) * float64(attempt+1))
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}

func (t *Transport) logRetry(attempt int, resp *http.Response, err error, backoff time.Duration) {
	if err != nil {
		return
	}

	status := "unknown"
	if resp != nil {
		status = string(resp.StatusCode)
	}
}
