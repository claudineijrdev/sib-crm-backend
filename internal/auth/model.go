package auth

import (
	"time"

	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
	"github.com/google/uuid"
)

// User represents the GORM model for a user.
type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null" json:"tenant_id"`
	Tenant        tenants.Tenant `gorm:"foreignKey:TenantID"`
	Email         string         `gorm:"type:varchar(255);not null;unique" json:"email"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
