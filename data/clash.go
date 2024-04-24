package data

type ClashYaml struct {
	Port               int    `yaml:"port"`
	SocksPort          int    `yaml:"socks-port"`
	RedirPort          int    `yaml:"redir-port"`
	TproxyPort         int    `yaml:"tproxy-port"`
	AllowLan           bool   `yaml:"allow-lan"`
	BindAddress        string `yaml:"bind-address"`
	Mode               string `yaml:"mode"`
	LogLevel           string `yaml:"log-level"`
	ExternalController string `yaml:"external-controller"`
	Sniffer            struct {
		Enable   bool     `yaml:"enable"`
		Sniffing []string `yaml:"sniffing"`
	} `yaml:"sniffer"`
	DNS struct {
		Enable            bool     `yaml:"enable"`
		Ipv6              bool     `yaml:"ipv6"`
		DefaultNameserver []string `yaml:"default-nameserver"`
		FakeIPRange       string   `yaml:"fake-ip-range"`
		FakeIPFilter      []string `yaml:"fake-ip-filter"`
		Nameserver        []string `yaml:"nameserver"`
		Fallback          []string `yaml:"fallback"`
		FallbackFilter    struct {
			Geoip  bool     `yaml:"geoip"`
			Ipcidr []string `yaml:"ipcidr"`
		} `yaml:"fallback-filter"`
	} `yaml:"dns"`
	Proxies     []map[string]interface{} `yaml:"proxies"`
	ProxyGroups []struct {
		Name    string   `yaml:"name"`
		Type    string   `yaml:"type"`
		Proxies []string `yaml:"proxies"`
	} `yaml:"proxy-groups"`
	RuleProviders map[string]struct {
		Type     string `yaml:"type"`
		Behavior string `yaml:"behavior"`
		Path     string `yaml:"path"`
		URL      string `yaml:"url"`
		Interval int    `yaml:"interval"`
	} `yaml:"rule-providers"`
	Rules []string `yaml:"rules"`
}
