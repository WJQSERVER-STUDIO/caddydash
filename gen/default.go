package gen

var (
	DefaultGlobalConfig = CaddyGlobalConfig{
		Debug: false,
		PortsConfig: CaddyGlobalPortsConfig{
			AdminPort: "localhost:2019",
			HTTPPort:  80,
			HTTPSPort: 443,
		},
		Metrics: true,
		LogConfig: CaddyGlobalLogConfig{
			Level:             "INFO",
			RotateSize:        "10MB",
			RotateKeep:        "10",
			RotateKeepForTime: "24h",
		},
	}
)
