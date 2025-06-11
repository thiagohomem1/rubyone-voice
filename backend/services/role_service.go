package services

import (
	"errors"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
)

type RoleService struct {
	DB *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{DB: db}
}

func (s *RoleService) CreateRole(tenantID uint, name string) (*models.Role, error) {
	role := models.Role{
		TenantID: tenantID,
		Name:     name,
	}

	if err := s.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) GetRolesByTenant(tenantID uint) ([]models.Role, error) {
	var roles []models.Role
	
	if err := s.DB.Where("tenant_id = ?", tenantID).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) GetRoleByID(tenantID, roleID uint) (*models.Role, error) {
	var role models.Role
	
	if err := s.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).
		Preload("Permissions").
		First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) DeleteRole(tenantID, roleID uint) error {
	var role models.Role
	
	if err := s.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).
		First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	// Check if role is assigned to users
	var userCount int64
	if err := s.DB.Model(&models.User{}).
		Where("role_id = ? AND tenant_id = ?", roleID, tenantID).
		Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("cannot delete role that is assigned to users")
	}

	// Delete role permissions first
	if err := s.DB.Where("role_id = ?", roleID).
		Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	// Delete the role
	if err := s.DB.Delete(&role).Error; err != nil {
		return err
	}

	return nil
}

func (s *RoleService) AssignPermissionsToRole(tenantID, roleID uint, permissionIDs []uint) error {
	// Verify role belongs to tenant
	var role models.Role
	if err := s.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).
		First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	// Verify all permissions exist
	var permissionCount int64
	if err := s.DB.Model(&models.Permission{}).
		Where("id IN ?", permissionIDs).
		Count(&permissionCount).Error; err != nil {
		return err
	}

	if int(permissionCount) != len(permissionIDs) {
		return errors.New("one or more permissions not found")
	}

	// Remove existing permissions for this role
	if err := s.DB.Where("role_id = ?", roleID).
		Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	// Assign new permissions
	for _, permissionID := range permissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		
		if err := s.DB.Create(&rolePermission).Error; err != nil {
			return err
		}
	}

	return nil
} 