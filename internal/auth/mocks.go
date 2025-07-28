package auth

import (
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
)

// MockUserRepository para testes
type MockUserRepository struct {
	CreateFunc      func(user *User) error
	FindByEmailFunc func(email string) (*User, error)
	FindByIDFunc    func(id string) (*User, error)
}

func (m *MockUserRepository) Create(user *User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(id string) (*User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// MockTenantRepository para testes
type MockTenantRepository struct {
	CreateFunc func(tenant *tenants.Tenant) error
	FindByIDFunc func(id string) (*tenants.Tenant, error)
}

func (m *MockTenantRepository) Create(tenant *tenants.Tenant) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(tenant)
	}
	return nil
}

func (m *MockTenantRepository) FindByID(id string) (*tenants.Tenant, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

// MockAuthService para testes
type MockAuthService struct {
	RegisterUserFunc func(req RegisterRequest) (*User, *tenants.Tenant, error)
	LoginUserFunc    func(req LoginRequest) (*User, error)
}

func (m *MockAuthService) RegisterUser(req RegisterRequest) (*User, *tenants.Tenant, error) {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(req)
	}
	return nil, nil, nil
}

func (m *MockAuthService) LoginUser(req LoginRequest) (*User, error) {
	if m.LoginUserFunc != nil {
		return m.LoginUserFunc(req)
	}
	return nil, nil
} 