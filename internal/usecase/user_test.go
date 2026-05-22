package usecase

import (
	"context"
	"testing"

	"github.com/gliedabrennung/messenger-core/internal/testutil"
)

func TestUserUseCase_SearchUsers(t *testing.T) {
	repo := testutil.NewMockUserRepo()
	au := NewAuthUseCase(repo, "secret", 0)
	uu := NewUserUseCase(repo)

	ctx := context.Background()
	_, _ = au.Register(ctx, "alice", "password123")
	_, _ = au.Register(ctx, "bob", "password123")

	t.Run("FindsBySubstring", func(t *testing.T) {
		users, err := uu.SearchUsers(ctx, "ali")
		if err != nil {
			t.Fatalf("SearchUsers failed: %v", err)
		}
		if len(users) != 1 {
			t.Fatalf("expected 1 user, got %d", len(users))
		}
		if users[0].Username != "alice" {
			t.Errorf("expected alice, got %s", users[0].Username)
		}
		if users[0].PasswordHash != "" {
			t.Error("expected PasswordHash to be cleared")
		}
	})

	t.Run("NoResults", func(t *testing.T) {
		users, err := uu.SearchUsers(ctx, "zzz")
		if err != nil {
			t.Fatalf("SearchUsers failed: %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})
}

func TestUserUseCase_GetUsersByIDs(t *testing.T) {
	repo := testutil.NewMockUserRepo()
	au := NewAuthUseCase(repo, "secret", 0)
	uu := NewUserUseCase(repo)

	ctx := context.Background()
	u1, _ := au.Register(ctx, "user1", "password123")
	_, _ = au.Register(ctx, "user2", "password123")

	t.Run("FindsExisting", func(t *testing.T) {
		users, err := uu.GetUsersByIDs(ctx, []int64{u1.ID})
		if err != nil {
			t.Fatalf("GetUsersByIDs failed: %v", err)
		}
		if len(users) != 1 {
			t.Fatalf("expected 1 user, got %d", len(users))
		}
		if users[0].PasswordHash != "" {
			t.Error("expected PasswordHash to be cleared")
		}
	})

	t.Run("EmptyIDs", func(t *testing.T) {
		users, err := uu.GetUsersByIDs(ctx, nil)
		if err != nil {
			t.Fatalf("GetUsersByIDs failed: %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})
}
