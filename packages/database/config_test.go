package database

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultConfigReadsEnvironment(t *testing.T) {
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "15432")
	t.Setenv("DB_USER", "anexis")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "anexis_test")
	t.Setenv("DB_SSLMODE", "require")

	cfg := DefaultConfig()

	if cfg.Host != "db" || cfg.Port != "15432" || cfg.User != "anexis" {
		t.Fatalf("unexpected db config: %#v", cfg)
	}
	if cfg.Password != "secret" || cfg.DBName != "anexis_test" || cfg.SSLMode != "require" {
		t.Fatalf("unexpected db credentials: %#v", cfg)
	}
	if cfg.MaxIdleConns != 10 || cfg.MaxOpenConns != 100 || cfg.ConnMaxLifetime != time.Hour {
		t.Fatalf("unexpected pool config: %#v", cfg)
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "secret",
		DBName:   "anexis",
		SSLMode:  "disable",
	}

	dsn := cfg.DSN()
	for _, part := range []string{
		"host=localhost",
		"user=postgres",
		"password=secret",
		"dbname=anexis",
		"port=5432",
		"sslmode=disable",
	} {
		if !strings.Contains(dsn, part) {
			t.Fatalf("expected DSN %q to contain %q", dsn, part)
		}
	}
}
