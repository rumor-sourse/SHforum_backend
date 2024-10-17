package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID               int64 `gorm:"primarykey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Content          string         `gorm:"type:varchar(1000);not null;comment:'评论内容'" json:"content"`
	PostID           int64          `gorm:"not null;index;comment:'所属帖子ID'" json:"post_id"`
	UserID           int64          `gorm:"not null;index;comment:'所属用户ID'" json:"user_id"`
	UserName         string         `gorm:"type:varchar(100);not null;comment:'所属用户名';column:username" json:"username"`
	CommentLikeCount int64          `gorm:"type:bigint;not null;default:0;comment:'点赞数'" json:"comment_like_count"`
	IsAdminComment   bool           `gorm:"type:tinyint;not null;default:0;comment:'是否是贴主本人评论'" json:"is_admin_comment"`
	//为-1则为顶级评论，否则则为父评论的ID
	ParentCommentID       int64  `gorm:"comment:'父评论ID':not null" json:"parent_comment_id"`
	ParentCommentUserName string `gorm:"type:varchar(100);comment:'父评论用户名';column:parent_comment_username" json:"parent_comment_username"`
}
