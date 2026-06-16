package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// UserIDKey is the key used to store user ID in context
	UserIDKey = "user_id"
	// UserEmailKey is the key used to store user email in context
	UserEmailKey = "user_email"
	// APIKeyKey is the key used to store API key in context
	APIKeyKey = "api_key"
)

var standaloneUserID = uuid.MustParse("00000000-0000-4000-8000-000000000001")

// StandaloneUserID returns the stable single-user ID used in API key mode.
func StandaloneUserID() uuid.UUID {
	return standaloneUserID
}

// APIKeyAuth middleware validates API key for standalone mode
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CONFIGURATION_ERROR",
					"message": "API key not configured",
				},
			})
			return
		}

		apiKeyHeader := c.GetHeader("X-API-Key")
		if apiKeyHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "X-API-Key header required",
				},
			})
			return
		}

		if apiKeyHeader != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid API key",
				},
			})
			return
		}

		c.Set(APIKeyKey, apiKeyHeader)
		c.Set(UserIDKey, standaloneUserID)
		c.Next()
	}
}

// OptionalAPIKeyAuth middleware validates API key but doesn't require it
func OptionalAPIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" {
			c.Next()
			return
		}

		apiKeyHeader := c.GetHeader("X-API-Key")
		if apiKeyHeader == "" {
			c.Next()
			return
		}

		if apiKeyHeader == apiKey {
			c.Set(APIKeyKey, apiKeyHeader)
			c.Set(UserIDKey, standaloneUserID)
		}

		c.Next()
	}
}

// JWTAuth middleware validates JWT tokens
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header required",
				},
			})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid authorization header format",
				},
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid or expired token",
				},
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid token claims",
				},
			})
			return
		}

		// Extract user ID (stored as string UUID)
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid user ID in token",
				},
			})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid user ID format",
				},
			})
			return
		}

		// Store in context
		c.Set(UserIDKey, userID)
		if email, ok := claims["user_email"].(string); ok {
			c.Set(UserEmailKey, email)
		}

		c.Next()
	}
}

// OptionalJWTAuth middleware validates JWT tokens but doesn't require them
func OptionalJWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.Next()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserIDKey, userID)
		if email, ok := claims["user_email"].(string); ok {
			c.Set(UserEmailKey, email)
		}

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetAPIKey extracts API key from context
func GetAPIKey(c *gin.Context) (string, bool) {
	apiKey, exists := c.Get(APIKeyKey)
	if !exists {
		return "", false
	}
	key, ok := apiKey.(string)
	return key, ok
}

// AuthMiddleware returns the appropriate auth middleware based on server mode
func AuthMiddleware(secret, apiKey string, isStandalone bool) gin.HandlerFunc {
	if isStandalone && apiKey != "" {
		return APIKeyAuth(apiKey)
	}
	return JWTAuth(secret)
}

// OptionalAuthMiddleware returns an optional auth middleware for public endpoints
func OptionalAuthMiddleware(secret, apiKey string, isStandalone bool) gin.HandlerFunc {
	if isStandalone && apiKey != "" {
		return OptionalAPIKeyAuth(apiKey)
	}
	return OptionalJWTAuth(secret)
}
