package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数校验
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//从c取到当发请求的用户ID
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLlogin)
		return
	}
	p.AuthorId = userID
	//2.创建帖子
	if err := logic.CreatePost(&p); err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

// 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数（从url中帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post with ivalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id取出帖子数据
	data, err := logic.GetPostbyID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostbyID() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

// 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostListHandler falied", zap.Error(err))
		return
	}
	ResponseSuccess(c, data)
}

// 升级版帖子列表接口
// 根据前端传来的参数（创建时间、得分）动态获取帖子列表
func GetPostListHandler2(c *gin.Context) {
	//Get请求参数：/api/v1/posts2?page=1&size=10&order=time
	//指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostListNew(p)
	//获取数据
	if err != nil {
		zap.L().Error("logic.GetPostListNew falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 根据社区查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	//获取数据
//	data,err:=logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("GetPostListHandler falied", zap.Error(err))
//		ResponseError(c,CodeServerBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//}
