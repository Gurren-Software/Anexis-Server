package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAPIKeyAuthSetsStandaloneUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(APIKeyAuth("secret"))
	router.GET("/files", func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			t.Fatalf("expected user id in context")
		}
		if userID != StandaloneUserID() {
			t.Fatalf("expected standalone user id, got %s", userID)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/files", nil)
	req.Header.Set("X-API-Key", "secret")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestOptionalAPIKeyAuthSetsStandaloneUserIDWhenKeyMatches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(OptionalAPIKeyAuth("secret"))
	router.GET("/links/:id", func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			t.Fatalf("expected user id in context")
		}
		if userID != StandaloneUserID() {
			t.Fatalf("expected standalone user id, got %s", userID)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/links/abc", nil)
	req.Header.Set("X-API-Key", "secret")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}
