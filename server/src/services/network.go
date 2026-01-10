package services

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"gorm.io/gorm"
)

type NetworkService struct {
	db      *gorm.DB
	clients map[model.ProtocolType]ProtocolClient
}

type ProtocolClient interface {
	Test(config *model.ProtocolConfig) (*model.ProtocolTestResponse, error)
	GetStatus(config *model.ProtocolConfig) (*model.ProtocolStatus, error)
}

type HTTPClient struct{}

func (c *HTTPClient) Test(config *model.ProtocolConfig) (*model.ProtocolTestResponse, error) {
	start := time.Now()

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &model.ProtocolTestResponse{
			Success: false,
			Message: fmt.Sprintf("HTTP request failed: %v", err),
			Latency: time.Since(start).Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	return &model.ProtocolTestResponse{
		Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
		Message: fmt.Sprintf("HTTP %d %s", resp.StatusCode, resp.Status),
		Latency: time.Since(start).Milliseconds(),
		Details: map[string]interface{}{
			"status_code": resp.StatusCode,
			"headers":     resp.Header,
		},
	}, nil
}

func (c *HTTPClient) GetStatus(config *model.ProtocolConfig) (*model.ProtocolStatus, error) {
	test, err := c.Test(config)
	if err != nil {
		return nil, err
	}

	status := "active"
	if !test.Success {
		status = "inactive"
	}

	return &model.ProtocolStatus{
		Protocol:  model.ProtocolHTTP,
		Status:    status,
		LastCheck: time.Now(),
		Message:   test.Message,
		Latency:   test.Latency,
	}, nil
}

type HTTPSClient struct{}

func (c *HTTPSClient) Test(config *model.ProtocolConfig) (*model.ProtocolTestResponse, error) {
	start := time.Now()

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := fmt.Sprintf("https://%s:%d", config.Host, config.Port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &model.ProtocolTestResponse{
			Success: false,
			Message: fmt.Sprintf("HTTPS request failed: %v", err),
			Latency: time.Since(start).Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	return &model.ProtocolTestResponse{
		Success: resp.StatusCode >= 200 && resp.StatusCode < 300,
		Message: fmt.Sprintf("HTTPS %d %s", resp.StatusCode, resp.Status),
		Latency: time.Since(start).Milliseconds(),
		Details: map[string]interface{}{
			"status_code": resp.StatusCode,
			"tls_version": resp.TLS,
			"headers":     resp.Header,
		},
	}, nil
}

func (c *HTTPSClient) GetStatus(config *model.ProtocolConfig) (*model.ProtocolStatus, error) {
	test, err := c.Test(config)
	if err != nil {
		return nil, err
	}

	status := "active"
	if !test.Success {
		status = "inactive"
	}

	return &model.ProtocolStatus{
		Protocol:  model.ProtocolHTTPS,
		Status:    status,
		LastCheck: time.Now(),
		Message:   test.Message,
		Latency:   test.Latency,
	}, nil
}

type TCPClient struct{}

func (c *TCPClient) Test(config *model.ProtocolConfig) (*model.ProtocolTestResponse, error) {
	start := time.Now()

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(config.Timeout)*time.Second)
	if err != nil {
		return &model.ProtocolTestResponse{
			Success: false,
			Message: fmt.Sprintf("TCP connection failed: %v", err),
			Latency: time.Since(start).Milliseconds(),
		}, nil
	}
	defer conn.Close()

	return &model.ProtocolTestResponse{
		Success: true,
		Message: "TCP connection successful",
		Latency: time.Since(start).Milliseconds(),
		Details: map[string]interface{}{
			"remote_addr": conn.RemoteAddr().String(),
			"local_addr":  conn.LocalAddr().String(),
		},
	}, nil
}

func (c *TCPClient) GetStatus(config *model.ProtocolConfig) (*model.ProtocolStatus, error) {
	test, err := c.Test(config)
	if err != nil {
		return nil, err
	}

	status := "active"
	if !test.Success {
		status = "inactive"
	}

	return &model.ProtocolStatus{
		Protocol:  model.ProtocolSSH,
		Status:    status,
		LastCheck: time.Now(),
		Message:   test.Message,
		Latency:   test.Latency,
	}, nil
}

func NewNetworkService(db *gorm.DB) *NetworkService {
	clients := make(map[model.ProtocolType]ProtocolClient)
	clients[model.ProtocolHTTP] = &HTTPClient{}
	clients[model.ProtocolHTTPS] = &HTTPSClient{}
	clients[model.ProtocolSSH] = &TCPClient{}
	clients[model.ProtocolFTP] = &TCPClient{}
	clients[model.ProtocolSFTP] = &TCPClient{}
	clients[model.ProtocolSMB] = &TCPClient{}
	clients[model.ProtocolNFS] = &TCPClient{}
	clients[model.ProtocolRsync] = &TCPClient{}
	clients[model.ProtocolGit] = &TCPClient{}
	clients[model.ProtocolWebDAV] = &HTTPClient{}
	clients[model.ProtocolCustom] = &TCPClient{}

	return &NetworkService{
		db:      db,
		clients: clients,
	}
}

func (s *NetworkService) GetNetworks() ([]model.Network, error) {
	var networks []model.Network
	if err := s.db.Find(&networks).Error; err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}
	return networks, nil
}

func (s *NetworkService) GetNetwork(id uint) (*model.Network, error) {
	var network model.Network
	if err := s.db.First(&network, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("network not found")
		}
		return nil, fmt.Errorf("failed to get network: %w", err)
	}
	return &network, nil
}

func (s *NetworkService) CreateNetwork(req *model.NetworkRequest) (*model.Network, error) {
	network := &model.Network{
		Name:   req.Name,
		Type:   req.Type,
		Status: "active",
	}

	if err := network.SetProtocolConfig(req.Config); err != nil {
		return nil, fmt.Errorf("failed to set protocol config: %w", err)
	}

	if err := s.db.Create(network).Error; err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	return network, nil
}

func (s *NetworkService) UpdateNetwork(id uint, req *model.NetworkRequest) (*model.Network, error) {
	network, err := s.GetNetwork(id)
	if err != nil {
		return nil, err
	}

	network.Name = req.Name
	network.Type = req.Type

	if err := network.SetProtocolConfig(req.Config); err != nil {
		return nil, fmt.Errorf("failed to set protocol config: %w", err)
	}

	if err := s.db.Save(network).Error; err != nil {
		return nil, fmt.Errorf("failed to update network: %w", err)
	}

	return network, nil
}

func (s *NetworkService) DeleteNetwork(id uint) error {
	if err := s.db.Delete(&model.Network{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete network: %w", err)
	}
	return nil
}

func (s *NetworkService) GetNetworkByName(name string) (*model.Network, error) {
	var network model.Network
	if err := s.db.Where("name = ?", name).First(&network).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("network not found")
		}
		return nil, fmt.Errorf("failed to get network by name: %w", err)
	}
	return &network, nil
}

func (s *NetworkService) TestProtocol(req *model.ProtocolTestRequest) (*model.ProtocolTestResponse, error) {
	client, exists := s.clients[req.Type]
	if !exists {
		return &model.ProtocolTestResponse{
			Success: false,
			Message: fmt.Sprintf("Protocol %s not supported", req.Type),
		}, nil
	}

	return client.Test(req.Config)
}

func (s *NetworkService) GetProtocolStatus(id uint) (*model.ProtocolStatus, error) {
	network, err := s.GetNetwork(id)
	if err != nil {
		return nil, err
	}

	client, exists := s.clients[network.Type]
	if !exists {
		return &model.ProtocolStatus{
			Protocol: network.Type,
			Status:   "unsupported",
			Message:  fmt.Sprintf("Protocol %s not supported", network.Type),
		}, nil
	}

	config, err := network.GetProtocolConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol config: %w", err)
	}

	return client.GetStatus(config)
}

func (s *NetworkService) GetSupportedProtocols() []model.ProtocolType {
	protocols := make([]model.ProtocolType, 0, len(s.clients))
	for protocol := range s.clients {
		protocols = append(protocols, protocol)
	}
	return protocols
}
