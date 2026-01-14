package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/types"
)

// CLIController handles CLI command endpoints
type CLIController struct {
	localClient *services.LocalClient
}

// NewCLIController creates a new CLI controller
func NewCLIController(localClient *services.LocalClient) *CLIController {
	return &CLIController{
		localClient: localClient,
	}
}

// ReadRequest represents a read request
type ReadRequest struct {
	Path   string `json:"path" binding:"required"`
	Field  string `json:"field,omitempty"`
	Format string `json:"format,omitempty"`
}

// WriteRequest represents a write request
type WriteRequest struct {
	Path   string                 `json:"path" binding:"required"`
	Data   map[string]interface{} `json:"data" binding:"required"`
	Force  bool                   `json:"force,omitempty"`
	Format string                 `json:"format,omitempty"`
}

// DeleteRequest represents a delete request
type DeleteRequest struct {
	Path      string `json:"path" binding:"required"`
	Versions  string `json:"versions,omitempty"`
	Force     bool   `json:"force,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}

// ListRequest represents a list request
type ListRequest struct {
	Path   string `json:"path,omitempty"`
	Format string `json:"format,omitempty"`
}

// CLIResponse represents a CLI response
type CLIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Read handles the read command
func (c *CLIController) Read(ctx *gin.Context) {
	var req ReadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Get secret
	secret, err := c.localClient.GetSecret(req.Path)
	if err != nil {
		ctx.JSON(http.StatusNotFound, CLIResponse{
			Success: false,
			Error:   "Secret not found: " + err.Error(),
		})
		return
	}

	var responseData interface{}
	// Handle field extraction
	if req.Field != "" {
		if fieldValue, exists := secret.Data[req.Field]; exists {
			responseData = fieldValue
		} else {
			ctx.JSON(http.StatusNotFound, CLIResponse{
				Success: false,
				Error:   "Field '" + req.Field + "' not found in secret",
			})
			return
		}
	} else {
		responseData = secret.Data
	}

	// Format response based on format
	switch req.Format {
	case "json":
		ctx.JSON(http.StatusOK, CLIResponse{
			Success: true,
			Data:    responseData,
		})
	case "yaml":
		ctx.Header("Content-Type", "application/x-yaml")
		ctx.String(http.StatusOK, formatYAML(responseData))
	default:
		ctx.JSON(http.StatusOK, CLIResponse{
			Success: true,
			Data:    responseData,
		})
	}
}

// Write handles the write command
func (c *CLIController) Write(ctx *gin.Context) {
	var req WriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Check if secret exists (if not force)
	if !req.Force {
		_, err := c.localClient.GetSecret(req.Path)
		if err == nil {
			ctx.JSON(http.StatusConflict, CLIResponse{
				Success: false,
				Error:   "Secret already exists. Use force=true to overwrite.",
			})
			return
		}
	}

	// Create secret
	secret := &types.Secret{
		Path: req.Path,
		Data: req.Data,
		Metadata: &types.SecretMetadata{
			CreatedBy: "cli-api",
			UpdatedBy: "cli-api",
		},
		Version: 1,
	}

	// Write secret
	err := c.localClient.SetSecret(req.Path, secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to write secret: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Message: "Successfully wrote to " + req.Path,
		Data: map[string]interface{}{
			"path":    req.Path,
			"version": secret.Version,
		},
	})
}

// Delete handles the delete command
func (c *CLIController) Delete(ctx *gin.Context) {
	var req DeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Check if secret exists
	_, err := c.localClient.GetSecret(req.Path)
	if err != nil {
		ctx.JSON(http.StatusNotFound, CLIResponse{
			Success: false,
			Error:   "Secret not found: " + err.Error(),
		})
		return
	}

	// Delete secret
	err = c.localClient.DeleteSecret(req.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to delete secret: " + err.Error(),
		})
		return
	}

	response := CLIResponse{
		Success: true,
		Message: "Successfully deleted " + req.Path,
		Data: map[string]interface{}{
			"path": req.Path,
		},
	}

	responseData := response.Data.(map[string]interface{})
	if req.Versions != "" {
		responseData["deleted_versions"] = req.Versions
	}

	if req.Recursive {
		responseData["recursive"] = true
	}

	ctx.JSON(http.StatusOK, response)
}

// List handles the list command
func (c *CLIController) List(ctx *gin.Context) {
	var req ListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Allow empty body for list requests
		req.Path = ""
		req.Format = "table"
	}

	// List secrets
	keys, err := c.localClient.ListSecrets(req.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to list secrets: " + err.Error(),
		})
		return
	}

	var responseData interface{}
	switch req.Format {
	case "json":
		responseData = map[string]interface{}{
			"keys": keys,
		}
		ctx.JSON(http.StatusOK, CLIResponse{
			Success: true,
			Data:    responseData,
		})
	case "yaml":
		yamlData := map[string]interface{}{
			"keys": keys,
		}
		ctx.Header("Content-Type", "application/x-yaml")
		ctx.String(http.StatusOK, formatYAML(yamlData))
	default:
		responseData = map[string]interface{}{
			"keys": keys,
		}
		ctx.JSON(http.StatusOK, CLIResponse{
			Success: true,
			Data:    responseData,
		})
	}
}

// Status handles the status command
func (c *CLIController) Status(ctx *gin.Context) {
	status, err := c.localClient.GetStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to get status: " + err.Error(),
		})
		return
	}

	responseData := map[string]interface{}{
		"mode":          string(status.Mode),
		"connected":     status.Connected,
		"authenticated": status.Authenticated,
		"local_path":    status.LocalPath,
		"server_url":    status.ServerURL,
	}

	if status.LastSync != nil {
		responseData["last_sync"] = *status.LastSync
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Data:    responseData,
	})
}

// Helper function to format data as YAML
func formatYAML(data interface{}) string {
	// Simple YAML formatting - in production, use a proper YAML library
	yaml := "---\n"
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			yaml += key + ": " + formatValue(value) + "\n"
		}
	case []interface{}:
		for _, item := range v {
			yaml += "- " + formatValue(item) + "\n"
		}
	default:
		yaml += formatValue(v) + "\n"
	}
	return yaml
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "\"" + v + "\""
	case map[string]interface{}, []interface{}:
		return formatYAML(v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return "\"" + fmt.Sprintf("%v", v) + "\""
	}
}
