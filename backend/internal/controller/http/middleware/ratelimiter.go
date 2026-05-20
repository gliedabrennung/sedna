package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gliedabrennung/messenger-core/internal/pkg/api"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var rl = &rateLimiter{
	visitors: make(map[string]*visitor),
}

func init() {
	go cleanupVisitors()
}

// cleanupVisitors acts as a background worker
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(time.Second), 20) // 20 requests per second
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func RateLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			api.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests, please try again later", nil)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
