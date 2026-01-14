package model

import (
	"time"
)

// TokenAuthRequest represents a token authentication request
type TokenAuthRequest struct {
	Role  string `json:"role" binding:"required"`
	Token string `json:"token" binding:"required"`
}

// UserPassAuthRequest represents a user/password authentication request
type UserPassAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents a Vault-compatible authentication response
type AuthResponse struct {
	Auth struct {
		ClientToken      string                 `json:"client_token"`
		Accessor         string                 `json:"accessor"`
		Policies         []string               `json:"policies"`
		TokenPolicies    []string               `json:"token_policies"`
		IdentityPolicies []string               `json:"identity_policies"`
		Metadata         map[string]interface{} `json:"metadata"`
		LeaseDuration    int                    `json:"lease_duration"`
		Renewable        bool                   `json:"renewable"`
		EntityID         string                 `json:"entity_id"`
		TokenType        string                 `json:"token_type"`
		Orphan           bool                   `json:"orphan"`
		Path             string                 `json:"path"`
		Namespace        string                 `json:"namespace"`
	} `json:"auth"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          interface{} `json:"data"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      []string    `json:"warnings"`
}

// LoginAttempt represents a login attempt record (for auditing)
type LoginAttempt struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Method    string    `json:"method"` // token, userpass, etc.
	Success   bool      `json:"success"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// SessionInfo represents an active session
type SessionInfo struct {
	Token         string                 `json:"token"`
	ClientToken   string                 `json:"client_token"`
	Accessor      string                 `json:"accessor"`
	Policies      []string               `json:"policies"`
	Metadata      map[string]interface{} `json:"metadata"`
	LeaseDuration int                    `json:"lease_duration"`
	Renewable     bool                   `json:"renewable"`
	EntityID      string                 `json:"entity_id"`
	TokenType     string                 `json:"token_type"`
	Path          string                 `json:"path"`
	CreatedAt     time.Time              `json:"created_at"`
	ExpiresAt     time.Time              `json:"expires_at"`
	LastActivity  time.Time              `json:"last_activity"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
}

// IsExpired checks if the session is expired
func (s *SessionInfo) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// UpdateLastActivity updates the last activity timestamp
func (s *SessionInfo) UpdateLastActivity() {
	s.LastActivity = time.Now()
}

// ValidateToken validates a token format
func ValidateToken(token string) bool {
	if len(token) < 8 {
		return false
	}

	// Basic validation - no spaces
	for _, char := range token {
		if char == ' ' {
			return false
		}
	}

	return true
}
