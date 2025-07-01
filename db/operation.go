package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// 用户校验操作
/*
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);`)
*/

// AddUser 向 'users' 表中添加一个新用户.
func (cdb *ConfigDB) AddUser(username, password string) error {
	insertSQL := `
	INSERT INTO users (username, password)
	VALUES (?, ?);
	`
	_, err := cdb.DB.Exec(insertSQL, username, password)
	if err != nil {
		return fmt.Errorf("db: failed to add user '%s': %w", username, err)
	}
	return nil
}

// GetUserByUsername 从 'users' 表中根据用户名获取用户信息.
func (cdb *ConfigDB) GetUserByUsername(username string) (*UsersTable, error) {
	querySQL := `SELECT username, password, created_at, updated_at FROM users WHERE username = ?;`
	row := cdb.DB.QueryRow(querySQL, username)

	user := &UsersTable{}
	err := row.Scan(&user.UserName, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db: user '%s' not found: %w", username, err)
		}
		return nil, fmt.Errorf("db: failed to get user '%s': %w", username, err)
	}
	return user, nil
}

// DeleteUser 从 'users' 表中删除一个用户.
func (cdb *ConfigDB) DeleteUser(username string) error {
	_, err := cdb.DB.Exec(`DELETE FROM users WHERE username = ?;`, username)
	if err != nil {
		return fmt.Errorf("db: failed to delete user '%s': %w", username, err)
	}
	return nil
}

// UpdateUserPassword 更新用户的密码.
func (cdb *ConfigDB) UpdateUserPassword(username, newPassword string) error {
	updateSQL := `
	UPDATE users
	SET password = ?, updated_at = strftime('%s', 'now')
	WHERE username = ?;
	`
	result, err := cdb.DB.Exec(updateSQL, newPassword, username)
	if err != nil {
		return fmt.Errorf("db: failed to update user '%s' password: %w", username, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("db: failed to get rows affected by update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("db: no user with username '%s' found to update", username)
	}
	return nil
}

// RangeUserNames 获取所有用户的用户名.
func (cdb *ConfigDB) RangeUserNames() ([]string, error) {
	querySQL := `SELECT username FROM users;`
	rows, err := cdb.DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get usernames from users: %w", err)
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("db: failed to scan username: %w", err)
		}
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error during user rows iteration: %w", err)
	}

	return usernames, nil
}

// HasAnyUser 检查 'users' 表中是否存在任何用户.
func (cdb *ConfigDB) HasAnyUser() (bool, error) {
	querySQL := `SELECT EXISTS(SELECT 1 FROM users LIMIT 1);`
	var exists bool
	err := cdb.DB.QueryRow(querySQL).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("db: failed to check if any user exists: %w", err)
	}
	return exists, nil
}

// IsUserExists 检查指定用户名的用户是否存在.
func (cdb *ConfigDB) IsUserExists(username string) (bool, error) {
	querySQL := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? LIMIT 1);`
	var exists bool
	err := cdb.DB.QueryRow(querySQL, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("db: failed to check if user '%s' exists: %w", username, err)
	}
	return exists, nil
}

// GetPasswordByUsername 从 'users' 表中根据用户名获取密码.
func (cdb *ConfigDB) GetPasswordByUsername(username string) (string, error) {
	querySQL := `SELECT password FROM users WHERE username = ?;`
	var password string
	err := cdb.DB.QueryRow(querySQL, username).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("db: user '%s' not found: %w", username, err)
		}
		return "", fmt.Errorf("db: failed to get password for user '%s': %w", username, err)
	}
	return password, nil
}

// --- 模板操作 (Templates Table) ---

// SaveTemplate 在 'templates' 表中保存或更新一个模板.
func (cdb *ConfigDB) SaveTemplate(entry TemplateEntry) error {
	insertSQL := `
	INSERT INTO templates (filename, template_type, content)
	VALUES (?, ?, ?)
	ON CONFLICT(filename) DO UPDATE SET
		template_type = EXCLUDED.template_type,
		content = EXCLUDED.content,
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, entry.Filename, entry.TemplateType, entry.Content)
	if err != nil {
		return fmt.Errorf("db: failed to save template '%s': %w", entry.Filename, err)
	}
	return nil
}

// GetTemplate 从 'templates' 表中获取一个模板内容.
func (cdb *ConfigDB) GetTemplate(filename string) (*TemplateEntry, error) {
	querySQL := `SELECT filename, template_type, content, created_at, updated_at FROM templates WHERE filename = ?;`
	row := cdb.DB.QueryRow(querySQL, filename)

	entry := &TemplateEntry{}
	err := row.Scan(&entry.Filename, &entry.TemplateType, &entry.Content, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db: template '%s' not found: %w", filename, err)
		}
		return nil, fmt.Errorf("db: failed to get template '%s': %w", filename, err)
	}
	return entry, nil
}

// DeleteTemplate 从 'templates' 表中删除一个模板.
// 请注意: 此操作不会级联删除 'config_params' 或 'rendered_configs' 中的关联数据;
// 因为 'templates' 表不再是它们的外键父表.
func (cdb *ConfigDB) DeleteTemplate(filename string) error {
	_, err := cdb.DB.Exec(`DELETE FROM templates WHERE filename = ?;`, filename)
	if err != nil {
		return fmt.Errorf("db: failed to delete template '%s': %w", filename, err)
	}
	return nil
}

// RangeTempaltes 获取所有模板的名称
func (cdb *ConfigDB) RangeTemplates() ([]string, error) {
	querySQL := `SELECT filename FROM templates;`
	rows, err := cdb.DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get filenames from templates: %w", err)
	}
	defer rows.Close()

	var filenames []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, fmt.Errorf("db: failed to scan template filename: %w", err)
		}
		filenames = append(filenames, filename)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error during template rows iteration: %w", err)
	}

	return filenames, nil
}

// GetAllTempaltes
func (cdb *ConfigDB) GetAllTemplates() ([]TemplateEntry, error) {
	querySQL := `SELECT filename, template_type, content, created_at, updated_at FROM templates;`
	rows, err := cdb.DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get all templates: %w", err)
	}
	defer rows.Close()

	var templates []TemplateEntry
	for rows.Next() {
		var entry TemplateEntry
		if err := rows.Scan(&entry.Filename, &entry.TemplateType, &entry.Content, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db: failed to scan template entry: %w", err)
		}
		templates = append(templates, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error during templates rows iteration: %w", err)
	}

	return templates, nil
}

// --- 参数操作 (Config_Params Table) ---

/*
	filename        TEXT PRIMARY KEY,
	template_type   TEXT NOT NULL,
	params_gob      BLOB NOT NULL,
	params_origin   BLOB NOT NULL,
	created_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	updated_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
*/

// SaveParams 在 'config_params' 表中保存或更新一个模板的参数.
// entry.ParamsGOB 应该是一个经过 GOB 编码的字节切片.
func (cdb *ConfigDB) SaveParams(entry ParamsEntry) error {
	insertSQL := `
	INSERT INTO config_params (filename, template_type, params_gob, params_origin)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(filename) DO UPDATE SET
		template_type = EXCLUDED.template_type,
		params_gob = EXCLUDED.params_gob,
		params_origin = EXCLUDED.params_origin,
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, entry.Filename, entry.TemplateType, entry.ParamsGOB, entry.ParamsOrigin)
	if err != nil {
		return fmt.Errorf("db: failed to save params for '%s': %w", entry.Filename, err)
	}
	return nil
}

// GetParams 从 'config_params' 表中获取一个模板的参数.
// 返回的 ParamsGOB 是 GOB 编码的字节切片; 调用方需要自行解码.
func (cdb *ConfigDB) GetParams(filename string) (*ParamsEntry, error) {
	querySQL := `SELECT filename, template_type, params_gob, params_origin, created_at, updated_at FROM config_params WHERE filename = ?;`
	row := cdb.DB.QueryRow(querySQL, filename)

	entry := &ParamsEntry{}
	err := row.Scan(&entry.Filename, &entry.TemplateType, &entry.ParamsGOB, &entry.ParamsOrigin, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db: params for '%s' not found: %w", filename, err)
		}
		return nil, fmt.Errorf("db: failed to get params for '%s': %w", filename, err)
	}
	return entry, nil
}

// DeleteParams 从 'config_params' 表中删除一个模板的参数.
// 此操作将级联删除 'rendered_configs' 表中与该 filename 关联的所有渲染产物.
func (cdb *ConfigDB) DeleteParams(filename string) error {
	_, err := cdb.DB.Exec(`DELETE FROM config_params WHERE filename = ?;`, filename)
	if err != nil {
		return fmt.Errorf("db: failed to delete params for '%s': %w", filename, err)
	}
	return nil
}

// GetFileNames
func (cdb *ConfigDB) GetFileNames() ([]string, error) {
	querySQL := `SELECT filename FROM config_params;`
	rows, err := cdb.DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get filenames from config_params: %w", err)
	}
	defer rows.Close()

	var filenames []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, fmt.Errorf("db: failed to scan filename: %w", err)
		}
		filenames = append(filenames, filename)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error during rows iteration: %w", err)
	}

	return filenames, nil
}

// RangeAllParams
func (cdb *ConfigDB) RangeAllParams() ([]ParamsEntry, error) {
	querySQL := `SELECT filename, template_type, params_gob, params_origin, created_at, updated_at FROM config_params;`
	rows, err := cdb.DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get all params: %w", err)
	}
	defer rows.Close()

	var params []ParamsEntry
	for rows.Next() {
		var entry ParamsEntry
		if err := rows.Scan(&entry.Filename, &entry.TemplateType, &entry.ParamsGOB, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db: failed to scan params entry: %w", err)
		}
		params = append(params, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error during params rows iteration: %w", err)
	}

	return params, nil
}

// --- 渲染产物操作 (Rendered_Configs Table) ---

// SaveRenderedConfig 在 'rendered_configs' 表中保存或更新一个渲染后的配置文件内容.
// 注意: 该操作依赖于 'config_params' 表中已存在对应的 filename; 否则会违反外键约束.
func (cdb *ConfigDB) SaveRenderedConfig(entry RenderedConfigEntry) error {
	insertSQL := `
	INSERT INTO rendered_configs (filename, rendered_content, rendered_at)
	VALUES (?, ?, strftime('%s', 'now'))
	ON CONFLICT(filename) DO UPDATE SET
		rendered_content = EXCLUDED.rendered_content,
		rendered_at = strftime('%s', 'now'),
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, entry.Filename, entry.RenderedContent)
	if err != nil {
		return fmt.Errorf("db: failed to save rendered config for '%s': %w", entry.Filename, err)
	}
	return nil
}

// GetRenderedConfig 从 'rendered_configs' 表中获取一个渲染后的配置文件内容.
func (cdb *ConfigDB) GetRenderedConfig(filename string) (*RenderedConfigEntry, error) {
	querySQL := `SELECT filename, rendered_content, rendered_at, updated_at FROM rendered_configs WHERE filename = ?;`
	row := cdb.DB.QueryRow(querySQL, filename)

	entry := &RenderedConfigEntry{}
	err := row.Scan(&entry.Filename, &entry.RenderedContent, &entry.RenderedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db: rendered config '%s' not found: %w", filename, err)
		}
		return nil, fmt.Errorf("db: failed to get rendered config '%s': %w", filename, err)
	}
	return entry, nil
}

// DeleteRenderedConfig 从 'rendered_configs' 表中删除一个渲染后的配置文件内容.
func (cdb *ConfigDB) DeleteRenderedConfig(filename string) error {
	_, err := cdb.DB.Exec(`DELETE FROM rendered_configs WHERE filename = ?;`, filename)
	if err != nil {
		return fmt.Errorf("db: failed to delete rendered config for '%s': %w", filename, err)
	}
	return nil
}

// RangeAllReandered
func (cdb *ConfigDB) RangeAllReandered() ([]RenderedConfigEntry, error) {
	querySQL := `SELECT filename, rendered_content, rendered_at, updated_at FROM rendered_configs;`
	rows, err := cdb.DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("db: failed to get all rendered configs: %w", err)
	}
	defer rows.Close()

	var renderedConfigs []RenderedConfigEntry
	for rows.Next() {
		var entry RenderedConfigEntry
		if err := rows.Scan(&entry.Filename, &entry.RenderedContent, &entry.RenderedAt, &entry.UpdatedAt); err != nil {
			return nil, fmt.Errorf("db: failed to scan rendered config entry: %w", err)
		}
		renderedConfigs = append(renderedConfigs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: error during rendered configs rows iteration: %w", err)
	}

	return renderedConfigs, nil
}

// --- 全局配置操作 (Global_Configs Table) ---
/*
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS global_configs (
		filename        TEXT PRIMARY KEY,
		params          BLOB NOT NULL,
		tmpl_content    BLOB NOT NULL,
		rendered_content BLOB NOT NULL,
		updated_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);`)
*/

// SaveGlobalConfig 在 'global_configs' 表中保存或更新全局配置.
func (cdb *ConfigDB) SaveGlobalConfig(entry GlobalConfig) error {
	insertSQL := `
	INSERT INTO global_configs (filename, params, tmpl_content, rendered_content)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(filename) DO UPDATE SET
		params = EXCLUDED.params,
		tmpl_content = EXCLUDED.tmpl_content,
		rendered_content = EXCLUDED.rendered_content,
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, entry.Filename, entry.Params, entry.TmplContent, entry.RenderedContent)
	if err != nil {
		return fmt.Errorf("db: failed to save global config for '%s': %w", entry.Filename, err)
	}
	return nil
}

// GetGlobalConfig 从 'global_configs' 表中获取全局配置.
func (cdb *ConfigDB) GetGlobalConfig(filename string) (*GlobalConfig, error) {
	querySQL := `SELECT filename, params, tmpl_content, rendered_content, updated_at FROM global_configs WHERE filename = ?;`
	row := cdb.DB.QueryRow(querySQL, filename)

	entry := &GlobalConfig{}
	err := row.Scan(&entry.Filename, &entry.Params, &entry.TmplContent, &entry.RenderedContent, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db: global config '%s' not found: %w", filename, err)
		}
		return nil, fmt.Errorf("db: failed to get global config '%s': %w", filename, err)
	}
	return entry, nil
}

// DeleteGlobalConfig 从 'global_configs' 表中删除全局配置.
func (cdb *ConfigDB) DeleteGlobalConfig(filename string) error {
	_, err := cdb.DB.Exec(`DELETE FROM global_configs WHERE filename = ?;`, filename)
	if err != nil {
		return fmt.Errorf("db: failed to delete global config for '%s': %w", filename, err)
	}
	return nil
}

// SaveGlobalParams
func (cdb *ConfigDB) SaveGlobalParams(filename string, params []byte) error {
	insertSQL := `
	INSERT INTO global_configs (filename, params)
	VALUES (?, ?)
	ON CONFLICT(filename) DO UPDATE SET
		params = EXCLUDED.params,
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, filename, params)
	if err != nil {
		return fmt.Errorf("db: failed to save global params for '%s': %w", filename, err)
	}
	return nil
}

// SaveGlobalRenderedContent
func (cdb *ConfigDB) SaveGlobalRenderedContent(filename string, renderedContent []byte) error {
	insertSQL := `
	INSERT INTO global_configs (filename, rendered_content)
	VALUES (?, ?)
	ON CONFLICT(filename) DO UPDATE SET
		rendered_content = EXCLUDED.rendered_content,
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, filename, renderedContent)
	if err != nil {
		return fmt.Errorf("db: failed to save global rendered content for '%s': %w", filename, err)
	}
	return nil
}

// SaveGlobalTemplate
func (cdb *ConfigDB) SaveGlobalTemplate(filename string, tmplContent []byte) error {
	insertSQL := `
	INSERT INTO global_configs (filename, tmpl_content)
	VALUES (?, ?)
	ON CONFLICT(filename) DO UPDATE SET
		tmpl_content = EXCLUDED.tmpl_content,
		updated_at = strftime('%s', 'now');
	`
	_, err := cdb.DB.Exec(insertSQL, filename, tmplContent)
	if err != nil {
		return fmt.Errorf("db: failed to save global template for '%s': %w", filename, err)
	}
	return nil
}

// GetGlobalTemplate
func (cdb *ConfigDB) GetGlobalTemplate(filename string) ([]byte, error) {
	querySQL := `SELECT tmpl_content FROM global_configs WHERE filename = ?;`
	var tmplContent []byte
	err := cdb.DB.QueryRow(querySQL, filename).Scan(&tmplContent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db: global template '%s' not found: %w", filename, err)
		}
		return nil, fmt.Errorf("db: failed to get global template '%s': %w", filename, err)
	}
	return tmplContent, nil
}
