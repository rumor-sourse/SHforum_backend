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
	"math/rand"
	"time"
)

const (
	Charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	NumberCharset = "0123456789"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1、获取参数和参数校验
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
	// 2、业务处理
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

func randomInteger(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = NumberCharset[rand.Intn(len(NumberCharset))]
	}

	return string(b)
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = Charset[rand.Intn(len(Charset))]
	}

	return string(b)
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
	err := logic.MQSendCodeMessage(email, code)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	logic.MQReceiveCodeMessage()
	ResponseSuccess(c, nil)
}
