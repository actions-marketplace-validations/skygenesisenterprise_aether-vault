package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type ProtocolType string

const (
	ProtocolHTTP   ProtocolType = "http"
	ProtocolHTTPS  ProtocolType = "https"
	ProtocolSSH    ProtocolType = "ssh"
	ProtocolFTP    ProtocolType = "ftp"
	ProtocolSFTP   ProtocolType = "sftp"
	ProtocolWebDAV ProtocolType = "webdav"
	ProtocolSMB    ProtocolType = "smb"
	ProtocolNFS    ProtocolType = "nfs"
	ProtocolRsync  ProtocolType = "rsync"
	ProtocolGit    ProtocolType = "git"
	ProtocolCustom ProtocolType = "custom"
)

type Network struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;unique"`
	Type      ProtocolType   `json:"type" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:'active'"`
	Config    string         `json:"config" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ProtocolConfig struct {
	Host        string                 `json:"host"`
	Port        int                    `json:"port"`
	Username    string                 `json:"username,omitempty"`
	Password    string                 `json:"password,omitempty"`
	PrivateKey  string                 `json:"private_key,omitempty"`
	Certificate string                 `json:"certificate,omitempty"`
	Timeout     int                    `json:"timeout,omitempty"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

type NetworkRequest struct {
	Name   string          `json:"name" binding:"required"`
	Type   ProtocolType    `json:"type" binding:"required"`
	Config *ProtocolConfig `json:"config"`
}

type NetworkResponse struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	Type      ProtocolType    `json:"type"`
	Status    string          `json:"status"`
	Config    *ProtocolConfig `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ProtocolStatus struct {
	Protocol  ProtocolType `json:"protocol"`
	Status    string       `json:"status"`
	LastCheck time.Time    `json:"last_check"`
	Message   string       `json:"message,omitempty"`
	Latency   int64        `json:"latency,omitempty"`
}

type ProtocolTestRequest struct {
	Type   ProtocolType    `json:"type" binding:"required"`
	Config *ProtocolConfig `json:"config" binding:"required"`
}

type ProtocolTestResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Latency int64                  `json:"latency,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (n *Network) GetProtocolConfig() (*ProtocolConfig, error) {
	if n.Config == "" {
		return nil, nil
	}

	var config ProtocolConfig
	if err := json.Unmarshal([]byte(n.Config), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (n *Network) SetProtocolConfig(config *ProtocolConfig) error {
	if config == nil {
		n.Config = ""
		return nil
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	n.Config = string(data)
	return nil
}
