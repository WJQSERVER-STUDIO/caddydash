package api

import (
	"caddydash/config"
	"caddydash/db"
	"caddydash/user"
	"net/http"
	"strings"

	"github.com/fenthope/sessions"
	"github.com/infinite-iroha/touka"
)

var (
	exactMatchPaths = map[string]struct{}{
		"/login":              {},
		"/login.html":         {},
		"/v0/api/auth/login":  {},
		"/v0/api/auth/logout": {},
		"/v0/api/auth/init":   {},
		"/init.html":          {},
		"/favicon.ico":        {},
		"/v0/api/info":        {},
	}
	prefixMatchPaths = []string{ // 保持前缀匹配，因为数量少
		"/js/",
		"/css/",
		"/locales",
	}
	loginMatchPaths = map[string]struct{}{
		"/login":             {},
		"/login.html":        {},
		"/v0/api/auth/login": {},
	}
	initMatchPaths = map[string]struct{}{
		"/v0/api/auth/init": {},
		"/init.html":        {},
	}
)

func isPassPath(requestPath string) bool {
	// 前缀匹配
	for _, prefix := range prefixMatchPaths {
		if strings.HasPrefix(requestPath, prefix) {
			return true
		}
	}

	// 精确匹配
	if _, ok := exactMatchPaths[requestPath]; ok {
		return true
	}

	return false
}

func isLoginPath(requestPath string) bool {
	if _, ok := loginMatchPaths[requestPath]; ok {
		return true
	}
	return false
}

func isInitPath(requestPath string) bool {
	if _, ok := initMatchPaths[requestPath]; ok {
		return true
	}
	return false
}

func SessionMiddleware(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		session := sessions.Default(c)
		requestPath := c.Request.URL.Path
		pass := isPassPath(requestPath)
		if !user.IsAdminInit() && !pass || !user.IsAdminInit() && isLoginPath(requestPath) {
			c.Redirect(http.StatusFound, "/init.html")
			c.Abort()
			return
		} else if user.IsAdminInit() && isInitPath(requestPath) {
			c.Redirect(http.StatusFound, "/login.html")
			c.Abort()
			return
		}

		if session.Get("authenticated") != true && !pass {
			c.Redirect(http.StatusFound, "/login.html")
			c.Abort()
			return
		}
		c.Next()

	}
}

func AuthLogin(c *touka.Context, cfg *config.Config, cdb *db.ConfigDB) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 输入验证
	if username == "" || password == "" {
		c.Errorf("Username or password not provided")
		c.JSON(http.StatusBadRequest, touka.H{"error": "Need username and password"})
		return
	}
	c.Infof("user login: %s password: %s", username, password)

	// 验证账户密码
	pass, err := user.CheckLogin(username, password, cdb)
	if err != nil {
		c.Errorf("Failed to check login: %v", err)
		c.JSON(http.StatusInternalServerError, touka.H{"error": "Internal Auth Check Error"})
		return
	}
	if !pass {
		c.Errorf("Invalid username or password")
		c.JSON(http.StatusUnauthorized, touka.H{"error": "Invalid username or password"})
		return
	}
	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Save()
	c.Infof("Login successful for user: %s", username)
	c.JSON(http.StatusOK, touka.H{"success": true})
}

func AuthLogout(c *touka.Context) {

	session := sessions.Default(c)
	session.Set("authenticated", false)
	session.Clear()
	err := session.Save()
	if err != nil {
		c.Errorf("Failed to save session: %v", err)
		c.JSON(http.StatusInternalServerError, touka.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusFound, "/login.html")

}

func ResetPassword(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		username := c.PostForm("username")
		oldPassword := c.PostForm("old_password")
		newPassword := c.PostForm("new_password")
		// 验证是否为空
		if username == "" || oldPassword == "" || newPassword == "" {
			c.JSON(400, touka.H{"error": "username and password are required"})
			return
		}
		// 验证用户是否存在
		exist, err := cdb.IsUserExists(username)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		if !exist {
			//不正确的参数
			c.JSON(400, touka.H{"error": "user not exist"})
			return
		}
		// 是否可以重置
		ok, err := user.CheckLogin(username, oldPassword, cdb)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		if !ok {
			// 错误的密码
			c.JSON(400, touka.H{"error": "current password is not correct"})
			return
		}
		// 更新密码
		hashpwd, err := user.HashPassword(newPassword)
		if err != nil {
			c.Errorf("Failed to hash password: %v", err)
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		err = cdb.UpdateUserPassword(username, hashpwd)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		// 进行logout
		AuthLogout(c)
	}
}

func AuthInitHandle(cdb *db.ConfigDB) touka.HandlerFunc {
	return func(c *touka.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		// 验证是否为空
		if username == "" || password == "" {
			c.JSON(400, touka.H{"error": "username and password are required"})
			return
		}
		// 初始化管理员
		err := user.InitAdminUser(username, password, cdb)
		if err != nil {
			c.JSON(500, touka.H{"error": err.Error()})
			return
		}
		c.JSON(200, touka.H{"message": "admin initialized"})
	}
}

func AuthInitStatus() touka.HandlerFunc {
	return func(c *touka.Context) {
		// 返回是否init管理员
		isInit := user.IsAdminInit()
		if isInit {
			c.JSON(200, touka.H{"admin_init": true})
		} else {
			c.JSON(200, touka.H{"admin_init": false})
		}
	}
}
