package models

type UsageReport struct {
	TenantID       uint `json:"tenant_id"`
	ActiveUsers    uint `json:"active_users"`
	TotalCalls     uint `json:"total_calls"`
	MaxUsers       uint `json:"max_users"`
	MaxCalls       uint `json:"max_calls"`
	UsersRemaining uint `json:"users_remaining"`
	CallsRemaining uint `json:"calls_remaining"`
} 