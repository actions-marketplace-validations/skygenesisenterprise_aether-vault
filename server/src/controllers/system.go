package controllers

import (
	"aether-vault/src/model"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemController struct {
	db *gorm.DB
}

func NewSystemController(db *gorm.DB) *SystemController {
	return &SystemController{
		db: db,
	}
}

func (c *SystemController) Health(ctx *gin.Context) {
	status := "healthy"
	dbStatus := "connected"

	sqlDB, err := c.db.DB()
	if err != nil {
		status = "unhealthy"
		dbStatus = "disconnected"
	} else {
		if err := sqlDB.Ping(); err != nil {
			status = "unhealthy"
			dbStatus = "disconnected"
		}
	}

	response := model.HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Database:  dbStatus,
	}

	if status == "unhealthy" {
		ctx.JSON(http.StatusServiceUnavailable, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SystemController) Version(ctx *gin.Context) {
	response := model.VersionResponse{
		Version:   "1.0.0",
		BuildTime: "2024-01-01T00:00:00Z",
		GitCommit: "unknown",
		GoVersion: runtime.Version(),
	}

	ctx.JSON(http.StatusOK, response)
}
