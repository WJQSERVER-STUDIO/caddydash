package api

import (
	"caddydash/config"
	"caddydash/db"
	"caddydash/gen"
	"fmt"
	"os"

	"github.com/infinite-iroha/touka"
)

func GetConfig(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		filename := c.Param("filename")
		params, err := cdb.GetParams(filename)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		// 解码[]byte的gob数据
		var config gen.CaddyUniConfig
		err = gen.DecodeGobConfig(params.ParamsOrigin, &config)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, config)
	}
}

func PutConfig(cdb *db.ConfigDB, cfg *config.Config) touka.HandlerFunc {
	return func(c *touka.Context) {
		filename := c.Param("filename")
		var config gen.CaddyUniConfig
		err := c.ShouldBindJSON(&config)
		if err != nil {
			c.JSON(400, touka.H{"error": err.Error()})
			return
		}

		var paramsGOB []byte
		var paramsOrigin []byte

		// Mode标识符固定为uni, 模板已被统合为只有uni
		paramsGOB, err = gen.EncodeGobConfig(config)
		if err != nil {
			c.Warnf("encode gob config error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		// 把json变为gob []byte
		paramsOrigin, err = gen.EncodeGobConfig(config)
		if err != nil {
			c.Warnf("encode origin config error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		paramsEntry := db.ParamsEntry{
			Filename:     filename,
			TemplateType: config.Mode,
			ParamsGOB:    paramsGOB,
			ParamsOrigin: paramsOrigin,
		}

		err = WriteConfig(cdb, paramsEntry, cfg, filename)
		if err != nil {
			c.Warnf("write config error: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}

		c.JSON(200, touka.H{"message": "config saved and rendered"})
		return
	}
}

// 渲染并写入配置
func WriteConfig(cdb *db.ConfigDB, paramsEntry db.ParamsEntry, cfg *config.Config, filename string) error {
	var err error
	err = cdb.SaveParams(paramsEntry)
	if err != nil {
		err = fmt.Errorf("save params error: %w", err)
		return err
	}
	err = gen.RenderConfig(filename, cdb)
	if err != nil {
		err = fmt.Errorf("render config error: %w", err)
		return err
	}
	// 写入文件
	renderedEntry, err := cdb.GetRenderedConfig(filename)
	if err != nil {
		err = fmt.Errorf("get rendered config error: %w", err)
		return err
	}
	err = os.WriteFile(cfg.Server.CaddyDir+"config.d/"+filename, renderedEntry.RenderedContent, 0644)
	if err != nil {
		err = fmt.Errorf("write rendered config file error: %w", err)
		return err
	}
	return nil
}

func DeleteConfig(cdb *db.ConfigDB, cfg *config.Config) touka.HandlerFunc {
	return func(c *touka.Context) {
		filename := c.Param("filename")
		err := cdb.DeleteParams(filename)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		// 删除文件
		err = os.Remove(cfg.Server.CaddyDir + "config.d/" + filename)
		if err != nil {
			c.Warnf("delete rendered config file error: %v", err)
		}
		c.JSON(200, touka.H{"message": "config deleted"})
	}
}

func GetHeadersPreset() touka.HandlerFunc {
	return func(c *touka.Context) {
		presetName := c.Param("name")
		if presetName == "" {
			c.JSON(400, touka.H{"error": "presetName is required"})
			return
		}
		preset, found := GetHeaderSetByID(presetName)
		if !found {
			c.JSON(404, touka.H{"error": "preset not found"})
			return
		}
		c.JSON(200, preset)
	}
}
