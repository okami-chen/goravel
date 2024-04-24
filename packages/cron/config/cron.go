package config

import (
	"github.com/goravel/framework/facades"
	"goravel/packages/cron/job"
)

func init() {
	config := facades.Config()
	config.Add("cron", map[string]any{
		"test": &job.CronTestHandle{},
	})
}
