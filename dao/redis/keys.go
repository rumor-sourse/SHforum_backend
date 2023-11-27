package redis

// redis key 使用命名空间的方式，方便查询和拆分
const (
	KeyPrefix          = "SHforum_backend:"
	KeyPostTimeZSet    = "post:time"   // zset帖子以发帖时间为分数
	KeyPostScoreZSet   = "post:score"  // zset帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" // zset记录用户及投票类型;参数是post_id
	KeyCommunitySetPF  = "community:"  // set记录每个分区下帖子的id;参数是community_id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
