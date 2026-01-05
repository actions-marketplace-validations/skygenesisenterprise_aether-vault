package main

import (
	"github.com/skygenesisenterprise/aether-vault/server/src/config"
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/server/src/routes"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := migrateDatabase(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userService := services.NewUserService(db)
	auditService := services.NewAuditService(db)
	authService := services.NewAuthService(userService, &cfg.JWT)
	secretService := services.NewSecretService(db, cfg.Security.EncryptionKey, "default-salt", cfg.Security.KDFIterations, auditService)
	totpService := services.NewTOTPService(db, auditService)
	policyService := services.NewPolicyService(db)

	router := routes.NewRouter(db, authService, secretService, totpService, userService, policyService, auditService)
	router.SetupRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router.GetEngine(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	log.Printf("Aether Vault API server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
		dbConfig.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
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
	return db.AutoMigrate(
		&model.User{},
		&model.Secret{},
		&model.TOTP{},
		&model.Policy{},
		&model.AuditLog{},
	)
}
