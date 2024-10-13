package redis

// redis key
// redis key注意使用命名空间的方式，方便查询和拆分
const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"
	KeyPostScoreZSet   = "post:score"
	KeyPostVotedPrefix = "post:voted:"

	KeyCommunityVotedZSetPF = "community:"
)

func getRedisKey(key string) string {
	return Prefix + key
}
