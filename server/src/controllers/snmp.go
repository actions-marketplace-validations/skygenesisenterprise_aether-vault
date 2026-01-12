package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
)

type SNMPController struct {
	snmpService *services.SNMPService
}

func NewSNMPController(snmpService *services.SNMPService) *SNMPController {
	return &SNMPController{
		snmpService: snmpService,
	}
}

func (c *SNMPController) GetSNMPData(ctx *gin.Context) {
	var req services.SNMPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Target == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Target is required"})
		return
	}

	if req.Port == 0 {
		req.Port = 161
	}

	if len(req.OIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one OID is required"})
		return
	}

	response, err := c.snmpService.GetSNMPData(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SNMPController) WalkSNMP(ctx *gin.Context) {
	var req services.SNMPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Target == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Target is required"})
		return
	}

	if req.Port == 0 {
		req.Port = 161
	}

	response, err := c.snmpService.WalkSNMP(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *SNMPController) TestConnection(ctx *gin.Context) {
	var req services.SNMPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Target == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Target is required"})
		return
	}

	if req.Port == 0 {
		req.Port = 161
	}

	err := c.snmpService.TestConnection(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SNMP connection successful",
	})
}
