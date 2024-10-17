package logic

import (
	"SHforum_backend/internal/dao/redis"
	"SHforum_backend/internal/models"
	"go.uber.org/zap"
	"strconv"
)

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
