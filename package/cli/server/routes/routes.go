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
	networkController   *controllers.NetworkController
	snmpController      *controllers.SNMPController
	vaultController     *controllers.VaultController
	authMiddleware      *middleware.AuthMiddleware
	auditMiddleware     *middleware.AuditMiddleware
	rateLimitMiddleware *middleware.RateLimitMiddleware
	networkMiddleware   *middleware.NetworkMiddleware
	snmpMiddleware      *middleware.SNMPMiddleware
}

func NewRouter(
	db *gorm.DB,
	authService *services.AuthService,
	secretService *services.SecretService,
	totpService *services.TOTPService,
	userService *services.UserService,
	policyService *services.PolicyService,
	auditService *services.AuditService,
	networkService *services.NetworkService,
	snmpService *services.SNMPService,
	vaultController *controllers.VaultController,
) *Router {
	authController := controllers.NewAuthController(authService, auditService)
	secretController := controllers.NewSecretController(secretService)
	totpController := controllers.NewTOTPController(totpService)
	identityController := controllers.NewIdentityController(userService, policyService)
	auditController := controllers.NewAuditController(auditService)
	systemController := controllers.NewSystemController(db)
	userController := controllers.NewUserController(userService, auditService)
	networkController := controllers.NewNetworkController(networkService)
	snmpController := controllers.NewSNMPController(snmpService)

	authMiddleware := middleware.NewAuthMiddleware(authService)
	auditMiddleware := middleware.NewAuditMiddleware(auditService)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(100, 60) // 100 requests per minute

	networkConfig := &middleware.NetworkConfig{
		MaxRequestsPerMinute: 50,
		MaxConcurrent:        5,
		TimeoutSeconds:       30,
	}
	networkMiddleware := middleware.NewNetworkMiddleware(networkConfig)

	snmpConfig := &middleware.SNMPConfig{
		MaxRequestsPerMinute: 30,
		AllowedNetworks:      []string{"127.0.0.1/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		TimeoutSeconds:       30,
	}
	snmpMiddleware := middleware.NewSNMPMiddleware(snmpConfig)

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
		networkController:   networkController,
		snmpController:      snmpController,
		authMiddleware:      authMiddleware,
		auditMiddleware:     auditMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		networkMiddleware:   networkMiddleware,
		snmpMiddleware:      snmpMiddleware,
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

	network := v1.Group("/network")
	network.Use(r.authMiddleware.RequireAuth())
	network.Use(r.networkMiddleware.ValidateProtocol())
	network.Use(r.networkMiddleware.NetworkRateLimit())
	network.Use(r.networkMiddleware.ProtocolSecurity())
	network.Use(r.networkMiddleware.NetworkLogging())
	{
		network.GET("", r.networkController.GetNetworks)
		network.POST("", r.networkController.CreateNetwork)
		network.GET("/:id", r.networkController.GetNetwork)
		network.PUT("/:id", r.networkController.UpdateNetwork)
		network.DELETE("/:id", r.networkController.DeleteNetwork)

		network.GET("/protocols", r.networkController.GetSupportedProtocols)
		network.POST("/test", r.networkController.TestProtocol)
		network.GET("/:id/status", r.networkController.GetProtocolStatus)
	}

	snmp := v1.Group("/snmp")
	snmp.Use(r.authMiddleware.RequireAuth())
	snmp.Use(r.snmpMiddleware.RateLimit())
	snmp.Use(r.snmpMiddleware.SecurityHeaders())
	snmp.Use(r.snmpMiddleware.ValidateTarget())
	snmp.Use(r.snmpMiddleware.LogSNMPRequest())
	{
		snmp.POST("/get", r.snmpController.GetSNMPData)
		snmp.POST("/walk", r.snmpController.WalkSNMP)
		snmp.POST("/test", r.snmpController.TestConnection)
	}

	system := v1.Group("/system")
	{
		system.GET("/health", r.systemController.Health)
		system.GET("/version", r.systemController.Version)
	}

	// Vault-compatible API routes
	r.setupVaultRoutes()
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
