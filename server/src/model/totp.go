package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TOTP struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Secret      string         `gorm:"type:text;not null" json:"-"`
	Algorithm   string         `gorm:"default:SHA1" json:"algorithm"`
	Digits      int            `gorm:"default:6" json:"digits"`
	Period      int            `gorm:"default:30" json:"period"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (t *TOTP) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
