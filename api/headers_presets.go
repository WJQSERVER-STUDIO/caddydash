package api

// HeaderSet 定义了一个可复用的、完整的HTTP头预设。
// JSON标签用于API响应的序列化。
type HeaderSet struct {
	ID          string              `json:"id"`          // 唯一ID, e.g., "real_ip_cloudflare"
	Name        string              `json:"name"`        // UI上显示的名称 (为简化, 此处不使用i18n key)
	Description string              `json:"description"` // UI上的提示文本
	Target      string              `json:"target"`      // 目标: "global" 或 "upstream"
	Headers     map[string][]string `json:"headers"`     // 预设的请求头键值对
}

// HeaderSetMetadata 是HeaderSet的轻量级版本, 用于在列表中显示, 不包含具体的Headers数据。
type HeaderSetMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      string `json:"target"`
}

// registry 是一个私有变量, 作为所有Header预设的“数据库”。
// 在这里添加、修改或删除预设。
var registry = []HeaderSet{
	{
		ID:          "real_ip_cloudflare",
		Name:        "Cloudflare 真实IP",
		Description: "添加从Cloudflare获取真实客户端IP所需的请求头。适用于通过Cloudflare代理的流量。",
		Target:      "upstream", // 此预设仅适用于上游请求头
		Headers: map[string][]string{
			"X-Forwarded-For":   {"{http.request.header.CF-Connecting-IP}"},
			"X-Forwarded-Proto": {"{http.request.header.CF-Visitor}"},
			"X-Real-IP":         {"{http.request.header.CF-Connecting-IP}"},
		},
	},
	// 真实IP-直接
	{
		ID:          "real_ip_direct",
		Name:        "真实IP (直接)",
		Description: "当Caddy直接暴露在互联网上时, 使用此预设获取客户端真实IP。",
		Target:      "upstream",
		Headers: map[string][]string{
			"X-Forwarded-For":   {"{http.request.remote.host}"},
			"X-Forwarded-Proto": {"{http.request.scheme}"},
			"X-Real-IP":         {"{http.request.remote.host}"},
		},
	},
	// 真实IP-中间层
	{
		ID:          "real_ip_intermediate",
		Name:        "真实IP (中间层)",
		Description: "当Caddy位于反向代理或负载均衡器之后时, 使用此预设获取客户端真实IP。",
		Target:      "upstream",
		Headers: map[string][]string{
			"X-Forwarded-For":   {"{http.request.header.X-Forwarded-For}"},
			"X-Forwarded-Proto": {"{http.request.header.X-Forwarded-Proto}"},
			"X-Real-IP":         {"{http.request.header.X-Real-IP}"},
		},
	},
	{
		ID:          "common_security_headers",
		Name:        "通用安全响应头",
		Description: "添加一系列推荐的HTTP安全头以增强站点安全性 (HSTS, X-Frame-Options等)。",
		Target:      "global", // 此预设适用于全局响应头
		Headers: map[string][]string{
			"Strict-Transport-Security": {"max-age=31536000; includeSubDomains; preload"},
			"X-Frame-Options":           {"SAMEORIGIN"},
			"X-Content-Type-Options":    {"nosniff"},
			"Referrer-Policy":           {"strict-origin-when-cross-origin"},
			"Permissions-Policy":        {"geolocation=(), microphone=()"},
		},
	},
	{
		ID:          "cors_allow_all",
		Name:        "CORS (允许所有来源)",
		Description: "添加允许所有来源跨域请求的响应头。警告: 生产环境请谨慎使用。",
		Target:      "global",
		Headers: map[string][]string{
			"Access-Control-Allow-Origin":      {"*"},
			"Access-Control-Allow-Methods":     {"GET, POST, PUT, DELETE, OPTIONS"},
			"Access-Control-Allow-Headers":     {"Content-Type, Authorization, X-Requested-With"},
			"Access-Control-Allow-Credentials": {"true"},
		},
	},
}

// GetHeaderSetMetadataList 返回所有可用预设的元数据列表。
// 这个函数是线程安全的, 因为它只读取全局只读变量 registry。
func GetHeaderSetMetadataList() []HeaderSetMetadata {
	metadata := make([]HeaderSetMetadata, len(registry))
	for i, set := range registry {
		metadata[i] = HeaderSetMetadata{
			ID:          set.ID,
			Name:        set.Name,
			Description: set.Description,
			Target:      set.Target,
		}
	}
	return metadata
}

// GetHeaderSetByID 通过其唯一ID查找并返回一个完整的预设。
// 返回找到的预设和一个布尔值, 表示是否找到。
func GetHeaderSetByID(id string) (*HeaderSet, bool) {
	for _, set := range registry {
		if set.ID == id {
			// 返回该结构体的副本指针, 避免外部修改全局变量
			foundSet := set
			return &foundSet, true
		}
	}
	return nil, false
}
