package container

import (
	"github.com/claudineijrdev/sib-crm-backend/internal/auth"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/cache"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/telemetry"
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
	"gorm.io/gorm"
)

type Container struct {
	// Infraestrutura compartilhada
	DB             *gorm.DB
	Telemetry      telemetry.TelemetryService
	Cache          cache.CacheService
	
	// Repositórios
	UserRepo       auth.UserRepository
	TenantRepo     tenants.TenantRepository
	
	// Services
	AuthService    auth.AuthService
	
	// Handlers
	AuthHandler    *auth.AuthHandler
}

func NewContainer(db *gorm.DB) *Container {
	// Inicializar serviços de infraestrutura
	telemetryService := telemetry.NewTelemetryService(true) // enabled
	cacheService := cache.NewCacheService(nil) // nil client por enquanto
	
	// Criar repositórios (com cache e telemetria)
	userRepo := auth.NewUserRepository(db, cacheService, telemetryService)
	tenantRepo := tenants.NewTenantRepository(db, cacheService, telemetryService)
	
	// Criar services (só telemetria, cache já está no repository)
	authService := auth.NewAuthService(userRepo, tenantRepo, db, telemetryService)
	
	// Criar handlers
	authHandler := auth.NewAuthHandler(authService)
	
	return &Container{
		// Infraestrutura
		DB:        db,
		Telemetry: telemetryService,
		Cache:     cacheService,
		
		// Repositórios
		UserRepo:   userRepo,
		TenantRepo: tenantRepo,
		
		// Services
		AuthService: authService,
		
		// Handlers
		AuthHandler: authHandler,
	}
} 