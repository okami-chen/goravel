package jwt

import (
	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "jwt"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Singleton(Binding, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		return NewJwt(config.GetString("auth.defaults.guard"),
			app.MakeCache(), config, app.MakeOrm()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {

}
