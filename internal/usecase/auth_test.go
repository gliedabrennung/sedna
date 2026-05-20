package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/gliedabrennung/messenger-core/internal/apperr"
	"github.com/gliedabrennung/messenger-core/internal/testutil"
)

func TestAuthUseCase_Register(t *testing.T) {
	repo := testutil.NewMockUserRepo()
	au := NewAuthUseCase(repo, "secret", time.Hour)

	ctx := context.Background()
	user, err := au.Register(ctx, "testuser", "password")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", user.Username)
	}

	_, err = au.Register(ctx, "testuser", "password")
	if err != apperr.ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	repo := testutil.NewMockUserRepo()
	au := NewAuthUseCase(repo, "secret", time.Hour)

	ctx := context.Background()
	_, _ = au.Register(ctx, "testuser", "password")

	t.Run("Success", func(t *testing.T) {
		user, token, err := au.Login(ctx, "testuser", "password")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", user.Username)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		_, _, err := au.Login(ctx, "testuser", "wrongpassword")
		if err != apperr.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		_, _, err := au.Login(ctx, "nonexistent", "password")
		if err != apperr.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}
