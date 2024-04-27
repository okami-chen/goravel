package routes

import (
	"github.com/goravel/framework/facades"
)

func Web() {
	facades.Route().StaticFile("/favicon.ico", "./public/favicon.ico")
}
