package middleware

import (
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserMiddleware struct {
	userService *services.UserService
}

func NewUserMiddleware(userService *services.UserService) *UserMiddleware {
	return &UserMiddleware{
		userService: userService,
	}
}

func (m *UserMiddleware) RequireOwnership() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_UNAUTHORIZED",
					Message: "Unauthorized",
				},
			})
			ctx.Abort()
			return
		}

		targetUserIDStr := ctx.Param("id")
		if targetUserIDStr == "" {
			targetUserIDStr = ctx.Param("user_id")
		}

		if targetUserIDStr == "" {
			ctx.Next()
			return
		}

		targetUserID, err := uuid.Parse(targetUserIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_INVALID_ID",
					Message: "Invalid user ID",
				},
			})
			ctx.Abort()
			return
		}

		currentUserUUID := currentUserID.(uuid.UUID)

		if currentUserUUID != targetUserID {
			ctx.JSON(http.StatusForbidden, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_ACCESS_DENIED",
					Message: "Access denied: insufficient permissions",
				},
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *UserMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_UNAUTHORIZED",
					Message: "Unauthorized",
				},
			})
			ctx.Abort()
			return
		}

		user, err := m.userService.GetUserByID(currentUserID.(uuid.UUID))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_UNAUTHORIZED",
					Message: "Unauthorized",
				},
			})
			ctx.Abort()
			return
		}

		if !m.isAdmin(user) {
			ctx.JSON(http.StatusForbidden, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_ACCESS_DENIED",
					Message: "Access denied: admin privileges required",
				},
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *UserMiddleware) RequireActiveUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_UNAUTHORIZED",
					Message: "Unauthorized",
				},
			})
			ctx.Abort()
			return
		}

		user, err := m.userService.GetUserByID(currentUserID.(uuid.UUID))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_UNAUTHORIZED",
					Message: "Unauthorized",
				},
			})
			ctx.Abort()
			return
		}

		if !user.IsActive {
			ctx.JSON(http.StatusForbidden, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_ACCOUNT_DISABLED",
					Message: "Account is disabled",
				},
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *UserMiddleware) ValidateUserExists() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr := ctx.Param("id")
		if userIDStr == "" {
			userIDStr = ctx.Param("user_id")
		}

		if userIDStr == "" {
			ctx.Next()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_INVALID_ID",
					Message: "Invalid user ID",
				},
			})
			ctx.Abort()
			return
		}

		_, err = m.userService.GetUserByID(userID)
		if err != nil {
			if err == services.ErrUserNotFound {
				ctx.JSON(http.StatusNotFound, model.ErrorResponse{
					Error: model.ErrorDetail{
						Code:    "VAULT_USER_NOT_FOUND",
						Message: "User not found",
					},
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: model.ErrorDetail{
					Code:    "VAULT_INTERNAL_ERROR",
					Message: "Failed to validate user",
				},
			})
			return
		}

		ctx.Next()
	}
}

func (m *UserMiddleware) isAdmin(user *model.User) bool {
	return user.Email == "admin@aether-vault.local"
}

func (m *UserMiddleware) CanAccessResource(ctx *gin.Context, resourceType string, resourceID uuid.UUID) bool {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		return false
	}

	user, err := m.userService.GetUserByID(currentUserID.(uuid.UUID))
	if err != nil {
		return false
	}

	if m.isAdmin(user) {
		return true
	}

	switch resourceType {
	case "secret", "totp", "policy":
		return true
	case "user":
		targetUserIDStr := ctx.Param("id")
		if targetUserIDStr == "" {
			return false
		}
		targetUserID, err := uuid.Parse(targetUserIDStr)
		if err != nil {
			return false
		}
		return currentUserID.(uuid.UUID) == targetUserID
	default:
		return false
	}
}
