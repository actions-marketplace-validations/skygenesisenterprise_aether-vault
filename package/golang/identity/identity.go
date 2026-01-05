package identity

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type Identity struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	FirstName string                 `json:"first_name,omitempty"`
	LastName  string                 `json:"last_name,omitempty"`
	Status    string                 `json:"status"`
	Roles     []string               `json:"roles,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	LastLogin *time.Time             `json:"last_login,omitempty"`
}

type CreateIdentityRequest struct {
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	FirstName string                 `json:"first_name,omitempty"`
	LastName  string                 `json:"last_name,omitempty"`
	Password  string                 `json:"password"`
	Roles     []string               `json:"roles,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateIdentityRequest struct {
	Email     *string                `json:"email,omitempty"`
	FirstName *string                `json:"first_name,omitempty"`
	LastName  *string                `json:"last_name,omitempty"`
	Status    *string                `json:"status,omitempty"`
	Roles     []string               `json:"roles,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ListIdentitiesRequest struct {
	Username string   `json:"username,omitempty"`
	Email    string   `json:"email,omitempty"`
	Status   string   `json:"status,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

type ListIdentitiesResponse struct {
	Identities []Identity `json:"identities"`
	Total      int        `json:"total"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

type Session struct {
	ID           string                 `json:"id"`
	IdentityID   string                 `json:"identity_id"`
	Token        string                 `json:"token,omitempty"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
	IsActive     bool                   `json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	ExpiresAt    *time.Time             `json:"expires_at,omitempty"`
	LastActivity *time.Time             `json:"last_activity,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type ListSessionsRequest struct {
	IdentityID string `json:"identity_id,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

type ListSessionsResponse struct {
	Sessions []Session `json:"sessions"`
	Total    int       `json:"total"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}

type IdentityClient struct {
	client *client.Client
}

func NewIdentityClient(client *client.Client) *IdentityClient {
	return &IdentityClient{
		client: client,
	}
}

func (i *IdentityClient) Create(ctx context.Context, req *CreateIdentityRequest) (*Identity, error) {
	resp, err := i.client.Post(ctx, "/api/v1/identities", req)
	if err != nil {
		return nil, err
	}

	var identity Identity
	if err := resp.Decode(&identity); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode identity response")
	}

	return &identity, nil
}

func (i *IdentityClient) Get(ctx context.Context, id string) (*Identity, error) {
	resp, err := i.client.Get(ctx, "/api/v1/identities/"+id)
	if err != nil {
		return nil, err
	}

	var identity Identity
	if err := resp.Decode(&identity); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode identity response")
	}

	return &identity, nil
}

func (i *IdentityClient) GetByUsername(ctx context.Context, username string) (*Identity, error) {
	resp, err := i.client.Get(ctx, "/api/v1/identities/username/"+username)
	if err != nil {
		return nil, err
	}

	var identity Identity
	if err := resp.Decode(&identity); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode identity response")
	}

	return &identity, nil
}

func (i *IdentityClient) Update(ctx context.Context, id string, req *UpdateIdentityRequest) (*Identity, error) {
	resp, err := i.client.Put(ctx, "/api/v1/identities/"+id, req)
	if err != nil {
		return nil, err
	}

	var identity Identity
	if err := resp.Decode(&identity); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode identity response")
	}

	return &identity, nil
}

func (i *IdentityClient) Delete(ctx context.Context, id string) error {
	_, err := i.client.Delete(ctx, "/api/v1/identities/"+id)
	return err
}

func (i *IdentityClient) List(ctx context.Context, req *ListIdentitiesRequest) (*ListIdentitiesResponse, error) {
	resp, err := i.client.Post(ctx, "/api/v1/identities/search", req)
	if err != nil {
		return nil, err
	}

	var listResp ListIdentitiesResponse
	if err := resp.Decode(&listResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode identities list response")
	}

	return &listResp, nil
}

func (i *IdentityClient) ChangePassword(ctx context.Context, id string, req *ChangePasswordRequest) error {
	_, err := i.client.Post(ctx, "/api/v1/identities/"+id+"/change-password", req)
	return err
}

func (i *IdentityClient) Enable(ctx context.Context, id string) error {
	_, err := i.client.Post(ctx, "/api/v1/identities/"+id+"/enable", nil)
	return err
}

func (i *IdentityClient) Disable(ctx context.Context, id string) error {
	_, err := i.client.Post(ctx, "/api/v1/identities/"+id+"/disable", nil)
	return err
}

func (i *IdentityClient) Lock(ctx context.Context, id string) error {
	_, err := i.client.Post(ctx, "/api/v1/identities/"+id+"/lock", nil)
	return err
}

func (i *IdentityClient) Unlock(ctx context.Context, id string) error {
	_, err := i.client.Post(ctx, "/api/v1/identities/"+id+"/unlock", nil)
	return err
}

func (i *IdentityClient) GetSessions(ctx context.Context, req *ListSessionsRequest) (*ListSessionsResponse, error) {
	resp, err := i.client.Post(ctx, "/api/v1/identities/sessions/search", req)
	if err != nil {
		return nil, err
	}

	var listResp ListSessionsResponse
	if err := resp.Decode(&listResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode sessions list response")
	}

	return &listResp, nil
}

func (i *IdentityClient) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := i.client.Delete(ctx, "/api/v1/identities/sessions/"+sessionID)
	return err
}

func (i *IdentityClient) RevokeAllSessions(ctx context.Context, identityID string) error {
	_, err := i.client.Delete(ctx, "/api/v1/identities/"+identityID+"/sessions")
	return err
}
