package services

import (
	"errors"
	"gorm.io/gorm"
	"github.com/your-module/backend/models"
)

type CallService struct {
	DB *gorm.DB
}

func NewCallService(db *gorm.DB) *CallService {
	return &CallService{DB: db}
}

func (s *CallService) CreateCall(tenantID uint, caller string, callee string, duration uint, recordingURL string) (*models.Call, error) {
	call := models.Call{
		TenantID:     tenantID,
		Caller:       caller,
		Callee:       callee,
		Billsec:      int(duration),
		RecordingURL: recordingURL,
	}

	if err := s.DB.Create(&call).Error; err != nil {
		return nil, err
	}

	return &call, nil
}

func (s *CallService) GetAllCalls(tenantID uint) ([]models.Call, error) {
	var calls []models.Call
	
	if err := s.DB.Where("tenant_id = ?", tenantID).
		Find(&calls).Error; err != nil {
		return nil, err
	}

	return calls, nil
}

func (s *CallService) GetCallByID(tenantID uint, callID uint) (*models.Call, error) {
	var call models.Call
	
	if err := s.DB.Where("id = ? AND tenant_id = ?", callID, tenantID).
		First(&call).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("call not found")
		}
		return nil, err
	}

	return &call, nil
}

func (s *CallService) DeleteCall(tenantID uint, callID uint) error {
	var call models.Call
	
	if err := s.DB.Where("id = ? AND tenant_id = ?", callID, tenantID).
		First(&call).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("call not found")
		}
		return err
	}

	if err := s.DB.Delete(&call).Error; err != nil {
		return err
	}

	return nil
} 