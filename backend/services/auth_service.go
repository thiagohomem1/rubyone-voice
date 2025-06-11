// backend/services/auth_service.go
package services

import (
	"errors"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
	"github.com/your-module/backend/utils"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) RegisterTenant(tenantName, domain, username, password string) (*models.User, string, error) {
	var existingTenant models.Tenant
	if err := s.DB.Where("domain = ?", domain).First(&existingTenant).Error; err == nil {
		return nil, "", errors.New("tenant domain already exists")
	}

	var existingUser models.User
	if err := s.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, "", errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	tx := s.DB.Begin()

	tenant := models.Tenant{
		Name:   tenantName,
		Domain: domain,
	}
	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	adminRole := models.Role{
		TenantID: tenant.ID,
		Name:     "Admin",
	}
	if err := tx.Create(&adminRole).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	user := models.User{
		TenantID:     tenant.ID,
		Username:     username,
		PasswordHash: hashedPassword,
		RoleID:       adminRole.ID,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	userTenant := models.UserTenant{
		UserID:   user.ID,
		TenantID: tenant.ID,
		IsActive: true,
	}
	if err := tx.Create(&userTenant).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	userRole := models.UserRole{
		UserID:   user.ID,
		RoleID:   adminRole.ID,
		TenantID: tenant.ID,
		IsActive: true,
	}
	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.RoleID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) RegisterUser(tenantID uint, username, password string, roleID uint) (*models.User, string, error) {
	var existingUser models.User
	if err := s.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, "", errors.New("username already exists")
	}

	var tenant models.Tenant
	if err := s.DB.First(&tenant, tenantID).Error; err != nil {
		return nil, "", errors.New("tenant not found")
	}

	var role models.Role
	if err := s.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
		return nil, "", errors.New("role not found for this tenant")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	tx := s.DB.Begin()

	user := models.User{
		TenantID:     tenantID,
		Username:     username,
		PasswordHash: hashedPassword,
		RoleID:       roleID,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	userTenant := models.UserTenant{
		UserID:   user.ID,
		TenantID: tenantID,
		IsActive: true,
	}
	if err := tx.Create(&userTenant).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	userRole := models.UserRole{
		UserID:   user.ID,
		RoleID:   roleID,
		TenantID: tenantID,
		IsActive: true,
	}
	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.RoleID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	var user models.User
	if err := s.DB.Preload("Tenant").Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	var userTenant models.UserTenant
	if err := s.DB.Where("user_id = ? AND tenant_id = ? AND is_active = ?", user.ID, user.TenantID, true).First(&userTenant).Error; err != nil {
		return nil, "", errors.New("user not active in tenant")
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.RoleID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
} 