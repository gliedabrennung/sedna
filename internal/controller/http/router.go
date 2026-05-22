package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gliedabrennung/messenger-core/internal/common/api"
	"github.com/gliedabrennung/messenger-core/internal/controller/http/middleware"
	"github.com/gliedabrennung/messenger-core/internal/domain"
)

type Deps struct {
	Auth      AuthService
	Users     UserService
	MsgRepo   domain.MessageRepository
	WsHandler app.HandlerFunc
	JWTSecret string
	Cookie    CookieConfig
}

func SetupRouter(h *server.Hertz, deps Deps) {
	h.Use(api.CustomErrorHandler())

	h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.Path())
		if !strings.HasPrefix(path, "/auth") &&
			!strings.HasPrefix(path, "/users") &&
			!strings.HasPrefix(path, "/messages") &&
			!strings.HasPrefix(path, "/ws") &&
			!strings.HasPrefix(path, "/health") {
			c.File("./frontend/dist/index.html")
			return
		}
		api.ErrorResponse(c, http.StatusNotFound,
			"NOT_FOUND",
			"Page not found",
			nil)
	})

	h.NoMethod(func(ctx context.Context, c *app.RequestContext) {
		api.ErrorResponse(c, http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"Method not allowed",
			nil)
	})

	authHandler := NewAuthHandler(deps.Auth, deps.Cookie)
	userHandler := NewUserHandler(deps.Users)
	authLimiter := middleware.NewRateLimiter(5, 10)
	authMiddleware := middleware.JWTAuth(deps.JWTSecret, deps.Cookie.Name)

	h.GET("/health", ServeHealth)

	auth := h.Group("/auth")
	auth.Use(authLimiter.Handler())
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)

	users := h.Group("/users", authMiddleware)
	users.GET("/search", userHandler.Search)
	users.GET("/bulk", userHandler.GetBulk)

	msgHandler := NewMessageHandler(deps.MsgRepo)
	h.GET("/messages", authMiddleware, msgHandler.GetHistory)

	h.GET("/ws", authMiddleware, deps.WsHandler)

	h.StaticFS("/", &app.FS{
		Root:       "./frontend/dist",
		IndexNames: []string{"index.html"},
		PathNotFound: func(ctx context.Context, c *app.RequestContext) {
			path := string(c.Request.Path())
			if strings.HasPrefix(path, "/auth") ||
				strings.HasPrefix(path, "/users") ||
				strings.HasPrefix(path, "/messages") ||
				strings.HasPrefix(path, "/ws") ||
				strings.HasPrefix(path, "/health") {
				c.JSON(http.StatusNotFound, map[string]interface{}{
					"error_code": "NOT_FOUND",
					"message":    "Page not found",
				})
				return
			}
			c.File("./frontend/dist/index.html")
		},
	})
}
