package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成post id
	p.PostID = int64(snowflake.GenID())
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.PostID, p.CommunityID)
	return
}

func GetPostbyID(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合接口需要的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostbyID() failed", zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserbyid(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetPostbyID() failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
		return
	}
	//根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunitybyID() failed", zap.Int64("Community_id", post.CommunityID), zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := mysql.GetUserbyid(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetPostbyID() failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunitybyID() failed", zap.Int64("Community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2、从redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data")
		return
	}
	// 3、根据id从mysql数据库查询帖子详细信息
	//返回的数据要按照给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//提取查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将贴子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		user, err := mysql.GetUserbyid(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetPostbyID() failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunitybyID() failed", zap.Int64("Community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2、从redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data")
		return
	}
	// 3、根据id从mysql数据库查询帖子详细信息
	//返回的数据要按照给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//提取查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将贴子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		user, err := mysql.GetUserbyid(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetPostbyID() failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailbyID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunitybyID() failed", zap.Int64("Community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

// 将两个查询接口合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//根据请求参数不同，执行不同的逻辑
	if p.CommunityID == 0 {
		//查所有
		data, err = GetPostList2(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
	}
	return
}
