package cloudflare

import (
	"encoding/json"
	"fmt"
	"goravel/app/services/cloudflare/types"
	"goravel/pkg/tool"
	"strings"
)

type Waf struct {
	Zone  string
	Token struct {
		Update string
		Read   string
	}
	Host string
}

func NewWAF(zone, readToken, updateToken string) *Waf {
	return &Waf{
		Zone: zone,
		Token: struct {
			Update string
			Read   string
		}{
			Update: updateToken,
			Read:   readToken,
		},
		Host: "https://api.cloudflare.com/client/v4/zones",
	}
}

func (r Waf) getEntrypointUrl(ruleSet, ruleId string) string {
	return fmt.Sprintf("%s/%s/rulesets/phases/%s/entrypoint", r.Host, ruleSet, ruleId)
}

func (r Waf) getRuleList(ruleSet, ruleId string) (*types.Rule, error) {
	h := make(map[string]string)
	h["Authorization"] = fmt.Sprintf("Bearer %s", r.Token.Read)
	h["Content-Type"] = fmt.Sprintf("application/json")
	data := make(map[string]interface{})
	url := r.getEntrypointUrl(ruleSet, ruleId)

	resp, err := tool.Request("PATCH", url, data, h)
	if err != nil {
		return nil, err
	}
	var ret map[string]interface{}
	err = json.NewDecoder(strings.NewReader(string(resp))).Decode(&ret)

	rule := types.Rule{}
	// 打印特定规则的信息
	rules := ret["result"].(map[string]interface{})["rules"].([]interface{})
	for _, rs := range rules {
		ruleMap := rs.(map[string]interface{})
		if ruleMap["id"] == ruleId {
			rule.ID = ruleMap["id"].(string)
			rule.Action = "managed_challenge"
			rule.Expression = ruleMap["expression"].(string)
			rule.Description = ruleMap["description"].(string)
			rule.Version = ruleMap["version"].(string)
			rule.Ref = ruleMap["ref"].(string)
			rule.LastUpdated = ruleMap["last_updated"].(string)
			rule.Enabled = ruleMap["enabled"].(bool)
			break
		}
	}
	return &rule, nil
}
