package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gliedabrennung/messenger-core/internal/entity"
	"github.com/gliedabrennung/messenger-core/internal/pkg/api"
	"github.com/gliedabrennung/messenger-core/internal/usecase"
)

type AuthHandler struct {
	useCase *usecase.AuthUseCase
}

func NewAuthHandler(useCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
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

// Register godoc
// @Summary Register a new user
// @Description Creates a new user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authRequest true "Register details"
// @Success 201 {object} registerResponse
// @Failure 400 {object} api.Error
// @Failure 409 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req authRequest
	if err := c.BindAndValidate(&req); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST", "invalid request body", err.Error())
		return
	}
	if err := validateCredentials(req.Username, req.Password); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_CREDENTIALS", err.Error(), nil)
		return
	}

	user, err := h.useCase.Register(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
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

// Login godoc
// @Summary Login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authRequest true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /auth/login [post]
func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req authRequest
	if err := c.BindAndValidate(&req); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST", "invalid request body", err.Error())
		return
	}

	user, token, err := h.useCase.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			api.ErrorResponse(c, http.StatusUnauthorized,
				"INVALID_CREDENTIALS", "invalid username or password", nil)
			return
		}
		api.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR", "failed to login", nil)
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token, User: user})
}

func validateCredentials(username, password string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 24 {
		return errors.New("username must be between 3 and 24 characters")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

type usersResponse struct {
	Users []entity.User `json:"users"`
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves a list of all registered users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} usersResponse
// @Failure 500 {object} api.Error
// @Router /users [get]
func (h *AuthHandler) GetAllUsers(ctx context.Context, c *app.RequestContext) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 100 // default limit
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	users, err := h.useCase.GetAllUsers(ctx, limit, offset)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get users: %v", err)
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get users", nil)
		return
	}
	
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, usersResponse{Users: users})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /admin/users/{id} [delete]
func (h *AuthHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_PARAM", "invalid user id", nil)
		return
	}

	if err := h.useCase.DeleteUser(ctx, id); err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete user", err.Error())
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true})
}
