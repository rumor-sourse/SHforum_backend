package mysql

import (
	"SHforum_backend/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

const secret = "SHforum_backend"

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	//select count(userid) from user where username=?
	var user *models.User
	var count int64
	db.Debug().Model(user).Where("username=?", username).Select("userid").Count(&count)
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//insert into user(userid, username, password, email) values(?, ?, ?, ?)
	user.Password = encryptPassword(user.Password)
	db.Debug().Create(user)
	return
}

// encryptPassword 对密码进行加密
func encryptPassword(opassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(opassword)))
}

// Login 用户登录
func Login(user *models.User) (err error) {
	oPassword := user.Password
	//select userid, username, password from user where username=?
	result := db.Debug().Select("userid", "username", "password").Where("username=?", user.Username).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserById 根据用户ID获取用户信息
func GetUserById(uid int64) (data *models.User, err error) {
	data = new(models.User)
	//select user_id, username from user where user_id = ?
	result := db.Debug().Select("userid", "username").Where("userid = ?", uid).First(data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = ErrorInvalidID
	}
	return
}
