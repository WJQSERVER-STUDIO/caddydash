package gen

import (
	"bytes"
	"caddydash/db"
	"encoding/gob"
	"fmt"
	"text/template"
)

func RenderConfig(site string, cdb *db.ConfigDB) error {
	// 检索site config
	paramsEntry, err := cdb.GetParams(site)
	if err != nil {
		return err
	}

	var caddycfg CaddyUniConfig
	err = DecodeGobConfig(paramsEntry.ParamsGOB, &caddycfg)
	if err != nil {
		return err
	}

	// 读取模板
	//tmplEntry, err := cdb.GetTemplate("reverse_proxy")
	tmplEntry, err := cdb.GetTemplate(caddycfg.Mode)
	if err != nil {
		return err
	}

	rpTmpl := string(tmplEntry.Content)
	// 使用caddycfg渲染最终产物
	parsedTmpl, parseErr := template.New(tmplEntry.Filename).Parse(rpTmpl)
	if parseErr != nil {
		return fmt.Errorf("db: failed to parse template content for '%s': %w", tmplEntry.Filename, parseErr)
	}
	var renderedContentBuilder bytes.Buffer
	if err := parsedTmpl.Execute(&renderedContentBuilder, caddycfg); err != nil {
		return fmt.Errorf("db: failed to render template '%s': %w", tmplEntry.Filename, err)
	}

	// 保存渲染结果
	renderedEntry := db.RenderedConfigEntry{
		Filename:        caddycfg.DomainConfig.Domain, // 使用域名作为文件名
		RenderedContent: renderedContentBuilder.Bytes(),
	}

	// 保存渲染产物
	if err := cdb.SaveRenderedConfig(renderedEntry); err != nil {
		return fmt.Errorf("db: failed to save rendered config for '%s': %w", caddycfg.DomainConfig.Domain, err)
	}

	return nil
}

func RenderGlobalConfig(paramsGob []byte, tmplContent []byte) ([]byte, error) {
	// 渲染caddyfile
	var globalConfig CaddyGlobalConfig
	err := DecodeGobConfig(paramsGob, &globalConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode global config params: %w", err)
	}

	parsedTmpl, parseErr := template.New("caddyfile").Parse(string(tmplContent))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse global caddyfile template: %w", parseErr)
	}

	var renderedContentBuilder bytes.Buffer
	if err := parsedTmpl.Execute(&renderedContentBuilder, globalConfig); err != nil {
		return nil, fmt.Errorf("failed to render global caddyfile template: %w", err)
	}

	return renderedContentBuilder.Bytes(), nil
}

// 把caddycfg内容转为GOB
func EncodeGobConfig(caddycfg any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(caddycfg); err != nil {
		return nil, fmt.Errorf("db: failed to Encode struct to GOB: %w", err)
	}
	return buf.Bytes(), nil // 返回编码后的字节切片.
}

func DecodeGobConfig(data []byte, tmplStruct any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(tmplStruct); err != nil {
		return fmt.Errorf("db: failed to Decode GOB to struct: %w", err)
	}
	return nil
}

/*
func RenderConfig(site string, cdb *db.ConfigDB) error {
	// 检索site config
	paramsEntry, err := cdb.GetParams(site)
	if err != nil {
		return err
	}

	var caddycfg any
	var specificConfigType interface{}

	if paramsEntry.TemplateType == "reverse_proxy" {
		specificConfigType = &CaddyReverseProxyConfig{}
	} else if paramsEntry.TemplateType == "file_server" {
		specificConfigType = &CaddyFileServerConfig{}
	} else {
		return fmt.Errorf("unknown template type: %s", paramsEntry.TemplateType)
	}

	err = DecodeGobConfig(paramsEntry.ParamsGOB, specificConfigType)
	if err != nil {
		log.Printf("decode gob config error: %v", err)
		return err
	}

	caddycfg = specificConfigType

	// 读取模板
	//tmplEntry, err := cdb.GetTemplate("reverse_proxy")
	tmplEntry, err := cdb.GetTemplate(paramsEntry.TemplateType)
	if err != nil {
		log.Printf("get template error: %v", err)
		tmplList, err := cdb.RangeTemplates()
		if err != nil {
			return err
		}
		log.Printf("template list: %v, GetName %s", tmplList, paramsEntry.TemplateType)
		return err
	}

	rpTmpl := string(tmplEntry.Content)
	// 使用caddycfg渲染最终产物
	parsedTmpl, parseErr := template.New(tmplEntry.Filename).Parse(rpTmpl)
	if parseErr != nil {
		return fmt.Errorf("db: failed to parse template content for '%s': %w", tmplEntry.Filename, parseErr)
	}
	var renderedContentBuilder bytes.Buffer
	if err := parsedTmpl.Execute(&renderedContentBuilder, caddycfg); err != nil {
		return fmt.Errorf("db: failed to render template '%s': %w", tmplEntry.Filename, err)
	}

	// 类型断言获得domain
	var domain string
	switch cfg := caddycfg.(type) {
	case *CaddyReverseProxyConfig:
		domain = cfg.Domain
	case *CaddyFileServerConfig:
		domain = cfg.Domain
	default:
		return fmt.Errorf("unknown config type for domain extraction")
	}

	// 保存渲染结果
	renderedEntry := db.RenderedConfigEntry{
		Filename:        domain, // 使用域名作为文件名
		RenderedContent: renderedContentBuilder.Bytes(),
	}

	// 保存渲染产物
	if err := cdb.SaveRenderedConfig(renderedEntry); err != nil {
		return fmt.Errorf("db: failed to save rendered config for '%s': %w", domain, err)
	}

	return nil
}
*/
