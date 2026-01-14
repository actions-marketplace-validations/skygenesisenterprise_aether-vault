package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/types"
)

// KVController handles key-value operations
type KVController struct {
	localClient *services.LocalClient
}

// NewKVController creates a new KV controller
func NewKVController(localClient *services.LocalClient) *KVController {
	return &KVController{
		localClient: localClient,
	}
}

// KVRequest represents a KV request
type KVRequest struct {
	Path    string                 `json:"path" binding:"required"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Version int64                  `json:"version,omitempty"`
	Method  string                 `json:"method,omitempty"`
}

// PatchRequest represents a patch request
type PatchRequest struct {
	Path string                 `json:"path" binding:"required"`
	Data map[string]interface{} `json:"data" binding:"required"`
	Op   string                 `json:"op,omitempty"`
}

// KVGet handles KV get operations
func (c *KVController) Get(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Path parameter is required",
		})
		return
	}

	secret, err := c.localClient.GetSecret(path)
	if err != nil {
		ctx.JSON(http.StatusNotFound, CLIResponse{
			Success: false,
			Error:   "Secret not found: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"data": secret.Data,
			"metadata": map[string]interface{}{
				"created_time": time.Unix(secret.Metadata.CreatedAt, 0).Format(time.RFC3339),
				"updated_time": time.Unix(secret.Metadata.UpdatedAt, 0).Format(time.RFC3339),
				"version":      secret.Version,
				"created_by":   secret.Metadata.CreatedBy,
				"updated_by":   secret.Metadata.UpdatedBy,
			},
		},
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Data:    response,
	})
}

// KVPut handles KV put operations
func (c *KVController) Put(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Path parameter is required",
		})
		return
	}

	var req struct {
		Data map[string]interface{} `json:"data" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Get existing secret to determine version
	existingSecret, _ := c.localClient.GetSecret(path)
	var version int64 = 1
	if existingSecret != nil {
		version = existingSecret.Version + 1
	}

	secret := &types.Secret{
		Path: path,
		Data: req.Data,
		Metadata: &types.SecretMetadata{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			CreatedBy: "kv-api",
			UpdatedBy: "kv-api",
		},
		Version: version,
	}

	err := c.localClient.SetSecret(path, secret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to write secret: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"data": req.Data,
			"metadata": map[string]interface{}{
				"created_time": time.Unix(secret.Metadata.CreatedAt, 0).Format(time.RFC3339),
				"version":      secret.Version,
				"created_by":   secret.Metadata.CreatedBy,
			},
		},
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Message: "Successfully wrote to " + path,
		Data:    response,
	})
}

// KVDelete handles KV delete operations
func (c *KVController) Delete(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Path parameter is required",
		})
		return
	}

	versions := ctx.Query("versions")
	force := ctx.Query("force") == "true"

	// Check if secret exists
	existingSecret, err := c.localClient.GetSecret(path)
	if err != nil {
		ctx.JSON(http.StatusNotFound, CLIResponse{
			Success: false,
			Error:   "Secret not found: " + err.Error(),
		})
		return
	}

	// Delete secret
	err = c.localClient.DeleteSecret(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to delete secret: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"path":    path,
		"version": existingSecret.Version,
	}

	if versions != "" {
		response["deleted_versions"] = versions
	}

	if force {
		response["force"] = true
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Message: "Successfully deleted " + path,
		Data:    response,
	})
}

// KVPatch handles KV patch operations
func (c *KVController) Patch(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Path parameter is required",
		})
		return
	}

	var req PatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Get existing secret
	existingSecret, err := c.localClient.GetSecret(path)
	if err != nil {
		ctx.JSON(http.StatusNotFound, CLIResponse{
			Success: false,
			Error:   "Secret not found: " + err.Error(),
		})
		return
	}

	// Apply patch based on operation
	switch req.Op {
	case "replace":
		for key, value := range req.Data {
			existingSecret.Data[key] = value
		}
	case "add":
		for key, value := range req.Data {
			existingSecret.Data[key] = value
		}
	case "remove":
		for key := range req.Data {
			delete(existingSecret.Data, key)
		}
	default:
		// Default: merge/patch
		for key, value := range req.Data {
			existingSecret.Data[key] = value
		}
	}

	// Update metadata
	existingSecret.Metadata.UpdatedAt = time.Now().Unix()
	existingSecret.Metadata.UpdatedBy = "kv-patch-api"
	existingSecret.Version++

	err = c.localClient.SetSecret(path, existingSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to patch secret: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"data": existingSecret.Data,
			"metadata": map[string]interface{}{
				"updated_time": time.Unix(existingSecret.Metadata.UpdatedAt, 0).Format(time.RFC3339),
				"version":      existingSecret.Version,
				"updated_by":   existingSecret.Metadata.UpdatedBy,
			},
		},
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Message: "Successfully patched " + path,
		Data:    response,
	})
}

// KVList handles KV list operations
func (c *KVController) List(ctx *gin.Context) {
	path := ctx.Param("path")

	keys, err := c.localClient.ListSecrets(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CLIResponse{
			Success: false,
			Error:   "Failed to list secrets: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"keys": keys,
		},
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Data:    response,
	})
}

// KVMetadata handles KV metadata operations
func (c *KVController) Metadata(ctx *gin.Context) {
	path := ctx.Param("path")
	if path == "" {
		ctx.JSON(http.StatusBadRequest, CLIResponse{
			Success: false,
			Error:   "Path parameter is required",
		})
		return
	}

	secret, err := c.localClient.GetSecret(path)
	if err != nil {
		ctx.JSON(http.StatusNotFound, CLIResponse{
			Success: false,
			Error:   "Secret not found: " + err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"metadata": map[string]interface{}{
				"created_time": time.Unix(secret.Metadata.CreatedAt, 0).Format(time.RFC3339),
				"updated_time": time.Unix(secret.Metadata.UpdatedAt, 0).Format(time.RFC3339),
				"version":      secret.Version,
				"created_by":   secret.Metadata.CreatedBy,
				"updated_by":   secret.Metadata.UpdatedBy,
				"path":         secret.Path,
			},
		},
	}

	ctx.JSON(http.StatusOK, CLIResponse{
		Success: true,
		Data:    response,
	})
}
