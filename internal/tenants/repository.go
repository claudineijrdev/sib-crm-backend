package tenants

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/claudineijrdev/sib-crm-backend/internal/platform/cache"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/telemetry"
	"gorm.io/gorm"
)

// Repository base (sem cache/telemetria)
type tenantRepositoryBase struct {
	db *gorm.DB
}

func newTenantRepositoryBase(db *gorm.DB) *tenantRepositoryBase {
	return &tenantRepositoryBase{db: db}
}

func (r *tenantRepositoryBase) create(tenant *Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepositoryBase) findByID(id string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// Repository com cache e telemetria (decorator)
type tenantRepository struct {
	base      *tenantRepositoryBase
	cache     cache.CacheService
	telemetry telemetry.TelemetryService
}

func NewTenantRepository(db *gorm.DB, cache cache.CacheService, telemetry telemetry.TelemetryService) TenantRepository {
	return &tenantRepository{
		base:      newTenantRepositoryBase(db),
		cache:     cache,
		telemetry: telemetry,
	}
}

func (r *tenantRepository) Create(tenant *Tenant) error {
	ctx := context.Background()
	span, ctx := r.telemetry.StartSpan(ctx, "repository.tenant.create")
	defer span.End()

	span.SetTag("tenant_id", tenant.ID.String())
	span.SetTag("tenant_name", tenant.Name)

	err := r.base.create(tenant)
	if err != nil {
		span.SetError(err)
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("tenant:id:%s", tenant.ID.String())
	r.cache.Delete(ctx, cacheKey)

	// Track metric
	r.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "repository.tenant.create.success",
		Value: 1,
		Tags:  map[string]string{"tenant_name": tenant.Name},
	})

	return nil
}

func (r *tenantRepository) FindByID(id string) (*Tenant, error) {
	ctx := context.Background()
	span, ctx := r.telemetry.StartSpan(ctx, "repository.tenant.find_by_id")
	defer span.End()

	span.SetTag("tenant_id", id)

	// Try cache first
	cacheKey := fmt.Sprintf("tenant:id:%s", id)
	if cached, err := r.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		span.SetTag("cache_hit", "true")
		if tenant, ok := cached.(*Tenant); ok {
			r.telemetry.TrackMetric(ctx, telemetry.Metric{
				Name:  "repository.tenant.find_by_id.cache_hit",
				Value: 1,
				Tags:  map[string]string{"tenant_id": id},
			})
			return tenant, nil
		}
	}

	// Cache miss - query database
	span.SetTag("cache_hit", "false")
	tenant, err := r.base.findByID(id)
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	// Cache the result (even if nil to avoid cache penetration)
	ttl := 10 * time.Minute
	if tenant == nil {
		ttl = 5 * time.Minute // Shorter TTL for non-existent tenants
	}
	r.cache.Set(ctx, cacheKey, tenant, ttl)

	r.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "repository.tenant.find_by_id.cache_miss",
		Value: 1,
		Tags:  map[string]string{"tenant_id": id},
	})

	return tenant, nil
} 