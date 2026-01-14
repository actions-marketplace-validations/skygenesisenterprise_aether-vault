package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/types"
)

// VaultController handles Vault-compatible API endpoints
type VaultController struct {
	localClient *services.LocalClient
	config      *types.ClientConfig
}

// NewVaultController creates a new Vault controller
func NewVaultController(config *types.ClientConfig) (*VaultController, error) {
	clientOptions := &types.ClientOptions{
		Config: config,
	}

	localClient, err := services.NewLocalClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create local client: %w", err)
	}

	return &VaultController{
		localClient: localClient,
		config:      config,
	}, nil
}

// handleHealth handles health check requests
func (c *VaultController) HandleHealth(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	response := map[string]interface{}{
		"initialized":                  true,
		"sealed":                       false,
		"standby":                      false,
		"performance_standby":          false,
		"replication_performance_mode": "none",
		"replication_dr_mode":          "disabled",
		"server_time_utc":              time.Now().UTC().Unix(),
		"version":                      "1.0.0",
		"cluster_name":                 "aether-vault",
		"cluster_id":                   "local",
	}
	ctx.JSON(http.StatusOK, response)
}

// handleTokenLogin handles token authentication
func (c *VaultController) HandleTokenLogin(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var req struct {
		Role  string `json:"role"`
		Token string `json:"token"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Simple token validation - in production, this should be more secure
	if req.Token != "dev-token" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	response := map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": req.Token,
			"policies":     []string{"root"},
			"metadata": map[string]interface{}{
				"role": req.Role,
			},
			"lease_duration": 0,
			"renewable":      false,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// handleSecret handles secret CRUD operations
func (c *VaultController) HandleSecret(ctx *gin.Context) {
	path := ctx.Param("path")
	ctx.Header("Content-Type", "application/json")

	switch ctx.Request.Method {
	case "GET":
		c.handleGetSecret(ctx, path)
	case "POST":
		c.handleSetSecret(ctx, path)
	case "DELETE":
		c.handleDeleteSecret(ctx, path)
	}
}

// handleGetSecret retrieves a secret
func (c *VaultController) handleGetSecret(ctx *gin.Context, path string) {
	secret, err := c.localClient.GetSecret(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get secret: %v", err),
		})
		return
	}

	response := map[string]interface{}{
		"data":     secret.Data,
		"metadata": secret.Metadata,
	}
	ctx.JSON(http.StatusOK, response)
}

// handleSetSecret stores a secret
func (c *VaultController) handleSetSecret(ctx *gin.Context, path string) {
	var req struct {
		Data map[string]interface{} `json:"data"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	secret := &types.Secret{
		Path: path,
		Data: req.Data,
		Metadata: &types.SecretMetadata{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			CreatedBy: "vault-server",
			UpdatedBy: "vault-server",
		},
		Version: 1,
	}

	if err := c.localClient.SetSecret(path, secret); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to set secret: %v", err),
		})
		return
	}

	response := map[string]interface{}{
		"data": req.Data,
	}
	ctx.JSON(http.StatusOK, response)
}

// handleDeleteSecret removes a secret
func (c *VaultController) handleDeleteSecret(ctx *gin.Context, path string) {
	if err := c.localClient.DeleteSecret(path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to delete secret: %v", err),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// handleSecretList handles secret listing
func (c *VaultController) HandleSecretList(ctx *gin.Context) {
	prefix := ctx.Param("path")

	secrets, err := c.localClient.ListSecrets(prefix)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to list secrets: %v", err),
		})
		return
	}

	ctx.Header("Content-Type", "application/json")
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"keys": secrets,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// handleStatus handles status requests
func (c *VaultController) HandleStatus(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	status, err := c.localClient.GetStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get status: %v", err),
		})
		return
	}

	response := map[string]interface{}{
		"initialized": true,
		"sealed":      false,
		"standby":     false,
		"version":     "1.0.0",
		"server_time": time.Now().UTC().Unix(),
		"mode":        string(status.Mode),
		"local_path":  status.LocalPath,
	}

	ctx.JSON(http.StatusOK, response)
}

// Close closes the controller and its resources
func (c *VaultController) Close() error {
	if c.localClient != nil {
		return c.localClient.Close()
	}
	return nil
}
