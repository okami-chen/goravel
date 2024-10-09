package controllers

import (
	"encoding/json"
	"github.com/Dreamacro/clash/adapter"
	"github.com/gin-gonic/gin"
	resp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/wangluozhe/requests"
	"github.com/wangluozhe/requests/url"
	"goravel/pkg/healthcheck"
	"goravel/pkg/tool"
	"strings"
	"time"
)

func mapToJson(data map[string]any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)

}

type ProxyController struct {
	//Dependent services
}

func NewProxyController() *ProxyController {
	return &ProxyController{
		//Inject services
	}
}

func (r *ProxyController) Index(ctx resp.Context) resp.Response {

	if ctx.Request().Input("key") != "123" {
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"message": "秘钥错误",
		})
	}

	uri := ctx.Request().Input("url")
	if uri == "" {
		uri = "http://ip-api.com/json"
	}

	filePath := facades.Storage().Path("all.yaml")
	p := tool.RandomProxySimple(filePath)

	facades.Log().Info(mapToJson(gin.H{
		"Name":   p["name"].(string),
		"Type":   p["type"].(string),
		"Server": p["server"].(string),
		"Port":   p["port"].(int),
		"Url":    uri,
	}))

	clash, err := adapter.ParseProxy(p)

	if err != nil {
		facades.Log().Error(mapToJson(gin.H{
			"Name":   p["name"].(string),
			"Server": p["server"].(string),
			"Url":    uri,
			"Error":  err.Error(),
		}))
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	}

	addr, err := healthcheck.UrlToMetadata(uri)
	if err != nil {
		facades.Log().Error(mapToJson(gin.H{
			"Name":   p["name"].(string),
			"Server": p["server"].(string),
			"Url":    uri,
			"Error":  err.Error(),
		}))
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"message": err,
		})
	}
	conn, err := clash.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"message": err,
		})
	}
	defer conn.Close()
	session := requests.NewSession(conn)
	req := url.NewRequest()
	headers := url.NewHeaders()
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Headers = headers
	req.Timeout = time.Second * 5
	req.Ja3 = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
	req.Ja3 = "772,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,43-11-65281-5-18-10-51-65037-0-23-13-17513-16-27-45-35,25497-29-23-24,0"
	ret, err := session.Get(uri, req)
	if err != nil {
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"message": err,
		})
	}

	if strings.Contains(ret.Text, "403 Forbidden") {
		return ctx.Response().Json(200, gin.H{
			"code":    500,
			"name":    clash.Name(),
			"message": "ip被封",
			"proxy": map[string]interface{}{
				"name":   p["name"].(string),
				"type":   p["type"].(string),
				"server": p["server"].(string),
				"port":   p["port"].(int),
			},
			"data": "",
		})
	}

	return ctx.Response().Json(200, gin.H{
		"code": 0,
		"j3":   req.Ja3,
		"proxy": map[string]interface{}{
			"name":   p["name"].(string),
			"type":   p["type"].(string),
			"server": p["server"].(string),
			"port":   p["port"].(int),
		},
		"data": ret.Text,
	})
}
