package redis

import (
	"SHforum_backend/models"
	"github.com/redis/go-redis/v9"
	"math"
	"strconv"
	"time"
)

func CreateComment(commentID int64, postID int64) error {
	// 1、评论发布的时候要设置一个有效期
	pipeline := client.TxPipeline()
	pipeline.ZAdd(ctx, getRedisKey(KeyCommentTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: commentID,
	})
	// 2、评论发布的时候要初始化分数
	pipeline.ZAdd(ctx, getRedisKey(keyCommentScoreZSet), redis.Z{
		Score:  0,
		Member: commentID,
	})
	// 3、评论发布的时候要把评论id添加到帖子set里面
	cKey := getRedisKey(KeyCommentInPostSetPF + strconv.Itoa(int(postID)))
	pipeline.SAdd(ctx, cKey, commentID)
	_, err := pipeline.Exec(ctx)
	return err
}

// CommentLike 评论点赞
func CommentLike(userID, commentID string, value float64) error {
	// 1、判断投票限制
	//去redis取评论发布时间
	commentTime := client.ZScore(ctx, getRedisKey(KeyCommentTimeZSet), commentID).Val()
	if float64(time.Now().Unix())-commentTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2、更新分数
	// 先查当前用户给当前评论的投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyCommentLikedZSetPF+commentID), userID).Val()
	//如果这一次投票的值和之前的值一样，就提示不允许重复投票
	if value == ov {
		return ErrorVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(keyCommentScoreZSet), op*diff*scorePerLike, commentID)
	// 3、记录投票
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyCommentLikedZSetPF+commentID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyCommentLikedZSetPF+commentID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	// 4、返回结果
	_, err := pipeline.Exec(ctx)
	return err
}

// GetCommentLikeCount 获取评论点赞数
func GetCommentLikeCount(commentID string) (count int64, err error) {
	pipeline := client.Pipeline()
	key := getRedisKey(KeyCommentLikedZSetPF + commentID)
	cmd := pipeline.ZCount(ctx, key, "1", "1")
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return 0, err
	}
	count, err = cmd.Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetPostCommentIDsInOrder(p *models.ParamCommentList) ([]string, error) {
	//1、根据用户请求中携带的order参数确定要查询的redis key
	orderkey := getRedisKey(KeyCommentTimeZSet)
	if p.Order == models.OrderScore {
		orderkey = getRedisKey(keyCommentScoreZSet)
	}

	//贴子的key
	pKey := getRedisKey(KeyCommentInPostSetPF + strconv.Itoa(int(p.PostID)))
	// 利用缓存key减少zinterstore的执行次数
	key := orderkey + strconv.Itoa(int(p.PostID))
	if client.Exists(ctx, key).Val() < 1 {
		//不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{pKey, orderkey},
			Aggregate: "MAX",
		})
		pipeline.Expire(ctx, key, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	//存在的话直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}

// UpdateHotComment 缓存存储热评信息
func UpdateHotComment(postID string, hotcomment map[string]string) error {
	pkey := getRedisKey(KeyHotCommentHashPF + postID)
	_, err := client.HSet(ctx, pkey, hotcomment).Result()
	if err != nil {
		return err
	}
	//过期时间设置为半小时
	_, err = client.Expire(ctx, pkey, 30*time.Minute).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetHotComment 获取热评信息
func GetHotComment(postID string) (data map[string]string, err error) {
	key := getRedisKey(KeyHotCommentHashPF + postID)
	data, err = client.HGetAll(ctx, key).Result()
	if len(data) == 0 {
		return nil, Nil
	} else if err != nil {
		return nil, err
	}
	return data, nil
}
