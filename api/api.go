package api

import (
	"caddydash/apic"
	"caddydash/config"
	"caddydash/db"
	"caddydash/gen"

	"github.com/infinite-iroha/touka"
)

func ApiGroup(v0 touka.IRouter, cdb *db.ConfigDB, cfg *config.Config, version string) {
	api := v0.Group("/api")
	api.GET("/config/filenames", func(c *touka.Context) {
		filenames, err := cdb.GetFileNames()
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, filenames)
	})

	api.GET("/info", infoHandle(version))

	// 配置参数相关
	cfgr := api.Group("/config")
	{
		cfgr.GET("/file/:filename", GetConfig(cdb))            // 读取配置(与写入一致)
		cfgr.PUT("/file/:filename", PutConfig(cdb, cfg))       // 写入配置
		cfgr.DELETE("/file/:filename", DeleteConfig(cdb, cfg)) //删除配置

		cfgr.GET("/files/params", FilesParams(cdb))       // 获取所有配置, 需进行decode
		cfgr.GET("/files/templates", FilesTemplates(cdb)) // 获取所有模板
		cfgr.GET("/files/rendered", FilesRendered(cdb))   // 获取所有渲染产物

		cfgr.GET("/templates", GetTemplates(cdb)) // 获取可用模板名称

		cfgr.GET("/headers-presets", func(c *touka.Context) {
			c.JSON(200, GetHeaderSetMetadataList())
		})
		cfgr.GET("/headers-presets/:name", GetHeadersPreset())

		glbr := api.Group("/global")
		{
			glbr.GET("/log/levels", func(c *touka.Context) {
				c.JSON(200, gen.LogLevelList)
			})
			glbr.GET("/tls/providers", func(c *touka.Context) {
				c.JSON(200, gen.ProviderList)
			})
			glbr.PUT("/config", PutGlobalConfig(cdb, cfg))
			glbr.GET("/config", GetGlobalConfig(cdb))
		}
	}

	// caddy实例相关
	caddy := api.Group("/caddy")
	{
		caddy.POST("/stop", apic.StopCaddy()) // 无需payload
		caddy.POST("/run", apic.StartCaddy(cfg))
		caddy.POST("/restart", apic.RestartCaddy(cfg))
		caddy.GET("/status", apic.IsCaddyRunning())

		logs := caddy.Group("/logs")
		{
			logs.GET("/stdout", StreamLog(cfg.Server.CaddyDir+"caddystdout.log"))
			logs.GET("/caddy", StreamLog(cfg.Server.CaddyDir+"log/caddy.log"))
			logs.GET("/sites", ListSiteLogs(cfg.Server.CaddyDir+"log"))
			logs.GET("/site/:sitename", SiteLog(cfg.Server.CaddyDir+"log"))
		}
	}

	// 鉴权相关
	auth := api.Group("/auth")
	{
		auth.POST("/login", func(c *touka.Context) {
			AuthLogin(c, cfg, cdb)
		})
		auth.POST("/logout", func(c *touka.Context) {
			AuthLogout(c)
		})
		auth.GET("/logout", func(c *touka.Context) {
			AuthLogout(c)
		})
		auth.GET("/init", AuthInitStatus())
		auth.POST("/init", AuthInitHandle(cdb))
		auth.POST("resetpwd", ResetPassword(cdb))
	}

	// 运行时相关
	rtr := api.Group("/runtime")
	{
		rtr.GET("/", runtimeInfo())
	}
}

// GetTemplates 获取可用的tmpls name
func GetTemplates(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		templates, err := cdb.RangeTemplates()
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, templates)
	}
}
