package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Content     string `gorm:"type:varchar(100);not null" json:"content"`
	AuthorID    int64  `gorm:"not null" json:"author_id"`
	CommunityID int64  `gorm:"not null" json:"community_id"`
	Status      int32  `gorm:"type:tinyint;not null;default:1" json:"status"`
}
