package services

import (
	"errors"
	"time"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
)

type SubscriptionService struct {
	DB *gorm.DB
}

func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{DB: db}
}

func (s *SubscriptionService) CreatePlan(name string, maxUsers uint, maxCalls uint, price float64) (*models.Plan, error) {
	plan := models.Plan{
		Name:     name,
		MaxUsers: maxUsers,
		MaxCalls: maxCalls,
		Price:    price,
	}

	if err := s.DB.Create(&plan).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *SubscriptionService) GetAllPlans() ([]models.Plan, error) {
	var plans []models.Plan
	
	if err := s.DB.Find(&plans).Error; err != nil {
		return nil, err
	}

	return plans, nil
}

func (s *SubscriptionService) GetPlanByID(planID uint) (*models.Plan, error) {
	var plan models.Plan
	
	if err := s.DB.First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plan not found")
		}
		return nil, err
	}

	return &plan, nil
}

func (s *SubscriptionService) DeletePlan(planID uint) error {
	var plan models.Plan
	
	if err := s.DB.First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("plan not found")
		}
		return err
	}

	// Check if plan has active subscriptions
	var subscriptionCount int64
	if err := s.DB.Model(&models.Subscription{}).
		Where("plan_id = ? AND is_active = ?", planID, true).
		Count(&subscriptionCount).Error; err != nil {
		return err
	}

	if subscriptionCount > 0 {
		return errors.New("cannot delete plan with active subscriptions")
	}

	// Delete the plan
	if err := s.DB.Delete(&plan).Error; err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) SubscribeTenant(tenantID uint, planID uint) error {
	// Verify tenant exists
	var tenant models.Tenant
	if err := s.DB.First(&tenant, tenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tenant not found")
		}
		return err
	}

	// Verify plan exists
	var plan models.Plan
	if err := s.DB.First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("plan not found")
		}
		return err
	}

	// End current active subscription if exists
	var currentSubscription models.Subscription
	if err := s.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		First(&currentSubscription).Error; err == nil {
		// Subscription exists, end it
		now := time.Now()
		currentSubscription.IsActive = false
		currentSubscription.EndedAt = &now
		
		if err := s.DB.Save(&currentSubscription).Error; err != nil {
			return err
		}
	}

	// Create new subscription
	subscription := models.Subscription{
		TenantID:  tenantID,
		PlanID:    planID,
		IsActive:  true,
		StartedAt: time.Now(),
	}

	if err := s.DB.Create(&subscription).Error; err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) GetTenantSubscription(tenantID uint) (*models.Subscription, error) {
	var subscription models.Subscription
	
	if err := s.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Preload("Plan").
		First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}

	return &subscription, nil
} 