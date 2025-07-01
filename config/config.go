package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig `toml:"database"`
}

/*
[server]
host = "0.0.0.0"
port = 8080
debug = false
caddyDir = "./"
*/
type ServerConfig struct {
	Port     int    `toml:"port"`
	Host     string `toml:"host"`
	Debug    bool   `toml:"debug"`
	CaddyDir string `toml:"caddyDir"`
}

/*
[database]
filepath = "sqlite.db"
*/
type DatabaseConfig struct {
	Filepath string `toml:"filepath"`
}

// LoadConfig 从 TOML 配置文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	if !FileExists(filePath) {
		// 楔入配置文件
		err := DefaultConfig().WriteConfig(filePath)
		if err != nil {
			return nil, err
		}
		return DefaultConfig(), nil
	}

	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// 写入配置文件
func (c *Config) WriteConfig(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(c)
}

// 检测文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 默认配置结构体
func DefaultConfig() *Config {
	/*
		[server]
		host = "0.0.0.0"
		port = 81
		debug = false
		caddyDir = "./"

		[tmpl]
		path = "./tmpl"

		[database]
		filepath = "caddydash.db"
	*/
	return &Config{
		Server: ServerConfig{
			Host:     "0.0.0.0",
			Port:     81,
			Debug:    false,
			CaddyDir: "./",
		},
		DB: DatabaseConfig{
			Filepath: "./db/caddydash.db",
		},
	}
}
