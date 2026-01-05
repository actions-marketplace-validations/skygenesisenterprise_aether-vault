package policies

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type Policy struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Actions     []string               `json:"actions"`
	Effect      string                 `json:"effect"`
	Priority    int                    `json:"priority"`
	IsActive    bool                   `json:"is_active"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	UserID      string                 `json:"user_id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type CreatePolicyRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Actions     []string               `json:"actions"`
	Effect      string                 `json:"effect"`
	Priority    int                    `json:"priority"`
	IsActive    bool                   `json:"is_active"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type UpdatePolicyRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Resource    *string                `json:"resource,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Effect      *string                `json:"effect,omitempty"`
	Priority    *int                   `json:"priority,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type ListPoliciesRequest struct {
	Name     string `json:"name,omitempty"`
	Resource string `json:"resource,omitempty"`
	Effect   string `json:"effect,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

type ListPoliciesResponse struct {
	Policies []Policy `json:"policies"`
	Total    int      `json:"total"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

type PolicyCheckRequest struct {
	IdentityID string                 `json:"identity_id"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

type PolicyCheckResponse struct {
	Allowed   bool      `json:"allowed"`
	Reason    string    `json:"reason,omitempty"`
	Policies  []string  `json:"policies,omitempty"`
	Evaluated time.Time `json:"evaluated"`
}

type PolicyAssignment struct {
	ID         string                 `json:"id"`
	IdentityID string                 `json:"identity_id"`
	PolicyID   string                 `json:"policy_id"`
	IsActive   bool                   `json:"is_active"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type CreateAssignmentRequest struct {
	IdentityID string                 `json:"identity_id"`
	PolicyID   string                 `json:"policy_id"`
	IsActive   bool                   `json:"is_active"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type ListAssignmentsRequest struct {
	IdentityID string `json:"identity_id,omitempty"`
	PolicyID   string `json:"policy_id,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

type ListAssignmentsResponse struct {
	Assignments []PolicyAssignment `json:"assignments"`
	Total       int                `json:"total"`
	Limit       int                `json:"limit"`
	Offset      int                `json:"offset"`
}

type PolicyClient struct {
	client *client.Client
}

func NewPolicyClient(client *client.Client) *PolicyClient {
	return &PolicyClient{
		client: client,
	}
}

func (p *PolicyClient) Create(ctx context.Context, req *CreatePolicyRequest) (*Policy, error) {
	resp, err := p.client.Post(ctx, "/api/v1/policies", req)
	if err != nil {
		return nil, err
	}

	var policy Policy
	if err := resp.Decode(&policy); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode policy response")
	}

	return &policy, nil
}

func (p *PolicyClient) Get(ctx context.Context, id string) (*Policy, error) {
	resp, err := p.client.Get(ctx, "/api/v1/policies/"+id)
	if err != nil {
		return nil, err
	}

	var policy Policy
	if err := resp.Decode(&policy); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode policy response")
	}

	return &policy, nil
}

func (p *PolicyClient) Update(ctx context.Context, id string, req *UpdatePolicyRequest) (*Policy, error) {
	resp, err := p.client.Put(ctx, "/api/v1/policies/"+id, req)
	if err != nil {
		return nil, err
	}

	var policy Policy
	if err := resp.Decode(&policy); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode policy response")
	}

	return &policy, nil
}

func (p *PolicyClient) Delete(ctx context.Context, id string) error {
	_, err := p.client.Delete(ctx, "/api/v1/policies/"+id)
	return err
}

func (p *PolicyClient) List(ctx context.Context, req *ListPoliciesRequest) (*ListPoliciesResponse, error) {
	resp, err := p.client.Post(ctx, "/api/v1/policies/search", req)
	if err != nil {
		return nil, err
	}

	var listResp ListPoliciesResponse
	if err := resp.Decode(&listResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode policies list response")
	}

	return &listResp, nil
}

func (p *PolicyClient) Check(ctx context.Context, req *PolicyCheckRequest) (*PolicyCheckResponse, error) {
	resp, err := p.client.Post(ctx, "/api/v1/policies/check", req)
	if err != nil {
		return nil, err
	}

	var checkResp PolicyCheckResponse
	if err := resp.Decode(&checkResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode policy check response")
	}

	return &checkResp, nil
}

func (p *PolicyClient) Enable(ctx context.Context, id string) error {
	_, err := p.client.Post(ctx, "/api/v1/policies/"+id+"/enable", nil)
	return err
}

func (p *PolicyClient) Disable(ctx context.Context, id string) error {
	_, err := p.client.Post(ctx, "/api/v1/policies/"+id+"/disable", nil)
	return err
}

func (p *PolicyClient) Assign(ctx context.Context, req *CreateAssignmentRequest) (*PolicyAssignment, error) {
	resp, err := p.client.Post(ctx, "/api/v1/policies/assignments", req)
	if err != nil {
		return nil, err
	}

	var assignment PolicyAssignment
	if err := resp.Decode(&assignment); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode assignment response")
	}

	return &assignment, nil
}

func (p *PolicyClient) RemoveAssignment(ctx context.Context, id string) error {
	_, err := p.client.Delete(ctx, "/api/v1/policies/assignments/"+id)
	return err
}

func (p *PolicyClient) ListAssignments(ctx context.Context, req *ListAssignmentsRequest) (*ListAssignmentsResponse, error) {
	resp, err := p.client.Post(ctx, "/api/v1/policies/assignments/search", req)
	if err != nil {
		return nil, err
	}

	var listResp ListAssignmentsResponse
	if err := resp.Decode(&listResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode assignments list response")
	}

	return &listResp, nil
}

func (p *PolicyClient) GetIdentityPolicies(ctx context.Context, identityID string) ([]Policy, error) {
	resp, err := p.client.Get(ctx, "/api/v1/policies/identity/"+identityID)
	if err != nil {
		return nil, err
	}

	var policies []Policy
	if err := resp.Decode(&policies); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode identity policies response")
	}

	return policies, nil
}
