package redis

import (
	"SHforum_backend/util/email"
	"strconv"
	"time"
)

// SaveCode 将验证码存入redis
func SaveCode(em string, code string) {
	emailExpireTime, err := strconv.Atoi(email.EmailExpireTime)
	if err != nil {
		return
	}
	client.Set(ctx, em, code, time.Duration(emailExpireTime)*time.Minute)
}

// GetCode 从redis中获取验证码
func GetCode(email string) (code string, err error) {
	return client.Get(ctx, email).Result()
}
