package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gliedabrennung/messenger-core/internal/common/api"
	"github.com/gliedabrennung/messenger-core/internal/common/authctx"
	"github.com/gliedabrennung/messenger-core/internal/common/logger"
	"github.com/gliedabrennung/messenger-core/internal/domain"
	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type MessageHandler struct {
	repo domain.MessageRepository
}

func NewMessageHandler(repo domain.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

type historyResponse struct {
	Messages   []*entity.Message `json:"messages"`
	NextCursor string            `json:"next_cursor"`
}

func (h *MessageHandler) GetHistory(ctx context.Context, c *app.RequestContext) {
	if h.repo == nil {
		c.JSON(http.StatusOK, historyResponse{Messages: []*entity.Message{}})
		return
	}

	userID, ok := authctx.UserID(c)
	if !ok {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing user context", nil)
		return
	}

	partnerIDStr := c.Query("partner_id")
	if partnerIDStr == "" {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "partner_id is required", nil)
		return
	}

	partnerID, err := strconv.ParseInt(partnerIDStr, 10, 64)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid partner_id", nil)
		return
	}

	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	cursor := c.Query("cursor")

	chatID := entity.MakeChatID(userID, partnerID)
	messages, nextCursor, err := h.repo.GetChatHistory(ctx, chatID, limit, cursor)
	if err != nil {
		logger.CtxErrorf(ctx, "failed to get chat history: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to retrieve messages", nil)
		return
	}

	if messages == nil {
		messages = make([]*entity.Message, 0)
	}

	c.JSON(http.StatusOK, historyResponse{
		Messages:   messages,
		NextCursor: nextCursor,
	})
}
