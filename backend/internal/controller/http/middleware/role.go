package middleware

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gliedabrennung/messenger-core/internal/pkg/api"
	"github.com/gliedabrennung/messenger-core/internal/pkg/authctx"
	"github.com/gliedabrennung/messenger-core/internal/usecase"
)

func RoleAuth(authUseCase *usecase.AuthUseCase, allowedRole string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		userID, ok := authctx.UserID(c)
		if !ok {
			api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authentication context", nil)
			c.Abort()
			return
		}

		user, err := authUseCase.GetUserByID(ctx, userID)
		if err != nil || user == nil {
			api.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not found", nil)
			c.Abort()
			return
		}

		if user.Role != allowedRole {
			api.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "You don't have permission to access this resource", nil)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
