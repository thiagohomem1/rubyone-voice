package services

import (
	"saas-backend/models"
	"gorm.io/gorm"
)

type UsageService struct {
	db *gorm.DB
}

func NewUsageService(db *gorm.DB) *UsageService {
	return &UsageService{db: db}
}

func (s *UsageService) GetTenantUsage(tenantID uint) (*models.UsageReport, error) {
	var activeUsers int64
	var totalCalls int64
	var subscription models.Subscription
	var plan models.Plan

	// Count active users for the tenant
	err := s.db.Model(&models.UserTenant{}).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Count(&activeUsers).Error
	if err != nil {
		return nil, err
	}

	// Count total calls for the tenant
	err = s.db.Model(&models.Call{}).
		Where("tenant_id = ?", tenantID).
		Count(&totalCalls).Error
	if err != nil {
		return nil, err
	}

	// Get current active subscription
	err = s.db.Where("tenant_id = ? AND status = ?", tenantID, "active").
		Order("created_at DESC").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}

	// Get plan details
	err = s.db.Where("id = ?", subscription.PlanID).First(&plan).Error
	if err != nil {
		return nil, err
	}

	// Calculate remaining quotas
	usersRemaining := uint(0)
	if plan.MaxUsers > uint(activeUsers) {
		usersRemaining = plan.MaxUsers - uint(activeUsers)
	}

	callsRemaining := uint(0)
	if plan.MaxCalls > uint(totalCalls) {
		callsRemaining = plan.MaxCalls - uint(totalCalls)
	}

	usageReport := &models.UsageReport{
		TenantID:       tenantID,
		ActiveUsers:    uint(activeUsers),
		TotalCalls:     uint(totalCalls),
		MaxUsers:       plan.MaxUsers,
		MaxCalls:       plan.MaxCalls,
		UsersRemaining: usersRemaining,
		CallsRemaining: callsRemaining,
	}

	return usageReport, nil
} 