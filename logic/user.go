package logic

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/dao/redis"
	"SHforum_backend/models"
	"SHforum_backend/models/response"
	"SHforum_backend/pkg/jwt"
	"SHforum_backend/pkg/snowflake"
	"SHforum_backend/rabbitmq"
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
	rcode, err := redis.GetCode(p.Email)
	if rcode != p.Code {
		return err
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

/*func SendCode(email string, code string) (err error) {
	err = util.SendEmailWithCode([]string{email}, code)
	if err != nil {
		return err
	}
	redis.SaveCode(email, code)
	return
}*/

func MQSendCodeMessage(email string, code string) {
	rmq := rabbitmq.NewRabbitMQSimple("send_code")
	rmq.PublishSendCodeMessage(email, code)
}

func MQReceiveCodeMessage() {
	rmq := rabbitmq.NewRabbitMQSimple("send_code")
	rmq.ConsumeCodeMessage()
}

func Follow(userId int64, followeduser int64) (err error) {
	return mysql.Follow(userId, followeduser)
}

func UnFollow(userId int64, followeduser int64) (err error) {
	return mysql.UnFollow(userId, followeduser)
}

func GetFanList(userId int64) (fans []int64, err error) {
	list, err := mysql.GetFanList(userId)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		fans = append(fans, v.FanUser)
	}
	return
}
