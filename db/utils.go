package db

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

// --- GOB 编码/解码辅助函数 ---

// gobEncode 将 Go 值编码为 GOB 格式的字节切片.
// 'data' 必须是可被 GOB 编码的值; 例如基本类型, 切片, 映射, 结构体.
func gobEncode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, fmt.Errorf("db: failed to GOB encode data: %w", err)
	}
	return buf.Bytes(), nil
}

// gobDecode 将 GOB 格式的字节切片解码到 Go 值.
// 'data' 是 GOB 编码的字节切片.
// 'valuePtr' 必须是指向目标 Go 值的指针; 其类型必须与编码时的数据类型兼容.
func gobDecode(data []byte, valuePtr interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(valuePtr); err != nil {
		return fmt.Errorf("db: failed to GOB decode data: %w", err)
	}
	return nil
}

// --- 业务逻辑: 渲染并保存 (由使用者调用) ---
// RenderAndSaveConfig 从数据库获取模板和参数; 渲染后保存到 'rendered_configs' 表.
// 这是一个组合操作; 通常由应用程序逻辑在需要时调用.
// filename: 要渲染的配置文件的唯一标识.
// templateParser: 一个实现了 TemplateParser 接口的模板解析器实例.
// dynamicParams: 运行时动态提供的参数; 它们会覆盖存储在数据库中的同名参数.
func (cdb *ConfigDB) RenderAndSaveConfig(filename string, dynamicParams map[string]interface{}) error {
	// 1. 获取模板内容.
	tmplEntry, err := cdb.GetTemplate(filename)
	if err != nil {
		return fmt.Errorf("db: failed to get template '%s' for rendering: %w", filename, err)
	}

	// 2. 获取存储的参数.
	paramsEntry, err := cdb.GetParams(filename)
	var storedParams map[string]interface{}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			storedParams = make(map[string]interface{}) // 参数不存在; 使用空 map.
		} else {
			return fmt.Errorf("db: failed to get parameters for '%s': %w", filename, err)
		}
	} else {
		// 解码 GOB 参数到 map.
		if err := gobDecode(paramsEntry.ParamsGOB, &storedParams); err != nil {
			return fmt.Errorf("db: failed to decode stored parameters for '%s': %w", filename, err)
		}
	}

	// 合并参数: 动态传入的参数覆盖存储的参数.
	if dynamicParams != nil {
		for k, v := range dynamicParams {
			storedParams[k] = v
		}
	}

	// 3. 渲染模板.
	var parsedTmpl *template.Template
	var parseErr error

	// 使用传入的 templateParser 实例来解析模板内容.
	// 注意: templateParser.Parse 是在提供的实例上调用; 以解析特定内容.
	//parsedTmpl, parseErr = templateParser.Parse(string(tmplEntry.Content))
	parsedTmpl, parseErr = template.New(tmplEntry.Filename).Parse(string(tmplEntry.Content))

	if parseErr != nil {
		return fmt.Errorf("db: failed to parse template content for '%s': %w", tmplEntry.Filename, parseErr)
	}

	var renderedContentBuilder strings.Builder
	if err := parsedTmpl.Execute(&renderedContentBuilder, storedParams); err != nil {
		return fmt.Errorf("db: failed to render template '%s': %w", tmplEntry.Filename, err)
	}

	// 4. 保存渲染结果.
	renderedEntry := RenderedConfigEntry{
		Filename:        filename,
		RenderedContent: []byte(renderedContentBuilder.String()),
	}
	if err := cdb.SaveRenderedConfig(renderedEntry); err != nil {
		return fmt.Errorf("db: failed to save rendered config for '%s': %w", filename, err)
	}
	return nil
}
