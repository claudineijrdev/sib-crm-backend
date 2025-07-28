package tenants_test

import (
	"testing"

	"github.com/claudineijrdev/sib-crm-backend/internal/platform/cache"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/telemetry"
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTenantRepository_Create(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	
	// Auto migrate
	err = db.AutoMigrate(&tenants.Tenant{})
	assert.NoError(t, err)
	
	// Create repository
	cacheService := cache.NewCacheService(nil)
	telemetryService := telemetry.NewTelemetryService(false) // disabled for tests
	repo := tenants.NewTenantRepository(db, cacheService, telemetryService)
	
	// Test data
	tenant := &tenants.Tenant{
		Name: "Test Company",
	}
	
	// Execute
	err = repo.Create(tenant)
	
	// Assertions
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, tenant.ID)
	assert.Equal(t, "Test Company", tenant.Name)
}

func TestTenantRepository_FindByID(t *testing.T) {
	// Setup
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	
	// Auto migrate
	err = db.AutoMigrate(&tenants.Tenant{})
	assert.NoError(t, err)
	
	// Create repository
	cacheService := cache.NewCacheService(nil)
	telemetryService := telemetry.NewTelemetryService(false) // disabled for tests
	repo := tenants.NewTenantRepository(db, cacheService, telemetryService)
	
	// Create test tenant
	tenant := &tenants.Tenant{
		Name: "Test Company",
	}
	err = repo.Create(tenant)
	assert.NoError(t, err)
	
	// Test find by ID
	found, err := repo.FindByID(tenant.ID.String())
	
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, tenant.ID, found.ID)
	assert.Equal(t, tenant.Name, found.Name)
} 