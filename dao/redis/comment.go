package redis

import (
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

func CreateComment(commentID int64, postID int64) error {
	// 1、评论发布的时候要设置一个有效期
	pipeline := client.TxPipeline()
	pipeline.ZAdd(getRedisKey(KeyCommentTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: commentID,
	})
	// 2、评论发布的时候要初始化分数
	pipeline.ZAdd(getRedisKey(keyCommentScoreZSet), redis.Z{
		Score:  0,
		Member: commentID,
	})
	// 3、评论发布的时候要把评论id添加到帖子set里面
	cKey := getRedisKey(KeyCommentInPostSetPF + strconv.Itoa(int(postID)))
	pipeline.SAdd(cKey, commentID)
	_, err := pipeline.Exec()
	return err
}

// CommentLike 评论点赞
func CommentLike(userID, commentID string, value float64) error {
	// 1、判断投票限制
	//去redis取评论发布时间
	commentTime := client.ZScore(getRedisKey(KeyCommentTimeZSet), commentID).Val()
	if float64(time.Now().Unix())-commentTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2、更新分数
	// 先查当前用户给当前评论的投票记录
	ov := client.ZScore(getRedisKey(KeyCommentLikedZSetPF+commentID), userID).Val()
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
	pipeline.ZIncrBy(getRedisKey(keyCommentScoreZSet), op*diff*scorePerLike, commentID)
	// 3、记录投票
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyCommentLikedZSetPF+commentID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyCommentLikedZSetPF+commentID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	// 4、返回结果
	_, err := pipeline.Exec()
	return err
}
