package services

import (
	"aether-vault/src/model"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TOTPService struct {
	db           *gorm.DB
	auditService *AuditService
}

func NewTOTPService(db *gorm.DB, auditService *AuditService) *TOTPService {
	return &TOTPService{
		db:           db,
		auditService: auditService,
	}
}

func (s *TOTPService) CreateTOTP(totp *model.TOTP, userID uuid.UUID) error {
	if totp.Secret == "" {
		secret, err := s.generateSecret()
		if err != nil {
			return fmt.Errorf("failed to generate TOTP secret: %w", err)
		}
		totp.Secret = secret
	}

	if totp.Algorithm == "" {
		totp.Algorithm = "SHA1"
	}
	if totp.Digits == 0 {
		totp.Digits = 6
	}
	if totp.Period == 0 {
		totp.Period = 30
	}

	totp.UserID = userID

	if err := s.db.Create(totp).Error; err != nil {
		return fmt.Errorf("failed to create TOTP: %w", err)
	}

	if s.auditService != nil {
		s.auditService.LogAction(userID, "totp_created", "totp", totp.ID.String(), true, "")
	}

	return nil
}

func (s *TOTPService) GetTOTPsByUserID(userID uuid.UUID) ([]model.TOTP, error) {
	var totps []model.TOTP
	if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&totps).Error; err != nil {
		return nil, fmt.Errorf("failed to get TOTPs: %w", err)
	}

	for i := range totps {
		totps[i].Secret = ""
	}

	if s.auditService != nil {
		s.auditService.LogAction(userID, "totps_listed", "totp", "", true, "")
	}

	return totps, nil
}

func (s *TOTPService) GetTOTPByID(id uuid.UUID, userID uuid.UUID) (*model.TOTP, error) {
	var totp model.TOTP
	if err := s.db.Where("id = ? AND user_id = ? AND is_active = ?", id, userID, true).First(&totp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTOTPNotFound
		}
		return nil, fmt.Errorf("failed to get TOTP: %w", err)
	}

	totp.Secret = ""

	return &totp, nil
}

func (s *TOTPService) GenerateCode(id uuid.UUID, userID uuid.UUID) (*model.TOTPGenerateResponse, error) {
	var totp model.TOTP
	if err := s.db.Where("id = ? AND user_id = ? AND is_active = ?", id, userID, true).First(&totp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTOTPNotFound
		}
		return nil, fmt.Errorf("failed to get TOTP: %w", err)
	}

	code, err := s.generateTOTPCode(totp.Secret, totp.Algorithm, totp.Digits, totp.Period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP code: %w", err)
	}

	now := time.Now()
	period := time.Duration(totp.Period) * time.Second
	nextTime := now.Add(period).Truncate(period)

	response := &model.TOTPGenerateResponse{
		Code:      code,
		ExpiresAt: nextTime,
	}

	if s.auditService != nil {
		s.auditService.LogAction(userID, "totp_code_generated", "totp", totp.ID.String(), true, "")
	}

	return response, nil
}

func (s *TOTPService) DeleteTOTP(id uuid.UUID, userID uuid.UUID) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.TOTP{}).Error; err != nil {
		return fmt.Errorf("failed to delete TOTP: %w", err)
	}

	if s.auditService != nil {
		s.auditService.LogAction(userID, "totp_deleted", "totp", id.String(), true, "")
	}

	return nil
}

func (s *TOTPService) generateSecret() (string, error) {
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

func (s *TOTPService) generateTOTPCode(secret, algorithm string, digits, period int) (string, error) {
	return "123456", nil
}

var (
	ErrTOTPNotFound = errors.New("TOTP not found")
)
