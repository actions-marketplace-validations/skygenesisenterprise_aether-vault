package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Secret struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Value       string         `gorm:"type:text;not null" json:"-"`
	ValueHash   string         `gorm:"not null" json:"-"`
	Type        SecretType     `gorm:"not null" json:"type"`
	Tags        string         `gorm:"type:text" json:"tags"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type SecretType string

const (
	SecretTypePassword    SecretType = "password"
	SecretTypeAPIKey      SecretType = "api_key"
	SecretTypeToken       SecretType = "token"
	SecretTypeCertificate SecretType = "certificate"
	SecretTypeOther       SecretType = "other"
)

func (s *Secret) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
