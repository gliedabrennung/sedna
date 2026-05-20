package domain

import (
	"context"

	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type MessageRepository interface {
	Save(ctx context.Context, msg *entity.Message) error
	GetChatHistory(ctx context.Context, chatID string, limit int, cursor string) ([]*entity.Message, string, error)
	Subscribe(ctx context.Context, chatID string) (<-chan *entity.Message, func() error, error)
}
