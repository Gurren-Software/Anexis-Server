package auth

import (
	"errors"
	"time"

	"github.com/Treefle-labs/anexis-server/packages/database/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

// Service handles auth business logic
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

// Register registers a new user
func (s *Service) Register(req *RegisterRequest) (*models.User, error) {
	// Check if user exists
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
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

	return s.generateTokens(user)
}

// RefreshToken refreshes access token using refresh token
func (s *Service) RefreshToken(refreshToken string) (*TokenResponse, error) {
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.FindByID(userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokens(user)
}

// GetUserByID returns user by ID
func (s *Service) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

// ChangePassword changes user password
func (s *Service) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(userID, string(hashedPassword))
}

func (s *Service) generateTokens(user *models.User) (*TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiration)

	// Access token claims
	accessClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
		"type":    "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token (longer expiration)
	refreshExpiresAt := now.Add(s.jwtExpiration * 7)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     refreshExpiresAt.Unix(),
		"iat":     now.Unix(),
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtExpiration.Seconds()),
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *Service) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
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
