package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gliedabrennung/messenger-core/internal/pkg/api"
	"github.com/gliedabrennung/messenger-core/internal/pkg/authctx"
	"github.com/gliedabrennung/messenger-core/internal/usecase"
)

type ContactHandler struct {
	useCase *usecase.ContactUseCase
}

func NewContactHandler(uc *usecase.ContactUseCase) *ContactHandler {
	return &ContactHandler{useCase: uc}
}

type contactRequest struct {
	ContactID int64 `json:"contact_id"`
}

// AddContact godoc
// @Summary Add a new contact
// @Description Sends a contact request to another user
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body contactRequest true "Contact details"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /contacts [post]
func (h *ContactHandler) AddContact(ctx context.Context, c *app.RequestContext) {
	userID, ok := authctx.UserID(c)
	if !ok {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing user context", nil)
		return
	}

	var req contactRequest
	if err := c.BindAndValidate(&req); err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid body", err.Error())
		return
	}

	if err := h.useCase.SendRequest(ctx, userID, req.ContactID); err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to add contact", err.Error())
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true})
}

// AcceptContact godoc
// @Summary Accept a contact request
// @Description Accepts a pending contact request from another user
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /contacts/{id} [patch]
func (h *ContactHandler) AcceptContact(ctx context.Context, c *app.RequestContext) {
	userID, ok := authctx.UserID(c)
	if !ok {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing user context", nil)
		return
	}

	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_PARAM", "invalid contact id", nil)
		return
	}

	if err := h.useCase.AcceptRequest(ctx, userID, contactID); err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to accept contact", err.Error())
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true})
}

// GetContacts godoc
// @Summary Get all contacts
// @Description Retrieves a list of accepted contacts for the current user
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /contacts [get]
func (h *ContactHandler) GetContacts(ctx context.Context, c *app.RequestContext) {
	userID, ok := authctx.UserID(c)
	if !ok {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing user context", nil)
		return
	}

	users, err := h.useCase.GetUserContacts(ctx, userID)
	if err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get contacts", err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"contacts": users,
	})
}

// RemoveContact godoc
// @Summary Remove a contact
// @Description Removes a contact from the current user's contact list
// @Tags contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Contact ID"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /contacts/{id} [delete]
func (h *ContactHandler) RemoveContact(ctx context.Context, c *app.RequestContext) {
	userID, ok := authctx.UserID(c)
	if !ok {
		api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing user context", nil)
		return
	}

	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		api.ErrorResponse(c, http.StatusBadRequest, "INVALID_PARAM", "invalid contact id", nil)
		return
	}

	if err := h.useCase.RemoveContact(ctx, userID, contactID); err != nil {
		api.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to remove contact", err.Error())
		return
	}

	c.JSON(http.StatusOK, api.Response{Success: true})
}
