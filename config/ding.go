package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("ding", map[string]any{
		"default": map[string]any{
			"token":  config.Env("DING_DEFAULT_TOKEN", "0017582cd374573718939b3ba9f6fb0c62f6a5653db4b72e5dce2395f078ee14"),
			"secret": config.Env("DING_DEFAULT_SECRET", "SEC3e377c16b2ab637524ae9ace0348509270626427e6aa53a9fbb9a6a510184d1b"),
		},
	})
}
