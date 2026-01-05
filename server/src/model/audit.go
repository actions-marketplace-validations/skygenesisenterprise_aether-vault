package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Action     string     `gorm:"not null" json:"action"`
	Resource   string     `json:"resource"`
	ResourceID *string    `json:"resource_id"`
	IPAddress  string     `gorm:"not null" json:"ip_address"`
	UserAgent  string     `gorm:"type:text" json:"user_agent"`
	Success    bool       `gorm:"default:true" json:"success"`
	Details    string     `gorm:"type:text" json:"details"`
	CreatedAt  time.Time  `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
