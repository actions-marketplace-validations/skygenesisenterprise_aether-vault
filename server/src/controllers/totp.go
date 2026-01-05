package controllers

import (
	"aether-vault/src/model"
	"aether-vault/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TOTPController struct {
	totpService *services.TOTPService
}

func NewTOTPController(totpService *services.TOTPService) *TOTPController {
	return &TOTPController{
		totpService: totpService,
	}
}

func (c *TOTPController) GetTOTPs(ctx *gin.Context) {
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

	totps, err := c.totpService.GetTOTPsByUserID(userID.(uuid.UUID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to retrieve TOTPs",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"totps": totps})
}

func (c *TOTPController) CreateTOTP(ctx *gin.Context) {
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

	var req model.CreateTOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_REQUEST",
				Message: "Invalid request format",
			},
		})
		return
	}

	totp := &model.TOTP{
		Name:        req.Name,
		Description: req.Description,
		Secret:      req.Secret,
		Algorithm:   req.Algorithm,
		Digits:      req.Digits,
		Period:      req.Period,
		IsActive:    true,
	}

	if err := c.totpService.CreateTOTP(totp, userID.(uuid.UUID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to create TOTP",
			},
		})
		return
	}

	totp.Secret = ""
	ctx.JSON(http.StatusCreated, totp)
}

func (c *TOTPController) GenerateCode(ctx *gin.Context) {
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

	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_ID",
				Message: "Invalid TOTP ID",
			},
		})
		return
	}

	response, err := c.totpService.GenerateCode(id, userID.(uuid.UUID))
	if err != nil {
		if err == services.ErrTOTPNotFound {
			ctx.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_TOTP_NOT_FOUND",
					Message: "TOTP not found",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to generate TOTP code",
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
