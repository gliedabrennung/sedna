package usecase

import (
	"context"
	"testing"

	"github.com/gliedabrennung/messenger-core/internal/entity"
)

// mockMessageRepo is an in-memory implementation of MessageRepository.
type mockMessageRepo struct {
	messages []entity.Message
	nextID   int64
}

func (m *mockMessageRepo) SaveMessage(_ context.Context, msg *entity.Message) error {
	m.nextID++
	msg.ID = m.nextID
	msg.Status = "sent"
	m.messages = append(m.messages, *msg)
	return nil
}

func (m *mockMessageRepo) GetMessagesBetween(_ context.Context, user1, user2 int64, limit, offset int) ([]entity.Message, error) {
	var result []entity.Message
	for _, msg := range m.messages {
		if (msg.FromID == user1 && msg.ToID == user2) || (msg.FromID == user2 && msg.ToID == user1) {
			result = append(result, msg)
		}
	}
	end := offset + limit
	if offset >= len(result) {
		return nil, nil
	}
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *mockMessageRepo) UpdateMessageStatuses(_ context.Context, fromID, toID int64, status string) error {
	for i, msg := range m.messages {
		if msg.FromID == fromID && msg.ToID == toID {
			m.messages[i].Status = status
		}
	}
	return nil
}

func (m *mockMessageRepo) GetUnreadCounts(_ context.Context, userID int64) (map[int64]int, error) {
	counts := make(map[int64]int)
	for _, msg := range m.messages {
		if msg.ToID == userID && msg.Status != "read" {
			counts[msg.FromID]++
		}
	}
	return counts, nil
}

// ---- Tests ----

func TestMessageUseCase_SendMessage(t *testing.T) {
	repo := &mockMessageRepo{}
	uc := NewMessageUseCase(repo)
	ctx := context.Background()

	msg, err := uc.SendMessage(ctx, 1, 2, "text", "Hello!")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}
	if msg.ID == 0 {
		t.Error("expected non-zero message ID")
	}
	if msg.Status != "sent" {
		t.Errorf("expected status 'sent', got '%s'", msg.Status)
	}
	if msg.Content != "Hello!" {
		t.Errorf("expected content 'Hello!', got '%s'", msg.Content)
	}
}

func TestMessageUseCase_GetChatHistory(t *testing.T) {
	repo := &mockMessageRepo{}
	uc := NewMessageUseCase(repo)
	ctx := context.Background()

	_, _ = uc.SendMessage(ctx, 1, 2, "text", "First message")
	_, _ = uc.SendMessage(ctx, 2, 1, "text", "Reply")
	_, _ = uc.SendMessage(ctx, 1, 2, "text", "Third message")

	msgs, err := uc.GetChatHistory(ctx, 1, 2, 10, 0)
	if err != nil {
		t.Fatalf("GetChatHistory failed: %v", err)
	}
	if len(msgs) != 3 {
		t.Errorf("expected 3 messages, got %d", len(msgs))
	}
}

func TestMessageUseCase_GetChatHistory_Pagination(t *testing.T) {
	repo := &mockMessageRepo{}
	uc := NewMessageUseCase(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, _ = uc.SendMessage(ctx, 1, 2, "text", "msg")
	}

	msgs, err := uc.GetChatHistory(ctx, 1, 2, 2, 1) // skip first, get 2
	if err != nil {
		t.Fatalf("GetChatHistory pagination failed: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages with limit=2 offset=1, got %d", len(msgs))
	}
}

func TestMessageUseCase_UpdateStatus(t *testing.T) {
	repo := &mockMessageRepo{}
	uc := NewMessageUseCase(repo)
	ctx := context.Background()

	_, _ = uc.SendMessage(ctx, 1, 2, "text", "Hello!")

	err := uc.UpdateStatus(ctx, 1, 2, "read")
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}
	if repo.messages[0].Status != "read" {
		t.Errorf("expected status 'read', got '%s'", repo.messages[0].Status)
	}
}

func TestMessageUseCase_GetUnreadCounts(t *testing.T) {
	repo := &mockMessageRepo{}
	uc := NewMessageUseCase(repo)
	ctx := context.Background()

	_, _ = uc.SendMessage(ctx, 1, 2, "text", "msg1")
	_, _ = uc.SendMessage(ctx, 1, 2, "text", "msg2")
	_, _ = uc.SendMessage(ctx, 3, 2, "text", "msg3") // from user 3 to user 2

	counts, err := uc.GetUnreadCounts(ctx, 2)
	if err != nil {
		t.Fatalf("GetUnreadCounts failed: %v", err)
	}
	if counts[1] != 2 {
		t.Errorf("expected 2 unread from user 1, got %d", counts[1])
	}
	if counts[3] != 1 {
		t.Errorf("expected 1 unread from user 3, got %d", counts[3])
	}
}
