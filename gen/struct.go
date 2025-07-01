package gen

func HeadersMapToHeadersUp(headers map[string][]string) []string {
	var headersUp []string
	for key, values := range headers {
		for _, value := range values {
			headersUp = append(headersUp, key+" "+value)
		}
	}
	return headersUp
}

type CaddyUniConfig struct {
	DomainConfig CaddyUniDomainConfig     `json:"domain_config"`
	Mode         string                   `json:"mode"`
	Upstream     CaddyUniUpstreamConfig   `json:"upstream_config"`
	FileServer   CaddyUniFileServerConfig `json:"file_server_config"`
	Headers      map[string][]string      `json:"headers"`
	Log          CaddyUniLogConfig        `json:"log_config"`
	ErrorPage    CaddyUniErrorPageConfig  `json:"error_page_config"`
	Encode       CaddyUniEncodeConfig     `json:"encode_config"`
}

type CaddyUniDomainConfig struct {
	Domain      string   `json:"domain"`
	MutiDomains bool     `json:"muti_domains"`
	Domains     []string `json:"domains"`
}

type CaddyUniUpstreamConfig struct {
	EnableUpStream  bool                `json:"enable_upstream"`
	UpStream        string              `json:"upstream"`
	MutiUpStreams   bool                `json:"muti_upstream"`
	UpStreams       []string            `json:"upstream_servers"`
	UpStreamHeaders map[string][]string `json:"upstream_headers"`
}

type CaddyUniFileServerConfig struct {
	EnableFileServer bool   `json:"enable_file_server"`
	FileDirPath      string `json:"file_dir_path"`
	EnableBrowser    bool   `json:"enable_browser"`
}

type CaddyUniLogConfig struct {
	EnableLog bool   `json:"enable_log"`
	LogDomain string `json:"log_domain"`
}

type CaddyUniErrorPageConfig struct {
	EnableErrorPage bool `json:"enable_error_page"`
}

type CaddyUniEncodeConfig struct {
	EnableEncode bool `json:"enable_encode"`
}

type CaddyGlobalConfig struct {
	Debug            bool                        `json:"debug"`
	PortsConfig      CaddyGlobalPortsConfig      `json:"ports_config"`
	Metrics          bool                        `json:"metrics"`
	LogConfig        CaddyGlobalLogConfig        `json:"log_config"`
	TLSConfig        CaddyGlobalTLSConfig        `json:"tls_config"`
	TLSSnippetConfig CaddyGlobalSnippetTLSConfig `json:"tls_snippet_config"`
}

type CaddyGlobalPortsConfig struct {
	AdminPort string `json:"admin_port"`
	HTTPPort  uint16 `json:"http_port"`
	HTTPSPort uint16 `json:"https_port"`
}

type CaddyGlobalLogConfig struct {
	Level string `json:"level"`
	// 日志滚动配置
	RotateSize        string `json:"rotate_size"`
	RotateKeep        string `json:"rotate_keep"`
	RotateKeepForTime string `json:"rotate_keep_for_time"`
}

// 维护一个日志等级列表
// Possible levels: DEBUG, INFO, WARN, ERROR, PANIC, and FATAL
var LogLevelList = map[string]struct{}{
	"DEBUG": {},
	"INFO":  {},
	"WARN":  {},
	"ERROR": {},
	"PANIC": {},
	"FATAL": {},
}

type CaddyGlobalTLSConfig struct {
	EnableDNSChallenge bool   `json:"enable_dns_challenge"`
	Provider           string `json:"provider"`
	Token              string `json:"token"`
	ECHOuterSNI        string `json:"echouter_sni"`
	Email              string `json:"email"`
}

type CaddyGlobalSnippetTLSConfig struct {
	EnableSiteTLSSnippet bool   `json:"enable_site_tls_snippet"`
	Email                string `json:"email"`
	Provider             string `json:"provider"`
	Token                string `json:"token"`
}

// 维护一个提供商列表
var ProviderList = map[string]struct{}{
	"cloudflare": {},
}
