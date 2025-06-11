package services

import (
	"errors"
	"saas-backend/models"
	"saas-backend/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(tenantID uint, roleID uint, email, password, firstName, lastName string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	tx := s.db.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	userTenant := &models.UserTenant{
		UserID:   user.ID,
		TenantID: tenantID,
	}

	if err := tx.Create(userTenant).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	userRole := &models.UserRole{
		UserID: user.ID,
		RoleID: roleID,
	}

	if err := tx.Create(userRole).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *UserService) GetAllUsers(tenantID uint) ([]models.User, error) {
	var users []models.User

	err := s.db.Table("users").
		Select("users.id, users.email, users.first_name, users.last_name, users.is_active, users.created_at, users.updated_at").
		Joins("JOIN user_tenants ON user_tenants.user_id = users.id").
		Where("user_tenants.tenant_id = ?", tenantID).
		Find(&users).Error

	return users, err
}

func (s *UserService) GetUserByID(tenantID, userID uint) (*models.User, error) {
	var user models.User

	err := s.db.Table("users").
		Select("users.id, users.email, users.first_name, users.last_name, users.is_active, users.created_at, users.updated_at").
		Joins("JOIN user_tenants ON user_tenants.user_id = users.id").
		Where("user_tenants.tenant_id = ? AND users.id = ?", tenantID, userID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) DeleteUser(tenantID, userID uint) error {
	var count int64

	err := s.db.Table("users").
		Joins("JOIN user_tenants ON user_tenants.user_id = users.id").
		Where("user_tenants.tenant_id = ? AND users.id = ?", tenantID, userID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user not found")
	}

	tx := s.db.Begin()

	if err := tx.Where("user_id = ?", userID).Delete(&models.UserTenant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.User{}, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
} 