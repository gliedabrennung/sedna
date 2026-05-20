package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
)

func TestRateLimiter_AllowsBurst(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.GET("/test", rl.Handler(), func(ctx context.Context, c *app.RequestContext) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 5; i++ {
		w := ut.PerformRequest(engine, http.MethodGet, "/test", nil)
		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, w.Code)
		}
	}
}

func TestRateLimiter_BlocksAfterBurst(t *testing.T) {
	rl := NewRateLimiter(0.1, 2)

	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.GET("/test", rl.Handler(), func(ctx context.Context, c *app.RequestContext) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		w := ut.PerformRequest(engine, http.MethodGet, "/test", nil)
		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i, w.Code)
		}
	}

	w := ut.PerformRequest(engine, http.MethodGet, "/test", nil)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}
