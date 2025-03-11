// internal/auth/domain/auth.go
package domain

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	OrganizationName string `json:"organization_name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8"`
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
}

type AuthResponse struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenType      string    `json:"token_type"`
	ExpiresIn      int64     `json:"expires_in"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Role           string    `json:"role"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
