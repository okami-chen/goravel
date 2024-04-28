package clash

import (
	"github.com/goravel/framework/contracts/foundation"
	router "github.com/goravel/framework/contracts/route"
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

	route.Prefix("/geo").Group(func(route router.Router) {
		infoCtr := controllers.InfoController{}
		route.Get("/info", infoCtr.Index)

		locCtr := controllers.LocController{}
		route.Get("/loc", locCtr.Index)
	})

	clashCtr := controllers.ClashController{}
	route.Get("/clash", clashCtr.Index)
	route.Get("/bee", clashCtr.Index)

	nodeCtr := controllers.NodeController{}
	route.Get("/node", nodeCtr.Index)

	qxCtr := controllers.QXController{}
	route.Get("/qx", qxCtr.Index)

	subCtr := controllers.SubController{}
	route.Get("/sub", subCtr.Index)

	pingCtr := controllers.PingController{}
	route.Get("/ping", pingCtr.Index)

}
