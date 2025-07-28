package tenants

type TenantRepository interface {
	Create(tenant *Tenant) error
	FindByID(id string) (*Tenant, error)
} 