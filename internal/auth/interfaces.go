package auth

import (
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
)

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}

type AuthService interface {
	RegisterUser(req RegisterRequest) (*User, *tenants.Tenant, error)
	LoginUser(req LoginRequest) (*User, error)
} 