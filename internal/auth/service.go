package auth

import (
	"context"
	"errors"
	"time"

	"github.com/claudineijrdev/sib-crm-backend/internal/platform/telemetry"
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	userRepo   UserRepository
	tenantRepo tenants.TenantRepository
	db         *gorm.DB
	telemetry  telemetry.TelemetryService
}

func NewAuthService(
	userRepo UserRepository, 
	tenantRepo tenants.TenantRepository, 
	db *gorm.DB,
	telemetry telemetry.TelemetryService,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		db:         db,
		telemetry:  telemetry,
	}
}

func (s *authService) RegisterUser(req RegisterRequest) (*User, *tenants.Tenant, error) {
	ctx := context.Background()
	span, ctx := s.telemetry.StartSpan(ctx, "auth.register_user")
	defer span.End()

	// Track metric
	s.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "auth.register.attempt",
		Value: 1,
		Tags:  map[string]string{"email": req.Email},
	})

	// Verificar se email já existe (cache já está no repository)
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		span.SetError(err)
		return nil, nil, err
	}
	if existingUser != nil {
		span.SetTag("user_exists", "true")
		return nil, nil, errors.New("email already exists")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		span.SetError(err)
		return nil, nil, err
	}

	// Criar tenant usando o repository
	tenant := &tenants.Tenant{Name: req.Name}
	err = s.tenantRepo.Create(tenant)
	if err != nil {
		span.SetError(err)
		return nil, nil, err
	}

	// Criar user
	user := &User{
		TenantID:     tenant.ID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}
	err = s.userRepo.Create(user)
	if err != nil {
		span.SetError(err)
		return nil, nil, err
	}

	// Track success event
	s.telemetry.TrackEvent(ctx, telemetry.Event{
		Name: "auth.register.success",
		Properties: map[string]interface{}{
			"user_id":   user.ID.String(),
			"tenant_id": tenant.ID.String(),
			"email":     req.Email,
		},
		Timestamp: time.Now(),
	})

	return user, tenant, nil
}

func (s *authService) LoginUser(req LoginRequest) (*User, error) {
	ctx := context.Background()
	span, ctx := s.telemetry.StartSpan(ctx, "auth.login_user")
	defer span.End()

	// Track metric
	s.telemetry.TrackMetric(ctx, telemetry.Metric{
		Name:  "auth.login.attempt",
		Value: 1,
		Tags:  map[string]string{"email": req.Email},
	})

	// Buscar user por email (cache já está no repository)
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		span.SetError(err)
		return nil, err
	}
	if user == nil {
		span.SetTag("user_not_found", "true")
		return nil, errors.New("invalid credentials")
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		span.SetTag("invalid_password", "true")
		return nil, errors.New("invalid credentials")
	}

	// Track success event
	s.telemetry.TrackEvent(ctx, telemetry.Event{
		Name: "auth.login.success",
		Properties: map[string]interface{}{
			"user_id": user.ID.String(),
			"email":   req.Email,
		},
		Timestamp: time.Now(),
	})

	return user, nil
} 