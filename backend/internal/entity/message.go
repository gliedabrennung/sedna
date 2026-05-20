package entity

import "time"

type Message struct {
	ID        int64     `json:"id"`
	FromID    int64     `json:"from_id"`
	ToID      int64     `json:"to_id"`
	Type      string    `json:"type"` // "text" or "audio"
	Content   string    `json:"content"`
	Status    string    `json:"status"` // "sent", "delivered", "read"
	CreatedAt time.Time `json:"created_at"`
}
