package logic

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/models"
	"SHforum_backend/models/response"
	"SHforum_backend/pkg/jwt"
	"SHforum_backend/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//1、判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错
		return err
	}
	//2、生成UID
	userId := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userId,
		Username: p.Username,
		Password: p.Password,
		Email:    p.Email,
	}
	//3、保存用户信息
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (userresp *response.UserResponse, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的指针，拿到userID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	userresp = &response.UserResponse{
		UserID: user.UserID,
		Name:   user.Username,
		Token:  token,
	}
	return
}
