package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"goravel/packages/clash/services"
	"strings"
)

type QXController struct {
	//Dependent services
}

func NewQXController() *QXController {
	return &QXController{
		//Inject services
	}
}

func (r *QXController) Index(ctx http.Context) http.Response {
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
	c := request.Input("c")
	if cls == "" && cache.Has(cacheKey) {
		return ctx.Response().
			Header("subscription-userinfo", getSubscriptInfo(c)).
			Data(200, contentType, []byte(cache.Get(cacheKey).(string)))
	}
	proxies := make([]models.Proxy, 0)
	q := request.Input("q")
	if q != "" {
		emojis := services.FindEmojiByCode(q, ctx)
		proxies = services.List(emojis, c, ctx)
		proxies = services.SortByEmoji(emojis, proxies)
	}
	s := request.Input("s")
	if s != "" {
		emojis := services.FindEmojiByCode(s, ctx)
		proxies = services.List(nil, c, ctx)
		proxies = services.SortByEmoji(emojis, proxies)
	}
	if q == "" && s == "" {
		proxies = services.List(nil, c, ctx)
	}
	resp := getQuantumultX(proxies, "")
	cache.Put(cacheKey, resp, 600)
	return ctx.Response().
		Header("subscription-userinfo", getSubscriptInfo(c)).
		Data(200, contentType, []byte(resp))
}
