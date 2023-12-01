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
	UserID       int64 `gorm:"not null;" json:"userid"`       //用户
	FollowedUser int64 `gorm:"not null;" json:"followeduser"` //关注的用户
	gorm.Model
}
