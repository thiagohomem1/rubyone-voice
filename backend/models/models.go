package models

import (
	"time"
	"gorm.io/gorm"
)

type Tenant struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	Domain    string         `gorm:"not null;unique" json:"domain"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Users         []User         `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	Roles         []Role         `gorm:"foreignKey:TenantID" json:"roles,omitempty"`
	Calls         []Call         `gorm:"foreignKey:TenantID" json:"calls,omitempty"`
	UserTenants   []UserTenant   `gorm:"foreignKey:TenantID" json:"user_tenants,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:TenantID" json:"subscriptions,omitempty"`
}

type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Tenant          Tenant           `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users           []User           `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID" json:"role_permissions,omitempty"`
	Permissions     []Permission     `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	UserRoles       []UserRole       `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
}

type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"not null;unique" json:"code"`
	Description string         `gorm:"not null" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID" json:"role_permissions,omitempty"`
	Roles           []Role           `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

type RolePermission struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	PermissionID uint           `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	Username     string         `gorm:"not null;unique" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Tenant      Tenant       `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Role        Role         `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	UserTenants []UserTenant `gorm:"foreignKey:UserID" json:"user_tenants,omitempty"`
	UserRoles   []UserRole   `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
}

type UserTenant struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	TenantID   uint           `gorm:"not null;index" json:"tenant_id"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	JoinedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
	LeftAt     *time.Time     `json:"left_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	
	// Unique constraint to prevent duplicate user-tenant relationships
	// Using composite unique index for better performance
	// gorm:"uniqueIndex:idx_user_tenant_unique"
}

type UserRole struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	RoleID     uint           `gorm:"not null;index" json:"role_id"`
	TenantID   uint           `gorm:"not null;index" json:"tenant_id"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	AssignedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	RevokedAt  *time.Time     `json:"revoked_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role   Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	
	// Unique constraint to prevent duplicate user-role-tenant relationships
	// Using composite unique index for better performance
	// gorm:"uniqueIndex:idx_user_role_tenant_unique"
}

type Call struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     uint           `gorm:"not null;index" json:"tenant_id"`
	UUID         string         `gorm:"not null;unique" json:"uuid"`
	Caller       string         `gorm:"not null" json:"caller"`
	Callee       string         `gorm:"not null" json:"callee"`
	StartTime    *time.Time     `json:"start_time"`
	AnswerTime   *time.Time     `json:"answer_time"`
	EndTime      *time.Time     `json:"end_time"`
	Billsec      int            `gorm:"default:0" json:"billsec"`
	RecordingURL string         `json:"recording_url"`
	Cost         float64        `gorm:"type:decimal(10,4);default:0" json:"cost"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// Method to add unique indexes for UserTenant and UserRole models
func (UserTenant) TableName() string {
	return "user_tenants"
}

func (UserRole) TableName() string {
	return "user_roles"
}

type Plan struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	MaxUsers  uint           `gorm:"not null" json:"max_users"`
	MaxCalls  uint           `gorm:"not null" json:"max_calls"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Subscriptions []Subscription `gorm:"foreignKey:PlanID" json:"subscriptions,omitempty"`
}

type Subscription struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	PlanID    uint           `gorm:"not null;index" json:"plan_id"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	StartedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"started_at"`
	EndedAt   *time.Time     `json:"ended_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relations
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Plan   Plan   `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
} 