package models

import "gorm.io/gorm"

type Message struct {
	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Content     string `gorm:"type:varchar(100);not null" json:"content"`
	SendUser    int64  `gorm:"not null" json:"send_user"`
	ReceiveUser int64  `gorm:"not null" json:"receive_user"`
	HadRead     string `gorm:"type:varchar(100);not null" json:"had_read" binding:"oneof=1 0"`
	gorm.Model
}
