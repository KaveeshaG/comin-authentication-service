package service

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/repository"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
)

type AuthService struct {
	repo       repository.AuthRepository
	jwtService *jwt.JWTService
}

func NewAuthService(repo repository.AuthRepository, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Create organization
	org := &domain.Organization{
		Name:   req.OrganizationName,
		Slug:   generateSlug(req.OrganizationName),
		Status: "active",
	}
	if err := s.repo.CreateOrganization(org); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		OrganizationID: org.ID,
		Email:          req.Email,
		Password:       string(hashedPassword),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Role:           "admin",
		Status:         "active",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(user)
}

func (s *AuthService) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokens(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*domain.AuthResponse, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateTokens(user)
}

func (s *AuthService) ValidateToken(token string) (*jwt.Claims, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *AuthService) generateTokens(user *domain.User) (*domain.AuthResponse, error) {
	// Generate access token
	accessToken, err := s.jwtService.GenerateToken(
		user.ID,
		user.OrganizationID,
		user.Email,
		user.Role,
		time.Hour*24, // 24 hour expiry
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateToken(
		user.ID,
		user.OrganizationID,
		user.Email,
		user.Role,
		time.Hour*24*30, // 30 days expiry
	)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		TokenType:      "Bearer",
		ExpiresIn:      int64(time.Hour * 24 / time.Second),
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
	}, nil
}

func generateSlug(name string) string {
	// TODO: Implement proper slug generation with validation
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
