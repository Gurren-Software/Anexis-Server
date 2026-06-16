package http

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewServerConfiguresRouterAndAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := NewServer("127.0.0.1", "9090", true)

	if server.Router() == nil {
		t.Fatalf("expected router to be configured")
	}
	if server.server.Addr != "127.0.0.1:9090" {
		t.Fatalf("unexpected server address %q", server.server.Addr)
	}
}
