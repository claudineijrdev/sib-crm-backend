package auth

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
type userRepositoryBase struct {
	db *gorm.DB
}

func newUserRepositoryBase(db *gorm.DB) *userRepositoryBase {
	return &userRepositoryBase{db: db}
}

func (r *userRepositoryBase) create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryBase) findByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryBase) findByID(id string) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Repository com cache e telemetria (decorator)
type userRepository struct {
	base      *userRepositoryBase
	cache     cache.CacheService
	telemetry telemetry.TelemetryService
}

func NewUserRepository(db *gorm.DB, cache cache.CacheService, telemetry telemetry.TelemetryService) UserRepository {
	return &userRepository{
		base:      newUserRepositoryBase(db),
		cache:     cache,
		telemetry: telemetry,
	}
}

func (r *userRepository) Create(user *User) error {
	ctx := context.Background()
	span, ctx := r.telemetry.StartSpan(ctx, "repository.user.create")
	defer span.End()

	span.SetTag("user_id", user.ID.String())
	span.SetTag("email", user.Email)

	err := r.base.create(user)
	if err != nil {
		span.SetError(err)
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:email:%s", user.Email)
	r.cache.Delete(ctx, cacheKey)

	// Track metric
	r.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "repository.user.create.success",
		Value: 1,
		Tags:  map[string]string{"email": user.Email},
	})

	return nil
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	ctx := context.Background()
	span, ctx := r.telemetry.StartSpan(ctx, "repository.user.find_by_email")
	defer span.End()

	span.SetTag("email", email)

	// Try cache first
	cacheKey := fmt.Sprintf("user:email:%s", email)
	if cached, err := r.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		span.SetTag("cache_hit", "true")
		if user, ok := cached.(*User); ok {
			r.telemetry.TrackMetric(ctx, telemetry.Metric{
				Name:  "repository.user.find_by_email.cache_hit",
				Value: 1,
				Tags:  map[string]string{"email": email},
			})
			return user, nil
		}
	}

	// Cache miss - query database
	span.SetTag("cache_hit", "false")
	user, err := r.base.findByEmail(email)
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	// Cache the result (even if nil to avoid cache penetration)
	ttl := 10 * time.Minute
	if user == nil {
		ttl = 5 * time.Minute // Shorter TTL for non-existent users
	}
	r.cache.Set(ctx, cacheKey, user, ttl)

	r.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "repository.user.find_by_email.cache_miss",
		Value: 1,
		Tags:  map[string]string{"email": email},
	})

	return user, nil
}

func (r *userRepository) FindByID(id string) (*User, error) {
	ctx := context.Background()
	span, ctx := r.telemetry.StartSpan(ctx, "repository.user.find_by_id")
	defer span.End()

	span.SetTag("user_id", id)

	// Try cache first
	cacheKey := fmt.Sprintf("user:id:%s", id)
	if cached, err := r.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		span.SetTag("cache_hit", "true")
		if user, ok := cached.(*User); ok {
			r.telemetry.TrackMetric(ctx, telemetry.Metric{
				Name:  "repository.user.find_by_id.cache_hit",
				Value: 1,
				Tags:  map[string]string{"user_id": id},
			})
			return user, nil
		}
	}

	// Cache miss - query database
	span.SetTag("cache_hit", "false")
	user, err := r.base.findByID(id)
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	// Cache the result
	if user != nil {
		r.cache.Set(ctx, cacheKey, user, 10*time.Minute)
	}

	r.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "repository.user.find_by_id.cache_miss",
		Value: 1,
		Tags:  map[string]string{"user_id": id},
	})

	return user, nil
} 