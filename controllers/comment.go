package controllers

import (
	"SHforum_backend/logic"
	"SHforum_backend/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// GetCommentByPostIdHandler 根据帖子id获取评论列表
func GetCommentByPostIdHandler(c *gin.Context) {
	//获取参数
	pid := c.Param("id")
	postID, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		zap.L().Error("GetCommentByPostIdHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询数据
	data, err := logic.GetCommentList(postID)
	if err != nil {
		zap.L().Error("logic.GetCommentList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSuccess(c, data)
}

// CreateCommentHandler 创建评论
func CreateCommentHandler(c *gin.Context) {
	//获取参数及参数校验
	p := new(models.Comment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.UserID = userID
	//创建评论
	if err := logic.CreateComment(p); err != nil {
		zap.L().Error("logic.CreateComment() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(c *gin.Context) {
	//获取参数
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		zap.L().Error("DeleteComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//删除评论
	if err := logic.DeleteComment(commentID); err != nil {
		zap.L().Error("logic.DeleteComment() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}

// CanEditCommentHandler 判断是否有权限编辑评论
func CanEditCommentHandler(c *gin.Context) {
	//获取参数
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		zap.L().Error("CanEditCommentHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//根据id获取评论数据
	data, err := logic.GetCommentByID(commentID)
	//从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	//判断是否有权限
	if ok := data.UserID != userID; !ok {
		ResponseError(c, CodeNotPermission)
		return
	}
	//返回响应
	ResponseSuccess(c, nil)
}
