package healthcheck

import (
	"context"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"goravel/pkg/proxy"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
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

	addr = C.Metadata{
		Host:    u.Hostname(),
		DstIP:   nil,
		DstPort: port,
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
	req = req.WithContext(ctx)

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
	req = req.WithContext(ctx)

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
	if resp.StatusCode >= 400 {
		return fmt.Errorf("%d %s for proxy %s %s", resp.StatusCode, resp.Status, clashProxy.Name(), clashProxy.Addr())
	}
	resp.Body.Close()
	return nil
}

func getRandomUserAgent() string {
	var userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1) AppleWebKit/533.34.1 (KHTML, like Gecko) Version/5.0 Safari/533.34.1",
		"Mozilla/5.0 (Windows; U; Windows NT 5.0) AppleWebKit/533.50.5 (KHTML, like Gecko) Version/4.0.3 Safari/533.50.5",
		"Mozilla/5.0 (Windows; U; Windows NT 6.0) AppleWebKit/531.31.7 (KHTML, like Gecko) Version/4.0.3 Safari/531.31.7",
		"Mozilla/5.0 (compatible; MSIE 7.0; Windows CE; Trident/5.0)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.2) AppleWebKit/535.46.6 (KHTML, like Gecko) Version/4.1 Safari/535.46.6",
		"Opera/9.80.(Windows 98; Win 9x 4.90; lg-UG) Presto/2.9.167 Version/12.00",
		"Opera/9.67.(Windows NT 5.0; nan-TW) Presto/2.9.172 Version/12.00",
		"Opera/8.35.(Windows NT 6.0; zu-ZA) Presto/2.9.172 Version/11.00",
		"Mozilla/5.0 (Windows 98) AppleWebKit/534.2 (KHTML, like Gecko) Chrome/20.0.839.0 Safari/534.2",
		"Mozilla/5.0 (compatible; MSIE 5.0; Windows 95; Trident/5.0)",
		"Mozilla/5.0 (Windows; U; Windows 98; Win 9x 4.90) AppleWebKit/531.38.2 (KHTML, like Gecko) Version/5.0 Safari/531.38.2",
		"Mozilla/5.0 (compatible; MSIE 7.0; Windows NT 4.0; Trident/5.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows CE; Trident/3.1)",
		"Mozilla/5.0 (compatible; MSIE 5.0; Windows NT 5.1; Trident/5.1)",
		"Opera/8.97.(Windows NT 5.1; ln-CD) Presto/2.9.177 Version/12.00",
		"Opera/8.60.(Windows NT 6.2; nn-NO) Presto/2.9.177 Version/11.00",
		"Opera/9.83.(Windows CE; fr-CA) Presto/2.9.174 Version/11.00",
		"Mozilla/5.0 (compatible; MSIE 7.0; Windows NT 6.2; Trident/4.0)",
		"Mozilla/5.0 (compatible; MSIE 7.0; Windows NT 5.01; Trident/5.0)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/535.35.2 (KHTML, like Gecko) Version/5.0.4 Safari/535.35.2",
		"Mozilla/5.0 (compatible; MSIE 5.0; Windows NT 6.2; Trident/3.0)",
		"Mozilla/5.0 (compatible; MSIE 7.0; Windows NT 4.0; Trident/3.1)",
		"Mozilla/5.0 (Windows NT 5.01; yi-US; rv:1.9.0.20) Gecko/2010-12-06 01:56:51 Firefox/3.8",
		"Mozilla/5.0 (compatible; MSIE 7.0; Windows 95; Trident/4.1)",
		"Mozilla/5.0 (compatible; MSIE 5.0; Windows 98; Trident/4.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.2; Trident/3.1)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.0) AppleWebKit/532.18.7 (KHTML, like Gecko) Version/5.1 Safari/532.18.7",
		"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.2; Trident/4.1)",
		"Mozilla/5.0 (Windows; U; Windows 98; Win 9x 4.90) AppleWebKit/532.1.2 (KHTML, like Gecko) Version/4.0.5 Safari/532.1.2",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.0 (KHTML, like Gecko) Chrome/46.0.830.0 Safari/536.0",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_5_7 rv:5.0; byn-ER) AppleWebKit/531.41.2 (KHTML, like Gecko) Version/5.0.1 Safari/531.41.2",
		"Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10_8_5 rv:3.0; mhr-RU) AppleWebKit/534.35.5 (KHTML, like Gecko) Version/4.0.4 Safari/534.35.5",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_10_8; rv:1.9.4.20) Gecko/2021-06-02 15:37:33 Firefox/5.0",
		"Mozilla/5.0 (iPod; U; CPU iPhone OS 3_1 like Mac OS X; hy-AM) AppleWebKit/533.49.3 (KHTML, like Gecko) Version/3.0.5 Mobile/8B115 Safari/6533.49.3",
		"Mozilla/5.0 (iPod; U; CPU iPhone OS 4_2 like Mac OS X; ln-CD) AppleWebKit/531.7.6 (KHTML, like Gecko) Version/3.0.5 Mobile/8B113 Safari/6531.7.6",
		"Mozilla/5.0 (iPod; U; CPU iPhone OS 3_2 like Mac OS X; ka-GE) AppleWebKit/531.38.5 (KHTML, like Gecko) Version/4.0.5 Mobile/8B112 Safari/6531.38.5",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2 rv:5.0; tl-PH) AppleWebKit/535.45.1 (KHTML, like Gecko) Version/4.0 Safari/535.45.1",
		"Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10_11_8 rv:6.0; wa-BE) AppleWebKit/535.37.7 (KHTML, like Gecko) Version/4.0 Safari/535.37.7",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_9_0 rv:6.0; ce-RU) AppleWebKit/533.14.2 (KHTML, like Gecko) Version/4.0.1 Safari/533.14.2",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_9_8; rv:1.9.4.20) Gecko/2020-10-12 09:24:52 Firefox/3.6.13",
		"Mozilla/5.0 (Macintosh; PPC Mac OS X 10_6_4 rv:5.0; hne-IN) AppleWebKit/532.13.7 (KHTML, like Gecko) Version/4.0.4 Safari/532.13.7",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5396.2 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5410.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5410.0 Safari/537.36",
	}
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
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
		// Note: Dial specifies the dial function for creating unencrypted TCP connections.
		// When httpClient sets this transport, it will use the tcp/udp connection returned from
		// function Dial instead of default tcp/udp connection. It's the key to set custom proxy for http transport
		DialContext: func(ctx context.Context, network, url string) (net.Conn, error) {
			return conn, nil
		},
		// from http.DefaultTransport
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

	transport := &http.Transport{
		// Note: Dial specifies the dial function for creating unencrypted TCP connections.
		// When httpClient sets this transport, it will use the tcp/udp connection returned from
		// function Dial instead of default tcp/udp connection. It's the key to set custom proxy for http transport
		Dial: func(string, string) (net.Conn, error) {
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
