package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
)

type UserController struct {
	userService  *services.UserService
	auditService *services.AuditService
	db           *gorm.DB
}

func NewUserController(userService *services.UserService, auditService *services.AuditService) *UserController {
	return &UserController{
		userService:  userService,
		auditService: auditService,
		db:           userService.GetDB(),
	}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	var users []model.User
	if err := c.db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to retrieve users",
			},
		})
		return
	}

	for i := range users {
		users[i].Password = ""
		users[i].Secrets = nil
		users[i].TOTPs = nil
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_ID",
				Message: "Invalid user ID",
			},
		})
		return
	}

	user, err := c.userService.GetUserByID(id)
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
				Message: "Failed to retrieve user",
			},
		})
		return
	}

	user.Password = ""
	user.Secrets = nil
	user.TOTPs = nil

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_REQUEST",
				Message: "Invalid request format",
			},
		})
		return
	}

	user := &model.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := c.userService.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to create user",
			},
		})
		return
	}

	user.Password = ""

	if c.auditService != nil {
		c.auditService.LogAction(user.ID, "user_created", "user", user.ID.String(), true, "")
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_ID",
				Message: "Invalid user ID",
			},
		})
		return
	}

	var req struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		IsActive  *bool   `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_REQUEST",
				Message: "Invalid request format",
			},
		})
		return
	}

	user, err := c.userService.GetUserByID(id)
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
				Message: "Failed to retrieve user",
			},
		})
		return
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := c.userService.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to update user",
			},
		})
		return
	}

	user.Password = ""
	user.Secrets = nil
	user.TOTPs = nil

	if c.auditService != nil {
		c.auditService.LogAction(id, "user_updated", "user", user.ID.String(), true, "")
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INVALID_ID",
				Message: "Invalid user ID",
			},
		})
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: model.ErrorDetail{
				Code:    "VAULT_INTERNAL_ERROR",
				Message: "Failed to delete user",
			},
		})
		return
	}

	if c.auditService != nil {
		c.auditService.LogAction(id, "user_deleted", "user", id.String(), true, "")
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
