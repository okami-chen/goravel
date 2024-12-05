package clash

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	router "github.com/goravel/framework/contracts/route"
	"goravel/packages/clash/commands"
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

	registerCommand()
	registerRouter()
}

func registerCommand() {
	App.MakeArtisan().Register([]console.Command{
		&commands.RequestCommand{},
	})
}

func registerRouter() {

	route := App.MakeRoute()

	route.Prefix("/geo").Group(func(route router.Router) {
		infoCtr := controllers.NewInfoController()
		route.Get("/info", infoCtr.Index)

		locCtr := controllers.NewLocController()
		route.Get("/loc", locCtr.Index)

		hubCtr := controllers.NewHubController()
		route.Get("/hub", hubCtr.Index)

		geoCtr := controllers.NewGeoController()
		route.Get("/base", geoCtr.Index)

		scoreCtr := controllers.NewScoreController()
		route.Get("/score", scoreCtr.Index)

		ipCtr := controllers.NewIpController()
		route.Get("/ip", ipCtr.Index)
	})

	subCtr := controllers.NewSubController()
	route.Get("/sub", subCtr.Index)

	pingCtr := controllers.NewPingController()
	route.Post("/ping", pingCtr.Index)

	proxyCtr := controllers.NewProxyController()
	route.Get("/proxy", proxyCtr.Index)
	route.Post("/proxy", proxyCtr.Index)
}