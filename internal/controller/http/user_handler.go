package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gliedabrennung/messenger-core/internal/common/api"
	"github.com/gliedabrennung/messenger-core/internal/common/authctx"
	"github.com/gliedabrennung/messenger-core/internal/common/logger"
	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type UserService interface {
	SearchUsers(ctx context.Context, query string) ([]entity.User, error)
	GetUsersByIDs(ctx context.Context, ids []int64) ([]entity.User, error)
}

type UserHandler struct {
	users UserService
}

func NewUserHandler(users UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Search(ctx context.Context, c *app.RequestContext) {
	q := c.Query("q")
	if q == "" {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "query parameter 'q' is required", nil)
		return
	}

	users, err := h.users.SearchUsers(ctx, q)
	if err != nil {
		logger.CtxErrorf(ctx, "search failed: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to search users", nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetBulk(ctx context.Context, c *app.RequestContext) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		c.JSON(http.StatusOK, []entity.User{})
		return
	}

	parts := strings.Split(idsStr, ",")
	var ids []int64
	for _, p := range parts {
		if id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	users, err := h.users.GetUsersByIDs(ctx, ids)
	if err != nil {
		logger.CtxErrorf(ctx, "get bulk failed: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get users", nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Me(ctx context.Context, c *app.RequestContext) {
	userID, ok := authctx.UserID(c)
	if !ok {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing user context", nil)
		return
	}

	users, err := h.users.GetUsersByIDs(ctx, []int64{userID})
	if err != nil || len(users) == 0 {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "user not found", nil)
		return
	}

	c.JSON(http.StatusOK, users[0])
}
