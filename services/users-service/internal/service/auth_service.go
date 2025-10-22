package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/domain"
	"github.com/youngermaster/my-distributed-bookstore/services/users-service/internal/repository"
	customJWT "github.com/youngermaster/my-distributed-bookstore/services/users-service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(ctx context.Context, email, password, fullName string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
	ValidateToken(ctx context.Context, token string) (*customJWT.Claims, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

type authService struct {
	userRepo repository.UserRepository
	jwtManager *customJWT.JWTManager
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.UserRepository, jwtManager *customJWT.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, email, password, fullName string) (*domain.User, error) {
	// Validate input
	if email == "" || password == "" || fullName == "" {
		return nil, ErrInvalidInput
	}

	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default "customer" role (assuming it exists)
	// This is a simplified approach; in production, ensure the role exists first
	// For now, we'll skip this and handle it in migrations or seed data

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if !verifyPassword(user.PasswordHash, password) {
		return "", nil, ErrInvalidCredentials
	}

	// Extract role names
	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*customJWT.Claims, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (string, error) {
	newToken, err := s.jwtManager.RefreshToken(token)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	return newToken, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	// In a real implementation, you might want to blacklist the token
	// or delete sessions from a session store
	// For now, this is a placeholder
	return nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// verifyPassword checks if the provided password matches the hashed password
func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// hashToken creates a SHA-256 hash of a token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
