package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   int64  `gorm:"not null;" json:"userid"`
	Username string `gorm:"type:varchar(100);not null;" json:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
	Email    string `gorm:"type:varchar(100);not null;unique" json:"email"`
}

type Follow struct {
	UserID       int64 `gorm:"not null;" json:"userid"`        //用户
	FollowedUser int64 `gorm:"not null;" json:"followed_user"` //关注的用户
	gorm.Model
}

type Fan struct {
	UserID  int64 `gorm:"not null;" json:"userid"`   //用户
	FanUser int64 `gorm:"not null;" json:"fan_user"` //该用户的粉丝
	gorm.Model
}
