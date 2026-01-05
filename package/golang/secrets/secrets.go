package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type Secret struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Value       string                 `json:"value,omitempty"`
	Description string                 `json:"description,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	Version     int                    `json:"version"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type CreateSecretRequest struct {
	Name        string                 `json:"name"`
	Value       string                 `json:"value"`
	Description string                 `json:"description,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateSecretRequest struct {
	Value       *string                `json:"value,omitempty"`
	Description *string                `json:"description,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type ListSecretsRequest struct {
	Prefix string            `json:"prefix,omitempty"`
	Tags   map[string]string `json:"tags,omitempty"`
	Limit  int               `json:"limit,omitempty"`
	Offset int               `json:"offset,omitempty"`
}

type ListSecretsResponse struct {
	Secrets []Secret `json:"secrets"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

type SecretVersion struct {
	ID        string                 `json:"id"`
	SecretID  string                 `json:"secret_id"`
	Version   int                    `json:"version"`
	Value     string                 `json:"value"`
	CreatedAt time.Time              `json:"created_at"`
	CreatedBy string                 `json:"created_by"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ListVersionsResponse struct {
	Versions []SecretVersion `json:"versions"`
	Total    int             `json:"total"`
}

type SecretsClient struct {
	client *client.Client
}

func NewSecretsClient(client *client.Client) *SecretsClient {
	return &SecretsClient{
		client: client,
	}
}

func (s *SecretsClient) Create(ctx context.Context, req *CreateSecretRequest) (*Secret, error) {
	resp, err := s.client.Post(ctx, "/api/v1/secrets", req)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode secret response")
	}

	return &secret, nil
}

func (s *SecretsClient) Get(ctx context.Context, name string, version ...int) (*Secret, error) {
	endpoint := "/api/v1/secrets/" + name
	if len(version) > 0 && version[0] > 0 {
		endpoint += "?version=" + fmt.Sprintf("%d", version[0])
	}

	resp, err := s.client.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode secret response")
	}

	return &secret, nil
}

func (s *SecretsClient) Update(ctx context.Context, name string, req *UpdateSecretRequest) (*Secret, error) {
	resp, err := s.client.Put(ctx, "/api/v1/secrets/"+name, req)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode secret response")
	}

	return &secret, nil
}

func (s *SecretsClient) Delete(ctx context.Context, name string) error {
	_, err := s.client.Delete(ctx, "/api/v1/secrets/"+name)
	return err
}

func (s *SecretsClient) List(ctx context.Context, req *ListSecretsRequest) (*ListSecretsResponse, error) {
	resp, err := s.client.Post(ctx, "/api/v1/secrets/search", req)
	if err != nil {
		return nil, err
	}

	var listResp ListSecretsResponse
	if err := resp.Decode(&listResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode secrets list response")
	}

	return &listResp, nil
}

func (s *SecretsClient) ListVersions(ctx context.Context, name string) (*ListVersionsResponse, error) {
	resp, err := s.client.Get(ctx, "/api/v1/secrets/"+name+"/versions")
	if err != nil {
		return nil, err
	}

	var versionsResp ListVersionsResponse
	if err := resp.Decode(&versionsResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode versions response")
	}

	return &versionsResp, nil
}

func (s *SecretsClient) Restore(ctx context.Context, name string, version int) (*Secret, error) {
	req := map[string]interface{}{
		"version": version,
	}

	resp, err := s.client.Post(ctx, "/api/v1/secrets/"+name+"/restore", req)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := resp.Decode(&secret); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode restored secret response")
	}

	return &secret, nil
}

func (s *SecretsClient) Exists(ctx context.Context, name string) (bool, error) {
	resp, err := s.client.Get(ctx, "/api/v1/secrets/"+name+"/exists")
	if err != nil {
		if errors.GetErrorCode(err) == errors.ErrCodeNotFound {
			return false, nil
		}
		return false, err
	}

	var existsResp struct {
		Exists bool `json:"exists"`
	}

	if err := resp.Decode(&existsResp); err != nil {
		return false, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode exists response")
	}

	return existsResp.Exists, nil
}
