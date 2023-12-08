package logic

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/dao/redis"
	"SHforum_backend/logic/rabbitmq"
	"SHforum_backend/models"
	"SHforum_backend/models/response"
	"SHforum_backend/pkg/snowflake"
	"fmt"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成post_id
	p.ID = uint(snowflake.GenID())
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	//把帖子保存到redis
	err = redis.CreatePost(int64(p.ID), p.CommunityID)
	if err != nil {
		return err
	}
	return
}

func GetPostById(pid int64) (data *response.PostDetailResponse, err error) {
	//查询帖子详情
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	//查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("authorid", post.AuthorID),
			zap.Error(err))
		return
	}
	//查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("communityid", post.CommunityID),
			zap.Error(err))
		return
	}
	//组合数据
	data = &response.PostDetailResponse{
		AuthorName: user.Username,
		CommunityResponse: &response.CommunityResponse{
			ID:   community.ID,
			Name: community.Name,
		},
		PostResponse: &response.PostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Content:     post.Content,
			AuthorID:    post.AuthorID,
			CommunityID: post.CommunityID,
			Status:      post.Status,
		},
	}
	return
}

func GetPostList(page, size int64) (data []*response.PostDetailResponse, err error) {
	//查询帖子列表
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList(page, size) failed",
			zap.Int64("page", page),
			zap.Int64("size", size),
			zap.Error(err))
		return
	}
	//遍历每个帖子，查询对应的作者信息
	for _, post := range posts {
		//查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		//组合数据
		postDetail := &response.PostDetailResponse{
			AuthorName: user.Username,
			CommunityResponse: &response.CommunityResponse{
				ID:   community.ID,
				Name: community.Name,
			},
			PostResponse: &response.PostResponse{
				ID:          post.ID,
				Title:       post.Title,
				Content:     post.Content,
				AuthorID:    post.AuthorID,
				CommunityID: post.CommunityID,
				Status:      post.Status,
			},
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*response.PostDetailResponse, err error) {
	//从redis拿到所有的id
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	//根据id去数据库查询帖子详细信息
	// 返回的数据还要按照给定的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇贴子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//遍历每个帖子，查询对应的作者信息
	for idx, post := range posts {
		//查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		//组合数据
		postDetail := &response.PostDetailResponse{
			AuthorName: user.Username,
			VoteNum:    voteData[idx],
			CommunityResponse: &response.CommunityResponse{
				ID:   community.ID,
				Name: community.Name,
			},
			PostResponse: &response.PostResponse{
				ID:          post.ID,
				Title:       post.Title,
				Content:     post.Content,
				AuthorID:    post.AuthorID,
				CommunityID: post.CommunityID,
				Status:      post.Status,
			},
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*response.PostDetailResponse, err error) {
	//从redis拿到所有的id
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	//根据id去数据库查询帖子详细信息
	// 返回的数据还要按照给定的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇贴子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//遍历每个帖子，查询对应的作者信息
	for idx, post := range posts {
		//查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		//组合数据
		postDetail := &response.PostDetailResponse{
			AuthorName: user.Username,
			VoteNum:    voteData[idx],
			CommunityResponse: &response.CommunityResponse{
				ID:   community.ID,
				Name: community.Name,
			},
			PostResponse: &response.PostResponse{
				ID:          post.ID,
				Title:       post.Title,
				Content:     post.Content,
				AuthorID:    post.AuthorID,
				CommunityID: post.CommunityID,
				Status:      post.Status,
			},
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将两个查询列表的函数合二为一
func GetPostListNew(p *models.ParamPostList) (data []*response.PostDetailResponse, err error) {
	if p.CommunityID == 0 {
		//查询所有的帖子
		data, err = GetPostList2(p)
	} else {
		//查询社区的帖子
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}

func MQSendCreatePostMessage(userID int64) {
	rmq := rabbitmq.NewRabbitMQPubSub("new_post")
	defer rmq.Destroy()
	msg := fmt.Sprintf("您关注的用户%d创建了一条新帖子", userID)
	rmq.PublishCreatePostMessage(msg)
}

func MQReceiveCreatePostMessage(userID int64) {
	//查找该用户所有粉丝
	fans, err := GetFanList(userID)
	if err != nil {
		zap.L().Error("logic.GetFanList() failed", zap.Error(err))
		return
	}
	rmq := rabbitmq.NewRabbitMQPubSub("new_post")
	defer rmq.Destroy()
	for _, fan := range fans {
		rmq.ConsumeCreatePostMessage(userID, fan)
	}
}
