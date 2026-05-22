package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
)

func TestServeHealth_StatusOK(t *testing.T) {
	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.GET("/health", ServeHealth)

	w := ut.PerformRequest(engine, http.MethodGet, "/health", nil)
	if got := w.Result().StatusCode(); got != http.StatusOK {
		t.Errorf("expected 200, got %d", got)
	}

	var resp healthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", resp.Status)
	}
}
