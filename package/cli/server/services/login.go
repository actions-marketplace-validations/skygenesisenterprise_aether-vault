package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/server/model"
)

// LoginService handles authentication logic
type LoginService struct {
	sessions map[string]*model.SessionInfo
}

// NewLoginService creates a new login service
func NewLoginService() *LoginService {
	return &LoginService{
		sessions: make(map[string]*model.SessionInfo),
	}
}

// AuthenticateToken validates a token and creates a session
func (s *LoginService) AuthenticateToken(req *model.TokenAuthRequest, ipAddress, userAgent string) (*model.AuthResponse, error) {
	// Validate token format
	if !model.ValidateToken(req.Token) {
		return nil, fmt.Errorf("invalid token format")
	}

	// Check against allowed development tokens
	allowedTokens := []string{
		"dev-token",
		"root-token",
		"vault-token",
		"test-token",
		"admin-token",
	}

	valid := false
	for _, token := range allowedTokens {
		if req.Token == token {
			valid = true
			break
		}
	}

	if !valid {
		// For development, also accept any token with "vault" in it
		if !strings.Contains(req.Token, "vault") && !strings.Contains(req.Token, "token") {
			return nil, fmt.Errorf("invalid token")
		}
	}

	// Generate session
	clientToken := s.generateClientToken()
	accessor := s.generateAccessor()

	// Determine policies based on role
	policies := s.getPoliciesForRole(req.Role)

	session := &model.SessionInfo{
		Token:       req.Token,
		ClientToken: clientToken,
		Accessor:    accessor,
		Policies:    policies,
		Metadata: map[string]interface{}{
			"role":           req.Role,
			"auth_method":    "token",
			"original_token": req.Token,
		},
		LeaseDuration: 86400,
		Renewable:     true,
		EntityID:      s.generateEntityID(),
		TokenType:     "batch",
		Path:          "auth/token/login",
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(time.Duration(86400) * time.Second),
		LastActivity:  time.Now(),
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	}

	// Store session
	s.sessions[clientToken] = session

	// Build Vault-compatible response
	response := &model.AuthResponse{}
	response.Auth.ClientToken = clientToken
	response.Auth.Accessor = accessor
	response.Auth.Policies = policies
	response.Auth.TokenPolicies = policies
	response.Auth.IdentityPolicies = policies
	response.Auth.Metadata = session.Metadata
	response.Auth.LeaseDuration = 86400
	response.Auth.Renewable = true
	response.Auth.EntityID = session.EntityID
	response.Auth.TokenType = "batch"
	response.Auth.Orphan = false
	response.Auth.Path = "auth/token/login"
	response.Auth.Namespace = ""

	return response, nil
}

// AuthenticateUserPass validates username/password and creates a session
func (s *LoginService) AuthenticateUserPass(req *model.UserPassAuthRequest, ipAddress, userAgent string) (*model.AuthResponse, error) {
	// Simple validation for development
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("username and password required")
	}

	// Check credentials (development)
	validCredentials := false
	if req.Username == "admin" && req.Password == "password" {
		validCredentials = true
	}
	if req.Username == "root" && req.Password == "root" {
		validCredentials = true
	}
	if req.Username == "vault" && req.Password == "vault" {
		validCredentials = true
	}

	if !validCredentials {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Generate session
	clientToken := s.generateClientToken()
	accessor := s.generateAccessor()

	// Determine policies based on username
	policies := s.getPoliciesForUsername(req.Username)

	session := &model.SessionInfo{
		Token:       req.Username,
		ClientToken: clientToken,
		Accessor:    accessor,
		Policies:    policies,
		Metadata: map[string]interface{}{
			"username":    req.Username,
			"auth_method": "userpass",
		},
		LeaseDuration: 86400,
		Renewable:     true,
		EntityID:      s.generateEntityID(),
		TokenType:     "service",
		Path:          "auth/userpass/login/" + req.Username,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(time.Duration(86400) * time.Second),
		LastActivity:  time.Now(),
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	}

	// Store session
	s.sessions[clientToken] = session

	// Build Vault-compatible response
	response := &model.AuthResponse{}
	response.Auth.ClientToken = clientToken
	response.Auth.Accessor = accessor
	response.Auth.Policies = policies
	response.Auth.TokenPolicies = policies
	response.Auth.IdentityPolicies = policies
	response.Auth.Metadata = session.Metadata
	response.Auth.LeaseDuration = 86400
	response.Auth.Renewable = true
	response.Auth.EntityID = session.EntityID
	response.Auth.TokenType = "service"
	response.Auth.Orphan = false
	response.Auth.Path = "auth/userpass/login/" + req.Username
	response.Auth.Namespace = ""

	return response, nil
}

// ValidateSession validates a client token
func (s *LoginService) ValidateSession(clientToken string) (*model.SessionInfo, error) {
	session, exists := s.sessions[clientToken]
	if !exists {
		return nil, fmt.Errorf("invalid session")
	}

	if session.IsExpired() {
		delete(s.sessions, clientToken)
		return nil, fmt.Errorf("session expired")
	}

	session.UpdateLastActivity()
	return session, nil
}

// RevokeSession revokes a session
func (s *LoginService) RevokeSession(clientToken string) error {
	if _, exists := s.sessions[clientToken]; !exists {
		return fmt.Errorf("session not found")
	}

	delete(s.sessions, clientToken)
	return nil
}

// GetSession returns session information
func (s *LoginService) GetSession(clientToken string) (*model.SessionInfo, error) {
	return s.ValidateSession(clientToken)
}

// CleanupExpiredSessions removes expired sessions
func (s *LoginService) CleanupExpiredSessions() {
	for token, session := range s.sessions {
		if session.IsExpired() {
			delete(s.sessions, token)
		}
	}
}

// GetActiveSessions returns all active sessions
func (s *LoginService) GetActiveSessions() map[string]*model.SessionInfo {
	s.CleanupExpiredSessions()
	return s.sessions
}

// Helper functions

func (s *LoginService) generateClientToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "s." + hex.EncodeToString(bytes)[:24]
}

func (s *LoginService) generateAccessor() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *LoginService) generateEntityID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "entity-" + hex.EncodeToString(bytes)
}

func (s *LoginService) getPoliciesForRole(role string) []string {
	switch role {
	case "root":
		return []string{"root", "default"}
	case "admin":
		return []string{"admin", "default"}
	default:
		return []string{"default"}
	}
}

func (s *LoginService) getPoliciesForUsername(username string) []string {
	switch username {
	case "admin", "root":
		return []string{"admin", "root", "default"}
	case "vault":
		return []string{"vault", "default"}
	default:
		return []string{"default"}
	}
}
