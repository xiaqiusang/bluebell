package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//投票几种情况
/*
direction=1时，两种情况：
1、之前没有投过票，现在投赞成票 差值绝对值：1 +432
2、之前投反对票，现在投赞成票		2 +432*2
direction=0，两种情况：
1、之前投赞成票，现在取消投票		1	-432
2、之前投反对票，现在取消投票		1  +432
direction=-1，两种情况：
1、之前没投过票，现在投反对票		1	-432
2、之前投赞成票，现在投反对票		2	-432*2

投票限制：
每个帖子自发表之日起一个星期内允许用户投票
1.到期后将redis中保存的赞成和反对票数保存到mysql中
2.到期之后删除保存的keyPostVotedZset

*/

const (
	OneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票432分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeat     = errors.New("不允许重复投票")
)

func CreatePost(postid, communityID int64) error {

	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postid,
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postid,
	})

	//把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunityVotedZSetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postid)

	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票的限制
	//从redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > OneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	//2和3需要放到一个事务中操作
	//2.更新分数
	//查询之前的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedPrefix+postID), userID).Val()

	//如果这一次投票的值和之前保存的值一致
	if value == ov {
		return ErrVoteRepeat
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value)
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	//3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedPrefix+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
