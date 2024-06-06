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
		infoCtr := controllers.NewInfoController()
		route.Get("/info", infoCtr.Index)

		locCtr := controllers.NewLocController()
		route.Get("/loc", locCtr.Index)

		geoCtr := controllers.NewGeoController()
		route.Get("/base", geoCtr.Index)

		scoreCtr := controllers.NewScoreController()
		route.Get("/score", scoreCtr.Index)
	})

	subCtr := controllers.NewSubController()
	route.Get("/sub", subCtr.Index)

	pingCtr := controllers.NewPingController()
	route.Post("/ping", pingCtr.Index)

}
