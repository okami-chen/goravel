package controllers

import (
	"encoding/json"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"goravel/pkg/proxy"
	"goravel/pkg/tool"
	"strings"
	"time"
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
	if cls == "" && cache.Has(cacheKey) {
		return ctx.Response().
			Header("subscription-userinfo", getSubscriptInfo(request.Input("c"))).
			Data(200, contentType, []byte(cache.Get(cacheKey).(string)))
	}
	query := facades.Orm().WithContext(ctx).Query()

	var obj []models.Proxy
	//query = buildCondition("name", tool.GetNameReplaces(), request.Input("q"), query)
	query = buildCondition("title", tool.GetGroupReplaces(), request.Input("t"), query)
	query = buildCondition("code", tool.GetGroupReplaces(), strings.ToUpper(request.Input("c")), query)
	query = buildNameQuery(request.Input("q"), query)

	tag := request.Input("tag")
	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	query.Order("name asc").Find(&obj)
	resp := ""
	s := ctx.Request().Input("s")
	if s == "" {
		s = ctx.Request().Input("q")
	}
	if s == "" {
		for _, item := range obj {
			var ret map[string]interface{}
			json.Unmarshal([]byte(item.Body), &ret)
			px, _ := proxy.ParseProxyFromClashProxy(ret)
			resp = resp + px.ToQuantumultX() + "\n"
		}
	} else {
		var mGroup = make(map[string][]models.Proxy)
		var mOther = make([]models.Proxy, 0)
		var emoji []models.Emoji
		t2 := strings.Split(s, ".")
		em := facades.Orm().WithContext(ctx).Query()
		c2 := models.Condition{}
		em.Where(c2.ConditionsEqOr("code", t2)).Find(&emoji)

		hmap := make(map[string]string)

		for _, s3 := range emoji {
			hmap[s3.Code] = s3.Country
		}

		for i2, s2 := range t2 {
			if _, ok := hmap[s2]; ok {
				t2[i2] = hmap[s2]
			}
		}

		for _, v2 := range obj {
			var found bool
			for _, v1 := range t2 {
				if strings.Contains(v2.Name, v1) {
					mGroup[v1] = append(mGroup[v1], v2)
					found = true
					break
				}
			}
			if !found {
				mOther = append(mOther, v2)
			}
		}
		// 分组
		for _, v3 := range t2 {
			//不存在下标
			if _, ok := mGroup[v3]; !ok {
				continue
			}
			resp = getQuantumultX(mGroup[v3], resp)
		}
		// 其他
		if len(mOther) > 0 {
			resp = getQuantumultX(mOther, resp)
		}
	}

	cache.Put(cacheKey, resp, time.Minute*60)
	return ctx.Response().
		Header("subscription-userinfo", getSubscriptInfo(request.Input("c"))).
		Data(200, contentType, []byte(resp))
}
