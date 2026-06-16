package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccessResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		write      func(*gin.Context)
		wantStatus int
	}{
		{
			name: "ok",
			write: func(c *gin.Context) {
				OK(c, gin.H{"id": "123"})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "created",
			write: func(c *gin.Context) {
				Created(c, gin.H{"id": "123"})
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "ok with meta",
			write: func(c *gin.Context) {
				OKWithMeta(c, []string{"a"}, &Meta{Page: 1, PerPage: 10, Total: 20, TotalPages: 2})
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)

			tt.write(c)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}

			var got Response
			if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if !got.Success {
				t.Fatalf("expected success response")
			}
		})
	}
}

func TestErrorResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		write      func(*gin.Context)
		wantStatus int
		wantCode   string
	}{
		{
			name: "bad request",
			write: func(c *gin.Context) {
				BadRequest(c, "BAD_INPUT", "bad input")
			},
			wantStatus: http.StatusBadRequest,
			wantCode:   "BAD_INPUT",
		},
		{
			name: "unauthorized",
			write: func(c *gin.Context) {
				Unauthorized(c, "missing token")
			},
			wantStatus: http.StatusUnauthorized,
			wantCode:   "UNAUTHORIZED",
		},
		{
			name: "validation",
			write: func(c *gin.Context) {
				ValidationError(c, "invalid", "name is required")
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)

			tt.write(c)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}

			var got Response
			if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if got.Success {
				t.Fatalf("expected error response")
			}
			if got.Error == nil || got.Error.Code != tt.wantCode {
				t.Fatalf("expected error code %q, got %#v", tt.wantCode, got.Error)
			}
		})
	}
}

func TestNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	NoContent(c)

	if c.Writer.Status() != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", c.Writer.Status())
	}
}
