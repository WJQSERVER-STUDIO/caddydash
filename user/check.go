package user

import (
	"caddydash/db"

	"golang.org/x/crypto/bcrypt"
)

// 判断是否可以登陆
func CheckLogin(username, password string, cdb *db.ConfigDB) (bool, error) {
	// 判断数据库内是否存在username
	userExist, err := cdb.IsUserExists(username)
	if err != nil {
		return false, err
	}
	if !userExist {
		return false, nil
	}
	passwordb, err := cdb.GetPasswordByUsername(username)
	if err != nil {
		return false, err
	}
	// 校验密码
	check, err := checkPasswordHash(password, passwordb)
	if err != nil {
		return false, err
	}
	return check, nil
}

func IsAdminInit() bool {
	return userStatus.IsUserInitialized()
}

// 校验密码, 避免时序攻击问题
func checkPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
