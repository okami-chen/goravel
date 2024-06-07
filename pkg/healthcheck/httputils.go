package healthcheck

import (
	"context"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"goravel/pkg/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func UrlToMetadata(rawURL string) (addr C.Metadata, err error) {
	return urlToMetadata(rawURL)
}

// DO NOT EDIT. Copied from clash because it's an unexported function
func urlToMetadata(rawURL string) (addr C.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}

	portInt, err := strconv.ParseUint(port, 10, 16)

	addr = C.Metadata{
		Host:    u.Hostname(),
		DstIP:   nil,
		DstPort: C.Port(portInt),
	}
	return
}

func HTTPGetViaProxy(clashProxy C.Proxy, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DelayTimeout)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", getRandomUserAgent())
	req = req.WithContext(ctx)

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func HTTPHeadViaProxy(clashProxy C.Proxy, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DelayTimeout)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", getRandomUserAgent())
	req = req.WithContext(ctx)

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("%d %s for proxy %s %s", resp.StatusCode, resp.Status, clashProxy.Name(), clashProxy.Addr())
	}
	resp.Body.Close()
	return nil
}

func HTTPGetBodyViaProxy(clashProxy C.Proxy, url string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), DelayTimeout)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return nil, err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", getRandomUserAgent())
	req = req.WithContext(ctx)

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read speedtest config file
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HTTPGetBodyViaProxyWithTimeRetry(proxy C.Proxy, url string, timeout time.Duration, retries int) ([]byte, error) {
	var err error
	var resp []byte
	for i := 0; i < retries; i++ {
		resp, err = HTTPGetBodyViaProxyWithTime(proxy, url, timeout)
		if err == nil {
			return resp, nil
		}
		time.Sleep(time.Second) // 重试之前等待一秒钟
	}

	return nil, err
}

func HTTPGetBodyViaProxyWithTime(clashProxy C.Proxy, url string, t time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return nil, err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", getRandomUserAgent())
	cookie := &http.Cookie{
		Name:  "g2g_regional",
		Value: "%7B%22country%22%3A%22US%22%2C%22currency%22%3A%22USD%22%2C%22language%22%3A%22en%22%7D",
	}
	req.Header.Add("Cookie", cookie.String())
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read speedtest config file
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Get body without return to save memory
func HTTPGetBodyViaProxyWithTimeNoReturn(clashProxy C.Proxy, url string, t time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", getRandomUserAgent())

	transport := &http.Transport{
		// Note: Dial specifies the dial function for creating unencrypted TCP connections.
		// When httpClient sets this transport, it will use the tcp/udp connection returned from
		// function Dial instead of default tcp/udp connection. It's the key to set custom proxy for http transport
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read speedtest config file
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func HTTPostBodyViaProxy(clashProxy C.Proxy, url string, param io.Reader) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), DelayTimeout)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return nil, err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodPost, url, param)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", getRandomUserAgent())
	req = req.WithContext(ctx)

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read speedtest config file
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HTTPPostBodyViaProxyWithTime(clashProxy C.Proxy, url string, t time.Duration, param io.Reader) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	addr, err := urlToMetadata(url)
	if err != nil {
		return nil, err
	}
	conn, err := clashProxy.DialContext(ctx, &addr) // 建立到proxy server的connection，对Proxy的类别做了自适应相当于泛型
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req, err := http.NewRequest(http.MethodPost, url, param)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", getRandomUserAgent())
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read speedtest config file
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func checkErrorProxies(proxies []proxy.Proxy) bool {
	if proxies == nil {
		return false
	}
	if len(proxies) == 0 {
		return false
	}
	if proxies[0] == nil {
		return false
	}
	return true
}
