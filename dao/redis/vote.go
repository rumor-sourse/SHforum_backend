package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

/* 投票功能：
 投一票就加432分  86400(一天时间）/200(200张赞成票可以给你的贴子续一天）=432 来自《redis实战》

	投票的几张情况：
direction=1时，有两种情况：
	1、之前没有投过票，现在投赞成票	-->更新分数和投票记录 差值的绝对值：1 +432
    2、之前投反对票，现在改投赞成票	-->更新分数和投票记录 差值的绝对值：2 +432*2
direction=0时，有两种情况：
	1、之前投过反对票，现在要取消投票	-->更新分数和投票记录 差值的绝对值：1 +432
	2、之前投过赞成票，现在要取消投票	-->更新分数和投票记录 差值的绝对值：1 -432
direction=-1时，有两种情况：
	1、之前没有投过票，现在投反对票	-->更新分数和投票记录 差值的绝对值：1 -432
	2、之前投赞成票，现在改投反对票	-->更新分数和投票记录 差值的绝对值：2 -432*2

投票的限制：
每个贴子自发表之日起一个星期之内允许投票，超过一个星期就不允许投票了
     1、到期之后将redis中保存的赞成票数和反对票数存储到mysql中
	 2、到期之后删除那个keyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeated = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	// 1、帖子发布的时候要设置一个有效期
	pipeline := client.TxPipeline()
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 2、帖子发布的时候要初始化分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postID,
	})
	// 3、帖子发布的时候要把帖子id添加到社区set里面
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1、判断投票限制
	//去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2、更新分数
	// 先查当前用户给当前贴子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
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
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3、记录用户为该贴子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
