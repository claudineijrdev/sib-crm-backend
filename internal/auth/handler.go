package auth

import (
	"net/http"

	"github.com/claudineijrdev/sib-crm-backend/internal/platform/database"
	"github.com/claudineijrdev/sib-crm-backend/internal/tenants"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var tenant tenants.Tenant
	var user User

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		tenant = tenants.Tenant{Name: req.Name}
		if err := tx.Create(&tenant).Error; err != nil {
			return err
		}

		user = User{
			TenantID:     tenant.ID,
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tenant_id": tenant.ID, "user_id": user.ID})
}
