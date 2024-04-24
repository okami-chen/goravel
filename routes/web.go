package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
	dir "net/http"
)

func Web() {
	facades.Route().StaticFS("/font", dir.Dir("./public/font"))
	facades.Route().StaticFS("/js", dir.Dir("./public/js"))
	facades.Route().StaticFS("/css", dir.Dir("./public/css"))
	facades.Route().StaticFS("/image", dir.Dir("./public/image"))
	facades.Route().StaticFS("/upload", dir.Dir("./public/upload"))
	facades.Route().StaticFS("/cdn-cgi", dir.Dir("./public/cdn-cgi"))
	facades.Route().StaticFile("/favicon.ico", "./public/favicon.ico")
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("index.tmpl", map[string]any{
			"version": support.Version,
		})
	})
}
