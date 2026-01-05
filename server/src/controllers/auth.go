package controllers

import (
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	authService  *services.AuthService
	auditService *services.AuditService
}

func NewAuthController(authService *services.AuthService, auditService *services.AuditService) *AuthController {
	return &AuthController{
		authService:  authService,
		auditService: auditService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_REQUEST",
				Message: "Invalid request format",
			},
		})
		return
	}

	response, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		if c.auditService != nil {
			c.auditService.LogAnonymousAction("login_failed", "auth", "", ctx.ClientIP(), ctx.GetHeader("User-Agent"), false, err.Error())
		}

		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_CREDENTIALS",
				Message: "Invalid email or password",
			},
		})
		return
	}

	if c.auditService != nil {
		c.auditService.LogAnonymousAction("login_success", "auth", "", ctx.ClientIP(), ctx.GetHeader("User-Agent"), true, "")
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) Logout(ctx *gin.Context) {
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

	if c.auditService != nil {
		c.auditService.LogAction(userID.(uuid.UUID), "logout", "auth", "", true, "")
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (c *AuthController) GetSession(ctx *gin.Context) {
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

	response, err := c.authService.GetSession(userID.(uuid.UUID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
