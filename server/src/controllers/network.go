package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"github.com/skygenesisenterprise/aether-vault/server/src/services"
)

type NetworkController struct {
	networkService *services.NetworkService
}

func NewNetworkController(networkService *services.NetworkService) *NetworkController {
	return &NetworkController{
		networkService: networkService,
	}
}

func (c *NetworkController) GetNetworks(ctx *gin.Context) {
	networks, err := c.networkService.GetNetworks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]model.NetworkResponse, len(networks))
	for i, network := range networks {
		config, _ := network.GetProtocolConfig()
		response[i] = model.NetworkResponse{
			ID:        network.ID,
			Name:      network.Name,
			Type:      network.Type,
			Status:    network.Status,
			Config:    config,
			CreatedAt: network.CreatedAt,
			UpdatedAt: network.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkController) GetNetwork(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid network ID"})
		return
	}

	network, err := c.networkService.GetNetwork(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	config, _ := network.GetProtocolConfig()
	response := model.NetworkResponse{
		ID:        network.ID,
		Name:      network.Name,
		Type:      network.Type,
		Status:    network.Status,
		Config:    config,
		CreatedAt: network.CreatedAt,
		UpdatedAt: network.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkController) CreateNetwork(ctx *gin.Context) {
	var req model.NetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network, err := c.networkService.CreateNetwork(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config, _ := network.GetProtocolConfig()
	response := model.NetworkResponse{
		ID:        network.ID,
		Name:      network.Name,
		Type:      network.Type,
		Status:    network.Status,
		Config:    config,
		CreatedAt: network.CreatedAt,
		UpdatedAt: network.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *NetworkController) UpdateNetwork(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid network ID"})
		return
	}

	var req model.NetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network, err := c.networkService.UpdateNetwork(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	config, _ := network.GetProtocolConfig()
	response := model.NetworkResponse{
		ID:        network.ID,
		Name:      network.Name,
		Type:      network.Type,
		Status:    network.Status,
		Config:    config,
		CreatedAt: network.CreatedAt,
		UpdatedAt: network.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkController) DeleteNetwork(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid network ID"})
		return
	}

	if err := c.networkService.DeleteNetwork(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *NetworkController) TestProtocol(ctx *gin.Context) {
	var req model.ProtocolTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.networkService.TestProtocol(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkController) GetProtocolStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid network ID"})
		return
	}

	status, err := c.networkService.GetProtocolStatus(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, status)
}

func (c *NetworkController) GetSupportedProtocols(ctx *gin.Context) {
	protocols := c.networkService.GetSupportedProtocols()
	ctx.JSON(http.StatusOK, gin.H{
		"protocols": protocols,
		"count":     len(protocols),
	})
}
