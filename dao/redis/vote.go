package redis

import (
	"github.com/redis/go-redis/v9"
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

func CreatePost(postID, communityID int64) error {
	// 1、帖子发布的时候要设置一个有效期
	pipeline := client.TxPipeline()
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 2、帖子发布的时候要初始化分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postID,
	})
	// 3、帖子发布的时候要把帖子id添加到社区set里面
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1、判断投票限制
	//去redis取帖子发布时间
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2、更新分数
	// 先查当前用户给当前贴子的投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
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
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3、记录用户为该贴子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}

/*论文写的
func VoteForPost(uID, pID string, v float64) error {
	// 查投票记录
	postkey := getRedisKey(KeyPostVotedZSetPF + pID)
	score := client.ZScore(ctx, postkey, uID).Val()
	//如果这一次投票的值和之前的值一样，就提示不允许重复投票
	if score == v {
		zap.L().Error("重复投票")
		return ErrorVoteRepeated
	}
	var dir float64
	if v > score {
		dir = 1
	} else {
		dir = -1
	}
	dif := math.Abs(score - v) //计算差值
	pipeline := client.TxPipeline()
	scorekey := getRedisKey(KeyPostScoreZSet)
	pipeline.ZIncrBy(ctx, scorekey, dir*dif*scorePerVote, pID)
	// 3、记录用户为该贴子投票的数据
	if v == 0 {
		pipeline.ZRem(ctx, postkey, pID)
	} else {
		pipeline.ZAdd(ctx, postkey, redis.Z{v, uID})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
*/
