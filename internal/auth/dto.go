package auth

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	TenantID string `json:"tenant_id"`
	UserID   string `json:"user_id"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
} 