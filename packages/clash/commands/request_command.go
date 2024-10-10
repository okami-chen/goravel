package commands

import (
	"context"
	"github.com/Dreamacro/clash/adapter"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/imroc/req/v3"
	"goravel/pkg/healthcheck"
	"goravel/pkg/tool"
	"net"
	"time"
)

type RequestCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *RequestCommand) Signature() string {
	return "clash:request"
}

// Description The console command description.
func (receiver *RequestCommand) Description() string {
	return "Command description"
}

// Extend The console command extend.
func (receiver *RequestCommand) Extend() command.Extend {
	return command.Extend{}
}

func dialContext(uri string, ctx context.Context) (net.Conn, error) {
	filePath := facades.Storage().Path("all.yaml")
	p := tool.RandomProxySimple(filePath)
	clash, _ := adapter.ParseProxy(p)
	addr, _ := healthcheck.UrlToMetadata(uri)
	return clash.DialContext(ctx, &addr)
}

// Handle Execute the console command.
func (receiver *RequestCommand) Handle(cx console.Context) error {
	uri := "https://www.futbin.com/25/player/49675/asier-villalibre-molina"
	//uri = "http://ip-api.com/json"
	uri = "https://tls.peet.ws/api/all"
	uri = "https://tls.browserleaks.com/json"
	//uri = "https://ascii2d.net/"
	client := req.C()
	client.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialContext(uri, ctx)
	}
	request := client.DevMode().
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36").
		SetTLSFingerprintChrome(). // 模拟 Chrome 浏览器的 TLS 握手指纹，让网站相信这是 Chrome 浏览器在访问，予以通行。
		ImpersonateChrome().
		EnableForceHTTP2().
		EnableDumpEachRequest().
		SetCommonRetryCount(5).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.StatusCode == http.StatusTooManyRequests
		}).
		SetCommonRetryHook(func(resp *req.Response, err error) {
			c := client.Clone()
			c.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialContext(uri, ctx)
			}
			resp.Request.SetClient(c) // Change the client of request dynamically.
		}).
		SetTimeout(time.Second * 5)

	request.R().MustGet(uri)
	return nil
}
