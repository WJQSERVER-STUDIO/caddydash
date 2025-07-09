package main

import (
	"caddydash/api"
	"caddydash/apic"
	"caddydash/config"
	"caddydash/db"
	"caddydash/gen"
	"caddydash/user"
	"crypto/rand"
	"encoding/gob"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fenthope/compress"
	"github.com/fenthope/reco"
	"github.com/fenthope/record"
	"github.com/fenthope/sessions"
	"github.com/fenthope/sessions/cookie"
	"github.com/infinite-iroha/touka"
	"github.com/klauspost/compress/zstd"
	_ "modernc.org/sqlite"
)

var (
	cfg        *config.Config
	cfgfile    string
	cdb        *db.ConfigDB
	sessionKey []byte
	version    string
)

func init() {
	parseFlags()
	loadConfig()
	loadDatabase(cfg.DB.Filepath)
	loadtmpltoDB(cfg.Server.CaddyDir+"tmpl", cdb)
	loadAdminStatus(cdb)
	initSessionKey()
}

func parseFlags() {
	//posix
	flag.StringVar(&cfgfile, "c", "./config.toml", "Path to the configuration file")
	flag.Parse()
}

func loadConfig() {
	var err error
	cfg, err = config.LoadConfig(cfgfile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		// 如果配置文件加载失败，也显示帮助信息并退出
		flag.Usage()
		os.Exit(1)
	}
	if cfg != nil && cfg.Server.Debug { // 确保 cfg 不为 nil
		fmt.Println("Config File Path: ", cfgfile)
		fmt.Printf("Loaded config: %v\n", cfg)
	}
	fmt.Printf("Loaded config: %v\n", cfg)
}

func loadDatabase(filepath string) {
	var err error
	cdb, err = db.InitDB(filepath)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
}

func loadAdminStatus(cdb *db.ConfigDB) {
	err := user.InitFormEnv(cdb)
	if err != nil {
		fmt.Printf("Failed to initialize admin user status: %v\n", err)
		os.Exit(1)
	}

	err = user.InitAdminUserStatus(cdb)
	if err != nil {
		fmt.Printf("Failed to initialize admin user status: %v\n", err)
		os.Exit(1)
	}
}

func loadtmpltoDB(path string, cdb *db.ConfigDB) {
	err := gen.ReadTmplToDB(path, cdb)
	if err != nil {
		fmt.Printf("Failed to load templates: %v\n", err)
		os.Exit(1)
	}
	err = gen.SetGlobalConfig(cfg, cdb)
	if err != nil {
		fmt.Printf("Failed to set global config: %v\n", err)
		os.Exit(1)
	}
	err = gen.Add80SiteConfig(cfg, cdb)
	if err != nil {
		fmt.Printf("Failed to add :80 site config: %v\n", err)
		os.Exit(1)
	}
}

func initSessionKey() {
	// crypto 生成随机
	sessionKey = make([]byte, 32)
	_, err := rand.Read(sessionKey)
	if err != nil {
		fmt.Printf("Failed to generate session key: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	defer cdb.CloseDB()

	r := touka.Default()
	logger, err := reco.New(reco.Config{
		Level:          reco.LevelInfo,
		Mode:           reco.ModeText,
		FilePath:       cfg.Server.CaddyDir + "log/caddydash.log",
		MaxFileSizeMB:  5,
		EnableRotation: true,
		Async:          true,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()
	r.SetLogger(logger)
	r.Use(record.Middleware())
	setCacheControl(r)

	r.Use(compress.Compression(compress.CompressOptions{
		// Algorithms: 配置每种压缩算法的级别和是否启用对象池
		Algorithms: map[string]compress.AlgorithmConfig{
			compress.EncodingGzip: {
				Level:       -1,
				PoolEnabled: true,
			},
			compress.EncodingDeflate: {
				Level:       -1,    // Deflate默认压缩比
				PoolEnabled: false, // Deflate不启用对象池
			},
			compress.EncodingZstd: {
				Level:       int(zstd.SpeedBestCompression), // Zstandard最佳压缩比
				PoolEnabled: true,                           // 启用Zstandard压缩器的对象池
			},
		},

		// MinContentLength: 响应内容达到此字节数才进行压缩 (例如 1KB)
		MinContentLength: 512,

		// CompressibleTypes: 只有响应的 Content-Type 匹配此列表中的MIME类型前缀才进行压缩
		CompressibleTypes: compress.DefaultCompressibleTypes,

		// EncodingPriority: 当客户端接受多种支持的压缩算法时，服务器选择的优先级顺序
		EncodingPriority: []string{
			compress.EncodingZstd,
			compress.EncodingGzip,
			compress.EncodingDeflate,
		},
	}))

	store := cookie.NewStore(sessionKey)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   10800, // 3 hours
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("mysession", store))
	// 应用 session 中间件
	r.Use(api.SessionMiddleware(cdb))

	v0 := r.Group("/v0")
	api.ApiGroup(v0, cdb, cfg, version)

	gob.Register(map[string]interface{}{})
	gob.Register(time.Time{})
	gob.Register(gen.CaddyUniConfig{})

	frontendFS := os.DirFS("frontend")
	r.SetUnMatchFS(http.FS(frontendFS))

	// 打印定义的路由
	fmt.Println("Registered Routes:")
	for _, info := range r.GetRouterInfo() {
		fmt.Printf("  Method: %-7s Path: %-25s Handler: %-40s Group: %s\n", info.Method, info.Path, info.Handler, info.Group)
	}

	go func() {
		err := apic.RunCaddy(cfg)
		if err != nil {
			fmt.Printf("Failed to start caddy: %v\n", err)
			os.Exit(1)
		}
	}()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	r.Run(addr)
}
