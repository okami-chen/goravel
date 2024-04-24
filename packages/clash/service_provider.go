package clash

import (
	"github.com/goravel/framework/contracts/foundation"
	"goravel/packages/clash/http/controllers"
)

const Binding = "clash"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return nil, nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {

	route := app.MakeRoute()
	clashCtr := controllers.ClashController{}
	route.Get("/clash", clashCtr.Index)
	route.Get("/bee", clashCtr.Index)

	qxCtr := controllers.QXController{}
	route.Get("/qx", qxCtr.Index)

	pingCtr := controllers.PingController{}
	route.Post("/ping", pingCtr.Index)
}
