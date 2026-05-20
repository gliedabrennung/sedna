package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "DSN=postgres://localhost/test\nJWT_SECRET=testsecret\nJWT_TTL=1h\nADDR=:9090\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(envPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.DSN != "postgres://localhost/test" {
		t.Errorf("expected DSN postgres://localhost/test, got %s", cfg.DSN)
	}
	if cfg.JWTSecret != "testsecret" {
		t.Errorf("expected JWTSecret testsecret, got %s", cfg.JWTSecret)
	}
	if cfg.Addr != ":9090" {
		t.Errorf("expected Addr :9090, got %s", cfg.Addr)
	}
}

func TestLoadConfig_MissingFile_FallsBackToEnv(t *testing.T) {
	t.Setenv("DSN", "postgres://envhost/db")
	t.Setenv("JWT_SECRET", "envsecret")

	cfg, err := LoadConfig("nonexistent.env")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.DSN != "postgres://envhost/db" {
		t.Errorf("expected DSN from env, got %s", cfg.DSN)
	}
}

func TestLoadConfig_MissingRequired(t *testing.T) {
	_, err := LoadConfig("nonexistent.env")
	if err == nil {
		t.Error("expected error for missing required fields")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	t.Setenv("DSN", "postgres://localhost/db")
	t.Setenv("JWT_SECRET", "secret")

	cfg, err := LoadConfig("nonexistent.env")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Addr != ":8080" {
		t.Errorf("expected default Addr :8080, got %s", cfg.Addr)
	}
	if cfg.JWTTTL.Hours() != 24 {
		t.Errorf("expected default JWT_TTL 24h, got %v", cfg.JWTTTL)
	}
}
