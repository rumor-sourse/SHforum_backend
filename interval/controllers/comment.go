package controllers

import (
	"SHforum_backend/interval/logic"
	"SHforum_backend/interval/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type CommentController struct {
}

// CreateCommentHandler 创建评论
func (cc *CommentController) CreateCommentHandler(c *gin.Context) {
	// 获取参数及参数校验
	p := new(models.Comment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.UserID = userID
	// 创建评论
	if err := logic.CreateComment(p); err != nil {
		zap.L().Error("CreateComment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// LikeCommentHandler 点赞评论
func (cc *CommentController) LikeCommentHandler(c *gin.Context) {
	// 获取参数
	p := new(models.ParamCommentLike)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	// 从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 点赞
	if err := logic.LikeComment(userID, p); err != nil {
		zap.L().Error("LikeComment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// DeleteCommentHandler 删除评论
func (cc *CommentController) DeleteCommentHandler(c *gin.Context) {
	// 获取参数
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		zap.L().Error("DeleteComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 删除评论
	if err := logic.DeleteComment(commentID); err != nil {
		zap.L().Error("DeleteComment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// CanEditCommentHandler 判断是否有权限编辑评论
func (cc *CommentController) CanEditCommentHandler(c *gin.Context) {
	// 获取参数
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		zap.L().Error("CanEditCommentHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id获取评论数据
	data, err := logic.GetCommentByID(commentID)
	if err != nil {
		zap.L().Error("GetCommentByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 判断是否有权限
	if ok := data.UserID != userID; !ok {
		ResponseError(c, CodeNotPermission)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}
