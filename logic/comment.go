package logic

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/dao/redis"
	"SHforum_backend/models"
	"SHforum_backend/models/response"
	"SHforum_backend/rabbitmq"
	"SHforum_backend/util/snowflake"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

// GetCommentList 根据帖子id获取评论列表
func GetCommentList(p *models.ParamCommentList) (data []*response.CommentResponse, err error) {
	//从redis拿到所有的id
	ids, err := redis.GetPostCommentIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostCommentIDs(p) return 0 data")
		return
	}
	// 根据id去数据库查询评论
	comments, err := mysql.GetCommentListByIDs(ids)
	if err != nil {
		return
	}
	//组合数据
	for _, comment := range comments {
		co := &response.CommentResponse{
			ID:                    comment.ID,
			Content:               comment.Content,
			UserName:              comment.UserName,
			CommentLikeCount:      comment.CommentLikeCount,
			IsAdminComment:        comment.IsAdminComment,
			ParentCommentUserName: comment.ParentCommentUserName,
		}
		data = append(data, co)
	}
	return
}

// CreateComment 创建评论
func CreateComment(p *models.Comment) (err error) {
	//获取comment_id
	p.ID = snowflake.GenID("comment")
	//获取UserName
	user, err := mysql.GetUserById(p.UserID)
	if err != nil {
		return err
	}
	p.UserName = user.Username
	//获取ParentCommentName
	if p.ParentCommentID != -1 {
		comment, err := mysql.GetCommentByID(p.ParentCommentID)
		if err != nil {
			return err
		}
		p.ParentCommentUserName = comment.UserName
	}
	//获取IsAdminComment
	post, err := mysql.GetPostById(p.PostID)
	if err != nil {
		return err
	}
	if post.AuthorID != p.UserID {
		p.IsAdminComment = false
	} else {
		p.IsAdminComment = true
	}
	err = mysql.CreateComment(p)
	if err != nil {
		return err

	}
	//把评论保存到redis
	err = redis.CreateComment(p.ID, post.ID)
	if err != nil {
		return err
	}
	//发送消息到rabbitmq
	MQCreateCommentMessage(p, post)
	return
}

// GetCommentByID 根据评论ID获取评论
func GetCommentByID(commentID int64) (data *response.CommentResponse, err error) {
	//查询评论
	comment, err := mysql.GetCommentByID(commentID)
	if err != nil {
		zap.L().Error("mysql.GetCommentByID(commentID) failed",
			zap.Int64("commentID", commentID),
			zap.Error(err))
		return
	}
	//组合数据
	data = &response.CommentResponse{
		ID:                    comment.ID,
		Content:               comment.Content,
		UserID:                comment.UserID,
		UserName:              comment.UserName,
		CommentLikeCount:      comment.CommentLikeCount,
		IsAdminComment:        comment.IsAdminComment,
		ParentCommentUserName: comment.ParentCommentUserName,
	}
	return
}

// DeleteComment 删除评论
func DeleteComment(commentID int64) (err error) {
	return mysql.DeleteComment(commentID)
}

// LikeComment 给评论点赞
func LikeComment(userID int64, p *models.ParamCommentLike) (err error) {
	err = redis.CommentLike(strconv.Itoa(int(userID)), p.CommentID, float64(p.Direction))
	if err != nil {
		return err
	}
	// 获取点赞数
	likeCount, err := redis.GetCommentLikeCount(p.CommentID)
	//发送消息到rabbitmq
	MQCreateUpdateCommentLikeMessgeByMysql(p.CommentID, likeCount)
	return
}

func GetHotComment(pid int64) (data *response.CommentResponse, err error) {
	postID := strconv.Itoa(int(pid))
	mutextname := fmt.Sprintf("hotcomment:%s", postID)
	mutex := redis.RedSync.NewMutex(mutextname)
	hotcomment, err := redis.GetHotComment(postID)
	//缓存过期，redis中没有数据
	if errors.Is(err, redis.Nil) {
		//申请分布式锁
		if err := mutex.Lock(); err != nil {
			panic(err)
		}
		//让其中一个线程从mysql中读取信息
		hc, err := mysql.GetHotComment(pid)
		if err != nil {
			return nil, err
		}
		hotcomment = map[string]string{
			"ID":                    strconv.FormatInt(hc.ID, 10),
			"CreatedAt":             hc.CreatedAt.Format(models.TimeFormatStr),
			"UpdatedAt":             hc.UpdatedAt.Format(models.TimeFormatStr),
			"Content":               hc.Content,
			"PostID":                strconv.FormatInt(hc.PostID, 10),
			"UserID":                strconv.FormatInt(hc.UserID, 10),
			"UserName":              hc.UserName,
			"CommentLikeCount":      strconv.FormatInt(hc.CommentLikeCount, 10),
			"IsAdminComment":        strconv.FormatBool(hc.IsAdminComment),
			"ParentCommentID":       strconv.FormatInt(hc.ParentCommentID, 10),
			"ParentCommentUserName": hc.ParentCommentUserName,
		}
		//更新缓存信息
		err = redis.UpdateHotComment(postID, hotcomment)
		//释放锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			panic("unlock failed")
		}

	}
	return nil, nil
}

func MQCreateCommentMessage(p *models.Comment, post *models.Post) {
	rmq := rabbitmq.NewRabbitMQSimple("new_comment")
	defer rmq.Destroy()
	rmq.PublishCreateCommentMessage(p, post)
}

func MQReceiveCreateCommentMessageByMysql() {
	rmq := rabbitmq.NewRabbitMQSimple("new_comment")
	defer rmq.Destroy()
	rmq.ConsumeCreateCommentMessageByMysql()
}

func MQCreateUpdateCommentLikeMessgeByMysql(CommentID string, likeCount int64) {
	rmq := rabbitmq.NewRabbitMQSimple("update_comment_like")
	defer rmq.Destroy()
	rmq.PublishUpdateCommentLikeMessge(CommentID, likeCount)
}

func MQReceiveUpdateCommentLikeMessgeByMysql() {
	rmq := rabbitmq.NewRabbitMQSimple("update_comment_like")
	defer rmq.Destroy()
	rmq.ConsumeUpdateCommentLikeMessgeByMysql()
}
