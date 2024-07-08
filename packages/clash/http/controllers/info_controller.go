package controllers

import (
	"encoding/json"
	"github.com/Dreamacro/clash/adapter"
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"goravel/pkg/healthcheck"
	"time"
)

type InfoController struct {
	//Dependent services
}

func NewInfoController() *InfoController {
	return &InfoController{
		//Inject services
	}
}

func (r *InfoController) Index(ctx http.Context) http.Response {
	request := ctx.Request()
	id := request.Input("id")
	if id == "" {
		return ctx.Response().Json(500, gin.H{
			"message": " error",
		})
	}
	var obj models.Proxy

	query := facades.Orm().WithContext(ctx).Query()
	query.WithTrashed().Where("id = ?", id).Find(&obj)

	if obj.Id == 0 {
		return ctx.Response().Json(500, gin.H{
			"message": "Not Found",
		})
	}
	var p map[string]interface{}
	json.Unmarshal([]byte(obj.Body), &p)
	px, err := adapter.ParseProxy(p)
	if err != nil {
		return ctx.Response().Json(500, gin.H{
			"message": err.Error(),
		})
	}

	url := "http://ip-api.com/json/?lang=zh-CN&fields=status,message,continent,continentCode,country,countryCode"
	url = url + ",region,regionName,city,district,zip,lat,lon,timezone,offset,currency,isp,org,as,asname,reverse"
	url = url + ",mobile,proxy,hosting,query"
	body, err := healthcheck.HTTPGetBodyViaProxyWithTime(px, url, time.Second*10)
	if err != nil {
		return ctx.Response().Json(500, gin.H{
			"message": err.Error(),
		})
	}

	return ctx.Response().Data(http.StatusOK, "application/json;charset=utf-8", body)

}
