package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
)

// AuthController handles authentication commands
type AuthController struct {
	localClient *services.LocalClient
}

// NewAuthController creates a new auth controller
func NewAuthController(localClient *services.LocalClient) *AuthController {
	return &AuthController{
		localClient: localClient,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Method     string `json:"method,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Token      string `json:"token,omitempty"`
	Role       string `json:"role,omitempty"`
	ServerURL  string `json:"server_url,omitempty"`
	AuthMethod string `json:"auth_method,omitempty"`
}

// InitRequest represents an init request
type InitRequest struct {
	Path        string `json:"path,omitempty"`
	Environment string `json:"environment,omitempty"`
	Force       bool   `json:"force,omitempty"`
}

// Login handles the login command
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Simple authentication logic
	// In production, this should validate against proper auth systems
	if req.Token == "dev-token" || (req.Username == "admin" && req.Password == "password") {
		ctx.JSON(http.StatusOK, CLIResponse{
			Success: true,
			Message: "Successfully authenticated",
			Data: map[string]interface{}{
				"token":    "session-token-" + req.Username,
				"policies": []string{"root"},
				"metadata": map[string]interface{}{
					"username": req.Username,
					"role":     req.Role,
				},
				"lease_duration": 3600,
				"renewable":      true,
			},
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, CLIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
	}
}

// Logout handles the logout command
func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Message: "Successfully logged out",
	})
}

// Init handles the init command
func (c *AuthController) Init(ctx *gin.Context) {
	var req InitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Allow empty body for init
		req.Path = "./data"
		req.Environment = "development"
		req.Force = false
	}

	// Initialize local storage
	err := c.localClient.Authenticate("local", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to initialize: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Message: "Successfully initialized vault",
		Data: map[string]interface{}{
			"path":        req.Path,
			"environment": req.Environment,
			"initialized": true,
		},
	})
}
