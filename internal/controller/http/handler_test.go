package http

import (
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
)

func TestServeHome_StatusOK(t *testing.T) {
	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.GET("/", ServeHome)

	w := ut.PerformRequest(engine, http.MethodGet, "/", nil)
	if got := w.Result().StatusCode(); got != http.StatusOK {
		t.Errorf("expected 200, got %d", got)
	}
}
