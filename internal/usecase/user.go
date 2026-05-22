package usecase

import (
	"context"
	"fmt"

	"github.com/gliedabrennung/messenger-core/internal/domain"
	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	users, err := u.repo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}

func (u *UserUseCase) GetUsersByIDs(ctx context.Context, ids []int64) ([]entity.User, error) {
	users, err := u.repo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get users by ids: %w", err)
	}
	for i := range users {
		users[i].PasswordHash = ""
	}
	return users, nil
}
