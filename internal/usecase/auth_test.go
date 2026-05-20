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

	t.Run("Success", func(t *testing.T) {
		user, err := au.Register(ctx, "testuser", "password123")
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", user.Username)
		}
		if user.PasswordHash != "" {
			t.Error("expected PasswordHash to be cleared")
		}
	})

	t.Run("DuplicateUser", func(t *testing.T) {
		_, err := au.Register(ctx, "testuser", "password123")
		if err != apperr.ErrUserAlreadyExists {
			t.Errorf("expected ErrUserAlreadyExists, got %v", err)
		}
	})

	t.Run("ShortUsername", func(t *testing.T) {
		_, err := au.Register(ctx, "ab", "password123")
		if err != apperr.ErrInvalidUsername {
			t.Errorf("expected ErrInvalidUsername, got %v", err)
		}
	})

	t.Run("ShortPassword", func(t *testing.T) {
		_, err := au.Register(ctx, "validuser", "short")
		if err != apperr.ErrInvalidPassword {
			t.Errorf("expected ErrInvalidPassword, got %v", err)
		}
	})

	t.Run("PasswordTooLong", func(t *testing.T) {
		longPass := string(make([]byte, 73))
		_, err := au.Register(ctx, "longpassuser", longPass)
		if err != apperr.ErrInvalidPassword {
			t.Errorf("expected ErrInvalidPassword for 73-byte password, got %v", err)
		}
	})

	t.Run("PasswordMaxLength", func(t *testing.T) {
		maxPass := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffff123456789012"
		if len(maxPass) != 72 {
			t.Skipf("test password is %d bytes, expected 72", len(maxPass))
		}
		_, err := au.Register(ctx, "maxpassuser", maxPass)
		if err != nil {
			t.Errorf("expected success for 72-byte password, got %v", err)
		}
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	repo := testutil.NewMockUserRepo()
	au := NewAuthUseCase(repo, "secret", time.Hour)

	ctx := context.Background()
	_, _ = au.Register(ctx, "testuser", "password123")

	t.Run("Success", func(t *testing.T) {
		user, token, err := au.Login(ctx, "testuser", "password123")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("expected username testuser, got %s", user.Username)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
		if user.PasswordHash != "" {
			t.Error("expected PasswordHash to be cleared")
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		_, _, err := au.Login(ctx, "testuser", "wrongpassword")
		if err != apperr.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		_, _, err := au.Login(ctx, "nonexistent", "password123")
		if err != apperr.ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}
