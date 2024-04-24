package facades

import (
	"log"

	"goravel/packages/cron"
	"goravel/packages/cron/contracts"
)

func Cron() contracts.Cron {
	instance, err := cron.App.Make(cron.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Cron)
}

func Crontab() contracts.Cron {
	instance, err := cron.App.Make(cron.Binding + ".core")
	if err != nil {
		log.Println(err)
		return nil
	}
	return instance.(contracts.Cron)
}
