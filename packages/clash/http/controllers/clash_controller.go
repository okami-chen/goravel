package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"gopkg.in/yaml.v3"
	"goravel/app/models"
	"goravel/data"
	"goravel/pkg/tool"
	"io"
	"os"
	"strings"
	"time"
)

type ClashController struct {
	//Dependent services
}

var countries = []string{"美国", "英国", "法国", "德国", "瑞典", "香港", "台湾", "日本", "韩国", "日本", "家宽", "新加坡"}

func NewClashController() *ClashController {
	return &ClashController{
		//Inject services
	}
}

func (r *ClashController) Index(ctx http.Context) http.Response {
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
	t := request.Input("t")
	q := request.Input("q")
	// 指定路由
	if strings.Contains(cacheKey, "bee") {
		c = "a.e"
	}

	if cls == "" && cache.Has(cacheKey) {
		return ctx.Response().
			Header("subscription-userinfo", getSubscriptInfo(c)).
			Data(200, contentType, []byte(cache.Get(cacheKey).(string)))
	}
	query := facades.Orm().WithContext(ctx).Query()

	var obj []models.Proxy
	query = buildCondition("title", tool.GetGroupReplaces(), t, query)
	query = buildCondition("code", tool.GetGroupReplaces(), c, query)
	if q != "" {
		query = buildNameQuery(q, query)
	}
	tag := request.Input("tag")
	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	query.Order("name asc").Find(&obj)

	// template
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

	s := ctx.Request().Input("s")
	if s == "" {
		s = q
	}
	if s != "" {
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
			for _, v4 := range mGroup[v3] {
				clashYaml = processProxy(clashYaml, v4)
			}
		}
		// 其他
		if len(mOther) > 0 {
			for _, v5 := range mOther {
				clashYaml = processProxy(clashYaml, v5)
			}
		}

	} else {
		for _, v := range obj {
			clashYaml = processProxy(clashYaml, v)
		}
	}

	bt, _ := yaml.Marshal(clashYaml)
	err = cache.Put(cacheKey, string(bt), time.Minute*60)
	if err != nil {
		facades.Log().Error(err.Error())
	}
	return ctx.Response().
		Header("subscription-userinfo", getSubscriptInfo(request.Input("c"))).
		Data(200, contentType, bt)
}
