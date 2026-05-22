package http

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

type healthResponse struct {
	Status string `json:"status"`
}

func ServeHealth(_ context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
