package usecase

import (
	"context"

	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, msg *entity.Message) error
	GetMessagesBetween(ctx context.Context, user1, user2 int64, limit, offset int) ([]entity.Message, error)
	UpdateMessageStatuses(ctx context.Context, fromID, toID int64, status string) error
	GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error)
}

type MessageUseCase struct {
	repo MessageRepository
}

func NewMessageUseCase(repo MessageRepository) *MessageUseCase {
	return &MessageUseCase{
		repo: repo,
	}
}

func (m *MessageUseCase) SendMessage(ctx context.Context, from, to int64, msgType, content string) (*entity.Message, error) {
	msg := &entity.Message{
		FromID:  from,
		ToID:    to,
		Type:    msgType,
		Content: content,
	}
	if err := m.repo.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (m *MessageUseCase) GetChatHistory(ctx context.Context, user1, user2 int64, limit, offset int) ([]entity.Message, error) {
	return m.repo.GetMessagesBetween(ctx, user1, user2, limit, offset)
}

func (m *MessageUseCase) UpdateStatus(ctx context.Context, fromID, toID int64, status string) error {
	return m.repo.UpdateMessageStatuses(ctx, fromID, toID, status)
}

func (m *MessageUseCase) GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error) {
	return m.repo.GetUnreadCounts(ctx, userID)
}
