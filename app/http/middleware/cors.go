package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

func Cors() http.Middleware {
	return func(ctx http.Context) {
		if ctx.Request().Method() == "OPTIONS" {
			ctx.Request().AbortWithStatus(204)
			return
		}
		ctx.Request().Next()
	}
}
