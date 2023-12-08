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
	db.Debug().Model(user).Where("username=?", username).Select("user_id").Count(&count)
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// CheckFollowExist 检查是否已经关注
func CheckFollowExist(userId int64, followeduser int64) (err error) {
	//select count(userid) from user where username=?
	var follow *models.Follow
	var count int64
	db.Debug().Model(follow).Where("user_id=? and followed_user=?", userId, followeduser).Select("user_id").Count(&count)
	if count > 0 {
		return ErrorFollowExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//insert into user(userid, username, password, email) values(?, ?, ?, ?)
	user.Password = encryptPassword(user.Password)
	result := db.Debug().Create(user)
	if result.Error != nil {
		return result.Error
	}
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
	result := db.Debug().Select("user_id", "username", "password").Where("username=?", user.Username).First(user)
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
	result := db.Debug().Select("user_id", "username").Where("user_id = ?", uid).First(data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = ErrorInvalidID
	}
	return
}

// Follow 关注：userId关注了followeduser
func Follow(userId int64, followeduser int64) (err error) {
	//检查是否已经关注
	if err := CheckFollowExist(userId, followeduser); err != nil {
		return err
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//insert into follow(user_id, followed_user) values(?, ?)
	follow := &models.Follow{
		UserID:       userId,
		FollowedUser: followeduser,
	}
	fan := &models.Fan{
		UserID:  followeduser,
		FanUser: userId,
	}
	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Debug().Create(follow).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Create(fan).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnFollow 取消关注：userId取消关注了followeduser
func UnFollow(userId int64, followeduser int64) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	//delete from follow where user_id=? and followed_user=?
	tx.Debug().Where("user_id=? and followed_user=?", userId, followeduser).Delete(&models.Follow{})
	//delete from fan where user_id=? and fan_user=?
	tx.Debug().Where("user_id=? and fan_user=?", followeduser, userId).Delete(&models.Fan{})
	return tx.Commit().Error
}

// GetFollowList 获取关注用户列表
func GetFollowList(userId int64) (followList []*models.Follow, err error) {
	//select followed_user from follow where user_id=?
	result := db.Debug().Select("followed_user").Where("user_id=?", userId).Find(&followList)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}

// GetFanList 获取粉丝列表
func GetFanList(userId int64) (fanList []*models.Fan, err error) {
	//select fan_user from fan where user_id=?
	result := db.Debug().Select("fan_user").Where("user_id=?", userId).Find(&fanList)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}
