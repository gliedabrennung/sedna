package domain

import (
	"context"

	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Search(ctx context.Context, query string) ([]entity.User, error)
	GetByIDs(ctx context.Context, ids []int64) ([]entity.User, error)
}

type MessageRepository interface {
	Save(ctx context.Context, msg *entity.Message) error
	GetChatHistory(ctx context.Context, chatID string, limit int, cursor string) ([]*entity.Message, string, error)
	Subscribe(ctx context.Context, chatID string) (<-chan *entity.Message, func() error, error)
}
