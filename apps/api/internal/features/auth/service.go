package auth

import (
	"errors"
	"time"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)

// Service handles authentication business logic
type Service struct {
	repo          *Repository
	jwtSecret     string
	jwtExpiration time.Duration
}

// NewService creates a new auth service
func NewService(repo *Repository, jwtSecret string, expirationHours int) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     jwtSecret,
		jwtExpiration: time.Duration(expirationHours) * time.Hour,
	}
}

// Register creates a new user account
func (s *Service) Register(req *RegisterRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.repo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(req *LoginRequest) (*TokenResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	return s.generateTokenResponse(user)
}

// RefreshToken generates new tokens from a refresh token
func (s *Service) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// Parse refresh token
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidCredentials
	}

	// Get user
	userIDStr := claims.Subject
	user, err := s.repo.FindByEmail(userIDStr)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokenResponse(user)
}

// GetUser returns user by ID
func (s *Service) GetUser(userID uint) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// ChangePassword changes user password
func (s *Service) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(userID, string(hashedPassword))
}

func (s *Service) generateTokenResponse(user *models.User) (*TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiration)

	// Access token claims
	email := user.Email

	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   email,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token (longer expiration)
	refreshExpiresAt := now.Add(s.jwtExpiration * 7) // 7x access token expiration
	refreshClaims := jwt.MapClaims{
		"sub": email,
		"exp": refreshExpiresAt.Unix(),
		"iat": now.Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtExpiration.Seconds()),
		RefreshToken: refreshTokenString,
	}, nil
}

// ToUserResponse converts user model to response DTO
func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		StorageQuota: user.StorageQuota,
		StorageUsed:  user.StorageUsed,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
	}
}
