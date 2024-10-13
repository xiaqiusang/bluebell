package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID))
	fmt.Println("flag")
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
