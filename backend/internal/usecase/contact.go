package usecase

import (
	"context"

	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type ContactRepository interface {
	AddContact(ctx context.Context, userID, contactID int64) error
	UpdateStatus(ctx context.Context, userID, contactID int64, status string) error
	GetContacts(ctx context.Context, userID int64) ([]entity.User, error)
	DeleteContact(ctx context.Context, userID, contactID int64) error
}

type ContactUseCase struct {
	repo ContactRepository
}

func NewContactUseCase(repo ContactRepository) *ContactUseCase {
	return &ContactUseCase{repo: repo}
}

func (c *ContactUseCase) SendRequest(ctx context.Context, userID, contactID int64) error {
	return c.repo.AddContact(ctx, userID, contactID)
}

func (c *ContactUseCase) AcceptRequest(ctx context.Context, userID, contactID int64) error {
	// The person accepting the request is the contactID of the original row.
	return c.repo.UpdateStatus(ctx, contactID, userID, "accepted")
}

func (c *ContactUseCase) GetUserContacts(ctx context.Context, userID int64) ([]entity.User, error) {
	return c.repo.GetContacts(ctx, userID)
}

func (c *ContactUseCase) RemoveContact(ctx context.Context, userID, contactID int64) error {
	return c.repo.DeleteContact(ctx, userID, contactID)
}
