package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/gliedabrennung/messenger-core/internal/apperr"
	"github.com/gliedabrennung/messenger-core/internal/common/api"
	"github.com/gliedabrennung/messenger-core/internal/common/logger"
	"github.com/gliedabrennung/messenger-core/internal/entity"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (*entity.User, error)
	Login(ctx context.Context, username, password string) (*entity.User, string, error)
}

type UserService interface {
	SearchUsers(ctx context.Context, query string) ([]entity.User, error)
	GetUsersByIDs(ctx context.Context, ids []int64) ([]entity.User, error)
}

type CookieConfig struct {
	Name     string
	MaxAge   int
	Secure   bool
	Domain   string
}

type AuthHandler struct {
	auth   AuthService
	cookie CookieConfig
}

func NewAuthHandler(auth AuthService, cookie CookieConfig) *AuthHandler {
	return &AuthHandler{auth: auth, cookie: cookie}
}

type UserHandler struct {
	users UserService
}

func NewUserHandler(users UserService) *UserHandler {
	return &UserHandler{users: users}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	User *entity.User `json:"user"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

func (h *AuthHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req authRequest
	if err := c.BindAndValidate(&req); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST", "invalid request body", err.Error())
		return
	}

	user, err := h.auth.Register(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidUsername) || errors.Is(err, apperr.ErrInvalidPassword) {
			api.ErrorResponse(c, http.StatusBadRequest,
				"INVALID_CREDENTIALS", err.Error(), nil)
			return
		}
		if errors.Is(err, apperr.ErrUserAlreadyExists) {
			api.ErrorResponse(c, http.StatusConflict,
				"USER_EXISTS", "username is already taken", nil)
			return
		}
		logger.CtxErrorf(ctx, "register failed: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR", "failed to register user", nil)
		return
	}

	c.JSON(http.StatusCreated, registerResponse{User: user})
}

func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req authRequest
	if err := c.BindAndValidate(&req); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST", "invalid request body", err.Error())
		return
	}

	user, token, err := h.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			api.ErrorResponse(c, http.StatusUnauthorized,
				"INVALID_CREDENTIALS", "invalid username or password", nil)
			return
		}
		logger.CtxErrorf(ctx, "login failed: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR", "failed to login", nil)
		return
	}

	h.setTokenCookie(c, token)
	c.JSON(http.StatusOK, loginResponse{Token: token, User: user})
}

func (h *AuthHandler) Logout(_ context.Context, c *app.RequestContext) {
	h.clearTokenCookie(c)
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthHandler) setTokenCookie(c *app.RequestContext, token string) {
	c.SetCookie(
		h.cookie.Name,
		token,
		h.cookie.MaxAge,
		"/",
		h.cookie.Domain,
		protocol.CookieSameSiteStrictMode,
		h.cookie.Secure,
		true,
	)
}

func (h *AuthHandler) clearTokenCookie(c *app.RequestContext) {
	c.SetCookie(
		h.cookie.Name,
		"",
		-1,
		"/",
		h.cookie.Domain,
		protocol.CookieSameSiteStrictMode,
		h.cookie.Secure,
		true,
	)
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
