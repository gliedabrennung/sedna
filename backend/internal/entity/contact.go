package entity

import "time"

type Contact struct {
	UserID    int64     `json:"user_id"`
	ContactID int64     `json:"contact_id"`
	Status    string    `json:"status"` // "pending", "accepted"
	CreatedAt time.Time `json:"created_at"`
}
