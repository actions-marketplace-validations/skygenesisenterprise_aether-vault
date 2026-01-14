package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/skygenesisenterprise/aether-vault/package/cli/server/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/src/routes"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/src/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Printf("🚀 Starting Aether Vault CLI Server")
	log.Printf("📍 Server will listen on %s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🔧 Environment: %s", cfg.Server.Environment)

	// Initialize database
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Printf("⚠️  Running in development mode without database")
		db = nil
	}

	if db != nil {
		if err := migrateDatabase(db); err != nil {
			log.Printf("⚠️  Database migration failed: %v", err)
			db = nil
		} else {
			log.Printf("✅ Database connected and migrated")
			if cfg.Database.Type == "sqlite" {
				log.Printf("📁 SQLite database: %s", cfg.Database.SQLitePath)
			} else {
				log.Printf("🐘 PostgreSQL database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
			}
		}
	}

	// Initialize services
	var userService *services.UserService
	var auditService *services.AuditService
	var secretService *services.SecretService
	var totpService *services.TOTPService
	var policyService *services.PolicyService
	var networkService *services.NetworkService
	var snmpService *services.SNMPService

	if db != nil {
		userService = services.NewUserService(db)
		auditService = services.NewAuditService(db)
		secretService = services.NewSecretService(db, cfg.Security.EncryptionKey, "cli-salt", cfg.Security.KDFIterations, auditService)
		totpService = services.NewTOTPService(db, auditService)
		policyService = services.NewPolicyService(db)
		networkService = services.NewNetworkService(db)
		snmpService = services.NewSNMPService()
		log.Printf("✅ All services initialized with database backend")
	} else {
		log.Printf("🔧 Initializing services in development mode")
		networkService = services.NewNetworkService(nil)
		snmpService = services.NewSNMPService()
		log.Printf("⚠️  Some features may be limited without database")
	}

	// Initialize auth service
	authService := services.NewAuthService(userService, &cfg.JWT)

	// Setup routes
	router := routes.NewRouter(db, authService, secretService, totpService, userService, policyService, auditService, networkService, snmpService)
	router.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router.GetEngine(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Server listening on http://%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("📋 API endpoints available at http://%s:%d/api/v1/", cfg.Server.Host, cfg.Server.Port)
		log.Printf("🏥 Health check: http://%s:%d/api/v1/system/health", cfg.Server.Host, cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("🛑 Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown:", err)
	}

	log.Printf("✅ Server exited")
}

func initDatabase(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if dbConfig.Type == "sqlite" {
		db, err = gorm.Open(sqlite.Open(dbConfig.SQLitePath), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to sqlite database: %w", err)
		}
		log.Printf("📁 Connecting to SQLite database: %s", dbConfig.SQLitePath)
	} else if dbConfig.Type == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.DBName,
			dbConfig.Port,
			dbConfig.SSLMode,
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
		}
		log.Printf("🐘 Connecting to PostgreSQL database: %s:%d/%s", dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	} else {
		return nil, fmt.Errorf("unsupported database type: %s", dbConfig.Type)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func migrateDatabase(db *gorm.DB) error {
	log.Printf("🔄 Running database migrations...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Secret{},
		&model.TOTP{},
		&model.Policy{},
		&model.AuditLog{},
	)

	if err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	log.Printf("✅ Database migrations completed")
	return nil
}
