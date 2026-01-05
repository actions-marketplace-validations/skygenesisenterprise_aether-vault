package auth

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TOTP     string `json:"totp,omitempty"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
}

type RefreshRequest struct {
	Token string `json:"token"`
}

type RefreshResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type RevokeRequest struct {
	Token string `json:"token"`
}

type RevokeResponse struct {
	Revoked bool `json:"revoked"`
}

type AuthClient struct {
	client *client.Client
}

func NewAuthClient(client *client.Client) *AuthClient {
	return &AuthClient{
		client: client,
	}
}

func (a *AuthClient) Login(ctx context.Context, username, password, totp string) (*AuthResponse, error) {
	req := &AuthRequest{
		Username: username,
		Password: password,
		TOTP:     totp,
	}

	resp, err := a.client.Post(ctx, "/api/v1/auth/login", req)
	if err != nil {
		return nil, err
	}

	var authResp AuthResponse
	if err := resp.Decode(&authResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode auth response")
	}

	return &authResp, nil
}

func (a *AuthClient) Refresh(ctx context.Context, token string) (*RefreshResponse, error) {
	req := &RefreshRequest{
		Token: token,
	}

	resp, err := a.client.Post(ctx, "/api/v1/auth/refresh", req)
	if err != nil {
		return nil, err
	}

	var refreshResp RefreshResponse
	if err := resp.Decode(&refreshResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode refresh response")
	}

	return &refreshResp, nil
}

func (a *AuthClient) Verify(ctx context.Context, token string) (*VerifyResponse, error) {
	req := &VerifyRequest{
		Token: token,
	}

	resp, err := a.client.Post(ctx, "/api/v1/auth/verify", req)
	if err != nil {
		return nil, err
	}

	var verifyResp VerifyResponse
	if err := resp.Decode(&verifyResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode verify response")
	}

	return &verifyResp, nil
}

func (a *AuthClient) Revoke(ctx context.Context, token string) (*RevokeResponse, error) {
	req := &RevokeRequest{
		Token: token,
	}

	resp, err := a.client.Post(ctx, "/api/v1/auth/revoke", req)
	if err != nil {
		return nil, err
	}

	var revokeResp RevokeResponse
	if err := resp.Decode(&revokeResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode revoke response")
	}

	return &revokeResp, nil
}

func (a *AuthClient) Logout(ctx context.Context) error {
	_, err := a.client.Post(ctx, "/api/v1/auth/logout", nil)
	return err
}
