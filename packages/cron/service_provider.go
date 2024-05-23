package cron

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/robfig/cron/v3"
	"goravel/packages/cron/console/commands"
	"goravel/packages/cron/job"
)

const Binding = "cron.core"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return nil, nil
	})

	app.Singleton(Binding, func(app foundation.Application) (any, error) {
		secondParser := cron.NewParser(cron.Second | cron.Minute |
			cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
		return cron.New(cron.WithParser(secondParser), cron.WithChain()), nil
	})
	job.InitJob()
	facades.Validation()
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Commands([]console.Command{
		&commands.TestCommand{},
	})
	job.Setup()
	c, _ := app.Make(Binding)
	c.(*cron.Cron).Start()

}
