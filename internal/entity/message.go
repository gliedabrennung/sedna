package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	ChatID    string    `json:"chat_id"`
	MessageID string    `json:"message_id"`
	FromID    int64     `json:"from_id"`
	ToID      int64     `json:"to_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func MakeChatID(userA, userB int64) string {
	if userA > userB {
		userA, userB = userB, userA
	}
	return fmt.Sprintf("%d:%d", userA, userB)
}

func ValidateChatAccess(chatID string, userID int64) bool {
	parts := strings.SplitN(chatID, ":", 2)
	if len(parts) != 2 {
		return false
	}
	id1, err1 := strconv.ParseInt(parts[0], 10, 64)
	id2, err2 := strconv.ParseInt(parts[1], 10, 64)
	if err1 != nil || err2 != nil {
		return false
	}
	return userID == id1 || userID == id2
}
