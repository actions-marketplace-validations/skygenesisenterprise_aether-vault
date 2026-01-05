package totp

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type GenerateRequest struct {
	AccountName string `json:"account_name"`
	Issuer      string `json:"issuer,omitempty"`
	Algorithm   string `json:"algorithm,omitempty"`
	Digits      int    `json:"digits,omitempty"`
	Period      int    `json:"period,omitempty"`
}

type GenerateResponse struct {
	Secret      string     `json:"secret"`
	QRCode      string     `json:"qr_code"`
	BackupCodes []string   `json:"backup_codes"`
	AccountName string     `json:"account_name"`
	Issuer      string     `json:"issuer"`
	Algorithm   string     `json:"algorithm"`
	Digits      int        `json:"digits"`
	Period      int        `json:"period"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type VerifyRequest struct {
	Secret string `json:"secret"`
	Token  string `json:"token"`
	Window int    `json:"window,omitempty"`
}

type VerifyResponse struct {
	Valid      bool       `json:"valid"`
	Counter    int64      `json:"counter,omitempty"`
	UsedBackup bool       `json:"used_backup,omitempty"`
	LastUsed   *time.Time `json:"last_used,omitempty"`
}

type ValidateRequest struct {
	Token  string `json:"token"`
	Window int    `json:"window,omitempty"`
}

type ValidateResponse struct {
	Valid     bool       `json:"valid"`
	Remaining int        `json:"remaining"`
	NextAt    *time.Time `json:"next_at,omitempty"`
}

type TOTPSecret struct {
	ID          string                 `json:"id"`
	AccountName string                 `json:"account_name"`
	Issuer      string                 `json:"issuer"`
	Secret      string                 `json:"secret,omitempty"`
	Algorithm   string                 `json:"algorithm"`
	Digits      int                    `json:"digits"`
	Period      int                    `json:"period"`
	BackupCodes []string               `json:"backup_codes"`
	IsActive    bool                   `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	LastUsed    *time.Time             `json:"last_used,omitempty"`
}

type ListRequest struct {
	Issuer      string `json:"issuer,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Offset      int    `json:"offset,omitempty"`
}

type ListResponse struct {
	Secrets []TOTPSecret `json:"secrets"`
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
}

type TOTPClient struct {
	client *client.Client
}

func NewTOTPClient(client *client.Client) *TOTPClient {
	return &TOTPClient{
		client: client,
	}
}

func (t *TOTPClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	resp, err := t.client.Post(ctx, "/api/v1/totp/generate", req)
	if err != nil {
		return nil, err
	}

	var generateResp GenerateResponse
	if err := resp.Decode(&generateResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP generate response")
	}

	return &generateResp, nil
}

func (t *TOTPClient) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	resp, err := t.client.Post(ctx, "/api/v1/totp/verify", req)
	if err != nil {
		return nil, err
	}

	var verifyResp VerifyResponse
	if err := resp.Decode(&verifyResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP verify response")
	}

	return &verifyResp, nil
}

func (t *TOTPClient) Validate(ctx context.Context, token string, window ...int) (*ValidateResponse, error) {
	req := &ValidateRequest{
		Token: token,
	}
	if len(window) > 0 {
		req.Window = window[0]
	}

	resp, err := t.client.Post(ctx, "/api/v1/totp/validate", req)
	if err != nil {
		return nil, err
	}

	var validateResp ValidateResponse
	if err := resp.Decode(&validateResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP validate response")
	}

	return &validateResp, nil
}

func (t *TOTPClient) Create(ctx context.Context, req *GenerateRequest) (*TOTPSecret, error) {
	resp, err := t.client.Post(ctx, "/api/v1/totp/secrets", req)
	if err != nil {
		return nil, err
	}

	var secret TOTPSecret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP secret response")
	}

	return &secret, nil
}

func (t *TOTPClient) Get(ctx context.Context, id string) (*TOTPSecret, error) {
	resp, err := t.client.Get(ctx, "/api/v1/totp/secrets/"+id)
	if err != nil {
		return nil, err
	}

	var secret TOTPSecret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP secret response")
	}

	return &secret, nil
}

func (t *TOTPClient) Update(ctx context.Context, id string, req map[string]interface{}) (*TOTPSecret, error) {
	resp, err := t.client.Put(ctx, "/api/v1/totp/secrets/"+id, req)
	if err != nil {
		return nil, err
	}

	var secret TOTPSecret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP secret response")
	}

	return &secret, nil
}

func (t *TOTPClient) Delete(ctx context.Context, id string) error {
	_, err := t.client.Delete(ctx, "/api/v1/totp/secrets/"+id)
	return err
}

func (t *TOTPClient) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	resp, err := t.client.Post(ctx, "/api/v1/totp/secrets/search", req)
	if err != nil {
		return nil, err
	}

	var listResp ListResponse
	if err := resp.Decode(&listResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP list response")
	}

	return &listResp, nil
}

func (t *TOTPClient) RegenerateBackupCodes(ctx context.Context, id string) (*TOTPSecret, error) {
	resp, err := t.client.Post(ctx, "/api/v1/totp/secrets/"+id+"/regenerate-backup-codes", nil)
	if err != nil {
		return nil, err
	}

	var secret TOTPSecret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode TOTP secret response")
	}

	return &secret, nil
}
