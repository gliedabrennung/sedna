package authctx

import (
	"github.com/cloudwego/hertz/pkg/app"
)

const userIDKey = "userID"

func SetUserID(c *app.RequestContext, id int64) {
	c.Set(userIDKey, id)
}

func UserID(c *app.RequestContext) (int64, bool) {
	v, ok := c.Get(userIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
