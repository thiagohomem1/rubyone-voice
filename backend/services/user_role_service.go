package services

import (
	"errors"
	"time"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
)

type UserRoleService struct {
	DB *gorm.DB
}

func NewUserRoleService(db *gorm.DB) *UserRoleService {
	return &UserRoleService{DB: db}
}

func (s *UserRoleService) AssignRolesToUser(tenantID, userID uint, roleIDs []uint) error {
	// Verify user belongs to tenant
	var user models.User
	if err := s.DB.Where("id = ? AND tenant_id = ?", userID, tenantID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify all roles belong to the same tenant
	var roleCount int64
	if err := s.DB.Model(&models.Role{}).
		Where("id IN ? AND tenant_id = ?", roleIDs, tenantID).
		Count(&roleCount).Error; err != nil {
		return err
	}

	if int(roleCount) != len(roleIDs) {
		return errors.New("one or more roles not found or don't belong to tenant")
	}

	// Remove existing active user roles for this tenant
	if err := s.DB.Model(&models.UserRole{}).
		Where("user_id = ? AND tenant_id = ?", userID, tenantID).
		Updates(map[string]interface{}{
			"is_active":  false,
			"revoked_at": time.Now(),
		}).Error; err != nil {
		return err
	}

	// Assign new roles
	for _, roleID := range roleIDs {
		userRole := models.UserRole{
			UserID:     userID,
			RoleID:     roleID,
			TenantID:   tenantID,
			IsActive:   true,
			AssignedAt: time.Now(),
		}
		
		if err := s.DB.Create(&userRole).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *UserRoleService) GetUserRoles(tenantID, userID uint) ([]models.Role, error) {
	// Verify user belongs to tenant
	var user models.User
	if err := s.DB.Where("id = ? AND tenant_id = ?", userID, tenantID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var roles []models.Role
	
	if err := s.DB.Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.tenant_id = ? AND user_roles.is_active = ?", 
			userID, tenantID, true).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *UserRoleService) RemoveRoleFromUser(tenantID, userID, roleID uint) error {
	// Verify user belongs to tenant
	var user models.User
	if err := s.DB.Where("id = ? AND tenant_id = ?", userID, tenantID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify role belongs to tenant
	var role models.Role
	if err := s.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).
		First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	// Find and deactivate the user role
	var userRole models.UserRole
	if err := s.DB.Where("user_id = ? AND role_id = ? AND tenant_id = ? AND is_active = ?", 
		userID, roleID, tenantID, true).
		First(&userRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user role assignment not found")
		}
		return err
	}

	// Update the user role to inactive
	now := time.Now()
	userRole.IsActive = false
	userRole.RevokedAt = &now

	if err := s.DB.Save(&userRole).Error; err != nil {
		return err
	}

	return nil
} 