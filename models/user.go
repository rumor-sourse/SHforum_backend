package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   int64  `gorm:"not null;" json:"userid"`
	Username string `gorm:"type:varchar(100);not null;" json:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
	Email    string `gorm:"type:varchar(100);not null;unique" json:"email"`
}
