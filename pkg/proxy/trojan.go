package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrorNotTrojanink = errors.New("not a correct trojan link")
)

// TODO unknown field
// Link: host, path
// Trojan: Network GrpcOpts

type Trojan struct {
	Base
	Password       string   `yaml:"password" json:"password"`
	ALPN           []string `yaml:"alpn,omitempty" json:"alpn,omitempty"`
	SNI            string   `yaml:"sni,omitempty" json:"sni,omitempty"`
	Transport      string   `yaml:"transport,omitempty" json:"transport,omitempty"`
	Host           string   `yaml:"host,omitempty" json:"host,omitempty"`
	Path           string   `yaml:"path,omitempty" json:"path,omitempty"`
	SkipCertVerify bool     `yaml:"skip-cert-verify,omitempty" json:"skip-cert-verify,omitempty"`
	UDP            bool     `yaml:"udp,omitempty" json:"udp,omitempty"`
	// Network        string      `yaml:"network,omitempty" json:"network,omitempty"`
	// GrpcOpts       GrpcOptions `yaml:"grpc-opts,omitempty" json:"grpc-opts,omitempty"`
}

/**
  - name: "trojan"
    type: trojan
    server: server
    port: 443
    password: yourpsk
    # udp: true
    # sni: example.com # aka server name
    # alpn:
    #   - h2
    #   - http/1.1
    # skip-cert-verify: true
*/

func (t Trojan) Identifier() string {
	return net.JoinHostPort(t.Server, strconv.Itoa(t.Port)) + t.Password
}

func (t Trojan) String() string {
	data, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(data)
}

func (t Trojan) ToClash() string {
	data, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return "- " + string(data)
}

func (t Trojan) ToLoon() string {
	verify := true
	if t.SkipCertVerify {
		verify = false
	}
	//trojan1 = trojan,example.com,443,"password",transport=tcp,skip-cert-verify=false,sni=example.com,udp=true
	text := fmt.Sprintf("%s = trojan,%s:%d,%s,transport=%s,skip-cert-verify=%t",
		t.Name, t.Server, t.Port, `"`+t.Password+`"`, t.Transport, verify)
	//trojan2 = trojan,example.com,443,"password",transport=ws,path=/,host=micsoft.com,skip-cert-verify=true,sni=example.com,udp=true
	if t.Transport == "ws" {
		text += fmt.Sprintf("path=%s,host=%s", t.Path, t.Host)
	}
	text += fmt.Sprintf(",sni=%s,udp=%t", t.SNI, t.UDP)
	return text
}

func (t Trojan) ToQuantumultX() string {
	vfi := true
	if t.SkipCertVerify {
		vfi = false
	}
	// trojan=example.com:443, password=pwd, over-tls=true, tls-verification=false, fast-open=false, udp-relay=false, tag=节点名称
	text := fmt.Sprintf("trojan = %s:%d, password=%s, over-tls=true, tls-host=%s, fast-open=true, udp-relay=%t, tls-verification=%t",
		t.Server, t.Port, t.Password, t.SNI, t.UDP, vfi)
	text = text + fmt.Sprintf(", tag=%s", t.Name)
	return text
}

func (t Trojan) ToSurge() string {
	// node1 = trojan, server, port,  password=, sni=, obfs-host=, skip-cert-verify=false, udp-relay=false
	return fmt.Sprintf("%s = trojan, %s, %d, password=%s, sni=%s, skip-cert-verify=false, tfo=false, udp-relay=%t",
		t.Name, t.Server, t.Port, t.Password, t.SNI, t.UDP)
}

func (t Trojan) Clone() Proxy {
	return &t
}

// https://p4gefau1t.github.io/trojan-go/developer/url/
func (t Trojan) Link() (link string) {
	query := url.Values{}
	if t.SNI != "" {
		query.Set("sni", url.QueryEscape(t.SNI))
	}

	uri := url.URL{
		Scheme:   "trojan",
		User:     url.User(url.QueryEscape(t.Password)),
		Host:     net.JoinHostPort(t.Server, strconv.Itoa(t.Port)),
		RawQuery: query.Encode(),
		Fragment: t.Name,
	}

	return uri.String()
}

func ParseTrojanLink(link string) (*Trojan, error) {
	if !strings.HasPrefix(link, "trojan://") && !strings.HasPrefix(link, "trojan-go://") {
		return nil, ErrorNotTrojanink
	}

	/**
	trojan-go://
	    $(trojan-password)
	    @
	    trojan-host
	    :
	    port
	/?
	    sni=$(tls-sni.com)&
	    type=$(original|ws|h2|h2+ws)&
	        host=$(websocket-host.com)&
	        path=$(/websocket/path)&
	    encryption=$(ss;aes-256-gcm;ss-password)&
	    plugin=$(...)
	#$(descriptive-text)
	*/

	uri, err := url.Parse(link)
	if err != nil {
		return nil, ErrorNotSSLink
	}

	password := uri.User.Username()
	password, _ = url.QueryUnescape(password)

	server := uri.Hostname()
	port, _ := strconv.Atoi(uri.Port())

	moreInfos := uri.Query()
	sni := moreInfos.Get("sni")
	sni, _ = url.QueryUnescape(sni)
	transformType := moreInfos.Get("type")
	transformType, _ = url.QueryUnescape(transformType)
	// host := moreInfos.Get("host")
	// host, _ = url.QueryUnescape(host)
	// path := moreInfos.Get("path")
	// path, _ = url.QueryUnescape(path)

	alpn := make([]string, 0)
	if transformType == "h2" {
		alpn = append(alpn, "h2")
	}

	if port == 0 {
		return nil, ErrorNotTrojanink
	}

	return &Trojan{
		Base: Base{
			Name:   "",
			Server: server,
			Port:   port,
			Type:   "trojan",
		},
		Password:       password,
		ALPN:           alpn,
		SNI:            sni,
		UDP:            true,
		SkipCertVerify: true,
	}, nil
}

var (
	trojanPlainRe = regexp.MustCompile("trojan(-go)?://([A-Za-z0-9+/_&?=@:%.-])+")
)

func GrepTrojanLinkFromString(text string) []string {
	results := make([]string, 0)
	if !strings.Contains(text, "trojan://") {
		return results
	}
	texts := strings.Split(text, "trojan://")
	for _, text := range texts {
		results = append(results, trojanPlainRe.FindAllString("trojan://"+text, -1)...)
	}
	return results
}
