package redis

import "errors"

// redis key 使用命名空间的方式，方便查询和拆分
const (
	KeyPrefix          = "SHforum_backend:"
	KeyPostTimeZSet    = "post:time"   // zset帖子以发帖时间为分数
	KeyPostScoreZSet   = "post:score"  // zset帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" // zset记录用户及投票类型;参数是post_id
	KeyCommunitySetPF  = "community:"  // set记录每个分区下帖子的id;参数是community_id

	KeyCommentTimeZSet    = "comment:time"    // zset评论以评论时间为分数
	keyCommentScoreZSet   = "comment:score"   // zset评论及投票的分数
	KeyCommentLikedZSetPF = "comment:liked:"  // zset记录用户及投票类型;参数是comment_id
	KeyCommentInPostSetPF = "comment:inpost:" // zset记录每个帖子下的评论id;参数是post_id
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票的分数
	scorePerLike     = 288 //每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeated = errors.New("不允许重复投票")
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
