package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
)

// LoginController handles authentication commands
type LoginController struct {
	localClient  *services.LocalClient
	loginService *services.LoginService
}

// NewLoginController creates a new login controller
func NewLoginController(localClient *services.LocalClient, loginService *services.LoginService) *LoginController {
	return &LoginController{
		localClient:  localClient,
		loginService: loginService,
	}
}

// TokenLoginRequest represents a token login request
type TokenLoginRequest struct {
	Role  string `json:"role" binding:"required"`
	Token string `json:"token" binding:"required"`
}

// TokenLoginResponse represents a token login response
type TokenLoginResponse struct {
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
	Auth          interface{} `json:"auth"`
}

// CLIResponse represents a standard CLI response
type CLIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TokenLogin handles the token login command
func (c *LoginController) TokenLogin(ctx *gin.Context) {
	var req TokenLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate token format
	if req.Token == "" {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Token is required",
		})
		return
	}

	// Simple authentication logic for development
	// In production, validate against proper auth systems
	var response TokenLoginResponse

	if req.Token == "dev-token" || req.Token == "root-token" || req.Token == "vault-token" {
		// Successful authentication
		response.Auth.ClientToken = "s." + generateToken()
		response.Auth.Accessor = generateToken()
		response.Auth.Policies = []string{"root", "default"}
		response.Auth.TokenPolicies = []string{"root"}
		response.Auth.IdentityPolicies = []string{"root"}
		response.Auth.Metadata = map[string]interface{}{
			"role":        req.Role,
			"username":    "admin",
			"auth_method": "token",
		}
		response.Auth.LeaseDuration = 86400
		response.Auth.Renewable = true
		response.Auth.EntityID = "entity-" + generateToken()
		response.Auth.TokenType = "batch"
		response.Auth.Orphan = false
		response.Auth.Path = "auth/token/login"
		response.Auth.Namespace = ""

		ctx.JSON(http.StatusOK, response)
	} else {
		// Failed authentication
		ctx.JSON(http.StatusUnauthorized, CLIResponse{
			Success: false,
			Error:   "invalid token",
		})
	}
}

// UserPassLogin handles user/password login
func (c *LoginController) UserPassLogin(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Simple authentication logic
	if req.Username == "admin" && req.Password == "password" {
		var response TokenLoginResponse
		response.Auth.ClientToken = "s." + generateToken()
		response.Auth.Accessor = generateToken()
		response.Auth.Policies = []string{"admin", "default"}
		response.Auth.TokenPolicies = []string{"admin"}
		response.Auth.IdentityPolicies = []string{"admin"}
		response.Auth.Metadata = map[string]interface{}{
			"username":    req.Username,
			"auth_method": "userpass",
		}
		response.Auth.LeaseDuration = 86400
		response.Auth.Renewable = true
		response.Auth.EntityID = "entity-" + generateToken()
		response.Auth.TokenType = "service"
		response.Auth.Orphan = false
		response.Auth.Path = "auth/userpass/login/" + req.Username
		response.Auth.Namespace = ""

		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusUnauthorized, CLIResponse{
			Success: false,
			Error:   "invalid username or password",
		})
	}
}

// helper function to generate a simple token
func generateToken() string {
	return "x" + "a" + "b" + "c" + "d" + "e" + "f" + "1" + "2" + "3" + "4" + "5" + "6"
}
