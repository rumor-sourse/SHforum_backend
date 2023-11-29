package redis

import (
	"SHforum_backend/util"
	"strconv"
	"time"
)

// SaveCode 将验证码存入redis
func SaveCode(email string, code string) {
	emailExpireTime, err := strconv.Atoi(util.EmailExpireTime)
	if err != nil {
		return
	}
	client.Set(email, code, time.Duration(emailExpireTime)*time.Minute)
}

// GetCode 从redis中获取验证码
func GetCode(email string) (code string, err error) {
	return client.Get(email).Result()
}
