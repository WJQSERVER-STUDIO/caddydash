package gen

import (
	"caddydash/config"
	"caddydash/db"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// 把指定目录下的文件作为模板读入
func ReadTmplToDB(dir string, cdb *db.ConfigDB) error {
	// 遍历目录下的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 读取文件内容
		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		tmplEntry := db.TemplateEntry{
			Filename:     file.Name(),
			TemplateType: file.Name(),
			Content:      content,
			UpdatedAt:    time.Now().Unix(),
		}

		// 存储到数据库
		err = cdb.SaveTemplate(tmplEntry)
		if err != nil {
			return err
		}

		// 输出tmpl名
		fmt.Printf("Read template: %s\n", file.Name())

	}
	return nil
}

func Add80SiteConfig(cfg *config.Config, cdb *db.ConfigDB) error {
	// 检查:80是否已存在于数据库中

	_, err := cdb.GetParams(":80")
	if err == nil {
		// 如果存在，则不添加
		return nil
	}

	siteConfig := CaddyUniConfig{
		DomainConfig: CaddyUniDomainConfig{
			Domain:      ":80",
			MutiDomains: false,
			Domains:     []string{":80"},
		},
		Mode: "uni",
		FileServer: CaddyUniFileServerConfig{
			EnableFileServer: true,
			FileDirPath:      cfg.Server.CaddyDir + "pages/demo",
			EnableBrowser:    false,
		},
		Log: CaddyUniLogConfig{
			EnableLog: true,
			LogDomain: ":80",
		},
		ErrorPage: CaddyUniErrorPageConfig{
			EnableErrorPage: true,
		},
		Encode: CaddyUniEncodeConfig{
			EnableEncode: true,
		},
	}

	// 制作db.ParamsEntry
	gobData, err := EncodeGobConfig(siteConfig)
	if err != nil {
		return err
	}

	originGobData, err := EncodeGobConfig(siteConfig)
	if err != nil {
		return err
	}

	paramsEntry := db.ParamsEntry{
		Filename:     ":80",
		TemplateType: "file_server",
		ParamsGOB:    gobData,
		ParamsOrigin: originGobData,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	filename := paramsEntry.Filename

	// 保存到数据库
	err = cdb.SaveParams(paramsEntry)
	if err != nil {
		return err
	}

	// 渲染配置
	err = RenderConfig(filename, cdb)
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

// 读入全局配置模板
func ReadGlobalTmpl(dir string) ([]byte, error) {
	// 读取目录下的caddyfile文件

	content, err := os.ReadFile(filepath.Join(dir, "gtmpl", "caddyfile"))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Read global template: %s\n", "caddyfile")
	return content, nil
}

func SetGlobalConfig(cfg *config.Config, cdb *db.ConfigDB) error {
	var config = DefaultGlobalConfig
	paramsGob, err := EncodeGobConfig(config)
	if err != nil {
		return fmt.Errorf("encode gob config error: %w", err)
	}

	// 取出数据库内的tmpl
	tmplContent, err := ReadGlobalTmpl(cfg.Server.CaddyDir)
	if err != nil {
		return fmt.Errorf("get global template error: %w", err)
	}

	renderedContent, err := RenderGlobalConfig(paramsGob, tmplContent)
	if err != nil {
		return fmt.Errorf("render global config error: %w", err)
	}

	//回写条目到数据库
	err = cdb.SaveGlobalConfig(db.GlobalConfig{
		Filename:        "caddyfile",
		Params:          paramsGob,
		TmplContent:     tmplContent,
		RenderedContent: renderedContent,
	})
	if err != nil {
		return fmt.Errorf("save global config error: %w", err)
	}

	err = os.WriteFile(cfg.Server.CaddyDir+"Caddyfile", renderedContent, 0644)
	if err != nil {
		return fmt.Errorf("write Caddyfile error: %w", err)
	}
	return nil
}
