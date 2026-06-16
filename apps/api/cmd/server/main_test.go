package main

import (
	"strings"
	"testing"

	swaggerdocs "github.com/Gurren-Software/Anexis-Server/apps/api/docs"
)

func TestConfigureSwaggerAuthUsesAPIKeyInStandaloneMode(t *testing.T) {
	originalTemplate := swaggerdocs.SwaggerInfo.SwaggerTemplate
	defer func() {
		swaggerdocs.SwaggerInfo.SwaggerTemplate = originalTemplate
	}()

	configureSwaggerAuth(true)

	template := swaggerdocs.SwaggerInfo.SwaggerTemplate
	if !strings.Contains(template, `"name": "X-API-Key"`) {
		t.Fatalf("expected standalone Swagger auth header to be X-API-Key")
	}
	if strings.Contains(template, `"name": "Authorization"`) {
		t.Fatalf("expected standalone Swagger auth header not to use Authorization")
	}
	if !strings.Contains(template, `"APIKeyAuth": []`) {
		t.Fatalf("expected protected endpoints to use APIKeyAuth")
	}
}

func TestConfigureSwaggerAuthKeepsBearerAuthOutsideStandaloneMode(t *testing.T) {
	originalTemplate := swaggerdocs.SwaggerInfo.SwaggerTemplate
	defer func() {
		swaggerdocs.SwaggerInfo.SwaggerTemplate = originalTemplate
	}()

	configureSwaggerAuth(false)

	if swaggerdocs.SwaggerInfo.SwaggerTemplate != originalTemplate {
		t.Fatalf("expected non-standalone Swagger template to stay unchanged")
	}
}
