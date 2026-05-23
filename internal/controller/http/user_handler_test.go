package http

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/gliedabrennung/sedna/internal/common/authctx"
	"github.com/gliedabrennung/sedna/internal/entity"
)

type mockUserService struct {
	searchUsers   func(ctx context.Context, query string) ([]entity.User, error)
	getUsersByIDs func(ctx context.Context, ids []int64) ([]entity.User, error)
}

func (m *mockUserService) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	if m.searchUsers != nil {
		return m.searchUsers(ctx, query)
	}
	return nil, nil
}

func (m *mockUserService) GetUsersByIDs(ctx context.Context, ids []int64) ([]entity.User, error) {
	if m.getUsersByIDs != nil {
		return m.getUsersByIDs(ctx, ids)
	}
	return nil, nil
}

func TestUserHandler_Search(t *testing.T) {
	mockSvc := &mockUserService{
		searchUsers: func(ctx context.Context, query string) ([]entity.User, error) {
			return []entity.User{{ID: 1, Username: "test"}}, nil
		},
	}
	h := NewUserHandler(mockSvc)
	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.GET("/search", h.Search)

	t.Run("Success", func(t *testing.T) {
		w := ut.PerformRequest(engine, http.MethodGet, "/search?q=test", nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var users []entity.User
		json.Unmarshal(w.Body.Bytes(), &users)
		if len(users) != 1 || users[0].Username != "test" {
			t.Errorf("unexpected response: %v", string(w.Body.Bytes()))
		}
	})

	t.Run("NoQuery", func(t *testing.T) {
		w := ut.PerformRequest(engine, http.MethodGet, "/search", nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestUserHandler_GetBulk(t *testing.T) {
	mockSvc := &mockUserService{
		getUsersByIDs: func(ctx context.Context, ids []int64) ([]entity.User, error) {
			return []entity.User{{ID: 1, Username: "test"}, {ID: 2, Username: "test2"}}, nil
		},
	}
	h := NewUserHandler(mockSvc)
	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.GET("/bulk", h.GetBulk)

	t.Run("Success", func(t *testing.T) {
		w := ut.PerformRequest(engine, http.MethodGet, "/bulk?ids=1,2", nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var users []entity.User
		json.Unmarshal(w.Body.Bytes(), &users)
		if len(users) != 2 {
			t.Errorf("unexpected response: %v", string(w.Body.Bytes()))
		}
	})

	t.Run("NoIDs", func(t *testing.T) {
		w := ut.PerformRequest(engine, http.MethodGet, "/bulk", nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestUserHandler_Me(t *testing.T) {
	mockSvc := &mockUserService{
		getUsersByIDs: func(ctx context.Context, ids []int64) ([]entity.User, error) {
			return []entity.User{{ID: 1, Username: "test"}}, nil
		},
	}
	h := NewUserHandler(mockSvc)
	engine := route.NewEngine(config.NewOptions([]config.Option{}))
	engine.Use(func(c context.Context, ctx *app.RequestContext) {
		if string(ctx.Request.URI().Path()) == "/me_auth" {
			authctx.SetUserID(ctx, 1)
		}
		ctx.Next(c)
	})
	engine.GET("/me", h.Me)
	engine.GET("/me_auth", h.Me)

	t.Run("Unauthorized", func(t *testing.T) {
		w := ut.PerformRequest(engine, http.MethodGet, "/me", nil)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("Success", func(t *testing.T) {
		w := ut.PerformRequest(engine, http.MethodGet, "/me_auth", nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var user entity.User
		json.Unmarshal(w.Body.Bytes(), &user)
		if user.ID != 1 {
			t.Errorf("unexpected response: %v", string(w.Body.Bytes()))
		}
	})
}
