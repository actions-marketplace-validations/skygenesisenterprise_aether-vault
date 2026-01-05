package routes

import (
	"github.com/skygenesisenterprise/aether-vault/server/src/controllers"
	"github.com/skygenesisenterprise/aether-vault/server/src/middleware"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	engine              *gin.Engine
	authController      *controllers.AuthController
	secretController    *controllers.SecretController
	totpController      *controllers.TOTPController
	identityController  *controllers.IdentityController
	auditController     *controllers.AuditController
	systemController    *controllers.SystemController
	userController      *controllers.UserController
	authMiddleware      *middleware.AuthMiddleware
	auditMiddleware     *middleware.AuditMiddleware
	rateLimitMiddleware *middleware.RateLimitMiddleware
}

func NewRouter(
	db *gorm.DB,
	authService *services.AuthService,
	secretService *services.SecretService,
	totpService *services.TOTPService,
	userService *services.UserService,
	policyService *services.PolicyService,
	auditService *services.AuditService,
) *Router {
	authController := controllers.NewAuthController(authService, auditService)
	secretController := controllers.NewSecretController(secretService)
	totpController := controllers.NewTOTPController(totpService)
	identityController := controllers.NewIdentityController(userService, policyService)
	auditController := controllers.NewAuditController(auditService)
	systemController := controllers.NewSystemController(db)
	userController := controllers.NewUserController(userService, auditService)

	authMiddleware := middleware.NewAuthMiddleware(authService)
	auditMiddleware := middleware.NewAuditMiddleware(auditService)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(100, 60) // 100 requests per minute

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.SecurityHeadersMiddleware())
	engine.Use(middleware.RequestIDMiddleware())
	engine.Use(rateLimitMiddleware.Limit())
	engine.Use(auditMiddleware.Audit())

	return &Router{
		engine:              engine,
		authController:      authController,
		secretController:    secretController,
		totpController:      totpController,
		identityController:  identityController,
		auditController:     auditController,
		systemController:    systemController,
		userController:      userController,
		authMiddleware:      authMiddleware,
		auditMiddleware:     auditMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

func (r *Router) SetupRoutes() {
	v1 := r.engine.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", r.authController.Login)
		auth.POST("/logout", r.authMiddleware.RequireAuth(), r.authController.Logout)
		auth.GET("/session", r.authMiddleware.RequireAuth(), r.authController.GetSession)
	}

	secrets := v1.Group("/secrets")
	secrets.Use(r.authMiddleware.RequireAuth())
	{
		secrets.GET("", r.secretController.GetSecrets)
		secrets.POST("", r.secretController.CreateSecret)
		secrets.GET("/:id", r.secretController.GetSecret)
		secrets.PUT("/:id", r.secretController.UpdateSecret)
		secrets.DELETE("/:id", r.secretController.DeleteSecret)
	}

	totp := v1.Group("/totp")
	totp.Use(r.authMiddleware.RequireAuth())
	{
		totp.GET("", r.totpController.GetTOTPs)
		totp.POST("", r.totpController.CreateTOTP)
		totp.POST("/:id/generate", r.totpController.GenerateCode)
	}

	identity := v1.Group("/identity")
	identity.Use(r.authMiddleware.RequireAuth())
	{
		identity.GET("/me", r.identityController.GetMe)
		identity.GET("/policies", r.identityController.GetPolicies)
	}

	users := v1.Group("/users")
	users.Use(r.authMiddleware.RequireAuth())
	{
		users.GET("", r.userController.GetUsers)
		users.GET("/:id", r.userController.GetUser)
		users.POST("", r.userController.CreateUser)
		users.PUT("/:id", r.userController.UpdateUser)
		users.DELETE("/:id", r.userController.DeleteUser)
	}

	audit := v1.Group("/audit")
	audit.Use(r.authMiddleware.RequireAuth())
	{
		audit.GET("/logs", r.auditController.GetAuditLogs)
	}

	system := v1.Group("/system")
	{
		system.GET("/health", r.systemController.Health)
		system.GET("/version", r.systemController.Version)
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
