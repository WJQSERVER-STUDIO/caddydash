package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
	//_ "github.com/mattn/go-sqlite3"
)

// TemplateEntry 对应 templates 表
type TemplateEntry struct {
	Filename     string `json:"filename"`
	TemplateType string `json:"template_type"`
	Content      []byte `json:"content"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// ParamsEntry 对应 config_params 表
type ParamsEntry struct {
	Filename     string `json:"filename"`
	TemplateType string `json:"template_type"`
	ParamsGOB    []byte `json:"params_gob"`
	ParamsOrigin []byte `json:"params_origin"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// RenderedConfigEntry 对应 rendered_configs 表
type RenderedConfigEntry struct {
	Filename        string `json:"filename"`
	RenderedContent []byte `json:"rendered_content"`
	RenderedAt      int64  `json:"rendered_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// UsersTable
type UsersTable struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// GlobalConfig 全局CaddyFile配置专用table
type GlobalConfig struct {
	Filename        string `json:"filename"`
	Params          []byte `json:"params"`
	TmplContent     []byte `json:"tmpl_content"`
	RenderedContent []byte `json:"rendered_content"`
	UpdatedAt       int64  `json:"updated_at"`
}

// ConfigDB
type ConfigDB struct {
	DB *sql.DB
}

func InitDB(filepath string) (*ConfigDB, error) {
	db, err := loadDB(filepath)
	if err != nil {
		return nil, err
	}
	cdb := &ConfigDB{DB: db}
	// 尝试 Ping 数据库以验证连接是否成功
	if err = db.Ping(); err != nil {
		db.Close() // 如果连接失败，确保关闭数据库句柄
		return nil, fmt.Errorf("db: failed to connect to database: %w", err)
	}
	err = cdb.createTables()
	if err != nil {
		return nil, err
	}
	return cdb, nil
}

func loadDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?cache=shared&_journal=WAL", filepath))
	if err != nil {
		log.Fatalf("Can not connect to database: %v", err)
		return nil, err
	}
	return db, nil
}

func (cdb *ConfigDB) createTables() error {
	tx, err := cdb.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for table creation: %w", err)
	}
	defer tx.Rollback() // 确保在出错时回滚事务

	// 创建 templates 表
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS templates (
		filename        TEXT PRIMARY KEY,
		template_type   TEXT NOT NULL,
		content         BLOB NOT NULL,
		created_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);`)
	if err != nil {
		return fmt.Errorf("failed to create 'templates' table: %w", err)
	}

	// 2. 创建 config_params 表
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS config_params (
		filename        TEXT PRIMARY KEY,
		template_type   TEXT NOT NULL,
		params_gob      BLOB NOT NULL,
		params_origin   BLOB NOT NULL,
		created_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		-- FOREIGN KEY(filename) REFERENCES templates(filename) ON DELETE CASCADE -- 已移除
	);`)
	if err != nil {
		return fmt.Errorf("failed to create 'config_params' table: %w", err)
	}

	// 3. 创建 rendered_configs 表 (外键指向 config_params 表)
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS rendered_configs (
		filename        TEXT PRIMARY KEY,
		rendered_content BLOB NOT NULL,
		rendered_at     INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		FOREIGN KEY(filename) REFERENCES config_params(filename) ON DELETE CASCADE -- 修改外键引用
	);`)
	if err != nil {
		return fmt.Errorf("failed to create 'rendered_configs' table: %w", err)
	}

	// 4. 创建 users 表
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL,
		created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);`)
	if err != nil {
		return fmt.Errorf("failed to create 'users' table: %w", err)
	}

	// 5. 创建 global_configs 表
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS global_configs (
		filename        TEXT PRIMARY KEY,
		params          BLOB NOT NULL DEFAULT '',
		tmpl_content    BLOB NOT NULL DEFAULT '',
		rendered_content BLOB NOT NULL DEFAULT '',
		updated_at      INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);`)
	if err != nil {
		return fmt.Errorf("failed to create 'global_configs' table: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for table creation: %w", err)
	}

	return nil
}

func (cdb *ConfigDB) CloseDB() error {
	if cdb.DB == nil {
		return nil // 数据库未打开或已关闭
	}
	err := cdb.DB.Close()
	if err != nil {
		return fmt.Errorf("db: failed to close database: %w", err)
	}
	return nil
}
