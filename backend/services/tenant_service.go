package services

import (
	"errors"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
)

type TenantService struct {
	DB *gorm.DB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{DB: db}
}

func (s *TenantService) CreateTenant(name, domain string) (*models.Tenant, error) {
	tenant := models.Tenant{
		Name:   name,
		Domain: domain,
	}

	if err := s.DB.Create(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (s *TenantService) GetAllTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	
	if err := s.DB.Find(&tenants).Error; err != nil {
		return nil, err
	}

	return tenants, nil
}

func (s *TenantService) GetTenantByID(tenantID uint) (*models.Tenant, error) {
	var tenant models.Tenant
	
	if err := s.DB.First(&tenant, tenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	return &tenant, nil
}

func (s *TenantService) DeleteTenant(tenantID uint) error {
	var tenant models.Tenant
	
	if err := s.DB.First(&tenant, tenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tenant not found")
		}
		return err
	}

	// Check if tenant has associated users
	var userCount int64
	if err := s.DB.Model(&models.User{}).
		Where("tenant_id = ?", tenantID).
		Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("cannot delete tenant that has associated users")
	}

	// Check if tenant has associated roles
	var roleCount int64
	if err := s.DB.Model(&models.Role{}).
		Where("tenant_id = ?", tenantID).
		Count(&roleCount).Error; err != nil {
		return err
	}

	if roleCount > 0 {
		return errors.New("cannot delete tenant that has associated roles")
	}

	// Check if tenant has associated calls
	var callCount int64
	if err := s.DB.Model(&models.Call{}).
		Where("tenant_id = ?", tenantID).
		Count(&callCount).Error; err != nil {
		return err
	}

	if callCount > 0 {
		return errors.New("cannot delete tenant that has associated calls")
	}

	// Delete the tenant
	if err := s.DB.Delete(&tenant).Error; err != nil {
		return err
	}

	return nil
} 