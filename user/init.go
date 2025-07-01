package user

import (
	"caddydash/db"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func InitAdminUser(username string, password string, cdb *db.ConfigDB) error {
	hasUser, err := cdb.HasAnyUser()
	if err != nil {
		return fmt.Errorf("failed to check if any user exists: %w", err)
	}
	if hasUser {
		userStatus.SetInitialized(true)
		return nil
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	err = cdb.AddUser(username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to add admin user: %w", err)
	}
	userStatus.SetInitialized(true)
	return nil
}

// bcrypt加密password串
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func InitAdminUserStatus(cdb *db.ConfigDB) error {
	hasUser, err := cdb.HasAnyUser()
	if err != nil {
		return fmt.Errorf("failed to check if any user exists: %w", err)
	}
	if hasUser {
		userStatus.SetInitialized(true)
		return nil
	} else {
		userStatus.SetInitialized(false)
		return nil
	}
}

func InitFormEnv(cdb *db.ConfigDB) error {
	username := os.Getenv("CADDYDASH_USERNAME")
	password := os.Getenv("CADDYDASH_PASSWORD")

	if username != "" && password != "" {
		// 检查是否已经有用户
		hasUser, err := cdb.HasAnyUser()
		if err != nil {
			return fmt.Errorf("failed to check if any user exists: %w", err)
		}
		if hasUser {
			// 如果已经有用户，则不执行初始化，但设置状态为已初始化
			userStatus.SetInitialized(true)
			return nil
		}

		// 执行初始化
		hashedPassword, err := HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password from env: %w", err)
		}
		err = cdb.AddUser(username, hashedPassword)
		if err != nil {
			return fmt.Errorf("failed to add admin user from env: %w", err)
		}
		userStatus.SetInitialized(true)
		return nil
	}

	// 如果环境变量不存在，则根据数据库状态设置初始化状态
	return InitAdminUserStatus(cdb)
}
