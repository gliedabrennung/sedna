package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gliedabrennung/messenger-core/internal/apperr"
	"github.com/gliedabrennung/messenger-core/internal/entity"
	"github.com/gliedabrennung/messenger-core/internal/pkg/api"
)

// AuthService is the consumer-side interface for auth operations.
type AuthService interface {
	Register(ctx context.Context, username, password string) (*entity.User, error)
	Login(ctx context.Context, username, password string) (*entity.User, string, error)
}

type AuthHandler struct {
	auth AuthService
}

func NewAuthHandler(auth AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
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

	username, err := sanitizeAndValidate(req.Username, req.Password)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_CREDENTIALS", err.Error(), nil)
		return
	}

	user, err := h.auth.Register(ctx, username, req.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrUserAlreadyExists) {
			api.ErrorResponse(c, http.StatusConflict,
				"USER_EXISTS", "username is already taken", nil)
			return
		}
		hlog.CtxErrorf(ctx, "register failed: %v", err)
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
		hlog.CtxErrorf(ctx, "login failed: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR", "failed to login", nil)
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token, User: user})
}

// sanitizeAndValidate trims the username and validates both credentials.
// Returns the sanitized username or an error.
func sanitizeAndValidate(username, password string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 24 {
		return "", errors.New("username must be between 3 and 24 characters")
	}
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}
	return username, nil
}
