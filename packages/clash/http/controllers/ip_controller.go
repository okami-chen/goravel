package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/oschwald/geoip2-golang"
	"github.com/xiaoqidun/qqwry"
	"net"
)

type IpController struct {
	//Dependent services
}

func NewIpController() *IpController {
	return &IpController{
		//Inject services
	}
}

func (r *IpController) Index(ctx http.Context) http.Response {
	request := ctx.Request()
	ip := request.Ip()
	filePath := facades.Storage().Path("qqwry.dat")
	if err := qqwry.LoadFile(filePath); err != nil {
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"ip":      ip,
			"message": err.Error(),
		})
	}

	// 从内存或缓存查询IP
	location, err := qqwry.QueryIP(ip)
	if err != nil {
		return ctx.Response().Json(200, gin.H{
			"code":    501,
			"ip":      ip,
			"message": err.Error(),
		})
	}

	hmap := make(map[string]string)
	hmap["src"] = request.Input("src", "wan1")
	hmap["ip"] = ip
	hmap["country"] = location.Country
	hmap["province"] = location.Province
	hmap["city"] = location.City
	hmap["isp"] = location.ISP

	bt, _ := json.Marshal(hmap)

	facades.Log().Infof("%s", string(bt))

	return ctx.Response().Json(200, gin.H{
		"code":    0,
		"ip":      ip,
		"message": location.Province + " - " + location.City,
	})
}

func (r *IpController) IndexOld(ctx http.Context) http.Response {
	request := ctx.Request()
	ip := request.Ip()
	facades.Log().Infof("src: %s, ip: %s", request.Input("src"), ip)
	filePath := facades.Storage().Path("Country.mmdb")
	db, err := geoip2.Open(filePath)
	if err != nil {
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"ip":      ip,
		})
	}
	defer db.Close()

	//Parse Ip
	parse := net.ParseIP("183.160.253.249")
	record, err := db.City(parse)
	if err != nil {
		return ctx.Response().Json(200, gin.H{
			"code":    501,
			"message": err.Error(),
			"ip":      ip,
		})
	}
	facades.Log().Infof("%#v", record.City)

	return ctx.Response().Json(200, gin.H{
		"code":    0,
		"ip":      ip,
		"message": record.City.Names["zh-CN"],
	})
}
