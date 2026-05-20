package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/gliedabrennung/messenger-core/internal/entity"
	"github.com/gliedabrennung/messenger-core/internal/repository"
)

type mockUserRepo struct {
	users  map[string]*entity.User
	byID   map[int64]*entity.User
	nextID int64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:  make(map[string]*entity.User),
		byID:   make(map[int64]*entity.User),
		nextID: 1,
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.Username]; exists {
		return repository.ErrUserAlreadyExists
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.Username] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, ok := m.users[username]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, ok := m.byID[id]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepo) GetAll(ctx context.Context, limit, offset int) ([]entity.User, error) {
	var all []entity.User
	for _, u := range m.users {
		all = append(all, *u)
	}
	if offset >= len(all) {
		return []entity.User{}, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id int64) error {
	user, ok := m.byID[id]
	if !ok {
		return repository.ErrUserNotFound
	}
	delete(m.users, user.Username)
	delete(m.byID, id)
	return nil
}

func TestAuthUseCase_Register(t *testing.T) {
	repo := newMockUserRepo()
	au := NewAuthUseCase(repo, "secret", time.Hour)
	ctx := context.Background()

	user, err := au.Register(ctx, "testuser", "password123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", user.Username)
	}
	if user.ID == 0 {
		t.Error("expected non-zero user ID")
	}

	_, err = au.Register(ctx, "testuser", "password123")
	if err != ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	repo := newMockUserRepo()
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
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		_, _, err := au.Login(ctx, "testuser", "wrongpassword")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		_, _, err := au.Login(ctx, "nonexistent", "password123")
		if err != ErrInvalidCredentials {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestAuthUseCase_GetAllUsers(t *testing.T) {
	repo := newMockUserRepo()
	au := NewAuthUseCase(repo, "secret", time.Hour)
	ctx := context.Background()

	_, _ = au.Register(ctx, "alice", "password123")
	_, _ = au.Register(ctx, "bob", "password123")

	users, err := au.GetAllUsers(ctx, 100, 0)
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestAuthUseCase_DeleteUser(t *testing.T) {
	repo := newMockUserRepo()
	au := NewAuthUseCase(repo, "secret", time.Hour)
	ctx := context.Background()

	user, _ := au.Register(ctx, "alice", "password123")

	err := au.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	_, err = repo.GetByUsername(ctx, "alice")
	if err != repository.ErrUserNotFound {
		t.Error("expected user to be deleted")
	}
}
