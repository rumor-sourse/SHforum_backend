package mysql

import (
	"SHforum_backend/models"
	"gorm.io/gorm/clause"
)

// GetCommentList 获取评论列表
func GetCommentList(postID int64) (comments []*models.Comment, err error) {
	comments = make([]*models.Comment, 0, 2)
	result := db.Debug().
		Select("id", "content", "username", "comment_like_count", "is_admin_comment", "parent_comment_username").
		Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}

// GetCommentByID 根据评论ID获取评论
func GetCommentByID(commentID int64) (comment *models.Comment, err error) {
	comment = new(models.Comment)
	result := db.Debug().Where("id=?", commentID).First(comment)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}

// CreateComment 创建评论
func CreateComment(p *models.Comment) (err error) {
	result := db.Debug().Create(p)
	if result.Error != nil {
		return result.Error
	}
	return
}

// DeleteComment 删除评论
func DeleteComment(commentID int64) (err error) {
	result := db.Debug().Where("id = ?", commentID).Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	return
}

// UpdateCommentLikeCount 更新评论点赞数
func UpdateCommentLikeCount(commentID int64, value int64) (err error) {
	result := db.Debug().Model(&models.Comment{}).Where("id = ?", commentID).Update("comment_like_count", value)
	if result.Error != nil {
		return result.Error
	}
	return
}

func GetCommentListByIDs(ids []string) (comments []*models.Comment, err error) {
	comments = make([]*models.Comment, len(ids))
	result := db.Debug().
		Select("id", "content", "username", "comment_like_count", "is_admin_comment", "parent_comment_username").
		Where("id IN ?", ids).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{ids}, WithoutParentheses: true},
	}).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}
