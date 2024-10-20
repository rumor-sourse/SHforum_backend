package logic

import (
	"SHforum_backend/internal/dao/mysql"
	"SHforum_backend/internal/dao/redis"
	"SHforum_backend/internal/models"
	"SHforum_backend/internal/models/response"
	"SHforum_backend/internal/settings"
	rpc "SHforum_backend/pb"
	"SHforum_backend/pkg/es"
	"SHforum_backend/pkg/rabbitmq"
	"SHforum_backend/pkg/snowflake"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"time"
)

func CreatePost(p *models.Post) (err error) {
	//生成post_id
	p.ID = snowflake.GenID("post")
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
	//创建生产者告诉其粉丝有新帖子产生
	MQSendCreatePostMessage(p.AuthorID, *p)
	return
}

func UpdatePost(p *models.Post) (err error) {
	return mysql.UpdatePost(p)
}
func DeletePost(pid int64) (err error) {
	return mysql.DeletePost(pid)
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

func GetPostListAll(p *models.ParamPostList) (data []*response.PostDetailResponse, err error) {
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

func GetPostListByCommunity(p *models.ParamPostList) (data []*response.PostDetailResponse, err error) {
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

// GetPostList 查询贴子列表
func GetPostList(p *models.ParamPostList) (data []*response.PostDetailResponse, err error) {
	if p.CommunityID == 0 {
		// 查询所有的帖子
		data, err = GetPostListAll(p)
	} else {
		// 查询社区的帖子
		data, err = GetPostListByCommunity(p)
	}
	if err != nil {
		zap.L().Error("GetPostList failed", zap.Error(err))
		return nil, err
	}
	return
}

func Share(longUrl string) (shortUrl string, err error) {
	fmt.Println(longUrl)
	target := settings.Conf.RpcServerConfig.Host + ":" + strconv.FormatInt(int64(settings.Conf.RpcServerConfig.Port), 10)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("grpc.Dial() failed", zap.Error(err))
		return "", err
	}
	defer conn.Close()
	c := rpc.NewConvertServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Convert(ctx, &rpc.ConvertReq{LongUrl: longUrl})
	if err != nil {
		zap.L().Error("c.Convert() failed", zap.Error(err))
		return "", err
	}
	return r.ShortUrl, nil
}

func MQSendCreatePostMessage(userID int64, post models.Post) {
	rmq := rabbitmq.NewRabbitMQPubSub("new_post")
	defer rmq.Destroy()
	//查找该用户所有粉丝
	fans, err := GetFanList(userID)
	if err != nil {
		zap.L().Error("logic.GetFanList() failed", zap.Error(err))
		return
	}
	msg := fmt.Sprintf("您关注的用户%d创建了一条新帖子", userID)
	rmq.PublishCreatePostMessage(userID, fans, msg, post)
}

func MQReceiveCreatePostMessageByMysql() {
	rmqMysql := rabbitmq.NewRabbitMQPubSub("new_post")
	defer rmqMysql.Destroy()
	rmqMysql.ConsumeCreatePostMessageByMysql()
}

func MQReceiveCreatePostMessageByEs() {
	rmqEs := rabbitmq.NewRabbitMQPubSub("new_post")
	defer rmqEs.Destroy()
	rmqEs.ConsumeCreatePostMessageByEs()
}

func SearchPost(keyWord string, page int) (data []*models.Post, err error) {
	//查询es
	posts, err := es.SearchPostIndex(keyWord, page)
	if err != nil {
		return
	}
	//查询作者信息
	for _, post := range posts {
		data = append(data, &post)
	}
	return
}
