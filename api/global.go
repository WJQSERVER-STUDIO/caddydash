package api

import (
	"caddydash/config"
	"caddydash/db"
	"caddydash/gen"
	"os"

	"github.com/infinite-iroha/touka"
)

func PutGlobalConfig(cdb *db.ConfigDB, cfg *config.Config) touka.HandlerFunc {
	return func(c *touka.Context) {
		var config gen.CaddyGlobalConfig
		err := c.ShouldBindJSON(&config)
		if err != nil {
			c.JSON(400, touka.H{"error": err.Error()})
			return
		}

		var paramsGob []byte

		paramsGob, err = gen.EncodeGobConfig(config)
		if err != nil {
			c.Warnf("encode gob config error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		// 取出数据库内的tmpl
		tmplContent, err := cdb.GetGlobalTemplate("caddyfile")
		if err != nil {
			c.Warnf("get global template error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		renderedContent, err := gen.RenderGlobalConfig(paramsGob, tmplContent)
		if err != nil {
			c.Warnf("render global config error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		//回写条目到数据库
		err = cdb.SaveGlobalConfig(db.GlobalConfig{
			Filename:        "caddyfile",
			Params:          paramsGob,
			TmplContent:     tmplContent,
			RenderedContent: renderedContent,
		})
		if err != nil {
			c.Warnf("save global config error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		err = os.WriteFile(cfg.Server.CaddyDir+"Caddyfile", renderedContent, 0644)
		if err != nil {
			c.Warnf("write Caddyfile error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		c.JSON(200, touka.H{"message": "global config saved"})
		return
	}
}

// GetGlobalConfig 检出已有配置
func GetGlobalConfig(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		globalConfig, err := cdb.GetGlobalConfig("caddyfile")
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		var config gen.CaddyGlobalConfig
		err = gen.DecodeGobConfig(globalConfig.Params, &config)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, config)
	}
}
