package controllers

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/logic"
	"SHforum_backend/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
)

const (
	Charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	NumberCharset = "0123456789"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型的错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 如果role字段为空就是common_user
	if len(p.Role) == 0 {
		p.Role = models.CommonUser
	}
	// 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3、返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录请求的函数
func LoginHandler(c *gin.Context) {
	// 1、获取参数和参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型的错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2、业务逻辑处理
	loginresp, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3、返回响应
	ResponseSuccess(c, gin.H{
		"userID":   fmt.Sprintf("%d", loginresp.UserID), //可能会失真
		"username": loginresp.Name,
		"token":    loginresp.Token,
	})
}

// SendCodeHandler 发送邮箱验证码
func SendCodeHandler(c *gin.Context) {
	email := c.Query("email")
	if len(email) == 0 {
		ResponseError(c, CodeEmailEmpty)
		return
	}

	code := randomInteger(6)
	/*	err := logic.SendCode(email, code)
		if err != nil {
			ResponseError(c, CodeServerBusy)
			return
		}*/
	logic.SendCode(email, code)
	ResponseSuccess(c, nil)
}

// FollowHandler 关注用户
func FollowHandler(c *gin.Context) {
	p := c.PostForm("followeduser")
	if len(p) == 0 {
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	followeduser, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	err = logic.Follow(userID, followeduser)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// UnFollowHandler 取消关注用户
func UnFollowHandler(c *gin.Context) {
	p := c.PostForm("followeduser")
	if len(p) == 0 {
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	followeduser, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	err = logic.UnFollow(userID, followeduser)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// UpdateUserByIDHandler 更新用户信息
func UpdateUserByIDHandler(c *gin.Context) {
	//获取参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		zap.L().Error("userid parse error", zap.Error(err))
		return
	}
	p := new(models.ParamUpdateUser)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("UpdateUserByID with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
	}
	//从请求中获取到当前发请求的用户的id
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
	}
	//业务处理
	fmt.Println(userID, currentUserID, p.Username)
	err = logic.UpdateUserByID(userID, currentUserID, p)
	if err != nil {
		return
	}
	ResponseSuccess(c, nil)
}

// SendMessageHandler 发送消息
func SendMessageHandler(c *gin.Context) {
	//获取参数
	userIDStr := c.Param("id")
	touserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		zap.L().Error("userid parse error", zap.Error(err))
		return
	}
	p := new(models.ParamSendMessage)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SendMessage with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
	}
	//从请求中获取到当前发请求的用户的id
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
	}
	//业务处理
	err = logic.SendMessage(currentUserID, touserID, p)
	if err != nil {
		return
	}
	ResponseSuccess(c, nil)
}
