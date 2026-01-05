package controllers

import (
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditController struct {
	auditService *services.AuditService
}

func NewAuditController(auditService *services.AuditService) *AuditController {
	return &AuditController{
		auditService: auditService,
	}
}

func (c *AuditController) GetAuditLogs(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_UNAUTHORIZED",
				Message: "Unauthorized",
			},
		})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "50")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	userIDUUID := userID.(uuid.UUID)
	logs, err := c.auditService.GetAuditLogs(&userIDUUID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to retrieve audit logs",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"logs": logs, "limit": limit, "offset": offset})
}
