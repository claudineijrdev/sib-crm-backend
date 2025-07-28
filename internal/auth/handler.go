package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tenant, err := h.authService.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := RegisterResponse{
		TenantID: tenant.ID.String(),
		UserID:   user.ID.String(),
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implementar geração de JWT token
	response := LoginResponse{
		UserID: user.ID.String(),
		Token:  "jwt-token-here", // Placeholder
	}

	c.JSON(http.StatusOK, response)
}
