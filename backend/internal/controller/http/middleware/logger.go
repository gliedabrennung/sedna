package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func RequestLogger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		latency := time.Since(start)
		
		status := c.Response.StatusCode()
		method := string(c.Request.Header.Method())
		path := string(c.Request.URI().Path())

		hlog.CtxInfof(ctx, "[HTTP] Method=%s Path=%s Status=%d Latency=%s IP=%s", method, path, status, latency, c.ClientIP())
	}
}
