package services

import (
	"aether-vault/src/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) LogAction(userID uuid.UUID, action, resource, resourceID string, success bool, details string) error {
	auditLog := &model.AuditLog{
		UserID:     &userID,
		Action:     action,
		Resource:   resource,
		ResourceID: &resourceID,
		Success:    success,
		Details:    details,
		CreatedAt:  time.Now(),
	}

	if err := s.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (s *AuditService) LogAnonymousAction(action, resource, resourceID, ipAddress, userAgent string, success bool, details string) error {
	auditLog := &model.AuditLog{
		UserID:     nil,
		Action:     action,
		Resource:   resource,
		ResourceID: &resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Success:    success,
		Details:    details,
		CreatedAt:  time.Now(),
	}

	if err := s.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (s *AuditService) GetAuditLogs(userID *uuid.UUID, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	query := s.db.Order("created_at DESC")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

func (s *AuditService) GetAuditLogsByResource(resource, resourceID string, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	if err := s.db.Where("resource = ? AND resource_id = ?", resource, resourceID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

func (s *AuditService) CleanupOldLogs(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	if err := s.db.Where("created_at < ?", cutoffDate).Delete(&model.AuditLog{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup old audit logs: %w", err)
	}

	return nil
}
