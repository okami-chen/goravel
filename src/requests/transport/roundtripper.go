package transport

import (
	"context"
	"errors"
	"fmt"
	"net"

	"strings"
	"sync"

	utls "github.com/refraction-networking/utls"
	http "github.com/wangluozhe/chttp"
	http2 "github.com/wangluozhe/chttp/http2"
	"golang.org/x/net/proxy"
)

var errProtocolNegotiated = errors.New("protocol negotiated")

type roundTripper struct {
	sync.RWMutex
	// fix typing
	JA3       string
	UserAgent string

	cachedConnections sync.Map
	cachedTransports  sync.Map

	dialer        proxy.ContextDialer
	config        *utls.Config
	tlsExtensions *TLSExtensions
	http2Settings *http2.HTTP2Settings
	forceHTTP1    bool
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("user-agent") == "" {
		req.Header.Set("user-agent", rt.UserAgent)
	}
	addr := rt.getDialTLSAddr(req)
	transport, err := rt.getTransport(req, addr)
	if err != nil {
		return nil, err
	}
	return transport.RoundTrip(req)
}

func (rt *roundTripper) getTransport(req *http.Request, addr string) (http.RoundTripper, error) {
	transport, ok := rt.cachedTransports.Load(addr)
	if ok {
		return transport.(http.RoundTripper), nil
	}

	switch strings.ToLower(req.URL.Scheme) {
	case "http":
		ts := &http.Transport{DialContext: rt.dialer.DialContext, DisableKeepAlives: true}
		rt.cachedTransports.Store(addr, ts)
		return ts, nil
	case "https":
	default:
		return nil, fmt.Errorf("invalid URL scheme: [%v]", req.URL.Scheme)
	}

	_, err := rt.dialTLS(context.Background(), "tcp", addr)
	switch err {
	case errProtocolNegotiated:
	case nil:
		// Should never happen.
		//panic("dialTLS returned no error when determining cachedTransports")
		return nil, fmt.Errorf("dialTLS returned no error when determining cachedTransports")
	default:
		return nil, err
	}
	transport, _ = rt.cachedTransports.Load(addr)

	return transport.(http.RoundTripper), nil
}

func (rt *roundTripper) dialTLS(ctx context.Context, network, addr string) (net.Conn, error) {
	rt.Lock()
	defer rt.Unlock()

	// If we have the connection from when we determined the HTTPS
	// cachedTransports to use, return that.
	conn, ok := rt.cachedConnections.Load(addr)
	if ok {
		return conn.(net.Conn), nil
	}
	rawConn, err := rt.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	var host string
	if host, _, err = net.SplitHostPort(addr); err != nil {
		host = addr
	}
	//////////////////

	spec, err := StringToSpec(rt.JA3, rt.UserAgent, rt.tlsExtensions, rt.forceHTTP1)
	if err != nil {
		return nil, err
	}

	rt.config.ServerName = host
	tlsConn := utls.UClient(rawConn, rt.config.Clone(), utls.HelloCustom)

	if err := tlsConn.ApplyPreset(spec); err != nil {
		return nil, err
	}

	if err = tlsConn.HandshakeContext(ctx); err != nil {
		_ = tlsConn.Close()

		if err.Error() == "tls: CurvePreferences includes unsupported curve" {
			//fix this
			return nil, fmt.Errorf("conn.Handshake() error for tls 1.3 (please retry request): %+v", err)
		}
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	//////////
	_, ok = rt.cachedTransports.Load(addr)
	if ok {
		return tlsConn, nil
	}

	// No http.Transport constructed yet, create one based on the results
	// of ALPN.
	switch tlsConn.ConnectionState().NegotiatedProtocol {
	case http2.NextProtoTLS:
		t2 := http2.Transport{
			DialTLS:         rt.dialTLSHTTP2,
			TLSClientConfig: rt.config,
		}
		if rt.http2Settings != nil {
			t2.HTTP2Settings = rt.http2Settings
		}
		rt.cachedTransports.Store(addr, &t2)
	default:
		// Assume the remote peer is speaking HTTP 1.x + TLS.
		rt.cachedTransports.Store(addr, &http.Transport{DialTLSContext: rt.dialTLS})

	}

	// Stash the connection just established for use servicing the
	// actual request (should be near-immediate).
	rt.cachedConnections.Store(addr, tlsConn)

	return nil, errProtocolNegotiated
}

func (rt *roundTripper) dialTLSHTTP2(network, addr string, _ *utls.Config) (net.Conn, error) {
	return rt.dialTLS(context.Background(), network, addr)
}

func (rt *roundTripper) getDialTLSAddr(req *http.Request) string {
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err == nil {
		return net.JoinHostPort(host, port)
	}
	return net.JoinHostPort(req.URL.Host, "443") // we can assume port is 443 at this point
}

func newRoundTripper(browser Browser, config *utls.Config, tlsExtensions *TLSExtensions, http2Settings *http2.HTTP2Settings, forceHTTP1 bool, dialer ...proxy.ContextDialer) http.RoundTripper {
	if config == nil {
		if strings.Index(strings.Split(browser.JA3, ",")[2], "-41") == -1 {
			config = &utls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			config = &utls.Config{
				InsecureSkipVerify: true,
				SessionTicketKey:   [32]byte{},
				ClientSessionCache: utls.NewLRUClientSessionCache(0),
				OmitEmptyPsk:       true,
			}
		}
	}
	if len(dialer) > 0 {

		return &roundTripper{
			dialer: dialer[0],

			JA3:               browser.JA3,
			UserAgent:         browser.UserAgent,
			cachedTransports:  sync.Map{},
			cachedConnections: sync.Map{},
			config:            config,
			tlsExtensions:     tlsExtensions,
			http2Settings:     http2Settings,
			forceHTTP1:        forceHTTP1,
		}
	}

	return &roundTripper{
		dialer: proxy.Direct,

		JA3:               browser.JA3,
		UserAgent:         browser.UserAgent,
		cachedTransports:  sync.Map{},
		cachedConnections: sync.Map{},
		config:            config,
		tlsExtensions:     tlsExtensions,
		http2Settings:     http2Settings,
		forceHTTP1:        forceHTTP1,
	}
}
