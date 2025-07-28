package tenants

// MockTenantRepository para testes
type MockTenantRepository struct {
	CreateFunc func(tenant *Tenant) error
	FindByIDFunc func(id string) (*Tenant, error)
}

func (m *MockTenantRepository) Create(tenant *Tenant) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(tenant)
	}
	return nil
}

func (m *MockTenantRepository) FindByID(id string) (*Tenant, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
} 