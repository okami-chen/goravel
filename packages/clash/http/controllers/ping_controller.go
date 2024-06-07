package controllers

import (
	"github.com/Dreamacro/clash/adapter"
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/contracts/http"
	"gopkg.in/yaml.v3"
	"goravel/pkg/healthcheck"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PingController struct {
}

func NewPingController() *PingController {
	return &PingController{
		//Inject services
	}
}

func (r *PingController) Index(ctx http.Context) http.Response {
	request := ctx.Request()
	if request.Input("cfg") == "" || request.Input("url") == "" {
		return ctx.Response().Success().Json(gin.H{
			"success":  false,
			"message":  "cfg -> is nil",
			"response": "",
		})
	}

	p, cur, total := PingRandomProxy(request.Input("cfg"))
	proxy, err := adapter.ParseProxy(p)
	if err != nil {
		ctx.Response().Success().Json(gin.H{
			"success": false,
			"proxy": map[string]interface{}{
				"server": p["name"].(string),
				"type":   p["type"].(string),
				"total":  total,
				"random": cur,
			},
			"message":  "proxy -> " + err.Error(),
			"response": "",
		})
	}

	resp, err := healthcheck.HTTPGetBodyViaProxyWithTimeRetry(proxy, request.Input("url"), time.Second*5, 2)
	if err != nil {
		return ctx.Response().Success().Json(gin.H{
			"success": false,
			"proxy": map[string]interface{}{
				"server": p["name"].(string),
				"type":   p["type"].(string),
				"total":  total,
				"random": cur,
			},
			"message":  "check -> " + err.Error(),
			"response": "",
		})
	}
	return ctx.Response().Success().Json(gin.H{
		"success": true,
		"proxy": map[string]interface{}{
			"server": p["name"].(string),
			"type":   p["type"].(string),
			"total":  total,
			"random": cur,
		},
		"message":  "ok",
		"response": string(resp),
	})
}

type ClashYaml struct {
	Proxy []map[string]interface{} `json:"proxies" yaml:"proxies"`
}

func PingRandomProxy(cfg string) (str map[string]interface{}, cur int, total int) {
	clash := ClashYaml{}
	filepath.Walk(cfg, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			content, _ := os.ReadFile(path)
			p := ClashYaml{}
			yaml.Unmarshal(content, &p)
			clash.Proxy = append(clash.Proxy, p.Proxy...)
		}
		return nil
	})
	l := len(clash.Proxy)
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(l)
	if randomInt >= l {
		randomInt = l - 1
	}
	return clash.Proxy[randomInt], randomInt, l
}
