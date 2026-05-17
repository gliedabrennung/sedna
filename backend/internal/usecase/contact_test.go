package usecase

import (
	"context"
	"testing"

	"github.com/gliedabrennung/messenger-core/internal/entity"
)

// mockContactRepo is an in-memory implementation of ContactRepository.
type mockContactRepo struct {
	contacts []entity.Contact
	users    map[int64]entity.User
}

func newMockContactRepo() *mockContactRepo {
	return &mockContactRepo{
		users: map[int64]entity.User{
			1: {ID: 1, Username: "alice", Role: "user"},
			2: {ID: 2, Username: "bob", Role: "user"},
			3: {ID: 3, Username: "charlie", Role: "user"},
		},
	}
}

func (m *mockContactRepo) AddContact(_ context.Context, userID, contactID int64) error {
	// avoid duplicates
	for _, c := range m.contacts {
		if c.UserID == userID && c.ContactID == contactID {
			return nil
		}
	}
	m.contacts = append(m.contacts, entity.Contact{
		UserID:    userID,
		ContactID: contactID,
		Status:    "pending",
	})
	return nil
}

func (m *mockContactRepo) UpdateStatus(_ context.Context, userID, contactID int64, status string) error {
	for i, c := range m.contacts {
		if c.UserID == userID && c.ContactID == contactID {
			m.contacts[i].Status = status
			return nil
		}
	}
	return nil
}

func (m *mockContactRepo) GetContacts(_ context.Context, userID int64) ([]entity.User, error) {
	var result []entity.User
	for _, c := range m.contacts {
		if c.Status != "accepted" {
			continue
		}
		if c.UserID == userID {
			if u, ok := m.users[c.ContactID]; ok {
				result = append(result, u)
			}
		} else if c.ContactID == userID {
			if u, ok := m.users[c.UserID]; ok {
				result = append(result, u)
			}
		}
	}
	return result, nil
}

func (m *mockContactRepo) DeleteContact(_ context.Context, userID, contactID int64) error {
	var remaining []entity.Contact
	for _, c := range m.contacts {
		if !((c.UserID == userID && c.ContactID == contactID) ||
			(c.UserID == contactID && c.ContactID == userID)) {
			remaining = append(remaining, c)
		}
	}
	m.contacts = remaining
	return nil
}

// ---- Tests ----

func TestContactUseCase_SendAndAcceptRequest(t *testing.T) {
	repo := newMockContactRepo()
	uc := NewContactUseCase(repo)
	ctx := context.Background()

	// alice sends bob a request
	err := uc.SendRequest(ctx, 1, 2)
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}
	if len(repo.contacts) != 1 {
		t.Fatalf("expected 1 contact entry, got %d", len(repo.contacts))
	}
	if repo.contacts[0].Status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", repo.contacts[0].Status)
	}

	// bob accepts alice's request
	err = uc.AcceptRequest(ctx, 2, 1) // bob (userID=2) accepts alice (contactID=1)
	if err != nil {
		t.Fatalf("AcceptRequest failed: %v", err)
	}
	if repo.contacts[0].Status != "accepted" {
		t.Errorf("expected status 'accepted', got '%s'", repo.contacts[0].Status)
	}
}

func TestContactUseCase_GetUserContacts(t *testing.T) {
	repo := newMockContactRepo()
	uc := NewContactUseCase(repo)
	ctx := context.Background()

	_ = uc.SendRequest(ctx, 1, 2)
	_ = uc.AcceptRequest(ctx, 2, 1)

	_ = uc.SendRequest(ctx, 1, 3)
	_ = uc.AcceptRequest(ctx, 3, 1)

	contacts, err := uc.GetUserContacts(ctx, 1)
	if err != nil {
		t.Fatalf("GetUserContacts failed: %v", err)
	}
	if len(contacts) != 2 {
		t.Errorf("expected 2 contacts for alice, got %d", len(contacts))
	}
}

func TestContactUseCase_RemoveContact(t *testing.T) {
	repo := newMockContactRepo()
	uc := NewContactUseCase(repo)
	ctx := context.Background()

	_ = uc.SendRequest(ctx, 1, 2)
	_ = uc.AcceptRequest(ctx, 2, 1)

	err := uc.RemoveContact(ctx, 1, 2)
	if err != nil {
		t.Fatalf("RemoveContact failed: %v", err)
	}

	contacts, _ := uc.GetUserContacts(ctx, 1)
	if len(contacts) != 0 {
		t.Errorf("expected 0 contacts after removal, got %d", len(contacts))
	}
}

func TestContactUseCase_DuplicateRequest(t *testing.T) {
	repo := newMockContactRepo()
	uc := NewContactUseCase(repo)
	ctx := context.Background()

	_ = uc.SendRequest(ctx, 1, 2)
	_ = uc.SendRequest(ctx, 1, 2) // duplicate — should not create another entry

	if len(repo.contacts) != 1 {
		t.Errorf("expected 1 contact entry (no duplicates), got %d", len(repo.contacts))
	}
}
