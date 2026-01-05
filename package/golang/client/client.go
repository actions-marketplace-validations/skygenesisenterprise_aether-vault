package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/config"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
	"github.com/skygenesisenterprise/aether-vault/package/golang/transport"
)

type Client struct {
	httpClient *http.Client
	config     *config.Config
	baseURL    string
}

type RequestOption func(*http.Request)

func NewClient(cfg *config.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	httpClient, err := cfg.ToHTTPClient()
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to create HTTP client")
	}

	if cfg.RetryCount > 0 {
		retryConfig := &transport.RetryConfig{
			MaxRetries:     cfg.RetryCount,
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     5 * time.Second,
		}

		transportLayer := transport.NewTransport(
			transport.WithRetryConfig(retryConfig),
			transport.WithTimeout(cfg.Timeout),
			transport.WithDebug(cfg.Debug),
		)

		httpClient.Transport = transportLayer
	}

	return &Client{
		httpClient: httpClient,
		config:     cfg,
		baseURL:    cfg.Endpoint,
	}, nil
}

func (c *Client) Do(ctx context.Context, method, endpoint string, body interface{}, opts ...RequestOption) (*Response, error) {
	req, err := c.buildRequest(ctx, method, endpoint, body)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInvalidRequest, "failed to build request")
	}

	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeUnavailable, "HTTP request failed")
	}
	defer resp.Body.Close()

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to read response body")
	}
	response.Body = respBody

	if resp.StatusCode >= 400 {
		var vaultErr errors.VaultError
		if err := json.Unmarshal(respBody, &vaultErr); err != nil {
			return response, errors.NewError(errors.ErrorCodeFromStatus(resp.StatusCode),
				fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(respBody)))
		}
		vaultErr.Status = resp.StatusCode
		return response, &vaultErr
	}

	return response, nil
}

func (c *Client) Get(ctx context.Context, endpoint string, opts ...RequestOption) (*Response, error) {
	return c.Do(ctx, http.MethodGet, endpoint, nil, opts...)
}

func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.Do(ctx, http.MethodPost, endpoint, body, opts...)
}

func (c *Client) Put(ctx context.Context, endpoint string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.Do(ctx, http.MethodPut, endpoint, body, opts...)
}

func (c *Client) Patch(ctx context.Context, endpoint string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.Do(ctx, http.MethodPatch, endpoint, body, opts...)
}

func (c *Client) Delete(ctx context.Context, endpoint string, opts ...RequestOption) (*Response, error) {
	return c.Do(ctx, http.MethodDelete, endpoint, nil, opts...)
}

func (c *Client) buildRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, errors.WrapError(err, errors.ErrCodeInvalidRequest, "failed to marshal request body")
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInvalidRequest, "failed to create HTTP request")
	}

	if c.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.Token)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.config.UserAgent)

	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func (r *Response) Decode(v interface{}) error {
	if len(r.Body) == 0 {
		return nil
	}
	return json.Unmarshal(r.Body, v)
}

func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

func WithQueryParam(req *http.Request, key, value string) {
	if req.URL == nil {
		return
	}
	q := req.URL.Query()
	q.Add(key, value)
	req.URL.RawQuery = q.Encode()
}
