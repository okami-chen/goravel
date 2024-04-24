package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("cloudflare", map[string]any{
		"default": map[string]any{
			"zone":     config.Env("CLOUDFLARE_DEFAULT_ZONE", "fabdeca4e6f947479bec9b47f9497f99"),
			"rule_set": config.Env("CLOUDFLARE_DEFAULT_RULE_ID", "539b9265146c4619b6e5ef0a62030c47"),
			"rule_id":  config.Env("CLOUDFLARE_DEFAULT_RULE_ID", "a0634671f23d4767b5ddcc9679b55646"),
			"token": map[string]any{
				"read":   config.Env("CLOUDFLARE_DEFAULT_TOKEN_READ", "iExtktyzu6cXdbcRVQA0N58EvTEM0gm7xsIGrI-k"),
				"update": config.Env("CLOUDFLARE_DEFAULT_TOKEN_READ", "eK-qaOsA4ooYMpPhQpCcmmwXmjDBko8XacBjsoyU"),
			},
		},
	})
}
