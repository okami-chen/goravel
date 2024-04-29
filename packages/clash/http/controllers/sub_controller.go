package controllers

import (
	"encoding/json"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"gopkg.in/yaml.v3"
	"goravel/app/models"
	"goravel/data"
	"goravel/packages/clash/services"
	"io"
	"os"
	"strings"
	"time"
)

type SubController struct {
	BaseController
}

func NewSubController() *SubController {
	return &SubController{}
}

type ClasNodeList struct {
	Proxies []map[string]interface{} `yaml:"proxies"`
}

func (r *SubController) Index(ctx http.Context) http.Response {
	r.Ctx = ctx
	cache := facades.Cache().WithContext(ctx)
	cacheKey := ctx.Request().Url()
	sec := ctx.Request().Input("key")
	cls := ctx.Request().Input("cls")
	if cls != "" {
		cacheKey = strings.Replace("&cls="+cls, "", cacheKey, -1)
	}
	contentType := "text/plain;charset=utf-8"
	if sec == "" || sec != "123" {
		return ctx.Response().Data(403, contentType, nil)
	}
	request := ctx.Request()
	in := request.Input("in")
	out := request.Input("out")
	if cls == "" && cache.Has(cacheKey) {
		return ctx.Response().
			Header("subscription-userinfo", r.getSubInfo(in)).
			Data(200, contentType, []byte(cache.Get(cacheKey).(string)))
	}
	proxies := make([]models.Proxy, 0)
	if strings.Contains(cacheKey, "bee") {
		in = "a.e.m.g"
	}
	q := request.Input("q")
	s := request.Input("s")

	if q != "" {
		emojis := services.FindEmojiByCode(q, ctx)
		proxies = services.List(emojis, in, out, ctx)
		proxies = services.SortByEmoji(emojis, proxies)
	}
	if s != "" {
		emojis := services.FindEmojiByCode(s, ctx)
		proxies = services.List(nil, in, out, ctx)
		proxies = services.SortByEmoji(emojis, proxies)
	}
	if q == "" && s == "" {
		proxies = services.List(nil, in, out, ctx)
	}

	flag := request.Input("flag")
	agent := strings.ToLower(request.Input("agent"))

	if strings.Contains(agent, "clash") {
		flag = "clash"
	}
	if strings.Contains(agent, "quantumult") {
		flag = "qx"
	}
	resp := ""
	if flag == "qx" || flag == "quantumultx" {
		resp = r.getQuantumultX(proxies, "")
		cache.Put(cacheKey, resp, 600)
	} else if flag == "node" {
		nodeList := ClasNodeList{}
		for _, item := range proxies {
			var ret map[string]interface{}
			json.Unmarshal([]byte(item.Body), &ret)
			nodeList.Proxies = append(nodeList.Proxies, ret)
		}
		bt, _ := yaml.Marshal(nodeList)
		cache.Put(cacheKey, string(bt), time.Minute*60)
		resp = string(bt)
	} else {
		fileName := "storage/clash/clash_v6.yaml"
		if request.Input("l") != "" {
			fileName = "storage/clash/clash_" + request.Input("l") + ".yaml"
		}
		clashYaml := data.ClashYaml{}
		file, _ := os.Open(fileName)
		defer file.Close()
		// 读取文件内容
		content, err := io.ReadAll(file)
		if err != nil {
			return ctx.Response().Success().Data("text/plain;charset=utf-8", []byte(""))
		}
		yaml.Unmarshal(content, &clashYaml)
		clashYaml = r.getClash(clashYaml, proxies)
		bt, _ := yaml.Marshal(clashYaml)
		err = cache.Put(cacheKey, string(bt), time.Minute*60)
		resp = string(bt)
	}

	return ctx.Response().
		Header("subscription-userinfo", r.getSubInfo(in)).
		Data(200, contentType, []byte(resp))
}
