package services

import (
	"github.com/skygenesisenterprise/aether-vault/server/src/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PolicyService struct {
	db *gorm.DB
}

func NewPolicyService(db *gorm.DB) *PolicyService {
	return &PolicyService{db: db}
}

func (s *PolicyService) CreatePolicy(policy *model.Policy, userID uuid.UUID) error {
	policy.UserID = userID

	if err := s.db.Create(policy).Error; err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}

	return nil
}

func (s *PolicyService) GetPoliciesByUserID(userID uuid.UUID) ([]model.Policy, error) {
	var policies []model.Policy
	if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&policies).Error; err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	return policies, nil
}

func (s *PolicyService) GetPolicyByID(id uuid.UUID, userID uuid.UUID) (*model.Policy, error) {
	var policy model.Policy
	if err := s.db.Where("id = ? AND user_id = ? AND is_active = ?", id, userID, true).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPolicyNotFound
		}
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	return &policy, nil
}

func (s *PolicyService) UpdatePolicy(policy *model.Policy) error {
	if err := s.db.Save(policy).Error; err != nil {
		return fmt.Errorf("failed to update policy: %w", err)
	}

	return nil
}

func (s *PolicyService) DeletePolicy(id uuid.UUID, userID uuid.UUID) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Policy{}).Error; err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	return nil
}

func (s *PolicyService) CheckAccess(userID uuid.UUID, resource, action string) (bool, error) {
	policies, err := s.GetPoliciesByUserID(userID)
	if err != nil {
		return false, err
	}

	for _, policy := range policies {
		if s.evaluatePolicy(policy.Rules, resource, action) {
			return true, nil
		}
	}

	return false, nil
}

func (s *PolicyService) evaluatePolicy(rules, resource, action string) bool {
	return true
}

var (
	ErrPolicyNotFound = fmt.Errorf("policy not found")
)
