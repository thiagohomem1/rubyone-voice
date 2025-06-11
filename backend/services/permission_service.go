// backend/services/permission_service.go
package services

import (
	"errors"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
)

type PermissionService struct {
	DB *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{DB: db}
}

func (s *PermissionService) CreatePermission(code, description string) (*models.Permission, error) {
	permission := models.Permission{
		Code:        code,
		Description: description,
	}

	if err := s.DB.Create(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	
	if err := s.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (s *PermissionService) GetPermissionByID(permissionID uint) (*models.Permission, error) {
	var permission models.Permission
	
	if err := s.DB.First(&permission, permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

func (s *PermissionService) DeletePermission(permissionID uint) error {
	var permission models.Permission
	
	if err := s.DB.First(&permission, permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return err
	}

	// Check if permission is assigned to roles
	var rolePermissionCount int64
	if err := s.DB.Model(&models.RolePermission{}).
		Where("permission_id = ?", permissionID).
		Count(&rolePermissionCount).Error; err != nil {
		return err
	}

	if rolePermissionCount > 0 {
		return errors.New("cannot delete permission that is assigned to roles")
	}

	// Delete the permission
	if err := s.DB.Delete(&permission).Error; err != nil {
		return err
	}

	return nil
}

func (s *PermissionService) AssignPermissionToRole(roleID, permissionID uint) error {
	var role models.Role
	if err := s.DB.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	var permission models.Permission
	if err := s.DB.First(&permission, permissionID).Error; err != nil {
		return errors.New("permission not found")
	}

	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	if err := s.DB.Create(&rolePermission).Error; err != nil {
		return err
	}

	return nil
}

func (s *PermissionService) HasPermission(userID uint, permissionCode string) (bool, error) {
	var user models.User
	if err := s.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	for _, permission := range user.Role.Permissions {
		if permission.Code == permissionCode {
			return true, nil
		}
	}

	return false, nil
}
