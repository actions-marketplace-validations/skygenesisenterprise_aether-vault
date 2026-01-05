package audit

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
)

type AuditEvent struct {
	ID         string                 `json:"id"`
	IdentityID *string                `json:"identity_id,omitempty"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID *string                `json:"resource_id,omitempty"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type SearchAuditRequest struct {
	IdentityID *string    `json:"identity_id,omitempty"`
	Action     string     `json:"action,omitempty"`
	Resource   string     `json:"resource,omitempty"`
	ResourceID *string    `json:"resource_id,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
	Status     string     `json:"status,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}

type SearchAuditResponse struct {
	Events []AuditEvent `json:"events"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type AuditStats struct {
	TotalEvents    int64                    `json:"total_events"`
	EventsByAction map[string]int64         `json:"events_by_action"`
	EventsByStatus map[string]int64         `json:"events_by_status"`
	EventsByHour   []map[string]interface{} `json:"events_by_hour"`
	TimeRange      struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	} `json:"time_range"`
}

type AuditClient struct {
	client *client.Client
}

func NewAuditClient(client *client.Client) *AuditClient {
	return &AuditClient{
		client: client,
	}
}

func (a *AuditClient) Search(ctx context.Context, req *SearchAuditRequest) (*SearchAuditResponse, error) {
	resp, err := a.client.Post(ctx, "/api/v1/audit/search", req)
	if err != nil {
		return nil, err
	}

	var searchResp SearchAuditResponse
	if err := resp.Decode(&searchResp); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode audit search response")
	}

	return &searchResp, nil
}

func (a *AuditClient) Get(ctx context.Context, id string) (*AuditEvent, error) {
	resp, err := a.client.Get(ctx, "/api/v1/audit/events/"+id)
	if err != nil {
		return nil, err
	}

	var event AuditEvent
	if err := resp.Decode(&event); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode audit event response")
	}

	return &event, nil
}

func (a *AuditClient) GetStats(ctx context.Context, startTime, endTime *time.Time) (*AuditStats, error) {
	req := map[string]interface{}{}
	if startTime != nil {
		req["start_time"] = startTime
	}
	if endTime != nil {
		req["end_time"] = endTime
	}

	resp, err := a.client.Post(ctx, "/api/v1/audit/stats", req)
	if err != nil {
		return nil, err
	}

	var stats AuditStats
	if err := resp.Decode(&stats); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode audit stats response")
	}

	return &stats, nil
}

func (a *AuditClient) Export(ctx context.Context, req *SearchAuditRequest) ([]byte, error) {
	resp, err := a.client.Post(ctx, "/api/v1/audit/export", req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (a *AuditClient) Retention(ctx context.Context) (map[string]interface{}, error) {
	resp, err := a.client.Get(ctx, "/api/v1/audit/retention")
	if err != nil {
		return nil, err
	}

	var retention map[string]interface{}
	if err := resp.Decode(&retention); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode retention response")
	}

	return retention, nil
}

func (a *AuditClient) SetRetention(ctx context.Context, days int) error {
	req := map[string]interface{}{
		"days": days,
	}

	_, err := a.client.Put(ctx, "/api/v1/audit/retention", req)
	return err
}
