package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claudineijrdev/sib-crm-backend/internal/auth"
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockAuthService := &auth.MockAuthService{}
	handler := auth.NewAuthHandler(mockAuthService)

	// Test data
	userID := uuid.New()
	tenantID := uuid.New()
	
	reqBody := auth.RegisterRequest{
		Name:     "Test Company",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &auth.User{
		ID:       userID,
		TenantID: tenantID,
		Email:    "test@example.com",
	}

	expectedTenant := &tenants.Tenant{
		ID:   tenantID,
		Name: "Test Company",
	}

	// Mock behavior
	mockAuthService.RegisterUserFunc = func(req auth.RegisterRequest) (*auth.User, *tenants.Tenant, error) {
		return expectedUser, expectedTenant, nil
	}

	// Create request
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Register(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response auth.RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, tenantID.String(), response.TenantID)
	assert.Equal(t, userID.String(), response.UserID)
}

func TestAuthHandler_Login(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockAuthService := &auth.MockAuthService{}
	handler := auth.NewAuthHandler(mockAuthService)

	// Test data
	userID := uuid.New()
	
	reqBody := auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &auth.User{
		ID:    userID,
		Email: "test@example.com",
	}

	// Mock behavior
	mockAuthService.LoginUserFunc = func(req auth.LoginRequest) (*auth.User, error) {
		return expectedUser, nil
	}

	// Create request
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Login(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response auth.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), response.UserID)
	assert.Equal(t, "jwt-token-here", response.Token)
} 